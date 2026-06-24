package repository

import (
	"context"
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

// UserRepository handles user data access
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByID finds user by CUID string ID
func (r *UserRepository) FindByID(ctx context.Context, userID string) (*domain.User, error) {
	user := &domain.User{}
	if err := r.db.WithContext(ctx).Where("id = ?", userID).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// FindByEmail finds user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Create creates new user
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Update updates user details
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", user.ID).Updates(user).Error
}

// UpdatePoints adds points to user (atomic operation)
func (r *UserRepository) UpdatePoints(ctx context.Context, userID string, points int) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).
		Where("id = ?", userID).
		Update("total_points", gorm.Expr("total_points + ?", points)).
		Error
}

// UpdateLevel updates user level
func (r *UserRepository) UpdateLevel(ctx context.Context, userID string, level int) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).
		Where("id = ?", userID).
		Update("level", level).
		Error
}

// ListAll lists all users with pagination
func (r *UserRepository) ListAll(ctx context.Context, page, limit int) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ListByRole returns paginated users filtered by role
func (r *UserRepository) ListByRole(ctx context.Context, role string, page, limit int) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.User{}).Where("role = ?", role).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).Where("role = ?", role).Offset(offset).Limit(limit).Order("name ASC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetTopByPoints returns top N active users ordered by total_points descending
func (r *UserRepository) GetTopByPoints(ctx context.Context, limit int) ([]domain.User, error) {
	var users []domain.User
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("total_points DESC").
		Limit(limit).
		Find(&users).Error
	return users, err
}
