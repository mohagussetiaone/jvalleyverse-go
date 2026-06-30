package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/dto"
	"jvalleyverse/internal/service"
)

type FaqHandler struct {
	faqSvc service.IFaqService
}

func NewFaqHandler(faqSvc service.IFaqService) *FaqHandler {
	return &FaqHandler{faqSvc: faqSvc}
}

// GET /api/faqs — public, no auth required, returns active FAQs only
func (h *FaqHandler) ListPublic(c *fiber.Ctx) error {
	faqs, err := h.faqSvc.ListPublicFAQs(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch FAQs"})
	}

	if faqs == nil {
		faqs = []dto.FAQItem{}
	}

	return c.JSON(fiber.Map{
		"data": faqs,
	})
}

// GET /api/admin/faqs — admin only, all FAQs with pagination
func (h *FaqHandler) ListAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	faqs, total, err := h.faqSvc.ListAllFAQs(c.UserContext(), page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if faqs == nil {
		faqs = []dto.FAQItem{}
	}

	return c.JSON(fiber.Map{
		"data": faqs,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GET /api/admin/faqs/:id — admin only, get single FAQ
func (h *FaqHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid FAQ ID"})
	}

	faq, err := h.faqSvc.GetFAQByID(c.UserContext(), id)
	if err != nil {
		return c.Status(mapFaqErr(err)).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(faq)
}

// POST /api/admin/faqs — admin only, create FAQ
func (h *FaqHandler) Create(c *fiber.Ctx) error {
	var input struct {
		Question   string `json:"question"`
		Answer     string `json:"answer"`
		Category   string `json:"category"`
		OrderIndex int    `json:"order_index"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if input.Question == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "question is required"})
	}
	if input.Answer == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "answer is required"})
	}

	faq, err := h.faqSvc.CreateFAQ(c.UserContext(), input.Question, input.Answer, input.Category, input.OrderIndex)
	if err != nil {
		return c.Status(mapFaqErr(err)).JSON(fiber.Map{"error": err.Error()})
	}

	// Log admin action
	userID := c.Locals("userID").(string)
	logAdminAction(c.UserContext(), userID, "create", "faq", faq.ID, "Question: "+input.Question)

	return c.Status(fiber.StatusCreated).JSON(faq)
}

// PUT /api/admin/faqs/:id — admin only, update FAQ
func (h *FaqHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid FAQ ID"})
	}

	var input struct {
		Question   string `json:"question"`
		Answer     string `json:"answer"`
		Category   string `json:"category"`
		OrderIndex int    `json:"order_index"`
		IsActive   *bool  `json:"is_active"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	faq, err := h.faqSvc.UpdateFAQ(c.UserContext(), id, input.Question, input.Answer, input.Category, input.OrderIndex, isActive)
	if err != nil {
		return c.Status(mapFaqErr(err)).JSON(fiber.Map{"error": err.Error()})
	}

	// Log admin action
	userID := c.Locals("userID").(string)
	logAdminAction(c.UserContext(), userID, "update", "faq", id, "")

	return c.JSON(faq)
}

// DELETE /api/admin/faqs/:id — admin only, delete FAQ
func (h *FaqHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid FAQ ID"})
	}

	if err := h.faqSvc.DeleteFAQ(c.UserContext(), id); err != nil {
		return c.Status(mapFaqErr(err)).JSON(fiber.Map{"error": err.Error()})
	}

	// Log admin action
	userID := c.Locals("userID").(string)
	logAdminAction(c.UserContext(), userID, "delete", "faq", id, "")

	return c.JSON(fiber.Map{"message": "FAQ deleted"})
}

func mapFaqErr(err error) int {
	switch err {
	case domain.ErrNotFound:
		return fiber.StatusNotFound
	case domain.ErrInvalidInput:
		return fiber.StatusBadRequest
	default:
		return fiber.StatusInternalServerError
	}
}
