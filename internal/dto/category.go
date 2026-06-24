package dto

import (
	"jvalleyverse/internal/domain"
)

// ──────────────────────────────────────────────
//  BRIEF DTOs
// ──────────────────────────────────────────────

type CategoryBrief struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type CategoryWithCourses struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Slug        string           `json:"slug"`
	Description string           `json:"description,omitempty"`
	Courses     []CourseListItem `json:"courses,omitempty"`
}

// ──────────────────────────────────────────────
//  CONVERTERS
// ──────────────────────────────────────────────

func ToCategoryBrief(c domain.Category) CategoryBrief {
	return CategoryBrief{
		ID:   c.ID,
		Name: c.Name,
		Slug: c.Slug,
	}
}
