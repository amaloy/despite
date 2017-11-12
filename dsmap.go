package main

import (
	"math/rand"
)

type dsmap struct {
	name    string
	width   int
	height  int
	xstart  int
	ystart  int
	players map[int]*player
}

const standardMapWidth = 52
const standardMapHeight = 100

func (m *dsmap) getRandomStartCoords() (x int, y int) {
	x = (rand.Intn(5) - 3) + m.xstart
	if x >= m.width {
		x = m.width - 1
	} else if x < 0 {
		x = 0
	}
	y = (rand.Intn(5) - 3) + m.ystart
	if y >= m.height {
		y = m.height - 1
	} else if y < 0 {
		y = 0
	}
	return
}

func (m *dsmap) addPlayer(p *player) {
	p.mapContext.currMap = m
	p.mapContext.currX, p.mapContext.currY = m.getRandomStartCoords()
	p.mapContext.dsCoords = string(toDSChar(p.mapContext.currX)) + string(toDSChar(p.mapContext.currY))

	playerWrite(p, "]"+m.name)
	p.playerWriteAt()

	// Show this player themself
	m.placePlayerBroadcast(p)
	// Show other players to this player
	for _, other := range m.players {
		playerWrite(p, getPlacePlayerString(other))
	}
	m.players[p.connID] = p

	p.resumeMapDraw()
}

func (m *dsmap) removePlayer(p *player) {
	delete(m.players, p.connID)
	chanBroadcast <- "<" + p.mapContext.dsCoords + " "
}

func (m *dsmap) placePlayerBroadcast(p *player) {
	chanBroadcast <- getPlacePlayerString(p)
}

func (m *dsmap) movePlayerBroadcast(p *player, oldDsCoords string) {
	chanBroadcast <- getPlacePlayerString(p) + oldDsCoords + " "
}

func (m *dsmap) nextxy(x int, y int, dir int) (int, int) {
	nx := x
	if dir == 3 || dir == 9 {
		if y%2 == 0 {
			nx++
		}
	} else if y%2 == 1 {
		nx--
	}
	if nx < 0 || nx >= m.width {
		nx = x
	}

	ny := y
	if dir == 7 || dir == 9 {
		ny--
	} else {
		ny++
	}
	if ny < 0 || ny >= m.height {
		ny = y
	}
	return nx, ny
}

func (p *player) haltMapDraw() {
	playerWrite(p, "~")
}

func (p *player) resumeMapDraw() {
	playerWrite(p, "=")
}

func (p *player) playerWriteAt() {
	playerWrite(p, "@"+p.mapContext.dsCoords)
}

func getPlacePlayerString(p *player) string {
	return "<" + p.mapContext.dsCoords + string(p.shape) + p.color
}
