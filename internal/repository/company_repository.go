package repository

import (
	"context"

	"gorm.io/gorm"

	"jvalleyverse/internal/domain"
)

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

// Get returns the single company record (creates default if not exists)
func (r *CompanyRepository) Get(ctx context.Context) (*domain.Company, error) {
	var company domain.Company
	err := r.db.WithContext(ctx).First(&company).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return empty record if not seeded yet
			return &domain.Company{}, nil
		}
		return nil, err
	}
	return &company, nil
}

// Upsert creates or updates the single company record
func (r *CompanyRepository) Upsert(ctx context.Context, company *domain.Company) error {
	// Try to find existing record
	var existing domain.Company
	err := r.db.WithContext(ctx).First(&existing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return r.db.WithContext(ctx).Create(company).Error
		}
		return err
	}
	// Update existing record
	company.ID = existing.ID
	company.CreatedAt = existing.CreatedAt
	company.UpdatedAt = existing.UpdatedAt
	return r.db.WithContext(ctx).Save(company).Error
}
