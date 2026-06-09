package handler

import (
	"jvalleyverse/internal/service"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	categorySvc service.ICategoryService
}

func NewCategoryHandler(categorySvc service.ICategoryService) *CategoryHandler {
	return &CategoryHandler{categorySvc: categorySvc}
}

// ListCategories returns all categories (public, no auth required)
func (h *CategoryHandler) ListCategories(c *fiber.Ctx) error {
	categories, err := h.categorySvc.ListCategories(c.UserContext())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(categories)
}

// GetCategoryBySlug returns category detail by slug (public, no auth required)
func (h *CategoryHandler) GetCategoryBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid slug"})
	}

	category, err := h.categorySvc.GetCategoryBySlug(c.UserContext(), slug)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Category not found"})
	}

	return c.JSON(category)
}

// ListProjectsByCategory returns projects belonging to a category (public, no auth required)
func (h *CategoryHandler) ListProjectsByCategory(c *fiber.Ctx) error {
	categoryID := c.Params("category_id")
	if categoryID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid category ID"})
	}

	projects, err := h.categorySvc.ListProjectsByCategoryID(c.UserContext(), categoryID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(projects)
}

// CreateCategory creates a new category (admin only)
func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var input struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if input.Name == "" || input.Slug == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name and slug are required"})
	}

	category, err := h.categorySvc.CreateCategory(c.UserContext(), input.Name, input.Slug, input.Description)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(category)
}

// UpdateCategory updates a category (admin only)
func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	categoryID := c.Params("id")
	if categoryID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid category ID"})
	}

	var input struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	category, err := h.categorySvc.UpdateCategory(c.UserContext(), categoryID, input.Name, input.Slug, input.Description)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(category)
}

// DeleteCategory deletes a category (admin only)
func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	categoryID := c.Params("id")
	if categoryID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid category ID"})
	}

	if err := h.categorySvc.DeleteCategory(c.UserContext(), categoryID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "category deleted"})
}
