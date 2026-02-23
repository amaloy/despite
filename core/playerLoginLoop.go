package core

import (
	"strings"
)

func playerLoginLoop(p *Player) (err error) {
	p.Send(MOTD)
	p.Send("\nDragonroar!\nV0026")

	isnew, err := cmdConnect(p)
	if err != nil {
		return
	}
	if isnew {
		err = cmdColor(p)
	}
	return
}

func cmdConnect(p *Player) (isnew bool, err error) {
	for {
		err = p.readLine()
		if err != nil {
			return
		}
		p.lastLine = strings.TrimSpace(p.lastLine)

		if strings.HasPrefix(p.lastLine, "connect") {
			split := strings.Split(p.lastLine, " ")
			if len(split) < 3 {
				p.Send("(Usage: connect <username> <password>)")
				continue
			}
			username := split[1]
			password := split[2]
			if authenticate(username, password) {
				p.Name = username
				p.Send("(Connected.)")
				return true, nil
			} else {
				p.Send("(Invalid username or password.)")
			}
		} else {
			p.Send("(Please connect with 'connect <username> <password>'.)")
		}
	}
}

func authenticate(username string, password string) bool {
	// TODO
	return true
}

func cmdColor(p *Player) (err error) {
	p.Send("cs")
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
		err = cmdDesc(p)
	}
	return
}

func cmdDesc(p *Player) (err error) {
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
	}
	return
}
