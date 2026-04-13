package service

import (
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/repository"
)

type ProjectService struct {
	projectRepo *repository.ProjectRepository
	classRepo   *repository.ClassRepository
	userRepo    *repository.UserRepository
}

func NewProjectService() *ProjectService {
	return &ProjectService{
		projectRepo: repository.NewProjectRepository(),
		classRepo:   repository.NewClassRepository(),
		userRepo:    repository.NewUserRepository(),
	}
}

// CreateProject creates new learning project (admin only)
func (s *ProjectService) CreateProject(adminID uint, title, description, thumbnail string, categoryID uint) (*domain.Project, error) {
	project := &domain.Project{
		Title:       title,
		Description: description,
		Thumbnail:   thumbnail,
		AdminID:     adminID,
		CategoryID:  categoryID,
		Visibility:  "public",
	}

	if err := s.projectRepo.Create(project); err != nil {
		return nil, err
	}

	return project, nil
}

// ListProjects returns all projects with pagination
func (s *ProjectService) ListProjects(page, limit int) ([]map[string]interface{}, error) {
	projects, _, err := s.projectRepo.ListAll(page, limit)
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
			"admin_id":    p.AdminID,
			"admin_name":  p.Admin.Name,
			"visibility":  p.Visibility,
			"created_at":  p.CreatedAt,
		}
	}

	return result, nil
}

// GetProject returns specific project with classes
func (s *ProjectService) GetProject(projectID uint) (map[string]interface{}, error) {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, err
	}

	// Get classes in project
	classes, _, err := s.classRepo.ListByProjectID(projectID, 1, 100)
	if err != nil {
		return nil, err
	}

	classData := make([]map[string]interface{}, len(classes))
	for i, c := range classes {
		classData[i] = map[string]interface{}{
			"id":        c.ID,
			"title":     c.Title,
			"difficulty": c.Difficulty,
		}
	}

	return map[string]interface{}{
		"id":          project.ID,
		"title":       project.Title,
		"description": project.Description,
		"thumbnail":   project.Thumbnail,
		"admin_id":    project.AdminID,
		"admin_name":  project.Admin.Name,
		"visibility":  project.Visibility,
		"classes":     classData,
		"created_at":  project.CreatedAt,
	}, nil
}

// UpdateProject updates project details (admin only)
func (s *ProjectService) UpdateProject(projectID, adminID uint, title, description string, visibility string) error {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return err
	}

	// Only admin owner can update
	if project.AdminID != adminID {
		return nil
	}

	project.Title = title
	project.Description = description
	project.Visibility = visibility

	return s.projectRepo.Update(project)
}

// DeleteProject deletes project and cascade delete classes (admin only)
func (s *ProjectService) DeleteProject(projectID, adminID uint) error {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return err
	}

	// Only admin owner can delete
	if project.AdminID != adminID {
		return nil
	}

	return s.projectRepo.DeleteByID(projectID)
}
