package commands

import (
	"context"

	"lesiw.io/command"
	"lesiw.io/command/sys"
)

func (Ops) Prune() error {
	ctx := context.Background() 
	sh := command.Shell(sys.Machine(), "docker")

	if err := sh.Exec(ctx, 
		"docker", 
		"image", 
		"prune", 
		"-f"); err != nil {
    	return nil
    }

	return nil
}