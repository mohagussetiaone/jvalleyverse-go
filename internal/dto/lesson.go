package dto

import (
	"jvalleyverse/internal/domain"
)

type LessonBrief struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Slug       string `json:"slug"`
	Difficulty string `json:"difficulty,omitempty"`
	Duration   int    `json:"duration"`
	OrderIndex int    `json:"order_index"`
	VideoURL   string `json:"video_url,omitempty"`
}

type LessonDetailResponse struct {
	Lesson     LessonBrief              `json:"lesson"`
	Details    *domain.LessonDetail      `json:"details,omitempty"`
	Progress   *domain.LessonProgress    `json:"progress,omitempty"`
	NextLesson *LessonBrief              `json:"next_lesson,omitempty"`
	Section    *SectionBrief             `json:"section,omitempty"`
	Course     *CourseListItem           `json:"course,omitempty"`
}

func ToLessonBrief(l domain.Lesson) LessonBrief {
	return LessonBrief{
		ID:         l.ID,
		Title:      l.Title,
		Slug:       l.Slug,
		Difficulty: l.Difficulty,
		Duration:   l.Duration,
		OrderIndex: l.OrderIndex,
		VideoURL:   l.VideoURL,
	}
}

func ToLessonBriefPtr(l *domain.Lesson) *LessonBrief {
	if l == nil {
		return nil
	}
	b := ToLessonBrief(*l)
	return &b
}
