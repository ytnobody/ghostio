package github

import "time"

type Issue struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	HTMLURL   string    `json:"html_url"`
	User      User      `json:"user"`
	CreatedAt time.Time `json:"created_at"`
}

type PullRequest struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	HTMLURL   string    `json:"html_url"`
	User      User      `json:"user"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	Login string `json:"login"`
}

type Repo struct {
	FullName string `json:"full_name"`
}

type Event struct {
	Type      string       `json:"type"`
	Actor     Actor        `json:"actor"`
	Repo      Repo         `json:"repo"`
	Payload   Payload      `json:"payload"`
	CreatedAt time.Time    `json:"created_at"`
}

type Actor struct {
	Login string `json:"login"`
}

type Payload struct {
	Action      string       `json:"action,omitempty"`
	Issue       *Issue       `json:"issue,omitempty"`
	PullRequest *PullRequest `json:"pull_request,omitempty"`
	Comment     *Comment     `json:"comment,omitempty"`
}

type Comment struct {
	Body    string `json:"body"`
	HTMLURL string `json:"html_url"`
	User    User   `json:"user"`
}
