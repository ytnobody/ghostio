package watcher

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// Fetcher fetches events from GitHub API using gh command
type Fetcher struct {
	repo string
}

// NewFetcher creates a new Fetcher for the given repository
func NewFetcher(repo string) *Fetcher {
	return &Fetcher{repo: repo}
}

// FetchEvents retrieves recent events from the repository
func (f *Fetcher) FetchEvents() ([]Event, error) {
	cmd := exec.Command("gh", "api", fmt.Sprintf("repos/%s/events", f.repo))
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("gh api failed: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("failed to execute gh command: %w", err)
	}

	var events []Event
	if err := json.Unmarshal(output, &events); err != nil {
		return nil, fmt.Errorf("failed to parse events: %w", err)
	}

	return events, nil
}

// FetchTargetEvents retrieves only the events we want to monitor
func (f *Fetcher) FetchTargetEvents() ([]Event, error) {
	events, err := f.FetchEvents()
	if err != nil {
		return nil, err
	}

	var targetEvents []Event
	for _, e := range events {
		if IsTargetEvent(e.Type) {
			targetEvents = append(targetEvents, e)
		}
	}

	return targetEvents, nil
}
