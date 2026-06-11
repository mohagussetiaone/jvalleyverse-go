package service

import (
	"context"
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/repository"
)

// ICertificateService defines the business logic for managing certificates
type ICertificateService interface {
	IssueCertificate(ctx context.Context, userID, classID string, code string) (*domain.Certificate, error)
	GetCertificate(ctx context.Context, certID string) (map[string]interface{}, error)
	GetCertificateByCode(ctx context.Context, code, requesterID, requesterRole string) (map[string]interface{}, error)
	ListUserCertificates(ctx context.Context, userID string, page, limit int) ([]map[string]interface{}, error)
}

type CertificateService struct {
	certRepo  *repository.CertificateRepository
	classRepo *repository.ClassRepository
	userRepo  *repository.UserRepository
}

func NewCertificateService(
	certRepo *repository.CertificateRepository,
	classRepo *repository.ClassRepository,
	userRepo *repository.UserRepository,
) *CertificateService {
	return &CertificateService{
		certRepo:  certRepo,
		classRepo: classRepo,
		userRepo:  userRepo,
	}
}

// IssueCertificate issues new certificate for class completion
func (s *CertificateService) IssueCertificate(ctx context.Context, userID, classID string, code string) (*domain.Certificate, error) {
	cert := &domain.Certificate{
		UserID:     userID,
		ClassID:    classID,
		UniqueCode: code,
	}

	if err := s.certRepo.Create(ctx, cert); err != nil {
		return nil, err
	}

	return cert, nil
}

// GetCertificate retrieves certificate details
func (s *CertificateService) GetCertificate(ctx context.Context, certID string) (map[string]interface{}, error) {
	cert, err := s.certRepo.FindByID(ctx, certID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":          cert.ID,
		"unique_code": cert.UniqueCode,
		"issued_at":   cert.IssuedAt,
		"user_id":     cert.UserID,
		"class_id":    cert.ClassID,
		"class_name":  cert.Class.Title,
		"user_name":   cert.User.Name,
	}, nil
}

// GetCertificateByCode retrieves certificate by unique code
func (s *CertificateService) GetCertificateByCode(ctx context.Context, code, requesterID, requesterRole string) (map[string]interface{}, error) {
	cert, err := s.certRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if cert.UserID != requesterID && requesterRole != "admin" {
		return nil, domain.ErrForbidden
	}

	return map[string]interface{}{
		"id":          cert.ID,
		"unique_code": cert.UniqueCode,
		"issued_at":   cert.IssuedAt,
		"user_id":     cert.UserID,
		"class_id":    cert.ClassID,
		"class_name":  cert.Class.Title,
		"user_name":   cert.User.Name,
		"achievement": map[string]interface{}{
			"type":        "certificate",
			"title":       cert.Class.Title,
			"unique_code": cert.UniqueCode,
		},
	}, nil
}

// ListUserCertificates returns certificates earned by user
func (s *CertificateService) ListUserCertificates(ctx context.Context, userID string, page, limit int) ([]map[string]interface{}, error) {
	certs, _, err := s.certRepo.ListByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(certs))
	for i, cert := range certs {
		result[i] = map[string]interface{}{
			"id":          cert.ID,
			"unique_code": cert.UniqueCode,
			"issued_at":   cert.IssuedAt,
			"class_id":    cert.ClassID,
			"class_name":  cert.Class.Title,
			"achievement": "certificate",
		}
	}

	return result, nil
}
