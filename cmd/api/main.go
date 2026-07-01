package main

import (
	"fmt"
	"jvalleyverse/internal/minio"
	"jvalleyverse/internal/service"
	"jvalleyverse/pkg/config"
	"jvalleyverse/pkg/database"
	"jvalleyverse/pkg/middleware"
	"jvalleyverse/pkg/redis"
	"jvalleyverse/pkg/routes"
	"jvalleyverse/pkg/swagger"
	"log"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Load config
	config.LoadConfig()
	// Connect DB (required)
	database.ConnectDB()
	// Initialize service layer (depends on database connection)
	service.InitServices(database.DB)

	// Auto-seed initial data if tables are empty
	if err := service.SeedInitialData(database.DB); err != nil {
		log.Fatalf("Failed to auto-seed data: %v", err)
	}
	// Connect Redis (optional - app still works without it)
	redis.ConnectRedis()
	// Connect MinIO (optional - app still works without it)
	minio.ConnectMinio()

	app := fiber.New(fiber.Config{
		AppName:   "JValleyVerse API v1.0.0",
		BodyLimit: 10 * 1024 * 1024,
	})

	// Logger middleware
	app.Use(logger.New())

	// Global middleware
	app.Use(middleware.SetupCORS())
	app.Use(middleware.SecurityHeaders())
	app.Use(middleware.RateLimiter())

	// ==================== PROMETHEUS METRICS ====================
	prometheus := fiberprometheus.New("jvalleyverse")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	// ==================== SWAGGER DOCUMENTATION ====================
	app.Get("/docs", swagger.SwaggerHandler)
	app.Get("/api/docs/openapi.json", swagger.OpenAPIHandler)

	// ==================== SETUP ALL ROUTES ====================
	routes.SetupRoutes(app)

	// Start server
	port := config.AppConfig.Port
	fmt.Printf("\n✅ Server starting on http://localhost:%s\n", port)
	fmt.Printf("📖 Swagger UI: http://localhost:%s/docs\n", port)
	fmt.Printf("📄 OpenAPI Spec: http://localhost:%s/api/docs/openapi.json\n\n", port)

	log.Fatal(app.Listen(":" + port))
}
