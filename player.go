package main

import (
	"bufio"
	"fmt"
)

type player struct {
	connID int
	reader *bufio.Reader
	writer *bufio.Writer
	name   string
}

func playerLoop(p *player) {
	playerWrite(p, motd)
	playerWrite(p, "\nDragonroar!\nV0026\n")

	var line string
	var err error
	for {
		line, err = playerReadLine(p)
		if err != nil {
			break
		}
		if line[0] == '"' {
			chanBroadcast <- fmt.Sprintf("(%s: %s", p.name, line[1:])
		} else {
			playerWrite(p, "\n(That just won't do.\n")
		}
	}

	chanDisconnPlayers <- p
}

func playerWrite(p *player, message string) {
	if _, err := p.writer.WriteString(message); err != nil {
		chanDisconnPlayers <- p
	} else {
		p.writer.Flush()
	}
}

func playerReadLine(p *player) (string, error) {
	return p.reader.ReadString('\n')
}
