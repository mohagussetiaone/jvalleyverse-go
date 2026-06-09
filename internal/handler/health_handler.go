package handler

import (
	"os"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health returns basic health status (local environment)
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
	})
}

// HealthDetailed returns system metrics for production/VPS
func (h *HealthHandler) HealthDetailed(c *fiber.Ctx) error {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = os.Getenv("ENV")
	}

	// Local environment - return basic status
	if env == "local" || env == "development" || env == "" {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"environment": "local",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
		})
	}

	// Production/VPS - return system metrics
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Get hostname
	hostname, _ := os.Hostname()

	return c.JSON(fiber.Map{
		"status":      "ok",
		"environment": env,
		"hostname":    hostname,
		"timestamp":   time.Now().UTC(),
		"version":     "1.0.0",
		"system": fiber.Map{
			"os":   runtime.GOOS,
			"arch": runtime.GOARCH,
		},
		"memory": fiber.Map{
			"allocated_mb":       float64(m.Alloc) / 1024 / 1024,
			"total_allocated_mb": float64(m.TotalAlloc) / 1024 / 1024,
			"system_mb":          float64(m.Sys) / 1024 / 1024,
			"num_gc":             m.NumGC,
		},
		"goroutines": runtime.NumGoroutine(),
	})
}
