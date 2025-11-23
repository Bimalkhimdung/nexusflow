package service

import (
	"context"
	"fmt"
	"regexp"

	"github.com/nexusflow/nexusflow/pkg/kafka"
	"github.com/nexusflow/nexusflow/pkg/logger"
	"github.com/nexusflow/nexusflow/services/user-service/internal/models"
	"github.com/nexusflow/nexusflow/services/user-service/internal/repository"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// UserService handles user business logic
type UserService struct {
	repo     *repository.UserRepository
	producer *kafka.Producer
	log      *logger.Logger
}

// NewUserService creates a new user service
func NewUserService(repo *repository.UserRepository, producer *kafka.Producer, log *logger.Logger) *UserService {
	return &UserService{
		repo:     repo,
		producer: producer,
		log:      log,
	}
}

// CreateUserInput represents input for creating a user
type CreateUserInput struct {
	OrganizationID string
	Email          string
	DisplayName    string
	AvatarURL      string
	Timezone       string
	Locale         string
	CreatedBy      string
}

// UpdateUserInput represents input for updating a user
type UpdateUserInput struct {
	DisplayName *string
	AvatarURL   *string
	Timezone    *string
	Locale      *string
	Status      *models.UserStatus
	UpdatedBy   string
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, input CreateUserInput) (*models.User, error) {
	// Validate input
	if err := s.validateCreateInput(input); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if user already exists
	existing, _ := s.repo.GetByEmail(ctx, input.Email)
	if existing != nil {
		return nil, fmt.Errorf("user with email %s already exists", input.Email)
	}

	// Create user model
	user := &models.User{
		Email:       input.Email,
		DisplayName: input.DisplayName,
		AvatarURL:   input.AvatarURL,
		Timezone:    input.Timezone,
		Locale:      input.Locale,
		Status:      models.UserStatusActive,
	}
	user.OrganizationID = input.OrganizationID
	user.CreatedBy = input.CreatedBy

	// Save to database
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Publish event
	s.publishUserCreatedEvent(ctx, user)

	s.log.Sugar().Infow("User created successfully", "user_id", user.ID, "email", user.Email)
	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, id string) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(ctx context.Context, id string, input UpdateUserInput) (*models.User, error) {
	// Get existing user
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update fields
	if input.DisplayName != nil {
		user.DisplayName = *input.DisplayName
	}
	if input.AvatarURL != nil {
		user.AvatarURL = *input.AvatarURL
	}
	if input.Timezone != nil {
		user.Timezone = *input.Timezone
	}
	if input.Locale != nil {
		user.Locale = *input.Locale
	}
	if input.Status != nil {
		user.Status = *input.Status
	}
	user.UpdatedBy = input.UpdatedBy

	// Save to database
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Publish event
	s.publishUserUpdatedEvent(ctx, user)

	s.log.Sugar().Infow("User updated successfully", "user_id", user.ID)
	return user, nil
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	// Check if user exists
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Delete user
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Publish event
	s.publishUserDeletedEvent(ctx, user)

	s.log.Sugar().Infow("User deleted successfully", "user_id", id)
	return nil
}

// ListUsers retrieves users with pagination
func (s *UserService) ListUsers(ctx context.Context, orgID string, page, pageSize int) ([]*models.User, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	users, total, err := s.repo.List(ctx, orgID, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}

// SearchUsers searches for users
func (s *UserService) SearchUsers(ctx context.Context, query string, orgIDs []string, page, pageSize int) ([]*models.User, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	users, total, err := s.repo.Search(ctx, query, orgIDs, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search users: %w", err)
	}

	return users, total, nil
}

// UpdateUserPreferences updates user preferences
func (s *UserService) UpdateUserPreferences(ctx context.Context, id string, preferences map[string]interface{}) error {
	if err := s.repo.UpdatePreferences(ctx, id, preferences); err != nil {
		return fmt.Errorf("failed to update preferences: %w", err)
	}

	s.log.Sugar().Infow("User preferences updated", "user_id", id)
	return nil
}

// UpdateLastLogin updates the last login timestamp
func (s *UserService) UpdateLastLogin(ctx context.Context, id string) error {
	if err := s.repo.UpdateLastLogin(ctx, id); err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	return nil
}

// validateCreateInput validates user creation input
func (s *UserService) validateCreateInput(input CreateUserInput) error {
	if input.OrganizationID == "" {
		return fmt.Errorf("organization_id is required")
	}
	if input.Email == "" {
		return fmt.Errorf("email is required")
	}
	if !emailRegex.MatchString(input.Email) {
		return fmt.Errorf("invalid email format")
	}
	if input.DisplayName == "" {
		return fmt.Errorf("display_name is required")
	}
	return nil
}

// publishUserCreatedEvent publishes a user created event
func (s *UserService) publishUserCreatedEvent(ctx context.Context, user *models.User) {
	event := kafka.Event{
		Type:           kafka.EventTypeUserCreated,
		OrganizationID: user.OrganizationID,
		UserID:         user.CreatedBy,
		Payload: map[string]interface{}{
			"user_id":      user.ID,
			"email":        user.Email,
			"display_name": user.DisplayName,
		},
	}

	if err := s.producer.PublishEvent(kafka.TopicUserEvents, event); err != nil {
		s.log.Sugar().Errorw("Failed to publish user created event", "error", err, "user_id", user.ID)
	}
}

// publishUserUpdatedEvent publishes a user updated event
func (s *UserService) publishUserUpdatedEvent(ctx context.Context, user *models.User) {
	event := kafka.Event{
		Type:           kafka.EventTypeUserUpdated,
		OrganizationID: user.OrganizationID,
		UserID:         user.UpdatedBy,
		Payload: map[string]interface{}{
			"user_id":      user.ID,
			"email":        user.Email,
			"display_name": user.DisplayName,
		},
	}

	if err := s.producer.PublishEvent(kafka.TopicUserEvents, event); err != nil {
		s.log.Sugar().Errorw("Failed to publish user updated event", "error", err, "user_id", user.ID)
	}
}

// publishUserDeletedEvent publishes a user deleted event
func (s *UserService) publishUserDeletedEvent(ctx context.Context, user *models.User) {
	event := kafka.Event{
		Type:           "user.deleted",
		OrganizationID: user.OrganizationID,
		UserID:         user.UpdatedBy,
		Payload: map[string]interface{}{
			"user_id": user.ID,
			"email":   user.Email,
		},
	}

	if err := s.producer.PublishEvent(kafka.TopicUserEvents, event); err != nil {
		s.log.Sugar().Errorw("Failed to publish user deleted event", "error", err, "user_id", user.ID)
	}
}
