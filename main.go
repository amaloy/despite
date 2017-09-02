package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
)

var chanCleanDisconns = make(chan *player)
var chanBroadcast = make(chan string)
var motd string

func main() {

	if motdBytes, err := ioutil.ReadFile("motd.txt"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		motd = string(motdBytes)
	}

	// Connection count increment (not needed once there are names)
	clientCount := 0
	allPlayers := make(map[int]*player)
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
			p.connID = clientCount
			p.reader = bufio.NewReader(conn)
			p.writer = bufio.NewWriter(conn)
			p.name = fmt.Sprintf("Player %v", clientCount)

			allPlayers[p.connID] = p
			clientCount++

			// Spawn independant player exec
			go playerExec(p)

		case message := <-chanBroadcast:
			// Broadcast to all player
			for _, p := range allPlayers {
				go playerWrite(p, message)
			}
			log.Printf("New message: %s", message)
			log.Printf("broadcast to %d players", len(allPlayers))

		case p := <-chanCleanDisconns:
			log.Printf("%s disconnected", p.name)
			delete(allPlayers, p.connID)
		}
	}
}
