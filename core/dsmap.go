package core

import (
	"context"
	"fmt"
	"math/rand"
	"os"

	"github.com/google/uuid"
)

type dsmapTile struct {
	hasPlayer        bool
	hasBlockingFloor bool
}

type DSMap struct {
	name           string
	width, height  int
	tiles          [][]*dsmapTile
	xstart, ystart int
	Players        map[uuid.UUID]*Player
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

func (m *DSMap) readMapFromFile(filename string) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open map file %s: %w", filename, err)
	}
	defer f.Close()
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
	return
}

func (m *DSMap) getRandomStartCoords() (x, y int) {
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

func (m *DSMap) addPlayer(ctx context.Context, p *Player) {
	p.mapContext.currMap = m
	p.mapContext.currX, p.mapContext.currY = m.getRandomStartCoords()
	m.tiles[p.mapContext.currX][p.mapContext.currY].hasPlayer = true
	p.mapContext.dsCoords = string(toDSChar(p.mapContext.currX)) + string(toDSChar(p.mapContext.currY))

	p.Send("]" + m.name)
	p.playerWriteAt()

	m.Players[p.ConnID] = p
	// Show this player
	m.placePlayer(ctx, p)
	// Show other players to this player
	for _, other := range m.Players {
		p.Send(getPlacePlayerString(other))
	}

	p.resumeMapDraw()
}

func (m *DSMap) removePlayer(ctx context.Context, p *Player) {
	delete(m.Players, p.ConnID)
	m.tiles[p.mapContext.currX][p.mapContext.currY].hasPlayer = false
	broadcastMapExclude(ctx, "<"+p.mapContext.dsCoords+" ", p)
}

func (m *DSMap) placePlayer(ctx context.Context, p *Player) {
	m.tiles[p.mapContext.currX][p.mapContext.currY].hasPlayer = true
	broadcastMap(ctx, getPlacePlayerString(p), p)
}

func (m *DSMap) movePlayer(ctx context.Context, p *Player, dir int) {
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
		p.Send(message)
		p.resumeMapDraw()
		broadcastMapExclude(ctx, message, p)
	} else {
		m.placePlayer(ctx, p)
	}
}

func (m *DSMap) nextxy(x, y, dir int) (nx, ny int) {
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

func (m *DSMap) tileIsBlocked(x, y int) bool {
	tile := m.tiles[x][y]
	return tile.hasBlockingFloor || tile.hasPlayer
}

func (p *Player) haltMapDraw() {
	p.Send("~")
}

func (p *Player) resumeMapDraw() {
	p.Send("=")
}

func (p *Player) playerWriteAt() {
	p.Send("@" + p.mapContext.dsCoords)
}

func getPlacePlayerString(p *Player) string {
	return "<" + p.mapContext.dsCoords + string(p.visibleShape) + p.color
}
