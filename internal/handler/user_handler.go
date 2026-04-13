package handler

import (
	"jvalleyverse/internal/service"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userSvc *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{userSvc: service.NewUserService()}
}

// GetProfile returns current user profile
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	return c.JSON(fiber.Map{
		"id":       userID,
		"email":    "",
		"name":     "",
		"role":     "user",
		"points":   0,
		"level":    1,
		"avatar":   "",
		"bio":      "",
		"is_active": true,
	})
}

// UpdateProfile updates current user profile
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var input struct {
		Name   string `json:"name"`
		Bio    string `json:"bio"`
		Avatar string `json:"avatar"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// TODO: Call service to update user
	return c.JSON(fiber.Map{
		"message": "Profile updated",
		"id":      userID,
	})
}

// GetPublicProfile returns public user profile
func (h *UserHandler) GetPublicProfile(c *fiber.Ctx) error {
	userID := c.Params("id")

	return c.JSON(fiber.Map{
		"id":                userID,
		"name":              "",
		"avatar":            "",
		"level":             1,
		"points":            0,
		"showcase_count":    0,
		"certificate_count": 0,
	})
}

// GetActivityLog returns user activity log
func (h *UserHandler) GetActivityLog(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"data": []fiber.Map{},
		"pagination": fiber.Map{
			"page":  1,
			"limit": 50,
			"total": 0,
		},
	})
}
