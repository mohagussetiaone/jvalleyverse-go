package handler

import (
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ReplyHandler struct {
	replySvc service.IReplyService
}

func NewReplyHandler() *ReplyHandler {
	return &ReplyHandler{replySvc: service.GetReplyService()}
}

// GetMyReplies returns current user's replies
func (h *ReplyHandler) GetMyReplies(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	page := clampPage(c.QueryInt("page", DefaultPage))
	limit := clampLimit(c.QueryInt("limit", DefaultLimit), DefaultLimit)

	replies, total, err := h.replySvc.ListRepliesByUser(c.UserContext(), userID, page, limit)
	if err != nil {
		return safeError(c, 500, err)
	}

	return c.JSON(fiber.Map{
		"data": replies,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
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
		return safeError(c, mapServiceErrorToStatus(err), err)
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
		if err == domain.ErrForbidden {
			return c.Status(403).JSON(fiber.Map{"error": "You do not own this reply"})
		}
		return safeError(c, mapServiceErrorToStatus(err), err)
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
		if err == domain.ErrForbidden {
			return c.Status(403).JSON(fiber.Map{"error": "You do not own this reply"})
		}
		return safeError(c, mapServiceErrorToStatus(err), err)
	}

	return c.JSON(fiber.Map{"message": "Reply deleted"})
}

// LikeReply likes a reply (POST /api/replies/:id/like)
func (h *ReplyHandler) LikeReply(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	replyID := c.Params("id")
	if replyID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid reply ID"})
	}

	if err := h.replySvc.LikeReply(c.UserContext(), userID, replyID); err != nil {
		if err == domain.ErrNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "Reply not found"})
		}
		return safeError(c, mapServiceErrorToStatus(err), err)
	}

	return c.JSON(fiber.Map{"message": "Reply liked"})
}

// MarkBestReply marks a reply as best answer (POST /api/replies/:id/best)
func (h *ReplyHandler) MarkBestReply(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	replyID := c.Params("id")
	if replyID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid reply ID"})
	}

	var input struct {
		DiscussionID string `json:"discussion_id"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if input.DiscussionID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "discussion_id is required"})
	}

	if err := h.replySvc.MarkBestReply(c.UserContext(), replyID, input.DiscussionID, userID); err != nil {
		if err == domain.ErrForbidden {
			return c.Status(403).JSON(fiber.Map{"error": "Only discussion owner can mark best answer"})
		}
		return safeError(c, mapServiceErrorToStatus(err), err)
	}

	return c.JSON(fiber.Map{"message": "Reply marked as best answer"})
}
