package handler

import (
	"jvalleyverse/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ProjectHandler struct {
	projectSvc *service.ProjectService
}

func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{projectSvc: service.NewProjectService()}
}

// CreateProject creates new project (admin only)
func (h *ProjectHandler) CreateProject(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Thumbnail   string `json:"thumbnail"`
		CategoryID  uint   `json:"category_id"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	project, err := h.projectSvc.CreateProject(userID, input.Title, input.Description, input.Thumbnail, input.CategoryID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(project)
}

// ListProjects lists all projects (admin only)
func (h *ProjectHandler) ListProjects(c *fiber.Ctx) error {
	page, _ := c.ParamsInt("page")
	if page < 1 { page = 1 }
	limit := 20

	projects, err := h.projectSvc.ListProjects(page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": projects,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": len(projects), // Note: service should return total count for real pagination
		},
	})
}

// UpdateProject updates project (admin only)
func (h *ProjectHandler) UpdateProject(c *fiber.Ctx) error {
	adminID := c.Locals("userID").(uint)
	projectID, _ := c.ParamsInt("id")

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Visibility  string `json:"visibility"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := h.projectSvc.UpdateProject(uint(projectID), adminID, input.Title, input.Description, input.Visibility)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Project updated"})
}

// DeleteProject deletes project (admin only)
func (h *ProjectHandler) DeleteProject(c *fiber.Ctx) error {
	adminID := c.Locals("userID").(uint)
	projectID, _ := c.ParamsInt("id")

	err := h.projectSvc.DeleteProject(uint(projectID), adminID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Project deleted"})
}
