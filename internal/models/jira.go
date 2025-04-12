package models

import "time"

// JiraTicket represents a Jira ticket
type JiraTicket struct {
	Fields JiraFields `json:"fields"`
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
	Name         string `json:"name"`
	Active       bool   `json:"active"`
}

type JiraFields struct {
	Description string        `json:"description"`
	Summary     string        `json:"summary"`
	Comments    []JiraComment `json:"comment"`
}
