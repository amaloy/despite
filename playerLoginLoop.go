package main

import (
	"errors"
	"log"
	"strings"
)

func playerLoginLoop(p *player) (err error) {
	p.send(motd)
	p.send("\nDragonroar!\nV0026")

	isnew, err := cmdConnect(p)
	if err != nil {
		return
	}
	if isnew {
		err = cmdColor(p)
		if err != nil {
			return
		}
		err = cmdDesc(p)
		if err != nil {
			return
		}
	}

	p.send("&")

	log.Printf("Logged in: %s", p.name)

	return
}

func cmdConnect(p *player) (isnew bool, err error) {
	for {
		err = p.readLine()
		if err != nil {
			return
		}
		p.lastLine = strings.TrimSpace(p.lastLine)

		if strings.HasPrefix(p.lastLine, "connect") {
			split := strings.Split(p.lastLine, " ")
			if len(split) != 3 {
				break
			}

			// TODO allow and keep track of multiple attempts
			if authenticate(split[1], split[2]) {
				p.name = split[1]
				isnew = true
				return
			}
			break
		} else {
			break
		}
	}
	return false, errors.New("cmdConnect")
}

func authenticate(username string, password string) bool {
	// TODO
	return true
}

func cmdColor(p *player) (err error) {
	p.send("cs")
	err = p.readLine()
	if err != nil {
		return
	}
	if strings.HasPrefix(p.lastLine, "color") {
		color := p.lastLine[6 : len(p.lastLine)-1]
		if len(color) != 4 {
			color = "   !"
		}
		p.color = color
		return
	}
	return errors.New("cmdColor")
}

func cmdDesc(p *player) (err error) {
	err = p.readLine()
	if err != nil {
		return
	}
	if strings.HasPrefix(p.lastLine, "desc") {
		desc := p.lastLine[6 : len(p.lastLine)-1]
		if len(desc) > 500 {
			desc = desc[:500]
		}
		p.desc = desc
		return
	}
	return errors.New("cmdDesc")
}
