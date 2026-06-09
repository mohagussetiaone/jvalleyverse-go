package repository

import (
	"context"
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

// UserLevelRepository handles user level data access
type UserLevelRepository struct {
	db *gorm.DB
}

func NewUserLevelRepository(db *gorm.DB) *UserLevelRepository {
	return &UserLevelRepository{db: db}
}

// Create creates new user level record
func (r *UserLevelRepository) Create(ctx context.Context, level *domain.UserLevel) error {
	return r.db.WithContext(ctx).Create(level).Error
}

// FindByLevel finds level config by level number
func (r *UserLevelRepository) FindByLevel(ctx context.Context, levelNum int) (*domain.UserLevel, error) {
	level := &domain.UserLevel{}
	if err := r.db.WithContext(ctx).Where("level = ?", levelNum).First(level).Error; err != nil {
		return nil, err
	}
	return level, nil
}

// ListAll lists all level configuration
func (r *UserLevelRepository) ListAll(ctx context.Context) ([]domain.UserLevel, error) {
	var levels []domain.UserLevel
	err := r.db.WithContext(ctx).Order("level ASC").Find(&levels).Error
	return levels, err
}

// Update updates level record
func (r *UserLevelRepository) Update(ctx context.Context, level *domain.UserLevel) error {
	return r.db.WithContext(ctx).Model(level).Updates(level).Error
}
