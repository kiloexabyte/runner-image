package commands

import (
	"context"
	"log"

	"lesiw.io/command"
	"lesiw.io/command/sys"
)

func (Ops) Build() {
	ctx := context.Background() 
	sh := command.Shell(sys.Machine(), "docker")

	if err := sh.Exec(ctx, "docker", 
		"build", 
		"-t", 
		"kiloexabyte/runner-image", "."); err != nil {
    	log.Fatal(err)
    }

	if err := sh.Exec(ctx, "docker",
		"images",
		"kiloexabyte/runner-image",
		"--format", 
		"Image Size: {{.Size}}"); err != nil {

    	log.Fatal(err)
    }
}