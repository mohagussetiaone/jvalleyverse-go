package dto

import (
	"jvalleyverse/internal/domain"
)

type SectionBrief struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description,omitempty"`
	OrderIndex  int           `json:"order_index"`
	Lessons     []LessonBrief `json:"lessons,omitempty"`
}

type SectionDetail struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description,omitempty"`
	CourseID    string        `json:"course_id"`
	OrderIndex  int           `json:"order_index"`
	Lessons     []LessonBrief `json:"lessons,omitempty"`
}

func ToSectionBrief(s domain.Section) SectionBrief {
	lessons := make([]LessonBrief, len(s.Lessons))
	for i, l := range s.Lessons {
		lessons[i] = ToLessonBrief(l)
	}
	return SectionBrief{
		ID:          s.ID,
		Title:       s.Title,
		Description: s.Description,
		OrderIndex:  s.OrderIndex,
		Lessons:     lessons,
	}
}

func ToSectionBriefPtr(s *domain.Section) *SectionBrief {
	if s == nil {
		return nil
	}
	b := ToSectionBrief(*s)
	return &b
}

func ToSectionDetail(s *domain.Section) *SectionDetail {
	if s == nil {
		return nil
	}
	lessons := make([]LessonBrief, len(s.Lessons))
	for i, l := range s.Lessons {
		lessons[i] = ToLessonBrief(l)
	}
	return &SectionDetail{
		ID:          s.ID,
		Title:       s.Title,
		Description: s.Description,
		CourseID:    s.CourseID,
		OrderIndex:  s.OrderIndex,
		Lessons:     lessons,
	}
}
