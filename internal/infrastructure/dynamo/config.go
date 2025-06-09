package dynamo

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// NewDynamoClientAndTable returns a DynamoDB client and table name using env vars (endpoint, region, table)
func NewDynamoClientAndTable(ctx context.Context) (*dynamodb.Client, string, error) {
	endpoint := os.Getenv("DYNAMO_ENDPOINT")
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "ap-northeast-1"
	}
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, "", err
	}
	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		if endpoint != "" {
			o.BaseEndpoint = &endpoint
		}
	})
	tableName := os.Getenv("DYNAMO_TABLE")
	if tableName == "" {
		return nil, "", errors.New("DYNAMO_TABLE environment variable is required")
	}
	return client, tableName, nil
}
