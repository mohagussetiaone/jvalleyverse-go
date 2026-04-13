package handler

import (
	"jvalleyverse/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ReplyHandler struct {
	replySvc *service.ReplyService
}

func NewReplyHandler() *ReplyHandler {
	return &ReplyHandler{replySvc: service.NewReplyService()}
}

// CreateReply creates reply to discussion
func (h *ReplyHandler) CreateReply(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	discussionID := c.Params("id")

	var input struct {
		Content  string `json:"content"`
		ParentID *uint  `json:"parent_id"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// TODO: Call service
	return c.Status(201).JSON(fiber.Map{
		"message":        "Reply created",
		"discussion_id":  discussionID,
		"user_id":        userID,
		"parent_id":      input.ParentID,
	})
}

// UpdateReply updates reply
func (h *ReplyHandler) UpdateReply(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	replyID := c.Params("id")

	var input struct {
		Content string `json:"content"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// TODO: Verify ownership and call service
	return c.JSON(fiber.Map{
		"message": "Reply updated",
		"id":      replyID,
		"user_id": userID,
	})
}

// DeleteReply deletes reply
func (h *ReplyHandler) DeleteReply(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	replyID := c.Params("id")

	// TODO: Verify ownership and call service
	return c.JSON(fiber.Map{
		"message": "Reply deleted",
		"id":      replyID,
		"user_id": userID,
	})
}
