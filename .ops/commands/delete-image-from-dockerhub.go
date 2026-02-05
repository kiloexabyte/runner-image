package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
		log.Println("Skipping loading .env file")
	}
}

func fetchDockerToken() (string, error) {
	user := os.Getenv("DOCKER_USERNAME")
	pass := os.Getenv("DOCKER_PASSWORD")

	if user == "" || pass == "" {
		return "", fmt.Errorf(
			"DOCKER_USERNAME or DOCKER_PASSWORD is not set",
		)
	}

	body := strings.NewReader(
		`{"username":"` + user + `","password":"` + pass + `"}`,
	)

	req, err := http.NewRequest(
		http.MethodPost,
		"https://hub.docker.com/v2/users/login/",
		body,
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf(
			"Docker Hub auth failed: %s (%s)",
			resp.Status,
			strings.TrimSpace(string(b)),
		)
	}

	var result struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Token == "" {
		return "", errors.New("Docker Hub returned an empty token")
	}

	return result.Token, nil
}


func (Ops) DeleteImage() error {
	tag := os.Getenv("IMAGE_TAG")

	if tag == "" {
		log.Fatal(
			"Please provide a tag with -tag flag " +
				"(e.g., op deleteimage -tag PR32)",
		)
	}

	ctx := context.Background()
	token, err := fetchDockerToken()

	if err != nil {
		return err
	}

	imageTag := "kiloexabyte/runner-image:" + tag
	log.Printf("Deleting image: %s\n", imageTag)

	url := "https://hub.docker.com/v2/repositories/" +
		"kiloexabyte/runner-image/tags/" +
		tag + "/"

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		url,
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Docker Hub returns 204 No Content on success
	if resp.StatusCode != http.StatusNoContent {
		b, _ := io.ReadAll(resp.Body)
		log.Fatalf(
			"Failed to delete image tag %s: %s (%s)",
			tag,
			resp.Status,
			strings.TrimSpace(string(b)),
		)
	}

	log.Printf("Successfully deleted image tag: %s\n", tag)
	return nil
}
