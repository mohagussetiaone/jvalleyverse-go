package repository

import (
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

// DiscussionRepository handles discussion data access
type DiscussionRepository struct {
	db *gorm.DB
}

func NewDiscussionRepository() *DiscussionRepository {
	return &DiscussionRepository{db: db}
}

// Create creates new discussion
func (r *DiscussionRepository) Create(discussion *domain.Discussion) error {
	return r.db.Create(discussion).Error
}

// FindByID finds discussion with user info
func (r *DiscussionRepository) FindByID(discussionID uint) (*domain.Discussion, error) {
	discussion := &domain.Discussion{}
	if err := r.db.Preload("User").Preload("Class").First(discussion, discussionID).Error; err != nil {
		return nil, err
	}
	return discussion, nil
}

// ListByClassID lists discussions in a class
func (r *DiscussionRepository) ListByClassID(classID uint, page, limit int) ([]domain.Discussion, int64, error) {
	var discussions []domain.Discussion
	var total int64

	if err := r.db.Model(&domain.Discussion{}).Where("class_id = ?", classID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.Where("class_id = ?", classID).Preload("User").Offset(offset).Limit(limit).Order("created_at DESC").Find(&discussions).Error; err != nil {
		return nil, 0, err
	}

	return discussions, total, nil
}

// ListByUserID lists discussions created by user
func (r *DiscussionRepository) ListByUserID(userID uint, page, limit int) ([]domain.Discussion, int64, error) {
	var discussions []domain.Discussion
	var total int64

	if err := r.db.Model(&domain.Discussion{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.Where("user_id = ?", userID).Preload("Class").Offset(offset).Limit(limit).Order("created_at DESC").Find(&discussions).Error; err != nil {
		return nil, 0, err
	}

	return discussions, total, nil
}

// Update updates discussion
func (r *DiscussionRepository) Update(discussion *domain.Discussion) error {
	return r.db.Model(discussion).Updates(discussion).Error
}

// DeleteByID deletes discussion and cascade replies
func (r *DiscussionRepository) DeleteByID(discussionID uint) error {
	return r.db.Delete(&domain.Discussion{}, discussionID).Error
}
