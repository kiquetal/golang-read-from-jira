package service

import (
	"fmt"
	"github.com/kiquetal/golang-read-from-jira/internal/clients"
	"github.com/kiquetal/golang-read-from-jira/internal/models"
	"log"
	"os"
)

// TicketService coordinates between the Sayori and Jira clients
type TicketService struct {
	sayoriClient *clients.SayoriClient
	jiraClient   *clients.JiraClient
	dynamoClient *clients.DynamoDBClient
	logger       *log.Logger
}

// NewTicketService creates a new TicketService
func NewTicketService(logger *log.Logger) (*TicketService, error) {
	// Get configuration from environment variables
	sayoriBaseURL := getEnvOrDefault("SAYORI_BASE_URL", "http://localhost:8080")
	jiraBaseURL := getEnvOrDefault("JIRA_BASE_URL", "https://jira.example.com")
	jiraToken := os.Getenv("JIRA_TOKEN")

	// Validate required environment variables
	if jiraToken == "" {
		return nil, fmt.Errorf("JIRA_TOKEN environment variable is required")
	}

	// Create clients
	sayoriClient := clients.NewSayoriClient(sayoriBaseURL, logger)
	jiraClient := clients.NewJiraClient(jiraBaseURL, jiraToken, logger)

	dynamoClient, err := clients.NewDynamoDBClient()

	if err != nil {
		return nil, fmt.Errorf("failed to create DynamoDB client: %w", err)
	}

	// Create the DynamoDB table if it doesn't exist
	err = dynamoClient.CreateTableLocal("roaster-for-slack-test-users", "pk", "sk")
	if err != nil {
		fmt.Printf("failed to create DynamoDB table locally: %w", err)
	}

	return &TicketService{
		sayoriClient: sayoriClient,
		jiraClient:   jiraClient,
		logger:       logger,
		dynamoClient: dynamoClient,
	}, nil
}

// GetCommentsByDisplayName fetches the Jira ticket directly by ID,
// then extracts comments made by the specified display name
func (s *TicketService) GetCommentsByDisplayName(ticketID, displayName string) ([]string, error) {
	// Step 1: Get the Jira ticket directly using the ticket ID
	s.logger.Printf("Fetching ticket %s from Jira", ticketID)
	jiraTicket, err := s.jiraClient.GetTicket(ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket from Jira: %w", err)
	}

	// Step 2: Get comments by the specified display name
	s.logger.Printf("Extracting comments by display name: %s", displayName)
	comments, err := s.jiraClient.GetCommentsByUser(jiraTicket, displayName)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments by display name: %w", err)
	}

	return comments, nil
}

// GetCommentsByUser groups users with their tickets from Sayori,
// then looks in Jira for the last comment issued by this user in the ticket
func (s *TicketService) GetCommentsByUser() (map[string]map[string]string, map[string]string, error) {
	// Step 1: Get all tickets from Sayori
	s.logger.Printf("Fetching all tickets from Sayori")
	sayoriTickets, err := s.sayoriClient.GetTickets()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get tickets from Sayori: %w", err)
	}

	// Step 2: Group tickets by user
	s.logger.Printf("Grouping tickets by user")
	userTickets := make(map[string][]models.SayoriResponse)
	mapBotUserDisplayName := make(map[string]string)
	for _, ticket := range sayoriTickets {
		displayName := ticket.BotUser.DisplayName
		bot_user_id := ticket.BotUserID
		if _, exists := mapBotUserDisplayName[displayName]; !exists {
			mapBotUserDisplayName[displayName] = bot_user_id
		}
		userTickets[displayName] = append(userTickets[displayName], ticket)
	}

	for displayName, tickets := range userTickets {

		s.logger.Printf("User: %s, Tickets: %d", displayName, len(tickets))
		for _, ticket := range tickets {
			s.logger.Printf("Ticket ID: %s, Link: %s", ticket.TicketID, ticket.Link)
		}
	}

	// Step 3: For each user, get the last comment from all their tickets
	s.logger.Printf("Getting last comment for each user's tickets")
	result := make(map[string]map[string]string)
	for displayName, tickets := range userTickets {
		s.logger.Printf("Processing tickets for user: %s", displayName)
		ticketLastComments := make(map[string]string)

		for _, ticket := range tickets {
			// Only process tickets where the link starts with "https://jira."
			if len(ticket.Link) < 13 || ticket.Link[:13] != "https://jira." {
				s.logger.Printf("Skipping ticket %s as link does not start with 'https://jira.'", ticket.TicketID)
				continue
			}

			jiraTicketID := ticket.TicketID
			s.logger.Printf("Fetching ticket %s from Jira for user %s", jiraTicketID, displayName)

			jiraTicket, err := s.jiraClient.GetTicket(jiraTicketID)
			if err != nil {
				s.logger.Printf("Failed to get ticket %s from Jira: %v", jiraTicketID, err)
				continue
			}

			// Get the last comment from the specific user in the Jira ticket
			lastComment, err := s.jiraClient.GetLastCommentByUser(jiraTicket, displayName)
			if err != nil {
				s.logger.Printf("Failed to get last comment by user %s for ticket %s: %v", displayName, jiraTicketID, err)
				continue
			}

			if lastComment != "" {
				ticketLastComments[jiraTicketID] = lastComment
			}
		}

		if len(ticketLastComments) > 0 {
			result[displayName] = ticketLastComments
		}
	}

	return result, mapBotUserDisplayName, nil
}

// getEnvOrDefault gets an environment variable or returns a default value if not set
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (s *TicketService) PutItemInDynamo(botUserId, ticketId, lastComment string) error {

	ticketIdd := fmt.Sprintf("#TICKET#%s", ticketId)
	err := s.dynamoClient.PutTicketCommentsInDynammo(botUserId, ticketIdd, lastComment)
	if err != nil {
		return fmt.Errorf("failed to put item in DynamoDB: %w", err)
	}
	return nil
}
