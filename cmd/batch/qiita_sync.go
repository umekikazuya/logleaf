package main

import (
	"context"
	"fmt"
	"os"
	"time"

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
	qiitaClient := qiita.NewQiitaClient(token, user)
	items, err := qiitaClient.FetchStocksAll(
		context.Background(),
	)
	if err != nil {
		fmt.Println("Qiita API取得エラー:", err)
		os.Exit(1)
	}

	_ = godotenv.Load() // 本番は.env不要なのでエラー無視

	dynamoClient, tableName, err := dynamo.NewDynamoClientAndTable(ctx)
	if err != nil {
		panic(err)
	}
	repo := dynamo.NewLeafDynamoRepository(dynamoClient, tableName)
	ctx = context.Background()

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

	countNew := 0
	for _, item := range items {
		if _, exists := existingURLs[item.URL]; exists {
			continue // 既存記事はスキップ
		}
		tags := make([]string, len(item.Tags))
		for i, t := range item.Tags {
			tags[i] = t.Name
		}
		leaf, err := domain.NewLeaf(item.Title, item.URL, "qiita", tags, false)
		if err != nil {
			fmt.Println("Leaf生成エラー:", err)
			continue
		}
		_, err = repo.Put(ctx, leaf)
		if err != nil {
			fmt.Println("DynamoDB保存エラー:", err)
		}
		countNew++
		time.Sleep(200 * time.Millisecond) // API制限対策
	}
	fmt.Printf("Qiitaストック記事の同期が完了しました（新規追加: %d件）\n", countNew)
}
