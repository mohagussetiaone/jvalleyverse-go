package handler

import (
	"jvalleyverse/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ReplyHandler struct {
	replySvc service.IReplyService
}

func NewReplyHandler() *ReplyHandler {
	return &ReplyHandler{replySvc: service.GetReplyService()}
}

// CreateReply creates reply to discussion
func (h *ReplyHandler) CreateReply(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	discussionID := c.Params("id")
	if discussionID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid discussion ID"})
	}

	var input struct {
		Content  string  `json:"content"`
		ParentID *string `json:"parent_id"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	reply, err := h.replySvc.CreateReply(c.UserContext(), userID, discussionID, input.Content, input.ParentID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(reply)
}

// UpdateReply updates reply
func (h *ReplyHandler) UpdateReply(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	replyID := c.Params("id")
	if replyID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid reply ID"})
	}

	var input struct {
		Content string `json:"content"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.replySvc.UpdateReply(c.UserContext(), replyID, userID, input.Content); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Reply updated"})
}

// DeleteReply deletes reply
func (h *ReplyHandler) DeleteReply(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	replyID := c.Params("id")
	if replyID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid reply ID"})
	}

	role, _ := c.Locals("role").(string)
	isAdmin := role == "admin"

	if err := h.replySvc.DeleteReply(c.UserContext(), replyID, userID, isAdmin); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Reply deleted"})
}
