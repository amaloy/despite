package main

import (
	"bufio"
	"context"
	"net"
	"os"

	"github.com/amaloy/despite/core"
	"github.com/google/uuid"
)

type Server struct {
	allPlayers     map[uuid.UUID]*core.Player
	newConnections chan net.Conn
}

func (s *Server) Start(ctx context.Context) error {
	err := core.InitMainMap()
	if err != nil {
		core.Logger.Error("failed to initialize main map", "err", err)
		return err
	}

	server, err := net.Listen("tcp", ":7734")
	if err != nil {
		core.Logger.Error("failed to listen on port 7734", "err", err)
		return err
	}

	go func() {
		for {
			// Accept new connections
			conn, err := server.Accept()
			if err != nil {
				core.Logger.Error("failed to accept connection", "err", err)
				return
			}
			select {
			case s.newConnections <- conn:
			case <-ctx.Done():
				return
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case conn := <-s.newConnections:
			core.Logger.Info("Accepted new connection", "addr", conn.RemoteAddr())

			p := new(core.Player)
			p.ConnID = uuid.New()
			p.Conn = conn
			p.Reader = bufio.NewReader(conn)
			p.Writer = bufio.NewWriter(conn)

			s.allPlayers[p.ConnID] = p

			// Spawn independant player exec
			go core.PlayerExec(ctx, p)

		case payload := <-core.ChanBroadcast:
			var targets map[uuid.UUID]*core.Player
			if payload.TargetMap == nil {
				targets = s.allPlayers
			} else {
				targets = payload.TargetMap.Players
			}

			for _, p := range targets {
				if p != payload.ExcludePlayer {
					go p.Send(*payload.Message)
				}
			}

		case p := <-core.ChanCleanDisconns:
			core.Logger.Info("Player disconnected", "name", p.Name, "addr", p.Conn.RemoteAddr())
			delete(s.allPlayers, p.ConnID)
			p.Conn.Close()
		}
	}
}

func main() {
	s := &Server{
		allPlayers:     make(map[uuid.UUID]*core.Player),
		newConnections: make(chan net.Conn),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := s.Start(ctx)
	if err != nil {
		os.Exit(1)
	}
}
