package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/nexusflow/nexusflow/pkg/database"
	"github.com/nexusflow/nexusflow/pkg/logger"
	"github.com/nexusflow/nexusflow/services/org-service/internal/models"
)

// InviteRepository handles invite data access
type InviteRepository struct {
	db  *database.DB
	log *logger.Logger
}

// NewInviteRepository creates a new invite repository
func NewInviteRepository(db *database.DB, log *logger.Logger) *InviteRepository {
	return &InviteRepository{
		db:  db,
		log: log,
	}
}

// Create creates a new invite
func (r *InviteRepository) Create(ctx context.Context, invite *models.Invite) error {
	if invite.ID == "" {
		invite.ID = uuid.New().String()
	}
	
	_, err := r.db.NewInsert().Model(invite).Exec(ctx)
	if err != nil {
		r.log.Sugar().Errorw("Failed to create invite", "error", err, "email", invite.Email)
		return fmt.Errorf("create invite: %w", err)
	}
	
	return nil
}

// GetByToken gets an invite by token
func (r *InviteRepository) GetByToken(ctx context.Context, token string) (*models.Invite, error) {
	invite := new(models.Invite)
	err := r.db.NewSelect().Model(invite).Where("token = ?", token).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.log.Sugar().Errorw("Failed to get invite by token", "error", err)
		return nil, fmt.Errorf("get invite by token: %w", err)
	}
	return invite, nil
}

// GetByID gets an invite by ID
func (r *InviteRepository) GetByID(ctx context.Context, id string) (*models.Invite, error) {
	invite := new(models.Invite)
	err := r.db.NewSelect().Model(invite).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.log.Sugar().Errorw("Failed to get invite by ID", "error", err, "id", id)
		return nil, fmt.Errorf("get invite: %w", err)
	}
	return invite, nil
}

// Update updates an invite
func (r *InviteRepository) Update(ctx context.Context, invite *models.Invite) error {
	_, err := r.db.NewUpdate().
		Model(invite).
		WherePK().
		Exec(ctx)
		
	if err != nil {
		r.log.Sugar().Errorw("Failed to update invite", "error", err, "id", invite.ID)
		return fmt.Errorf("update invite: %w", err)
	}
	
	return nil
}

// List lists invites for an organization
func (r *InviteRepository) List(ctx context.Context, orgID string, limit, offset int) ([]*models.Invite, int, error) {
	var invites []*models.Invite
	
	count, err := r.db.NewSelect().
		Model(&invites).
		Where("organization_id = ?", orgID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		ScanAndCount(ctx)
		
	if err != nil {
		r.log.Sugar().Errorw("Failed to list invites", "error", err, "org_id", orgID)
		return nil, 0, fmt.Errorf("list invites: %w", err)
	}
	
	return invites, count, nil
}
