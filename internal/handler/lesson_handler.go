package handler

import (
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type LessonHandler struct {
	lessonSvc service.ILessonService
}

func NewLessonHandler(lessonSvc service.ILessonService) *LessonHandler {
	return &LessonHandler{lessonSvc: lessonSvc}
}

func (h *LessonHandler) GetPublicLessonByID(c *fiber.Ctx) error {
	lessonID := c.Params("id")

	if lessonID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Lesson ID is required",
		})
	}

	data, err := h.lessonSvc.GetPublicLessonByID(
		c.UserContext(),
		lessonID,
	)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Lesson not found",
		})
	}

	return c.JSON(data)
}

func (h *LessonHandler) ListLessonsByCourse(c *fiber.Ctx) error {
	courseID := c.Params("course_id")
	if courseID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Course ID is required"})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset := (page - 1) * limit

	lessons, total, err := h.lessonSvc.ListLessonsByCourse(c.UserContext(), courseID, limit, offset)
	if err != nil {
		return safeError(c, 500, err)
	}

	return c.JSON(fiber.Map{
		"data": lessons,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *LessonHandler) ListLessonsBySection(c *fiber.Ctx) error {
	sectionID := c.Params("section_id")
	if sectionID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Section ID is required"})
	}

	lessons, total, err := h.lessonSvc.ListLessonsBySection(c.UserContext(), sectionID)
	if err != nil {
		return safeError(c, 500, err)
	}

	return c.JSON(fiber.Map{
		"data":  lessons,
		"total": total,
	})
}

func (h *LessonHandler) GetLessonBySlug(c *fiber.Ctx) error {
	courseID := c.Params("course_id")
	slug := c.Params("slug")
	userID, ok := c.Locals("userID").(string)
	if !ok {
		userID = ""
	}

	data, err := h.lessonSvc.GetLessonBySlug(c.UserContext(), courseID, slug, userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Lesson not found"})
	}

	return c.JSON(data)
}

func (h *LessonHandler) StartLesson(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	lessonID := c.Params("id")
	if lessonID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid lesson id"})
	}

	progress, err := h.lessonSvc.StartLesson(c.UserContext(), userID, lessonID)
	if err != nil {
		return safeError(c, mapServiceErrorToStatus(err), err)
	}

	return c.JSON(fiber.Map{
		"message":  "Lesson started!",
		"progress": progress,
	})
}

func (h *LessonHandler) UpdateProgress(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	lessonID := c.Params("id")
	if lessonID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid lesson id"})
	}

	var input struct {
		Percentage int    `json:"progress_percentage"`
		Notes      string `json:"notes"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	progress, err := h.lessonSvc.UpdateProgress(c.UserContext(), userID, lessonID, input.Percentage, input.Notes)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(progress)
}

func (h *LessonHandler) CompleteLesson(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	lessonID := c.Params("id")
	if lessonID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid lesson id"})
	}

	data, err := h.lessonSvc.CompleteLesson(c.UserContext(), userID, lessonID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(data)
}

func (h *LessonHandler) CreateLesson(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	var input domain.Lesson
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if input.CourseID == "" || input.SectionID == "" || input.Title == "" || input.Slug == "" {
		return c.Status(400).JSON(fiber.Map{"error": "course_id, section_id, title, and slug are required"})
	}

	lesson, err := h.lessonSvc.AdminCreateLesson(c.UserContext(), userID, input)
	if err != nil {
		return c.Status(mapServiceErrorToStatus(err)).JSON(fiber.Map{"error": err.Error()})
	}

	logAdminAction(c.UserContext(), userID, "create", "lesson", lesson.ID, "Title: "+input.Title)
	return c.Status(201).JSON(lesson)
}

func (h *LessonHandler) CreateLessonDetail(c *fiber.Ctx) error {
	lessonID := c.Params("id")
	if lessonID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid lesson id"})
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

	detail, err := h.lessonSvc.AdminCreateLessonDetail(c.UserContext(), lessonID, input.About, input.Rules, input.Tools, input.ResourceMedia, input.Resources)
	if err != nil {
		return c.Status(mapServiceErrorToStatus(err)).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(detail)
}

func (h *LessonHandler) UpdateLesson(c *fiber.Ctx) error {
	adminID := c.Locals("userID").(string)
	lessonID := c.Params("id")
	if lessonID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid lesson id"})
	}

	var input domain.Lesson
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	lesson, err := h.lessonSvc.AdminUpdateLesson(c.UserContext(), adminID, lessonID, input)
	if err != nil {
		return c.Status(mapServiceErrorToStatus(err)).JSON(fiber.Map{"error": err.Error()})
	}

	logAdminAction(c.UserContext(), adminID, "update", "lesson", lessonID, "Title: "+input.Title)
	return c.JSON(lesson)
}

func (h *LessonHandler) DeleteLesson(c *fiber.Ctx) error {
	adminID := c.Locals("userID").(string)
	lessonID := c.Params("id")
	if lessonID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid lesson id"})
	}

	if err := h.lessonSvc.AdminDeleteLesson(c.UserContext(), adminID, lessonID); err != nil {
		return c.Status(mapServiceErrorToStatus(err)).JSON(fiber.Map{"error": err.Error()})
	}

	logAdminAction(c.UserContext(), adminID, "delete", "lesson", lessonID, "")
	return c.JSON(fiber.Map{"message": "Lesson deleted"})
}
