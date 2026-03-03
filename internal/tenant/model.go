package tenant

import "time"

type Tenant struct {
	ID        string    `dynamodbav:"id"`
	Name      string    `dynamodbav:"name"`
	CreatedAt time.Time `dynamodbav:"created_at"`
}