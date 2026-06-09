package repository

import (
	"context"
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

// ShowcaseRepository handles showcase data access
type ShowcaseRepository struct {
	db *gorm.DB
}

func NewShowcaseRepository(db *gorm.DB) *ShowcaseRepository {
	return &ShowcaseRepository{db: db}
}

// Create creates new showcase
func (r *ShowcaseRepository) Create(ctx context.Context, showcase *domain.Showcase) error {
	return r.db.WithContext(ctx).Create(showcase).Error
}

// FindByID finds showcase by CUID with user info
func (r *ShowcaseRepository) FindByID(ctx context.Context, showcaseID string) (*domain.Showcase, error) {
	showcase := &domain.Showcase{}
	if err := r.db.WithContext(ctx).Where("id = ?", showcaseID).Preload("User").First(showcase).Error; err != nil {
		return nil, err
	}
	return showcase, nil
}

// ListAll lists showcases with pagination, optional category filter, ordered by newest
func (r *ShowcaseRepository) ListAll(ctx context.Context, page, limit int, categoryID *string) ([]domain.Showcase, int64, error) {
	var showcases []domain.Showcase
	var total int64

	// buildBase creates a fresh query each time to avoid GORM state mutation bug
	buildBase := func() *gorm.DB {
		q := r.db.WithContext(ctx).
			Model(&domain.Showcase{}).
			Where("visibility = ?", "public").
			Where("status = ?", "published")
		if categoryID != nil && *categoryID != "" {
			q = q.Where("category_id = ?", *categoryID)
		}
		return q
	}

	if err := buildBase().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := buildBase().
		Preload("User").
		Preload("Category").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&showcases).Error; err != nil {
		return nil, 0, err
	}

	return showcases, total, nil
}

// ListByUserID lists showcases by user
func (r *ShowcaseRepository) ListByUserID(ctx context.Context, userID string, page, limit int) ([]domain.Showcase, int64, error) {
	var showcases []domain.Showcase
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Showcase{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Preload("User").Offset(offset).Limit(limit).Find(&showcases).Error; err != nil {
		return nil, 0, err
	}

	return showcases, total, nil
}

// Update updates showcase
func (r *ShowcaseRepository) Update(ctx context.Context, showcase *domain.Showcase) error {
	return r.db.WithContext(ctx).Model(showcase).Updates(showcase).Error
}

// Delete deletes showcase (cascade to likes and comments)
func (r *ShowcaseRepository) Delete(ctx context.Context, showcaseID string) error {
	return r.db.WithContext(ctx).Where("id = ?", showcaseID).Delete(&domain.Showcase{}).Error
}

// IncrementLikes increments like count
func (r *ShowcaseRepository) IncrementLikes(ctx context.Context, showcaseID string) error {
	return r.db.WithContext(ctx).Model(&domain.Showcase{}).Where("id = ?", showcaseID).
		Update("likes_count", gorm.Expr("likes_count + 1")).Error
}

// DecrementLikes decrements like count
func (r *ShowcaseRepository) DecrementLikes(ctx context.Context, showcaseID string) error {
	return r.db.WithContext(ctx).Model(&domain.Showcase{}).Where("id = ?", showcaseID).
		Update("likes_count", gorm.Expr("likes_count - 1")).Error
}

// IsLikedByUser checks if showcase is liked by user
func (r *ShowcaseRepository) IsLikedByUser(ctx context.Context, showcaseID, userID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.ShowcaseLike{}).
		Where("showcase_id = ? AND user_id = ?", showcaseID, userID).
		Count(&count).Error
	return count > 0, err
}

// GetLeaderboard returns top showcases by likes
func (r *ShowcaseRepository) GetLeaderboard(ctx context.Context, limit int) ([]domain.Showcase, error) {
	var showcases []domain.Showcase
	err := r.db.WithContext(ctx).Preload("User").Order("likes_count DESC").Limit(limit).Find(&showcases).Error
	return showcases, err
}
