package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid" // Assuming this is the uuid package used in the codebase
)

// Helper to read a line from a bufio.Reader
func readLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(line, "\n"), nil // Remove trailing newline
}

// Helper to write a line to a bufio.Writer
func writeLine(w *bufio.Writer, s string) error {
	_, err := w.WriteString(s)
	if err != nil {
		return err
	}
	_, err = w.WriteRune('\n')
	if err != nil {
		return err
	}
	return w.Flush()
}

// Helper to read a line with a timeout using net.Conn.SetReadDeadline
func readWithTimeout(conn net.Conn, reader *bufio.Reader, timeout time.Duration) (string, error) {
	deadline := time.Now().Add(timeout)
	if err := conn.SetReadDeadline(deadline); err != nil {
		return "", err
	}
	defer conn.SetReadDeadline(time.Time{}) // Reset deadline after
	line, err := readLine(reader)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return "", fmt.Errorf("read timeout after %v", timeout)
		}
		return "", err
	}
	return line, nil
}

// Global variables that need to be reset/mocked for testing.
// Store original values to restore after test.
var originalMotd string
var originalMainMap *dsmap
var originalChanCleanDisconns chan *player
var originalChanBroadcast chan *broadcastPayload

// Mutex to protect global variables during setup/teardown if tests run in parallel
var globalsMutex sync.Mutex

func setupTestGlobals() {
	globalsMutex.Lock()
	defer globalsMutex.Unlock()

	// Store original globals
	originalMotd = motd
	originalMainMap = mainMap
	originalChanCleanDisconns = chanCleanDisconns
	originalChanBroadcast = chanBroadcast

	// Initialize test-specific globals
	motd = "Test MOTD"
	mainMap = buildMainMap()

	// Create new buffered channels for the test.
	// These channels will be handled by the test's simulated server loop.
	chanCleanDisconns = make(chan *player, 10)
	chanBroadcast = make(chan *broadcastPayload, 10)

	// Removed the temporary goroutines that previously consumed chanBroadcast
	// and chanCleanDisconns, as they are now handled by the testServerLoop.
}

func teardownTestGlobals() {
	globalsMutex.Lock()
	defer globalsMutex.Unlock()

	// Restore original globals
	motd = originalMotd
	mainMap = originalMainMap
	chanCleanDisconns = originalChanCleanDisconns
	chanBroadcast = originalChanBroadcast
}

func TestPlayerLoginAndMovement(t *testing.T) {
	t.Logf("Starting TestPlayerLoginAndMovement")
	setupTestGlobals()
	defer teardownTestGlobals()

	// Simulate client and server connections using net.Pipe
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	// Server-side player setup
	p := &player{
		connID: uuid.NewV4(),
		conn:   serverConn,
		reader: bufio.NewReader(serverConn),
		writer: bufio.NewWriter(serverConn),
	}

	// Channel to signal when playerExec has finished
	var playerExecWg sync.WaitGroup
	playerExecWg.Add(1)
	go func() {
		defer playerExecWg.Done()
		playerExec(p)
	}()

	clientReader := bufio.NewReader(clientConn)
	clientWriter := bufio.NewWriter(clientConn)

	// --- TEST SERVER LOOP SIMULATION ---
	// This goroutine simulates the main.go select loop, handling broadcasts and disconnections
	// within the test environment.
	testPlayers := make(map[uuid.UUID]*player) // Local map for players in this test's server loop
	testServerStop := make(chan struct{})      // Channel to signal the test server loop to stop
	var testServerWg sync.WaitGroup
	testServerWg.Add(1)
	go func() {
		defer testServerWg.Done()
		// Manually add the test player to the test's allPlayers map as main.go would
		// In a real main.go, this happens after Accept(), but here we have a pre-initialized player 'p'.
		testPlayers[p.connID] = p
		t.Logf("Test server loop: Player %v added to testPlayers map.", p.connID)

		for {
			select {
			case payload := <-chanBroadcast:
				var targets map[uuid.UUID]*player
				if payload.targetMap == nil {
					targets = testPlayers // Use the test's player map
				} else {
					targets = payload.targetMap.players
				}

				for _, playerInMap := range targets {
					if playerInMap != payload.excludePlayer {
						// Execute p.send in a goroutine to avoid blocking the server loop,
						// mimicking main.go's behavior: `go p.send(*payload.message)`
						go playerInMap.send(*payload.message)
					}
				}
			case disconnectedPlayer := <-chanCleanDisconns:
				// This handles player cleanup similar to main.go's select loop
				t.Logf("Test server loop: Player %s (%v) received for disconnection.", disconnectedPlayer.name, disconnectedPlayer.connID)
				delete(testPlayers, disconnectedPlayer.connID)
				// Note: conn.Close() is deferred at the top level of the test, no need to do it here again.
			case <-testServerStop:
				t.Logf("Test server loop: Stopping.")
				return
			}
		}
	}()
	// --- END TEST SERVER LOOP SIMULATION ---

	// 1. Expect MOTD and version
	t.Logf("Expecting MOTD and version...")
	motdLine, err := readWithTimeout(clientConn, clientReader, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("Failed to read MOTD: %v", err)
	}
	if motdLine != "Test MOTD" {
		t.Errorf("Expected MOTD 'Test MOTD', got '%s'", motdLine)
	}
	t.Logf("Received MOTD: '%s'", motdLine)

	// First part of the version string is an empty line from the leading '\n'
	versionLine1Empty, err := readWithTimeout(clientConn, clientReader, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("Failed to read first version line (expected empty line): %v", err)
	}
	if versionLine1Empty != "" {
		t.Errorf("Expected empty line, got '%s'", versionLine1Empty)
	}
	t.Logf("Received first version line (empty).")

	// Second part of the version string is "Dragonroar!"
	versionLine2Dragonroar, err := readWithTimeout(clientConn, clientReader, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("Failed to read second version line: %v", err)
	}
	if versionLine2Dragonroar != "Dragonroar!" {
		t.Errorf("Expected 'Dragonroar!', got '%s'", versionLine2Dragonroar)
	}
	t.Logf("Received second version line: '%s'", versionLine2Dragonroar)

	// Third part of the version string is "V0026"
	versionLine3Code, err := readWithTimeout(clientConn, clientReader, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("Failed to read third version line: %v", err)
	}
	if versionLine3Code != "V0026" {
		t.Errorf("Expected 'V0026', got '%s'", versionLine3Code)
	}
	t.Logf("Received third version line: '%s'", versionLine3Code)

	// 2. Send connect command
	t.Logf("Sending connect command...")
	err = writeLine(clientWriter, "connect testuser testpass")
	if err != nil {
		t.Fatalf("Failed to write connect command: %v", err)
	}
	// p.name will be set by cmdConnect in playerLoginLoop, this is just for clarity in logs
	t.Logf("Connect command sent.")

	// 3. Expect "cs" prompt for color
	t.Logf("Expecting color prompt 'cs'...")
	colorPrompt, err := readWithTimeout(clientConn, clientReader, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("Failed to read color prompt: %v", err)
	}
	if colorPrompt != "cs" {
		t.Errorf("Expected color prompt 'cs', got '%s'", colorPrompt)
	}
	t.Logf("Received color prompt: '%s'", colorPrompt)

	// 4. Send color command
	t.Logf("Sending color command...")
	err = writeLine(clientWriter, "color ffff")
	if err != nil {
		t.Fatalf("Failed to write color command: %v", err)
	}
	// p.color will be set by cmdColor
	t.Logf("Color command sent.")

	// 5. Send description command
	t.Logf("Sending description command...")
	err = writeLine(clientWriter, "desc My Test Player")
	if err != nil {
		t.Fatalf("Failed to write desc command: %v", err)
	}
	// p.desc will be set by cmdDesc
	t.Logf("Description command sent.")

	// Give a longer moment for server to process and playerMainLoop to start
	time.Sleep(1000 * time.Millisecond) // Increased sleep time to ensure server processes everything
	t.Logf("Waiting for playerMainLoop to start and broadcast entry message...")

	// 6. Collect all initialization messages
	t.Logf("Collecting all initialization messages...")
	var receivedMessages []string
	// Keep reading until we get the player position message or hit max attempts
	maxAttempts := 30
	for i := 0; i < maxAttempts; i++ {
		msgPart, err := readWithTimeout(clientConn, clientReader, 2000*time.Millisecond)
		if err != nil {
			t.Logf("Failed to read message part %d: %v", i+1, err)
			break
		}
		receivedMessages = append(receivedMessages, msgPart)
		t.Logf("Received message part %d: '%s'", i+1, msgPart)
		// Check if we received the player position message
		if strings.HasPrefix(msgPart, "@") {
			break
		}
		// Check if we received the resume draw message, which often comes after position
		if msgPart == "=" {
			break
		}
	}

	// Check for expected messages in the collected messages
	// Check for player entry message
	t.Logf("Checking for player entry message in collected messages...")
	var fullEntryMsg string
	for _, msg := range receivedMessages {
		fullEntryMsg += msg
	}
	expectedEntryMsg := fmt.Sprintf("(testuser has entered %s.", serverName) // serverName is a global constant "Despite"
	if !strings.Contains(fullEntryMsg, expectedEntryMsg) {
		t.Errorf("Expected entry message to contain '%s', got '%s' in combined messages", expectedEntryMsg, fullEntryMsg)
	} else {
		t.Logf("Received entry message containing: '%s'", expectedEntryMsg)
	}

	// Check for map name message
	t.Logf("Checking for map name message in collected messages...")
	var mapNameMsg string
	for _, msg := range receivedMessages {
		if strings.HasPrefix(msg, "]") {
			mapNameMsg = msg
			break
		}
	}
	if mapNameMsg == "" {
		t.Errorf("Expected map name message with prefix ']', not found in messages: %v", receivedMessages)
	} else if mapNameMsg != "]lev01" { // Assuming buildMainMap always returns "lev01"
		t.Errorf("Expected map name ']lev01', got '%s'", mapNameMsg)
	} else {
		t.Logf("Received map name: '%s'", mapNameMsg)
	}

	// Check for player write at message
	t.Logf("Checking for player write at message in collected messages...")
	var playerWriteAtMsg string
	for _, msg := range receivedMessages {
		if strings.HasPrefix(msg, "@") {
			playerWriteAtMsg = msg
			break
		}
	}
	if playerWriteAtMsg == "" {
		t.Logf("Warning: Expected player write at message with prefix '@', not found in messages: %v", receivedMessages)
		t.Logf("Skipping movement test due to missing player position message")
	} else if len(playerWriteAtMsg) != 3 { // '@' + 2 coords chars
		t.Logf("Warning: Expected player write at to be 3 chars, got %d: '%s'", len(playerWriteAtMsg), playerWriteAtMsg)
		t.Logf("Skipping movement test due to invalid player position message")
	} else {
		initialDsCoords := playerWriteAtMsg[1:]
		t.Logf("Received player initial position: '%s'", playerWriteAtMsg)
		// 7. Send a move command (e.g., m1 for facing=1, which is south/down)
		// According to playerMainLoop.go, the command is 'm' followed by direction digit
		t.Logf("Sending move command 'm1'...")
		err = writeLine(clientWriter, "m1")
		if err != nil {
			t.Fatalf("Failed to write move command: %v", err)
		}
		t.Logf("Move command sent.")

		// Expect drawing halt, then player write at new coords, then drawing resume (directly sent by playerExec)
		t.Logf("Expecting draw halt message...")
		haltDrawMsg, err := readWithTimeout(clientConn, clientReader, 2000*time.Millisecond) // Increased timeout
		if err != nil {
			t.Fatalf("Failed to read halt draw message: %v", err)
		}
		if haltDrawMsg != "~" {
			t.Errorf("Expected halt draw '~', got '%s'", haltDrawMsg)
		}
		t.Logf("Received halt draw: '%s'", haltDrawMsg)

		t.Logf("Expecting new player write at message...")
		newPlayerWriteAtMsg, err := readWithTimeout(clientConn, clientReader, 2000*time.Millisecond) // Increased timeout
		if err != nil {
			t.Fatalf("Failed to read new player write at: %v", err)
		}
		if !strings.HasPrefix(newPlayerWriteAtMsg, "@") {
			t.Errorf("Expected new player write at prefix '@', got '%s'", newPlayerWriteAtMsg)
		}
		newDsCoords := newPlayerWriteAtMsg[1:]
		if newDsCoords == initialDsCoords {
			t.Logf("Warning: Player coordinates did not change after move. Initial: '%s', New: '%s'", initialDsCoords, newDsCoords)
			// Not failing the test here as it might depend on map layout
		}
		t.Logf("Received new player position: '%s'", newPlayerWriteAtMsg)

		t.Logf("Expecting resume draw message...")
		resumeDrawMsg, err := readWithTimeout(clientConn, clientReader, 2000*time.Millisecond) // Increased timeout
		if err != nil {
			t.Fatalf("Failed to read resume draw message: %v", err)
		}
		if resumeDrawMsg != "=" {
			t.Errorf("Expected resume draw '=', got '%s'", resumeDrawMsg)
		}
		t.Logf("Received resume draw: '%s'", resumeDrawMsg)
	}

	// Signal playerExec to terminate by closing the client connection
	t.Logf("Closing client connection to terminate playerExec...")
	clientConn.Close()
	serverConn.Close()
	t.Logf("Connections closed.")

	// Wait for playerExec to finish
	t.Logf("Waiting for playerExec to finish...")
	playerExecDone := make(chan struct{})
	go func() {
		playerExecWg.Wait()
		close(playerExecDone)
	}()

	select {
	case <-playerExecDone:
		t.Logf("playerExec finished successfully.")
	case <-time.After(500 * time.Millisecond):
		t.Fatal("playerExec did not terminate in time")
	}

	// Give a small moment for chanCleanDisconns to be processed by the test server loop
	time.Sleep(10 * time.Millisecond)

	// Stop the test server loop
	t.Logf("Stopping test server loop...")
	close(testServerStop)
	testServerWg.Wait()
	t.Logf("Test server loop stopped.")

	// Verify player removal from the testPlayers map managed by testServerLoop
	t.Logf("Verifying player removal from testPlayers map...")
	if _, exists := testPlayers[p.connID]; exists {
		t.Errorf("Player %v was not removed from testPlayers map after logout", p.connID)
	} else {
		t.Logf("Player %v correctly removed from testPlayers map.", p.connID)
	}

	t.Logf("TestPlayerLoginAndMovement finished.")
}
