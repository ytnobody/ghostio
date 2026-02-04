package formatter

import (
	"strings"
	"testing"
	"time"

	"github.com/ytnobody/ghostio/internal/github"
)

func TestFormat_IssuesEvent(t *testing.T) {
	event := &github.Event{
		Type: "IssuesEvent",
		Actor: github.Actor{
			Login: "testuser",
		},
		Repo: github.Repo{
			FullName: "owner/repo",
		},
		Payload: github.Payload{
			Action: "opened",
			Issue: &github.Issue{
				Number:  42,
				Title:   "Bug in login",
				Body:    "Login fails with error 500",
				HTMLURL: "https://github.com/owner/repo/issues/42",
			},
		},
		CreatedAt: time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC),
	}

	result := Format(event)

	checks := []string{
		"[2024-01-15 10:30:45]",
		"IssuesEvent (opened)",
		"@testuser",
		"#42: Bug in login",
		"https://github.com/owner/repo/issues/42",
		"Login fails with error 500",
	}

	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf("expected output to contain %q, got:\n%s", check, result)
		}
	}
}

func TestFormat_PullRequestEvent(t *testing.T) {
	event := &github.Event{
		Type: "PullRequestEvent",
		Actor: github.Actor{
			Login: "pruser",
		},
		Payload: github.Payload{
			Action: "opened",
			PullRequest: &github.PullRequest{
				Number:  123,
				Title:   "Add new feature",
				Body:    "This PR adds a cool feature",
				HTMLURL: "https://github.com/owner/repo/pull/123",
			},
		},
		CreatedAt: time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC),
	}

	result := Format(event)

	checks := []string{
		"PullRequestEvent (opened)",
		"#123: Add new feature",
		"https://github.com/owner/repo/pull/123",
		"This PR adds a cool feature",
	}

	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf("expected output to contain %q, got:\n%s", check, result)
		}
	}
}

func TestFormat_IssueCommentEvent(t *testing.T) {
	event := &github.Event{
		Type: "IssueCommentEvent",
		Actor: github.Actor{
			Login: "commenter",
		},
		Payload: github.Payload{
			Action: "created",
			Issue: &github.Issue{
				Number: 42,
				Title:  "Bug in login",
			},
			Comment: &github.Comment{
				Body:    "I can reproduce this issue",
				HTMLURL: "https://github.com/owner/repo/issues/42#issuecomment-1",
			},
		},
		CreatedAt: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
	}

	result := Format(event)

	checks := []string{
		"IssueCommentEvent (created)",
		"on #42: Bug in login",
		"https://github.com/owner/repo/issues/42#issuecomment-1",
		"I can reproduce this issue",
	}

	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf("expected output to contain %q, got:\n%s", check, result)
		}
	}
}
