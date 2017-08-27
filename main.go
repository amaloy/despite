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
	allPlayers := make(map[net.Conn]*player)
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

			p := new(player)
			p.conn = conn
			p.name = fmt.Sprintf("Player %v", clientCount)

			allPlayers[conn] = p
			clientCount++

			// Spawn independant player loop
			go playerLoop(p)

		case message := <-chanBroadcast:
			// Broadcast to all player
			for _, p := range allPlayers {
				go playerWrite(p, message)
			}
			log.Printf("New message: %s", message)
			log.Printf("broadcast to %d players", len(allPlayers))

		case conn := <-chanDeadConns:
			log.Printf("%s disconnected", allPlayers[conn].name)
			delete(allPlayers, conn)
		}
	}
}

func playerWrite(p *player, message string) {
	playerWriteBytes(p, []byte(message))
}

func playerWriteBytes(p *player, message []byte) {
	_, err := p.conn.Write(message)
	if err != nil {
		chanDeadConns <- p.conn
	}
}
