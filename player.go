package main

import (
	"bufio"
	"net"

	"github.com/satori/go.uuid"
)

type player struct {
	connID          uuid.UUID
	conn            net.Conn
	reader          *bufio.Reader
	writer          *bufio.Writer
	lastLine        string
	name            string
	color           string
	desc            string
	facing          int
	facingShapeBase rune
	visibleShape    rune
	shapeMoveCycle  int
	mapContext      *playerMapContext
}

type playerMapContext struct {
	currMap      *dsmap
	currX, currY int
	dsCoords     string
}

var longShapeStart = [][]int{
	{2, 2, 6, 10, 10, 6, 10, 14, 14},
	{2, 2, 7, 12, 12, 7, 12, 17, 17},
	{2, 2, 5, 8, 8, 5, 8, 11, 11}}

var moveCycleLoop = []int{-1, 0, 1, 0}

func playerExec(p *player) {
	var err error
	err = playerLoginLoop(p)
	if err != nil {
		p.logOut()
		return
	}
	playerMainLoop(p)
	p.logOut()
}

func (p *player) logOut() {
	p.mapContext.currMap.removePlayer(p)
	chanCleanDisconns <- p
}

func (p *player) send(message string) (err error) {
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

func (p *player) readLine() (err error) {
	p.lastLine, err = p.reader.ReadString('\n')
	return
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
