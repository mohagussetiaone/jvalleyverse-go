package dto

import "time"

type PortfolioItem struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Type        string   `json:"type"` // "certificate", "showcase", "study_case"
	URL         string   `json:"url,omitempty"`
	ImageURL    string   `json:"image_url,omitempty"`
	Technologies []string `json:"technologies,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type PortfolioResponse struct {
	User          UserBrief       `json:"user"`
	TotalPoints   int             `json:"total_points"`
	Level         int             `json:"level"`
	Items         []PortfolioItem `json:"items"`
	CertCount     int             `json:"cert_count"`
	ShowcaseCount int             `json:"showcase_count"`
	StudyCaseCount int            `json:"study_case_count"`
}
