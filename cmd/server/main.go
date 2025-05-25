package main

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
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

	godotenv.Load()
	endpoint := os.Getenv("DYNAMO_ENDPOINT")
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "ap-northeast-1"
	}
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithEndpointResolver(
			aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				if endpoint != "" && service == dynamodb.ServiceID {
					return aws.Endpoint{
						URL:           endpoint,
						SigningRegion: region,
					}, nil
				}
				return aws.Endpoint{}, &aws.EndpointNotFoundError{}
			}),
		),
	)
	if err != nil {
		panic("failed to load AWS config: " + err.Error())
	}
	client := dynamodb.NewFromConfig(cfg)

	leafRepo := dynamo.NewLeafDynamoRepository(client, os.Getenv("DYNAMO_TABLE"))
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
	r.Run(":" + port)
}
