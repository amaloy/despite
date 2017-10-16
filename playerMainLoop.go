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
		case '"':
			// Typed input
			chanBroadcast <- fmt.Sprintf("(%s: %s", p.name, p.readLine[1:len(p.readLine)-1])
		case '<':
			// Rotate left
			playerRotateLeft(p)
		case '>':
			// Rotate right
			playerRotateRight(p)
		default:
			playerWrite(p, "(That just won't do.")
		}
	}
}

func playerRotateLeft(p *player) {
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

func playerRotateRight(p *player) {
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
