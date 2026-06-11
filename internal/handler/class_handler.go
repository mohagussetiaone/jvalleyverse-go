package handler

import (
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ClassHandler struct {
	classSvc service.IClassService
}

func NewClassHandler(classSvc service.IClassService) *ClassHandler {
	return &ClassHandler{classSvc: classSvc}
}

// GetPublicClassByID
func (h *ClassHandler) GetPublicClassByID(c *fiber.Ctx) error {
	classID := c.Params("id")

	if classID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Class ID is required",
		})
	}

	data, err := h.classSvc.GetPublicClassByID(
		c.UserContext(),
		classID,
	)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Class not found",
		})
	}

	return c.JSON(data)
}

func (h *ClassHandler) ListClassesByProject(c *fiber.Ctx) error {
	projectID := c.Params("project_id")
	if projectID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Project ID is required"})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset := (page - 1) * limit

	classes, total, err := h.classSvc.ListClassesByProject(c.UserContext(), projectID, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": classes,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// ListClassesByPhase lists classes belonging to a specific phase (public)
// GET /api/projects/:project_id/phases/:phase_id/classes
func (h *ClassHandler) ListClassesByPhase(c *fiber.Ctx) error {
	phaseID := c.Params("phase_id")
	if phaseID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Phase ID is required"})
	}

	classes, total, err := h.classSvc.ListClassesByPhase(c.UserContext(), phaseID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":  classes,
		"total": total,
	})
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
	if input.ProjectID == "" || input.PhaseID == "" || input.Title == "" || input.Slug == "" {
		return c.Status(400).JSON(fiber.Map{"error": "project_id, phase_id, title, and slug are required"})
	}

	class, err := h.classSvc.AdminCreateClass(c.UserContext(), userID, input)
	if err != nil {
		return c.Status(mapServiceErrorToStatus(err)).JSON(fiber.Map{"error": err.Error()})
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
		return c.Status(mapServiceErrorToStatus(err)).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(detail)
}

func (h *ClassHandler) UpdateClass(c *fiber.Ctx) error {
	adminID := c.Locals("userID").(string)
	classID := c.Params("id")
	if classID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid class id"})
	}

	var input domain.Class
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	class, err := h.classSvc.AdminUpdateClass(c.UserContext(), adminID, classID, input)
	if err != nil {
		return c.Status(mapServiceErrorToStatus(err)).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(class)
}

func (h *ClassHandler) DeleteClass(c *fiber.Ctx) error {
	adminID := c.Locals("userID").(string)
	classID := c.Params("id")
	if classID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid class id"})
	}

	if err := h.classSvc.AdminDeleteClass(c.UserContext(), adminID, classID); err != nil {
		return c.Status(mapServiceErrorToStatus(err)).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Class deleted"})
}
