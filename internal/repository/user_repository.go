package repository

import (
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

// UserRepository handles user data access
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db: db}
}

// FindByID finds user by ID
func (r *UserRepository) FindByID(userID uint) (*domain.User, error) {
	user := &domain.User{}
	if err := r.db.First(user, userID).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// FindByEmail finds user by email
func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	if err := r.db.Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Create creates new user
func (r *UserRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

// Update updates user details
func (r *UserRepository) Update(user *domain.User) error {
	return r.db.Model(&domain.User{}).Where("id = ?", user.ID).Updates(user).Error
}

// UpdatePoints adds points to user (atomic operation)
func (r *UserRepository) UpdatePoints(userID uint, points int) error {
	return r.db.Model(&domain.User{}).
		Where("id = ?", userID).
		Update("total_points", gorm.Expr("total_points + ?", points)).
		Error
}

// UpdateLevel updates user level
func (r *UserRepository) UpdateLevel(userID uint, level int) error {
	return r.db.Model(&domain.User{}).
		Where("id = ?", userID).
		Update("level", level).
		Error
}

// ListAll lists all users with pagination
func (r *UserRepository) ListAll(page, limit int) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	if err := r.db.Model(&domain.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
