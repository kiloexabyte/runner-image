package commands

import (
	"context"
	"fmt"

	"lesiw.io/command"
	"lesiw.io/command/sys"
	"lesiw.io/fs"
)

func (Ops) Lint() error {
	sh := command.Shell(sys.Machine(), "golangci-lint", "go")
	ctx := fs.WithWorkDir(context.Background(), ".ops")

	if err := sh.Exec(ctx, "golangci-lint", "run"); err != nil {
		return fmt.Errorf("golangci-lint: %w", err)
	}

	if err := sh.Exec(ctx, "go", "fmt"); err != nil {
		return fmt.Errorf("go fmt: %w", err)
	}

	return nil
}