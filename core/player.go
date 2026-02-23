package core

import (
	"bufio"
	"context"
	"net"

	"github.com/google/uuid"
)

type Player struct {
	ConnID          uuid.UUID
	Conn            net.Conn
	Reader          *bufio.Reader
	Writer          *bufio.Writer
	lastLine        string
	Name            string
	color           string
	desc            string
	facing          int
	facingShapeBase rune
	visibleShape    rune
	shapeMoveCycle  int
	mapContext      *playerMapContext
	Ctx             context.Context
}

type playerMapContext struct {
	currMap      *DSMap
	currX, currY int
	dsCoords     string
}

var longShapeStart = [][]int{
	{2, 2, 6, 10, 10, 6, 10, 14, 14},
	{2, 2, 7, 12, 12, 7, 12, 17, 17},
	{2, 2, 5, 8, 8, 5, 8, 11, 11}}

var moveCycleLoop = []int{-1, 0, 1, 0}

func PlayerExec(ctx context.Context, p *Player) {
	p.Ctx = ctx
	var err error
	err = playerLoginLoop(p)
	if err != nil {
		p.logOut()
		return
	}
	playerMainLoop(ctx, p)
	p.logOut()
}

func (p *Player) logOut() {
	if p.mapContext != nil && p.mapContext.currMap != nil {
		p.mapContext.currMap.removePlayer(p.Ctx, p)
	}
	ChanCleanDisconns <- p
}

func (p *Player) Send(message string) (err error) {
	_, err = p.Writer.WriteString(message)
	if err != nil {
		return err
	}
	_, err = p.Writer.WriteRune('\n')
	if err != nil {
		return err
	}
	err = p.Writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (p *Player) readLine() (err error) {
	p.lastLine, err = p.Reader.ReadString('\n')
	return
}

func (p *Player) setShapeStanding() {
	p.facingShapeBase = toDSChar(longShapeStart[1][p.facing-1])
	p.visibleShape = p.facingShapeBase
}

func (p *Player) setShapeCycleMove() {
	p.facingShapeBase = toDSChar(longShapeStart[1][p.facing-1])
	p.visibleShape = p.facingShapeBase + rune(moveCycleLoop[p.shapeMoveCycle])
	p.shapeMoveCycle++
	if p.shapeMoveCycle == 4 {
		p.shapeMoveCycle = 0
	}
}
