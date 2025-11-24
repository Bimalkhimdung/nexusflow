package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nexusflow/nexusflow/pkg/kafka"
	"github.com/nexusflow/nexusflow/pkg/logger"
	"github.com/nexusflow/nexusflow/services/org-service/internal/models"
	"github.com/nexusflow/nexusflow/services/org-service/internal/repository"
)

// OrgService handles organization business logic
type OrgService struct {
	orgRepo    *repository.OrgRepository
	teamRepo   *repository.TeamRepository
	inviteRepo *repository.InviteRepository
	producer   *kafka.Producer
	log        *logger.Logger
}

// NewOrgService creates a new organization service
func NewOrgService(
	orgRepo *repository.OrgRepository,
	teamRepo *repository.TeamRepository,
	inviteRepo *repository.InviteRepository,
	producer *kafka.Producer,
	log *logger.Logger,
) *OrgService {
	return &OrgService{
		orgRepo:    orgRepo,
		teamRepo:   teamRepo,
		inviteRepo: inviteRepo,
		producer:   producer,
		log:        log,
	}
}

// CreateOrgInput represents input for creating an organization
type CreateOrgInput struct {
	Name        string
	Slug        string
	Description string
	UserID      string // Creator
}

// CreateOrganization creates a new organization
func (s *OrgService) CreateOrganization(ctx context.Context, input CreateOrgInput) (*models.Organization, error) {
	// Validate input
	if input.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if input.Slug == "" {
		return nil, fmt.Errorf("slug is required")
	}
	if input.UserID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	// Check if slug exists
	existing, _ := s.orgRepo.GetBySlug(ctx, input.Slug)
	if existing != nil {
		return nil, fmt.Errorf("organization with slug %s already exists", input.Slug)
	}

	org := &models.Organization{
		Name:        input.Name,
		Slug:        input.Slug,
		Description: input.Description,
		Status:      models.OrgStatusActive,
		Plan:        models.OrgPlanFree,
	}

	// Create org and add creator as owner
	if err := s.orgRepo.CreateWithMember(ctx, org, input.UserID); err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	// Publish event
	s.publishEvent(kafka.EventTypeOrgCreated, org.ID, input.UserID, map[string]interface{}{
		"name": org.Name,
		"slug": org.Slug,
	})

	return org, nil
}

// GetOrganization gets an organization by ID
func (s *OrgService) GetOrganization(ctx context.Context, id string) (*models.Organization, error) {
	return s.orgRepo.GetByID(ctx, id)
}

// UpdateOrgInput represents input for updating an organization
type UpdateOrgInput struct {
	ID          string
	Name        *string
	Description *string
	LogoURL     *string
	Settings    map[string]string
}

// UpdateOrganization updates an organization
func (s *OrgService) UpdateOrganization(ctx context.Context, input UpdateOrgInput) (*models.Organization, error) {
	org, err := s.orgRepo.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, fmt.Errorf("organization not found")
	}

	if input.Name != nil {
		org.Name = *input.Name
	}
	if input.Description != nil {
		org.Description = *input.Description
	}
	if input.LogoURL != nil {
		org.LogoURL = *input.LogoURL
	}
	if input.Settings != nil {
		if org.Settings == nil {
			org.Settings = make(map[string]string)
		}
		for k, v := range input.Settings {
			org.Settings[k] = v
		}
	}

	if err := s.orgRepo.Update(ctx, org); err != nil {
		return nil, fmt.Errorf("failed to update organization: %w", err)
	}

	// Publish event
	s.publishEvent(kafka.EventTypeOrgUpdated, org.ID, "", map[string]interface{}{
		"name": org.Name,
	})

	return org, nil
}

// DeleteOrganization deletes an organization
func (s *OrgService) DeleteOrganization(ctx context.Context, id string) error {
	if err := s.orgRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}
	
	// Publish event (using empty user ID as we don't have context here easily without fetching)
	s.publishEvent("org.deleted", id, "", nil)
	
	return nil
}

// ListOrganizations lists organizations for a user
func (s *OrgService) ListOrganizations(ctx context.Context, userID string, page, pageSize int) ([]*models.Organization, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	
	return s.orgRepo.List(ctx, userID, pageSize, offset)
}

// AddMember adds a member to an organization
func (s *OrgService) AddMember(ctx context.Context, orgID, userID string, role models.OrgRole) (*models.OrgMember, error) {
	// Check if already member
	existing, _ := s.orgRepo.GetMember(ctx, orgID, userID)
	if existing != nil {
		return nil, fmt.Errorf("user is already a member")
	}

	member := &models.OrgMember{
		OrganizationID: orgID,
		UserID:         userID,
		Role:           role,
		JoinedAt:       time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.orgRepo.AddMember(ctx, member); err != nil {
		return nil, fmt.Errorf("failed to add member: %w", err)
	}

	// Publish event
	s.publishEvent("member.added", orgID, userID, map[string]interface{}{
		"role": role,
	})

	return member, nil
}

// RemoveMember removes a member from an organization
func (s *OrgService) RemoveMember(ctx context.Context, orgID, userID string) error {
	if err := s.orgRepo.RemoveMember(ctx, orgID, userID); err != nil {
		return fmt.Errorf("failed to remove member: %w", err)
	}

	// Publish event
	s.publishEvent("member.removed", orgID, userID, nil)

	return nil
}

// UpdateMemberRole updates a member's role
func (s *OrgService) UpdateMemberRole(ctx context.Context, orgID, userID string, role models.OrgRole) (*models.OrgMember, error) {
	if err := s.orgRepo.UpdateMemberRole(ctx, orgID, userID, role); err != nil {
		return nil, fmt.Errorf("failed to update member role: %w", err)
	}

	member, err := s.orgRepo.GetMember(ctx, orgID, userID)
	if err != nil {
		return nil, err
	}

	// Publish event
	s.publishEvent("member.updated", orgID, userID, map[string]interface{}{
		"role": role,
	})

	return member, nil
}

// ListMembers lists members of an organization
func (s *OrgService) ListMembers(ctx context.Context, orgID string, page, pageSize int) ([]*models.OrgMember, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	
	return s.orgRepo.ListMembers(ctx, orgID, pageSize, offset)
}

// GetMemberRole gets a member's role in an organization
func (s *OrgService) GetMemberRole(ctx context.Context, orgID, userID string) (models.OrgRole, bool, error) {
	return s.orgRepo.GetMemberRole(ctx, orgID, userID)
}

// CreateTeam creates a new team
func (s *OrgService) CreateTeam(ctx context.Context, orgID, name, description string) (*models.Team, error) {
	team := &models.Team{
		OrganizationID: orgID,
		Name:           name,
		Description:    description,
	}

	if err := s.teamRepo.Create(ctx, team); err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	return team, nil
}

// GetTeam gets a team by ID
func (s *OrgService) GetTeam(ctx context.Context, id string) (*models.Team, error) {
	return s.teamRepo.GetByID(ctx, id)
}

// UpdateTeam updates a team
func (s *OrgService) UpdateTeam(ctx context.Context, id string, name, description *string) (*models.Team, error) {
	team, err := s.teamRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, fmt.Errorf("team not found")
	}

	if name != nil {
		team.Name = *name
	}
	if description != nil {
		team.Description = *description
	}

	if err := s.teamRepo.Update(ctx, team); err != nil {
		return nil, fmt.Errorf("failed to update team: %w", err)
	}

	return team, nil
}

// DeleteTeam deletes a team
func (s *OrgService) DeleteTeam(ctx context.Context, id string) error {
	return s.teamRepo.Delete(ctx, id)
}

// ListTeams lists teams in an organization
func (s *OrgService) ListTeams(ctx context.Context, orgID string, page, pageSize int) ([]*models.Team, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	
	return s.teamRepo.List(ctx, orgID, pageSize, offset)
}

// AddTeamMember adds a user to a team
func (s *OrgService) AddTeamMember(ctx context.Context, teamID, userID string) error {
	return s.teamRepo.AddMember(ctx, teamID, userID)
}

// RemoveTeamMember removes a user from a team
func (s *OrgService) RemoveTeamMember(ctx context.Context, teamID, userID string) error {
	return s.teamRepo.RemoveMember(ctx, teamID, userID)
}

// CreateInvite creates a new invite
func (s *OrgService) CreateInvite(ctx context.Context, orgID, email string, role models.OrgRole, invitedBy string) (*models.Invite, error) {
	// Generate token
	token := uuid.New().String()

	invite := &models.Invite{
		OrganizationID: orgID,
		Email:          email,
		Role:           role,
		InvitedBy:      invitedBy,
		Token:          token,
		Status:         models.InviteStatusPending,
		ExpiresAt:      time.Now().Add(7 * 24 * time.Hour), // 7 days expiry
	}

	if err := s.inviteRepo.Create(ctx, invite); err != nil {
		return nil, fmt.Errorf("failed to create invite: %w", err)
	}

	// Publish event
	s.publishEvent("invite.created", orgID, invitedBy, map[string]interface{}{
		"email": email,
		"role":  role,
		"token": token,
	})

	return invite, nil
}

// AcceptInvite accepts an invite
func (s *OrgService) AcceptInvite(ctx context.Context, token, userID string) (*models.OrgMember, error) {
	invite, err := s.inviteRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if invite == nil {
		return nil, fmt.Errorf("invite not found")
	}

	if invite.Status != models.InviteStatusPending {
		return nil, fmt.Errorf("invite is not pending")
	}
	if time.Now().After(invite.ExpiresAt) {
		invite.Status = models.InviteStatusExpired
		_ = s.inviteRepo.Update(ctx, invite)
		return nil, fmt.Errorf("invite expired")
	}

	// Add member
	member, err := s.AddMember(ctx, invite.OrganizationID, userID, invite.Role)
	if err != nil {
		return nil, err
	}

	// Update invite status
	invite.Status = models.InviteStatusAccepted
	invite.AcceptedAt = time.Now()
	if err := s.inviteRepo.Update(ctx, invite); err != nil {
		s.log.Sugar().Errorw("Failed to update invite status", "error", err)
	}

	return member, nil
}

// RevokeInvite revokes an invite
func (s *OrgService) RevokeInvite(ctx context.Context, id string) error {
	invite, err := s.inviteRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if invite == nil {
		return fmt.Errorf("invite not found")
	}

	invite.Status = models.InviteStatusRevoked
	return s.inviteRepo.Update(ctx, invite)
}

// ListInvites lists invites for an organization
func (s *OrgService) ListInvites(ctx context.Context, orgID string, page, pageSize int) ([]*models.Invite, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	
	return s.inviteRepo.List(ctx, orgID, pageSize, offset)
}

// publishEvent publishes a Kafka event
func (s *OrgService) publishEvent(eventType, orgID, userID string, payload map[string]interface{}) {
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

	if err := s.producer.PublishEvent(kafka.TopicOrgEvents, event); err != nil {
		s.log.Sugar().Errorw("Failed to publish event", "error", err, "type", eventType)
	}
}
