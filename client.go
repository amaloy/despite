package main

import (
	"bufio"
	"fmt"
	"net"
)

func clientLoop(conn net.Conn, clientId int) {
	clientWriteBytes(conn, motd)
	clientWrite(conn, "\nDragonroar!\nV0026\n")

	reader := bufio.NewReader(conn)
	for {
		incoming, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if incoming[0] == '"' {
			chanBroadcast <- fmt.Sprintf("Client %d: %s", clientId, incoming[1:])
		} else {
			clientWrite(conn, "\n(That just won't do.)\n")
		}
	}

	chanDeadConns <- conn
}
