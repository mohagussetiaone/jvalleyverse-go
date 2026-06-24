package repository

import (
	"context"
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

type LessonProgressRepository struct {
	db *gorm.DB
}

func NewLessonProgressRepository(db *gorm.DB) *LessonProgressRepository {
	return &LessonProgressRepository{db: db}
}

func (r *LessonProgressRepository) Create(ctx context.Context, progress *domain.LessonProgress) error {
	return r.db.WithContext(ctx).Create(progress).Error
}

func (r *LessonProgressRepository) FindByUserAndLesson(ctx context.Context, userID, lessonID string) (*domain.LessonProgress, error) {
	progress := &domain.LessonProgress{}
	if err := r.db.WithContext(ctx).Where("user_id = ? AND lesson_id = ?", userID, lessonID).First(progress).Error; err != nil {
		return nil, err
	}
	return progress, nil
}

func (r *LessonProgressRepository) Update(ctx context.Context, progress *domain.LessonProgress) error {
	return r.db.WithContext(ctx).Save(progress).Error
}

func (r *LessonProgressRepository) ListByUserID(ctx context.Context, userID string) ([]domain.LessonProgress, error) {
	var progresses []domain.LessonProgress
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Preload("Lesson").Find(&progresses).Error; err != nil {
		return nil, err
	}
	return progresses, nil
}
