package commands

import (
	"context"
	"fmt"
	"log"
	"os"

	"lesiw.io/command"
	"lesiw.io/command/sys"
)

const ghBin = `C:\Program Files\GitHub CLI\gh.exe`

func (Ops) LoopCopilot() error {
	ctx := context.Background()
	m := sys.Machine()
	sh := command.Shell(m, ghBin, "op")

	task := os.Getenv("TASK")
	if task == "" {
		return fmt.Errorf("TASK environment variable is required")
	}

	maxIterations := 50
	var lastError string

	for i := range maxIterations {
		log.Printf("=== Iteration %d ===", i+1)

		// Build prompt with previous error context
		prompt := task
		if lastError != "" {
			prompt = task +
				"\n\nPrevious attempt failed with:\n" +
				lastError +
				"\n\nTry a different approach."
		}

		// Run copilot with full autonomy
		err := sh.Exec(ctx, ghBin, "copilot",
			"-p", prompt,
			"--yolo",
			"--no-ask-user",
		)
		if err != nil {
			lastError = fmt.Sprintf("copilot failed: %v", err)
			log.Printf("Attempt %d failed: %v", i+1, err)
			continue
		}

		// Check success conditions
		if err := checkSuccess(ctx, m); err != nil {
			lastError = err.Error()
			log.Printf("Success check failed: %v", err)
			continue
		}

		log.Printf("Task completed successfully")
		return nil
	}

	return fmt.Errorf("failed after %d iterations", maxIterations)
}
