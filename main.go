package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
)

var chanDeadConns = make(chan net.Conn)
var chanBroadcast = make(chan string)
var motd []byte

func main() {

	var err error

	motd, err = ioutil.ReadFile("motd.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Connection count increment (not needed once there are names)
	clientCount := 0

	// All people who are connected; a map wherein
	// the keys are net.Conn objects and the values
	// are client "ids", an integer.
	//
	allPlayers := make(map[net.Conn]int)

	newConnections := make(chan net.Conn)

	server, err := net.Listen("tcp", ":7734")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go func() {
		for {
			// Accept new connections
			conn, err := server.Accept()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			newConnections <- conn
		}
	}()

	for {

		select {

		case conn := <-newConnections:

			log.Printf("Accepted new player, #%d", clientCount)

			allPlayers[conn] = clientCount
			clientCount++

			// Spawn independant player loop
			go playerLoop(conn, allPlayers[conn])

		case message := <-chanBroadcast:
			// Broadcast to all player
			for conn := range allPlayers {
				go playerWrite(conn, message)
			}
			log.Printf("New message: %s", message)
			log.Printf("broadcast to %d players", len(allPlayers))

		case conn := <-chanDeadConns:
			log.Printf("Player %d disconnected", allPlayers[conn])
			delete(allPlayers, conn)
		}
	}
}

func playerWrite(conn net.Conn, message string) {
	playerWriteBytes(conn, []byte(message))
}

func playerWriteBytes(conn net.Conn, message []byte) {
	_, err := conn.Write(message)
	if err != nil {
		chanDeadConns <- conn
	}
}
