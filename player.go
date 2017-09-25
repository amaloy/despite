package main

import (
	"bufio"
	"net"
)

type player struct {
	connID   int
	conn     net.Conn
	reader   *bufio.Reader
	writer   *bufio.Writer
	name     string
	color    string
	pstring  string
	desc     string
	readLine string
}

func playerExec(p *player) {
	var err error
	err = playerLoginLoop(p)
	if err != nil {
		playerLogOut(p)
		return
	}
	playerMainLoop(p)
	playerLogOut(p)
}

func playerLogOut(p *player) {
	chanCleanDisconns <- p
}

func playerWrite(p *player, message string) (err error) {
	_, err = p.writer.WriteString(message)
	if err != nil {
		return err
	}
	_, err = p.writer.WriteRune('\n')
	if err != nil {
		return err
	}
	err = p.writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func playerReadLine(p *player) (string, error) {
	return p.reader.ReadString('\n')
}
