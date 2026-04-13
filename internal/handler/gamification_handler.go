package handler

import (
	"jvalleyverse/internal/service"

	"github.com/gofiber/fiber/v2"
)

type GamificationHandler struct {
	gamificationSvc *service.GamificationService
}

func NewGamificationHandler() *GamificationHandler {
	return &GamificationHandler{gamificationSvc: service.NewGamificationService()}
}

// GetLevels returns all level information
func (h *GamificationHandler) GetLevels(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"data": []fiber.Map{
			{
				"level":       1,
				"badge_name": "Beginner",
				"min_points":  0,
				"max_points":  199,
				"badge_icon": "",
				"description": "Welcome to the community!",
			},
			{
				"level":       2,
				"badge_name": "Learner",
				"min_points":  200,
				"max_points":  499,
				"badge_icon": "",
				"description": "You're making progress!",
			},
			{
				"level":       3,
				"badge_name": "Contributor",
				"min_points":  500,
				"max_points":  999,
				"badge_icon": "",
				"description": "Great contributions!",
			},
			{
				"level":       4,
				"badge_name": "Expert",
				"min_points":  1000,
				"max_points":  1999,
				"badge_icon": "",
				"description": "You're an expert!",
			},
			{
				"level":       5,
				"badge_name": "Master",
				"min_points":  2000,
				"max_points":  nil,
				"badge_icon": "",
				"description": "Master of the platform!",
			},
		},
	})
}

// GetUserPoints returns user points and level
func (h *GamificationHandler) GetUserPoints(c *fiber.Ctx) error {
	userID := c.Params("id")

	return c.JSON(fiber.Map{
		"user_id": userID,
		"name":    "",
		"points":  0,
		"level":   1,
		"rank":    0,
		"level_info": fiber.Map{
			"level":       1,
			"badge_name": "Beginner",
			"min_points":  0,
			"max_points":  199,
			"progress":    0,
		},
	})
}
