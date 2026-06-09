package repository

import (
	"context"
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

// CommunityPointRepository handles community points data access
type CommunityPointRepository struct {
	db *gorm.DB
}

func NewCommunityPointRepository(db *gorm.DB) *CommunityPointRepository {
	return &CommunityPointRepository{db: db}
}

// Create records points transaction
func (r *CommunityPointRepository) Create(ctx context.Context, point *domain.CommunityPoint) error {
	return r.db.WithContext(ctx).Create(point).Error
}

// ListByUserID lists points history for user
func (r *CommunityPointRepository) ListByUserID(ctx context.Context, userID string, page, limit int) ([]domain.CommunityPoint, int64, error) {
	var points []domain.CommunityPoint
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.CommunityPoint{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Offset(offset).Limit(limit).Order("created_at DESC").Find(&points).Error; err != nil {
		return nil, 0, err
	}

	return points, total, nil
}

// GetTotalPointsByUser gets accumulated points for user
func (r *CommunityPointRepository) GetTotalPointsByUser(ctx context.Context, userID string) (int, error) {
	var total int
	err := r.db.WithContext(ctx).Model(&domain.CommunityPoint{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(points_earned), 0)").
		Row().
		Scan(&total)
	return total, err
}
