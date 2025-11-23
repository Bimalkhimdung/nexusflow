package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nexusflow/nexusflow/pkg/database"
	"github.com/nexusflow/nexusflow/pkg/logger"
	"github.com/nexusflow/nexusflow/services/project-service/internal/models"
	"github.com/uptrace/bun"
)

// ProjectRepository handles project data access
type ProjectRepository struct {
	db  *database.DB
	log *logger.Logger
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(db *database.DB, log *logger.Logger) *ProjectRepository {
	return &ProjectRepository{
		db:  db,
		log: log,
	}
}

// Create creates a new project
func (r *ProjectRepository) Create(ctx context.Context, project *models.Project) error {
	if project.ID == "" {
		project.ID = uuid.New().String()
	}
	
	_, err := r.db.NewInsert().Model(project).Exec(ctx)
	if err != nil {
		r.log.Sugar().Errorw("Failed to create project", "error", err, "key", project.Key)
		return fmt.Errorf("create project: %w", err)
	}
	
	return nil
}

// GetByID gets a project by ID
func (r *ProjectRepository) GetByID(ctx context.Context, id string) (*models.Project, error) {
	project := new(models.Project)
	err := r.db.NewSelect().Model(project).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.log.Sugar().Errorw("Failed to get project by ID", "error", err, "id", id)
		return nil, fmt.Errorf("get project: %w", err)
	}
	return project, nil
}

// GetByKey gets a project by key and org ID
func (r *ProjectRepository) GetByKey(ctx context.Context, orgID, key string) (*models.Project, error) {
	project := new(models.Project)
	err := r.db.NewSelect().
		Model(project).
		Where("organization_id = ? AND key = ?", orgID, key).
		Scan(ctx)
		
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.log.Sugar().Errorw("Failed to get project by key", "error", err, "key", key)
		return nil, fmt.Errorf("get project by key: %w", err)
	}
	return project, nil
}

// Update updates a project
func (r *ProjectRepository) Update(ctx context.Context, project *models.Project) error {
	project.UpdatedAt = time.Now()
	
	_, err := r.db.NewUpdate().
		Model(project).
		WherePK().
		Exec(ctx)
		
	if err != nil {
		r.log.Sugar().Errorw("Failed to update project", "error", err, "id", project.ID)
		return fmt.Errorf("update project: %w", err)
	}
	
	return nil
}

// Delete soft deletes a project
func (r *ProjectRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().
		Model((*models.Project)(nil)).
		Where("id = ?", id).
		Exec(ctx)
		
	if err != nil {
		r.log.Sugar().Errorw("Failed to delete project", "error", err, "id", id)
		return fmt.Errorf("delete project: %w", err)
	}
	
	return nil
}

// List lists projects for an organization
func (r *ProjectRepository) List(ctx context.Context, orgID string, limit, offset int) ([]*models.Project, int, error) {
	var projects []*models.Project
	
	count, err := r.db.NewSelect().
		Model(&projects).
		Where("organization_id = ?", orgID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		ScanAndCount(ctx)
		
	if err != nil {
		r.log.Sugar().Errorw("Failed to list projects", "error", err, "org_id", orgID)
		return nil, 0, fmt.Errorf("list projects: %w", err)
	}
	
	return projects, count, nil
}

// AddMember adds a member to a project
func (r *ProjectRepository) AddMember(ctx context.Context, member *models.ProjectMember) error {
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
func (r *ProjectRepository) RemoveMember(ctx context.Context, projectID, userID string) error {
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

// UpdateMemberRole updates a member's role
func (r *ProjectRepository) UpdateMemberRole(ctx context.Context, projectID, userID string, role models.ProjectRole) error {
	_, err := r.db.NewUpdate().
		Model((*models.ProjectMember)(nil)).
		Set("role = ?", role).
		Set("updated_at = ?", time.Now()).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Exec(ctx)
		
	if err != nil {
		r.log.Sugar().Errorw("Failed to update project member role", "error", err, "project_id", projectID, "user_id", userID)
		return fmt.Errorf("update project member role: %w", err)
	}
	
	return nil
}

// GetMember gets a member by project ID and user ID
func (r *ProjectRepository) GetMember(ctx context.Context, projectID, userID string) (*models.ProjectMember, error) {
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

// ListMembers lists members of a project
func (r *ProjectRepository) ListMembers(ctx context.Context, projectID string, limit, offset int) ([]*models.ProjectMember, int, error) {
	var members []*models.ProjectMember
	
	count, err := r.db.NewSelect().
		Model(&members).
		Where("project_id = ?", projectID).
		Order("joined_at DESC").
		Limit(limit).
		Offset(offset).
		ScanAndCount(ctx)
		
	if err != nil {
		r.log.Sugar().Errorw("Failed to list project members", "error", err, "project_id", projectID)
		return nil, 0, fmt.Errorf("list project members: %w", err)
	}
	
	return members, count, nil
}

// CreateWithMember creates a project and adds the creator as admin in a transaction
func (r *ProjectRepository) CreateWithMember(ctx context.Context, project *models.Project, userID string) error {
	return r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		if project.ID == "" {
			project.ID = uuid.New().String()
		}
		
		// Create project
		if _, err := tx.NewInsert().Model(project).Exec(ctx); err != nil {
			return fmt.Errorf("create project: %w", err)
		}
		
		// Add member as admin
		member := &models.ProjectMember{
			ID:        uuid.New().String(),
			ProjectID: project.ID,
			UserID:    userID,
			Role:      models.ProjectRoleAdmin,
			JoinedAt:  time.Now(),
			UpdatedAt: time.Now(),
		}
		
		if _, err := tx.NewInsert().Model(member).Exec(ctx); err != nil {
			return fmt.Errorf("add admin: %w", err)
		}
		
		return nil
	})
}
