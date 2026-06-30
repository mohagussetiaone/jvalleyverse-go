package dto

import (
	"time"

	"jvalleyverse/internal/domain"
)

// FAQItem is the response DTO for a single FAQ
type FAQItem struct {
	ID         string    `json:"id"`
	Question   string    `json:"question"`
	Answer     string    `json:"answer"`
	Category   string    `json:"category"`
	OrderIndex int       `json:"order_index"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ToFAQItem converts a domain FAQ to a DTO FAQItem
func ToFAQItem(faq domain.FAQ) FAQItem {
	return FAQItem{
		ID:         faq.ID,
		Question:   faq.Question,
		Answer:     faq.Answer,
		Category:   faq.Category,
		OrderIndex: faq.OrderIndex,
		IsActive:   faq.IsActive,
		CreatedAt:  faq.CreatedAt,
		UpdatedAt:  faq.UpdatedAt,
	}
}
