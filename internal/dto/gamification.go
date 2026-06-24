package dto

import (
	"time"
)

type LeaderboardItem struct {
	Rank        int    `json:"rank"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar,omitempty"`
	TotalPoints int    `json:"total_points"`
	Level       int    `json:"level"`
}

type ActivityItem struct {
	ID        string    `json:"id"`
	Activity  string    `json:"activity"`
	Points    int       `json:"points"`
	Timestamp time.Time `json:"timestamp"`
}

type LevelInfo struct {
	Name        string `json:"name"`
	Threshold   int    `json:"threshold"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

type UserStats struct {
	UserID         string         `json:"user_id"`
	Name           string         `json:"name"`
	TotalPoints    int            `json:"total_points"`
	CurrentLevel   int            `json:"current_level"`
	RecentActivity []ActivityItem `json:"recent_activity,omitempty"`
}
