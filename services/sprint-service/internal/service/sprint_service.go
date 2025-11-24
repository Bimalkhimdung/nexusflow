package service

import (
	"context"
	"fmt"
	"time"

	"github.com/nexusflow/nexusflow/pkg/kafka"
	"github.com/nexusflow/nexusflow/pkg/logger"
	pb "github.com/nexusflow/nexusflow/pkg/proto/sprint/v1"
	"github.com/nexusflow/nexusflow/services/sprint-service/internal/models"
	"github.com/nexusflow/nexusflow/services/sprint-service/internal/repository"
)

type SprintService struct {
	repo     *repository.SprintRepository
	producer *kafka.Producer
	log      *logger.Logger
}

func NewSprintService(repo *repository.SprintRepository, producer *kafka.Producer, log *logger.Logger) *SprintService {
	return &SprintService{repo: repo, producer: producer, log: log}
}

// CreateSprint creates a new sprint
func (s *SprintService) CreateSprint(ctx context.Context, req *pb.CreateSprintRequest) (*models.Sprint, error) {
	sprint := &models.Sprint{
		ProjectID: req.ProjectId,
		Name:      req.Name,
		Goal:      req.Goal,
		Status:    models.SprintStatusPlanned,
	}

	if req.StartDate != "" {
		startDate, err := time.Parse(time.RFC3339, req.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start date: %w", err)
		}
		sprint.StartDate = startDate
	}

	if req.EndDate != "" {
		endDate, err := time.Parse(time.RFC3339, req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date: %w", err)
		}
		sprint.EndDate = endDate
	}

	// Validate dates
	if !sprint.StartDate.IsZero() && !sprint.EndDate.IsZero() && sprint.StartDate.After(sprint.EndDate) {
		return nil, fmt.Errorf("start date must be before end date")
	}

	if err := s.repo.CreateSprint(ctx, sprint); err != nil {
		return nil, fmt.Errorf("create sprint: %w", err)
	}

	s.publishEvent("sprint.created", sprint.ProjectID, map[string]interface{}{
		"sprint_id": sprint.ID,
		"name":      sprint.Name,
	})

	return sprint, nil
}

func (s *SprintService) GetSprint(ctx context.Context, id string) (*models.Sprint, error) {
	return s.repo.GetSprint(ctx, id)
}

func (s *SprintService) ListSprints(ctx context.Context, projectID string, status models.SprintStatus) ([]*models.Sprint, error) {
	return s.repo.ListSprints(ctx, projectID, status)
}

func (s *SprintService) UpdateSprint(ctx context.Context, req *pb.UpdateSprintRequest) (*models.Sprint, error) {
	sprint, err := s.repo.GetSprint(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		sprint.Name = req.Name
	}
	if req.Goal != "" {
		sprint.Goal = req.Goal
	}
	if req.StartDate != "" {
		startDate, err := time.Parse(time.RFC3339, req.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start date: %w", err)
		}
		sprint.StartDate = startDate
	}
	if req.EndDate != "" {
		endDate, err := time.Parse(time.RFC3339, req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date: %w", err)
		}
		sprint.EndDate = endDate
	}

	if err := s.repo.UpdateSprint(ctx, sprint); err != nil {
		return nil, fmt.Errorf("update sprint: %w", err)
	}

	s.publishEvent("sprint.updated", sprint.ProjectID, map[string]interface{}{"sprint_id": sprint.ID})
	return sprint, nil
}

func (s *SprintService) DeleteSprint(ctx context.Context, id string) error {
	sprint, err := s.repo.GetSprint(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.DeleteSprint(ctx, id); err != nil {
		return fmt.Errorf("delete sprint: %w", err)
	}

	s.publishEvent("sprint.deleted", sprint.ProjectID, map[string]interface{}{"sprint_id": id})
	return nil
}

// Sprint lifecycle
func (s *SprintService) StartSprint(ctx context.Context, sprintID string) (*models.Sprint, error) {
	sprint, err := s.repo.GetSprint(ctx, sprintID)
	if err != nil {
		return nil, err
	}

	if sprint.Status != models.SprintStatusPlanned {
		return nil, fmt.Errorf("sprint must be in planned status to start")
	}

	// Check if project already has an active sprint
	hasActive, err := s.repo.HasActiveSprint(ctx, sprint.ProjectID)
	if err != nil {
		return nil, err
	}
	if hasActive {
		return nil, fmt.Errorf("project already has an active sprint")
	}

	sprint.Status = models.SprintStatusActive
	if err := s.repo.UpdateSprint(ctx, sprint); err != nil {
		return nil, fmt.Errorf("start sprint: %w", err)
	}

	s.publishEvent("sprint.started", sprint.ProjectID, map[string]interface{}{"sprint_id": sprint.ID})
	return sprint, nil
}

func (s *SprintService) CompleteSprint(ctx context.Context, sprintID string) (*models.Sprint, error) {
	sprint, err := s.repo.GetSprint(ctx, sprintID)
	if err != nil {
		return nil, err
	}

	if sprint.Status != models.SprintStatusActive {
		return nil, fmt.Errorf("sprint must be active to complete")
	}

	sprint.Status = models.SprintStatusCompleted
	if err := s.repo.UpdateSprint(ctx, sprint); err != nil {
		return nil, fmt.Errorf("complete sprint: %w", err)
	}

	s.publishEvent("sprint.completed", sprint.ProjectID, map[string]interface{}{"sprint_id": sprint.ID})
	return sprint, nil
}

// Sprint issues
func (s *SprintService) AddIssueToSprint(ctx context.Context, sprintID, issueID string) error {
	sprint, err := s.repo.GetSprint(ctx, sprintID)
	if err != nil {
		return err
	}

	if sprint.Status == models.SprintStatusCompleted {
		return fmt.Errorf("cannot add issues to completed sprint")
	}

	if err := s.repo.AddIssueToSprint(ctx, sprintID, issueID); err != nil {
		return fmt.Errorf("add issue to sprint: %w", err)
	}

	s.publishEvent("sprint.issue_added", sprint.ProjectID, map[string]interface{}{
		"sprint_id": sprintID,
		"issue_id":  issueID,
	})
	return nil
}

func (s *SprintService) RemoveIssueFromSprint(ctx context.Context, sprintID, issueID string) error {
	sprint, err := s.repo.GetSprint(ctx, sprintID)
	if err != nil {
		return err
	}

	if err := s.repo.RemoveIssueFromSprint(ctx, sprintID, issueID); err != nil {
		return fmt.Errorf("remove issue from sprint: %w", err)
	}

	s.publishEvent("sprint.issue_removed", sprint.ProjectID, map[string]interface{}{
		"sprint_id": sprintID,
		"issue_id":  issueID,
	})
	return nil
}

func (s *SprintService) GetSprintIssues(ctx context.Context, sprintID string) ([]string, error) {
	return s.repo.ListSprintIssues(ctx, sprintID)
}

func (s *SprintService) publishEvent(eventType, projectID string, payload map[string]interface{}) {
	if s.producer == nil {
		return
	}
	event := kafka.Event{Type: eventType, Timestamp: time.Now(), Payload: payload}
	if projectID != "" {
		payload["project_id"] = projectID
	}
	_ = s.producer.PublishEvent("sprint-events", event)
}
