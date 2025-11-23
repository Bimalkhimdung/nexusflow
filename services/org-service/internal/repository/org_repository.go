package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nexusflow/nexusflow/pkg/database"
	"github.com/nexusflow/nexusflow/pkg/logger"
	"github.com/nexusflow/nexusflow/services/org-service/internal/models"
	"github.com/uptrace/bun"
)

// OrgRepository handles organization data access
type OrgRepository struct {
	db  *database.DB
	log *logger.Logger
}

// NewOrgRepository creates a new organization repository
func NewOrgRepository(db *database.DB, log *logger.Logger) *OrgRepository {
	return &OrgRepository{
		db:  db,
		log: log,
	}
}

// Create creates a new organization
func (r *OrgRepository) Create(ctx context.Context, org *models.Organization) error {
	if org.ID == "" {
		org.ID = uuid.New().String()
	}
	
	_, err := r.db.NewInsert().Model(org).Exec(ctx)
	if err != nil {
		r.log.Sugar().Errorw("Failed to create organization", "error", err, "slug", org.Slug)
		return fmt.Errorf("create organization: %w", err)
	}
	
	return nil
}

// GetByID gets an organization by ID
func (r *OrgRepository) GetByID(ctx context.Context, id string) (*models.Organization, error) {
	org := new(models.Organization)
	err := r.db.NewSelect().Model(org).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.log.Sugar().Errorw("Failed to get organization by ID", "error", err, "id", id)
		return nil, fmt.Errorf("get organization: %w", err)
	}
	return org, nil
}

// GetBySlug gets an organization by slug
func (r *OrgRepository) GetBySlug(ctx context.Context, slug string) (*models.Organization, error) {
	org := new(models.Organization)
	err := r.db.NewSelect().Model(org).Where("slug = ?", slug).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.log.Sugar().Errorw("Failed to get organization by slug", "error", err, "slug", slug)
		return nil, fmt.Errorf("get organization by slug: %w", err)
	}
	return org, nil
}

// Update updates an organization
func (r *OrgRepository) Update(ctx context.Context, org *models.Organization) error {
	org.UpdatedAt = time.Now()
	
	_, err := r.db.NewUpdate().
		Model(org).
		WherePK().
		Exec(ctx)
		
	if err != nil {
		r.log.Sugar().Errorw("Failed to update organization", "error", err, "id", org.ID)
		return fmt.Errorf("update organization: %w", err)
	}
	
	return nil
}

// Delete soft deletes an organization
func (r *OrgRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().
		Model((*models.Organization)(nil)).
		Where("id = ?", id).
		Exec(ctx)
		
	if err != nil {
		r.log.Sugar().Errorw("Failed to delete organization", "error", err, "id", id)
		return fmt.Errorf("delete organization: %w", err)
	}
	
	return nil
}

// List lists organizations for a user
func (r *OrgRepository) List(ctx context.Context, userID string, limit, offset int) ([]*models.Organization, int, error) {
	var orgs []*models.Organization
	
	// Join with org_members to filter by user
	q := r.db.NewSelect().
		Model(&orgs).
		Join("JOIN org_members AS om ON om.organization_id = o.id").
		Where("om.user_id = ?", userID).
		Order("o.created_at DESC").
		Limit(limit).
		Offset(offset)
		
	count, err := q.ScanAndCount(ctx)
	if err != nil {
		r.log.Sugar().Errorw("Failed to list organizations", "error", err, "user_id", userID)
		return nil, 0, fmt.Errorf("list organizations: %w", err)
	}
	
	return orgs, count, nil
}

// AddMember adds a member to an organization
func (r *OrgRepository) AddMember(ctx context.Context, member *models.OrgMember) error {
	if member.ID == "" {
		member.ID = uuid.New().String()
	}
	
	_, err := r.db.NewInsert().Model(member).Exec(ctx)
	if err != nil {
		r.log.Sugar().Errorw("Failed to add member", "error", err, "org_id", member.OrganizationID, "user_id", member.UserID)
		return fmt.Errorf("add member: %w", err)
	}
	
	return nil
}

// RemoveMember removes a member from an organization
func (r *OrgRepository) RemoveMember(ctx context.Context, orgID, userID string) error {
	_, err := r.db.NewDelete().
		Model((*models.OrgMember)(nil)).
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		Exec(ctx)
		
	if err != nil {
		r.log.Sugar().Errorw("Failed to remove member", "error", err, "org_id", orgID, "user_id", userID)
		return fmt.Errorf("remove member: %w", err)
	}
	
	return nil
}

// UpdateMemberRole updates a member's role
func (r *OrgRepository) UpdateMemberRole(ctx context.Context, orgID, userID string, role models.OrgRole) error {
	_, err := r.db.NewUpdate().
		Model((*models.OrgMember)(nil)).
		Set("role = ?", role).
		Set("updated_at = ?", time.Now()).
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		Exec(ctx)
		
	if err != nil {
		r.log.Sugar().Errorw("Failed to update member role", "error", err, "org_id", orgID, "user_id", userID)
		return fmt.Errorf("update member role: %w", err)
	}
	
	return nil
}

// GetMember gets a member by org ID and user ID
func (r *OrgRepository) GetMember(ctx context.Context, orgID, userID string) (*models.OrgMember, error) {
	member := new(models.OrgMember)
	err := r.db.NewSelect().
		Model(member).
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		Scan(ctx)
		
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.log.Sugar().Errorw("Failed to get member", "error", err, "org_id", orgID, "user_id", userID)
		return nil, fmt.Errorf("get member: %w", err)
	}
	
	return member, nil
}

// ListMembers lists members of an organization
func (r *OrgRepository) ListMembers(ctx context.Context, orgID string, limit, offset int) ([]*models.OrgMember, int, error) {
	var members []*models.OrgMember
	
	count, err := r.db.NewSelect().
		Model(&members).
		Where("organization_id = ?", orgID).
		Order("joined_at DESC").
		Limit(limit).
		Offset(offset).
		ScanAndCount(ctx)
		
	if err != nil {
		r.log.Sugar().Errorw("Failed to list members", "error", err, "org_id", orgID)
		return nil, 0, fmt.Errorf("list members: %w", err)
	}
	
	return members, count, nil
}

// CreateWithMember creates an organization and adds the creator as owner in a transaction
func (r *OrgRepository) CreateWithMember(ctx context.Context, org *models.Organization, userID string) error {
	return r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if org.ID == "" {
			org.ID = uuid.New().String()
		}
		
		// Create organization
		if _, err := tx.NewInsert().Model(org).Exec(ctx); err != nil {
			return fmt.Errorf("create org: %w", err)
		}
		
		// Add member as owner
		member := &models.OrgMember{
			ID:             uuid.New().String(),
			OrganizationID: org.ID,
			UserID:         userID,
			Role:           models.OrgRoleOwner,
			JoinedAt:       time.Now(),
			UpdatedAt:      time.Now(),
		}
		
		if _, err := tx.NewInsert().Model(member).Exec(ctx); err != nil {
			return fmt.Errorf("add owner: %w", err)
		}
		
		return nil
	})
}
