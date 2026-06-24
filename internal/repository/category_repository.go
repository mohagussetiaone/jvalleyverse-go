package repository

import (
	"context"
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

// CategoryRepository handles category data access
type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create creates new category
func (r *CategoryRepository) Create(ctx context.Context, category *domain.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

// ListAll lists all categories
func (r *CategoryRepository) ListAll(ctx context.Context) ([]domain.Category, error) {
	var categories []domain.Category
	err := r.db.WithContext(ctx).Order("name ASC").Find(&categories).Error
	return categories, err
}

// FindByID finds category by ID
func (r *CategoryRepository) FindByID(ctx context.Context, id string) (*domain.Category, error) {
	category := &domain.Category{}
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

// FindBySlug finds category by slug
func (r *CategoryRepository) FindBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	category := &domain.Category{}
	if err := r.db.WithContext(ctx).
		Where("slug = ?", slug).
		Preload("Courses", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Preload("Courses.Admin").
		Preload("Courses.Mentor").
		First(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

// Update updates category
func (r *CategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	return r.db.WithContext(ctx).Model(category).Updates(category).Error
}

// DeleteByID soft-deletes category
func (r *CategoryRepository) DeleteByID(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Unscoped().
		Where("id = ?", id).
		Delete(&domain.Category{}).Error
}

// ListCoursesByCategoryID lists courses belonging to a category
func (r *CategoryRepository) ListCoursesByCategoryID(ctx context.Context, categoryID string) ([]domain.Course, error) {
	var courses []domain.Course
	err := r.db.WithContext(ctx).
		Where("category_id = ?", categoryID).
		Preload("Admin").
		Preload("Category").
		Preload("Mentor").
		Preload("Sections").
		Order("created_at DESC").
		Find(&courses).Error
	return courses, err
}
