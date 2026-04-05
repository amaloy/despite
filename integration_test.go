package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"testing"
	"time"
)

// sendCommand sends a request over the TCP connection and waits briefly for processing.
func sendCommand(t *testing.T, conn net.Conn, request string) {
	_, err := conn.Write([]byte(request))
	if err != nil {
		t.Fatalf("Failed to send command: %v", err)
	}
	t.Logf("OK ->: %s", request)
	time.Sleep(100 * time.Millisecond)
}

// assertMatch checks if msg matches pattern, using regex if isRegex is true, otherwise exact match.
func assertMatch(t *testing.T, pattern string, msg string, isRegex bool) bool {
	if isRegex {
		matched, err := regexp.MatchString(pattern, msg)
		if err != nil {
			t.Fatalf("Invalid regex %q: %v", pattern, err)
		}
		return matched
	} else {
		return msg == pattern
	}
}

// assertResponse asserts that all strings in expected appear in the reader in order.
// If expectedAnyOrder is non-nil, unexpected values that match expectedAnyOrder are allowed, each at most once.
// If isRegex is true, patterns are treated as regex; if false, literal matches.
func assertResponse(t *testing.T, reader *bufio.Reader, expected []string, expectedAnyOrder []string, isRegex bool) {
	expectedAnyOrderCopy := make([]string, len(expectedAnyOrder))
	copy(expectedAnyOrderCopy, expectedAnyOrder)
	for _, exp := range expected {
		msg, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}
		msg = strings.TrimSpace(msg)
		if assertMatch(t, exp, msg, isRegex) {
			t.Logf("OK <-: %s", msg)
		} else {
			// check if in expectedAnyOrder
			found := false
			if expectedAnyOrderCopy != nil {
				for j, anyExp := range expectedAnyOrderCopy {
					if assertMatch(t, anyExp, msg, isRegex) {
						t.Logf("OK <- (any order): %s", msg)
						expectedAnyOrderCopy = append(expectedAnyOrderCopy[:j], expectedAnyOrderCopy[j+1:]...)
						found = true
						break
					}
				}
			}
			if !found {
				t.Fatalf("Expected %v or any order, got %q", exp, msg)
			}
		}
	}
	// Continue reading for any remaining expectedAnyOrder
	for len(expectedAnyOrderCopy) > 0 {
		msg, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("Failed to read response for expectedAnyOrder: %v", err)
		}
		msg = strings.TrimSpace(msg)
		found := false
		for j, anyExp := range expectedAnyOrderCopy {
			if assertMatch(t, anyExp, msg, isRegex) {
				t.Logf("OK <- (any order): %s", msg)
				expectedAnyOrderCopy = append(expectedAnyOrderCopy[:j], expectedAnyOrderCopy[j+1:]...)
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Unexpected message while waiting for expectedAnyOrder: %q", msg)
		}
	}
}

// sendCommandAndAssertResponse sends a request over the TCP connection, reads the next line of response,
// and asserts that all strings in expected appear in the response message in the order given.
func sendCommandAndAssertResponse(t *testing.T, conn net.Conn, reader *bufio.Reader, request string, expected []string) {
	sendCommand(t, conn, request)
	assertResponse(t, reader, expected, nil, false)
}

// TestIntegration_SinglePlayer tests the service by launching it and connecting a single player using a TCP connection
func TestIntegration_SinglePlayer(t *testing.T) {
	// Set an overall timeout for the test
	timeout := time.After(15 * time.Second)
	done := make(chan bool)
	errorChan := make(chan error)

	go func() {
		defer func() { done <- true }()

		// Start the server on a dynamically allocated port
		server, allPlayers, newConnections := startServer(":0", "Despite Server Test MOTD")
		defer server.Close()

		// Get the dynamically assigned port
		serverAddr := server.Addr().String()
		log.Printf("Server started on %s", serverAddr)

		// Handle new connections and broadcasts
		go runServerLoop(allPlayers, newConnections)

		// Give the server some time to start
		time.Sleep(200 * time.Millisecond)

		// Connect to the server as a client using the dynamically assigned port
		conn, err := net.Dial("tcp", serverAddr)
		if err != nil {
			errorChan <- fmt.Errorf("Failed to connect to server: %v", err)
			return
		}
		defer conn.Close()

		// Create a reader for server responses
		reader := bufio.NewReader(conn)

		assertResponse(t, reader, []string{motd, "", "Dragonroar!", "V0026"}, nil, false)

		username := "testplayer"
		password := "testpass"
		connectCmd := fmt.Sprintf("connect %s %s\n", username, password)
		sendCommandAndAssertResponse(t, conn, reader, connectCmd, []string{"cs"})

		colorCmd := "color #FFF!\n"
		sendCommand(t, conn, colorCmd)

		enteredMsg := fmt.Sprintf("\\(%s has entered Despite.", username)

		descCmd := "desc Test description\n"
		sendCommand(t, conn, descCmd)
		assertResponse(t, reader, []string{"&", "]lev01", "^\\@.{2}$", "^\\<.{7}$", "\\="}, []string{enteredMsg, "^\\<.{7}$"}, true)

		// Test move command
		moveCmd := "m 1\n"
		sendCommand(t, conn, moveCmd)
		// Assert successful move: ~ @<coords> <place><old_coords>  =
		assertResponse(t, reader, []string{"~", "^\\@.{2}$", "^\\<.{9}$", "="}, nil, true)

		// Test rotate left
		rotateLeftCmd := "<\n"
		sendCommand(t, conn, rotateLeftCmd)
		assertResponse(t, reader, []string{"^\\<.{7}$"}, nil, true)

		// Test rotate right
		rotateRightCmd := ">\n"
		sendCommand(t, conn, rotateRightCmd)
		assertResponse(t, reader, []string{"^\\<.{7}$"}, nil, true)

		// Test chat command
		chatCmd := "\"Hello world\n"
		sendCommand(t, conn, chatCmd)
		expectedChat := fmt.Sprintf("(%s: Hello world", username)
		assertResponse(t, reader, []string{expectedChat}, nil, false)

		conn.Close()
		t.Log("Closed client connection")
		// Wait for a short time to ensure server processes the disconnection
		time.Sleep(100 * time.Millisecond)
	}()

	select {
	case <-timeout:
		t.Fatal("Test timed out after 10 seconds")
	case err := <-errorChan:
		t.Fatal(err)
	case <-done:
		// Test completed normally
		t.Log("Test completed successfully")
	}
}
