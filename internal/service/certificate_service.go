package service

import (
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/repository"
)

type CertificateService struct {
	certRepo  *repository.CertificateRepository
	classRepo *repository.ClassRepository
	userRepo  *repository.UserRepository
}

func NewCertificateService() *CertificateService {
	return &CertificateService{
		certRepo:  repository.NewCertificateRepository(),
		classRepo: repository.NewClassRepository(),
		userRepo:  repository.NewUserRepository(),
	}
}

// IssueCertificate issues new certificate for class completion
func (s *CertificateService) IssueCertificate(userID, classID uint, code string) (*domain.Certificate, error) {
	cert := &domain.Certificate{
		UserID:     userID,
		ClassID:    classID,
		UniqueCode: code,
	}

	if err := s.certRepo.Create(cert); err != nil {
		return nil, err
	}

	return cert, nil
}

// GetCertificate retrieves certificate details
func (s *CertificateService) GetCertificate(certID uint) (map[string]interface{}, error) {
	cert, err := s.certRepo.FindByID(certID)
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
func (s *CertificateService) GetCertificateByCode(code string) (map[string]interface{}, error) {
	cert, err := s.certRepo.FindByCode(code)
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

// ListUserCertificates returns certificates earned by user
func (s *CertificateService) ListUserCertificates(userID uint, page, limit int) ([]map[string]interface{}, error) {
	certs, _, err := s.certRepo.ListByUserID(userID, page, limit)
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
		}
	}

	return result, nil
}
