package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"lesiw.io/command"
	"lesiw.io/command/sys"
)

func (Ops) Loop() error {
	ctx := context.Background()
	m := sys.Machine()
	sh := command.Shell(m, "gh", "go", "op")

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
		err := sh.Exec(ctx, "gh", "copilot",
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

func checkSuccess(ctx context.Context, m command.Machine) error {
	sh := command.Shell(m, "gh", "go", "op")
	var failures []string

	// Check 1: op lint
	log.Printf("Checking: op lint")
	if err := sh.Exec(ctx, "op", "lint"); err != nil {
		failures = append(failures, fmt.Sprintf("op lint failed: %v", err))
	}

	// Check 2: go build
	log.Printf("Checking: go build")
	if err := sh.Exec(ctx, "go", "build", "./..."); err != nil {
		failures = append(failures, fmt.Sprintf("go build failed: %v", err))
	}

	// Check 3: CI status on PR
	log.Printf("Checking: CI status")
	if err := checkCI(ctx, m); err != nil {
		failures = append(failures, fmt.Sprintf("CI check failed: %v", err))
	}

	if len(failures) > 0 {
		return fmt.Errorf(
			"success conditions not met:\n%s",
			strings.Join(failures, "\n"),
		)
	}

	return nil
}

func checkCI(ctx context.Context, m command.Machine) error {
	sh := command.Shell(m, "gh")

	// Wait a bit for CI to start
	time.Sleep(10 * time.Second)

	// Check PR status - this will fail if no PR or checks failing
	output, err := sh.Read(ctx, "gh", "pr", "checks", "--watch", "--fail-fast")
	if err != nil {
		return fmt.Errorf("%v: %s", err, output)
	}

	return nil
}