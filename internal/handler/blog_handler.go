package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/service"
)

type BlogHandler struct {
	blogSvc service.IBlogService
}

func NewBlogHandler(blogSvc service.IBlogService) *BlogHandler {
	return &BlogHandler{blogSvc: blogSvc}
}

// POST /api/admin/blogs
func (h *BlogHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req service.CreateBlogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Title is required"})
	}

	blog, err := h.blogSvc.CreateBlog(c.UserContext(), userID, req)
	if err != nil {
		return c.Status(mapBlogErr(err)).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(blog)
}

// GET /api/blogs
func (h *BlogHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")
	categoryID := c.Query("category_id", "")
	tag := c.Query("tag", "")

	items, pagination, err := h.blogSvc.ListBlogs(c.UserContext(), page, limit, search, categoryID, tag)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":       items,
		"pagination": pagination,
	})
}

// GET /api/blogs/:id
func (h *BlogHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	blog, err := h.blogSvc.GetBlogByID(c.UserContext(), id)
	if err != nil {
		return c.Status(mapBlogErr(err)).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(blog)
}

// PUT /api/admin/blogs/:id
func (h *BlogHandler) AdminUpdate(c *fiber.Ctx) error {
	blogID := c.Params("id")

	var req service.UpdateBlogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.blogSvc.AdminUpdateBlog(c.UserContext(), blogID, req); err != nil {
		return c.Status(mapBlogErr(err)).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Blog updated successfully"})
}

// DELETE /api/admin/blogs/:id
func (h *BlogHandler) AdminDelete(c *fiber.Ctx) error {
	blogID := c.Params("id")

	if err := h.blogSvc.AdminDeleteBlog(c.UserContext(), blogID); err != nil {
		return c.Status(mapBlogErr(err)).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Blog deleted successfully"})
}

// GET /api/users/me/blogs
func (h *BlogHandler) ListMyBlogs(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	status := c.Query("status", "")

	items, pagination, err := h.blogSvc.ListMyBlogs(c.UserContext(), userID, page, limit, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":       items,
		"pagination": pagination,
	})
}

func mapBlogErr(err error) int {
	switch err {
	case domain.ErrNotFound:
		return fiber.StatusNotFound
	case domain.ErrForbidden:
		return fiber.StatusForbidden
	default:
		return fiber.StatusInternalServerError
	}
}
