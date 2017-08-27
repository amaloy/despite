package main

import (
	"bufio"
	"fmt"
	"net"
)

func playerLoop(conn net.Conn, clientId int) {
	playerWriteBytes(conn, motd)
	playerWrite(conn, "\nDragonroar!\nV0026\n")

	reader := bufio.NewReader(conn)
	for {
		incoming, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if incoming[0] == '"' {
			chanBroadcast <- fmt.Sprintf("(Player %d: %s", clientId, incoming[1:])
		} else {
			playerWrite(conn, "\n(That just won't do.\n")
		}
	}

	chanDeadConns <- conn
}
