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

	// Read task from file or env var
	task, err := getTask()
	if err != nil {
		return err
	}

	// Setup output file if specified
	outputFile, err := setupOutputFile()
	if err != nil {
		return err
	}
	if outputFile != nil {
		defer outputFile.Close()
	}

	maxIterations := 50
	var lastError string

	for i := range maxIterations {
		logOutput(outputFile, "=== Iteration %d ===", i+1)

		// Build prompt with previous error context
		prompt := task
		if lastError != "" {
			prompt = task +
				"\n\nPrevious attempt failed with:\n" +
				lastError +
				"\n\nTry a different approach."
		}

		// Run claude and capture output
		output, err := sh.Read(ctx, "claude", "-p", prompt,
			"--dangerously-skip-permissions",
			"--verbose",
		)

		// Write claude's output
		if output != "" {
			logOutput(outputFile, "%s", output)
		}

		if err != nil {
			lastError = fmt.Sprintf("claude failed: %v", err)
			logOutput(outputFile, "Attempt %d failed: %v", i+1, err)
			continue
		}

		// Check success conditions
		if err := checkSuccess(ctx, m); err != nil {
			lastError = err.Error()
			logOutput(outputFile, "Success check failed: %v", err)
			continue
		}

		logOutput(outputFile, "Task completed successfully")
		return nil
	}

	return fmt.Errorf("failed after %d iterations", maxIterations)
}

func getTask() (string, error) {
	// Check env var first
	if task := os.Getenv("TASK"); task != "" {
		return task, nil
	}

	// Try task file (env var or default)
	taskFile := os.Getenv("TASK_FILE")
	if taskFile == "" {
		taskFile = "task.txt"
	}

	content, err := os.ReadFile(taskFile)
	if err != nil {
		return "", fmt.Errorf(
			"read task file %s: %w (set TASK env var or create task.txt)",
			taskFile,
			err,
		)
	}

	task := strings.TrimSpace(string(content))
	if task == "" {
		return "", fmt.Errorf("task file %s is empty", taskFile)
	}

	return task, nil
}

func setupOutputFile() (*os.File, error) {
	outputPath := os.Getenv("OUTPUT_FILE")
	if outputPath == "" {
		return nil, nil
	}

	f, err := os.OpenFile(
		outputPath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, fmt.Errorf("open output file: %w", err)
	}
	return f, nil
}

func logOutput(f *os.File, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	log.Print(msg)
	if f != nil {
		fmt.Fprintln(f, msg)
	}
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