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

	// Showcase & leaderboard (public read)
	showcaseHandler := handler.NewShowcaseHandler(service.GetShowcaseService())
	app.Get("/api/leaderboard", handler.NewGamificationHandler().GetLeaderboard)
	app.Get("/api/showcases", showcaseHandler.ListShowcases)
	app.Get("/api/showcases/:id", showcaseHandler.GetShowcase)

	// Categories (public read)
	categoryHandler := handler.NewCategoryHandler(service.GetCategoryService())
	app.Get("/api/categories", categoryHandler.ListCategories)
	app.Get("/api/categories/:slug", categoryHandler.GetCategoryBySlug)
	app.Get("/api/categories/:category_id/projects", categoryHandler.ListProjectsByCategory)

	// Projects (public read) — returns project with all phases + classes
	phaseHandler := handler.NewPhaseHandler(service.GetPhaseService())
	projectHandler := handler.NewProjectHandler(service.GetProjectService())

	app.Get("/api/projects", projectHandler.ListPublicProjects)
	app.Get("/api/projects/:project_id", phaseHandler.GetProjectWithPhases)

	// Phases (public read)
	app.Get("/api/projects/:project_id/phases", phaseHandler.ListPhasesByProject)
	app.Get("/api/projects/:project_id/phases/:phase_id", phaseHandler.GetPhase)

	// Classes within a phase (public read)
	publicClassHandler := handler.NewClassHandler(service.GetClassService())

	app.Get("/api/classes/:id", publicClassHandler.GetPublicClassByID)
	app.Get("/api/projects/:project_id/classes", publicClassHandler.ListClassesByProject)
	app.Get("/api/projects/:project_id/phases/:phase_id/classes", publicClassHandler.ListClassesByPhase)
	app.Get("/api/projects/:project_id/classes/:slug", publicClassHandler.GetClassBySlug)

	// Health check
	healthHandler := handler.NewHealthHandler()
	app.Get("/api/health", healthHandler.Health)
	app.Get("/api/health/detailed", healthHandler.HealthDetailed)

	// ==================== PROTECTED ROUTES (individual, pre-group) ====================

	userHandler := handler.NewUserHandler(service.GetUserService())
	app.Get("/api/users/me", middleware.JWTAuth(), userHandler.GetProfile)
	app.Put("/api/users/me", middleware.JWTAuth(), userHandler.UpdateProfile)
	app.Get("/api/users/me/activity", middleware.JWTAuth(), userHandler.GetActivityLog)
	app.Get("/api/users/:id", userHandler.GetPublicProfile)

	// ==================== PROTECTED ROUTES (JWT group) ====================

	api := app.Group("/api", middleware.JWTAuth())

	// Showcase (authenticated actions)
	api.Post("/showcases", showcaseHandler.Create)
	api.Put("/showcases/:id", showcaseHandler.Update)
	api.Delete("/showcases/:id", showcaseHandler.Delete)
	api.Post("/showcases/:id/like", showcaseHandler.Like)
	api.Delete("/showcases/:id/like", showcaseHandler.Unlike)

	// Certificates
	certificateHandler := handler.NewCertificateHandler()
	api.Get("/certificates", certificateHandler.ListCertificates)
	api.Get("/certificates/:code", certificateHandler.GetCertificate)

	// Discussions
	discussionHandler := handler.NewDiscussionHandler()
	api.Post("/discussions", discussionHandler.CreateDiscussion)
	api.Get("/discussions", discussionHandler.ListDiscussions)
	api.Get("/discussions/:id", discussionHandler.GetDiscussion)
	api.Put("/discussions/:id", discussionHandler.UpdateDiscussion)

	// Replies
	replyHandler := handler.NewReplyHandler()
	api.Post("/discussions/:id/replies", replyHandler.CreateReply)
	api.Put("/replies/:id", replyHandler.UpdateReply)
	api.Delete("/replies/:id", replyHandler.DeleteReply)

	// Gamification
	gamificationHandler := handler.NewGamificationHandler()
	api.Get("/levels", gamificationHandler.GetLevels)
	api.Get("/users/:id/points", gamificationHandler.GetUserPoints)

	// Class progress (user actions on a class — uses class_id directly)
	classHandler := handler.NewClassHandler(service.GetClassService())
	api.Post("/classes/:id/start", classHandler.StartClass)
	api.Put("/classes/:id/progress", classHandler.UpdateProgress)
	api.Post("/classes/:id/complete", classHandler.Complete)

	// ==================== ADMIN ROUTES ====================

	admin := api.Group("/admin", middleware.RequireRole("admin"))

	admin.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome admin"})
	})

	// Admin: Users
	admin.Get("/users", userHandler.GetAllUsers)

	// Admin: Projects
	admin.Post("/projects", projectHandler.CreateProject)
	admin.Get("/projects", projectHandler.ListProjects)
	admin.Put("/projects/:id", projectHandler.UpdateProject)
	admin.Delete("/projects/:id", projectHandler.DeleteProject)

	// Admin: Phases (nested under project for create, standalone for update/delete)
	admin.Post("/projects/:project_id/phases", phaseHandler.CreatePhase)
	admin.Put("/phases/:phase_id", phaseHandler.UpdatePhase)
	admin.Delete("/phases/:phase_id", phaseHandler.DeletePhase)

	// Admin: Classes (class belongs to a phase — phase_id required in body on create)
	admin.Post("/classes", classHandler.CreateClass)
	admin.Post("/classes/:id/details", classHandler.CreateClassDetail)
	admin.Put("/classes/:id", classHandler.UpdateClass)
	admin.Delete("/classes/:id", classHandler.DeleteClass)

	// Admin: Categories
	admin.Post("/categories", categoryHandler.CreateCategory)
	admin.Get("/categories", categoryHandler.ListCategories)
	admin.Put("/categories/:id", categoryHandler.UpdateCategory)
	admin.Delete("/categories/:id", categoryHandler.DeleteCategory)
}
