package watcher

import (
	"slices"
	"time"
)

// Event represents a GitHub repository event
type Event struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Actor     Actor     `json:"actor"`
	Repo      Repo      `json:"repo"`
	Payload   Payload   `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

// Actor represents the user who triggered the event
type Actor struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

// Repo represents the repository where the event occurred
type Repo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Payload contains event-specific data
type Payload struct {
	Action string `json:"action"`

	// Issue events
	Issue *Issue `json:"issue,omitempty"`

	// PR events
	PullRequest *PullRequest `json:"pull_request,omitempty"`

	// Comment events
	Comment *Comment `json:"comment,omitempty"`

	// Release events
	Release *Release `json:"release,omitempty"`
}

// Issue represents a GitHub issue
type Issue struct {
	Number  int    `json:"number"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	State   string `json:"state"`
	HTMLURL string `json:"html_url"`
}

// PullRequest represents a GitHub pull request
type PullRequest struct {
	Number  int    `json:"number"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	State   string `json:"state"`
	HTMLURL string `json:"html_url"`
	Merged  bool   `json:"merged"`
}

// Comment represents a GitHub comment
type Comment struct {
	ID      int    `json:"id"`
	Body    string `json:"body"`
	HTMLURL string `json:"html_url"`
}

// Release represents a GitHub release
type Release struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	HTMLURL string `json:"html_url"`
}

// TargetEventTypes are the event types we want to monitor
var TargetEventTypes = []string{
	"IssuesEvent",
	"PullRequestEvent",
	"IssueCommentEvent",
	"PullRequestReviewCommentEvent",
	"ReleaseEvent",
}

// IsTargetEvent checks if the event type is one we want to monitor
func IsTargetEvent(eventType string) bool {
	return slices.Contains(TargetEventTypes, eventType)
}
