package models

import (
	"fmt"
	"time"
)

type SayoriWrapperResponse struct {
	Data SayoriResponses `json:"data"`
}

// SayoriResponses represents an array of responses from the Sayori API
type SayoriResponses []SayoriResponse

// SayoriResponse represents a single ticket from the Sayori API
type SayoriResponse struct {
	ID             int        `json:"id"`
	Link           string     `json:"link"`
	Estimate       float64    `json:"estimate"`
	Difficulty     int        `json:"dificulty"` // Note: API returns "dificulty" (misspelled)
	TicketID       string     `json:"ticket"`
	BotUserID      string     `json:"bot_user_id"`
	TicketType     int        `json:"ticket_type"`
	CurrentProject string     `json:"current_project"`
	Comments       string     `json:"comments"`
	CreateDate     CustomTime `json:"create_date"`
	TicketDate     CustomTime `json:"ticket_date"`
	BotUser        BotUser    `json:"BotUser"`
}

// BotUser represents a user in the Sayori API response
type BotUser struct {
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	Project     string `json:"project"`
	Picture     string `json:"picture"`
	DisplayName string `json:"display_name"`
}

type CustomTime struct {
	time.Time
}

// UnmarshalJSON parses the custom timestamp format
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	// Remove quotes from the JSON string
	str := string(b)
	str = str[1 : len(str)-1]

	// Parse the timestamp
	parsedTime, err := time.Parse("2006-01-02T15:04:05.999999", str)
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}

	ct.Time = parsedTime
	return nil
}
