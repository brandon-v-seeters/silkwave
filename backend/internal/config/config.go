package config

import (
	"os"
)

type Config struct {
	Environment    string
	ServerPort     string
	ArangoEndpoint string
	ArangoDatabase string
	ArangoUsername string
	ArangoPassword string
	JWTSecret      string
	PasswordSecret string

	// Cloudflare R2 Storage
	R2AccountID       string
	R2AccessKeyID     string
	R2SecretAccessKey string
	R2BucketName      string
}

func (c *Config) IsDev() bool {
	return c.Environment == "development" || c.Environment == "dev"
}

func Load() *Config {
	return &Config{
		Environment:    getEnv("GO_ENV", "development"),
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		ArangoEndpoint: getEnv("ARANGO_ENDPOINT", "http://localhost:8529"),
		ArangoDatabase: getEnv("ARANGO_DATABASE", "silk_wave"),
		ArangoUsername: getEnv("ARANGO_USERNAME", "root"),
		ArangoPassword: getEnv("ARANGO_PASSWORD", ""),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		PasswordSecret: getEnv("PASSWORD_SECRET", "your-password-secret-change-in-production"),

		// Cloudflare R2 Storage
		R2AccountID:       getEnv("R2_ACCOUNT_ID", ""),
		R2AccessKeyID:     getEnv("R2_ACCESS_KEY_ID", ""),
		R2SecretAccessKey: getEnv("R2_SECRET_ACCESS_KEY", ""),
		R2BucketName:      getEnv("R2_BUCKET_NAME", "silkwave"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
