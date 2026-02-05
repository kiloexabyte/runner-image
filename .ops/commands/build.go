package commands

import (
	"context"
	"os"

	"lesiw.io/command"
	"lesiw.io/command/sys"
)

func (Ops) Build() error {
	ctx := context.Background()
	sh := command.Shell(sys.Machine(), "docker")

	tag := os.Getenv("IMAGE_TAG")
	if tag == "" {
		tag = "latest"
	}
	imageTag := "kiloexabyte/runner-image:" + tag

	err := sh.Exec(ctx, "docker", "build", "-t", imageTag, ".");
	if err != nil {
		return err
	}

	err = sh.Exec(ctx, "docker", "images", imageTag, "--format", 
		"Image Size: {{.Size}}");
	if err != nil {
		return err
	}

	return nil
}
