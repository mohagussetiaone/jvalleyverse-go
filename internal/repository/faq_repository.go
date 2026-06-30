package repository

import (
	"context"

	"gorm.io/gorm"

	"jvalleyverse/internal/domain"
)

type FAQRepository struct {
	db *gorm.DB
}

func NewFAQRepository(db *gorm.DB) *FAQRepository {
	return &FAQRepository{db: db}
}

func (r *FAQRepository) Create(ctx context.Context, faq *domain.FAQ) error {
	return r.db.WithContext(ctx).Create(faq).Error
}

func (r *FAQRepository) FindByID(ctx context.Context, id string) (*domain.FAQ, error) {
	var faq domain.FAQ
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&faq).Error
	if err != nil {
		return nil, err
	}
	return &faq, nil
}

func (r *FAQRepository) ListAll(ctx context.Context, page, limit int) ([]domain.FAQ, int64, error) {
	var faqs []domain.FAQ
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.FAQ{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	err := r.db.WithContext(ctx).
		Order("order_index ASC, created_at ASC").
		Offset(offset).
		Limit(limit).
		Find(&faqs).Error
	return faqs, total, err
}

func (r *FAQRepository) ListActive(ctx context.Context) ([]domain.FAQ, error) {
	var faqs []domain.FAQ
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("order_index ASC, created_at ASC").
		Find(&faqs).Error
	return faqs, err
}

func (r *FAQRepository) Update(ctx context.Context, faq *domain.FAQ) error {
	return r.db.WithContext(ctx).Save(faq).Error
}

func (r *FAQRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.FAQ{}).Error
}
