package service

import (
	"context"
	"fmt"
	"time"

	"github.com/nexusflow/nexusflow/pkg/kafka"
	"github.com/nexusflow/nexusflow/pkg/logger"
	"github.com/nexusflow/nexusflow/services/project-service/internal/models"
	"github.com/nexusflow/nexusflow/services/project-service/internal/repository"
)

// ProjectService handles project business logic
type ProjectService struct {
	repo     *repository.ProjectRepository
	producer *kafka.Producer
	log      *logger.Logger
}

// NewProjectService creates a new project service
func NewProjectService(
	repo *repository.ProjectRepository,
	producer *kafka.Producer,
	log *logger.Logger,
) *ProjectService {
	return &ProjectService{
		repo:     repo,
		producer: producer,
		log:      log,
	}
}

// CreateProjectInput represents input for creating a project
type CreateProjectInput struct {
	OrganizationID string
	Key            string
	Name           string
	Description    string
	Type           models.ProjectType
	LeadID         string
	UserID         string // Creator
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(ctx context.Context, input CreateProjectInput) (*models.Project, error) {
	// Validate input
	if input.OrganizationID == "" {
		return nil, fmt.Errorf("organization_id is required")
	}
	if input.Key == "" {
		return nil, fmt.Errorf("key is required")
	}
	if input.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if input.UserID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	// Check if key exists in org
	existing, _ := s.repo.GetByKey(ctx, input.OrganizationID, input.Key)
	if existing != nil {
		return nil, fmt.Errorf("project with key %s already exists in organization", input.Key)
	}

	project := &models.Project{
		OrganizationID: input.OrganizationID,
		Key:            input.Key,
		Name:           input.Name,
		Description:    input.Description,
		Type:           input.Type,
		Status:         models.ProjectStatusActive,
		LeadID:         input.LeadID,
	}
	if project.Type == "" {
		project.Type = models.ProjectTypeKanban
	}

	// Create project and add creator as admin
	if err := s.repo.CreateWithMember(ctx, project, input.UserID); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	// Publish event
	s.publishEvent("project.created", project.OrganizationID, input.UserID, map[string]interface{}{
		"project_id": project.ID,
		"key":        project.Key,
		"name":       project.Name,
	})

	return project, nil
}

// GetProject gets a project by ID
func (s *ProjectService) GetProject(ctx context.Context, id string) (*models.Project, error) {
	return s.repo.GetByID(ctx, id)
}

// GetProjectByKey gets a project by key
func (s *ProjectService) GetProjectByKey(ctx context.Context, orgID, key string) (*models.Project, error) {
	return s.repo.GetByKey(ctx, orgID, key)
}

// UpdateProjectInput represents input for updating a project
type UpdateProjectInput struct {
	ID          string
	Name        *string
	Description *string
	AvatarURL   *string
	LeadID      *string
	Settings    map[string]string
}

// UpdateProject updates a project
func (s *ProjectService) UpdateProject(ctx context.Context, input UpdateProjectInput) (*models.Project, error) {
	project, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, fmt.Errorf("project not found")
	}

	if input.Name != nil {
		project.Name = *input.Name
	}
	if input.Description != nil {
		project.Description = *input.Description
	}
	if input.AvatarURL != nil {
		project.AvatarURL = *input.AvatarURL
	}
	if input.LeadID != nil {
		project.LeadID = *input.LeadID
	}
	if input.Settings != nil {
		if project.Settings == nil {
			project.Settings = make(map[string]string)
		}
		for k, v := range input.Settings {
			project.Settings[k] = v
		}
	}

	if err := s.repo.Update(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	// Publish event
	s.publishEvent("project.updated", project.OrganizationID, "", map[string]interface{}{
		"project_id": project.ID,
	})

	return project, nil
}

// DeleteProject deletes a project
func (s *ProjectService) DeleteProject(ctx context.Context, id string) error {
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if project == nil {
		return fmt.Errorf("project not found")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	
	// Publish event
	s.publishEvent("project.deleted", project.OrganizationID, "", map[string]interface{}{
		"project_id": id,
	})
	
	return nil
}

// ListProjects lists projects for an organization
func (s *ProjectService) ListProjects(ctx context.Context, orgID string, page, pageSize int) ([]*models.Project, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	
	return s.repo.List(ctx, orgID, pageSize, offset)
}

// ArchiveProject archives a project
func (s *ProjectService) ArchiveProject(ctx context.Context, id string) (*models.Project, error) {
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, fmt.Errorf("project not found")
	}

	project.Status = models.ProjectStatusArchived
	if err := s.repo.Update(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to archive project: %w", err)
	}

	return project, nil
}

// AddMember adds a member to a project
func (s *ProjectService) AddMember(ctx context.Context, projectID, userID string, role models.ProjectRole) (*models.ProjectMember, error) {
	// Check if already member
	existing, _ := s.repo.GetMember(ctx, projectID, userID)
	if existing != nil {
		return nil, fmt.Errorf("user is already a member")
	}

	member := &models.ProjectMember{
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
		JoinedAt:  time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.AddMember(ctx, member); err != nil {
		return nil, fmt.Errorf("failed to add project member: %w", err)
	}

	// Publish event
	project, _ := s.repo.GetByID(ctx, projectID)
	if project != nil {
		s.publishEvent("project.member.added", project.OrganizationID, userID, map[string]interface{}{
			"project_id": projectID,
			"role":       role,
		})
	}

	return member, nil
}

// RemoveMember removes a member from a project
func (s *ProjectService) RemoveMember(ctx context.Context, projectID, userID string) error {
	if err := s.repo.RemoveMember(ctx, projectID, userID); err != nil {
		return fmt.Errorf("failed to remove project member: %w", err)
	}

	// Publish event
	project, _ := s.repo.GetByID(ctx, projectID)
	if project != nil {
		s.publishEvent("project.member.removed", project.OrganizationID, userID, map[string]interface{}{
			"project_id": projectID,
		})
	}

	return nil
}

// UpdateMemberRole updates a member's role
func (s *ProjectService) UpdateMemberRole(ctx context.Context, projectID, userID string, role models.ProjectRole) (*models.ProjectMember, error) {
	if err := s.repo.UpdateMemberRole(ctx, projectID, userID, role); err != nil {
		return nil, fmt.Errorf("failed to update project member role: %w", err)
	}

	member, err := s.repo.GetMember(ctx, projectID, userID)
	if err != nil {
		return nil, err
	}

	return member, nil
}

// ListMembers lists members of a project
func (s *ProjectService) ListMembers(ctx context.Context, projectID string, page, pageSize int) ([]*models.ProjectMember, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	
	return s.repo.ListMembers(ctx, projectID, pageSize, offset)
}

// publishEvent publishes a Kafka event
func (s *ProjectService) publishEvent(eventType, orgID, userID string, payload map[string]interface{}) {
	if s.producer == nil {
		return
	}

	event := kafka.Event{
		Type:           eventType,
		OrganizationID: orgID,
		UserID:         userID,
		Timestamp:      time.Now(),
		Payload:        payload,
	}

	if err := s.producer.PublishEvent("project-events", event); err != nil {
		s.log.Sugar().Errorw("Failed to publish event", "error", err, "type", eventType)
	}
}
