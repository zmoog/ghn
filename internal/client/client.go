package client

import (
	"fmt"
	"os"

	"github.com/google/go-github/v69/github"
)

// NewClient creates a new GitHub client
func NewClient() (*github.Client, error) {
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN is not set")
	}

	return github.NewClient(nil).WithAuthToken(githubToken), nil
}
