package watcher

import (
	"testing"
	"time"
)

func TestNewPoller(t *testing.T) {
	t.Run("uses default interval when not specified", func(t *testing.T) {
		p := NewPoller(PollerConfig{Repo: "owner/repo"})
		if p.config.Interval != DefaultInterval {
			t.Errorf("expected interval %v, got %v", DefaultInterval, p.config.Interval)
		}
	})

	t.Run("uses custom interval when specified", func(t *testing.T) {
		customInterval := 30 * time.Second
		p := NewPoller(PollerConfig{Repo: "owner/repo", Interval: customInterval})
		if p.config.Interval != customInterval {
			t.Errorf("expected interval %v, got %v", customInterval, p.config.Interval)
		}
	})
}

func TestDefaultInterval(t *testing.T) {
	if DefaultInterval != 90*time.Second {
		t.Errorf("expected default interval 90s, got %v", DefaultInterval)
	}
}

func TestParseResponse(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedStatus string
		expectedETag   string
		expectedBody   string
	}{
		{
			name: "standard response with CRLF",
			input: "HTTP/2 200 OK\r\nETag: \"abc123\"\r\nContent-Type: application/json\r\n\r\n[{\"id\":\"1\"}]",
			expectedStatus: "HTTP/2 200 OK",
			expectedETag:   "\"abc123\"",
			expectedBody:   "[{\"id\":\"1\"}]",
		},
		{
			name: "response with LF only",
			input: "HTTP/2 200 OK\nETag: \"def456\"\n\n[{\"id\":\"2\"}]",
			expectedStatus: "HTTP/2 200 OK",
			expectedETag:   "\"def456\"",
			expectedBody:   "[{\"id\":\"2\"}]",
		},
		{
			name: "304 not modified",
			input: "HTTP/2 304 Not Modified\nETag: \"same\"\n\n",
			expectedStatus: "HTTP/2 304 Not Modified",
			expectedETag:   "\"same\"",
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers, body := parseResponse([]byte(tt.input))

			if headers["status"] != tt.expectedStatus {
				t.Errorf("expected status %q, got %q", tt.expectedStatus, headers["status"])
			}
			if headers["etag"] != tt.expectedETag {
				t.Errorf("expected etag %q, got %q", tt.expectedETag, headers["etag"])
			}
			if body != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, body)
			}
		})
	}
}

func TestFilterNewEvents(t *testing.T) {
	p := NewPoller(PollerConfig{Repo: "owner/repo"})

	t.Run("returns all events when no last seen ID", func(t *testing.T) {
		events := []Event{
			{ID: "3"},
			{ID: "2"},
			{ID: "1"},
		}
		result := p.filterNewEvents(events)
		if len(result) != 3 {
			t.Errorf("expected 3 events, got %d", len(result))
		}
	})

	t.Run("returns only new events after last seen", func(t *testing.T) {
		p.lastSeenID = "2"
		events := []Event{
			{ID: "4"},
			{ID: "3"},
			{ID: "2"},
			{ID: "1"},
		}
		result := p.filterNewEvents(events)
		if len(result) != 2 {
			t.Errorf("expected 2 events, got %d", len(result))
		}
		if result[0].ID != "4" || result[1].ID != "3" {
			t.Error("expected events 4 and 3")
		}
	})

	t.Run("returns empty when last seen is first", func(t *testing.T) {
		p.lastSeenID = "3"
		events := []Event{
			{ID: "3"},
			{ID: "2"},
			{ID: "1"},
		}
		result := p.filterNewEvents(events)
		if len(result) != 0 {
			t.Errorf("expected 0 events, got %d", len(result))
		}
	})

	t.Run("returns all when last seen not found", func(t *testing.T) {
		p.lastSeenID = "999"
		events := []Event{
			{ID: "3"},
			{ID: "2"},
			{ID: "1"},
		}
		result := p.filterNewEvents(events)
		if len(result) != 3 {
			t.Errorf("expected 3 events, got %d", len(result))
		}
	})
}

func TestPollerConfig(t *testing.T) {
	t.Run("can set repo", func(t *testing.T) {
		config := PollerConfig{Repo: "octocat/hello-world"}
		p := NewPoller(config)
		if p.config.Repo != "octocat/hello-world" {
			t.Errorf("expected repo octocat/hello-world, got %s", p.config.Repo)
		}
	})

	t.Run("can set custom interval", func(t *testing.T) {
		config := PollerConfig{Repo: "owner/repo", Interval: 60 * time.Second}
		p := NewPoller(config)
		if p.config.Interval != 60*time.Second {
			t.Errorf("expected interval 60s, got %v", p.config.Interval)
		}
	})
}
