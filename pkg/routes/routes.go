package routes

import (
	"jvalleyverse/internal/handler"
	"jvalleyverse/internal/service"
	"jvalleyverse/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App) {
	// ==================== PUBLIC ROUTES ====================
	// Auth endpoints
	authHandler := handler.NewAuthHandler(service.GetUserService())
	app.Post("/api/auth/register", authHandler.Register)
	app.Post("/api/auth/login", authHandler.Login)
	// For backward compatibility
	app.Post("/auth/register", authHandler.Register)
	app.Post("/auth/login", authHandler.Login)

	// List public content
	showcaseHandler := handler.NewShowcaseHandler(service.GetShowcaseService())
	app.Get("/api/leaderboard", handler.NewGamificationHandler().GetLeaderboard)
	app.Get("/api/showcases", showcaseHandler.ListShowcases)
	app.Get("/api/showcases/:id", showcaseHandler.GetShowcase)

	// ==================== PUBLIC ROUTES (no auth) ====================
	// Publicly accessible class details
	publicClassHandler := handler.NewClassHandler(service.GetClassService())
	app.Get("/api/projects/:project_id/classes/:slug", publicClassHandler.GetClassBySlug)

	// Health check (public, no auth required)
	healthHandler := handler.NewHealthHandler()
	app.Get("/api/health", healthHandler.Health)
	app.Get("/api/health/detailed", healthHandler.HealthDetailed)

	// Public Category endpoints (no auth required)
	categoryHandler := handler.NewCategoryHandler(service.GetCategoryService())
	app.Get("/api/categories", categoryHandler.ListCategories)
	app.Get("/api/categories/:slug", categoryHandler.GetCategoryBySlug)
	app.Get("/api/categories/:category_id/projects", categoryHandler.ListProjectsByCategory)

	// ==================== PROTECTED ROUTES ====================

	// User endpoints (MUST be before /api/users/:id to avoid route conflicts)
	userHandler := handler.NewUserHandler(service.GetUserService())
	app.Get("/api/users/me", middleware.JWTAuth(), userHandler.GetProfile)
	app.Put("/api/users/me", middleware.JWTAuth(), userHandler.UpdateProfile)
	app.Get("/api/users/me/activity", middleware.JWTAuth(), userHandler.GetActivityLog)

	// Public user profile (no auth required, placed before protected group to bypass JWTAuth)
	app.Get("/api/users/:id", userHandler.GetPublicProfile)

	// ==================== PROTECTED ROUTES ====================

	api := app.Group("/api", middleware.JWTAuth())

	// Showcase endpoints
	api.Post("/showcases", showcaseHandler.Create)
	api.Put("/showcases/:id", showcaseHandler.Update)
	api.Delete("/showcases/:id", showcaseHandler.Delete)
	api.Post("/showcases/:id/like", showcaseHandler.Like)
	api.Delete("/showcases/:id/like", showcaseHandler.Unlike)

	// Certificate endpoints
	certificateHandler := handler.NewCertificateHandler()
	api.Get("/certificates", certificateHandler.ListCertificates)
	api.Get("/certificates/:code", certificateHandler.GetCertificate)

	// Discussion endpoints
	discussionHandler := handler.NewDiscussionHandler()
	api.Post("/discussions", discussionHandler.CreateDiscussion)
	api.Get("/discussions", discussionHandler.ListDiscussions)
	api.Get("/discussions/:id", discussionHandler.GetDiscussion)
	api.Put("/discussions/:id", discussionHandler.UpdateDiscussion)

	// Reply endpoints
	replyHandler := handler.NewReplyHandler()
	api.Post("/discussions/:id/replies", replyHandler.CreateReply)
	api.Put("/replies/:id", replyHandler.UpdateReply)
	api.Delete("/replies/:id", replyHandler.DeleteReply)

	// Gamification endpoints
	gamificationHandler := handler.NewGamificationHandler()
	api.Get("/levels", gamificationHandler.GetLevels)
	api.Get("/users/:id/points", gamificationHandler.GetUserPoints)

	// Class endpoints
	classHandler := handler.NewClassHandler(service.GetClassService())
	api.Get("/projects/:project_id/classes/:slug", classHandler.GetClassBySlug)
	api.Post("/classes/:id/start", classHandler.StartClass)
	api.Put("/classes/:id/progress", classHandler.UpdateProgress)
	api.Post("/classes/:id/complete", classHandler.Complete)

	// ==================== ADMIN ROUTES ====================
	admin := api.Group("/admin", middleware.RequireRole("admin"))

	admin.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome admin"})
	})

	// Admin User endpoints
	admin.Get("/users", userHandler.GetAllUsers)

	// Admin Project endpoints
	projectHandler := handler.NewProjectHandler(service.GetProjectService())
	admin.Post("/projects", projectHandler.CreateProject)
	admin.Get("/projects", projectHandler.ListProjects)
	admin.Put("/projects/:id", projectHandler.UpdateProject)
	admin.Delete("/projects/:id", projectHandler.DeleteProject)

	// Admin Class endpoints
	admin.Post("/classes", classHandler.CreateClass)
	admin.Post("/classes/:id/details", classHandler.CreateClassDetail)
	admin.Put("/classes/:id", classHandler.UpdateClass)
	admin.Delete("/classes/:id", classHandler.DeleteClass)

	// Admin Category endpoints
	admin.Post("/categories", categoryHandler.CreateCategory)
	admin.Get("/categories", categoryHandler.ListCategories)
	admin.Put("/categories/:id", categoryHandler.UpdateCategory)
	admin.Delete("/categories/:id", categoryHandler.DeleteCategory)
}
