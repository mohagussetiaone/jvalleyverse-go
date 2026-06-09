package repository

import (
	"context"
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

// ShowcaseLikeRepository handles showcase likes
type ShowcaseLikeRepository struct {
	db *gorm.DB
}

func NewShowcaseLikeRepository(db *gorm.DB) *ShowcaseLikeRepository {
	return &ShowcaseLikeRepository{db: db}
}

// Create adds like
func (r *ShowcaseLikeRepository) Create(ctx context.Context, like *domain.ShowcaseLike) error {
	return r.db.WithContext(ctx).Create(like).Error
}

// DeleteByUserShowcase removes like
func (r *ShowcaseLikeRepository) DeleteByUserShowcase(ctx context.Context, userID, showcaseID string) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND showcase_id = ?", userID, showcaseID).Delete(&domain.ShowcaseLike{}).Error
}

// Exists checks if like exists
func (r *ShowcaseLikeRepository) Exists(ctx context.Context, userID, showcaseID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.ShowcaseLike{}).
		Where("user_id = ? AND showcase_id = ?", userID, showcaseID).
		Count(&count).Error
	return count > 0, err
}
