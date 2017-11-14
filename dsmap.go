package main

import (
	"math/rand"
)

type dsmapTile struct {
	hasPlayer bool
}

type dsmap struct {
	name    string
	width   int
	height  int
	tiles   [][]*dsmapTile
	xstart  int
	ystart  int
	players map[int]*player
}

const standardMapWidth = 52
const standardMapHeight = 100

func (m *dsmap) getRandomStartCoords() (x int, y int) {
	for {
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
		if !m.tileIsBlocked(x, y) {
			break
		}
	}
	return
}

func (m *dsmap) addPlayer(p *player) {
	p.mapContext.currMap = m
	p.mapContext.currX, p.mapContext.currY = m.getRandomStartCoords()
	m.tiles[p.mapContext.currX][p.mapContext.currY].hasPlayer = true
	p.mapContext.dsCoords = string(toDSChar(p.mapContext.currX)) + string(toDSChar(p.mapContext.currY))

	playerWrite(p, "]"+m.name)
	p.playerWriteAt()

	m.players[p.connID] = p
	// Show this player
	m.placePlayer(p)
	// Show other players to this player
	for _, other := range m.players {
		playerWrite(p, getPlacePlayerString(other))
	}

	p.resumeMapDraw()
}

func (m *dsmap) removePlayer(p *player) {
	delete(m.players, p.connID)
	m.tiles[p.mapContext.currX][p.mapContext.currY].hasPlayer = false
	broadcastMapExclude("<"+p.mapContext.dsCoords+" ", p)
}

func (m *dsmap) placePlayer(p *player) {
	m.tiles[p.mapContext.currX][p.mapContext.currY].hasPlayer = true
	broadcastMap(getPlacePlayerString(p), p)
}

func (m *dsmap) movePlayer(p *player, dir int) {
	nx, ny := p.mapContext.currMap.nextxy(
		p.mapContext.currX, p.mapContext.currY, dir)
	if !m.tileIsBlocked(nx, ny) {
		m.tiles[p.mapContext.currX][p.mapContext.currY].hasPlayer = false
		p.mapContext.currX, p.mapContext.currY = nx, ny
		m.tiles[p.mapContext.currX][p.mapContext.currY].hasPlayer = true
		oldDsCoords := p.mapContext.dsCoords
		p.mapContext.dsCoords = string(toDSChar(p.mapContext.currX)) + string(toDSChar(p.mapContext.currY))
		p.haltMapDraw()
		p.playerWriteAt()
		message := getPlacePlayerString(p) + oldDsCoords + " "
		playerWrite(p, message)
		p.resumeMapDraw()
		broadcastMapExclude(message, p)
	} else {
		m.placePlayer(p)
	}
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

func (m *dsmap) tileIsBlocked(x int, y int) bool {
	return m.tiles[x][y].hasPlayer
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
