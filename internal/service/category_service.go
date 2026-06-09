package service

import (
	"context"
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/repository"
)

// ICategoryService defines the business logic for category management
type ICategoryService interface {
	ListCategories(ctx context.Context) ([]domain.Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*domain.Category, error)
	GetCategoryByID(ctx context.Context, id string) (*domain.Category, error)
	ListProjectsByCategoryID(ctx context.Context, categoryID string) ([]domain.Project, error)
	CreateCategory(ctx context.Context, name, slug, description string) (*domain.Category, error)
	UpdateCategory(ctx context.Context, id, name, slug, description string) (*domain.Category, error)
	DeleteCategory(ctx context.Context, id string) error
}

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

// ListCategories returns all categories (public)
func (s *CategoryService) ListCategories(ctx context.Context) ([]domain.Category, error) {
	return s.categoryRepo.ListAll(ctx)
}

// GetCategoryBySlug returns category by slug (public)
func (s *CategoryService) GetCategoryBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	return s.categoryRepo.FindBySlug(ctx, slug)
}

// GetCategoryByID returns category by ID
func (s *CategoryService) GetCategoryByID(ctx context.Context, id string) (*domain.Category, error) {
	return s.categoryRepo.FindByID(ctx, id)
}

// ListProjectsByCategoryID returns projects in a category (public)
func (s *CategoryService) ListProjectsByCategoryID(ctx context.Context, categoryID string) ([]domain.Project, error) {
	return s.categoryRepo.ListProjectsByCategoryID(ctx, categoryID)
}

// CreateCategory creates a new category (admin only)
func (s *CategoryService) CreateCategory(ctx context.Context, name, slug, description string) (*domain.Category, error) {
	category := &domain.Category{
		Name:        name,
		Slug:        slug,
		Description: description,
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// UpdateCategory updates an existing category (admin only)
func (s *CategoryService) UpdateCategory(ctx context.Context, id, name, slug, description string) (*domain.Category, error) {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		category.Name = name
	}
	if slug != "" {
		category.Slug = slug
	}
	if description != "" {
		category.Description = description
	}

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory soft-deletes a category (admin only)
func (s *CategoryService) DeleteCategory(ctx context.Context, id string) error {
	_, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return s.categoryRepo.DeleteByID(ctx, id)
}
