package dto

import (
	"encoding/json"
	"time"

	"jvalleyverse/internal/domain"
)

type StudyCaseListItem struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	ImgURL      string       `json:"img_url,omitempty"`
	YoutubeURL  string       `json:"youtube_url,omitempty"`
	Tags        []string     `json:"tags,omitempty"`
	Category    CategoryBrief `json:"category,omitempty"`
	User        UserBrief    `json:"user"`
	CreatedAt   time.Time    `json:"created_at"`
}

type StudyCaseDetail struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	ImgURL      string             `json:"img_url,omitempty"`
	YoutubeURL  string             `json:"youtube_url,omitempty"`
	Tags        []string           `json:"tags,omitempty"`
	Category    CategoryBrief      `json:"category,omitempty"`
	User        UserBrief          `json:"user"`
	CreatedAt   time.Time          `json:"created_at"`
	Discussions []DiscussionBrief  `json:"discussions,omitempty"`
}

func ToStudyCaseListItem(sc domain.StudyCase) StudyCaseListItem {
	var tags []string
	if sc.Tags != nil {
		json.Unmarshal(sc.Tags, &tags)
	}
	var category CategoryBrief
	if sc.CategoryID != nil {
		category = ToCategoryBrief(sc.Category)
	}
	return StudyCaseListItem{
		ID:          sc.ID,
		Name:        sc.Name,
		Description: sc.Description,
		ImgURL:      sc.ImgURL,
		YoutubeURL:  sc.YoutubeURL,
		Tags:        tags,
		Category:    category,
		User:        ToUserBrief(sc.User),
		CreatedAt:   sc.CreatedAt,
	}
}
