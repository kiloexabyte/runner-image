package commands

import (
	"context"
	"log"

	"lesiw.io/command"
	"lesiw.io/command/sys"
)

func (Ops) Prune() {
	ctx := context.Background() 
	sh := command.Shell(sys.Machine(), "docker")

	if err := sh.Exec(ctx, 
		"docker", 
		"image", 
		"prune", 
		"-f"); err != nil {
    	log.Fatal(err)
    }
}