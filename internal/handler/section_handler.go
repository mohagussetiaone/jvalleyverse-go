package handler

import (
	"jvalleyverse/internal/service"

	"github.com/gofiber/fiber/v2"
)

type SectionHandler struct {
	sectionSvc service.ISectionService
}

func NewSectionHandler(sectionSvc service.ISectionService) *SectionHandler {
	return &SectionHandler{sectionSvc: sectionSvc}
}

func (h *SectionHandler) GetCourseWithSections(c *fiber.Ctx) error {
	courseID := c.Params("course_id")
	if courseID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Course ID is required"})
	}

	userID, _ := c.Locals("userID").(string)

	course, err := h.sectionSvc.GetCourseWithSections(c.UserContext(), courseID, userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Course not found"})
	}

	return c.JSON(course)
}

func (h *SectionHandler) ListSectionsByCourse(c *fiber.Ctx) error {
	courseID := c.Params("course_id")
	if courseID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Course ID is required"})
	}

	sections, err := h.sectionSvc.ListSectionsByCourse(c.UserContext(), courseID)
	if err != nil {
		return safeError(c, 500, err)
	}

	return c.JSON(fiber.Map{"data": sections})
}

func (h *SectionHandler) GetSection(c *fiber.Ctx) error {
	sectionID := c.Params("section_id")
	if sectionID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Section ID is required"})
	}

	section, err := h.sectionSvc.GetSection(c.UserContext(), sectionID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Section not found"})
	}

	return c.JSON(section)
}

func (h *SectionHandler) CreateSection(c *fiber.Ctx) error {
	adminID := c.Locals("userID").(string)
	courseID := c.Params("course_id")
	if courseID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Course ID is required"})
	}

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		OrderIndex  int    `json:"order_index"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if input.Title == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Title is required"})
	}

	section, err := h.sectionSvc.CreateSection(c.UserContext(), adminID, courseID, input.Title, input.Description, input.OrderIndex)
	if err != nil {
		return c.Status(mapServiceErrorToStatus(err)).JSON(fiber.Map{"error": err.Error()})
	}

	logAdminAction(c.UserContext(), adminID, "create", "section", section.ID, "Title: "+input.Title)
	return c.Status(201).JSON(section)
}

func (h *SectionHandler) UpdateSection(c *fiber.Ctx) error {
	adminID := c.Locals("userID").(string)
	sectionID := c.Params("section_id")
	if sectionID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Section ID is required"})
	}

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		OrderIndex  int    `json:"order_index"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	section, err := h.sectionSvc.UpdateSection(c.UserContext(), adminID, sectionID, input.Title, input.Description, input.OrderIndex)
	if err != nil {
		return c.Status(mapServiceErrorToStatus(err)).JSON(fiber.Map{"error": err.Error()})
	}

	logAdminAction(c.UserContext(), adminID, "update", "section", sectionID, "Title: "+input.Title)
	return c.JSON(section)
}

func (h *SectionHandler) DeleteSection(c *fiber.Ctx) error {
	adminID := c.Locals("userID").(string)
	sectionID := c.Params("section_id")
	if sectionID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Section ID is required"})
	}

	if err := h.sectionSvc.DeleteSection(c.UserContext(), adminID, sectionID); err != nil {
		return c.Status(mapServiceErrorToStatus(err)).JSON(fiber.Map{"error": err.Error()})
	}

	logAdminAction(c.UserContext(), adminID, "delete", "section", sectionID, "")
	return c.JSON(fiber.Map{"message": "Section deleted"})
}
