package commands

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
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

func fetchDockerToken() string {
	user := os.Getenv("DOCKER_USERNAME")
	pass := os.Getenv("DOCKER_PASSWORD")

	if user == "" || pass == "" {
		log.Fatal("DOCKER_USERNAME or DOCKER_PASSWORD is not set")
	}

	body := strings.NewReader(
		`{"username":"` + user + `","password":"` + pass + `"}`,
	)

	req, err := http.NewRequest(
		"POST",
		"https://hub.docker.com/v2/users/login/",
		body,
	)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		log.Fatalf(
			"Docker Hub auth failed: %s (%s)",
			resp.Status,
			strings.TrimSpace(string(b)),
		)
	}

	var result struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal(err)
	}

	if result.Token == "" {
		log.Fatal("Docker Hub returned an empty token")
	}

	return result.Token
}

func (Ops) DeleteImage() {
	tag := os.Getenv("IMAGE_TAG")

	if tag == "" {
		log.Fatal(
			"Please provide a tag with -tag flag " +
				"(e.g., op deleteimage -tag PR32)",
		)
	}

	ctx := context.Background()

	m := sys.Machine()
	token := fetchDockerToken()

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
		"-H", "Authorization: Bearer "+token,
		url,
	); err != nil {
		log.Fatal(err)
	}

	log.Printf("Successfully deleted image tag: %s\n", tag)
}
