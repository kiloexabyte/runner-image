package commands

import (
	"context"

	"lesiw.io/command"
	"lesiw.io/command/sys"
)

func (Ops) Prune() error {
	ctx := context.Background() 
	sh := command.Shell(sys.Machine(), "docker")

	err := sh.Exec(ctx, "docker", "image", "prune", "-f")
	if  err != nil {
    	return nil
    }

	return nil
}