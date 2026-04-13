package repository

import (
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

// ProjectRepository handles project data access
type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository() *ProjectRepository {
	return &ProjectRepository{db: db}
}

// Create creates new project
func (r *ProjectRepository) Create(project *domain.Project) error {
	return r.db.Create(project).Error
}

// FindByID finds project with admin and classes
func (r *ProjectRepository) FindByID(projectID uint) (*domain.Project, error) {
	project := &domain.Project{}
	if err := r.db.Preload("Admin").First(project, projectID).Error; err != nil {
		return nil, err
	}
	return project, nil
}

// ListAll lists all projects with pagination
func (r *ProjectRepository) ListAll(page, limit int) ([]domain.Project, int64, error) {
	var projects []domain.Project
	var total int64

	if err := r.db.Model(&domain.Project{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.Preload("Admin").Offset(offset).Limit(limit).Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

// ListByAdminID lists projects created by admin
func (r *ProjectRepository) ListByAdminID(adminID uint, page, limit int) ([]domain.Project, int64, error) {
	var projects []domain.Project
	var total int64

	if err := r.db.Model(&domain.Project{}).Where("admin_id = ?", adminID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.Where("admin_id = ?", adminID).Offset(offset).Limit(limit).Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

// Update updates project
func (r *ProjectRepository) Update(project *domain.Project) error {
	return r.db.Model(project).Updates(project).Error
}

// DeleteByID deletes project and cascade classes
func (r *ProjectRepository) DeleteByID(projectID uint) error {
	return r.db.Delete(&domain.Project{}, projectID).Error
}
