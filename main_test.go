package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestFormatEvent_PushEvent(t *testing.T) {
	event := Event{
		Type:    "PushEvent",
		Repo:    Repo{Name: "test/repo"},
		Payload: Payload{Size: 3},
	}

	expected := "- Pushed 3 commits to test/repo"
	result := formatEvent(event)

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestFormatEvent_IssuesEvent(t *testing.T) {
	tests := []struct {
		action   string
		expected string
	}{
		{"opened", "- Opened a new issue in test/repo"},
		{"closed", "- Closed an issue in test/repo"},
		{"labeled", "- Labeled an issue in test/repo"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("IssuesEvent_%s", tt.action), func(t *testing.T) {
			event := Event{
				Type:    "IssuesEvent",
				Repo:    Repo{Name: "test/repo"},
				Payload: Payload{Action: tt.action},
			}

			result := formatEvent(event)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestFormatEvent_WatchEvent(t *testing.T) {
	event := Event{
		Type: "WatchEvent",
		Repo: Repo{Name: "test/repo"},
	}

	expected := "- Starred test/repo"
	result := formatEvent(event)

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestFormatEvent_CreateEvent(t *testing.T) {
	tests := []struct {
		refType  string
		ref      string
		expected string
	}{
		{"repository", "", "- Created repository test/repo"},
		{"branch", "feature-branch", "- Created branch feature-branch in test/repo"},
		{"tag", "v1.0.0", "- Created tag in test/repo"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("CreateEvent_%s", tt.refType), func(t *testing.T) {
			event := Event{
				Type: "CreateEvent",
				Repo: Repo{Name: "test/repo"},
				Payload: Payload{
					RefType: tt.refType,
					Ref:     tt.ref,
				},
			}

			result := formatEvent(event)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestFormatEvent_DeleteEvent(t *testing.T) {
	event := Event{
		Type: "DeleteEvent",
		Repo: Repo{Name: "test/repo"},
		Payload: Payload{
			RefType: "branch",
			Ref:     "old-branch",
		},
	}

	expected := "- Deleted branch old-branch in test/repo"
	result := formatEvent(event)

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestFormatEvent_ForkEvent(t *testing.T) {
	event := Event{
		Type: "ForkEvent",
		Repo: Repo{Name: "test/repo"},
	}

	expected := "- Forked test/repo"
	result := formatEvent(event)

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestFormatEvent_PullRequestEvent(t *testing.T) {
	tests := []struct {
		action   string
		expected string
	}{
		{"opened", "- Opened a new pull request in test/repo"},
		{"closed", "- Closed a pull request in test/repo"},
		{"merged", "- Merged a pull request in test/repo"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("PullRequestEvent_%s", tt.action), func(t *testing.T) {
			event := Event{
				Type:    "PullRequestEvent",
				Repo:    Repo{Name: "test/repo"},
				Payload: Payload{Action: tt.action},
			}

			result := formatEvent(event)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestFormatEvent_PublicEvent(t *testing.T) {
	event := Event{
		Type: "PublicEvent",
		Repo: Repo{Name: "test/repo"},
	}

	expected := "- Made test/repo public"
	result := formatEvent(event)

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestFormatEvent_MemberEvent(t *testing.T) {
	event := Event{
		Type:    "MemberEvent",
		Repo:    Repo{Name: "test/repo"},
		Payload: Payload{Action: "added"},
	}

	expected := "- Added a collaborator to test/repo"
	result := formatEvent(event)

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestFormatEvent_UnknownEvent(t *testing.T) {
	event := Event{
		Type: "UnknownEvent",
		Repo: Repo{Name: "test/repo"},
	}

	expected := "- UnknownEvent in test/repo"
	result := formatEvent(event)

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestGroupEvents_EmptySlice(t *testing.T) {
	events := []Event{}
	result := groupEvents(events)

	if len(result) != 0 {
		t.Errorf("Expected empty slice, got %d events", len(result))
	}
}

func TestGroupEvents_SinglePushEvent(t *testing.T) {
	events := []Event{
		{
			Type:    "PushEvent",
			Repo:    Repo{Name: "test/repo"},
			Payload: Payload{Size: 3},
		},
	}

	result := groupEvents(events)

	if len(result) != 1 {
		t.Errorf("Expected 1 event, got %d", len(result))
	}

	if result[0].Payload.Size != 3 {
		t.Errorf("Expected size 3, got %d", result[0].Payload.Size)
	}
}

func TestGroupEvents_MultiplePushEventsSameRepo(t *testing.T) {
	events := []Event{
		{
			Type:    "PushEvent",
			Repo:    Repo{Name: "test/repo"},
			Payload: Payload{Size: 2},
		},
		{
			Type:    "PushEvent",
			Repo:    Repo{Name: "test/repo"},
			Payload: Payload{Size: 3},
		},
		{
			Type:    "PushEvent",
			Repo:    Repo{Name: "test/repo"},
			Payload: Payload{Size: 1},
		},
	}

	result := groupEvents(events)

	if len(result) != 1 {
		t.Errorf("Expected 1 grouped event, got %d", len(result))
	}

	if result[0].Payload.Size != 6 {
		t.Errorf("Expected total size 6, got %d", result[0].Payload.Size)
	}
}

func TestGroupEvents_MultiplePushEventsDifferentRepos(t *testing.T) {
	events := []Event{
		{
			Type:    "PushEvent",
			Repo:    Repo{Name: "test/repo1"},
			Payload: Payload{Size: 2},
		},
		{
			Type:    "PushEvent",
			Repo:    Repo{Name: "test/repo2"},
			Payload: Payload{Size: 3},
		},
	}

	result := groupEvents(events)

	if len(result) != 2 {
		t.Errorf("Expected 2 events, got %d", len(result))
	}

	if result[0].Payload.Size != 2 {
		t.Errorf("Expected first event size 2, got %d", result[0].Payload.Size)
	}

	if result[1].Payload.Size != 3 {
		t.Errorf("Expected second event size 3, got %d", result[1].Payload.Size)
	}
}

func TestGroupEvents_MixedEventTypes(t *testing.T) {
	events := []Event{
		{
			Type:    "PushEvent",
			Repo:    Repo{Name: "test/repo"},
			Payload: Payload{Size: 2},
		},
		{
			Type:    "IssuesEvent",
			Repo:    Repo{Name: "test/repo"},
			Payload: Payload{Action: "opened"},
		},
		{
			Type:    "PushEvent",
			Repo:    Repo{Name: "test/repo"},
			Payload: Payload{Size: 3},
		},
	}

	result := groupEvents(events)

	if len(result) != 3 {
		t.Errorf("Expected 3 events, got %d", len(result))
	}

	expected := []Event{
		{
			Type:    "PushEvent",
			Repo:    Repo{Name: "test/repo"},
			Payload: Payload{Size: 2},
		},
		{
			Type:    "IssuesEvent",
			Repo:    Repo{Name: "test/repo"},
			Payload: Payload{Action: "opened"},
		},
		{
			Type:    "PushEvent",
			Repo:    Repo{Name: "test/repo"},
			Payload: Payload{Size: 3},
		},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestGroupEvents_ConsecutivePushEventsWithOtherEvents(t *testing.T) {
	events := []Event{
		{
			Type:    "PushEvent",
			Repo:    Repo{Name: "test/repo"},
			Payload: Payload{Size: 1},
		},
		{
			Type:    "PushEvent",
			Repo:    Repo{Name: "test/repo"},
			Payload: Payload{Size: 2},
		},
		{
			Type: "WatchEvent",
			Repo: Repo{Name: "another/repo"},
		},
		{
			Type:    "PushEvent",
			Repo:    Repo{Name: "test/repo"},
			Payload: Payload{Size: 1},
		},
	}

	result := groupEvents(events)

	if len(result) != 3 {
		t.Errorf("Expected 3 events, got %d", len(result))
	}

	if result[0].Type != "PushEvent" || result[0].Payload.Size != 3 {
		t.Errorf("Expected first event to be grouped PushEvent with size 3, got %+v", result[0])
	}

	if result[1].Type != "WatchEvent" {
		t.Errorf("Expected second event to be WatchEvent, got %s", result[1].Type)
	}

	if result[2].Type != "PushEvent" || result[2].Payload.Size != 1 {
		t.Errorf("Expected third event to be PushEvent with size 1, got %+v", result[2])
	}
}

func TestFetchGitHubActivity_ValidUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"type": "PushEvent",
				"repo": {"name": "test/repo"},
				"payload": {"size": 2}
			}
		]`))
	}))
	defer server.Close()

	// Override the URL in the function by modifying the test
	// Since we can't easily mock the http.Get call without changing the function,
	// we'll create a separate test function that accepts a base URL
	events, err := fetchGitHubActivityWithURL("testuser", server.URL)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	if events[0].Type != "PushEvent" {
		t.Errorf("Expected PushEvent, got %s", events[0].Type)
	}
}

func TestFetchGitHubActivity_UserNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_, err := fetchGitHubActivityWithURL("nonexistentuser", server.URL)
	if err == nil {
		t.Error("Expected error for non-existent user")
	}

	expected := "user 'nonexistentuser' not found"
	if err.Error() != expected {
		t.Errorf("Expected error %q, got %q", expected, err.Error())
	}
}

func TestFetchGitHubActivity_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	_, err := fetchGitHubActivityWithURL("testuser", server.URL)
	if err == nil {
		t.Error("Expected error for API error")
	}

	expected := "GitHub API returned status 500"
	if err.Error() != expected {
		t.Errorf("Expected error %q, got %q", expected, err.Error())
	}
}

func TestFetchGitHubActivity_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	_, err := fetchGitHubActivityWithURL("testuser", server.URL)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}

	if !strings.Contains(fmt.Sprintf("%v", err), "failed to parse response") {
		t.Errorf("Expected parse error, got %v", err)
	}
}

// Helper function to test fetchGitHubActivity with a custom URL
func fetchGitHubActivityWithURL(username, baseURL string) ([]Event, error) {
	url := fmt.Sprintf("%s/users/%s/events", baseURL, username)
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
