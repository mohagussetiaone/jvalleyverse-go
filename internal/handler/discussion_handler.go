package handler

import (
	"jvalleyverse/internal/service"

	"github.com/gofiber/fiber/v2"
)

type DiscussionHandler struct {
	discussionSvc *service.DiscussionService
}

func NewDiscussionHandler() *DiscussionHandler {
	return &DiscussionHandler{discussionSvc: service.NewDiscussionService()}
}

// CreateDiscussion creates new discussion
func (h *DiscussionHandler) CreateDiscussion(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var input struct {
		Title      string `json:"title"`
		Content    string `json:"content"`
		ClassID    *uint  `json:"class_id"`
		CategoryID uint   `json:"category_id"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// TODO: Call service
	return c.Status(201).JSON(fiber.Map{
		"message": "Discussion created",
		"user_id": userID,
	})
}

// ListDiscussions returns list of discussions
func (h *DiscussionHandler) ListDiscussions(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"data": []fiber.Map{},
		"pagination": fiber.Map{
			"page":  1,
			"limit": 20,
			"total": 0,
		},
	})
}

// GetDiscussion returns specific discussion with replies
func (h *DiscussionHandler) GetDiscussion(c *fiber.Ctx) error {
	discussionID := c.Params("id")

	return c.JSON(fiber.Map{
		"id":          discussionID,
		"title":       "",
		"content":     "",
		"user":        fiber.Map{"id": 0, "name": ""},
		"replies":     []fiber.Map{},
		"views_count": 0,
		"status":      "open",
	})
}

// UpdateDiscussion updates discussion
func (h *DiscussionHandler) UpdateDiscussion(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	discussionID := c.Params("id")

	// TODO: Verify ownership and call service
	return c.JSON(fiber.Map{
		"message": "Discussion updated",
		"id":      discussionID,
		"user_id": userID,
	})
}
