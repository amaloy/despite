package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
)

const serverName string = "Despite"

var chanCleanDisconns = make(chan *player)
var chanBroadcast = make(chan string)
var motd string
var mainMap *dsmap

func main() {

	if motdBytes, err := ioutil.ReadFile("motd.txt"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		motd = string(motdBytes)
	}

	mainMap = buildMainMap()

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

			log.Printf("Accepted new connection, #%d", clientCount)

			p := new(player)
			p.connID = clientCount
			p.conn = conn
			p.reader = bufio.NewReader(conn)
			p.writer = bufio.NewWriter(conn)

			allPlayers[p.connID] = p
			clientCount++

			// Spawn independant player exec
			go playerExec(p)

		case message := <-chanBroadcast:
			// Broadcast to all players
			for _, p := range allPlayers {
				go playerWrite(p, message)
			}

		case p := <-chanCleanDisconns:
			log.Printf("%s disconnected", p.name)
			delete(allPlayers, p.connID)
			p.conn.Close()
		}
	}
}

func toDSChar(i int) rune {
	return (rune)(i + 32)
}

func buildMainMap() (m *dsmap) {
	m = new(dsmap)
	m.name = "lev01"
	m.width = standardMapWidth
	m.height = standardMapHeight
	m.xstart = 26
	m.ystart = 41
	m.players = make(map[int]*player)
	return
}
