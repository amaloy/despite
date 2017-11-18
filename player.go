package main

import (
	"bufio"
	"net"
)

type player struct {
	connID          int
	conn            net.Conn
	reader          *bufio.Reader
	writer          *bufio.Writer
	name            string
	color           string
	pstring         string
	desc            string
	facing          int
	facingShapeBase rune
	visibleShape    rune
	shapeMoveCycle  int
	readLine        string
	mapContext      *playerMapContext
}

type playerMapContext struct {
	currMap  *dsmap
	currX    int
	currY    int
	dsCoords string
}

var longShapeStart = [][]int{
	{2, 2, 6, 10, 10, 6, 10, 14, 14},
	{2, 2, 7, 12, 12, 7, 12, 17, 17},
	{2, 2, 5, 8, 8, 5, 8, 11, 11}}

var moveCycleLoop = [4]int{-1, 0, 1, 0}

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
	p.mapContext.currMap.removePlayer(p)
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

func (p *player) setShapeStanding() {
	p.facingShapeBase = toDSChar(longShapeStart[1][p.facing-1])
	p.visibleShape = p.facingShapeBase
}

func (p *player) setShapeCycleMove() {
	p.facingShapeBase = toDSChar(longShapeStart[1][p.facing-1])
	p.visibleShape = p.facingShapeBase + rune(moveCycleLoop[p.shapeMoveCycle])
	p.shapeMoveCycle++
	if p.shapeMoveCycle == 4 {
		p.shapeMoveCycle = 0
	}
}
