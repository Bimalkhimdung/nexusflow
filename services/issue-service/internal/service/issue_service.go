package service

import (
	"context"
	"fmt"
	"time"

	"github.com/nexusflow/nexusflow/pkg/kafka"
	"github.com/nexusflow/nexusflow/pkg/logger"
	pb "github.com/nexusflow/nexusflow/pkg/proto/project/v1"
	"github.com/nexusflow/nexusflow/services/issue-service/internal/models"
	"github.com/nexusflow/nexusflow/services/issue-service/internal/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// IssueService handles issue business logic
type IssueService struct {
	repo          *repository.IssueRepository
	producer      *kafka.Producer
	log           *logger.Logger
	projectClient pb.ProjectServiceClient
}

// NewIssueService creates a new issue service
func NewIssueService(
	repo *repository.IssueRepository,
	producer *kafka.Producer,
	log *logger.Logger,
	projectServiceAddr string,
) (*IssueService, error) {
	// Connect to project service
	conn, err := grpc.Dial(projectServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to project service: %w", err)
	}
	projectClient := pb.NewProjectServiceClient(conn)

	return &IssueService{
		repo:          repo,
		producer:      producer,
		log:           log,
		projectClient: projectClient,
	}, nil
}

// CreateIssueInput represents input for creating an issue
type CreateIssueInput struct {
	ProjectID   string
	Summary     string
	Description string
	Type        models.IssueType
	Priority    models.IssuePriority
	AssigneeID  string
	ReporterID  string
	ParentID    string
}

// CreateIssue creates a new issue
func (s *IssueService) CreateIssue(ctx context.Context, input CreateIssueInput) (*models.Issue, error) {
	// 1. Get Project Key from Project Service
	projectResp, err := s.projectClient.GetProject(ctx, &pb.GetProjectRequest{Id: input.ProjectID})
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	projectKey := projectResp.Project.Key

	// 2. Create Issue
	issue := &models.Issue{
		ProjectID:   input.ProjectID,
		Summary:     input.Summary,
		Description: input.Description,
		Type:        input.Type,
		Priority:    input.Priority,
		AssigneeID:  input.AssigneeID,
		ReporterID:  input.ReporterID,
		ParentID:    input.ParentID,
	}
	if issue.Type == "" {
		issue.Type = models.IssueTypeTask
	}
	if issue.Priority == "" {
		issue.Priority = models.IssuePriorityMedium
	}

	if err := s.repo.Create(ctx, issue, projectKey); err != nil {
		return nil, fmt.Errorf("failed to create issue: %w", err)
	}

	// 3. Publish Event
	s.publishEvent("issue.created", input.ProjectID, input.ReporterID, map[string]interface{}{
		"issue_id": issue.ID,
		"key":      issue.Key,
		"summary":  issue.Summary,
	})

	return issue, nil
}

// GetIssue gets an issue by ID
func (s *IssueService) GetIssue(ctx context.Context, id string) (*models.Issue, error) {
	return s.repo.GetByID(ctx, id)
}

// GetIssueByKey gets an issue by key
func (s *IssueService) GetIssueByKey(ctx context.Context, key string) (*models.Issue, error) {
	return s.repo.GetByKey(ctx, key)
}

// ListIssues lists issues
func (s *IssueService) ListIssues(ctx context.Context, projectID string, page, pageSize int) ([]*models.Issue, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	return s.repo.List(ctx, projectID, pageSize, offset)
}

// publishEvent publishes a Kafka event
func (s *IssueService) publishEvent(eventType, projectID, userID string, payload map[string]interface{}) {
	if s.producer == nil {
		return
	}

	event := kafka.Event{
		Type:      eventType,
		UserID:    userID,
		Timestamp: time.Now(),
		Payload:   payload,
		// TODO: Add ProjectID to event struct if needed, or put in payload
	}
	// Hack: Add project_id to payload for now as Event struct might not have it top-level
	payload["project_id"] = projectID

	if err := s.producer.PublishEvent("issue-events", event); err != nil {
		s.log.Sugar().Errorw("Failed to publish event", "error", err, "type", eventType)
	}
}
