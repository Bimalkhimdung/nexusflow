package handler

import (
	"context"
	"time"

	"github.com/nexusflow/nexusflow/pkg/logger"
	pb "github.com/nexusflow/nexusflow/pkg/proto/notification/v1"
	"github.com/nexusflow/nexusflow/services/notification-service/internal/models"
	"github.com/nexusflow/nexusflow/services/notification-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotificationHandler struct {
	pb.UnimplementedNotificationServiceServer
	svc *service.NotificationService
	log *logger.Logger
}

func NewNotificationHandler(svc *service.NotificationService, log *logger.Logger) *NotificationHandler {
	return &NotificationHandler{svc: svc, log: log}
}

// Helper conversions
func notificationToProto(n *models.Notification) *pb.Notification {
	if n == nil {
		return nil
	}
	return &pb.Notification{
		Id:        n.ID,
		UserId:    n.UserID,
		Type:      n.Type,
		Title:     n.Title,
		Message:   n.Message,
		Link:      n.Link,
		Metadata:  string(n.Metadata),
		Read:      n.Read,
		CreatedAt: n.CreatedAt.Format(time.RFC3339),
	}
}

func preferenceToProto(p *models.NotificationPreference) *pb.NotificationPreference {
	if p == nil {
		return nil
	}
	return &pb.NotificationPreference{
		Id:               p.ID,
		UserId:           p.UserID,
		NotificationType: p.NotificationType,
		InAppEnabled:     p.InAppEnabled,
		EmailEnabled:     p.EmailEnabled,
	}
}

// RPC Methods
func (h *NotificationHandler) ListNotifications(ctx context.Context, req *pb.ListNotificationsRequest) (*pb.ListNotificationsResponse, error) {
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 20
	}
	offset := int(req.Offset)

	notifications, total, err := h.svc.ListNotifications(ctx, req.UserId, limit, offset)
	if err != nil {
		h.log.Sugar().Errorw("Failed to list notifications", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to list notifications: %v", err)
	}

	var pbNotifications []*pb.Notification
	for _, n := range notifications {
		pbNotifications = append(pbNotifications, notificationToProto(n))
	}

	return &pb.ListNotificationsResponse{
		Notifications: pbNotifications,
		Total:         int32(total),
	}, nil
}

func (h *NotificationHandler) MarkAsRead(ctx context.Context, req *pb.MarkAsReadRequest) (*pb.MarkAsReadResponse, error) {
	if err := h.svc.MarkAsRead(ctx, req.NotificationId); err != nil {
		h.log.Sugar().Errorw("Failed to mark as read", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to mark as read: %v", err)
	}
	return &pb.MarkAsReadResponse{}, nil
}

func (h *NotificationHandler) MarkAllAsRead(ctx context.Context, req *pb.MarkAllAsReadRequest) (*pb.MarkAllAsReadResponse, error) {
	if err := h.svc.MarkAllAsRead(ctx, req.UserId); err != nil {
		h.log.Sugar().Errorw("Failed to mark all as read", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to mark all as read: %v", err)
	}
	return &pb.MarkAllAsReadResponse{}, nil
}

func (h *NotificationHandler) GetUnreadCount(ctx context.Context, req *pb.GetUnreadCountRequest) (*pb.GetUnreadCountResponse, error) {
	count, err := h.svc.GetUnreadCount(ctx, req.UserId)
	if err != nil {
		h.log.Sugar().Errorw("Failed to get unread count", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get unread count: %v", err)
	}
	return &pb.GetUnreadCountResponse{Count: int32(count)}, nil
}

func (h *NotificationHandler) GetPreferences(ctx context.Context, req *pb.GetPreferencesRequest) (*pb.GetPreferencesResponse, error) {
	prefs, err := h.svc.GetPreferences(ctx, req.UserId)
	if err != nil {
		h.log.Sugar().Errorw("Failed to get preferences", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get preferences: %v", err)
	}

	var pbPrefs []*pb.NotificationPreference
	for _, p := range prefs {
		pbPrefs = append(pbPrefs, preferenceToProto(p))
	}

	return &pb.GetPreferencesResponse{Preferences: pbPrefs}, nil
}

func (h *NotificationHandler) UpdatePreference(ctx context.Context, req *pb.UpdatePreferenceRequest) (*pb.UpdatePreferenceResponse, error) {
	pref, err := h.svc.UpdatePreference(ctx, req)
	if err != nil {
		h.log.Sugar().Errorw("Failed to update preference", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to update preference: %v", err)
	}
	return &pb.UpdatePreferenceResponse{Preference: preferenceToProto(pref)}, nil
}
