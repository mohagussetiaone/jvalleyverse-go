package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"jvalleyverse/internal/domain"
)

type IBlogRepository interface {
	Create(ctx context.Context, blog *domain.Blog) error
	FindByID(ctx context.Context, id string) (*domain.Blog, error)
	FindBySlug(ctx context.Context, slug string) (*domain.Blog, error)
	ListAll(ctx context.Context, page, limit int, search, status, categoryID, tag string) ([]domain.Blog, int64, error)
	ListByUserID(ctx context.Context, userID string, page, limit int, status string) ([]domain.Blog, int64, error)
	Update(ctx context.Context, blog *domain.Blog) error
	Delete(ctx context.Context, id string) error
}

type BlogRepository struct {
	db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) IBlogRepository {
	return &BlogRepository{db: db}
}

func (r *BlogRepository) Create(ctx context.Context, blog *domain.Blog) error {
	return r.db.WithContext(ctx).Create(blog).Error
}

func (r *BlogRepository) FindByID(ctx context.Context, id string) (*domain.Blog, error) {
	var blog domain.Blog
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Category").
		Where("id = ?", id).
		First(&blog).Error
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

func (r *BlogRepository) FindBySlug(ctx context.Context, slug string) (*domain.Blog, error) {
	var blog domain.Blog
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Category").
		Where("slug = ?", slug).
		First(&blog).Error
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

func (r *BlogRepository) buildBase(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Model(&domain.Blog{}).Preload("User").Preload("Category")
}

func (r *BlogRepository) ListAll(ctx context.Context, page, limit int, search, status, categoryID, tag string) ([]domain.Blog, int64, error) {
	// Count query
	countTx := r.buildBase(ctx)
	if search != "" {
		like := "%" + search + "%"
		countTx = countTx.Where("title ILIKE ? OR description ILIKE ?", like, like)
	}
	if status != "" {
		countTx = countTx.Where("status = ?", status)
	}
	if categoryID != "" {
		countTx = countTx.Where("category_id = ?", categoryID)
	}
	if tag != "" {
		like := fmt.Sprintf(`%%"%s"%%`, tag)
		countTx = countTx.Where("tags::text ILIKE ?", like)
	}

	var total int64
	if err := countTx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Data query (fresh chain)
	dataTx := r.buildBase(ctx)
	if search != "" {
		like := "%" + search + "%"
		dataTx = dataTx.Where("title ILIKE ? OR description ILIKE ?", like, like)
	}
	if status != "" {
		dataTx = dataTx.Where("status = ?", status)
	}
	if categoryID != "" {
		dataTx = dataTx.Where("category_id = ?", categoryID)
	}
	if tag != "" {
		like := fmt.Sprintf(`%%"%s"%%`, tag)
		dataTx = dataTx.Where("tags::text ILIKE ?", like)
	}

	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	var blogs []domain.Blog
	err := dataTx.Order("created_at DESC").Offset(offset).Limit(limit).Find(&blogs).Error
	return blogs, total, err
}

func (r *BlogRepository) ListByUserID(ctx context.Context, userID string, page, limit int, status string) ([]domain.Blog, int64, error) {
	var total int64

	countTx := r.db.WithContext(ctx).Model(&domain.Blog{}).Where("user_id = ?", userID)
	if status != "" {
		countTx = countTx.Where("status = ?", status)
	}
	if err := countTx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	dataTx := r.db.WithContext(ctx).
		Preload("User").
		Preload("Category").
		Where("user_id = ?", userID)
	if status != "" {
		dataTx = dataTx.Where("status = ?", status)
	}

	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	var blogs []domain.Blog
	err := dataTx.Order("created_at DESC").Offset(offset).Limit(limit).Find(&blogs).Error
	return blogs, total, err
}

func (r *BlogRepository) Update(ctx context.Context, blog *domain.Blog) error {
	return r.db.WithContext(ctx).Save(blog).Error
}

func (r *BlogRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Blog{}).Error
}
