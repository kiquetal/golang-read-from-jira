# Golang Jira Comment Reader

A Go application that fetches ticket information from Jira and extracts comments by specific users. It can either fetch all tickets from the Sayori API and group them by user, or directly fetch a specific Jira ticket by ID.

## Project Structure

```
.
├── cmd/
│   └── main.go           # Application entry point
├── internal/
│   ├── clients/
│   │   ├── jira.go       # Jira API client
│   │   └── sayori.go     # Sayori API client
│   ├── models/
│   │   ├── jira.go       # Jira data models
│   │   └── sayori.go     # Sayori data models
│   └── service/
│       └── ticket_service.go # Business logic
├── go.mod                # Go module file
└── README.md             # This file
```

## Requirements

- Go 1.22.2 or later
- Jira API credentials
- Access to the Sayori API

## Configuration

The application is configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `SAYORI_BASE_URL` | Base URL for the Sayori API | `http://localhost:8080` |
| `JIRA_BASE_URL` | Base URL for the Jira API | `https://jira.example.com` |
| `JIRA_API_TOKEN` | API token for Jira API authentication | (required) |
| `AWS_ACCESS_KEY_ID`       | Your AWS access key ID                           | (required)            |
| `AWS_SECRET_ACCESS_KEY`   | Your AWS secret access key                       | (required)            |
| `AWS_REGION`              | The AWS region where your DynamoDB table is located | `us-east-1`         |

## Usage

### Building the Application

```bash
go build -o jira-reader ./cmd/main.go
```

### Running the Application

The application can be run in two modes:

#### Mode 1: Fetch all tickets from Sayori and group by user (default)

This mode fetches all tickets from the Sayori API, groups them by user, and then gets the last comment from each user's tickets in Jira.

```bash
# Set required environment variables
export JIRA_USERNAME=your-username
export JIRA_API_TOKEN=your-api-token

# Run the application
./jira-reader
```

