package dto

import (
	"time"

	"jvalleyverse/internal/domain"
)

type ReviewItem struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	UserName   string    `json:"user_name"`
	UserAvatar string    `json:"user_avatar,omitempty"`
	CourseID   string    `json:"course_id,omitempty"`
	LessonID   string    `json:"lesson_id,omitempty"`
	Rating     int       `json:"rating"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
}

func ToReviewItems(reviews []domain.Review) []ReviewItem {
	result := make([]ReviewItem, len(reviews))
	for i, r := range reviews {
		result[i] = ReviewItem{
			ID:         r.ID,
			UserID:     r.UserID,
			UserName:   r.User.Name,
			UserAvatar: r.User.Avatar,
			CourseID:   r.CourseID,
			LessonID:   r.LessonID,
			Rating:     r.Rating,
			Message:    r.Message,
			CreatedAt:  r.CreatedAt,
		}
	}
	return result
}
