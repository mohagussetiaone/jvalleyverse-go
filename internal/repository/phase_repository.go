package repository

import (
	"context"
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

// PhaseRepository handles phase data access
type PhaseRepository struct {
	db *gorm.DB
}

func NewPhaseRepository(db *gorm.DB) *PhaseRepository {
	return &PhaseRepository{db: db}
}

// Create creates a new phase
func (r *PhaseRepository) Create(ctx context.Context, phase *domain.Phase) error {
	return r.db.WithContext(ctx).Create(phase).Error
}

// FindByID finds phase by ID, optionally preloading classes
func (r *PhaseRepository) FindByID(ctx context.Context, phaseID string) (*domain.Phase, error) {
	phase := &domain.Phase{}
	err := r.db.WithContext(ctx).
		Where("id = ?", phaseID).
		Preload("Project").
		First(phase).Error
	return phase, err
}

// FindByIDWithClasses finds phase including its ordered classes
func (r *PhaseRepository) FindByIDWithClasses(ctx context.Context, phaseID string) (*domain.Phase, error) {
	phase := &domain.Phase{}
	err := r.db.WithContext(ctx).
		Where("id = ?", phaseID).
		Preload("Project").
		Preload("Classes", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC")
		}).
		First(phase).Error
	return phase, err
}

// ListByProjectID lists all phases for a project ordered by order_index
func (r *PhaseRepository) ListByProjectID(ctx context.Context, projectID string) ([]domain.Phase, error) {
	var phases []domain.Phase
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("order_index ASC").
		Find(&phases).Error
	return phases, err
}

// ListByProjectIDWithClasses lists all phases with their classes for a project
func (r *PhaseRepository) ListByProjectIDWithClasses(ctx context.Context, projectID string) ([]domain.Phase, error) {
	var phases []domain.Phase
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("order_index ASC").
		Preload("Classes", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC")
		}).
		Find(&phases).Error
	return phases, err
}

// Update updates a phase
func (r *PhaseRepository) Update(ctx context.Context, phase *domain.Phase) error {
	return r.db.WithContext(ctx).Model(phase).Updates(phase).Error
}

// DeleteByID soft-deletes a phase (cascades to classes via DB constraint)
func (r *PhaseRepository) DeleteByID(ctx context.Context, phaseID string) error {
	return r.db.WithContext(ctx).Where("id = ?", phaseID).Delete(&domain.Phase{}).Error
}

// ExistsInProject checks if a phase belongs to a specific project
func (r *PhaseRepository) ExistsInProject(ctx context.Context, phaseID, projectID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Phase{}).
		Where("id = ? AND project_id = ?", phaseID, projectID).
		Count(&count).Error
	return count > 0, err
}
