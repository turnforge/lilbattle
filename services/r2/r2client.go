package r2

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2Client struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucket        string
	publicURL     string // Optional: for public bucket access (e.g., "https://pub-xxx.r2.dev")
}

type R2Config struct {
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	PublicURL       string // Optional: only needed if bucket is public
}

func NewR2Client(cfg R2Config) (*R2Client, error) {
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID))
	})

	return &R2Client{
		client:        client,
		presignClient: s3.NewPresignClient(client),
		bucket:        cfg.Bucket,
		publicURL:     cfg.PublicURL,
	}, nil
}

// Upload uploads data and returns the key (not URL - use GetPublicURL or GetPresignedURL for access)
func (r *R2Client) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to R2: %w", err)
	}
	return key, nil
}

// GetPublicURL returns the public URL for a key (only works if bucket has public access enabled)
func (r *R2Client) GetPublicURL(key string) string {
	if r.publicURL == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s", r.publicURL, key)
}

// GetPresignedURL generates a time-limited signed URL for private bucket access
func (r *R2Client) GetPresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	presigned, err := r.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return presigned.URL, nil
}

// GetPresignedURLDefault generates a presigned URL with 1 hour expiry
func (r *R2Client) GetPresignedURLDefault(ctx context.Context, key string) (string, error) {
	return r.GetPresignedURL(ctx, key, time.Hour)
}

// Delete removes an object from R2
func (r *R2Client) Delete(ctx context.Context, key string) error {
	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete from R2: %w", err)
	}
	return nil
}

// HeadObject returns object metadata (size, content-type, last modified, etc.)
func (r *R2Client) HeadObject(ctx context.Context, key string) (*s3.HeadObjectOutput, error) {
	output, err := r.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return output, nil
}

// Exists checks if an object exists in R2
func (r *R2Client) Exists(ctx context.Context, key string) (bool, error) {
	_, err := r.HeadObject(ctx, key)
	if err != nil {
		return false, nil
	}
	return true, nil
}

// ObjectInfo contains metadata about an S3/R2 object
type ObjectInfo struct {
	Key          string
	Size         int64
	LastModified *time.Time
}

// List lists objects with a given prefix (returns only keys)
func (r *R2Client) List(ctx context.Context, prefix string, maxKeys int32) ([]string, error) {
	objects, err := r.ListObjects(ctx, prefix, maxKeys)
	if err != nil {
		return nil, err
	}

	keys := make([]string, len(objects))
	for i, obj := range objects {
		keys[i] = obj.Key
	}
	return keys, nil
}

// ListObjects lists objects with a given prefix and returns full metadata
func (r *R2Client) ListObjects(ctx context.Context, prefix string, maxKeys int32) ([]ObjectInfo, error) {
	output, err := r.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:  aws.String(r.bucket),
		Prefix:  aws.String(prefix),
		MaxKeys: aws.Int32(maxKeys),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	objects := make([]ObjectInfo, 0, len(output.Contents))
	for _, obj := range output.Contents {
		info := ObjectInfo{
			Key: *obj.Key,
		}
		if obj.Size != nil {
			info.Size = *obj.Size
		}
		if obj.LastModified != nil {
			info.LastModified = obj.LastModified
		}
		objects = append(objects, info)
	}
	return objects, nil
}
