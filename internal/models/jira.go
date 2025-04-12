package models

import "time"

// JiraTicket represents a Jira ticket
type JiraTicket struct {
	ID          string        `json:"id"`
	Key         string        `json:"key"`
	Summary     string        `json:"summary"`
	Description string        `json:"description"`
	Created     time.Time     `json:"created"`
	Updated     time.Time     `json:"updated"`
	Status      string        `json:"status"`
	Comments    []JiraComment `json:"comments"`
}

// JiraComment represents a comment on a Jira ticket
type JiraComment struct {
	ID      string    `json:"id"`
	Body    string    `json:"body"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Author  JiraUser  `json:"author"`
}

// JiraUser represents a user in Jira
type JiraUser struct {
	AccountID    string `json:"accountId"`
	EmailAddress string `json:"emailAddress"`
	DisplayName  string `json:"displayName"`
	Active       bool   `json:"active"`
}
