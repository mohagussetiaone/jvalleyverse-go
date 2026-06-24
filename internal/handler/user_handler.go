package handler

import (
	"jvalleyverse/internal/service"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userSvc       service.IUserService
	dashboardSvc  service.IDashboardService
}

func NewUserHandler(userSvc service.IUserService, dashboardSvc service.IDashboardService) *UserHandler {
	return &UserHandler{
		userSvc:      userSvc,
		dashboardSvc: dashboardSvc,
	}
}

// GetProfile returns current user profile
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	user, err := h.userSvc.GetProfile(c.UserContext(), userID)
	if err != nil {
		return safeError(c, mapServiceErrorToStatus(err), err)
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
		return safeError(c, mapServiceErrorToStatus(err), err)
	}

	return c.JSON(fiber.Map{"message": "Profile updated"})
}

// GetDashboard returns dashboard widgets and stats
func (h *UserHandler) GetDashboard(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	dashboard, err := h.dashboardSvc.GetDashboard(c.UserContext(), userID)
	if err != nil {
		return safeError(c, 500, err)
	}

	return c.JSON(dashboard)
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

	logs, total, err := h.userSvc.GetUserActivityLog(c.UserContext(), userID, page, limit)
	if err != nil {
		return safeError(c, mapServiceErrorToStatus(err), err)
	}

	return c.JSON(fiber.Map{
		"data": logs,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// ListMentors returns paginated list of mentors
func (h *UserHandler) ListMentors(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	mentors, total, err := h.userSvc.ListMentors(c.UserContext(), page, limit)
	if err != nil {
		return safeError(c, 500, err)
	}

	return c.JSON(fiber.Map{
		"data": mentors,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetAllUsers returns paginated list of all users (admin only)
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	users, total, err := h.userSvc.ListAllUsers(c.UserContext(), page, limit)
	if err != nil {
		return safeError(c, 500, err)
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
