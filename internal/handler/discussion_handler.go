package handler

import (
	"jvalleyverse/internal/service"

	"github.com/gofiber/fiber/v2"
)

type DiscussionHandler struct {
	discussionSvc service.IDiscussionService
}

func NewDiscussionHandler() *DiscussionHandler {
	return &DiscussionHandler{discussionSvc: service.GetDiscussionService()}
}

// CreateDiscussion creates new discussion
func (h *DiscussionHandler) CreateDiscussion(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var input struct {
		Title      string  `json:"title"`
		Content    string  `json:"content"`
		ClassID    *string `json:"class_id"`
		CategoryID string  `json:"category_id"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	discussion, err := h.discussionSvc.CreateDiscussion(c.UserContext(), userID, input.Title, input.Content, input.ClassID, input.CategoryID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(discussion)
}

// ListDiscussions returns list of discussions
func (h *DiscussionHandler) ListDiscussions(c *fiber.Ctx) error {
	classID := c.Query("class_id")
	var classIDPtr *string
	if classID != "" {
		classIDPtr = &classID
	}

	discussions, err := h.discussionSvc.ListDiscussions(c.UserContext(), 1, 20, classIDPtr, nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": discussions,
		"pagination": fiber.Map{
			"page":  1,
			"limit": 20,
			"total": len(discussions),
		},
	})
}

// GetDiscussion returns specific discussion with replies
func (h *DiscussionHandler) GetDiscussion(c *fiber.Ctx) error {
	discussionID := c.Params("id")
	if discussionID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid discussion ID"})
	}

	data, err := h.discussionSvc.GetDiscussionWithReplies(c.UserContext(), discussionID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Discussion not found"})
	}

	return c.JSON(data)
}

// UpdateDiscussion updates discussion
func (h *DiscussionHandler) UpdateDiscussion(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	discussionID := c.Params("id")
	if discussionID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid discussion ID"})
	}

	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.discussionSvc.UpdateDiscussion(c.UserContext(), discussionID, userID, input.Title, input.Content); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Discussion updated"})
}
