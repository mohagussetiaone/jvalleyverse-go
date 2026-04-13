package repository

import (
	"jvalleyverse/internal/domain"
	"gorm.io/gorm"
)

type ClassDetailRepository struct {
	db *gorm.DB
}

func NewClassDetailRepository() *ClassDetailRepository {
	return &ClassDetailRepository{db: db}
}

func (r *ClassDetailRepository) Create(detail *domain.ClassDetail) error {
	return r.db.Create(detail).Error
}

func (r *ClassDetailRepository) FindByClassID(classID uint) (*domain.ClassDetail, error) {
	detail := &domain.ClassDetail{}
	if err := r.db.Where("class_id = ?", classID).First(detail).Error; err != nil {
		return nil, err
	}
	return detail, nil
}

func (r *ClassDetailRepository) Update(detail *domain.ClassDetail) error {
	return r.db.Save(detail).Error
}
