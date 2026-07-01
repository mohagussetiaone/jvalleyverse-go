package handler

import (
	"jvalleyverse/internal/service"

	"github.com/gofiber/fiber/v2"
)

type StatusHandler struct{}

func NewStatusHandler() *StatusHandler {
	return &StatusHandler{}
}

// SystemStatus returns comprehensive system health status (public, no auth)
// GET /api/system/status
func (h *StatusHandler) SystemStatus(c *fiber.Ctx) error {
	status := service.GetSystemStatus(c.UserContext())
	return c.JSON(status)
}
