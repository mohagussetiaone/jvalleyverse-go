package service

import (
	"context"
	"fmt"
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/dto"
	"jvalleyverse/internal/repository"
	"jvalleyverse/pkg/config"
)

type ICertificateService interface {
	IssueCertificate(ctx context.Context, userID, lessonID string, code string) (*domain.Certificate, error)
	GetCertificate(ctx context.Context, certID string) (*dto.CertificateItem, error)
	GetCertificateByCode(ctx context.Context, code, requesterID, requesterRole string) (*dto.CertificateItem, error)
	ListUserCertificates(ctx context.Context, userID string, page, limit int) ([]dto.CertificateItem, int64, error)
	VerifyCertificateByCode(ctx context.Context, code string) (*dto.CertificateItem, error)
}

type CertificateService struct {
	certRepo   *repository.CertificateRepository
	lessonRepo *repository.LessonRepository
	userRepo   *repository.UserRepository
}

func NewCertificateService(
	certRepo *repository.CertificateRepository,
	lessonRepo *repository.LessonRepository,
	userRepo *repository.UserRepository,
) *CertificateService {
	return &CertificateService{
		certRepo:   certRepo,
		lessonRepo: lessonRepo,
		userRepo:   userRepo,
	}
}

func (s *CertificateService) IssueCertificate(ctx context.Context, userID, lessonID string, code string) (*domain.Certificate, error) {
	// Build verification URL
	baseURL := "https://jvalleyverse.com"
	if config.AppConfig != nil && config.AppConfig.FrontendURL != "" {
		baseURL = config.AppConfig.FrontendURL
	}
	verificationURL := fmt.Sprintf("%s/certificates/verify/%s", baseURL, code)
	qrCodeURL := fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=%s", verificationURL)

	cert := &domain.Certificate{
		UserID:          userID,
		LessonID:        lessonID,
		UniqueCode:      code,
		VerificationURL: verificationURL,
		QRCodeURL:       qrCodeURL,
	}

	if err := s.certRepo.Create(ctx, cert); err != nil {
		return nil, err
	}

	return cert, nil
}

func (s *CertificateService) GetCertificate(ctx context.Context, certID string) (*dto.CertificateItem, error) {
	cert, err := s.certRepo.FindByID(ctx, certID)
	if err != nil {
		return nil, err
	}
	return dto.ToCertificateItem(cert), nil
}

func (s *CertificateService) GetCertificateByCode(ctx context.Context, code, requesterID, requesterRole string) (*dto.CertificateItem, error) {
	cert, err := s.certRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if cert.UserID != requesterID && requesterRole != "admin" {
		return nil, domain.ErrForbidden
	}

	item := dto.ToCertificateItem(cert)
	item.Achievement = &dto.AchievementInfo{
		Type:       "certificate",
		Title:      cert.Lesson.Title,
		UniqueCode: cert.UniqueCode,
	}
	item.QRCodeURL = cert.QRCodeURL
	item.VerificationURL = cert.VerificationURL
	return item, nil
}

// VerifyCertificateByCode returns certificate verification info (public, no auth required)
func (s *CertificateService) VerifyCertificateByCode(ctx context.Context, code string) (*dto.CertificateItem, error) {
	cert, err := s.certRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	item := dto.ToCertificateItem(cert)
	item.Achievement = &dto.AchievementInfo{
		Type:       "certificate",
		Title:      cert.Lesson.Title,
		UniqueCode: cert.UniqueCode,
	}
	item.QRCodeURL = cert.QRCodeURL
	item.VerificationURL = cert.VerificationURL
	return item, nil
}

func (s *CertificateService) ListUserCertificates(ctx context.Context, userID string, page, limit int) ([]dto.CertificateItem, int64, error) {
	certs, total, err := s.certRepo.ListByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.CertificateItem, len(certs))
	for i, cert := range certs {
		item := dto.ToCertificateItem(&cert)
		item.Achievement = &dto.AchievementInfo{Type: "certificate"}
		result[i] = *item
	}

	return result, total, nil
}
