package main

import "fmt"

func playerMainLoop(p *player) (err error) {
	chanBroadcast <- fmt.Sprintf("(%s has entered %s.", p.name, serverName)

	p.mapContext = new(playerMapContext)
	p.shape = '"'
	mainMap.addPlayer(p)

	for {
		p.readLine, err = playerReadLine(p)
		if err != nil {
			return
		}
		if p.readLine[0] == '"' {
			chanBroadcast <- fmt.Sprintf("(%s: %s", p.name, p.readLine[1:len(p.readLine)-1])
		} else {
			playerWrite(p, "(That just won't do.")
		}
	}
}
