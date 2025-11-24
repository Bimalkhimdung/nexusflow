package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/nexusflow/nexusflow/pkg/database"
	"github.com/nexusflow/nexusflow/pkg/logger"
	"github.com/nexusflow/nexusflow/services/notification-service/internal/models"
)

type NotificationRepository struct {
	db  *database.DB
	log *logger.Logger
}

func NewNotificationRepository(db *database.DB, log *logger.Logger) *NotificationRepository {
	return &NotificationRepository{db: db, log: log}
}

// Notification CRUD
func (r *NotificationRepository) CreateNotification(ctx context.Context, notification *models.Notification) error {
	notification.ID = ""
	notification.CreatedAt = time.Now()
	_, err := r.db.NewInsert().Model(notification).Exec(ctx)
	if err != nil {
		return fmt.Errorf("create notification: %w", err)
	}
	return nil
}

func (r *NotificationRepository) GetNotification(ctx context.Context, id string) (*models.Notification, error) {
	n := new(models.Notification)
	err := r.db.NewSelect().Model(n).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("get notification: %w", err)
	}
	return n, nil
}

func (r *NotificationRepository) ListNotifications(ctx context.Context, userID string, limit, offset int) ([]*models.Notification, int, error) {
	var notifications []*models.Notification
	
	query := r.db.NewSelect().Model(&notifications).
		Where("user_id = ?", userID).
		Order("created_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	err := query.Scan(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("list notifications: %w", err)
	}
	
	// Get total count
	total, err := r.db.NewSelect().Model((*models.Notification)(nil)).
		Where("user_id = ?", userID).
		Count(ctx)
	if err != nil {
		return notifications, 0, fmt.Errorf("count notifications: %w", err)
	}
	
	return notifications, total, nil
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, id string) error {
	_, err := r.db.NewUpdate().Model((*models.Notification)(nil)).
		Set("read = ?", true).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("mark as read: %w", err)
	}
	return nil
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	_, err := r.db.NewUpdate().Model((*models.Notification)(nil)).
		Set("read = ?", true).
		Where("user_id = ? AND read = ?", userID, false).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("mark all as read: %w", err)
	}
	return nil
}

func (r *NotificationRepository) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	count, err := r.db.NewSelect().Model((*models.Notification)(nil)).
		Where("user_id = ? AND read = ?", userID, false).
		Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("get unread count: %w", err)
	}
	return count, nil
}

// Preferences
func (r *NotificationRepository) GetPreferences(ctx context.Context, userID string) ([]*models.NotificationPreference, error) {
	var prefs []*models.NotificationPreference
	err := r.db.NewSelect().Model(&prefs).
		Where("user_id = ?", userID).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("get preferences: %w", err)
	}
	return prefs, nil
}

func (r *NotificationRepository) UpdatePreference(ctx context.Context, pref *models.NotificationPreference) error {
	pref.UpdatedAt = time.Now()
	_, err := r.db.NewInsert().Model(pref).
		On("CONFLICT (user_id, notification_type) DO UPDATE").
		Set("in_app_enabled = EXCLUDED.in_app_enabled").
		Set("email_enabled = EXCLUDED.email_enabled").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("update preference: %w", err)
	}
	return nil
}
