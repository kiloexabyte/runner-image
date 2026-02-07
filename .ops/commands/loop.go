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
			"--verbose",
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
	var failures []string

	if err := generalCheck(ctx, m); err != nil {
		failures = append(failures, err.Error())
	}

	if err := taskCheck(ctx, m); err != nil {
		failures = append(failures, err.Error())
	}

	if len(failures) > 0 {
		return fmt.Errorf(
			"success conditions not met:\n%s",
			strings.Join(failures, "\n"),
		)
	}

	return nil
}

func generalCheck(ctx context.Context, m command.Machine) error {
	sh := command.Shell(m, "op")

	log.Printf("Checking: op lint")
	if err := sh.Exec(ctx, "op", "lint"); err != nil {
		return fmt.Errorf("op lint failed: %w", err)
	}

	return nil
}

func taskCheck(_ context.Context, _ command.Machine) error {
	// TODO: implement task-specific checks
	return nil
}