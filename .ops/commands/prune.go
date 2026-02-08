package commands

import (
	"context"
	"fmt"

	"lesiw.io/command"
	"lesiw.io/command/sys"
)

func (Ops) Prune() error {
	ctx := context.Background()
	sh := command.Shell(sys.Machine(), "docker")

	err := sh.Exec(ctx, "docker", "image", "prune", "-f")
	if err != nil {
		return fmt.Errorf("prune docker images: %w", err)
	}

	return nil
}
