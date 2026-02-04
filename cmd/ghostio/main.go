package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ytnobody/ghostio/internal/watcher"
)

func main() {
	if len(os.Args) < 3 || os.Args[1] != "watch" {
		fmt.Fprintln(os.Stderr, "Usage: ghostio watch owner/repo")
		os.Exit(1)
	}

	repo := os.Args[2]
	fmt.Fprintf(os.Stderr, "Watching %s...\n", repo)

	// Setup context with signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Fprintln(os.Stderr, "\nShutting down...")
		cancel()
	}()

	// Setup poller
	poller := watcher.NewPoller(watcher.PollerConfig{
		Repo: repo,
	})

	// Event channel
	eventCh := make(chan watcher.Event)

	// Start event consumer
	go func() {
		for event := range eventCh {
			if watcher.IsTargetEvent(event.Type) {
				fmt.Println(watcher.FormatEvent(event))
				fmt.Println("---")
			}
		}
	}()

	// Start polling
	if err := poller.Start(ctx, eventCh); err != nil {
		if err == context.Canceled {
			return
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
