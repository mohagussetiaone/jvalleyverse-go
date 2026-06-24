package repository

import (
	"context"
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

type AdminAuditLogRepository struct {
	db *gorm.DB
}

func NewAdminAuditLogRepository(db *gorm.DB) *AdminAuditLogRepository {
	return &AdminAuditLogRepository{db: db}
}

func (r *AdminAuditLogRepository) Create(ctx context.Context, log *domain.AdminAuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *AdminAuditLogRepository) List(ctx context.Context, page, limit int) ([]domain.AdminAuditLog, int64, error) {
	var logs []domain.AdminAuditLog
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.AdminAuditLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).Preload("Admin").Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
