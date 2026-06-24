package repository

import (
	"context"
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(ctx context.Context, review *domain.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *ReviewRepository) FindByID(ctx context.Context, id string) (*domain.Review, error) {
	review := &domain.Review{}
	if err := r.db.WithContext(ctx).Preload("User").Where("id = ?", id).First(review).Error; err != nil {
		return nil, err
	}
	return review, nil
}

func (r *ReviewRepository) ListByCourse(ctx context.Context, courseID string) ([]domain.Review, error) {
	var reviews []domain.Review
	if err := r.db.WithContext(ctx).
		Where("course_id = ?", courseID).
		Preload("User").
		Order("created_at DESC").
		Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *ReviewRepository) ListByLesson(ctx context.Context, lessonID string) ([]domain.Review, error) {
	var reviews []domain.Review
	if err := r.db.WithContext(ctx).
		Where("lesson_id = ?", lessonID).
		Preload("User").
		Order("created_at DESC").
		Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *ReviewRepository) Update(ctx context.Context, review *domain.Review) error {
	return r.db.WithContext(ctx).Model(review).Updates(review).Error
}

func (r *ReviewRepository) DeleteByID(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Review{}).Error
}
