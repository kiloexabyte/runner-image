package commands

import (
	"context"
	"log"

	"lesiw.io/command"
	"lesiw.io/command/sys"
	"lesiw.io/fs"
)

func (Ops) Lint() {
	sh := command.Shell(sys.Machine(), "golangci-lint")
	ctx := fs.WithWorkDir(context.Background(), ".ops")

	if err := sh.Exec(ctx, "golangci-lint", "run"); err != nil {
		log.Fatal(err)
	}
}