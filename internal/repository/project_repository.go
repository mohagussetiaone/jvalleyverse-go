package repository

import (
	"context"
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

// ProjectRepository handles project data access
type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// Create creates new project
func (r *ProjectRepository) Create(ctx context.Context, project *domain.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

// FindByID finds project with admin and classes
func (r *ProjectRepository) FindByID(ctx context.Context, projectID string) (*domain.Project, error) {
	project := &domain.Project{}
	if err := r.db.WithContext(ctx).Where("id = ?", projectID).
		Preload("Admin").
		Preload("Category").
		First(project).Error; err != nil {
		return nil, err
	}
	return project, nil
}

// FindByIDWithPhases finds project including all phases and their classes
func (r *ProjectRepository) FindByIDWithPhases(ctx context.Context, projectID string) (*domain.Project, error) {
	project := &domain.Project{}
	if err := r.db.WithContext(ctx).Where("id = ?", projectID).
		Preload("Admin").
		Preload("Category").
		Preload("Phases", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC")
		}).
		Preload("Phases.Classes", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC")
		}).
		Preload("Phases.Classes.Details").
		Preload("Phases.Classes.NextClass").
		First(project).Error; err != nil {
		return nil, err
	}
	return project, nil
}

// ListAll lists all projects with pagination
func (r *ProjectRepository) ListAll(ctx context.Context, page, limit int) ([]domain.Project, int64, error) {
	var projects []domain.Project
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Project{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).
		Preload("Admin").
		Preload("Category").
		Preload("Phases").
		Offset(offset).
		Limit(limit).
		Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (r *ProjectRepository) ListPublic(
	ctx context.Context,
	page,
	limit int,
) ([]domain.Project, int64, error) {

	var projects []domain.Project
	var total int64

	query := r.db.WithContext(ctx).
		Model(&domain.Project{}).
		Where("visibility = ?", "public")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	if err := query.
		Preload("Admin").
		Preload("Category").
		Preload("Phases").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

// ListByAdminID lists projects created by admin
func (r *ProjectRepository) ListByAdminID(ctx context.Context, adminID string, page, limit int) ([]domain.Project, int64, error) {
	var projects []domain.Project
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Project{}).Where("admin_id = ?", adminID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).
		Where("admin_id = ?", adminID).
		Preload("Category").
		Preload("Phases").
		Offset(offset).
		Limit(limit).
		Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

// Update updates project
func (r *ProjectRepository) Update(ctx context.Context, project *domain.Project) error {
	return r.db.WithContext(ctx).Model(project).Updates(project).Error
}

// DeleteByID deletes project and cascade classes
func (r *ProjectRepository) DeleteByID(ctx context.Context, projectID string) error {
	return r.db.WithContext(ctx).Where("id = ?", projectID).Delete(&domain.Project{}).Error
}
