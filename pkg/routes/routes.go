package routes

import (
	"jvalleyverse/internal/handler"
	"jvalleyverse/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App) {
	// ==================== PUBLIC ROUTES ====================
	// Auth endpoints
	authHandler := handler.NewAuthHandler()
	app.Post("/api/auth/register", authHandler.Register)
	app.Post("/api/auth/login", authHandler.Login)
	// For backward compatibility
	app.Post("/auth/register", authHandler.Register)
	app.Post("/auth/login", authHandler.Login)

	// List public content
	showcaseHandler := handler.NewShowcaseHandler()
	app.Get("/api/leaderboard", showcaseHandler.GetLeaderboard)
	app.Get("/api/showcases", showcaseHandler.ListShowcases)
	app.Get("/api/showcases/:id", showcaseHandler.GetShowcase)

	// ==================== PUBLIC ROUTES (no auth) ====================
	// Publicly accessible user profile and class details
	publicUserHandler := handler.NewUserHandler()
	app.Get("/api/users/:id", publicUserHandler.GetPublicProfile)
	publicClassHandler := handler.NewClassHandler()
	app.Get("/api/projects/:project_id/classes/:slug", publicClassHandler.GetClassBySlug)

	// ==================== PROTECTED ROUTES ====================

	api := app.Group("/api", middleware.JWTAuth())

	// User endpoints
	userHandler := handler.NewUserHandler()
	api.Get("/users/me", userHandler.GetProfile)
	api.Put("/users/me", userHandler.UpdateProfile)
	api.Get("/users/:id", userHandler.GetPublicProfile)
	api.Get("/users/me/activity", userHandler.GetActivityLog)

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
	classHandler := handler.NewClassHandler()
	api.Get("/projects/:project_id/classes/:slug", classHandler.GetClassBySlug)
	api.Post("/classes/:id/start", classHandler.StartClass)
	api.Put("/classes/:id/progress", classHandler.UpdateProgress)
	api.Post("/classes/:id/complete", classHandler.Complete)

	// ==================== ADMIN ROUTES ====================
	admin := api.Group("/admin", middleware.RequireRole("admin"))

	admin.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome admin"})
	})

	// Admin Project endpoints
	projectHandler := handler.NewProjectHandler()
	admin.Post("/projects", projectHandler.CreateProject)
	admin.Get("/projects", projectHandler.ListProjects)
	admin.Put("/projects/:id", projectHandler.UpdateProject)
	admin.Delete("/projects/:id", projectHandler.DeleteProject)

	// Admin Class endpoints
	admin.Post("/classes", classHandler.CreateClass)
	admin.Post("/classes/:id/details", classHandler.CreateClassDetail)
	admin.Put("/classes/:id", classHandler.UpdateClass)
	admin.Delete("/classes/:id", classHandler.DeleteClass)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
}
