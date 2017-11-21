package main

import (
	"math/rand"
	"os"

	"github.com/satori/go.uuid"
)

type dsmapTile struct {
	hasPlayer        bool
	hasBlockingFloor bool
}

type dsmap struct {
	name           string
	width, height  int
	tiles          [][]*dsmapTile
	xstart, ystart int
	players        map[uuid.UUID]*player
}

const standardMapWidth = 52
const standardMapHeight = 100

var floorwalk = []int{
	0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 0, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 0,
	0, 1, 0, 0, 1, 0, 1, 1, 1,
	1, 0, 0, 0, 1, 0, 1, 1, 1,
	1, 1, 1, 0, 0, 0, 1, 1, 0}

func (m *dsmap) readMapFromFile(filename string) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	buff := make([]byte, m.height*2)
	var temp int
	// Read floor tiles
	for x := 0; x < m.width; x++ {
		f.Read(buff)
		for y := 0; y < m.height*2; y += 2 {
			temp = int(buff[y])*95 + int(buff[y+1])
			m.tiles[x][y/2].hasBlockingFloor = floorwalk[temp] == 1
		}
	}
	// TODO Read items
	f.Close()
	return
}

func (m *dsmap) getRandomStartCoords() (x, y int) {
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

	p.send("]" + m.name)
	p.playerWriteAt()

	m.players[p.connID] = p
	// Show this player
	m.placePlayer(p)
	// Show other players to this player
	for _, other := range m.players {
		p.send(getPlacePlayerString(other))
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
		p.send(message)
		p.resumeMapDraw()
		broadcastMapExclude(message, p)
	} else {
		m.placePlayer(p)
	}
}

func (m *dsmap) nextxy(x, y, dir int) (nx, ny int) {
	nx = x
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

	ny = y
	if dir == 7 || dir == 9 {
		ny--
	} else {
		ny++
	}
	if ny < 0 || ny >= m.height {
		ny = y
	}
	return
}

func (m *dsmap) tileIsBlocked(x, y int) bool {
	tile := m.tiles[x][y]
	return tile.hasBlockingFloor || tile.hasPlayer
}

func (p *player) haltMapDraw() {
	p.send("~")
}

func (p *player) resumeMapDraw() {
	p.send("=")
}

func (p *player) playerWriteAt() {
	p.send("@" + p.mapContext.dsCoords)
}

func getPlacePlayerString(p *player) string {
	return "<" + p.mapContext.dsCoords + string(p.visibleShape) + p.color
}
