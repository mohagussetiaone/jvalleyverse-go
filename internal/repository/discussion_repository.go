package repository

import (
	"context"
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

// DiscussionRepository handles discussion data access
type DiscussionRepository struct {
	db *gorm.DB
}

func NewDiscussionRepository(db *gorm.DB) *DiscussionRepository {
	return &DiscussionRepository{db: db}
}

// Create creates new discussion
func (r *DiscussionRepository) Create(ctx context.Context, discussion *domain.Discussion) error {
	return r.db.WithContext(ctx).Create(discussion).Error
}

// FindByID finds discussion with user info
func (r *DiscussionRepository) FindByID(ctx context.Context, discussionID string) (*domain.Discussion, error) {
	discussion := &domain.Discussion{}
	if err := r.db.WithContext(ctx).Where("id = ?", discussionID).Preload("User").Preload("Lesson").First(discussion).Error; err != nil {
		return nil, err
	}
	return discussion, nil
}

// ListAll lists all discussions with pagination
func (r *DiscussionRepository) ListAll(ctx context.Context, page, limit int) ([]domain.Discussion, int64, error) {
	var discussions []domain.Discussion
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Discussion{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).Preload("User").Preload("Replies").Offset(offset).Limit(limit).Order("created_at DESC").Find(&discussions).Error; err != nil {
		return nil, 0, err
	}

	return discussions, total, nil
}

func (r *DiscussionRepository) ListByLessonID(ctx context.Context, lessonID string, page, limit int) ([]domain.Discussion, int64, error) {
	var discussions []domain.Discussion
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Discussion{}).Where("lesson_id = ?", lessonID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).Where("lesson_id = ?", lessonID).Preload("User").Preload("Replies").Offset(offset).Limit(limit).Order("created_at DESC").Find(&discussions).Error; err != nil {
		return nil, 0, err
	}

	return discussions, total, nil
}

func (r *DiscussionRepository) ListByStudyCaseID(ctx context.Context, studyCaseID string, page, limit int) ([]domain.Discussion, int64, error) {
	var discussions []domain.Discussion
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Discussion{}).Where("study_case_id = ?", studyCaseID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).Where("study_case_id = ?", studyCaseID).Preload("User").Preload("Replies").Offset(offset).Limit(limit).Order("created_at DESC").Find(&discussions).Error; err != nil {
		return nil, 0, err
	}

	return discussions, total, nil
}

// ListByUserID lists discussions created by user
func (r *DiscussionRepository) ListByUserID(ctx context.Context, userID string, page, limit int) ([]domain.Discussion, int64, error) {
	var discussions []domain.Discussion
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Discussion{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Preload("User").Preload("Lesson").Preload("Replies").Offset(offset).Limit(limit).Order("created_at DESC").Find(&discussions).Error; err != nil {
		return nil, 0, err
	}

	return discussions, total, nil
}

// FindStudyCaseByID finds a study case by ID (for notification purposes)
func (r *DiscussionRepository) FindStudyCaseByID(ctx context.Context, id string) (*domain.StudyCase, error) {
	studyCase := &domain.StudyCase{}
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(studyCase).Error; err != nil {
		return nil, err
	}
	return studyCase, nil
}

// FindLessonByID finds a lesson by ID (for notification purposes)
func (r *DiscussionRepository) FindLessonByID(ctx context.Context, id string) (*domain.Lesson, error) {
	lesson := &domain.Lesson{}
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(lesson).Error; err != nil {
		return nil, err
	}
	return lesson, nil
}

// Update updates discussion
func (r *DiscussionRepository) Update(ctx context.Context, discussion *domain.Discussion) error {
	return r.db.WithContext(ctx).Model(discussion).Updates(discussion).Error
}

// DeleteByID deletes discussion and cascade replies
func (r *DiscussionRepository) DeleteByID(ctx context.Context, discussionID string) error {
	return r.db.WithContext(ctx).Where("id = ?", discussionID).Delete(&domain.Discussion{}).Error
}
