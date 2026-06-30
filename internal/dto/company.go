package dto

import (
	"time"

	"jvalleyverse/internal/domain"
)

// CompanyItem is the response DTO for company profile
type CompanyItem struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	BrandName string `json:"brand_name"`
	Tagline   string `json:"tagline"`
	Vision    string `json:"vision"`
	Mission   string `json:"mission"`
	LogoURL   string `json:"logo_url"`

	Domain    string `json:"domain"`
	Email     string `json:"email"`
	Facebook  string `json:"facebook"`
	Instagram string `json:"instagram"`
	Twitter   string `json:"twitter"`
	TikTok    string `json:"tiktok"`
	Youtube   string `json:"youtube"`
	LinkedIn  string `json:"linkedin"`
	WhatsApp  string `json:"whatsapp"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
}

// ToCompanyItem converts a domain Company to a DTO CompanyItem
func ToCompanyItem(c domain.Company) CompanyItem {
	return CompanyItem{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,

		BrandName: c.BrandName,
		Tagline:   c.Tagline,
		Vision:    c.Vision,
		Mission:   c.Mission,
		LogoURL:   c.LogoURL,

		Domain:    c.Domain,
		Email:     c.Email,
		Facebook:  c.Facebook,
		Instagram: c.Instagram,
		Twitter:   c.Twitter,
		TikTok:    c.TikTok,
		Youtube:   c.Youtube,
		LinkedIn:  c.LinkedIn,
		WhatsApp:  c.WhatsApp,
		Address:   c.Address,
		Phone:     c.Phone,
	}
}
