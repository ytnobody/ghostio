package formatter

import (
	"fmt"
	"strings"

	"github.com/ytnobody/ghostio/internal/github"
)

func Format(event *github.Event) string {
	var b strings.Builder

	timestamp := event.CreatedAt.Format("2006-01-02 15:04:05")
	action := event.Payload.Action
	if action == "" {
		action = "unknown"
	}

	b.WriteString(fmt.Sprintf("[%s] %s (%s) by @%s\n",
		timestamp, event.Type, action, event.Actor.Login))

	switch event.Type {
	case "IssuesEvent":
		if issue := event.Payload.Issue; issue != nil {
			b.WriteString(fmt.Sprintf("#%d: %s\n", issue.Number, issue.Title))
			b.WriteString(fmt.Sprintf("%s\n", issue.HTMLURL))
			if issue.Body != "" {
				b.WriteString(fmt.Sprintf("\n%s\n", issue.Body))
			}
		}

	case "PullRequestEvent":
		if pr := event.Payload.PullRequest; pr != nil {
			b.WriteString(fmt.Sprintf("#%d: %s\n", pr.Number, pr.Title))
			b.WriteString(fmt.Sprintf("%s\n", pr.HTMLURL))
			if pr.Body != "" {
				b.WriteString(fmt.Sprintf("\n%s\n", pr.Body))
			}
		}

	case "IssueCommentEvent", "PullRequestReviewCommentEvent":
		if issue := event.Payload.Issue; issue != nil {
			b.WriteString(fmt.Sprintf("on #%d: %s\n", issue.Number, issue.Title))
		}
		if pr := event.Payload.PullRequest; pr != nil {
			b.WriteString(fmt.Sprintf("on #%d: %s\n", pr.Number, pr.Title))
		}
		if comment := event.Payload.Comment; comment != nil {
			b.WriteString(fmt.Sprintf("%s\n", comment.HTMLURL))
			if comment.Body != "" {
				b.WriteString(fmt.Sprintf("\n%s\n", comment.Body))
			}
		}
	}

	return b.String()
}
