package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/nexusflow/nexusflow/pkg/database"
	"github.com/nexusflow/nexusflow/pkg/logger"
	"github.com/nexusflow/nexusflow/services/workflow-service/internal/models"
)

// WorkflowRepository handles workflow data access
type WorkflowRepository struct {
	db  *database.DB
	log *logger.Logger
}

// NewWorkflowRepository creates a new workflow repository
func NewWorkflowRepository(db *database.DB, log *logger.Logger) *WorkflowRepository {
	return &WorkflowRepository{
		db:  db,
		log: log,
	}
}

// CreateWorkflow creates a new workflow
func (r *WorkflowRepository) CreateWorkflow(ctx context.Context, workflow *models.Workflow) error {
	_, err := r.db.NewInsert().Model(workflow).Exec(ctx)
	if err != nil {
		return fmt.Errorf("create workflow: %w", err)
	}
	return nil
}

// GetWorkflow gets a workflow by ID
func (r *WorkflowRepository) GetWorkflow(ctx context.Context, id string) (*models.Workflow, error) {
	workflow := new(models.Workflow)
	err := r.db.NewSelect().Model(workflow).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get workflow: %w", err)
	}
	return workflow, nil
}

// ListWorkflows lists workflows for a project
func (r *WorkflowRepository) ListWorkflows(ctx context.Context, projectID string) ([]*models.Workflow, error) {
	var workflows []*models.Workflow
	err := r.db.NewSelect().
		Model(&workflows).
		Where("project_id = ?", projectID).
		Order("created_at DESC").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("list workflows: %w", err)
	}
	return workflows, nil
}

// CreateStatus creates a new status
func (r *WorkflowRepository) CreateStatus(ctx context.Context, status *models.WorkflowStatus) error {
	_, err := r.db.NewInsert().Model(status).Exec(ctx)
	if err != nil {
		return fmt.Errorf("create status: %w", err)
	}
	return nil
}

// ListStatuses lists statuses for a workflow
func (r *WorkflowRepository) ListStatuses(ctx context.Context, workflowID string) ([]*models.WorkflowStatus, error) {
	var statuses []*models.WorkflowStatus
	err := r.db.NewSelect().
		Model(&statuses).
		Where("workflow_id = ?", workflowID).
		Order("position ASC").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("list statuses: %w", err)
	}
	return statuses, nil
}

// CreateTransition creates a new transition
func (r *WorkflowRepository) CreateTransition(ctx context.Context, transition *models.WorkflowTransition) error {
	_, err := r.db.NewInsert().Model(transition).Exec(ctx)
	if err != nil {
		return fmt.Errorf("create transition: %w", err)
	}
	return nil
}

// ListTransitions lists transitions for a workflow
func (r *WorkflowRepository) ListTransitions(ctx context.Context, workflowID string) ([]*models.WorkflowTransition, error) {
	var transitions []*models.WorkflowTransition
	err := r.db.NewSelect().
		Model(&transitions).
		Where("workflow_id = ?", workflowID).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("list transitions: %w", err)
	}
	return transitions, nil
}

// GetTransition gets a transition by ID
func (r *WorkflowRepository) GetTransition(ctx context.Context, id string) (*models.WorkflowTransition, error) {
	transition := new(models.WorkflowTransition)
	err := r.db.NewSelect().Model(transition).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get transition: %w", err)
	}
	return transition, nil
}

// GetStatus gets a status by ID
func (r *WorkflowRepository) GetStatus(ctx context.Context, id string) (*models.WorkflowStatus, error) {
	status := new(models.WorkflowStatus)
	err := r.db.NewSelect().Model(status).Where("id = ?", id).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get status: %w", err)
	}
	return status, nil
}
