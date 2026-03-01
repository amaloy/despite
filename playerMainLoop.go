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
		if len(p.lastLine) == 0 {
			continue
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
	// Check if the command has the correct format (e.g., "m1")
	if len(p.lastLine) < 2 {
		p.send("(Invalid move command format.")
		return
	}

	// Parse the direction from the command
	dir := int(p.lastLine[1]) - 48
	// Validate direction is within acceptable range (1,3,7,9 based on rotation logic)
	if dir != 1 && dir != 3 && dir != 7 && dir != 9 {
		p.send("(Invalid direction for move.")
		return
	}

	p.facing = dir
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
