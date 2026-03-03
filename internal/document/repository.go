package document

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Repository struct {
	client    *dynamodb.Client
	tableName string
}

func NewRepository(client *dynamodb.Client, tableName string) *Repository {
	return &Repository{
		client:    client,
		tableName: tableName,
	}
}

func (r *Repository) Create(ctx context.Context, d *Document) error {
	d.CreatedAt = time.Now()

	item := map[string]interface{}{
		"PK":         fmt.Sprintf("TENANT#%s", d.TenantID),
		"SK":         fmt.Sprintf("DOC#%s", d.ID),
		"Type":       "DOC",
		"id":         d.ID,
		"tenant_id":  d.TenantID,
		"file_name":  d.FileName,
		"s3_key":     d.S3Key,
		"size":       d.Size,
		"created_at": d.CreatedAt,
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("marshal doc: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.tableName,
		Item:      av,
	})

	if err != nil {
		return fmt.Errorf("put doc: %w", err)
	}

	return nil
}

func (r *Repository) GetByTenant(ctx context.Context, tenantID string) ([]Document, error) {
	pk := fmt.Sprintf("TENANT#%s", tenantID)

	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              &r.tableName,
		KeyConditionExpression: awsString("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: pk},
			":sk": &types.AttributeValueMemberS{Value: "DOC#"},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("query docs: %w", err)
	}

	var docs []Document
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &docs); err != nil {
		return nil, fmt.Errorf("unmarshal docs: %w", err)
	}

	return docs, nil
}

func awsString(s string) *string {
	return &s
}
