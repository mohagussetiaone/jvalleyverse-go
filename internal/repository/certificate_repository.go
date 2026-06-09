package repository

import (
	"context"
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

// CertificateRepository handles certificate data access
type CertificateRepository struct {
	db *gorm.DB
}

func NewCertificateRepository(db *gorm.DB) *CertificateRepository {
	return &CertificateRepository{db: db}
}

// Create creates new certificate
func (r *CertificateRepository) Create(ctx context.Context, cert *domain.Certificate) error {
	return r.db.WithContext(ctx).Create(cert).Error
}

// FindByID finds certificate by CUID
func (r *CertificateRepository) FindByID(ctx context.Context, certID string) (*domain.Certificate, error) {
	cert := &domain.Certificate{}
	if err := r.db.WithContext(ctx).Where("id = ?", certID).Preload("User").Preload("Class").First(cert).Error; err != nil {
		return nil, err
	}
	return cert, nil
}

// FindByCode finds certificate by unique code
func (r *CertificateRepository) FindByCode(ctx context.Context, code string) (*domain.Certificate, error) {
	cert := &domain.Certificate{}
	if err := r.db.WithContext(ctx).Where("unique_code = ?", code).Preload("User").Preload("Class").First(cert).Error; err != nil {
		return nil, err
	}
	return cert, nil
}

// ListByUserID lists all certificates of user
func (r *CertificateRepository) ListByUserID(ctx context.Context, userID string, page, limit int) ([]domain.Certificate, int64, error) {
	var certs []domain.Certificate
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Certificate{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Preload("Class").Offset(offset).Limit(limit).Find(&certs).Error; err != nil {
		return nil, 0, err
	}

	return certs, total, nil
}

// ListByClassID lists certificates for a class
func (r *CertificateRepository) ListByClassID(ctx context.Context, classID string, page, limit int) ([]domain.Certificate, int64, error) {
	var certs []domain.Certificate
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Certificate{}).Where("class_id = ?", classID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).Where("class_id = ?", classID).Preload("User").Offset(offset).Limit(limit).Find(&certs).Error; err != nil {
		return nil, 0, err
	}

	return certs, total, nil
}
