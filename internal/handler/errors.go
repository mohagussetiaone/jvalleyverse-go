package handler

import (
	"errors"
	"jvalleyverse/internal/domain"

	"github.com/gofiber/fiber/v2"
)

func mapServiceErrorToStatus(err error) int {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return 400
	case errors.Is(err, domain.ErrUnauthorized):
		return 401
	case errors.Is(err, domain.ErrForbidden):
		return 403
	case errors.Is(err, domain.ErrNotFound),
		errors.Is(err, domain.ErrCourseNotFound),
		errors.Is(err, domain.ErrLessonNotFound),
		errors.Is(err, domain.ErrUserNotFound),
		errors.Is(err, domain.ErrStudyCaseNotFound):
		return 404
	case errors.Is(err, domain.ErrEmailExists):
		return 409
	default:
		return 500
	}
}

func safeError(c *fiber.Ctx, status int, err error) error {
	msg := map[int]string{
		400: "Invalid request",
		401: "Unauthorized",
		403: "Forbidden",
		404: "Resource not found",
		409: "Conflict",
		500: "Internal server error",
	}[status]

	if msg == "" {
		msg = "Internal server error"
	}

	// Only leak error details for non-500 errors
	if status < 500 {
		msg = err.Error()
	}

	return c.Status(status).JSON(fiber.Map{"error": msg})
}
