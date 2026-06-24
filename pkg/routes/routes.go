package routes

import (
	"jvalleyverse/internal/handler"
	"jvalleyverse/internal/service"
	"jvalleyverse/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	handler.InitAuditService(service.GetAuditService())

	// ==================== PUBLIC ROUTES ====================

	authHandler := handler.NewAuthHandler(service.GetUserService())
	app.Post("/api/auth/register", middleware.AuthRateLimiter(), middleware.IdempotencyMiddleware(), authHandler.Register)
	app.Post("/api/auth/login", middleware.AuthRateLimiter(), middleware.IdempotencyMiddleware(), authHandler.Login)
	app.Post("/api/auth/google", middleware.AuthRateLimiter(), middleware.IdempotencyMiddleware(), authHandler.GoogleLogin)
	app.Post("/api/auth/refresh", middleware.IdempotencyMiddleware(), authHandler.Refresh)
	app.Post("/api/auth/logout", middleware.JWTAuth(), authHandler.Logout)

	// ── PUBLIC ROUTES (with anti-scraping guard on content) ──

	showcaseHandler := handler.NewShowcaseHandler(service.GetShowcaseService())
	app.Get("/api/leaderboard", handler.NewGamificationHandler().GetLeaderboard)

	contentGuard := middleware.ScraperGuard()
	contentRateLimit := middleware.ContentRateLimiter()

	app.Get("/api/showcases", contentGuard, contentRateLimit, showcaseHandler.ListShowcases)
	app.Get("/api/showcases/:id", contentGuard, contentRateLimit, showcaseHandler.GetShowcase)

	categoryHandler := handler.NewCategoryHandler(service.GetCategoryService())
	app.Get("/api/categories", contentGuard, contentRateLimit, categoryHandler.ListCategories)
	app.Get("/api/categories/:slug", contentGuard, contentRateLimit, categoryHandler.GetCategoryBySlug)
	app.Get("/api/categories/:category_id/courses", contentGuard, contentRateLimit, middleware.OptionalJWTAuth(), categoryHandler.ListCoursesByCategory)

	sectionHandler := handler.NewSectionHandler(service.GetSectionService())
	courseHandler := handler.NewCourseHandler(service.GetCourseService())

	app.Get("/api/courses", contentGuard, contentRateLimit, middleware.OptionalJWTAuth(), courseHandler.ListPublicCourses)
	app.Get("/api/courses/:course_id", contentGuard, contentRateLimit, middleware.OptionalJWTAuth(), sectionHandler.GetCourseWithSections)
	app.Get("/api/courses/:course_id/sections", contentGuard, contentRateLimit, sectionHandler.ListSectionsByCourse)
	app.Get("/api/courses/:course_id/sections/:section_id", contentGuard, contentRateLimit, sectionHandler.GetSection)

	reviewHandler := handler.NewReviewHandler(service.GetReviewService())
	app.Get("/api/courses/:course_id/reviews", contentGuard, contentRateLimit, reviewHandler.ListCourseReviews)
	app.Get("/api/lessons/:id/reviews", contentGuard, contentRateLimit, reviewHandler.ListLessonReviews)

	lessonHandler := handler.NewLessonHandler(service.GetLessonService())
	app.Get("/api/lessons/:id", contentGuard, contentRateLimit, lessonHandler.GetPublicLessonByID)
	app.Get("/api/courses/:course_id/lessons", contentGuard, contentRateLimit, lessonHandler.ListLessonsByCourse)
	app.Get("/api/courses/:course_id/sections/:section_id/lessons", contentGuard, contentRateLimit, lessonHandler.ListLessonsBySection)
	app.Get("/api/courses/:course_id/lessons/:slug", contentGuard, contentRateLimit, lessonHandler.GetLessonBySlug)

	studyCaseHandler := handler.NewStudyCaseHandler(service.GetStudyCaseService())
	app.Get("/api/study-cases", contentGuard, contentRateLimit, studyCaseHandler.ListStudyCases)
	app.Get("/api/study-cases/:id", contentGuard, contentRateLimit, studyCaseHandler.GetStudyCase)

	healthHandler := handler.NewHealthHandler()
	app.Get("/api/health", healthHandler.Health)
	app.Get("/api/health/detailed", healthHandler.HealthDetailed)

	// --- Blogs ---
	blogHandler := handler.NewBlogHandler(service.GetBlogService())

	// Public blog routes
	blogs := app.Group("/api/blogs", contentGuard, contentRateLimit)
	blogs.Get("/", blogHandler.List)
	blogs.Get("/:id", blogHandler.GetByID)

	// --- Public Discussions (read-only, optional auth) ---
	discussionHandler := handler.NewDiscussionHandler()
	replyHandler := handler.NewReplyHandler()

	discussions := app.Group("/api/discussions", contentGuard, contentRateLimit, middleware.OptionalJWTAuth())
	discussions.Get("/", discussionHandler.ListDiscussions)
	discussions.Get("/:id", discussionHandler.GetDiscussion)
	// ==================== PROTECTED ROUTES ====================

	userHandler := handler.NewUserHandler(service.GetUserService(), service.GetDashboardService())
	app.Get("/api/mentors", userHandler.ListMentors)
	app.Get("/api/users/me", middleware.JWTAuth(), userHandler.GetProfile)
	app.Put("/api/users/me", middleware.JWTAuth(), userHandler.UpdateProfile)
	app.Get("/api/users/me/activity", middleware.JWTAuth(), userHandler.GetActivityLog)
	app.Get("/api/users/me/dashboard", middleware.JWTAuth(), userHandler.GetDashboard)
	app.Get("/api/users/:id", userHandler.GetPublicProfile)

	api := app.Group("/api", middleware.JWTAuth(), middleware.XSRFProtection(), middleware.IdempotencyMiddleware())
	// --- My Blogs ---
	api.Get("/users/me/blogs", blogHandler.ListMyBlogs)

	api.Post("/showcases", showcaseHandler.Create)
	api.Put("/showcases/:id", showcaseHandler.Update)
	api.Delete("/showcases/:id", showcaseHandler.Delete)
	api.Post("/showcases/:id/like", showcaseHandler.Like)
	api.Delete("/showcases/:id/like", showcaseHandler.Unlike)
	api.Get("/users/me/showcases", showcaseHandler.ListMyShowcases)

	certificateHandler := handler.NewCertificateHandler()
	api.Get("/users/me/certificates", certificateHandler.ListCertificates)
	api.Get("/users/me/certificates/:code", certificateHandler.GetCertificate)
	api.Get("/users/me/discussions", discussionHandler.ListMyDiscussions)
	api.Post("/discussions", discussionHandler.CreateDiscussion)
	api.Put("/discussions/:id", discussionHandler.UpdateDiscussion)
	api.Delete("/discussions/:id", discussionHandler.DeleteDiscussion)
	api.Post("/discussions/:id/close", discussionHandler.CloseDiscussion)
	replyHandler = handler.NewReplyHandler()
	api.Post("/discussions/:id/replies", replyHandler.CreateReply)
	api.Put("/replies/:id", replyHandler.UpdateReply)
	api.Delete("/replies/:id", replyHandler.DeleteReply)
	api.Post("/replies/:id/like", replyHandler.LikeReply)
	api.Post("/replies/:id/best", replyHandler.MarkBestReply)
	api.Get("/users/me/replies", replyHandler.GetMyReplies)
	api.Get("/users/me/study-cases", studyCaseHandler.ListMyStudyCases)

	api.Post("/courses/:id/enroll", courseHandler.EnrollCourse)
	api.Put("/courses/:id/last-lesson", courseHandler.SetLastLesson)
	api.Get("/users/me/courses", courseHandler.ListEnrolledCourses)

	// Notification routes
	notifHandler := handler.NewNotificationHandler(service.GetNotificationService())
	sseHandler := handler.NewSSEHandler()
	api.Get("/notifications/stream", sseHandler.StreamNotifications)
	api.Get("/users/me/notifications", notifHandler.ListNotifications)
	api.Get("/users/me/notifications/count", notifHandler.CountUnread)
	api.Put("/users/me/notifications/:id/read", notifHandler.MarkAsRead)
	api.Put("/users/me/notifications/read-all", notifHandler.MarkAllAsRead)
	api.Delete("/users/me/notifications/:id", notifHandler.DeleteNotification)

	gamificationHandler := handler.NewGamificationHandler()
	api.Get("/levels", gamificationHandler.GetLevels)
	api.Get("/users/:id/points", gamificationHandler.GetUserPoints)

	api.Post("/reviews", reviewHandler.CreateReview)
	api.Put("/reviews/:id", reviewHandler.UpdateReview)

	api.Delete("/reviews/:id", reviewHandler.DeleteReview)

	lessonProgressHandler := handler.NewLessonHandler(service.GetLessonService())
	api.Post("/lessons/:id/start", lessonProgressHandler.StartLesson)
	api.Put("/lessons/:id/progress", lessonProgressHandler.UpdateProgress)
	api.Post("/lessons/:id/complete", lessonProgressHandler.CompleteLesson)

	// ==================== ADMIN ROUTES ====================

	uploadHandler := handler.NewUploadHandler()
	api.Post("/upload", uploadHandler.Upload)

	admin := api.Group("/admin", middleware.RequireRole("admin"))
	// --- Admin Blogs (CRUD) ---
	blogAdmin := admin.Group("/blogs")
	blogAdmin.Post("/", blogHandler.Create)
	blogAdmin.Put("/:id", blogHandler.AdminUpdate)
	blogAdmin.Delete("/:id", blogHandler.AdminDelete)
	blogAdmin.Get("/", blogHandler.List)

	admin.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome admin"})
	})

	admin.Get("/users", userHandler.GetAllUsers)

	admin.Post("/courses", courseHandler.CreateCourse)
	admin.Put("/courses/:id", courseHandler.UpdateCourse)
	admin.Delete("/courses/:id", courseHandler.DeleteCourse)

	admin.Post("/courses/:course_id/sections", sectionHandler.CreateSection)
	admin.Put("/sections/:section_id", sectionHandler.UpdateSection)
	admin.Delete("/sections/:section_id", sectionHandler.DeleteSection)

	admin.Post("/lessons", lessonHandler.CreateLesson)
	admin.Post("/lessons/:id/details", lessonHandler.CreateLessonDetail)
	admin.Put("/lessons/:id", lessonHandler.UpdateLesson)
	admin.Delete("/lessons/:id", lessonHandler.DeleteLesson)

	admin.Post("/study-cases", studyCaseHandler.CreateStudyCase)
	admin.Put("/study-cases/:id", studyCaseHandler.UpdateStudyCase)
	admin.Delete("/study-cases/:id", studyCaseHandler.DeleteStudyCase)

	admin.Post("/categories", categoryHandler.CreateCategory)
	admin.Get("/categories", categoryHandler.ListCategories)
	admin.Put("/categories/:id", categoryHandler.UpdateCategory)
	admin.Delete("/categories/:id", categoryHandler.DeleteCategory)

}
