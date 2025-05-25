package dynamo

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/umekikazuya/logleaf/internal/domain"
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
	output, err := r.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &r.TableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "USER#me"},
			"sk": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}
	if output.Item == nil {
		return nil, nil
	}
	var leaf domain.Leaf
	if err := attributevalue.UnmarshalMap(output.Item, &leaf); err != nil {
		return nil, err
	}
	return &leaf, nil
}

func (r *LeafDynamoRepository) List(ctx context.Context, opts domain.ListOptions) ([]domain.Leaf, error) {
	// QueryInputの作成
	queryInput := &dynamodb.QueryInput{
		TableName:              &r.TableName,
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "USER#me"},
		},
	}
	// Limitの適用
	if opts.Limit > 0 {
		queryInput.Limit = aws.Int32(int32(opts.Limit))
	}
	// フィルタリングの適用
	queryOut, err := r.Client.Query(ctx, queryInput)
	if err != nil {
		return nil, err
	}
	var leaves []domain.Leaf
	if err := attributevalue.UnmarshalListOfMaps(queryOut.Items, &leaves); err != nil {
		return nil, err
	}
	return leaves, nil
}

func (r *LeafDynamoRepository) Put(ctx context.Context, leaf *domain.Leaf) (*domain.Leaf, error) {
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
		return nil, err
	}
	_, err = r.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.TableName,
		Item:      item,
	})
	if err != nil {
		return nil, err
	}
	return leaf, nil
}

func (r *LeafDynamoRepository) Update(ctx context.Context, id string, update *domain.Leaf) error {
	updateExpr := "SET #t = :title, #u = :url, #p = :platform, #tg = :tags, #r = :read"
	_, err := r.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &r.TableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "USER#me"},
			"sk": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression: aws.String(updateExpr),
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
	old, err := r.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &r.TableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "USER#me"},
			"sk": &types.AttributeValueMemberS{Value: id},
		},
		ReturnValues: types.ReturnValueAllOld,
	})
	if err != nil {
		return err
	}
	if old.Attributes == nil {
		return errors.New("leaf not found")
	}
	return nil
}
