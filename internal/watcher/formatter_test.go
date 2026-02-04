package watcher

import (
	"testing"
	"time"
)

func TestFormatEvent(t *testing.T) {
	baseTime := time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC)

	tests := []struct {
		name     string
		event    Event
		expected string
	}{
		{
			name: "Issue opened with body",
			event: Event{
				Type:  "IssuesEvent",
				Actor: Actor{Login: "username"},
				Payload: Payload{
					Action: "opened",
					Issue: &Issue{
						Number:  42,
						Title:   "Bug in login",
						Body:    "Login fails with error 500",
						HTMLURL: "https://github.com/owner/repo/issues/42",
					},
				},
				CreatedAt: baseTime,
			},
			expected: "[2024-01-15 10:30:45] Issue opened: #42 \"Bug in login\" by @username\nhttps://github.com/owner/repo/issues/42\n\nLogin fails with error 500",
		},
		{
			name: "Issue closed without body",
			event: Event{
				Type:  "IssuesEvent",
				Actor: Actor{Login: "user2"},
				Payload: Payload{
					Action: "closed",
					Issue: &Issue{
						Number:  10,
						Title:   "Fix needed",
						HTMLURL: "https://github.com/owner/repo/issues/10",
					},
				},
				CreatedAt: baseTime,
			},
			expected: "[2024-01-15 10:30:45] Issue closed: #10 \"Fix needed\" by @user2\nhttps://github.com/owner/repo/issues/10",
		},
		{
			name: "PR opened with body",
			event: Event{
				Type:  "PullRequestEvent",
				Actor: Actor{Login: "developer"},
				Payload: Payload{
					Action: "opened",
					PullRequest: &PullRequest{
						Number:  99,
						Title:   "Add feature",
						Body:    "This PR adds a new feature",
						HTMLURL: "https://github.com/owner/repo/pull/99",
					},
				},
				CreatedAt: baseTime,
			},
			expected: "[2024-01-15 10:30:45] PR opened: #99 \"Add feature\" by @developer\nhttps://github.com/owner/repo/pull/99\n\nThis PR adds a new feature",
		},
		{
			name: "PR merged",
			event: Event{
				Type:  "PullRequestEvent",
				Actor: Actor{Login: "username"},
				Payload: Payload{
					Action: "closed",
					PullRequest: &PullRequest{
						Number:  55,
						Title:   "Fix typo",
						Merged:  true,
						HTMLURL: "https://github.com/owner/repo/pull/55",
					},
				},
				CreatedAt: baseTime.Add(35 * time.Second),
			},
			expected: "[2024-01-15 10:31:20] PR merged: #55 \"Fix typo\" by @username\nhttps://github.com/owner/repo/pull/55",
		},
		{
			name: "PR closed without merge",
			event: Event{
				Type:  "PullRequestEvent",
				Actor: Actor{Login: "reviewer"},
				Payload: Payload{
					Action: "closed",
					PullRequest: &PullRequest{
						Number:  77,
						Title:   "WIP",
						Merged:  false,
						HTMLURL: "https://github.com/owner/repo/pull/77",
					},
				},
				CreatedAt: baseTime,
			},
			expected: "[2024-01-15 10:30:45] PR closed: #77 \"WIP\" by @reviewer\nhttps://github.com/owner/repo/pull/77",
		},
		{
			name: "Issue comment with body",
			event: Event{
				Type:  "IssueCommentEvent",
				Actor: Actor{Login: "username"},
				Payload: Payload{
					Action: "created",
					Issue: &Issue{
						Number:  123,
						Title:   "Some issue",
						HTMLURL: "https://github.com/owner/repo/issues/123",
					},
					Comment: &Comment{
						Body: "This is my comment",
					},
				},
				CreatedAt: baseTime.Add(75 * time.Second),
			},
			expected: "[2024-01-15 10:32:00] Comment on #123 \"Some issue\" by @username\nhttps://github.com/owner/repo/issues/123\n\nThis is my comment",
		},
		{
			name: "PR review comment with body",
			event: Event{
				Type:  "PullRequestReviewCommentEvent",
				Actor: Actor{Login: "reviewer"},
				Payload: Payload{
					Action: "created",
					PullRequest: &PullRequest{
						Number:  456,
						Title:   "Some PR",
						HTMLURL: "https://github.com/owner/repo/pull/456",
					},
					Comment: &Comment{
						Body: "LGTM",
					},
				},
				CreatedAt: baseTime,
			},
			expected: "[2024-01-15 10:30:45] PR comment on #456 \"Some PR\" by @reviewer\nhttps://github.com/owner/repo/pull/456\n\nLGTM",
		},
		{
			name: "Release published",
			event: Event{
				Type:  "ReleaseEvent",
				Actor: Actor{Login: "releaser"},
				Payload: Payload{
					Action: "published",
					Release: &Release{
						TagName: "v1.0.0",
						HTMLURL: "https://github.com/owner/repo/releases/tag/v1.0.0",
					},
				},
				CreatedAt: baseTime.Add(135 * time.Second),
			},
			expected: "[2024-01-15 10:33:00] Release v1.0.0 published\nhttps://github.com/owner/repo/releases/tag/v1.0.0",
		},
		{
			name: "Issue event without issue data",
			event: Event{
				Type:      "IssuesEvent",
				Actor:     Actor{Login: "user"},
				Payload:   Payload{Action: "opened"},
				CreatedAt: baseTime,
			},
			expected: `[2024-01-15 10:30:45] Issue opened by @user`,
		},
		{
			name: "Unknown event type",
			event: Event{
				Type:      "WatchEvent",
				Actor:     Actor{Login: "watcher"},
				CreatedAt: baseTime,
			},
			expected: `[2024-01-15 10:30:45] WatchEvent by @watcher`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatEvent(tt.event)
			if result != tt.expected {
				t.Errorf("FormatEvent() = %q, want %q", result, tt.expected)
			}
		})
	}
}
