package repository

import (
	"context"
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

func (r *NotificationRepository) ListByUserID(ctx context.Context, userID string, page, limit int) ([]domain.Notification, int64, error) {
	var notifications []domain.Notification
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Notification{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Offset(offset).Limit(limit).Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func (r *NotificationRepository) CountUnreadByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, notificationID, userID string) error {
	return r.db.WithContext(ctx).Model(&domain.Notification{}).Where("id = ? AND user_id = ?", notificationID, userID).Update("is_read", true).Error
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Model(&domain.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true).Error
}

func (r *NotificationRepository) DeleteByID(ctx context.Context, notificationID, userID string) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", notificationID, userID).Delete(&domain.Notification{}).Error
}
