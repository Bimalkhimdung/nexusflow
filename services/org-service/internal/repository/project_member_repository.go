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
)

// ProjectMemberRepository handles project member data access
type ProjectMemberRepository struct {
	db  *database.DB
	log *logger.Logger
}

// NewProjectMemberRepository creates a new project member repository
func NewProjectMemberRepository(db *database.DB, log *logger.Logger) *ProjectMemberRepository {
	return &ProjectMemberRepository{
		db:  db,
		log: log,
	}
}

// AddMember adds a member to a project
func (r *ProjectMemberRepository) AddMember(ctx context.Context, member *models.ProjectMember) error {
	if member.ID == "" {
		member.ID = uuid.New().String()
	}
	
	_, err := r.db.NewInsert().Model(member).Exec(ctx)
	if err != nil {
		r.log.Sugar().Errorw("Failed to add project member", "error", err, "project_id", member.ProjectID, "user_id", member.UserID)
		return fmt.Errorf("add project member: %w", err)
	}
	
	return nil
}

// RemoveMember removes a member from a project
func (r *ProjectMemberRepository) RemoveMember(ctx context.Context, projectID, userID string) error {
	_, err := r.db.NewDelete().
		Model((*models.ProjectMember)(nil)).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Exec(ctx)
		
	if err != nil {
		r.log.Sugar().Errorw("Failed to remove project member", "error", err, "project_id", projectID, "user_id", userID)
		return fmt.Errorf("remove project member: %w", err)
	}
	
	return nil
}

// GetMember gets a project member by project ID and user ID
func (r *ProjectMemberRepository) GetMember(ctx context.Context, projectID, userID string) (*models.ProjectMember, error) {
	member := new(models.ProjectMember)
	err := r.db.NewSelect().
		Model(member).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Scan(ctx)
		
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.log.Sugar().Errorw("Failed to get project member", "error", err, "project_id", projectID, "user_id", userID)
		return nil, fmt.Errorf("get project member: %w", err)
	}
	
	return member, nil
}

// IsMember checks if a user is a member of a project
func (r *ProjectMemberRepository) IsMember(ctx context.Context, projectID, userID string) (bool, error) {
	member, err := r.GetMember(ctx, projectID, userID)
	if err != nil {
		return false, err
	}
	return member != nil, nil
}

// ListMembers lists members of a project
func (r *ProjectMemberRepository) ListMembers(ctx context.Context, projectID string, limit, offset int) ([]*models.ProjectMember, int, error) {
	var members []*models.ProjectMember
	
	count, err := r.db.NewSelect().
		Model(&members).
		Where("project_id = ?", projectID).
		Order("added_at DESC").
		Limit(limit).
		Offset(offset).
		ScanAndCount(ctx)
		
	if err != nil {
		r.log.Sugar().Errorw("Failed to list project members", "error", err, "project_id", projectID)
		return nil, 0, fmt.Errorf("list project members: %w", err)
	}
	
	return members, count, nil
}

// GetUserProjects gets all projects a user is a member of
func (r *ProjectMemberRepository) GetUserProjects(ctx context.Context, userID string) ([]string, error) {
	var projectIDs []string
	
	err := r.db.NewSelect().
		Model((*models.ProjectMember)(nil)).
		Column("project_id").
		Where("user_id = ?", userID).
		Scan(ctx, &projectIDs)
		
	if err != nil {
		r.log.Sugar().Errorw("Failed to get user projects", "error", err, "user_id", userID)
		return nil, fmt.Errorf("get user projects: %w", err)
	}
	
	return projectIDs, nil
}
