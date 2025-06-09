package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/joho/godotenv"
	"github.com/umekikazuya/logleaf/internal/domain"
	"github.com/umekikazuya/logleaf/internal/infrastructure/dynamo"
	"github.com/umekikazuya/logleaf/internal/infrastructure/qiita"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	token := os.Getenv("QIITA_TOKEN")
	user := os.Getenv("QIITA_USER")

	if token == "" || user == "" {
		fmt.Println("QIITA_TOKENとQIITA_USERを環境変数で指定してください")
		os.Exit(1)
	}
	client := qiita.NewQiitaClient(token, user)
	items, err := client.FetchStocksAll(
		context.Background(),
	)
	if err != nil {
		fmt.Println("Qiita API取得エラー:", err)
		os.Exit(1)
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
	dynamoClient := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		if endpoint != "" {
			o.BaseEndpoint = &endpoint
		}
	})
	tableName := os.Getenv("DYNAMO_TABLE")
	if tableName == "" {
		panic("DYNAMO_TABLE environment variable is required")
	}
	repo := dynamo.NewLeafDynamoRepository(dynamoClient, tableName)

	// 既存LeafのURL一覧を取得して差分同期
	leaves, err := repo.List(ctx, domain.ListOptions{Limit: 1000})
	if err != nil {
		fmt.Println("DynamoDB取得エラー:", err)
		os.Exit(1)
	}
	existingURLs := make(map[string]struct{})
	for _, leaf := range leaves {
		existingURLs[leaf.URL().String()] = struct{}{}
	}
	fmt.Println(existingURLs)

	countNew := 0
	for _, item := range items {
		if _, exists := existingURLs[item.URL]; exists {
			continue
		}
		tags := make([]string, len(item.Tags))
		for i, t := range item.Tags {
			tags[i] = t.Name
		}
		leaf, err := domain.NewLeaf(
			item.Title, item.URL, "qiita", tags, false,
		)
		if err != nil {
			fmt.Println("Leaf生成エラー:", err)
			continue
		}
		_, err = repo.Put(ctx, leaf)
		if err != nil {
			fmt.Println("DynamoDB保存エラー:", err)
		}
		countNew++
		time.Sleep(1 * time.Second) // API制限対策
	}
	fmt.Printf("Qiitaストック記事の同期が完了しました（新規追加: %d件）\n", countNew)
}
