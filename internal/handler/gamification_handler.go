package handler

import (
	"jvalleyverse/internal/service"

	"github.com/gofiber/fiber/v2"
)

type GamificationHandler struct {
	gamificationSvc service.IGamificationService
}

func NewGamificationHandler() *GamificationHandler {
	return &GamificationHandler{gamificationSvc: service.GetGamificationService()}
}

// GetLevels returns all level information
func (h *GamificationHandler) GetLevels(c *fiber.Ctx) error {
	levels := h.gamificationSvc.GetLevelInfo()
	return c.JSON(fiber.Map{"data": levels})
}

// GetUserPoints returns user points and level
func (h *GamificationHandler) GetUserPoints(c *fiber.Ctx) error {
	userID := c.Params("id") // String CUID
	if userID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	stats, err := h.gamificationSvc.GetUserStats(c.UserContext(), userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(stats)
}

// GetLeaderboard returns top users by points
func (h *GamificationHandler) GetLeaderboard(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)

	leaderboard, err := h.gamificationSvc.GetLeaderboard(c.UserContext(), limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": leaderboard})
}
