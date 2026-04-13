package repository

import (
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

// ClassRepository handles class data access
type ClassRepository struct {
	db *gorm.DB
}

func NewClassRepository() *ClassRepository {
	return &ClassRepository{db: db}
}

// Create creates new class
func (r *ClassRepository) Create(class *domain.Class) error {
	return r.db.Create(class).Error
}

// FindByID finds class with project, admin and details
func (r *ClassRepository) FindByID(classID uint) (*domain.Class, error) {
	class := &domain.Class{}
	if err := r.db.Preload("Project").Preload("Admin").Preload("Details").First(class, classID).Error; err != nil {
		return nil, err
	}
	return class, nil
}

// FindBySlug finds class by slug within a project
func (r *ClassRepository) FindBySlug(projectID uint, slug string) (*domain.Class, error) {
	class := &domain.Class{}
	if err := r.db.Where("project_id = ? AND slug = ?", projectID, slug).
		Preload("Project").
		Preload("Admin").
		Preload("Details").
		Preload("NextClass").
		First(class).Error; err != nil {
		return nil, err
	}
	return class, nil
}

// FindNextClass finds the next class in sequence
func (r *ClassRepository) FindNextClass(nextClassID uint) (*domain.Class, error) {
	class := &domain.Class{}
	if err := r.db.First(class, nextClassID).Error; err != nil {
		return nil, err
	}
	return class, nil
}

// ListByProjectID lists classes in a project
func (r *ClassRepository) ListByProjectID(projectID uint, page, limit int) ([]domain.Class, int64, error) {
	var classes []domain.Class
	var total int64

	if err := r.db.Model(&domain.Class{}).Where("project_id = ?", projectID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.Where("project_id = ?", projectID).Preload("Admin").Offset(offset).Limit(limit).Find(&classes).Error; err != nil {
		return nil, 0, err
	}

	return classes, total, nil
}

// ListAll lists all classes with pagination
func (r *ClassRepository) ListAll(page, limit int, difficulty *string) ([]domain.Class, int64, error) {
	var classes []domain.Class
	var total int64

	query := r.db
	if difficulty != nil {
		query = query.Where("difficulty = ?", *difficulty)
	}

	if err := query.Model(&domain.Class{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Preload("Project").Preload("Admin").Offset(offset).Limit(limit).Find(&classes).Error; err != nil {
		return nil, 0, err
	}

	return classes, total, nil
}

// Update updates class
func (r *ClassRepository) Update(class *domain.Class) error {
	return r.db.Model(class).Updates(class).Error
}

// DeleteByID deletes class and cascade discussions/certificates
func (r *ClassRepository) DeleteByID(classID uint) error {
	return r.db.Delete(&domain.Class{}, classID).Error
}
