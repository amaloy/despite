package core

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

const serverName string = "Despite"

type broadcastPayload struct {
	Message       *string
	TargetMap     *DSMap
	ExcludePlayer *Player
}

var Logger = slog.Default()
var ChanCleanDisconns = make(chan *Player)
var ChanBroadcast = make(chan *broadcastPayload)
var MainMap *DSMap

func InitMainMap() error {
	var err error
	MainMap, err = BuildMainMap()
	return err
}

func broadcast(ctx context.Context, message *string, targetMap *DSMap, excludePlayer *Player) {
	ChanBroadcast <- &broadcastPayload{message, targetMap, excludePlayer}
}

func broadcastAll(ctx context.Context, message string) {
	broadcast(ctx, &message, nil, nil)
}

func broadcastMap(ctx context.Context, message string, p *Player) {
	broadcast(ctx, &message, p.mapContext.currMap, nil)
}

func broadcastMapExclude(ctx context.Context, message string, p *Player) {
	broadcast(ctx, &message, p.mapContext.currMap, p)
}

func toDSChar(i int) rune {
	return (rune)(i + 32)
}

func BuildMainMap() (*DSMap, error) {
	m := new(DSMap)
	m.name = "lev01"
	m.width = standardMapWidth
	m.height = standardMapHeight
	m.tiles = make([][]*dsmapTile, m.width)
	for x := range m.tiles {
		row := make([]*dsmapTile, m.height)
		for y := range row {
			row[y] = new(dsmapTile)
		}
		m.tiles[x] = row
	}
	m.xstart = 26
	m.ystart = 41
	err := m.readMapFromFile("core/resources/" + m.name + ".dsmap")
	if err != nil {
		return nil, err
	}
	m.Players = make(map[uuid.UUID]*Player)
	return m, nil
}
