package dto

import (
	"encoding/json"
	"time"

	"jvalleyverse/internal/domain"
)

type ShowcaseListItem struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description,omitempty"`
	MediaURLs   []string     `json:"media_urls,omitempty"`
	LikesCount  int          `json:"likes_count"`
	ViewsCount  int          `json:"views_count"`
	Visibility  string       `json:"visibility,omitempty"`
	User        UserBrief    `json:"user"`
	Category    CategoryBrief `json:"category"`
	CreatedAt   time.Time    `json:"created_at"`
}

type ShowcaseDetail struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description,omitempty"`
	MediaURLs   []string     `json:"media_urls,omitempty"`
	LikesCount  int          `json:"likes_count"`
	ViewsCount  int          `json:"views_count"`
	User        UserBrief    `json:"user"`
	Category    CategoryBrief `json:"category"`
	IsLikedByMe bool         `json:"is_liked_by_me,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
}

func ToShowcaseListItem(sc domain.Showcase) ShowcaseListItem {
	mediaURLs := make([]string, 0)
	if sc.MediaURLs != nil {
		json.Unmarshal(sc.MediaURLs, &mediaURLs)
	}
	return ShowcaseListItem{
		ID:          sc.ID,
		Title:       sc.Title,
		Description: sc.Description,
		MediaURLs:   mediaURLs,
		LikesCount:  sc.LikesCount,
		ViewsCount:  sc.ViewsCount,
		Visibility:  sc.Visibility,
		User:        ToUserBrief(sc.User),
		Category:    ToCategoryBrief(sc.Category),
		CreatedAt:   sc.CreatedAt,
	}
}

func ToShowcaseDetail(sc *domain.Showcase) *ShowcaseDetail {
	if sc == nil {
		return nil
	}
	mediaURLs := make([]string, 0)
	if sc.MediaURLs != nil {
		json.Unmarshal(sc.MediaURLs, &mediaURLs)
	}
	return &ShowcaseDetail{
		ID:          sc.ID,
		Title:       sc.Title,
		Description: sc.Description,
		MediaURLs:   mediaURLs,
		LikesCount:  sc.LikesCount,
		ViewsCount:  sc.ViewsCount,
		User:        ToUserBrief(sc.User),
		Category:    ToCategoryBrief(sc.Category),
		CreatedAt:   sc.CreatedAt,
	}
}
