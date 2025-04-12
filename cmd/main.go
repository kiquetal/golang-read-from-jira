package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kiquetal/golang-read-from-jira/internal/service"
)

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "[JIRA-READER] ", log.LstdFlags|log.Lshortfile)

	// Create service
	svc, err := service.NewTicketService(logger)
	if err != nil {
		logger.Fatalf("Failed to create ticket service: %v", err)
	}

	// Get comments by user
	userComments, err := svc.GetCommentsByUser()
	if err != nil {
		logger.Fatalf("Failed to get comments: %v", err)
	}

	// Print last comment for each user and ticket
	if len(userComments) == 0 {
		fmt.Println("No comments found for any user")
		return
	}

	for user, ticketComments := range userComments {
		fmt.Printf("\nLast comments by %s:\n", user)
		fmt.Println(strings.Repeat("-", 40))

		for ticketID, lastComment := range ticketComments {
			fmt.Printf("\nTicket: %s\n", ticketID)
			fmt.Printf("Last comment: %s\n", lastComment)
		}
		fmt.Println(strings.Repeat("=", 60))
	}
}
