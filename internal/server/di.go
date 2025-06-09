package server

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/umekikazuya/logleaf/internal/application"
	"github.com/umekikazuya/logleaf/internal/infrastructure/dynamo"
	"github.com/umekikazuya/logleaf/internal/interface/handler"
)

// アプリケーションの依存関係を初期化
func InitializeDependencies() (*handler.LeafHandler, string) {
	_ = godotenv.Load()

	client, tableName, err := dynamo.NewDynamoClientAndTable(context.Background())
	if err != nil {
		panic(err)
	}

	leafRepo := dynamo.NewLeafDynamoRepository(client, tableName)
	leafUsecase := application.NewLeafUsecase(leafRepo)
	leafHandler := handler.NewLeafHandler(leafUsecase)

	// Portを環境変数から取得（デフォルト8080）
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	return leafHandler, port
}
