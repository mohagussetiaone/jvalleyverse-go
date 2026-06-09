package handler

import (
	"jvalleyverse/internal/service"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userSvc service.IUserService
}

func NewUserHandler(userSvc service.IUserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

// GetProfile returns current user profile
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	user, err := h.userSvc.GetProfile(c.UserContext(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}

// UpdateProfile updates current user profile
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var input struct {
		Name   string `json:"name"`
		Bio    string `json:"bio"`
		Avatar string `json:"avatar"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.userSvc.UpdateProfile(c.UserContext(), userID, input.Name, input.Bio, input.Avatar); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Profile updated"})
}

// GetPublicProfile returns public user profile
func (h *UserHandler) GetPublicProfile(c *fiber.Ctx) error {
	userID := c.Params("id") // Direct string CUID
	if userID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	user, err := h.userSvc.GetProfile(c.UserContext(), userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Filter public info only or use a DTO
	return c.JSON(fiber.Map{
		"id":     user.ID,
		"name":   user.Name,
		"avatar": user.Avatar,
		"level":  user.Level,
		"points": user.Points,
	})
}

// GetActivityLog returns user activity log
func (h *UserHandler) GetActivityLog(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	logs, err := h.userSvc.GetUserActivityLog(c.UserContext(), userID, page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": logs,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": len(logs), 
		},
	})
}

// GetAllUsers returns paginated list of all users (admin only)
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	users, total, err := h.userSvc.ListAllUsers(c.UserContext(), page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": users,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}
