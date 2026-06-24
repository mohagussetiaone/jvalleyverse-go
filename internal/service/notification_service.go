package service

import (
	"context"
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/repository"
)

type INotificationService interface {
	CreateNotification(ctx context.Context, userID, nType, title, message, link string) error
	ListNotifications(ctx context.Context, userID string, page, limit int) ([]map[string]interface{}, int64, error)
	CountUnread(ctx context.Context, userID string) (int64, error)
	MarkAsRead(ctx context.Context, notificationID, userID string) error
	MarkAllAsRead(ctx context.Context, userID string) error
	DeleteNotification(ctx context.Context, notificationID, userID string) error
}

type NotificationService struct {
	notifRepo *repository.NotificationRepository
}

func NewNotificationService(notifRepo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{notifRepo: notifRepo}
}

func (s *NotificationService) CreateNotification(ctx context.Context, userID, nType, title, message, link string) error {
	notif := &domain.Notification{
		UserID:  userID,
		Type:    nType,
		Title:   title,
		Message: message,
		Link:    link,
		IsRead:  false,
	}

	if err := s.notifRepo.Create(ctx, notif); err != nil {
		return err
	}

	// Push real-time via SSE hub
	hub := GetNotificationHub()
	hub.Publish(userID, SSEEvent{
		Type:    nType,
		Title:   title,
		Message: message,
		Link:    link,
		Payload: map[string]interface{}{
			"id":         notif.ID,
			"created_at": notif.CreatedAt,
			"is_read":    false,
		},
	})

	return nil
}

func (s *NotificationService) ListNotifications(ctx context.Context, userID string, page, limit int) ([]map[string]interface{}, int64, error) {
	notifs, total, err := s.notifRepo.ListByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	result := make([]map[string]interface{}, len(notifs))
	for i, n := range notifs {
		result[i] = map[string]interface{}{
			"id":         n.ID,
			"type":       n.Type,
			"title":      n.Title,
			"message":    n.Message,
			"is_read":    n.IsRead,
			"link":       n.Link,
			"created_at": n.CreatedAt,
		}
	}

	return result, total, nil
}

func (s *NotificationService) CountUnread(ctx context.Context, userID string) (int64, error) {
	return s.notifRepo.CountUnreadByUserID(ctx, userID)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID, userID string) error {
	return s.notifRepo.MarkAsRead(ctx, notificationID, userID)
}

func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID string) error {
	return s.notifRepo.MarkAllAsRead(ctx, userID)
}

func (s *NotificationService) DeleteNotification(ctx context.Context, notificationID, userID string) error {
	return s.notifRepo.DeleteByID(ctx, notificationID, userID)
}
