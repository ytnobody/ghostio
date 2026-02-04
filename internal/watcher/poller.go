package watcher

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// PollerConfig holds configuration for the poller
type PollerConfig struct {
	Repo     string
	Interval time.Duration
}

// DefaultInterval is the default polling interval
const DefaultInterval = 90 * time.Second

// Poller polls GitHub events with ETag caching
type Poller struct {
	config    PollerConfig
	etag      string
	lastSeenID string
}

// NewPoller creates a new Poller with the given configuration
func NewPoller(config PollerConfig) *Poller {
	if config.Interval == 0 {
		config.Interval = DefaultInterval
	}
	return &Poller{config: config}
}

// PollResult represents the result of a poll
type PollResult struct {
	Events     []Event
	NotModified bool
	Error      error
}

// Start begins polling and sends events to the provided channel
func (p *Poller) Start(ctx context.Context, eventCh chan<- Event) error {
	ticker := time.NewTicker(p.config.Interval)
	defer ticker.Stop()

	// Initial fetch
	if err := p.poll(eventCh); err != nil {
		return fmt.Errorf("initial poll failed: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := p.poll(eventCh); err != nil {
				// Log error but continue polling
				fmt.Printf("poll error: %v\n", err)
			}
		}
	}
}

// poll fetches events and sends new ones to the channel
func (p *Poller) poll(eventCh chan<- Event) error {
	result := p.fetchWithETag()
	if result.Error != nil {
		return result.Error
	}

	if result.NotModified {
		return nil
	}

	// Send new events (in reverse order, oldest first)
	newEvents := p.filterNewEvents(result.Events)
	for i := len(newEvents) - 1; i >= 0; i-- {
		eventCh <- newEvents[i]
	}

	// Update last seen ID
	if len(result.Events) > 0 {
		p.lastSeenID = result.Events[0].ID
	}

	return nil
}

// fetchWithETag fetches events using ETag for caching
func (p *Poller) fetchWithETag() PollResult {
	args := []string{"api", fmt.Sprintf("repos/%s/events", p.config.Repo), "-i"}
	if p.etag != "" {
		args = append(args, "-H", fmt.Sprintf("If-None-Match: %s", p.etag))
	}

	cmd := exec.Command("gh", args...)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// Check if 304 Not Modified
			if strings.Contains(string(output), "304") || strings.Contains(string(exitErr.Stderr), "304") {
				return PollResult{NotModified: true}
			}
			return PollResult{Error: fmt.Errorf("gh api failed: %s", string(exitErr.Stderr))}
		}
		return PollResult{Error: fmt.Errorf("failed to execute gh command: %w", err)}
	}

	// Parse response (includes headers)
	headers, body := parseResponse(output)

	// Update ETag
	if etag := headers["etag"]; etag != "" {
		p.etag = etag
	}

	// Check for 304 in headers
	if strings.Contains(headers["status"], "304") {
		return PollResult{NotModified: true}
	}

	var events []Event
	if err := json.Unmarshal([]byte(body), &events); err != nil {
		return PollResult{Error: fmt.Errorf("failed to parse events: %w", err)}
	}

	return PollResult{Events: events}
}

// parseResponse splits headers and body from gh api -i output
func parseResponse(output []byte) (map[string]string, string) {
	parts := strings.SplitN(string(output), "\r\n\r\n", 2)
	if len(parts) < 2 {
		// Try with just \n\n
		parts = strings.SplitN(string(output), "\n\n", 2)
	}

	headers := make(map[string]string)
	if len(parts) >= 1 {
		for _, line := range strings.Split(parts[0], "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "HTTP/") {
				headers["status"] = line
				continue
			}
			if idx := strings.Index(line, ":"); idx > 0 {
				key := strings.ToLower(strings.TrimSpace(line[:idx]))
				value := strings.TrimSpace(line[idx+1:])
				headers[key] = value
			}
		}
	}

	body := ""
	if len(parts) >= 2 {
		body = parts[1]
	}

	return headers, body
}

// filterNewEvents returns only events newer than the last seen
func (p *Poller) filterNewEvents(events []Event) []Event {
	if p.lastSeenID == "" {
		return []Event{}
	}

	var newEvents []Event
	for _, e := range events {
		if e.ID == p.lastSeenID {
			break
		}
		newEvents = append(newEvents, e)
	}

	return newEvents
}

// FetchOnce does a single fetch (for testing/manual use)
func (p *Poller) FetchOnce() ([]Event, error) {
	result := p.fetchWithETag()
	if result.Error != nil {
		return nil, result.Error
	}
	if result.NotModified {
		return nil, nil
	}
	return p.filterNewEvents(result.Events), nil
}
