package dto

import "time"

type ActivityHistoryItem struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Link        string    `json:"link,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}
