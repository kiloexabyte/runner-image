package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"lesiw.io/command"
	"lesiw.io/command/sys"
)

func (Ops) Loop() error {
	ctx := context.Background()
	m := sys.Machine()
	sh := command.Shell(m, "claude", "op")

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

		// Run claude with full autonomy
		err := sh.Exec(ctx, "claude", "-p", prompt,
			"--dangerously-skip-permissions",
		)
		if err != nil {
			lastError = fmt.Sprintf("claude failed: %v", err)
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

func checkSuccess(ctx context.Context, m command.Machine) error {
	sh := command.Shell(m, "op")
	var failures []string

	// Check 1: op lint
	log.Printf("Checking: op lint")
	if err := sh.Exec(ctx, "op", "lint"); err != nil {
		failures = append(failures, fmt.Sprintf("op lint failed: %v", err))
	}

	// Check 2: op version (verifies build works)
	log.Printf("Checking: op version")
	if err := sh.Exec(ctx, "op", "version"); err != nil {
		failures = append(failures, fmt.Sprintf("op version failed: %v", err))
	}


	if len(failures) > 0 {
		return fmt.Errorf(
			"success conditions not met:\n%s",
			strings.Join(failures, "\n"),
		)
	}

	return nil
}