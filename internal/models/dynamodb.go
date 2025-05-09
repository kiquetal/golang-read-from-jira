package models

type DynamoDBTicket struct {
	Pk        string `dynamodbav:"user_id"`
	Sk        string `dynamodbav:"sk"`
	CreatedAt string `dynamodbav:"created_at"`
	UpdatedAt string `dynamodbav:"updated_at"`
	Comments  string `dynamodbav:"comments"`
	Summary   string `dynamodbav:"summary"`
}
