package dto

import (
	"time"

	"jvalleyverse/internal/domain"
)

type CertificateItem struct {
	ID          string           `json:"id"`
	UniqueCode  string           `json:"unique_code"`
	IssuedAt    time.Time        `json:"issued_at"`
	UserID      string           `json:"user_id,omitempty"`
	LessonID    string           `json:"lesson_id,omitempty"`
	LessonName  string           `json:"lesson_name,omitempty"`
	UserName    string           `json:"user_name,omitempty"`
	Achievement *AchievementInfo `json:"achievement,omitempty"`
}

type AchievementInfo struct {
	Type       string `json:"type"`
	Title      string `json:"title"`
	UniqueCode string `json:"unique_code"`
}

func ToCertificateItem(cert *domain.Certificate) *CertificateItem {
	if cert == nil {
		return nil
	}
	return &CertificateItem{
		ID:          cert.ID,
		UniqueCode:  cert.UniqueCode,
		IssuedAt:    cert.IssuedAt,
		UserID:      cert.UserID,
		LessonID:    cert.LessonID,
		LessonName:  cert.Lesson.Title,
		UserName:    cert.User.Name,
	}
}
