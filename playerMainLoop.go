package main

import "fmt"

func playerMainLoop(p *player) (err error) {
	broadcastAll(fmt.Sprintf("(%s has entered %s.", p.name, serverName))

	p.mapContext = new(playerMapContext)
	p.facing = 1
	p.shapeMoveCycle = 0
	p.setShapeStanding()
	mainMap.addPlayer(p)

	for {
		err = p.readLine()
		if err != nil {
			return
		}
		switch p.lastLine[0] {
		case 'm':
			p.move()
		case '"':
			// Typed input
			broadcastMap(fmt.Sprintf("(%s: %s", p.name, p.lastLine[1:len(p.lastLine)-1]), p)
		case '<':
			// Rotate left
			p.rotateLeft()
		case '>':
			// Rotate right
			p.rotateRight()
		default:
			p.send("(That just won't do.")
		}
	}
}

func (p *player) move() {
	p.facing = int(p.lastLine[2]) - 48
	p.setShapeCycleMove()
	p.mapContext.currMap.movePlayer(p, p.facing)
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
	p.shapeMoveCycle = 0
	p.setShapeStanding()
	p.mapContext.currMap.placePlayer(p)
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
	p.shapeMoveCycle = 0
	p.setShapeStanding()
	p.mapContext.currMap.placePlayer(p)
}
