package dto

import (
	"time"

	"jvalleyverse/internal/domain"
)

type CourseListItem struct {
	ID           string         `json:"id"`
	Title        string         `json:"title"`
	Description  string         `json:"description,omitempty"`
	Thumbnail    string         `json:"thumbnail,omitempty"`
	Price        float64        `json:"price"`
	Category     CategoryBrief  `json:"category"`
	AdminName    string         `json:"admin_name"`
	Mentor       *UserBrief     `json:"mentor,omitempty"`
	Hours        int            `json:"hours"`
	Visibility   string         `json:"visibility,omitempty"`
	SectionCount int            `json:"section_count"`
	IsEnrolled   bool           `json:"is_enrolled,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
}

type CourseDetailWithSections struct {
	ID                 string         `json:"id"`
	Title              string         `json:"title"`
	Description        string         `json:"description"`
	Thumbnail          string         `json:"thumbnail,omitempty"`
	Price              float64        `json:"price"`
	Category           CategoryBrief  `json:"category"`
	AdminID            string         `json:"admin_id"`
	AdminName          string         `json:"admin_name"`
	Mentor             *UserBrief     `json:"mentor,omitempty"`
	Hours              int            `json:"hours"`
	TotalDurationHours int            `json:"total_duration_hours"`
	Visibility         string         `json:"visibility"`
	Sections           []SectionBrief `json:"sections"`
	IsEnrolled         bool           `json:"is_enrolled,omitempty"`
	CreatedAt          time.Time      `json:"created_at"`
}

type EnrolledCourseItem struct {
	CourseListItem
	EnrolledAt  time.Time `json:"enrolled_at"`
	LastLessonID *string  `json:"last_lesson_id"`
}

func CourseToListItem(c domain.Course) CourseListItem {
	return CourseListItem{
		ID:           c.ID,
		Title:        c.Title,
		Description:  c.Description,
		Thumbnail:    c.Thumbnail,
		Price:        c.Price,
		Category:     ToCategoryBrief(c.Category),
		AdminName:    c.Admin.Name,
		Mentor:       ToUserBriefPtr(&c.Mentor, c.MentorID),
		Hours:        c.Hours,
		Visibility:   c.Visibility,
		SectionCount: len(c.Sections),
		CreatedAt:    c.CreatedAt,
	}
}

func CourseToListItemPtr(c *domain.Course) *CourseListItem {
	if c == nil {
		return nil
	}
	item := CourseToListItem(*c)
	return &item
}
