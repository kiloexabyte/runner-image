package commands

import (
	"log"

	"lesiw.io/cmdio/sys"
)

func (Ops) Prune() {
	var rnr = sys.Runner()

	defer rnr.Close()

	// Prune unused Docker images
	// Note: does not remove images that are tagged with "latest"
	err := rnr.Run("docker", "image", "prune", "-f")
	if err != nil {
		log.Fatal(err)
	}
}