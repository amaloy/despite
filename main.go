package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/satori/go.uuid"
)

const serverName string = "Despite"

type broadcastPayload struct {
	message       string
	targetMap     *dsmap
	excludePlayer *player
}

var chanCleanDisconns = make(chan *player)
var chanBroadcast = make(chan broadcastPayload)
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

	allPlayers := make(map[uuid.UUID]*player)
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

			log.Printf("Accepted new connection, %v", conn.RemoteAddr())

			p := new(player)
			p.connID = uuid.NewV4()
			p.conn = conn
			p.reader = bufio.NewReader(conn)
			p.writer = bufio.NewWriter(conn)

			allPlayers[p.connID] = p

			// Spawn independant player exec
			go playerExec(p)

		case payload := <-chanBroadcast:
			var targets map[uuid.UUID]*player
			if payload.targetMap == nil {
				targets = allPlayers
			} else {
				targets = payload.targetMap.players
			}

			for _, p := range targets {
				if p != payload.excludePlayer {
					go p.send(payload.message)
				}
			}

		case p := <-chanCleanDisconns:
			log.Printf("%s (%v) disconnected", p.name, p.conn.RemoteAddr())
			delete(allPlayers, p.connID)
			p.conn.Close()
		}
	}
}

func broadcast(message string, targetMap *dsmap, excludePlayer *player) {
	chanBroadcast <- broadcastPayload{message, targetMap, excludePlayer}
}

func broadcastAll(message string) {
	broadcast(message, nil, nil)
}

func broadcastMap(message string, p *player) {
	broadcast(message, p.mapContext.currMap, nil)
}

func broadcastMapExclude(message string, p *player) {
	broadcast(message, p.mapContext.currMap, p)
}

func toDSChar(i int) rune {
	return (rune)(i + 32)
}

func buildMainMap() (m *dsmap) {
	m = new(dsmap)
	m.name = "lev01"
	m.width = standardMapWidth
	m.height = standardMapHeight
	m.tiles = make([][]*dsmapTile, m.width)
	for x := range m.tiles {
		row := make([]*dsmapTile, m.height)
		for y := range row {
			row[y] = new(dsmapTile)
		}
		m.tiles[x] = row
	}
	m.xstart = 26
	m.ystart = 41
	m.players = make(map[uuid.UUID]*player)
	return
}
