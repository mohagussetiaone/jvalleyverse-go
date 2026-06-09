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
	err := r.db.WithContext(ctx).Find(&categories).Error
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
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(category).Error; err != nil {
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
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Category{}).Error
}

// ListProjectsByCategoryID lists projects belonging to a category
func (r *CategoryRepository) ListProjectsByCategoryID(ctx context.Context, categoryID string) ([]domain.Project, error) {
	var projects []domain.Project
	err := r.db.WithContext(ctx).Where("category_id = ?", categoryID).Find(&projects).Error
	return projects, err
}
