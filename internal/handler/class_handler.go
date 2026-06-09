package handler

import (
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ClassHandler struct {
	classSvc service.IClassService
}

func NewClassHandler(classSvc service.IClassService) *ClassHandler {
	return &ClassHandler{classSvc: classSvc}
}

func (h *ClassHandler) GetClassBySlug(c *fiber.Ctx) error {
	projectID := c.Params("project_id") // String CUID
	slug := c.Params("slug")
	userID, ok := c.Locals("userID").(string)
	if !ok {
		userID = "" // Anonymous
	}

	data, err := h.classSvc.GetClassBySlug(c.UserContext(), projectID, slug, userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Class not found"})
	}

	return c.JSON(data)
}

func (h *ClassHandler) StartClass(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	classID := c.Params("id") // String CUID
	if classID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid class id"})
	}

	progress, err := h.classSvc.StartClass(c.UserContext(), userID, classID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message":  "Class started!",
		"progress": progress,
	})
}

func (h *ClassHandler) UpdateProgress(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	classID := c.Params("id") // String CUID
	if classID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid class id"})
	}

	var input struct {
		Percentage int    `json:"progress_percentage"`
		Notes      string `json:"notes"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	progress, err := h.classSvc.UpdateProgress(c.UserContext(), userID, classID, input.Percentage, input.Notes)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(progress)
}

func (h *ClassHandler) Complete(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	classID := c.Params("id") // String CUID
	if classID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid class id"})
	}

	data, err := h.classSvc.CompleteClass(c.UserContext(), userID, classID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(data)
}

func (h *ClassHandler) CreateClass(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	var input domain.Class
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	class, err := h.classSvc.AdminCreateClass(c.UserContext(), userID, input)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(class)
}

func (h *ClassHandler) CreateClassDetail(c *fiber.Ctx) error {
	classID := c.Params("id") // String CUID
	if classID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid class id"})
	}

	var input struct {
		About         string      `json:"about"`
		Rules         string      `json:"rules"`
		Tools         interface{} `json:"tools"`
		ResourceMedia interface{} `json:"resource_media"`
		Resources     interface{} `json:"resources"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	detail, err := h.classSvc.AdminCreateClassDetail(c.UserContext(), classID, input.About, input.Rules, input.Tools, input.ResourceMedia, input.Resources)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(detail)
}

func (h *ClassHandler) UpdateClass(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented"})
}

func (h *ClassHandler) DeleteClass(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Not implemented"})
}
