package service

import (
	"context"
	"fmt"

	"github.com/lucsky/cuid"

	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/dto"
	"jvalleyverse/internal/repository"
)

type ICompanyService interface {
	GetCompany(ctx context.Context) (*dto.CompanyItem, error)
	UpdateCompany(ctx context.Context, input domain.Company) (*dto.CompanyItem, error)
}

type CompanyService struct {
	companyRepo *repository.CompanyRepository
}

func NewCompanyService(companyRepo *repository.CompanyRepository) ICompanyService {
	return &CompanyService{companyRepo: companyRepo}
}

func (s *CompanyService) GetCompany(ctx context.Context) (*dto.CompanyItem, error) {
	company, err := s.companyRepo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("get company: %w", err)
	}
	item := dto.ToCompanyItem(*company)
	return &item, nil
}

func (s *CompanyService) UpdateCompany(ctx context.Context, input domain.Company) (*dto.CompanyItem, error) {
	company, err := s.companyRepo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("get company: %w", err)
	}

	// If no existing record, set a new ID
	if company.ID == "" {
		company.ID = cuid.New()
	}

	// Update only non-empty fields
	if input.BrandName != "" {
		company.BrandName = input.BrandName
	}
	if input.Tagline != "" {
		company.Tagline = input.Tagline
	}
	if input.Vision != "" {
		company.Vision = input.Vision
	}
	if input.Mission != "" {
		company.Mission = input.Mission
	}
	if input.LogoURL != "" {
		company.LogoURL = input.LogoURL
	}
	if input.Domain != "" {
		company.Domain = input.Domain
	}
	if input.Email != "" {
		company.Email = input.Email
	}
	if input.Facebook != "" {
		company.Facebook = input.Facebook
	}
	if input.Instagram != "" {
		company.Instagram = input.Instagram
	}
	if input.Twitter != "" {
		company.Twitter = input.Twitter
	}
	if input.TikTok != "" {
		company.TikTok = input.TikTok
	}
	if input.Youtube != "" {
		company.Youtube = input.Youtube
	}
	if input.LinkedIn != "" {
		company.LinkedIn = input.LinkedIn
	}
	if input.WhatsApp != "" {
		company.WhatsApp = input.WhatsApp
	}
	if input.Address != "" {
		company.Address = input.Address
	}
	if input.Phone != "" {
		company.Phone = input.Phone
	}

	if err := s.companyRepo.Upsert(ctx, company); err != nil {
		return nil, fmt.Errorf("update company: %w", err)
	}

	item := dto.ToCompanyItem(*company)
	return &item, nil
}
