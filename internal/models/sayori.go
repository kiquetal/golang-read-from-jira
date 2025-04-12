package models

import "time"

// SayoriResponses represents an array of responses from the Sayori API
type SayoriResponses []SayoriResponse

// SayoriResponse represents a single ticket from the Sayori API
type SayoriResponse struct {
	ID             int       `json:"id"`
	Link           string    `json:"link"`
	Estimate       float64   `json:"estimate"`
	Difficulty     int       `json:"dificulty"` // Note: API returns "dificulty" (misspelled)
	TicketID       string    `json:"ticket"`
	BotUserID      string    `json:"bot_user_id"`
	TicketType     int       `json:"ticket_type"`
	CurrentProject string    `json:"current_project"`
	Comments       string    `json:"comments"`
	CreateDate     time.Time `json:"create_date"`
	TicketDate     time.Time `json:"ticket_date"`
	BotUser        BotUser   `json:"BotUser"`
}

// BotUser represents a user in the Sayori API response
type BotUser struct {
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	Project     string `json:"project"`
	Picture     string `json:"picture"`
	DisplayName string `json:"display_name"`
}
