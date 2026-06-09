package handler

import (
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ShowcaseHandler struct {
	showcaseSvc service.IShowcaseService
}

func NewShowcaseHandler(showcaseSvc service.IShowcaseService) *ShowcaseHandler {
	return &ShowcaseHandler{showcaseSvc: showcaseSvc}
}

// Create creates a new showcase (POST /api/showcases)
func (h *ShowcaseHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var input struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		MediaURLs   []string `json:"media_urls"`
		CategoryID  string   `json:"category_id"`
		Visibility  string   `json:"visibility"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if input.Title == "" {
		return c.Status(400).JSON(fiber.Map{"error": "title is required"})
	}
	if input.CategoryID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "category_id is required — use GET /api/categories to get valid IDs"})
	}

	showcase, err := h.showcaseSvc.CreateShowcase(c.UserContext(), userID, input.Title, input.Description, input.MediaURLs, input.CategoryID, input.Visibility)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(showcase)
}

// ListShowcases returns paginated public showcases (GET /api/showcases)
func (h *ShowcaseHandler) ListShowcases(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	categoryID := c.Query("category_id")
	sort := c.Query("sort", "newest")

	showcases, total, err := h.showcaseSvc.ListShowcases(c.UserContext(), page, limit, categoryID, sort)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": showcases,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetShowcase returns a single showcase by ID (GET /api/showcases/:id)
func (h *ShowcaseHandler) GetShowcase(c *fiber.Ctx) error {
	showcaseID := c.Params("id")
	if showcaseID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid showcase ID"})
	}

	showcase, err := h.showcaseSvc.GetShowcaseByID(c.UserContext(), showcaseID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Showcase not found"})
	}
	return c.JSON(showcase)
}

// Update updates a showcase (PUT /api/showcases/:id)
func (h *ShowcaseHandler) Update(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	showcaseID := c.Params("id")

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Visibility  string `json:"visibility"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	showcase, err := h.showcaseSvc.UpdateShowcase(c.UserContext(), showcaseID, userID, input.Title, input.Description, input.Visibility)
	if err != nil {
		if err == domain.ErrForbidden {
			return c.Status(403).JSON(fiber.Map{"error": "You do not own this showcase"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(showcase)
}

// Delete deletes a showcase (DELETE /api/showcases/:id)
func (h *ShowcaseHandler) Delete(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	showcaseID := c.Params("id")

	if err := h.showcaseSvc.DeleteShowcase(c.UserContext(), showcaseID, userID); err != nil {
		if err == domain.ErrForbidden {
			return c.Status(403).JSON(fiber.Map{"error": "You do not own this showcase"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Showcase deleted"})
}

// Like likes a showcase (POST /api/showcases/:id/like)
func (h *ShowcaseHandler) Like(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	showcaseID := c.Params("id")
	if showcaseID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid showcase id"})
	}
	if err := h.showcaseSvc.LikeShowcase(c.UserContext(), userID, showcaseID); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Liked successfully"})
}

// Unlike removes a like (DELETE /api/showcases/:id/like)
func (h *ShowcaseHandler) Unlike(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	showcaseID := c.Params("id")
	if showcaseID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid showcase id"})
	}
	if err := h.showcaseSvc.UnlikeShowcase(c.UserContext(), userID, showcaseID); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Showcase unliked"})
}

// GetLeaderboard is kept for backward compat but real leaderboard is in GamificationHandler
func (h *ShowcaseHandler) GetLeaderboard(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"data": []fiber.Map{}})
}
