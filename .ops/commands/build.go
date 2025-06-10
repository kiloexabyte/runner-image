package commands

import (
	"log"

	"lesiw.io/cmdio/sys"
)

func (Ops) Build() {
	var rnr = sys.Runner()
	defer rnr.Close()
	var err error

	err = rnr.Run("docker", "build", "-t", "kiloexabyte/runner-image", ".")
	if err != nil {
		log.Fatal(err)
	}

	err = rnr.Run("docker", "images", "kiloexabyte/runner-image", 
		"--format", "Image Size: {{.Size}}")
	if err != nil {
		log.Fatal(err)
	}
}
