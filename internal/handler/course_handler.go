package handler

import (
	"encoding/json"
	"jvalleyverse/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
)

type CourseHandler struct {
	courseSvc service.ICourseService
}

func NewCourseHandler(courseSvc service.ICourseService) *CourseHandler {
	return &CourseHandler{courseSvc: courseSvc}
}

// ListPublicCourses godoc
// @Summary      List public courses
// @Description  Get paginated list of public courses with optional filters by category_id, min_price, and max_price
// @Tags         Courses
// @Param        page         query int     false  "Page number (default: 1)"
// @Param        limit        query int     false  "Items per page (default: 10)"
// @Param        category_id  query string  false  "Filter by category ID"
// @Param        min_price    query number  false  "Minimum price filter"
// @Param        max_price    query number  false  "Maximum price filter"
// @Success      200  {object}  map[string]interface{}
// @Router       /courses [get]
func (h *CourseHandler) ListPublicCourses(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}

	limit := c.QueryInt("limit", 10)
	if limit < 1 {
		limit = 10
	}

	// Optional filters
	categoryID := c.Query("category_id", "")
	minPrice := c.Query("min_price", "")
	maxPrice := c.Query("max_price", "")

	var filter *service.CourseListFilter

	if categoryID != "" || minPrice != "" || maxPrice != "" {
		f := &service.CourseListFilter{}
		if categoryID != "" {
			f.CategoryID = &categoryID
		}
		if minPrice != "" {
			if val, err := strconv.ParseFloat(minPrice, 64); err == nil {
				f.MinPrice = &val
			}
		}
		if maxPrice != "" {
			if val, err := strconv.ParseFloat(maxPrice, 64); err == nil {
				f.MaxPrice = &val
			}
		}
		filter = f
	}

	userID, hasAuth := c.Locals("userID").(string)

	var total int64
	var err error
	var courses interface{}

	if hasAuth && userID != "" {
		courses, total, err = h.courseSvc.ListPublicCoursesWithEnrollment(
			c.UserContext(),
			userID,
			page,
			limit,
			filter,
		)
	} else {
		courses, total, err = h.courseSvc.ListPublicCourses(
			c.UserContext(),
			page,
			limit,
			filter,
		)
	}

	if err != nil {
		return safeError(c, 500, err)
	}

	return c.JSON(fiber.Map{
		"data": courses,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *CourseHandler) CreateCourse(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var input struct {
		Title              string          `json:"title"`
		Description        string          `json:"description"`
		Thumbnail          string          `json:"thumbnail"`
		CategoryID         string          `json:"category_id"`
		MentorID           string          `json:"mentor_id"`
		Price              float64         `json:"price"`
		Hours              int             `json:"hours"`
		LearningObjectives json.RawMessage `json:"learning_objectives"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if input.Title == "" || input.CategoryID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "title and category_id are required"})
	}
	if input.Price < 0 {
		return c.Status(400).JSON(fiber.Map{"error": "price must be >= 0"})
	}

	course, err := h.courseSvc.CreateCourse(c.UserContext(), userID, input.Title, input.Description, input.Thumbnail, input.CategoryID, input.MentorID, input.Price, input.Hours, datatypes.JSON(input.LearningObjectives))
	if err != nil {
		return c.Status(mapServiceErrorToStatus(err)).JSON(fiber.Map{"error": err.Error()})
	}

	logAdminAction(c.UserContext(), userID, "create", "course", course.ID, "Title: "+input.Title)
	return c.Status(201).JSON(course)
}

func (h *CourseHandler) UpdateCourse(c *fiber.Ctx) error {
	adminID := c.Locals("userID").(string)
	courseID := c.Params("id")
	if courseID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid course ID"})
	}

	var input struct {
		Title              string          `json:"title"`
		Description        string          `json:"description"`
		Price              float64         `json:"price"`
		Visibility         string          `json:"visibility"`
		LearningObjectives json.RawMessage `json:"learning_objectives"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if input.Price < 0 {
		return c.Status(400).JSON(fiber.Map{"error": "price must be >= 0"})
	}

	if err := h.courseSvc.UpdateCourse(c.UserContext(), courseID, adminID, input.Title, input.Description, input.Price, input.Visibility, datatypes.JSON(input.LearningObjectives)); err != nil {
		return c.Status(mapServiceErrorToStatus(err)).JSON(fiber.Map{"error": err.Error()})
	}

	logAdminAction(c.UserContext(), adminID, "update", "course", courseID, "Title: "+input.Title)
	return c.JSON(fiber.Map{"message": "Course updated"})
}

func (h *CourseHandler) DeleteCourse(c *fiber.Ctx) error {
	adminID := c.Locals("userID").(string)
	courseID := c.Params("id")
	if courseID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid course ID"})
	}

	if err := h.courseSvc.DeleteCourse(c.UserContext(), courseID, adminID); err != nil {
		return c.Status(mapServiceErrorToStatus(err)).JSON(fiber.Map{"error": err.Error()})
	}

	logAdminAction(c.UserContext(), adminID, "delete", "course", courseID, "")
	return c.JSON(fiber.Map{"message": "Course deleted"})
}

func (h *CourseHandler) EnrollCourse(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	courseID := c.Params("id")
	if courseID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid course ID"})
	}

	if err := h.courseSvc.EnrollCourse(c.UserContext(), userID, courseID); err != nil {
		return c.Status(mapServiceErrorToStatus(err)).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Successfully enrolled in course"})
}

func (h *CourseHandler) ListEnrolledCourses(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	limit := c.QueryInt("limit", 10)
	if limit < 1 {
		limit = 10
	}

	courses, total, err := h.courseSvc.ListEnrolledCourses(c.UserContext(), userID, page, limit)
	if err != nil {
		return safeError(c, 500, err)
	}

	return c.JSON(fiber.Map{
		"data": courses,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *CourseHandler) SetLastLesson(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	courseID := c.Params("id")
	if courseID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid course ID"})
	}

	var req struct {
		LessonID string `json:"lesson_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if req.LessonID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "lesson_id is required"})
	}

	if err := h.courseSvc.SetLastLesson(c.UserContext(), userID, courseID, req.LessonID); err != nil {
		return safeError(c, 500, err)
	}

	return c.JSON(fiber.Map{"message": "Last lesson updated"})
}
