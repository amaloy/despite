package main

func playerLoginLoop(p *player) error {
	playerWrite(p, motd)
	playerWrite(p, "\nDragonroar!\nV0026\n")
	return nil
}
