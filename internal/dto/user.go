package dto

import (
	"time"

	"jvalleyverse/internal/domain"
)

type UserBrief struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar,omitempty"`
	Bio    string `json:"bio,omitempty"`
	Role   string `json:"role,omitempty"`
}

type UserListItem struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	Avatar      string    `json:"avatar,omitempty"`
	Role        string    `json:"role"`
	Level       int       `json:"level"`
	TotalPoints int       `json:"total_points"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type MentorItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar,omitempty"`
	Bio         string `json:"bio,omitempty"`
	Level       int    `json:"level"`
	TotalPoints int    `json:"total_points"`
}

func ToUserBrief(u domain.User) UserBrief {
	return UserBrief{
		ID:     u.ID,
		Name:   u.Name,
		Avatar: u.Avatar,
		Bio:    u.Bio,
		Role:   u.Role,
	}
}

func ToUserBriefPtr(u *domain.User, id string) *UserBrief {
	if u == nil || id == "" {
		return nil
	}
	b := ToUserBrief(*u)
	return &b
}
