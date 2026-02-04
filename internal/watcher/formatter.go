package watcher

import (
	"fmt"
)

// FormatEvent formats an event into a human-readable string
func FormatEvent(e Event) string {
	timestamp := e.CreatedAt.Format("2006-01-02 15:04:05")

	switch e.Type {
	case "IssuesEvent":
		return formatIssueEvent(timestamp, e)
	case "PullRequestEvent":
		return formatPREvent(timestamp, e)
	case "IssueCommentEvent":
		return formatIssueCommentEvent(timestamp, e)
	case "PullRequestReviewCommentEvent":
		return formatPRCommentEvent(timestamp, e)
	case "ReleaseEvent":
		return formatReleaseEvent(timestamp, e)
	default:
		return fmt.Sprintf("[%s] %s by @%s", timestamp, e.Type, e.Actor.Login)
	}
}

func formatIssueEvent(timestamp string, e Event) string {
	if e.Payload.Issue == nil {
		return fmt.Sprintf("[%s] Issue %s by @%s", timestamp, e.Payload.Action, e.Actor.Login)
	}
	issue := e.Payload.Issue
	result := fmt.Sprintf("[%s] Issue %s: #%d \"%s\" by @%s\n%s",
		timestamp, e.Payload.Action, issue.Number, issue.Title, e.Actor.Login, issue.HTMLURL)
	if issue.Body != "" {
		result += fmt.Sprintf("\n\n%s", issue.Body)
	}
	return result
}

func formatPREvent(timestamp string, e Event) string {
	if e.Payload.PullRequest == nil {
		return fmt.Sprintf("[%s] PR %s by @%s", timestamp, e.Payload.Action, e.Actor.Login)
	}

	pr := e.Payload.PullRequest
	action := e.Payload.Action
	if action == "closed" && pr.Merged {
		action = "merged"
	}

	result := fmt.Sprintf("[%s] PR %s: #%d \"%s\" by @%s\n%s",
		timestamp, action, pr.Number, pr.Title, e.Actor.Login, pr.HTMLURL)
	if pr.Body != "" {
		result += fmt.Sprintf("\n\n%s", pr.Body)
	}
	return result
}

func formatIssueCommentEvent(timestamp string, e Event) string {
	if e.Payload.Issue == nil {
		return fmt.Sprintf("[%s] Comment by @%s", timestamp, e.Actor.Login)
	}
	issue := e.Payload.Issue
	result := fmt.Sprintf("[%s] Comment on #%d \"%s\" by @%s\n%s",
		timestamp, issue.Number, issue.Title, e.Actor.Login, issue.HTMLURL)
	if e.Payload.Comment != nil && e.Payload.Comment.Body != "" {
		result += fmt.Sprintf("\n\n%s", e.Payload.Comment.Body)
	}
	return result
}

func formatPRCommentEvent(timestamp string, e Event) string {
	if e.Payload.PullRequest == nil {
		return fmt.Sprintf("[%s] PR comment by @%s", timestamp, e.Actor.Login)
	}
	pr := e.Payload.PullRequest
	result := fmt.Sprintf("[%s] PR comment on #%d \"%s\" by @%s\n%s",
		timestamp, pr.Number, pr.Title, e.Actor.Login, pr.HTMLURL)
	if e.Payload.Comment != nil && e.Payload.Comment.Body != "" {
		result += fmt.Sprintf("\n\n%s", e.Payload.Comment.Body)
	}
	return result
}

func formatReleaseEvent(timestamp string, e Event) string {
	if e.Payload.Release == nil {
		return fmt.Sprintf("[%s] Release %s by @%s", timestamp, e.Payload.Action, e.Actor.Login)
	}
	release := e.Payload.Release
	result := fmt.Sprintf("[%s] Release %s %s\n%s",
		timestamp, release.TagName, e.Payload.Action, release.HTMLURL)
	return result
}
