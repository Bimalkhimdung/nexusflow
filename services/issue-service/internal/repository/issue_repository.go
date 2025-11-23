package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nexusflow/nexusflow/pkg/database"
	"github.com/nexusflow/nexusflow/pkg/logger"
	"github.com/nexusflow/nexusflow/services/issue-service/internal/models"
	"github.com/uptrace/bun"
)

// IssueRepository handles issue data access
type IssueRepository struct {
	db  *database.DB
	log *logger.Logger
}

// NewIssueRepository creates a new issue repository
func NewIssueRepository(db *database.DB, log *logger.Logger) *IssueRepository {
	return &IssueRepository{
		db:  db,
		log: log,
	}
}

// Create creates a new issue with atomic key generation
func (r *IssueRepository) Create(ctx context.Context, issue *models.Issue, projectKey string) error {
	return r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		// 1. Get and increment project counter
		counter := &models.ProjectCounter{ProjectID: issue.ProjectID}
		
		// Upsert counter: increment if exists, insert 1 if not
		_, err := tx.NewInsert().
			Model(counter).
			On("CONFLICT (project_id) DO UPDATE").
			Set("next_issue_number = pc.next_issue_number + 1").
			Returning("next_issue_number").
			Exec(ctx)
			
		if err != nil {
			return fmt.Errorf("failed to increment project counter: %w", err)
		}
		
		// 2. Generate key
		// Note: The returned counter value is the *next* number if we used RETURNING, 
		// but Bun's behavior with ON CONFLICT UPDATE RETURNING can be tricky.
		// Let's fetch the updated value to be safe, or trust the return.
		// Actually, let's use a simpler approach: 
		// INSERT ... ON CONFLICT DO UPDATE ... RETURNING next_issue_number
		// But we need the value *before* increment if we want 1, 2, 3... 
		// Wait, if default is 1, and we increment, we get 2. So first issue is 1?
		// Let's adjust: Default 1. First insert uses 1. Next update increments to 2.
		// Actually, let's just read-and-update with locking for simplicity and correctness.
		
		// Better approach with atomic update:
		// UPDATE project_counters SET next_issue_number = next_issue_number + 1 WHERE project_id = ? RETURNING next_issue_number
		// If no rows, INSERT.
		
		var issueNum int64
		
		// Try update first
		err = tx.NewRaw(`
			UPDATE project_counters 
			SET next_issue_number = next_issue_number + 1, updated_at = now() 
			WHERE project_id = ? 
			RETURNING next_issue_number - 1`, issue.ProjectID).Scan(ctx, &issueNum)
			
		if err != nil {
			if err == sql.ErrNoRows {
				// Create counter if not exists
				counter = &models.ProjectCounter{
					ProjectID:       issue.ProjectID,
					NextIssueNumber: 2, // Next will be 2
				}
				if _, err := tx.NewInsert().Model(counter).Exec(ctx); err != nil {
					return fmt.Errorf("failed to create project counter: %w", err)
				}
				issueNum = 1
			} else {
				return fmt.Errorf("failed to update project counter: %w", err)
			}
		}
		
		// Set key
		issue.Key = fmt.Sprintf("%s-%d", projectKey, issueNum)
		
		if issue.ID == "" {
			issue.ID = uuid.New().String()
		}
		
		// 3. Create issue
		if _, err := tx.NewInsert().Model(issue).Exec(ctx); err != nil {
			return fmt.Errorf("failed to create issue: %w", err)
		}
		
		return nil
	})
}

// GetByID gets an issue by ID
func (r *IssueRepository) GetByID(ctx context.Context, id string) (*models.Issue, error) {
	issue := new(models.Issue)
	err := r.db.NewSelect().Model(issue).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get issue: %w", err)
	}
	return issue, nil
}

// GetByKey gets an issue by key
func (r *IssueRepository) GetByKey(ctx context.Context, key string) (*models.Issue, error) {
	issue := new(models.Issue)
	err := r.db.NewSelect().Model(issue).Where("key = ?", key).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get issue by key: %w", err)
	}
	return issue, nil
}

// Update updates an issue
func (r *IssueRepository) Update(ctx context.Context, issue *models.Issue) error {
	issue.UpdatedAt = time.Now()
	_, err := r.db.NewUpdate().Model(issue).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("update issue: %w", err)
	}
	return nil
}

// Delete soft deletes an issue
func (r *IssueRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().Model((*models.Issue)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("delete issue: %w", err)
	}
	return nil
}

// List lists issues
func (r *IssueRepository) List(ctx context.Context, projectID string, limit, offset int) ([]*models.Issue, int, error) {
	var issues []*models.Issue
	count, err := r.db.NewSelect().
		Model(&issues).
		Where("project_id = ?", projectID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("list issues: %w", err)
	}
	return issues, count, nil
}
