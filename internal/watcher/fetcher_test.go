package watcher

import (
	"testing"
)

func TestIsTargetEvent(t *testing.T) {
	tests := []struct {
		eventType string
		expected  bool
	}{
		{"IssuesEvent", true},
		{"PullRequestEvent", true},
		{"IssueCommentEvent", true},
		{"PullRequestReviewCommentEvent", true},
		{"ReleaseEvent", true},
		{"PushEvent", false},
		{"CreateEvent", false},
		{"WatchEvent", false},
		{"ForkEvent", false},
	}

	for _, tt := range tests {
		t.Run(tt.eventType, func(t *testing.T) {
			if got := IsTargetEvent(tt.eventType); got != tt.expected {
				t.Errorf("IsTargetEvent(%q) = %v, want %v", tt.eventType, got, tt.expected)
			}
		})
	}
}

func TestNewFetcher(t *testing.T) {
	repo := "owner/repo"
	f := NewFetcher(repo)
	if f.repo != repo {
		t.Errorf("NewFetcher(%q).repo = %q, want %q", repo, f.repo, repo)
	}
}
