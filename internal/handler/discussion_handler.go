package handler

import (
	"jvalleyverse/internal/domain"
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
		Title       string  `json:"title"`
		Content     string  `json:"content"`
		LessonID    *string `json:"lesson_id"`
		StudyCaseID *string `json:"study_case_id"`
		CategoryID  string  `json:"category_id"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	discussion, err := h.discussionSvc.CreateDiscussion(c.UserContext(), userID, input.Title, input.Content, input.LessonID, input.StudyCaseID, input.CategoryID)
	if err != nil {
		return safeError(c, mapServiceErrorToStatus(err), err)
	}

	return c.Status(201).JSON(discussion)
}

// ListDiscussions returns list of discussions
func (h *DiscussionHandler) ListDiscussions(c *fiber.Ctx) error {
	studyCaseID := c.Query("study_case_id")
	var studyCaseIDPtr *string
	if studyCaseID != "" {
		studyCaseIDPtr = &studyCaseID
	}

	lessonID := c.Query("lesson_id")
	var lessonIDPtr *string
	if lessonID != "" {
		lessonIDPtr = &lessonID
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	discussions, total, err := h.discussionSvc.ListDiscussions(c.UserContext(), page, limit, lessonIDPtr, studyCaseIDPtr, nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": discussions,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
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
		if err == domain.ErrForbidden {
			return c.Status(403).JSON(fiber.Map{"error": "You do not own this discussion"})
		}
		return safeError(c, mapServiceErrorToStatus(err), err)
	}

	return c.JSON(fiber.Map{"message": "Discussion updated"})
}

// DeleteDiscussion deletes a discussion (owner or admin)
func (h *DiscussionHandler) DeleteDiscussion(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	discussionID := c.Params("id")
	if discussionID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid discussion ID"})
	}

	role, _ := c.Locals("role").(string)
	isAdmin := role == "admin"

	if err := h.discussionSvc.DeleteDiscussion(c.UserContext(), discussionID, userID, isAdmin); err != nil {
		if err == domain.ErrForbidden {
			return c.Status(403).JSON(fiber.Map{"error": "You do not own this discussion"})
		}
		return safeError(c, mapServiceErrorToStatus(err), err)
	}

	return c.JSON(fiber.Map{"message": "Discussion deleted"})
}

// ListMyDiscussions returns current user's discussions
func (h *DiscussionHandler) ListMyDiscussions(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	discussions, total, err := h.discussionSvc.ListUserDiscussions(c.UserContext(), userID, page, limit)
	if err != nil {
		return safeError(c, 500, err)
	}

	return c.JSON(fiber.Map{
		"data": discussions,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// CloseDiscussion closes a discussion (owner only)
func (h *DiscussionHandler) CloseDiscussion(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	discussionID := c.Params("id")
	if discussionID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid discussion ID"})
	}

	if err := h.discussionSvc.CloseDiscussion(c.UserContext(), discussionID, userID); err != nil {
		if err == domain.ErrForbidden {
			return c.Status(403).JSON(fiber.Map{"error": "You do not own this discussion"})
		}
		return safeError(c, mapServiceErrorToStatus(err), err)
	}

	return c.JSON(fiber.Map{"message": "Discussion closed"})
}
