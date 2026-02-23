package core

import (
	_ "embed"
	"strings"
)

//go:embed resources/motd.txt
var motdBytes []byte

// MOTD is the embedded server message-of-the-day.
var MOTD string

func init() {
	MOTD = strings.TrimSpace(string(motdBytes))
	if MOTD == "" {
		MOTD = "Welcome to DragonSpires (fallback MOTD)"
		Logger.Warn("MOTD file was empty or missing - using fallback")
	}
}
