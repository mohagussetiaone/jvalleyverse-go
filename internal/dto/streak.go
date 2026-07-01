package dto

import (
	"time"
)

type StreakItem struct {
	StreakCount      int       `json:"streak_count"`
	LongestStreak    int       `json:"longest_streak"`
	LastActivityDate time.Time `json:"last_activity_date"`
}
