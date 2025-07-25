package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Event struct {
	Type    string  `json:"type"`
	Repo    Repo    `json:"repo"`
	Payload Payload `json:"payload"`
}

type Repo struct {
	Name string `json:"name"`
}

type Payload struct {
	Ref     string `json:"ref"`
	RefType string `json:"ref_type"`
	Size    int    `json:"size"`
	Action  string `json:"action"`
}

// createHyperlink creates a clickable hyperlink using ANSI escape codes
func createHyperlink(url, text string) string {
	return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", url, text)
}

// fetchGitHubActivity retrieves the recent public GitHub events for the specified user.
// It returns a slice of Event structs representing the user's activity, or an error if the request fails.
func fetchGitHubActivity(username string) ([]Event, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("user '%s' not found", username)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return events, nil
}

// formatEvent returns a human-readable string describing a GitHub event.
// It formats the event based on its type and relevant payload details.
func formatEvent(event Event) string {
	repoName := event.Repo.Name
	repoURL := fmt.Sprintf("https://github.com/%s", repoName)
	linkedRepoName := createHyperlink(repoURL, repoName)

	switch event.Type {
	case "PushEvent":
		commitCount := event.Payload.Size
		return fmt.Sprintf("- Pushed %d commits to %s", commitCount, linkedRepoName)

	case "IssuesEvent":
		action := event.Payload.Action
		switch action {
		case "opened":
			return fmt.Sprintf("- Opened a new issue in %s", linkedRepoName)
		case "closed":
			return fmt.Sprintf("- Closed an issue in %s", linkedRepoName)
		}
		return fmt.Sprintf("- %s an issue in %s", strings.Title(action), linkedRepoName)

	case "WatchEvent":
		return fmt.Sprintf("- Starred %s", linkedRepoName)

	case "CreateEvent":
		refType := event.Payload.RefType
		switch refType {
		case "repository":
			return fmt.Sprintf("- Created repository %s", linkedRepoName)
		case "branch":
			ref := event.Payload.Ref
			return fmt.Sprintf("- Created branch %s in %s", ref, linkedRepoName)
		}
		return fmt.Sprintf("- Created %s in %s", refType, linkedRepoName)

	case "DeleteEvent":
		refType := event.Payload.RefType
		ref := event.Payload.Ref
		return fmt.Sprintf("- Deleted %s %s in %s", refType, ref, linkedRepoName)

	case "ForkEvent":
		return fmt.Sprintf("- Forked %s", linkedRepoName)

	case "PullRequestEvent":
		action := event.Payload.Action
		switch action {
		case "opened":
			return fmt.Sprintf("- Opened a new pull request in %s", linkedRepoName)
		case "closed":
			return fmt.Sprintf("- Closed a pull request in %s", linkedRepoName)
		}
		return fmt.Sprintf("- %s a pull request in %s", strings.Title(action), linkedRepoName)

	case "PublicEvent":
		return fmt.Sprintf("- Made %s public", linkedRepoName)

	case "MemberEvent":
		action := event.Payload.Action
		return fmt.Sprintf("- %s a collaborator to %s", strings.Title(action), linkedRepoName)

	default:
		return fmt.Sprintf("- %s in %s", event.Type, linkedRepoName)
	}
}

// groupEvents aggregates consecutive PushEvents for the same repository into a single event
// with the total number of commits. Other events are left unchanged.
// It returns a new slice of grouped events.
func groupEvents(events []Event) []Event {
	if len(events) == 0 {
		return events
	}

	var grouped []Event
	i := 0

	for i < len(events) {
		currentEvent := events[i]

		if currentEvent.Type == "PushEvent" {
			totalCommits := currentEvent.Payload.Size
			j := i + 1

			for j < len(events) && events[j].Type == "PushEvent" && events[j].Repo.Name == currentEvent.Repo.Name {
				totalCommits += events[j].Payload.Size
				j++
			}

			groupedEvent := currentEvent
			groupedEvent.Payload.Size = totalCommits
			grouped = append(grouped, groupedEvent)
			i = j
		} else {
			grouped = append(grouped, currentEvent)
			i++
		}
	}

	return grouped
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: github-activity <username>")
		fmt.Println("Example: github-activity nikitalobanov12")
		os.Exit(1)
	}

	username := os.Args[1]

	if strings.TrimSpace(username) == "" {
		fmt.Println("Error: username cannot be empty")
		os.Exit(1)
	}

	if username == "-h" || username == "--help" {
		fmt.Println("GitHub Activity CLI")
		fmt.Println("Usage: github-activity <username>")
		fmt.Println("Example: github-activity nikitalobanov12")
		os.Exit(0)
	}

	userURL := fmt.Sprintf("https://github.com/%s", username)
	linkedUsername := createHyperlink(userURL, username)
	fmt.Printf("Fetching activity for user: %s\n", linkedUsername)

	events, err := fetchGitHubActivity(username)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if len(events) == 0 {
		fmt.Printf("No recent activity found for user: %s\n", username)
		return
	}

	groupedEvents := groupEvents(events)

	fmt.Printf("\nRecent activity for %s:\n", linkedUsername)
	for i, event := range groupedEvents {
		if i >= 10 {
			break
		}
		fmt.Println(formatEvent(event))
	}
}
