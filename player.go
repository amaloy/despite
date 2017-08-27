package main

import (
	"bufio"
	"fmt"
	"net"
)

type player struct {
	conn net.Conn
	name string
}

func playerLoop(p *player) {
	playerWriteBytes(p, motd)
	playerWrite(p, "\nDragonroar!\nV0026\n")

	reader := bufio.NewReader(p.conn)
	for {
		incoming, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if incoming[0] == '"' {
			chanBroadcast <- fmt.Sprintf("(%s: %s", p.name, incoming[1:])
		} else {
			playerWrite(p, "\n(That just won't do.\n")
		}
	}

	chanDeadConns <- p.conn
}
