package clients

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/kiquetal/golang-read-from-jira/internal/models"
	"os"
	"time"
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

	if isLocal == "True" {

		fmt.Printf("Starting DynamoDB in local mode\n")
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

func (d *DynamoDBClient) CreateTableLocal(tableName, pk, sk string) error {
	// Create a table in DynamoDB
	// This is a placeholder and should be replaced with actual table creation code

	fmt.Printf("IS_LOCAL: %s\n", os.Getenv("IS_LOCAL"))
	if os.Getenv("IS_LOCAL") != "True" {

		return fmt.Errorf("not in local mode")
	}

	// Check if the table already exists

	describeTableInput := &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}
	_, err := d.client.DescribeTable(context.TODO(), describeTableInput)
	if err == nil {
		fmt.Printf("Table %s already exists\n", tableName)
		return nil
	}

	_, err = d.client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String(pk),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String(sk),
				KeyType:       types.KeyTypeRange,
			},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String(pk),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String(sk),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	// Return the table name and primary key
	return nil
}

func (d *DynamoDBClient) PutTicketCommentsInDynammo(userBot, ticket, description string) error {

	var dynamoDbTicket models.DynamoDBTicket = models.DynamoDBTicket{
		Pk:        userBot,
		Sk:        ticket,
		Comments:  description,
		UpdatedAt: time.Now().Format(time.RFC3339),
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	// Convert the ticket to a DynamoDB item
	item, err := attributevalue.MarshalMap(dynamoDbTicket)
	if err != nil {
		return fmt.Errorf("failed to marshal ticket: %w", err)
	}
	// Put the item in the DynamoDB table
	_, err = d.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("roaster-for-slack-test-users"),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to put item in DynamoDB: %w", err)
	}
	return nil

}
