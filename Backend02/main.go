package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Event struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	Payload struct {
		Commits []struct {
			Message string `json:"message"`
		} `json:"commits"`
	} `json:"payload"`
}

func main() {
	// Check if a username is provided as an argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: github-activity <username>")
		os.Exit(1)
	}

	username := os.Args[1]
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)

	// Make HTTP GET request to GitHub API
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: GitHub API returned status %d\n", resp.StatusCode)
		os.Exit(1)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		os.Exit(1)
	}

	// Parse JSON into a slice of Event structs
	var events []Event
	err = json.Unmarshal(body, &events)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	// Display the recent activity
	for _, event := range events {
		switch event.Type {
		case "PushEvent":
			commitCount := len(event.Payload.Commits)
			if commitCount > 0 {
				fmt.Printf("Pushed %d commits to %s\n", commitCount, event.Repo.Name)
			}
		case "IssuesEvent":
			fmt.Printf("Opened a new issue in %s\n", event.Repo.Name)
		case "WatchEvent":
			fmt.Printf("Starred %s\n", event.Repo.Name)
		default:
			fmt.Printf("Performed %s on %s\n", event.Type, event.Repo.Name)
		}
	}
}
