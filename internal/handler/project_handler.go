package handler

import (
	"jvalleyverse/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ProjectHandler struct {
	projectSvc service.IProjectService
}

func NewProjectHandler(projectSvc service.IProjectService) *ProjectHandler {
	return &ProjectHandler{projectSvc: projectSvc}
}

// ListPublicProjects lists public projects (no auth required)
func (h *ProjectHandler) ListPublicProjects(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}

	limit := c.QueryInt("limit", 10)
	if limit < 1 {
		limit = 10
	}

	projects, err := h.projectSvc.ListPublicProjects(
		c.UserContext(),
		page,
		limit,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": projects,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": len(projects),
		},
	})
}

// CreateProject creates new project (admin only)
func (h *ProjectHandler) CreateProject(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Thumbnail   string `json:"thumbnail"`
		CategoryID  string `json:"category_id"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if input.Title == "" || input.CategoryID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "title and category_id are required"})
	}

	project, err := h.projectSvc.CreateProject(c.UserContext(), userID, input.Title, input.Description, input.Thumbnail, input.CategoryID)
	if err != nil {
		return c.Status(mapServiceErrorToStatus(err)).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(project)
}

// ListProjects lists all projects (admin only)
func (h *ProjectHandler) ListProjects(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	limit := 20

	projects, err := h.projectSvc.ListProjects(c.UserContext(), page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": projects,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": len(projects),
		},
	})
}

// UpdateProject updates project (admin only)
func (h *ProjectHandler) UpdateProject(c *fiber.Ctx) error {
	adminID := c.Locals("userID").(string)
	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Visibility  string `json:"visibility"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.projectSvc.UpdateProject(c.UserContext(), projectID, adminID, input.Title, input.Description, input.Visibility); err != nil {
		return c.Status(mapServiceErrorToStatus(err)).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Project updated"})
}

// DeleteProject deletes project (admin only)
func (h *ProjectHandler) DeleteProject(c *fiber.Ctx) error {
	adminID := c.Locals("userID").(string)
	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	if err := h.projectSvc.DeleteProject(c.UserContext(), projectID, adminID); err != nil {
		return c.Status(mapServiceErrorToStatus(err)).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Project deleted"})
}
