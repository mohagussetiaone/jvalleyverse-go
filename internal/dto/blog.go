package dto

import (
	"encoding/json"
	"time"

	"jvalleyverse/internal/domain"
)

type BlogListItem struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Slug        string         `json:"slug"`
	Description string         `json:"description"`
	CoverImgURL string         `json:"cover_img_url"`
	Tags        []string       `json:"tags"`
	Status      string         `json:"status"`
	Author      UserBrief      `json:"author"`
	Category    CategoryBrief  `json:"category"`
	CreatedAt   time.Time      `json:"created_at"`
}

type BlogDetail struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Slug        string         `json:"slug"`
	Description string         `json:"description"`
	Content     string         `json:"content"`
	CoverImgURL string         `json:"cover_img_url"`
	Tags        []string       `json:"tags"`
	Status      string         `json:"status"`
	Author      UserBrief      `json:"author"`
	Category    CategoryBrief  `json:"category"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func ToBlogListItem(b domain.Blog) BlogListItem {
	tags := []string{}
	if b.Tags != nil {
		json.Unmarshal(b.Tags, &tags)
	}

	return BlogListItem{
		ID:          b.ID,
		Title:       b.Title,
		Slug:        b.Slug,
		Description: b.Description,
		CoverImgURL: b.CoverImgURL,
		Tags:        tags,
		Status:      b.Status,
		Author: UserBrief{
			ID:     b.User.ID,
			Name:   b.User.Name,
			Avatar: b.User.Avatar,
		},
		Category: CategoryBrief{
			ID:   b.Category.ID,
			Name: b.Category.Name,
			Slug: b.Category.Slug,
		},
		CreatedAt: b.CreatedAt,
	}
}

func ToBlogDetail(b *domain.Blog) *BlogDetail {
	tags := []string{}
	if b.Tags != nil {
		json.Unmarshal(b.Tags, &tags)
	}

	return &BlogDetail{
		ID:          b.ID,
		Title:       b.Title,
		Slug:        b.Slug,
		Description: b.Description,
		Content:     b.Content,
		CoverImgURL: b.CoverImgURL,
		Tags:        tags,
		Status:      b.Status,
		Author: UserBrief{
			ID:     b.User.ID,
			Name:   b.User.Name,
			Avatar: b.User.Avatar,
		},
		Category: CategoryBrief{
			ID:   b.Category.ID,
			Name: b.Category.Name,
			Slug: b.Category.Slug,
		},
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}
