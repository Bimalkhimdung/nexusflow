package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/nexusflow/nexusflow/pkg/database"
	"github.com/nexusflow/nexusflow/pkg/logger"
	"github.com/nexusflow/nexusflow/services/attachment-service/internal/models"
)

type AttachmentRepository struct {
	db  *database.DB
	log *logger.Logger
}

func NewAttachmentRepository(db *database.DB, log *logger.Logger) *AttachmentRepository {
	return &AttachmentRepository{db: db, log: log}
}

func (r *AttachmentRepository) CreateAttachment(ctx context.Context, attachment *models.Attachment) error {
	attachment.ID = ""
	attachment.CreatedAt = time.Now()
	_, err := r.db.NewInsert().Model(attachment).Exec(ctx)
	if err != nil {
		return fmt.Errorf("create attachment: %w", err)
	}
	return nil
}

func (r *AttachmentRepository) GetAttachment(ctx context.Context, id string) (*models.Attachment, error) {
	a := new(models.Attachment)
	err := r.db.NewSelect().Model(a).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("get attachment: %w", err)
	}
	return a, nil
}

func (r *AttachmentRepository) ListAttachments(ctx context.Context, entityType, entityID string) ([]*models.Attachment, error) {
	var attachments []*models.Attachment
	err := r.db.NewSelect().Model(&attachments).
		Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Order("created_at DESC").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("list attachments: %w", err)
	}
	return attachments, nil
}

func (r *AttachmentRepository) DeleteAttachment(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().Model((*models.Attachment)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("delete attachment: %w", err)
	}
	return nil
}

func (r *AttachmentRepository) GetAttachmentsByUploader(ctx context.Context, uploaderID string) ([]*models.Attachment, error) {
	var attachments []*models.Attachment
	err := r.db.NewSelect().Model(&attachments).
		Where("uploader_id = ?", uploaderID).
		Order("created_at DESC").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("get attachments by uploader: %w", err)
	}
	return attachments, nil
}
