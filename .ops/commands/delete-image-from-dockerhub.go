package commands

import (
	"bytes"
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

func fetchDockerToken(ctx context.Context) (string, error) {
	user := os.Getenv("DOCKER_USERNAME")
	pass := os.Getenv("DOCKER_PASSWORD")

	if user == "" || pass == "" {
		return "", errors.New(
			"docker username or docker password is not set",
		)
	}

	creds := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: user,
		Password: pass,
	}

	body, err := json.Marshal(creds)
	if err != nil {
		return "", fmt.Errorf("marshal credentials: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://hub.docker.com/v2/users/login/",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", fmt.Errorf("create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf(
			"docker hub auth failed: %s (%s)",
			resp.Status,
			strings.TrimSpace(string(b)),
		)
	}

	var result struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode docker hub response: %w", err)
	}

	if result.Token == "" {
		return "", errors.New("docker hub returned an empty token")
	}

	return result.Token, nil
}

func (Ops) DeleteImage() error {
	tag := os.Getenv("IMAGE_TAG")

	if tag == "" {
		return fmt.Errorf(
			"please provide a tag with -tag flag " +
				"(e.g., op deleteimage -tag PR32)",
		)
	}

	ctx := context.Background()
	token, err := fetchDockerToken(ctx)

	if err != nil {
		return fmt.Errorf("fetch docker token: %w", err)
	}

	imageTag := "kiloexabyte/runner-image:" + tag
	log.Printf("Deleting image: %s", imageTag)

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
		return fmt.Errorf("create delete request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("send delete request: %w", err)
	}
	defer resp.Body.Close()

	// Docker Hub returns 204 No Content on success
	if resp.StatusCode != http.StatusNoContent {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf(
			"delete image tag %s: %s (%s)",
			tag,
			resp.Status,
			strings.TrimSpace(string(b)),
		)
	}

	log.Printf("Successfully deleted image tag: %s", tag)
	return nil
}
