package repository

import (
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

// ShowcaseLikeRepository handles showcase likes
type ShowcaseLikeRepository struct {
	db *gorm.DB
}

func NewShowcaseLikeRepository() *ShowcaseLikeRepository {
	return &ShowcaseLikeRepository{db: db}
}

// Create adds like
func (r *ShowcaseLikeRepository) Create(like *domain.ShowcaseLike) error {
	return r.db.Create(like).Error
}

// DeleteByUserShowcase removes like
func (r *ShowcaseLikeRepository) DeleteByUserShowcase(userID, showcaseID uint) error {
	return r.db.Where("user_id = ? AND showcase_id = ?", userID, showcaseID).Delete(&domain.ShowcaseLike{}).Error
}

// Exists checks if like exists
func (r *ShowcaseLikeRepository) Exists(userID, showcaseID uint) (bool, error) {
	var count int64
	err := r.db.Model(&domain.ShowcaseLike{}).
		Where("user_id = ? AND showcase_id = ?", userID, showcaseID).
		Count(&count).Error
	return count > 0, err
}
