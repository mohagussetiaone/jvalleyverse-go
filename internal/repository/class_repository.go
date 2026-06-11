package repository

import (
	"context"
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

// ClassRepository handles class data access
type ClassRepository struct {
	db *gorm.DB
}

func NewClassRepository(db *gorm.DB) *ClassRepository {
	return &ClassRepository{db: db}
}

// FindPublicByID finds public class with all relations
func (r *ClassRepository) FindPublicByID(ctx context.Context, classID string) (*domain.Class, error) {
	class := &domain.Class{}

	if err := r.db.WithContext(ctx).
		Where("id = ?", classID).
		Where("visibility = ?", "public").
		Preload("Project").
		Preload("Phase").
		Preload("Admin").
		Preload("Details").
		Preload("NextClass").
		First(class).Error; err != nil {
		return nil, err
	}

	return class, nil
}

// Create creates new class
func (r *ClassRepository) Create(ctx context.Context, class *domain.Class) error {
	return r.db.WithContext(ctx).Create(class).Error
}

// FindByID finds class with project, admin and details
func (r *ClassRepository) FindByID(ctx context.Context, classID string) (*domain.Class, error) {
	class := &domain.Class{}
	if err := r.db.WithContext(ctx).Where("id = ?", classID).
		Preload("Project").Preload("Admin").Preload("Details").First(class).Error; err != nil {
		return nil, err
	}
	return class, nil
}

// FindBySlug finds class by slug within a project
func (r *ClassRepository) FindBySlug(ctx context.Context, projectID string, slug string) (*domain.Class, error) {
	class := &domain.Class{}
	if err := r.db.WithContext(ctx).Where("project_id = ? AND slug = ?", projectID, slug).
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
func (r *ClassRepository) FindNextClass(ctx context.Context, nextClassID string) (*domain.Class, error) {
	class := &domain.Class{}
	if err := r.db.WithContext(ctx).Where("id = ?", nextClassID).First(class).Error; err != nil {
		return nil, err
	}
	return class, nil
}

// ListByProjectID lists classes in a project
func (r *ClassRepository) ListByProjectID(ctx context.Context, projectID string, limit, offset int) ([]domain.Class, int64, error) {
	var classes []domain.Class
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Class{}).Where("project_id = ?", projectID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Where("project_id = ?", projectID).
		Preload("Admin").
		Order("order_index ASC").
		Offset(offset).Limit(limit).Find(&classes).Error; err != nil {
		return nil, 0, err
	}

	return classes, total, nil
}

// ListByPhaseID lists classes belonging to a specific phase
func (r *ClassRepository) ListByPhaseID(ctx context.Context, phaseID string) ([]domain.Class, int64, error) {
	var classes []domain.Class
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Class{}).Where("phase_id = ?", phaseID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Where("phase_id = ?", phaseID).
		Preload("Admin").
		Order("order_index ASC").
		Find(&classes).Error; err != nil {
		return nil, 0, err
	}

	return classes, total, nil
}

// ListAll lists all classes with pagination
func (r *ClassRepository) ListAll(ctx context.Context, page, limit int, difficulty *string) ([]domain.Class, int64, error) {
	var classes []domain.Class
	var total int64

	query := r.db.WithContext(ctx)
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
func (r *ClassRepository) Update(ctx context.Context, class *domain.Class) error {
	return r.db.WithContext(ctx).Model(class).Updates(class).Error
}

// DeleteByID deletes class and cascade discussions/certificates
func (r *ClassRepository) DeleteByID(ctx context.Context, classID string) error {
	return r.db.WithContext(ctx).Where("id = ?", classID).Delete(&domain.Class{}).Error
}
