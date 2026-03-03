package user

import "time"

type User struct {
	ID        string    `dynamodbav:"id"`
	TenantID  string    `dynamodbav:"tenant_id"`
	Email     string    `dynamodbav:"email"`
	CreatedAt time.Time `dynamodbav:"created_at"`
}