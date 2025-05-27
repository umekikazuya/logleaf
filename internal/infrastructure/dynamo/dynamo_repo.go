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
		return nil, errors.New("leaf not found")
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
	var record []LeafRecord
	if err := attributevalue.UnmarshalListOfMaps(queryOut.Items, &record); err != nil {
		return nil, err
	}
	var leaves []domain.Leaf
	for _, r := range record {
		leaf, err := RecordToLeaf(&r)
		if err != nil {
			return nil, err
		}
		leaves = append(leaves, *leaf)
	}
	return leaves, nil
}

func (r *LeafDynamoRepository) Put(ctx context.Context, leaf *domain.Leaf) (*domain.Leaf, error) {
	item, err := attributevalue.MarshalMap(LeafToRecord(leaf))
	if err != nil {
		return nil, err
	}
	putItem, err := r.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.TableName,
		Item:      item,
	})
	if err != nil {
		return nil, err
	}
	if putItem == nil {
		return nil, errors.New("DBに保存できませんでした")
	}
	var record LeafRecord
	if err := attributevalue.UnmarshalMap(putItem.Item, &record); err != nil {
		return nil, err
	}
	new, err := RecordToLeaf(&record)
	if err != nil {
		return nil, err
	}
	return new, nil
}

func (r *LeafDynamoRepository) Update(ctx context.Context, update *domain.Leaf) error {
	// UpdateItemで既存レコードのみ更新する
	input := &dynamodb.UpdateItemInput{
		TableName: &r.TableName,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "USER#me"},
			"sk": &types.AttributeValueMemberS{Value: update.ID().String()},
		},
		UpdateExpression: aws.String("SET note = :note, url = :url, platform = :platform, tags = :tags, read = :read, synced_at = :synced_at"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":note":     &types.AttributeValueMemberS{Value: update.Note()},
			":url":      &types.AttributeValueMemberS{Value: update.URL().String()},
			":platform": &types.AttributeValueMemberS{Value: update.Platform()},
			":tags": &types.AttributeValueMemberSS{Value: func() []string {
				tags := update.Tags()
				s := make([]string, len(tags))
				for i, t := range tags {
					s[i] = t.String()
				}
				return s
			}()},
			":read":      &types.AttributeValueMemberBOOL{Value: update.Read()},
			":synced_at": &types.AttributeValueMemberS{Value: update.SyncedAt().Format(time.RFC3339)},
		},
		ConditionExpression: aws.String("attribute_exists(pk) AND attribute_exists(sk)"),
	}
	_, err := r.Client.UpdateItem(ctx, input)
	if err != nil {
		var cce *types.ConditionalCheckFailedException
		if errors.As(err, &cce) {
			return errors.New("leaf not found")
		}
		return err
	}
	return nil
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

// DynamoDB永続化用レコード

type LeafRecord struct {
	PK       string   `dynamodbav:"pk"`
	SK       string   `dynamodbav:"sk"`
	ID       string   `dynamodbav:"id"`
	Note     string   `dynamodbav:"note"`
	URL      string   `dynamodbav:"url"`
	Platform string   `dynamodbav:"platform"`
	Tags     []string `dynamodbav:"tags"`
	Read     bool     `dynamodbav:"read"`
	SyncedAt string   `dynamodbav:"synced_at"`
}

// EntityをRecordに変換
func LeafToRecord(l *domain.Leaf) *LeafRecord {
	tags := make([]string, len(l.Tags()))
	for i, t := range l.Tags() {
		tags[i] = t.String()
	}
	return &LeafRecord{
		PK:       "USER#me",
		SK:       l.ID().String(),
		ID:       l.ID().String(),
		Note:     l.Note(),
		URL:      l.URL().String(),
		Platform: l.Platform(),
		Tags:     tags,
		Read:     l.Read(),
		SyncedAt: l.SyncedAt().Format(time.RFC3339),
	}
}

// RecordをEntityに変換
func RecordToLeaf(r *LeafRecord) (*domain.Leaf, error) {
	// id, url, syncedAtはNewLeafで使わないため削除
	tags := make([]string, 0, len(r.Tags))
	tagSet := make(map[string]struct{})
	for _, v := range r.Tags {
		if _, exists := tagSet[v]; exists {
			return nil, errors.New("タグが重複しています")
		}
		tagSet[v] = struct{}{}
		tags = append(tags, v)
	}
	if len(tags) > 10 {
		return nil, domain.ErrTagLimitExceeded
	}
	leaf, err := domain.NewLeaf(r.Note, r.URL, r.Platform, tags)
	if err != nil {
		return nil, err
	}
	if r.Read {
		_ = leaf.MarkAsRead()
	}
	// syncedAtはprivateなので反映できない。必要ならdomain.LeafにSetSyncedAtを追加すること。
	return leaf, nil
}
