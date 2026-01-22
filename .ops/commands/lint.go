package commands

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"lesiw.io/command"
	"lesiw.io/command/sys"
	"lesiw.io/fs"
)

func (Ops) Lint() {
	ctx := context.Background()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	opsDir := filepath.Join(wd, ".ops")

	machine := sys.Machine()
	sh := command.Shell(machine, "golangci-lint")

	ctx = fs.WithWorkDir(ctx, opsDir)

	if err := sh.Exec(ctx, "golangci-lint", "run"); err != nil {
		log.Fatal(err)
	}
}