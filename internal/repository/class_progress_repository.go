package repository

import (
	"context"
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

type ClassProgressRepository struct {
	db *gorm.DB
}

func NewClassProgressRepository(db *gorm.DB) *ClassProgressRepository {
	return &ClassProgressRepository{db: db}
}

func (r *ClassProgressRepository) Create(ctx context.Context, progress *domain.ClassProgress) error {
	return r.db.WithContext(ctx).Create(progress).Error
}

func (r *ClassProgressRepository) FindByUserAndClass(ctx context.Context, userID, classID string) (*domain.ClassProgress, error) {
	progress := &domain.ClassProgress{}
	if err := r.db.WithContext(ctx).Where("user_id = ? AND class_id = ?", userID, classID).First(progress).Error; err != nil {
		return nil, err
	}
	return progress, nil
}

func (r *ClassProgressRepository) Update(ctx context.Context, progress *domain.ClassProgress) error {
	return r.db.WithContext(ctx).Save(progress).Error
}

func (r *ClassProgressRepository) ListByUserID(ctx context.Context, userID string) ([]domain.ClassProgress, error) {
	var progresses []domain.ClassProgress
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Preload("Class").Find(&progresses).Error; err != nil {
		return nil, err
	}
	return progresses, nil
}
