package storage

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// R2Config holds Cloudflare R2 configuration.
type R2Config struct {
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	CustomDomain    string // Optional: e.g., "cdn.silkwave.io" - if set, used instead of R2 API endpoint
}

// Client wraps the S3 client for Cloudflare R2 operations.
type Client struct {
	s3Client   *s3.Client
	bucketName string
	resolver   *KeyResolver
}

// NewClient creates a new R2/S3 storage client.
func NewClient(ctx context.Context, cfg R2Config) (*Client, error) {
	cfg = normalizeR2Config(cfg)
	if err := validateR2Config(cfg); err != nil {
		return nil, err
	}

	// R2 API endpoint - required for presigned URLs and authenticated operations
	r2Endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID)

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               r2Endpoint,
			HostnameImmutable: true,
		}, nil
	})

	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &Client{
		s3Client:   s3Client,
		bucketName: cfg.BucketName,
		resolver:   NewKeyResolver(),
	}, nil
}

func normalizeR2Config(cfg R2Config) R2Config {
	cfg.AccountID = strings.TrimSpace(cfg.AccountID)
	cfg.AccessKeyID = strings.TrimSpace(cfg.AccessKeyID)
	cfg.SecretAccessKey = strings.TrimSpace(cfg.SecretAccessKey)
	cfg.BucketName = strings.TrimSpace(cfg.BucketName)
	cfg.CustomDomain = strings.TrimSpace(cfg.CustomDomain)
	return cfg
}

func validateR2Config(cfg R2Config) error {
	var missing []string
	if cfg.AccountID == "" {
		missing = append(missing, "R2_ACCOUNT_ID (or CF_ACCOUNT_ID)")
	}
	if cfg.AccessKeyID == "" {
		missing = append(missing, "R2_ACCESS_KEY_ID")
	}
	if cfg.SecretAccessKey == "" {
		missing = append(missing, "R2_SECRET_ACCESS_KEY")
	}
	if cfg.BucketName == "" {
		missing = append(missing, "R2_BUCKET_NAME")
	}
	if len(missing) > 0 {
		return fmt.Errorf("invalid R2 storage config: missing %s", strings.Join(missing, ", "))
	}
	return nil
}

// Resolver returns the key resolver for generating object keys.
func (c *Client) Resolver() *KeyResolver {
	return c.resolver
}

// BucketName returns the configured bucket name.
func (c *Client) BucketName() string {
	return c.bucketName
}

// --- Basic Operations ---

// Upload uploads data to the specified object key.
func (c *Client) Upload(ctx context.Context, key string, body io.Reader, contentType string) error {
	_, err := c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucketName),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("failed to upload object %s: %w", key, err)
	}
	return nil
}

// Download retrieves an object from storage.
func (c *Client) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	result, err := c.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download object %s: %w", key, err)
	}
	return result.Body, nil
}

// Delete removes an object from storage.
func (c *Client) Delete(ctx context.Context, key string) error {
	_, err := c.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object %s: %w", key, err)
	}
	return nil
}

// Exists checks if an object exists at the given key.
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	_, err := c.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		// Check if it's a "not found" error
		return false, nil
	}
	return true, nil
}

// List returns all object keys with the given prefix.
func (c *Client) List(ctx context.Context, prefix string) ([]string, error) {
	var keys []string

	paginator := s3.NewListObjectsV2Paginator(c.s3Client, &s3.ListObjectsV2Input{
		Bucket: aws.String(c.bucketName),
		Prefix: aws.String(prefix),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list objects with prefix %s: %w", prefix, err)
		}

		for _, obj := range page.Contents {
			keys = append(keys, *obj.Key)
		}
	}

	return keys, nil
}

// Copy copies an object from one key to another.
func (c *Client) Copy(ctx context.Context, sourceKey, destKey string) error {
	copySource := fmt.Sprintf("%s/%s", c.bucketName, sourceKey)

	_, err := c.s3Client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(c.bucketName),
		CopySource: aws.String(copySource),
		Key:        aws.String(destKey),
	})
	if err != nil {
		return fmt.Errorf("failed to copy object from %s to %s: %w", sourceKey, destKey, err)
	}
	return nil
}

// Move moves an object from one key to another (copy + delete).
func (c *Client) Move(ctx context.Context, sourceKey, destKey string) error {
	if err := c.Copy(ctx, sourceKey, destKey); err != nil {
		return err
	}
	return c.Delete(ctx, sourceKey)
}

// DeletePrefix deletes all objects with the given prefix.
func (c *Client) DeletePrefix(ctx context.Context, prefix string) error {
	keys, err := c.List(ctx, prefix)
	if err != nil {
		return err
	}

	for _, key := range keys {
		if err := c.Delete(ctx, key); err != nil {
			return err
		}
	}

	return nil
}

// GetPresignedURL generates a presigned URL for downloading an object.
func (c *Client) GetPresignedURL(ctx context.Context, key string, expireSeconds int64) (string, error) {
	presignClient := s3.NewPresignClient(c.s3Client)

	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(time.Duration(expireSeconds)*time.Second))

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL for %s: %w", key, err)
	}

	return request.URL, nil
}

// GetPresignedUploadURL generates a presigned URL for uploading an object.
func (c *Client) GetPresignedUploadURL(ctx context.Context, key, contentType string, expireSeconds int64) (string, error) {
	presignClient := s3.NewPresignClient(c.s3Client)

	request, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucketName),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACL("public-read"),
	}, s3.WithPresignExpires(time.Duration(expireSeconds)*time.Second))

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned upload URL for %s: %w", key, err)
	}

	return request.URL, nil
}
