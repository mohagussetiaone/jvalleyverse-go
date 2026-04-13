package repository

import (
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

// CategoryRepository handles category data access
type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create creates new category
func (r *CategoryRepository) Create(category *domain.Category) error {
	return r.db.Create(category).Error
}

// ListAll lists all categories
func (r *CategoryRepository) ListAll() ([]domain.Category, error) {
	var categories []domain.Category
	err := r.db.Find(&categories).Error
	return categories, err
}
