package clients

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"os"
)

type DynamoDBClient struct {
	// DynamoDB client configuration
	// For example, AWS region, credentials, etc.
	// This is a placeholder and should be replaced with actual configuration
	client *dynamodb.Client
}

func NewDynamoDBClient() (*DynamoDBClient, error) {
	// Initialize the DynamoDB client
	// This is a placeholder and should be replaced with actual initialization code

	isLocal := os.Getenv("IS_LOCAL")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"))

	if err != nil {
		return nil, fmt.Errorf("Unable to load AWS SDK config, "+"please check your AWS credentials and region: %v", err)

	}

	if isLocal == "true" {

		client := dynamodb.NewFromConfig(cfg,
			func(options *dynamodb.Options) {
				options.BaseEndpoint = aws.String("http://localhost:4566")
			})

		return &DynamoDBClient{
			client: client,
		}, nil
	}

	return &DynamoDBClient{
		client: dynamodb.NewFromConfig(cfg),
	}, nil

}
