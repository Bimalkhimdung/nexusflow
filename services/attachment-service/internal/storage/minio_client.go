package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nexusflow/nexusflow/pkg/logger"
)

type MinIOClient struct {
	client     *minio.Client
	bucketName string
	log        *logger.Logger
}

func NewMinIOClient(endpoint, accessKey, secretKey, bucketName string, useSSL bool, log *logger.Logger) (*MinIOClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	mc := &MinIOClient{
		client:     client,
		bucketName: bucketName,
		log:        log,
	}

	// Ensure bucket exists
	if err := mc.ensureBucket(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ensure bucket: %w", err)
	}

	log.Sugar().Infow("Connected to MinIO", "bucket", bucketName)
	return mc, nil
}

func (m *MinIOClient) ensureBucket(ctx context.Context) error {
	exists, err := m.client.BucketExists(ctx, m.bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket: %w", err)
	}

	if !exists {
		err = m.client.MakeBucket(ctx, m.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		m.log.Sugar().Infow("Created bucket", "bucket", m.bucketName)
	}

	return nil
}

// UploadFile uploads a file to MinIO
func (m *MinIOClient) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) error {
	_, err := m.client.PutObject(ctx, m.bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	m.log.Sugar().Infow("Uploaded file", "object", objectName, "size", size)
	return nil
}

// GetPresignedURL generates a presigned URL for downloading
func (m *MinIOClient) GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := m.client.PresignedGetObject(ctx, m.bucketName, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url.String(), nil
}

// DeleteFile deletes a file from MinIO
func (m *MinIOClient) DeleteFile(ctx context.Context, objectName string) error {
	err := m.client.RemoveObject(ctx, m.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	m.log.Sugar().Infow("Deleted file", "object", objectName)
	return nil
}

// GetFileInfo gets file information
func (m *MinIOClient) GetFileInfo(ctx context.Context, objectName string) (*minio.ObjectInfo, error) {
	info, err := m.client.StatObject(ctx, m.bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return &info, nil
}
