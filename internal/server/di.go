package server

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/joho/godotenv"
	"github.com/umekikazuya/logleaf/internal/infrastructure/dynamo"
	"github.com/umekikazuya/logleaf/internal/interface/handler"
	"github.com/umekikazuya/logleaf/internal/usecase"
)

// アプリケーションの依存関係を初期化
func InitializeDependencies() (*handler.LeafHandler, string) {
	// Env設定
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	// DynamoDB設定
	endpoint := os.Getenv("DYNAMO_ENDPOINT")
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "ap-northeast-1"
	}
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		panic("Failed to load AWS config: " + err.Error())
	}
	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		if endpoint != "" {
			o.BaseEndpoint = &endpoint
		}
	})
	tableName := os.Getenv("DYNAMO_TABLE")
	if tableName == "" {
		panic("DYNAMO_TABLE environment variable is required")
	}

	leafRepo := dynamo.NewLeafDynamoRepository(client, tableName)
	leafUsecase := usecase.NewLeafUsecase(leafRepo)
	leafHandler := handler.NewLeafHandler(leafUsecase)

	// Portを環境変数から取得（デフォルト8080）
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	return leafHandler, port
}
