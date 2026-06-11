package service

import (
	"context"
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/repository"
)

// IPhaseService defines business logic for managing learning phases
type IPhaseService interface {
	// Admin
	CreatePhase(ctx context.Context, adminID, projectID string, title, description string, orderIndex int) (*domain.Phase, error)
	UpdatePhase(ctx context.Context, adminID, phaseID string, title, description string, orderIndex int) (*domain.Phase, error)
	DeletePhase(ctx context.Context, adminID, phaseID string) error

	// Public / protected
	GetPhase(ctx context.Context, phaseID string) (*domain.Phase, error)
	ListPhasesByProject(ctx context.Context, projectID string) ([]domain.Phase, error)
	GetProjectWithPhases(ctx context.Context, projectID string) (*domain.Project, error)
}

type PhaseService struct {
	phaseRepo   *repository.PhaseRepository
	projectRepo *repository.ProjectRepository
}

func NewPhaseService(
	phaseRepo *repository.PhaseRepository,
	projectRepo *repository.ProjectRepository,
) *PhaseService {
	return &PhaseService{
		phaseRepo:   phaseRepo,
		projectRepo: projectRepo,
	}
}

// CreatePhase creates a new phase under a project (admin only)
func (s *PhaseService) CreatePhase(ctx context.Context, adminID, projectID string, title, description string, orderIndex int) (*domain.Phase, error) {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, domain.ErrProjectNotFound
	}
	if project.AdminID != adminID {
		return nil, domain.ErrForbidden
	}

	phase := &domain.Phase{
		ProjectID:   projectID,
		Title:       title,
		Description: description,
		OrderIndex:  orderIndex,
	}

	if err := s.phaseRepo.Create(ctx, phase); err != nil {
		return nil, err
	}

	return phase, nil
}

// UpdatePhase updates phase metadata (admin only)
func (s *PhaseService) UpdatePhase(ctx context.Context, adminID, phaseID string, title, description string, orderIndex int) (*domain.Phase, error) {
	phase, err := s.phaseRepo.FindByID(ctx, phaseID)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	project, err := s.projectRepo.FindByID(ctx, phase.ProjectID)
	if err != nil {
		return nil, domain.ErrProjectNotFound
	}
	if project.AdminID != adminID {
		return nil, domain.ErrForbidden
	}

	if title != "" {
		phase.Title = title
	}
	if description != "" {
		phase.Description = description
	}
	if orderIndex >= 0 {
		phase.OrderIndex = orderIndex
	}

	if err := s.phaseRepo.Update(ctx, phase); err != nil {
		return nil, err
	}

	return phase, nil
}

// DeletePhase deletes a phase and cascades to its classes (admin only)
func (s *PhaseService) DeletePhase(ctx context.Context, adminID, phaseID string) error {
	phase, err := s.phaseRepo.FindByID(ctx, phaseID)
	if err != nil {
		return domain.ErrNotFound
	}
	project, err := s.projectRepo.FindByID(ctx, phase.ProjectID)
	if err != nil {
		return domain.ErrProjectNotFound
	}
	if project.AdminID != adminID {
		return domain.ErrForbidden
	}
	return s.phaseRepo.DeleteByID(ctx, phaseID)
}

// GetPhase returns a phase with its classes
func (s *PhaseService) GetPhase(ctx context.Context, phaseID string) (*domain.Phase, error) {
	return s.phaseRepo.FindByIDWithClasses(ctx, phaseID)
}

// ListPhasesByProject returns all phases (without classes) for a project
func (s *PhaseService) ListPhasesByProject(ctx context.Context, projectID string) ([]domain.Phase, error) {
	return s.phaseRepo.ListByProjectID(ctx, projectID)
}

// GetProjectWithPhases returns a project with all its phases and classes (public)
func (s *PhaseService) GetProjectWithPhases(ctx context.Context, projectID string) (*domain.Project, error) {
	return s.projectRepo.FindByIDWithPhases(ctx, projectID)
}
