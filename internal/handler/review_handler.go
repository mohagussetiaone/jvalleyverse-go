package handler

import (
	"jvalleyverse/internal/dto"
	"jvalleyverse/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ReviewHandler struct {
	reviewSvc service.IReviewService
}

func NewReviewHandler(reviewSvc service.IReviewService) *ReviewHandler {
	return &ReviewHandler{reviewSvc: reviewSvc}
}

func (h *ReviewHandler) ListCourseReviews(c *fiber.Ctx) error {
	courseID := c.Params("course_id")
	if courseID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid course ID"})
	}

	page := clampPage(c.QueryInt("page", DefaultPage))
	limit := clampLimit(c.QueryInt("limit", DefaultLimit), DefaultLimit)

	reviews, total, err := h.reviewSvc.ListCourseReviews(c.UserContext(), courseID, page, limit)
	if err != nil {
		return safeError(c, 500, err)
	}

	if reviews == nil {
		reviews = []dto.ReviewItem{}
	}

	return c.JSON(fiber.Map{
		"data": reviews,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *ReviewHandler) ListLessonReviews(c *fiber.Ctx) error {
	lessonID := c.Params("id")
	if lessonID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid lesson ID"})
	}

	page := clampPage(c.QueryInt("page", DefaultPage))
	limit := clampLimit(c.QueryInt("limit", DefaultLimit), DefaultLimit)

	reviews, total, err := h.reviewSvc.ListLessonReviews(c.UserContext(), lessonID, page, limit)
	if err != nil {
		return safeError(c, 500, err)
	}

	if reviews == nil {
		reviews = []dto.ReviewItem{}
	}

	return c.JSON(fiber.Map{
		"data": reviews,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *ReviewHandler) CreateReview(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var input struct {
		CourseID string `json:"course_id"`
		LessonID string `json:"lesson_id"`
		Rating   int    `json:"rating"`
		Message  string `json:"message"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if input.Rating < 1 || input.Rating > 5 || input.Message == "" {
		return c.Status(400).JSON(fiber.Map{"error": "rating (1-5) and message are required"})
	}

	review, err := h.reviewSvc.CreateReview(c.UserContext(), userID, input.CourseID, input.LessonID, input.Rating, input.Message)
	if err != nil {
		return c.Status(mapServiceErrorToStatus(err)).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(review)
}

func (h *ReviewHandler) UpdateReview(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	reviewID := c.Params("id")
	if reviewID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid review ID"})
	}

	var input struct {
		Rating  int    `json:"rating"`
		Message string `json:"message"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	review, err := h.reviewSvc.UpdateReview(c.UserContext(), reviewID, userID, input.Rating, input.Message)
	if err != nil {
		return c.Status(mapServiceErrorToStatus(err)).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(review)
}

func (h *ReviewHandler) DeleteReview(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	reviewID := c.Params("id")
	if reviewID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid review ID"})
	}

	if err := h.reviewSvc.DeleteReview(c.UserContext(), reviewID, userID); err != nil {
		return c.Status(mapServiceErrorToStatus(err)).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Review deleted"})
}
