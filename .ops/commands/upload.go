package commands

import (
	"context"
	"io"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"lesiw.io/command"
	"lesiw.io/command/sys"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
		log.Println("Skipping loading .env file")
	}
}

func (Ops) Upload() {
	// Read env vars once
	user := os.Getenv("DOCKER_USERNAME")
	pass := os.Getenv("DOCKER_PASSWORD")

	// Build a context with environment variables
	ctx := context.Background()
	envVars := map[string]string{
		"DOCKER_USERNAME": user,
		"DOCKER_PASSWORD": pass,
	}
	ctx = command.WithEnv(ctx, envVars)

	m := sys.Machine()

	// 1) docker login --password-stdin
	// Create a writer command for docker login
	loginStdin := command.NewWriter(ctx, m,
		"docker", "login",
		"-u", user,
		"--password-stdin",
	)

	// Pipe the password into stdin
	if _, err := io.Copy(loginStdin, strings.NewReader(pass)); err != nil {
		log.Fatal(err)
	}

	// Close stdin and wait for the command to finish
	if err := loginStdin.Close(); err != nil {
		log.Fatal(err)
	}

	// 2) docker push
	if err := command.Do(ctx, m, "docker", 
		"push",
		"kiloexabyte/runner-image:latest"); err != nil {
		log.Fatal(err)
	}
}
