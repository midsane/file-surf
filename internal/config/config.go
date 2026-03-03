package config
/*
read environment variables, return config struct
*/
import (
	"log"
	"os"
)

type Config struct {
	Port        string
	AWSRegion   string
	S3Bucket    string
	TenantTable string
	UserTable   string
	DocTable    string
}

func Load() *Config {
	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		AWSRegion:   getEnv("AWS_REGION", "ap-south-1"),
		S3Bucket:    getEnv("S3_BUCKET", ""),
		TenantTable: getEnv("TENANT_TABLE", "tenants"),
		UserTable:   getEnv("USER_TABLE", "users"),
		DocTable:    getEnv("DOCUMENT_TABLE", "documents"),
	}

	if cfg.S3Bucket == "" {
		log.Fatal("S3_BUCKET must be set")
	}

	return cfg
}

func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

