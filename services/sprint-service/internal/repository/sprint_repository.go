package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/nexusflow/nexusflow/pkg/database"
	"github.com/nexusflow/nexusflow/pkg/logger"
	"github.com/nexusflow/nexusflow/services/sprint-service/internal/models"
)

type SprintRepository struct {
	db  *database.DB
	log *logger.Logger
}

func NewSprintRepository(db *database.DB, log *logger.Logger) *SprintRepository {
	return &SprintRepository{db: db, log: log}
}

// Sprint CRUD
func (r *SprintRepository) CreateSprint(ctx context.Context, sprint *models.Sprint) error {
	sprint.ID = ""
	sprint.CreatedAt = time.Now()
	sprint.UpdatedAt = time.Now()
	_, err := r.db.NewInsert().Model(sprint).Exec(ctx)
	if err != nil {
		return fmt.Errorf("create sprint: %w", err)
	}
	return nil
}

func (r *SprintRepository) GetSprint(ctx context.Context, id string) (*models.Sprint, error) {
	s := new(models.Sprint)
	err := r.db.NewSelect().Model(s).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("get sprint: %w", err)
	}
	return s, nil
}

func (r *SprintRepository) ListSprints(ctx context.Context, projectID string, status models.SprintStatus) ([]*models.Sprint, error) {
	var sprints []*models.Sprint
	query := r.db.NewSelect().Model(&sprints).Where("project_id = ?", projectID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("created_at DESC").Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("list sprints: %w", err)
	}
	return sprints, nil
}

func (r *SprintRepository) UpdateSprint(ctx context.Context, sprint *models.Sprint) error {
	sprint.UpdatedAt = time.Now()
	_, err := r.db.NewUpdate().Model(sprint).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("update sprint: %w", err)
	}
	return nil
}

func (r *SprintRepository) DeleteSprint(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().Model((*models.Sprint)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("delete sprint: %w", err)
	}
	return nil
}

// Sprint Issues
func (r *SprintRepository) AddIssueToSprint(ctx context.Context, sprintID, issueID string) error {
	si := &models.SprintIssue{
		SprintID: sprintID,
		IssueID:  issueID,
		AddedAt:  time.Now(),
	}
	_, err := r.db.NewInsert().Model(si).Exec(ctx)
	if err != nil {
		return fmt.Errorf("add issue to sprint: %w", err)
	}
	return nil
}

func (r *SprintRepository) RemoveIssueFromSprint(ctx context.Context, sprintID, issueID string) error {
	_, err := r.db.NewDelete().Model((*models.SprintIssue)(nil)).
		Where("sprint_id = ? AND issue_id = ?", sprintID, issueID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("remove issue from sprint: %w", err)
	}
	return nil
}

func (r *SprintRepository) ListSprintIssues(ctx context.Context, sprintID string) ([]string, error) {
	var issues []models.SprintIssue
	err := r.db.NewSelect().Model(&issues).Where("sprint_id = ?", sprintID).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("list sprint issues: %w", err)
	}
	
	issueIDs := make([]string, len(issues))
	for i, issue := range issues {
		issueIDs[i] = issue.IssueID
	}
	return issueIDs, nil
}

// Check if project has active sprint
func (r *SprintRepository) HasActiveSprint(ctx context.Context, projectID string) (bool, error) {
	count, err := r.db.NewSelect().Model((*models.Sprint)(nil)).
		Where("project_id = ? AND status = ?", projectID, models.SprintStatusActive).
		Count(ctx)
	if err != nil {
		return false, fmt.Errorf("check active sprint: %w", err)
	}
	return count > 0, nil
}
