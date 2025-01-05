package minio

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
)

const (
	DefaultBucketName = "gymsbro"
	DefaultExpiry     = 24 * time.Hour // 24 hour expiry for presigned URLs
)

type MinioService struct {
	client *minio.Client
}

func NewMinioService(client *minio.Client) MinioService {
	return MinioService{client: client}
}

// EnsureBucket creates the default bucket if it doesn't exist
func (s *MinioService) EnsureBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, DefaultBucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = s.client.MakeBucket(ctx, DefaultBucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return nil
}

// EnsureBucketWithPolicy creates the default bucket if it doesn't exist and sets a public read policy
func (s *MinioService) EnsureBucketWithPolicy(ctx context.Context, bucketName string) error {
	fullBucketName := fmt.Sprintf("%s-%s", DefaultBucketName, bucketName)
	exists, err := s.client.BucketExists(ctx, fullBucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = s.client.MakeBucket(ctx, fullBucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	// Set bucket policy to allow public read access
	policy := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["%s/*"]
			}
		]
	}`
	policy = fmt.Sprintf(policy, fullBucketName)

	err = s.client.SetBucketPolicy(ctx, fullBucketName, policy)
	if err != nil {
		return fmt.Errorf("failed to set bucket policy: %w", err)
	}

	return nil
}

// UploadFile uploads a file to MinIO and returns the object name
func (s *MinioService) UploadFile(ctx context.Context, reader io.Reader, bucketName string, objectName string, contentType string) error {
	fullBucketName := fmt.Sprintf("%s-%s", DefaultBucketName, bucketName)
	_, err := s.client.PutObject(ctx, fullBucketName, objectName, reader, -1, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

// GetFileURL generates a presigned URL for file download
func (s *MinioService) GetFileURL(ctx context.Context, bucketName string, objectName string) (string, error) {
	// Get presigned URL for object download
	fullBucketName := fmt.Sprintf("%s-%s", DefaultBucketName, bucketName)
	url, err := s.client.PresignedGetObject(ctx, fullBucketName, objectName, DefaultExpiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url.String(), nil
}

// DeleteFile removes a file from MinIO
func (s *MinioService) DeleteFile(ctx context.Context, bucketName string, objectName string) error {
	fullBucketName := fmt.Sprintf("%s-%s", DefaultBucketName, bucketName)
	err := s.client.RemoveObject(ctx, fullBucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}
