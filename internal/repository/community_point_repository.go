package repository

import (
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

// CommunityPointRepository handles community points data access
type CommunityPointRepository struct {
	db *gorm.DB
}

func NewCommunityPointRepository() *CommunityPointRepository {
	return &CommunityPointRepository{db: db}
}

// Create records points transaction
func (r *CommunityPointRepository) Create(point *domain.CommunityPoint) error {
	return r.db.Create(point).Error
}

// ListByUserID lists points history for user
func (r *CommunityPointRepository) ListByUserID(userID uint, page, limit int) ([]domain.CommunityPoint, int64, error) {
	var points []domain.CommunityPoint
	var total int64

	if err := r.db.Model(&domain.CommunityPoint{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.Where("user_id = ?", userID).Offset(offset).Limit(limit).Order("created_at DESC").Find(&points).Error; err != nil {
		return nil, 0, err
	}

	return points, total, nil
}

// GetTotalPointsByUser gets accumulated points for user
func (r *CommunityPointRepository) GetTotalPointsByUser(userID uint) (int, error) {
	var total int
	err := r.db.Model(&domain.CommunityPoint{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(points_earned), 0)").
		Row().
		Scan(&total)
	return total, err
}
