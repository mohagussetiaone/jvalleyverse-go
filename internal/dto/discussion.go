package dto

import (
	"time"

	"jvalleyverse/internal/domain"
)

type DiscussionBrief struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	User       UserBrief `json:"user"`
	Status     string    `json:"status"`
	ViewsCount int       `json:"views_count"`
	CreatedAt  time.Time `json:"created_at"`
}

type DiscussionListItem struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content,omitempty"`
	User        UserBrief `json:"user"`
	LessonID    *string   `json:"lesson_id,omitempty"`
	StudyCaseID *string   `json:"study_case_id,omitempty"`
	Status      string    `json:"status"`
	ViewsCount  int       `json:"view_count"`
	RepliesCount int       `json:"replies_count"`
	CreatedAt   time.Time `json:"created_at"`
}

type DiscussionDetail struct {
	ID         string              `json:"id"`
	Title      string              `json:"title"`
	Content    string              `json:"content"`
	User       UserBrief           `json:"user"`
	Status     string              `json:"status"`
	ViewsCount int                 `json:"view_count"`
	CreatedAt  time.Time           `json:"created_at"`
	Replies    []ReplyInDiscussion `json:"replies"`
}

type ReplyInDiscussion struct {
	ID        string              `json:"id"`
	Content   string              `json:"content"`
	User      UserBrief           `json:"user"`
	Likes     int                 `json:"likes"`
	IsBest    bool                `json:"is_best"`
	ParentID  *string             `json:"parent_id,omitempty"`
	CreatedAt time.Time           `json:"created_at"`
	Children  []ReplyInDiscussion `json:"children,omitempty"`
	Reactions []ReactionSummary   `json:"reactions,omitempty"`
}

// ToDiscussionBrief converts a domain Discussion to DiscussionBrief
func ToDiscussionBrief(d domain.Discussion) DiscussionBrief {
	return DiscussionBrief{
		ID:         d.ID,
		Title:      d.Title,
		User:       ToUserBrief(d.User),
		Status:     d.Status,
		ViewsCount: d.ViewsCount,
		CreatedAt:  d.CreatedAt,
	}
}

// ToDiscussionListItem converts a domain Discussion to DiscussionListItem
func ToDiscussionListItem(d domain.Discussion) DiscussionListItem {
	return DiscussionListItem{
		ID:           d.ID,
		Title:        d.Title,
		Content:      d.Content,
		User:         ToUserBrief(d.User),
		LessonID:     d.LessonID,
		StudyCaseID:  d.StudyCaseID,
		Status:       d.Status,
		ViewsCount:   d.ViewsCount,
		RepliesCount: len(d.Replies),
		CreatedAt:    d.CreatedAt,
	}
}

// ToReplyInDiscussion converts a domain Reply to ReplyInDiscussion
func ToReplyInDiscussion(r domain.Reply) ReplyInDiscussion {
	return ReplyInDiscussion{
		ID:        r.ID,
		Content:   r.Content,
		User:      ToUserBrief(r.User),
		Likes:     r.LikesCount,
		IsBest:    r.IsMarkedBest,
		ParentID:  r.ParentID,
		CreatedAt: r.CreatedAt,
	}
}
