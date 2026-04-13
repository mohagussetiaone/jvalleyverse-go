package repository

import (
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

// CertificateRepository handles certificate data access
type CertificateRepository struct {
	db *gorm.DB
}

func NewCertificateRepository() *CertificateRepository {
	return &CertificateRepository{db: db}
}

// Create creates new certificate
func (r *CertificateRepository) Create(cert *domain.Certificate) error {
	return r.db.Create(cert).Error
}

// FindByID finds certificate by ID
func (r *CertificateRepository) FindByID(certID uint) (*domain.Certificate, error) {
	cert := &domain.Certificate{}
	if err := r.db.Preload("User").Preload("Class").First(cert, certID).Error; err != nil {
		return nil, err
	}
	return cert, nil
}

// FindByCode finds certificate by unique code
func (r *CertificateRepository) FindByCode(code string) (*domain.Certificate, error) {
	cert := &domain.Certificate{}
	if err := r.db.Where("unique_code = ?", code).Preload("User").Preload("Class").First(cert).Error; err != nil {
		return nil, err
	}
	return cert, nil
}

// ListByUserID lists all certificates of user
func (r *CertificateRepository) ListByUserID(userID uint, page, limit int) ([]domain.Certificate, int64, error) {
	var certs []domain.Certificate
	var total int64

	if err := r.db.Model(&domain.Certificate{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.Where("user_id = ?", userID).Preload("Class").Offset(offset).Limit(limit).Find(&certs).Error; err != nil {
		return nil, 0, err
	}

	return certs, total, nil
}

// ListByClassID lists certificates for a class
func (r *CertificateRepository) ListByClassID(classID uint, page, limit int) ([]domain.Certificate, int64, error) {
	var certs []domain.Certificate
	var total int64

	if err := r.db.Model(&domain.Certificate{}).Where("class_id = ?", classID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.Where("class_id = ?", classID).Preload("User").Offset(offset).Limit(limit).Find(&certs).Error; err != nil {
		return nil, 0, err
	}

	return certs, total, nil
}
