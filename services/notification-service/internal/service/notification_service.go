package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nexusflow/nexusflow/pkg/logger"
	pb "github.com/nexusflow/nexusflow/pkg/proto/notification/v1"
	"github.com/nexusflow/nexusflow/services/notification-service/internal/models"
	"github.com/nexusflow/nexusflow/services/notification-service/internal/repository"
	"github.com/nexusflow/nexusflow/services/notification-service/internal/websocket"
)

type NotificationService struct {
	repo *repository.NotificationRepository
	hub  *websocket.Hub
	log  *logger.Logger
}

func NewNotificationService(repo *repository.NotificationRepository, hub *websocket.Hub, log *logger.Logger) *NotificationService {
	return &NotificationService{repo: repo, hub: hub, log: log}
}

// CreateNotification creates a notification and broadcasts it via WebSocket
func (s *NotificationService) CreateNotification(ctx context.Context, notification *models.Notification) error {
	if err := s.repo.CreateNotification(ctx, notification); err != nil {
		return fmt.Errorf("create notification: %w", err)
	}

	// Broadcast to WebSocket clients
	if s.hub != nil {
		s.hub.Broadcast(notification.UserID, notification)
	}

	return nil
}

func (s *NotificationService) ListNotifications(ctx context.Context, userID string, limit, offset int) ([]*models.Notification, int, error) {
	return s.repo.ListNotifications(ctx, userID, limit, offset)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, id string) error {
	return s.repo.MarkAsRead(ctx, id)
}

func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID string) error {
	return s.repo.MarkAllAsRead(ctx, userID)
}

func (s *NotificationService) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	return s.repo.GetUnreadCount(ctx, userID)
}

func (s *NotificationService) GetPreferences(ctx context.Context, userID string) ([]*models.NotificationPreference, error) {
	return s.repo.GetPreferences(ctx, userID)
}

func (s *NotificationService) UpdatePreference(ctx context.Context, req *pb.UpdatePreferenceRequest) (*models.NotificationPreference, error) {
	pref := &models.NotificationPreference{
		UserID:           req.UserId,
		NotificationType: req.NotificationType,
		InAppEnabled:     req.InAppEnabled,
		EmailEnabled:     req.EmailEnabled,
	}

	if err := s.repo.UpdatePreference(ctx, pref); err != nil {
		return nil, fmt.Errorf("update preference: %w", err)
	}

	return pref, nil
}

// ProcessEvent processes Kafka events and creates notifications
func (s *NotificationService) ProcessEvent(ctx context.Context, eventType string, payload map[string]interface{}) error {
	var notification *models.Notification

	switch eventType {
	case "comment.mention_created":
		notification = s.createMentionNotification(payload)
	case "issue.assigned":
		notification = s.createIssueAssignedNotification(payload)
	case "sprint.started":
		notification = s.createSprintStartedNotification(payload)
	case "sprint.completed":
		notification = s.createSprintCompletedNotification(payload)
	default:
		// Ignore unknown event types
		return nil
	}

	if notification != nil {
		return s.CreateNotification(ctx, notification)
	}

	return nil
}

func (s *NotificationService) createMentionNotification(payload map[string]interface{}) *models.Notification {
	userID, _ := payload["mentioned_user_id"].(string)
	commentID, _ := payload["comment_id"].(string)
	issueID, _ := payload["issue_id"].(string)

	metadata, _ := json.Marshal(payload)

	return &models.Notification{
		UserID:   userID,
		Type:     models.NotificationTypeCommentMention,
		Title:    "You were mentioned",
		Message:  "Someone mentioned you in a comment",
		Link:     fmt.Sprintf("/issues/%s#comment-%s", issueID, commentID),
		Metadata: metadata,
	}
}

func (s *NotificationService) createIssueAssignedNotification(payload map[string]interface{}) *models.Notification {
	userID, _ := payload["assignee_id"].(string)
	issueID, _ := payload["issue_id"].(string)
	issueKey, _ := payload["issue_key"].(string)

	metadata, _ := json.Marshal(payload)

	return &models.Notification{
		UserID:   userID,
		Type:     models.NotificationTypeIssueAssigned,
		Title:    "Issue assigned to you",
		Message:  fmt.Sprintf("You have been assigned to %s", issueKey),
		Link:     fmt.Sprintf("/issues/%s", issueID),
		Metadata: metadata,
	}
}

func (s *NotificationService) createSprintStartedNotification(payload map[string]interface{}) *models.Notification {
	// This would need to notify all team members
	// For now, we'll skip this as we don't have team membership info
	return nil
}

func (s *NotificationService) createSprintCompletedNotification(payload map[string]interface{}) *models.Notification {
	// Similar to sprint started
	return nil
}
