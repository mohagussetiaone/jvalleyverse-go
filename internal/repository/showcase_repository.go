package repository

import (
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

// ShowcaseRepository handles showcase data access
type ShowcaseRepository struct {
	db *gorm.DB
}

func NewShowcaseRepository() *ShowcaseRepository {
	return &ShowcaseRepository{db: db}
}

// Create creates new showcase
func (r *ShowcaseRepository) Create(showcase *domain.Showcase) error {
	return r.db.Create(showcase).Error
}

// FindByID finds showcase by ID with user and likes count
func (r *ShowcaseRepository) FindByID(showcaseID uint) (*domain.Showcase, error) {
	showcase := &domain.Showcase{}
	if err := r.db.Preload("User").First(showcase, showcaseID).Error; err != nil {
		return nil, err
	}
	return showcase, nil
}

// ListAll lists showcases with pagination and user info
func (r *ShowcaseRepository) ListAll(page, limit int, categoryID *uint) ([]domain.Showcase, int64, error) {
	var showcases []domain.Showcase
	var total int64

	query := r.db
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}

	if err := query.Model(&domain.Showcase{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Preload("User").Offset(offset).Limit(limit).Find(&showcases).Error; err != nil {
		return nil, 0, err
	}

	return showcases, total, nil
}

// ListByUserID lists showcases by user
func (r *ShowcaseRepository) ListByUserID(userID uint, page, limit int) ([]domain.Showcase, int64, error) {
	var showcases []domain.Showcase
	var total int64

	if err := r.db.Model(&domain.Showcase{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.Where("user_id = ?", userID).Preload("User").Offset(offset).Limit(limit).Find(&showcases).Error; err != nil {
		return nil, 0, err
	}

	return showcases, total, nil
}

// Update updates showcase
func (r *ShowcaseRepository) Update(showcase *domain.Showcase) error {
	return r.db.Model(showcase).Updates(showcase).Error
}

// Delete deletes showcase (cascade to likes and comments)
func (r *ShowcaseRepository) Delete(showcaseID uint) error {
	return r.db.Delete(&domain.Showcase{}, showcaseID).Error
}

// IncrementLikes increments like count
func (r *ShowcaseRepository) IncrementLikes(showcaseID uint) error {
	return r.db.Model(&domain.Showcase{}).Where("id = ?", showcaseID).
		Update("likes_count", gorm.Expr("likes_count + 1")).Error
}

// DecrementLikes decrements like count
func (r *ShowcaseRepository) DecrementLikes(showcaseID uint) error {
	return r.db.Model(&domain.Showcase{}).Where("id = ?", showcaseID).
		Update("likes_count", gorm.Expr("likes_count - 1")).Error
}

// IsLikedByUser checks if showcase is liked by user
func (r *ShowcaseRepository) IsLikedByUser(showcaseID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&domain.ShowcaseLike{}).
		Where("showcase_id = ? AND user_id = ?", showcaseID, userID).
		Count(&count).Error
	return count > 0, err
}

// GetLeaderboard returns top showcases by likes
func (r *ShowcaseRepository) GetLeaderboard(limit int) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}
