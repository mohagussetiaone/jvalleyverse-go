package repository

import (
	"context"
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

type LessonDetailRepository struct {
	db *gorm.DB
}

func NewLessonDetailRepository(db *gorm.DB) *LessonDetailRepository {
	return &LessonDetailRepository{db: db}
}

func (r *LessonDetailRepository) Create(ctx context.Context, detail *domain.LessonDetail) error {
	return r.db.WithContext(ctx).Create(detail).Error
}

func (r *LessonDetailRepository) FindByClassID(ctx context.Context, classID string) (*domain.LessonDetail, error) {
	detail := &domain.LessonDetail{}
	if err := r.db.WithContext(ctx).Where("class_id = ?", classID).First(detail).Error; err != nil {
		return nil, err
	}
	return detail, nil
}

func (r *LessonDetailRepository) Update(ctx context.Context, detail *domain.LessonDetail) error {
	return r.db.WithContext(ctx).Save(detail).Error
}
