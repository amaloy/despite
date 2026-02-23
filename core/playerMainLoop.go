package core

import (
	"context"
	"fmt"
)

func playerMainLoop(ctx context.Context, p *Player) (err error) {
	broadcastAll(ctx, fmt.Sprintf("(%s has entered %s.", p.Name, serverName))

	p.mapContext = new(playerMapContext)
	p.facing = 1
	p.shapeMoveCycle = 0
	p.setShapeStanding()
	MainMap.addPlayer(ctx, p)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err = p.readLine()
			if err != nil {
				return
			}
			switch p.lastLine[0] {
			case 'm':
				p.move()
			case '"':
				// Typed input
				broadcastMap(ctx, fmt.Sprintf("(%s: %s", p.Name, p.lastLine[1:len(p.lastLine)-1]), p)
			case '<':
				// Rotate left
				p.rotateLeft()
			case '>':
				// Rotate right
				p.rotateRight()
			default:
				p.Send("(That just won't do.")
			}
		}
	}
}

func (p *Player) move() {
	p.facing = int(p.lastLine[2]) - 48
	p.setShapeCycleMove()
	p.mapContext.currMap.movePlayer(p.Ctx, p, p.facing)
}

func (p *Player) rotateLeft() {
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
	p.mapContext.currMap.placePlayer(p.Ctx, p)
}

func (p *Player) rotateRight() {
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
	p.mapContext.currMap.placePlayer(p.Ctx, p)
}
