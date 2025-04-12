package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/kiquetal/golang-read-from-jira/internal/models"
)

// JiraClient is a client for the Jira API
type JiraClient struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
	logger     *log.Logger
}

// NewJiraClient creates a new Jira API client
func NewJiraClient(baseURL, apiToken string, logger *log.Logger) *JiraClient {
	return &JiraClient{
		baseURL:  baseURL,
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

// GetTicket retrieves a ticket from the Jira API
func (c *JiraClient) GetTicket(ticketID string) (*models.JiraTicket, error) {
	url := fmt.Sprintf("%s/rest/api/2/issue/%s?expand=comments", c.baseURL, ticketID)

	c.logger.Printf("Fetching ticket from Jira API: %s", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Basic "+c.apiToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var ticket models.JiraTicket
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Println("Body response: ", string(body))

	if err := json.Unmarshal(body, &ticket); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &ticket, nil
}

// GetCommentsByUser retrieves comments made by a specific user on a ticket
func (c *JiraClient) GetCommentsByUser(ticket *models.JiraTicket, displayName string) ([]string, error) {
	c.logger.Printf("Filtering comments by display name: %s", displayName)

	var userComments []string

	for _, comment := range ticket.Fields.Comments {
		if comment.Author.DisplayName == displayName {
			userComments = append(userComments, comment.Body)
		}
	}

	c.logger.Printf("Found %d comments by user %s", len(userComments), displayName)

	return userComments, nil
}

// GetLastCommentByUser retrieves the last comment made by a specific user on a ticket
func (c *JiraClient) GetLastCommentByUser(ticket *models.JiraTicket, displayName string) (string, error) {
	c.logger.Printf("Finding last comment by display name: %s", displayName)

	var lastComment *models.JiraComment

	for i := range ticket.Fields.Comments {
		comment := &ticket.Fields.Comments[i]
		if comment.Author.Name == displayName {
			if lastComment == nil || comment.Created.After(lastComment.Created) {
				lastComment = comment
			}
		}
	}

	if lastComment == nil {
		c.logger.Printf("No comments found by user %s", displayName)
		return "", nil
	}

	c.logger.Printf("Found last comment by user %s created at %s", displayName, lastComment.Created)
	return lastComment.Body, nil
}
