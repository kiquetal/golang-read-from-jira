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

// SayoriClient is a client for the Sayori API
type SayoriClient struct {
	baseURL    string
	httpClient *http.Client
	logger     *log.Logger
}

// NewSayoriClient creates a new Sayori API client
func NewSayoriClient(baseURL string, logger *log.Logger) *SayoriClient {
	return &SayoriClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

// GetTickets retrieves all tickets from the Sayori API
func (c *SayoriClient) GetTickets() (models.SayoriResponses, error) {
	url := fmt.Sprintf("%s/api/v1/tickets/APIGEE", c.baseURL)

	c.logger.Printf("Fetching tickets from Sayori API: %s", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var tickets models.SayoriWrapperResponse
	if err := json.NewDecoder(resp.Body).Decode(&tickets); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var arrayOfTickets models.SayoriResponses = tickets.Data
	return arrayOfTickets, nil
}

// FindTicket searches for a specific ticket by ID in the list of tickets
func (c *SayoriClient) FindTicket(tickets models.SayoriResponses, ticketID string) (*models.SayoriResponse, error) {
	c.logger.Printf("Searching for ticket with ID: %s", ticketID)

	for _, ticket := range tickets {
		if ticket.TicketID == ticketID {
			return &ticket, nil
		}
	}

	return nil, fmt.Errorf("ticket with ID %s not found", ticketID)
}
