package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

func playerLoginLoop(p *player) (err error) {
	playerWrite(p, motd)
	playerWrite(p, "\nDragonroar!\nV0026")

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

	playerWrite(p, "&")
	playerWrite(p, fmt.Sprintf("PY%s", p.pstring))

	log.Printf("Logged in: %s", p.name)

	return
}

func cmdConnect(p *player) (isnew bool, err error) {
	for {
		p.readLine, err = playerReadLine(p)
		if err != nil {
			return
		}
		p.readLine = strings.TrimSpace(p.readLine)

		if strings.HasPrefix(p.readLine, "connect") {
			split := strings.Split(p.readLine, " ")
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
	playerWrite(p, "cs")
	p.readLine, err = playerReadLine(p)
	if err != nil {
		return
	}
	if strings.HasPrefix(p.readLine, "color") {
		color := p.readLine[6 : len(p.readLine)-1]
		if len(color) != 4 {
			color = "   !"
		}
		// TODO implement color/pstring fully
		p.color = color
		p.pstring = "!'+!"
		return
	}
	return errors.New("cmdColor")
}

func cmdDesc(p *player) (err error) {
	p.readLine, err = playerReadLine(p)
	if err != nil {
		return
	}
	if strings.HasPrefix(p.readLine, "desc") {
		desc := p.readLine[6 : len(p.readLine)-1]
		if len(desc) > 500 {
			desc = desc[:500]
		}
		p.desc = desc
		return
	}
	return errors.New("cmdDesc")
}
