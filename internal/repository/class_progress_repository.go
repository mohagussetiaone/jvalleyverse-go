package repository

import (
	"jvalleyverse/internal/domain"
	"gorm.io/gorm"
)

type ClassProgressRepository struct {
	db *gorm.DB
}

func NewClassProgressRepository() *ClassProgressRepository {
	return &ClassProgressRepository{db: db}
}

func (r *ClassProgressRepository) Create(progress *domain.ClassProgress) error {
	return r.db.Create(progress).Error
}

func (r *ClassProgressRepository) FindByUserAndClass(userID, classID uint) (*domain.ClassProgress, error) {
	progress := &domain.ClassProgress{}
	if err := r.db.Where("user_id = ? AND class_id = ?", userID, classID).First(progress).Error; err != nil {
		return nil, err
	}
	return progress, nil
}

func (r *ClassProgressRepository) Update(progress *domain.ClassProgress) error {
	return r.db.Save(progress).Error
}

func (r *ClassProgressRepository) ListByUserID(userID uint) ([]domain.ClassProgress, error) {
	var progresses []domain.ClassProgress
	if err := r.db.Where("user_id = ?", userID).Preload("Class").Find(&progresses).Error; err != nil {
		return nil, err
	}
	return progresses, nil
}
