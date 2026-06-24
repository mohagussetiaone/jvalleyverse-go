package dto

import (
	"time"
)

type ReplyListItem struct {
	ID              string    `json:"id"`
	Content         string    `json:"content"`
	DiscussionID    string    `json:"discussion_id"`
	DiscussionTitle string    `json:"discussion_title"`
	ParentID        *string   `json:"parent_id,omitempty"`
	LikesCount      int       `json:"likes_count"`
	IsMarkedBest    bool      `json:"is_marked_best"`
	CreatedAt       time.Time `json:"created_at"`
}
