package main

import (
	"math/rand"
)

type dsmap struct {
	name   string
	width  int
	height int
	xstart int
	ystart int
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
	playerWrite(p, "@"+p.mapContext.dsCoords)
	playerWrite(p, "=")
}
