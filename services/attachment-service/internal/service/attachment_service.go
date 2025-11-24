package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/nexusflow/nexusflow/pkg/kafka"
	"github.com/nexusflow/nexusflow/pkg/logger"
	"github.com/nexusflow/nexusflow/services/attachment-service/internal/models"
	"github.com/nexusflow/nexusflow/services/attachment-service/internal/repository"
	"github.com/nexusflow/nexusflow/services/attachment-service/internal/storage"
)

type AttachmentService struct {
	repo     *repository.AttachmentRepository
	storage  *storage.MinIOClient
	producer *kafka.Producer
	log      *logger.Logger
}

func NewAttachmentService(repo *repository.AttachmentRepository, storage *storage.MinIOClient, producer *kafka.Producer, log *logger.Logger) *AttachmentService {
	return &AttachmentService{repo: repo, storage: storage, producer: producer, log: log}
}

// UploadAttachment handles file upload
func (s *AttachmentService) UploadAttachment(ctx context.Context, metadata *models.UploadMetadata, data []byte) (*models.Attachment, error) {
	// Generate unique filename
	ext := filepath.Ext(metadata.Filename)
	uniqueID := uuid.New().String()
	storagePath := fmt.Sprintf("%s/%s/%s%s", metadata.EntityType, metadata.EntityID, uniqueID, ext)

	// Upload to MinIO
	reader := bytes.NewReader(data)
	if err := s.storage.UploadFile(ctx, storagePath, reader, int64(len(data)), metadata.ContentType); err != nil {
		return nil, fmt.Errorf("upload to storage: %w", err)
	}

	// Create database record
	attachment := &models.Attachment{
		EntityType:       metadata.EntityType,
		EntityID:         metadata.EntityID,
		Filename:         fmt.Sprintf("%s%s", uniqueID, ext),
		OriginalFilename: metadata.Filename,
		ContentType:      metadata.ContentType,
		Size:             int64(len(data)),
		StoragePath:      storagePath,
		UploaderID:       metadata.UploaderID,
	}

	if err := s.repo.CreateAttachment(ctx, attachment); err != nil {
		// Cleanup: delete from storage if DB insert fails
		_ = s.storage.DeleteFile(ctx, storagePath)
		return nil, fmt.Errorf("create attachment record: %w", err)
	}

	s.publishEvent("attachment.uploaded", attachment.EntityID, map[string]interface{}{
		"attachment_id": attachment.ID,
		"entity_type":   attachment.EntityType,
		"filename":      attachment.OriginalFilename,
		"size":          attachment.Size,
	})

	return attachment, nil
}

// GetAttachment retrieves attachment metadata
func (s *AttachmentService) GetAttachment(ctx context.Context, id string) (*models.Attachment, error) {
	return s.repo.GetAttachment(ctx, id)
}

// GetDownloadURL generates a presigned URL for downloading
func (s *AttachmentService) GetDownloadURL(ctx context.Context, attachmentID string) (string, int64, error) {
	attachment, err := s.repo.GetAttachment(ctx, attachmentID)
	if err != nil {
		return "", 0, err
	}

	expiry := 1 * time.Hour
	url, err := s.storage.GetPresignedURL(ctx, attachment.StoragePath, expiry)
	if err != nil {
		return "", 0, fmt.Errorf("generate download URL: %w", err)
	}

	return url, int64(expiry.Seconds()), nil
}

// ListAttachments lists attachments for an entity
func (s *AttachmentService) ListAttachments(ctx context.Context, entityType, entityID string) ([]*models.Attachment, error) {
	return s.repo.ListAttachments(ctx, entityType, entityID)
}

// DeleteAttachment deletes an attachment
func (s *AttachmentService) DeleteAttachment(ctx context.Context, id string) error {
	attachment, err := s.repo.GetAttachment(ctx, id)
	if err != nil {
		return err
	}

	// Delete from storage
	if err := s.storage.DeleteFile(ctx, attachment.StoragePath); err != nil {
		s.log.Sugar().Warnw("Failed to delete from storage", "error", err, "path", attachment.StoragePath)
	}

	// Delete from database
	if err := s.repo.DeleteAttachment(ctx, id); err != nil {
		return fmt.Errorf("delete attachment record: %w", err)
	}

	s.publishEvent("attachment.deleted", attachment.EntityID, map[string]interface{}{
		"attachment_id": id,
		"entity_type":   attachment.EntityType,
	})

	return nil
}

func (s *AttachmentService) publishEvent(eventType, entityID string, payload map[string]interface{}) {
	if s.producer == nil {
		return
	}
	event := kafka.Event{Type: eventType, Timestamp: time.Now(), Payload: payload}
	if entityID != "" {
		payload["entity_id"] = entityID
	}
	_ = s.producer.PublishEvent("attachment-events", event)
}

// ValidateFile validates file size and type
func ValidateFile(size int64, contentType string, maxSize int64, allowedTypes []string) error {
	if size > maxSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", maxSize)
	}

	if len(allowedTypes) > 0 {
		allowed := false
		for _, t := range allowedTypes {
			if t == contentType {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("file type %s is not allowed", contentType)
		}
	}

	return nil
}

// ReadStream reads data from a stream
func ReadStream(stream io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, stream)
	if err != nil {
		return nil, fmt.Errorf("read stream: %w", err)
	}
	return buf.Bytes(), nil
}
