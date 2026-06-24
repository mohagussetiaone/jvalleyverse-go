package service

import (
	"context"
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/dto"
	"jvalleyverse/internal/repository"
)

// ICategoryService defines the business logic for category management
type ICategoryService interface {
	ListCategories(ctx context.Context) ([]dto.CategoryBrief, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*dto.CategoryWithCourses, error)
	GetCategoryByID(ctx context.Context, id string) (*domain.Category, error)
	ListCoursesByCategoryID(ctx context.Context, categoryID string) ([]domain.Course, error)
	ListCoursesByCategoryIDWithEnrollment(ctx context.Context, userID, categoryID string) ([]dto.CourseListItem, error)
	CreateCategory(ctx context.Context, name, slug, description string) (*domain.Category, error)
	UpdateCategory(ctx context.Context, id, name, slug, description string) (*domain.Category, error)
	DeleteCategory(ctx context.Context, id string) error
}

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
	enrollRepo   *repository.EnrollmentRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository, enrollRepo *repository.EnrollmentRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		enrollRepo:   enrollRepo,
	}
}

// ListCategories returns all categories (public)
func (s *CategoryService) ListCategories(ctx context.Context) ([]dto.CategoryBrief, error) {
	categories, err := s.categoryRepo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]dto.CategoryBrief, len(categories))
	for i, c := range categories {
		result[i] = dto.ToCategoryBrief(c)
	}
	return result, nil
}

// GetCategoryBySlug returns category by slug (public)
func (s *CategoryService) GetCategoryBySlug(ctx context.Context, slug string) (*dto.CategoryWithCourses, error) {
	category, err := s.categoryRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	courses := make([]dto.CourseListItem, len(category.Courses))
	for i, c := range category.Courses {
		courses[i] = dto.CourseToListItem(c)
	}

	return &dto.CategoryWithCourses{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		Courses:     courses,
	}, nil
}

// GetCategoryByID returns category by ID
func (s *CategoryService) GetCategoryByID(ctx context.Context, id string) (*domain.Category, error) {
	return s.categoryRepo.FindByID(ctx, id)
}

// ListCoursesByCategoryID returns courses in a category (public)
func (s *CategoryService) ListCoursesByCategoryID(ctx context.Context, categoryID string) ([]domain.Course, error) {
	return s.categoryRepo.ListCoursesByCategoryID(ctx, categoryID)
}

// ListCoursesByCategoryIDWithEnrollment returns courses with enrollment status for the user
func (s *CategoryService) ListCoursesByCategoryIDWithEnrollment(ctx context.Context, userID, categoryID string) ([]dto.CourseListItem, error) {
	courses, err := s.categoryRepo.ListCoursesByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.CourseListItem, len(courses))
	for i, c := range courses {
		item := dto.CourseToListItem(c)
		enrolled, _ := s.enrollRepo.Exists(ctx, userID, c.ID)
		item.IsEnrolled = enrolled
		result[i] = item
	}

	return result, nil
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
