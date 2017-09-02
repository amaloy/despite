package main

import "fmt"

func playerMainLoop(p *player) error {
	for {
		line, err := playerReadLine(p)
		if err != nil {
			return err
		}
		if line[0] == '"' {
			chanBroadcast <- fmt.Sprintf("(%s: %s", p.name, line[1:])
		} else {
			playerWrite(p, "\n(That just won't do.\n")
		}
	}
}
