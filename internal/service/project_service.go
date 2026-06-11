package service

import (
	"context"
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/repository"
)

// IProjectService defines the business logic for managing learning projects
type IProjectService interface {
	CreateProject(ctx context.Context, adminID string, title, description, thumbnail string, categoryID string) (*domain.Project, error)
	ListProjects(ctx context.Context, page, limit int) ([]map[string]interface{}, error)
	GetProject(ctx context.Context, projectID string) (map[string]interface{}, error)
	UpdateProject(ctx context.Context, projectID, adminID string, title, description string, visibility string) error
	DeleteProject(ctx context.Context, projectID, adminID string) error
	ListPublicProjects(ctx context.Context, page, limit int) ([]map[string]interface{}, error)
}

type ProjectService struct {
	projectRepo *repository.ProjectRepository
	classRepo   *repository.ClassRepository
	userRepo    *repository.UserRepository
}

func NewProjectService(
	projectRepo *repository.ProjectRepository,
	classRepo *repository.ClassRepository,
	userRepo *repository.UserRepository,
) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		classRepo:   classRepo,
		userRepo:    userRepo,
	}
}

// CreateProject creates new learning project (admin only)
func (s *ProjectService) CreateProject(ctx context.Context, adminID string, title, description, thumbnail string, categoryID string) (*domain.Project, error) {
	if title == "" || categoryID == "" {
		return nil, domain.ErrInvalidInput
	}

	project := &domain.Project{
		Title:       title,
		Description: description,
		Thumbnail:   thumbnail,
		AdminID:     adminID,
		CategoryID:  categoryID,
		Visibility:  "public",
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

// ListProjects returns all projects with pagination
func (s *ProjectService) ListProjects(ctx context.Context, page, limit int) ([]map[string]interface{}, error) {
	projects, _, err := s.projectRepo.ListAll(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(projects))
	for i, p := range projects {
		result[i] = map[string]interface{}{
			"id":          p.ID,
			"title":       p.Title,
			"description": p.Description,
			"thumbnail":   p.Thumbnail,
			"category":    p.Category,
			"admin_id":    p.AdminID,
			"admin_name":  p.Admin.Name,
			"visibility":  p.Visibility,
			"phase_count": len(p.Phases),
			"created_at":  p.CreatedAt,
		}
	}

	return result, nil
}

// GetProject returns specific project with classes
func (s *ProjectService) GetProject(ctx context.Context, projectID string) (map[string]interface{}, error) {
	project, err := s.projectRepo.FindByIDWithPhases(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":          project.ID,
		"title":       project.Title,
		"description": project.Description,
		"thumbnail":   project.Thumbnail,
		"category":    project.Category,
		"admin_id":    project.AdminID,
		"admin_name":  project.Admin.Name,
		"visibility":  project.Visibility,
		"phases":      project.Phases,
		"created_at":  project.CreatedAt,
	}, nil
}

func (s *ProjectService) ListPublicProjects(
	ctx context.Context,
	page,
	limit int,
) ([]map[string]interface{}, error) {

	projects, _, err := s.projectRepo.ListPublic(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(projects))

	for i, p := range projects {
		result[i] = map[string]interface{}{
			"id":          p.ID,
			"title":       p.Title,
			"description": p.Description,
			"thumbnail":   p.Thumbnail,
			"category":    p.Category,
			"admin_name":  p.Admin.Name,
			"visibility":  p.Visibility,
			"phase_count": len(p.Phases),
			"created_at":  p.CreatedAt,
		}
	}

	return result, nil
}

// UpdateProject updates project details (admin only)
func (s *ProjectService) UpdateProject(ctx context.Context, projectID, adminID string, title, description string, visibility string) error {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return err
	}

	// Only admin owner can update
	if project.AdminID != adminID {
		return domain.ErrForbidden
	}

	if title != "" {
		project.Title = title
	}
	if description != "" {
		project.Description = description
	}
	if visibility != "" {
		project.Visibility = visibility
	}

	return s.projectRepo.Update(ctx, project)
}

// DeleteProject deletes project and cascade delete classes (admin only)
func (s *ProjectService) DeleteProject(ctx context.Context, projectID, adminID string) error {
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return err
	}

	// Only admin owner can delete
	if project.AdminID != adminID {
		return domain.ErrForbidden
	}

	return s.projectRepo.DeleteByID(ctx, projectID)
}
