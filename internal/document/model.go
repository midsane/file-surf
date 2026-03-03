package document

import "time"

type Document struct {
	ID        string    `dynamodbav:"id"`
	TenantID  string    `dynamodbav:"tenant_id"`
	FileName  string    `dynamodbav:"file_name"`
	S3Key     string    `dynamodbav:"s3_key"`
	Size      int64     `dynamodbav:"size"`
	CreatedAt time.Time `dynamodbav:"created_at"`
}