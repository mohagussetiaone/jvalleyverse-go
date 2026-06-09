package repository

import (
	"context"
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

type ClassDetailRepository struct {
	db *gorm.DB
}

func NewClassDetailRepository(db *gorm.DB) *ClassDetailRepository {
	return &ClassDetailRepository{db: db}
}

func (r *ClassDetailRepository) Create(ctx context.Context, detail *domain.ClassDetail) error {
	return r.db.WithContext(ctx).Create(detail).Error
}

func (r *ClassDetailRepository) FindByClassID(ctx context.Context, classID string) (*domain.ClassDetail, error) {
	detail := &domain.ClassDetail{}
	if err := r.db.WithContext(ctx).Where("class_id = ?", classID).First(detail).Error; err != nil {
		return nil, err
	}
	return detail, nil
}

func (r *ClassDetailRepository) Update(ctx context.Context, detail *domain.ClassDetail) error {
	return r.db.WithContext(ctx).Save(detail).Error
}
