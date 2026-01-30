package commands

import (
	"context"
	"log"
	"os"

	"lesiw.io/command"
	"lesiw.io/command/sys"
)

func (Ops) Build() {
	ctx := context.Background()
	sh := command.Shell(sys.Machine(), "docker")

	tag := os.Getenv("IMAGE_TAG")
	if tag == "" {
		tag = "latest"
	}
	imageTag := "kiloexabyte/runner-image:" + tag

	if err := sh.Exec(ctx, "docker",
		"build",
		"-t",
		imageTag, "."); err != nil {
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
