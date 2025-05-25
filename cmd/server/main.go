package main

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/umekikazuya/logleaf/internal/infrastructure/dynamo"
	"github.com/umekikazuya/logleaf/internal/interface/handler"
	"github.com/umekikazuya/logleaf/internal/usecase"
)

func main() {
	r := gin.Default()

	// .envファイルの読み込み
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	// DynamoDBの設定
	endpoint := os.Getenv("DYNAMO_ENDPOINT")
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "ap-northeast-1"
	}
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		if endpoint != "" {
			o.BaseEndpoint = &endpoint
		}
	})

	// DynamoDBのテーブル名を環境変数から取得
	tableName := os.Getenv("DYNAMO_TABLE")
	if tableName == "" {
		panic("DYNAMO_TABLE environment variable is required")
	}

	leafRepo := dynamo.NewLeafDynamoRepository(client, tableName)
	leafUsecase := usecase.NewLeafUsecase(leafRepo)
	leafHandler := handler.NewLeafHandler(leafUsecase)
	api := r.Group("/api")
	{
		api.GET("/leaves", leafHandler.ListLeaves)
		api.POST("/leaves", leafHandler.AddLeaf)
		api.PATCH("/leaves/:id", leafHandler.UpdateLeaf)
		api.DELETE("/leaves/:id", leafHandler.DeleteLeaf)
	}

	// Portを環境変数から取得
	// デフォルトは8080
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		panic("failed to start server: " + err.Error())
	}
}
