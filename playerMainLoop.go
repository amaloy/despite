package main

import "fmt"

var longShapeStart = [][]int{
	{2, 2, 6, 10, 10, 6, 10, 14, 14},
	{2, 2, 7, 12, 12, 7, 12, 17, 17},
	{2, 2, 5, 8, 8, 5, 8, 11, 11}}

func playerMainLoop(p *player) (err error) {
	chanBroadcast <- fmt.Sprintf("(%s has entered %s.", p.name, serverName)

	p.mapContext = new(playerMapContext)
	p.shape = '"'
	p.facing = 1
	mainMap.addPlayer(p)

	for {
		p.readLine, err = playerReadLine(p)
		if err != nil {
			return
		}
		switch p.readLine[0] {
		case 'm':
			p.move()
		case '"':
			// Typed input
			chanBroadcast <- fmt.Sprintf("(%s: %s", p.name, p.readLine[1:len(p.readLine)-1])
		case '<':
			// Rotate left
			p.rotateLeft()
		case '>':
			// Rotate right
			p.rotateRight()
		default:
			playerWrite(p, "(That just won't do.")
		}
	}
}

func (p *player) move() {
	p.facing = int(p.readLine[2]) - 48
	p.shape = toDSChar(longShapeStart[1][p.facing-1])
	// TODO complete sprite handling

	nx, ny := p.mapContext.currMap.nextxy(
		p.mapContext.currX, p.mapContext.currY, p.facing)
	// TODO: Check if can move here

	p.mapContext.currX, p.mapContext.currY = nx, ny
	oldDsCoords := p.mapContext.dsCoords
	p.mapContext.dsCoords = string(toDSChar(p.mapContext.currX)) + string(toDSChar(p.mapContext.currY))
	p.haltMapDraw()
	p.playerWriteAt()
	// Maybe send to current player synchronously?
	p.mapContext.currMap.movePlayerBroadcast(p, oldDsCoords)
	p.resumeMapDraw()
}

func (p *player) rotateLeft() {
	switch p.facing {
	case 7:
		p.facing = 1
	case 9:
		p.facing = 7
	case 1:
		p.facing = 3
	case 3:
		p.facing = 9
	}
	p.shape = toDSChar(longShapeStart[1][p.facing-1])
	p.mapContext.currMap.placePlayerBroadcast(p)
}

func (p *player) rotateRight() {
	switch p.facing {
	case 7:
		p.facing = 9
	case 9:
		p.facing = 3
	case 1:
		p.facing = 7
	case 3:
		p.facing = 1
	}
	p.shape = toDSChar(longShapeStart[1][p.facing-1])
	p.mapContext.currMap.placePlayerBroadcast(p)
}
