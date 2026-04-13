package repository

import (
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

// UserLevelRepository handles user level data access
type UserLevelRepository struct {
	db *gorm.DB
}

func NewUserLevelRepository() *UserLevelRepository {
	return &UserLevelRepository{db: db}
}

// Create creates new user level record
func (r *UserLevelRepository) Create(level *domain.UserLevel) error {
	return r.db.Create(level).Error
}

// FindByID finds level by level number
func (r *UserLevelRepository) FindByID(levelID uint) (*domain.UserLevel, error) {
	level := &domain.UserLevel{}
	if err := r.db.First(level, levelID).Error; err != nil {
		return nil, err
	}
	return level, nil
}

// FindByUserID finds current level of user
func (r *UserLevelRepository) FindByUserID(userID uint) (*domain.UserLevel, error) {
	// This doesn't exist in schema - UserLevel is configuration
	// For user's current level, it's in User.Level field
	return nil, nil
}

// ListAll lists level configuration
func (r *UserLevelRepository) ListAll() ([]domain.UserLevel, error) {
	var levels []domain.UserLevel
	err := r.db.Order("level ASC").Find(&levels).Error
	return levels, err
}

// Update updates level record
func (r *UserLevelRepository) Update(level *domain.UserLevel) error {
	return r.db.Model(level).Updates(level).Error
}
