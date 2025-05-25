package dynamo

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/umekikazuya/logleaf/internal/domain"
	"github.com/umekikazuya/logleaf/internal/interface/repository"
)

type LeafDynamoRepository struct {
	Client    *dynamodb.Client
	TableName string
}

func NewLeafDynamoRepository(client *dynamodb.Client, tableName string) *LeafDynamoRepository {
	return &LeafDynamoRepository{
		Client:    client,
		TableName: tableName,
	}
}

func (r *LeafDynamoRepository) Get(ctx context.Context, id string) (*domain.Leaf, error) {
	out, err := r.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &r.TableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "USER#me"},
			"sk": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, nil
	}
	var leaf domain.Leaf
	if err := attributevalue.UnmarshalMap(out.Item, &leaf); err != nil {
		return nil, err
	}
	return &leaf, nil
}

func (r *LeafDynamoRepository) List(ctx context.Context, opts repository.ListOptions) ([]domain.Leaf, error) {
	// DynamoDBにそもそも接続できてるか確認
	if r.Client == nil {
		return nil, fmt.Errorf("dynamodb client is not initialized")
	}

	listTablesOut, err1 := r.Client.ListTables(ctx, &dynamodb.ListTablesInput{})
	if err1 != nil {
		fmt.Println("接続エラー:", err1)
	} else {
		fmt.Println("接続成功。テーブル一覧:", listTablesOut.TableNames)
		found := false
		for _, t := range listTablesOut.TableNames {
			if t == r.TableName {
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("テーブル '%s' が見つかりません\n", r.TableName)
		}
	}

	queryOut, err := r.Client.Query(ctx, &dynamodb.QueryInput{
		TableName:              &r.TableName,
		KeyConditionExpression: awsString("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "USER#me"},
		},
	})
	if err != nil {
		return nil, err
	}
	var leaves []domain.Leaf
	if err := attributevalue.UnmarshalListOfMaps(queryOut.Items, &leaves); err != nil {
		return nil, err
	}
	return leaves, nil
}

func (r *LeafDynamoRepository) Put(ctx context.Context, leaf *domain.Leaf) error {
	if r.Client == nil {
		return fmt.Errorf("dynamodb client is not initialized")
	}
	if leaf.ID == "" {
		return fmt.Errorf("leaf.ID cannot be empty")
	}

	item, err := attributevalue.MarshalMap(map[string]any{
		"pk":        "USER#me",
		"sk":        leaf.ID,
		"id":        leaf.ID,
		"title":     leaf.Title,
		"url":       leaf.URL,
		"platform":  leaf.Platform,
		"tags":      leaf.Tags,
		"read":      leaf.Read,
		"synced_at": time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}
	_, err = r.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.TableName,
		Item:      item,
	})
	return err
}

func (r *LeafDynamoRepository) Update(ctx context.Context, id string, update *domain.Leaf) error {
	updateExpr := "SET #t = :title, #u = :url, #p = :platform, #tg = :tags, #r = :read"
	_, err := r.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &r.TableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "USER#me"},
			"sk": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression: awsString(updateExpr),
		ExpressionAttributeNames: map[string]string{
			"#t":  "title",
			"#u":  "url",
			"#p":  "platform",
			"#tg": "tags",
			"#r":  "read",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":title":    &types.AttributeValueMemberS{Value: update.Title},
			":url":      &types.AttributeValueMemberS{Value: update.URL},
			":platform": &types.AttributeValueMemberS{Value: update.Platform},
			":tags":     &types.AttributeValueMemberSS{Value: update.Tags},
			":read":     &types.AttributeValueMemberBOOL{Value: update.Read},
		},
	})
	return err
}

func (r *LeafDynamoRepository) Delete(ctx context.Context, id string) error {
	_, err := r.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &r.TableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "USER#me"},
			"sk": &types.AttributeValueMemberS{Value: id},
		},
	})
	return err
}

// awsStringはstringのポインタを返すヘルパー
func awsString(s string) *string {
	return &s
}
