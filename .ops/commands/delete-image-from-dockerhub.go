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

func (Ops) Deleteimage() {
	tag := os.Getenv("IMAGE_TAG")

	if tag == "" {
		log.Fatal(
			"Please provide a tag with -tag flag " +
				"(e.g., op deleteimage -tag PR32)",
		)
	}

	// Read env vars
	user := os.Getenv("DOCKER_USERNAME")
	pass := os.Getenv("DOCKER_PASSWORD")

	ctx := context.Background()
	envVars := map[string]string{
		"DOCKER_USERNAME": user,
		"DOCKER_PASSWORD": pass,
	}
	ctx = command.WithEnv(ctx, envVars)

	m := sys.Machine()

	// Docker login
	loginStdin := command.NewWriter(ctx, m,
		"docker", "login",
		"-u", user,
		"--password-stdin",
	)

	if _, err := io.Copy(loginStdin, strings.NewReader(pass)); err != nil {
		log.Fatal(err)
	}

	if err := loginStdin.Close(); err != nil {
		log.Fatal(err)
	}

	// Delete image from Docker Hub using registry API
	imageTag := "kiloexabyte/runner-image:" + tag
	log.Printf("Deleting image: %s\n", imageTag)

	url := "https://hub.docker.com/v2/repositories/" +
		"kiloexabyte/runner-image/tags/" +
		tag + "/"

	if err := command.Do(
		ctx,
		m,
		"curl",
		"-X", "DELETE",
		"-H", "Authorization: Bearer "+pass,
		url,
	); err != nil {
		log.Fatal(err)
	}

	log.Printf("Successfully deleted image tag: %s\n", tag)
}
