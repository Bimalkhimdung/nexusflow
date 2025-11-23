package service

import (
	"context"
	"fmt"

	"github.com/nexusflow/nexusflow/pkg/logger"
	pb "github.com/nexusflow/nexusflow/pkg/proto/issue/v1"
	"github.com/nexusflow/nexusflow/services/workflow-service/internal/models"
	"github.com/nexusflow/nexusflow/services/workflow-service/internal/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// WorkflowService handles workflow business logic
type WorkflowService struct {
	repo        *repository.WorkflowRepository
	log         *logger.Logger
	issueClient pb.IssueServiceClient
}

// NewWorkflowService creates a new workflow service
func NewWorkflowService(
	repo *repository.WorkflowRepository,
	log *logger.Logger,
	issueServiceAddr string,
) (*WorkflowService, error) {
	// Connect to issue service
	conn, err := grpc.Dial(issueServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to issue service: %w", err)
	}
	issueClient := pb.NewIssueServiceClient(conn)

	return &WorkflowService{
		repo:        repo,
		log:         log,
		issueClient: issueClient,
	}, nil
}

// CreateWorkflow creates a new workflow
func (s *WorkflowService) CreateWorkflow(ctx context.Context, workflow *models.Workflow) (*models.Workflow, error) {
	if err := s.repo.CreateWorkflow(ctx, workflow); err != nil {
		return nil, err
	}
	return workflow, nil
}

// ListWorkflows lists workflows
func (s *WorkflowService) ListWorkflows(ctx context.Context, projectID string) ([]*models.Workflow, error) {
	return s.repo.ListWorkflows(ctx, projectID)
}

// CreateStatus creates a new status
func (s *WorkflowService) CreateStatus(ctx context.Context, status *models.WorkflowStatus) (*models.WorkflowStatus, error) {
	if err := s.repo.CreateStatus(ctx, status); err != nil {
		return nil, err
	}
	return status, nil
}

// CreateTransition creates a new transition
func (s *WorkflowService) CreateTransition(ctx context.Context, transition *models.WorkflowTransition) (*models.WorkflowTransition, error) {
	if err := s.repo.CreateTransition(ctx, transition); err != nil {
		return nil, err
	}
	return transition, nil
}

// ExecuteTransition executes a transition
func (s *WorkflowService) ExecuteTransition(ctx context.Context, issueID, transitionID, userID string) error {
	// 1. Get Transition
	transition, err := s.repo.GetTransition(ctx, transitionID)
	if err != nil {
		return fmt.Errorf("failed to get transition: %w", err)
	}
	if transition == nil {
		return fmt.Errorf("transition not found")
	}

	// 2. Validate Rules (TODO: Implement rule validation)
	// e.g., check permissions, conditions

	// 3. Update Issue Status in Issue Service
	// We need to get the target status ID.
	// But wait, Issue Service expects a status ID.
	// Does Issue Service know about Workflow Statuses?
	// Currently Issue Service has a status_id field.
	// So we just pass the ToStatusID.
	
	_, err = s.issueClient.UpdateIssue(ctx, &pb.UpdateIssueRequest{
		Id:       issueID,
		StatusId: &transition.ToStatusID,
	})
	if err != nil {
		return fmt.Errorf("failed to update issue status: %w", err)
	}

	// 4. Execute Post Functions (TODO: Implement post functions)
	// e.g., assign to user, add comment

	return nil
}

// GetAvailableTransitions gets available transitions for an issue
func (s *WorkflowService) GetAvailableTransitions(ctx context.Context, issueID string) ([]*models.WorkflowTransition, error) {
	// 1. Get Issue to find current status
	issueResp, err := s.issueClient.GetIssue(ctx, &pb.GetIssueRequest{Id: issueID})
	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}
	currentStatusID := issueResp.Issue.StatusId

	// 2. Get Workflow ID for the project (TODO: Need mapping from project to workflow)
	// For now, assume we can find transitions by FromStatusID?
	// But transitions are scoped to a workflow.
	// We need to know which workflow applies to this issue/project.
	// Let's assume for MVP we just search all transitions where FromStatusID matches.
	// This is inefficient but works if status IDs are unique globally (UUIDs).
	// Ideally: Project -> Workflow -> Transitions.
	
	// Wait, we don't have a method to list transitions by FromStatusID across all workflows?
	// Or we assume we know the workflow.
	// Let's implement a simple query in repo: Find transitions by FromStatusID.
	// But repo.ListTransitions takes workflowID.
	// We need a new repo method or change approach.
	
	// Better approach:
	// 1. Get Issue -> ProjectID.
	// 2. Get Workflow for ProjectID.
	// 3. List Transitions for Workflow.
	// 4. Filter by FromStatusID == CurrentStatusID.
	
	workflows, err := s.repo.ListWorkflows(ctx, issueResp.Issue.ProjectId)
	if err != nil {
		return nil, fmt.Errorf("failed to list workflows: %w", err)
	}
	if len(workflows) == 0 {
		return nil, fmt.Errorf("no workflow found for project")
	}
	// Assume first workflow is active
	workflow := workflows[0]
	
	transitions, err := s.repo.ListTransitions(ctx, workflow.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list transitions: %w", err)
	}
	
	var available []*models.WorkflowTransition
	for _, t := range transitions {
		if t.FromStatusID == currentStatusID {
			available = append(available, t)
		}
	}
	
	return available, nil
}
