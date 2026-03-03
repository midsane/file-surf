package user

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

func (r *Repository) Create(ctx context.Context, u *User) error {
	u.CreatedAt = time.Now()

	item := map[string]interface{}{
		"PK":         fmt.Sprintf("TENANT#%s", u.TenantID),
		"SK":         fmt.Sprintf("USER#%s", u.ID),
		"Type":       "USER",
		"id":         u.ID,
		"tenant_id":  u.TenantID,
		"email":      u.Email,
		"created_at": u.CreatedAt,
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("marshal user: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.tableName,
		Item:      av,
	})

	if err != nil {
		return fmt.Errorf("put user: %w", err)
	}

	return nil
}

func (r *Repository) GetByTenant(ctx context.Context, tenantID string) ([]User, error) {
	pk := fmt.Sprintf("TENANT#%s", tenantID)

	out, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              &r.tableName,
		KeyConditionExpression: awsString("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: pk},
			":sk": &types.AttributeValueMemberS{Value: "USER#"},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}

	var users []User
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &users); err != nil {
		return nil, fmt.Errorf("unmarshal users: %w", err)
	}

	return users, nil
}

func awsString(s string) *string {
	return &s
}
