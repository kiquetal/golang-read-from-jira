package models

type DynamoDBTicket struct {
	Pk        string `json:"pk"`
	SK        string `json:"sk"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Comments  string `json:"comments"`
	Summary   string `json:"summary"`
}
