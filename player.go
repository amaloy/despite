package main

import (
	"bufio"
	"log"
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
}

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
	log.Printf("Sending to %s: %s", p.name, message)
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
		log.Printf("Error flushing message to %s: %v", p.name, err)
		// Do not retry on short write; log the error and return it
		// This allows higher-level logic to handle connection issues
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
	if p.facing < 1 || p.facing > 9 {
		log.Printf("Invalid facing value for player %s: %d", p.name, p.facing)
		p.facing = 1 // Default to a valid facing value
	}
	p.facingShapeBase = toDSChar(longShapeStart[1][p.facing-1])
	p.visibleShape = p.facingShapeBase + rune(moveCycleLoop[p.shapeMoveCycle])
	p.shapeMoveCycle++
	if p.shapeMoveCycle == 4 {
		p.shapeMoveCycle = 0
	}
}
