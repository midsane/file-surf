package tenant

import (
	"context"
	"fmt"
	"time"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
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

func (r *Repository) Create(ctx context.Context, t *Tenant) error {
	t.CreatedAt = time.Now()

	item := map[string]interface{}{
		"PK":        fmt.Sprintf("TENANT#%s", t.ID),
		"SK":        "METADATA",
		"Type":      "TENANT",
		"id":        t.ID,
		"name":      t.Name,
		"created_at": t.CreatedAt,
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("marshal tenant: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.tableName,
		Item:      av,
		ConditionExpression: awsString("attribute_not_exists(PK)"),
	})

	if err != nil {
		return fmt.Errorf("put tenant: %w", err)
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, tenantID string) (*Tenant, error) {
	key := map[string]interface{}{
		"PK": fmt.Sprintf("TENANT#%s", tenantID),
		"SK": "METADATA",
	}

	av, err := attributevalue.MarshalMap(key)
	if err != nil {
		return nil, fmt.Errorf("marshal key: %w", err)
	}

	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &r.tableName,
		Key:       av,
	})

	if err != nil {
		return nil, fmt.Errorf("get tenant: %w", err)
	}

	if out.Item == nil {
		return nil, nil
	}

	var tenant Tenant
	if err := attributevalue.UnmarshalMap(out.Item, &tenant); err != nil {
		return nil, fmt.Errorf("unmarshal tenant: %w", err)
	}

	return &tenant, nil
}

func awsString(s string) *string {
	return &s
}