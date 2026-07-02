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

	// --- Public System Status (no auth) ---
	statusHandler := handler.NewStatusHandler()
	app.Get("/api/system/status", statusHandler.SystemStatus)

	// --- Public Company Profile (no auth) ---
	companyHandler := handler.NewCompanyHandler(service.GetCompanyService())
	app.Get("/api/company", companyHandler.GetCompany)

	// --- Public FAQs (no auth) ---
	faqHandler := handler.NewFaqHandler(service.GetFaqService())
	app.Get("/api/faqs", faqHandler.ListPublic)

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
	// ==================== PROTECTED ROUTES (JWT only — safe, no XSRF) ====================

	userHandler := handler.NewUserHandler(service.GetUserService(), service.GetDashboardService())
	app.Get("/api/mentors", userHandler.ListMentors)
	app.Get("/api/users/me", middleware.JWTAuth(), userHandler.GetProfile)
	app.Put("/api/users/me", middleware.JWTAuth(), userHandler.UpdateProfile)
	app.Post("/api/users/me/change-password", middleware.JWTAuth(), userHandler.ChangePassword)
	app.Post("/api/users/me/avatar", middleware.JWTAuth(), userHandler.UpdateProfilePicture)
	app.Get("/api/users/me/activity", middleware.JWTAuth(), userHandler.GetActivityLog)
	app.Get("/api/users/me/activity-history", middleware.JWTAuth(), userHandler.GetActivityHistory)
	app.Get("/api/users/me/dashboard", middleware.JWTAuth(), userHandler.GetDashboard)
	app.Get("/api/users/:id", userHandler.GetPublicProfile)

	// --- Public User Portfolio (no auth) ---
	app.Get("/api/users/:id/portfolio", userHandler.GetPortfolio)

	// Safe group — JWT only (no XSRF) for non-dangerous operations
	safe := app.Group("/api", middleware.JWTAuth(), middleware.IdempotencyMiddleware())

	// Certificates
	certificateHandler := handler.NewCertificateHandler()
	safe.Get("/users/me/certificates", certificateHandler.ListCertificates)
	safe.Get("/users/me/certificates/:code", certificateHandler.GetCertificate)

	// --- Public Certificate Verification (no auth) ---
	app.Get("/api/certificates/:code/verify", certificateHandler.VerifyCertificate)

	// Enrollment & Courses
	safe.Post("/courses/:id/enroll", courseHandler.EnrollCourse)
	safe.Put("/courses/:id/last-lesson", courseHandler.SetLastLesson)
	safe.Get("/users/me/courses", courseHandler.ListEnrolledCourses)

	// Learning Progress
	lessonProgressHandler := handler.NewLessonHandler(service.GetLessonService())
	safe.Post("/lessons/:id/start", lessonProgressHandler.StartLesson)
	safe.Put("/lessons/:id/progress", lessonProgressHandler.UpdateProgress)
	safe.Post("/lessons/:id/complete", lessonProgressHandler.CompleteLesson)

	// Notifications
	notifHandler := handler.NewNotificationHandler(service.GetNotificationService())
	sseHandler := handler.NewSSEHandler()
	safe.Get("/notifications/stream", sseHandler.StreamNotifications)
	safe.Get("/users/me/notifications", notifHandler.ListNotifications)
	safe.Get("/users/me/notifications/count", notifHandler.CountUnread)
	safe.Put("/users/me/notifications/:id/read", notifHandler.MarkAsRead)
	safe.Put("/users/me/notifications/read-all", notifHandler.MarkAllAsRead)
	safe.Delete("/users/me/notifications/:id", notifHandler.DeleteNotification)

	// My Learning Streak
	safe.Get("/users/me/streak", userHandler.GetMyStreak)

	// My Items (read)
	safe.Get("/users/me/discussions", discussionHandler.ListMyDiscussions)
	safe.Get("/users/me/replies", replyHandler.GetMyReplies)
	safe.Get("/users/me/study-cases", studyCaseHandler.ListMyStudyCases)
	safe.Get("/users/me/showcases", showcaseHandler.ListMyShowcases)
	safe.Get("/users/me/blogs", blogHandler.ListMyBlogs)

	// Gamification
	gamificationHandler := handler.NewGamificationHandler()
	safe.Get("/levels", gamificationHandler.GetLevels)
	safe.Get("/users/:id/points", gamificationHandler.GetUserPoints)

	// File Upload
	uploadHandler := handler.NewUploadHandler()
	safe.Post("/upload", uploadHandler.Upload)

	// ==================== DANGEROUS ROUTES (JWT + XSRF + Idempotency) ====================
	// Only content-modifying operations that need XSRF protection

	dangerous := app.Group("/api", middleware.JWTAuth(), middleware.XSRFProtection(), middleware.IdempotencyMiddleware())

	dangerous.Post("/showcases", showcaseHandler.Create)
	dangerous.Put("/showcases/:id", showcaseHandler.Update)
	dangerous.Delete("/showcases/:id", showcaseHandler.Delete)
	dangerous.Post("/showcases/:id/like", showcaseHandler.Like)
	dangerous.Delete("/showcases/:id/like", showcaseHandler.Unlike)

	dangerous.Post("/discussions", discussionHandler.CreateDiscussion)
	dangerous.Put("/discussions/:id", discussionHandler.UpdateDiscussion)
	dangerous.Delete("/discussions/:id", discussionHandler.DeleteDiscussion)
	dangerous.Post("/discussions/:id/close", discussionHandler.CloseDiscussion)

	replyHandler = handler.NewReplyHandler()
	dangerous.Post("/discussions/:id/replies", replyHandler.CreateReply)
	dangerous.Put("/replies/:id", replyHandler.UpdateReply)
	dangerous.Delete("/replies/:id", replyHandler.DeleteReply)
	dangerous.Post("/replies/:id/like", replyHandler.LikeReply)
	dangerous.Post("/replies/:id/react", replyHandler.ReactReply)
	dangerous.Delete("/replies/:id/react/:emoji", replyHandler.UnreactReply)
	dangerous.Post("/replies/:id/best", replyHandler.MarkBestReply)

	dangerous.Post("/reviews", reviewHandler.CreateReview)
	dangerous.Put("/reviews/:id", reviewHandler.UpdateReview)
	dangerous.Delete("/reviews/:id", reviewHandler.DeleteReview)

	// ==================== ADMIN ROUTES (JWT + XSRF + Idempotency + admin role) ====================

	admin := app.Group("/api/admin", middleware.JWTAuth(), middleware.XSRFProtection(), middleware.IdempotencyMiddleware(), middleware.RequireRole("admin"))
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

	// --- Admin Company Profile ---
	admin.Put("/company", companyHandler.UpdateCompany)

	// --- Admin FAQ CRUD ---
	admin.Get("/faqs", faqHandler.ListAll)
	admin.Get("/faqs/:id", faqHandler.GetByID)
	admin.Post("/faqs", faqHandler.Create)
	admin.Put("/faqs/:id", faqHandler.Update)
	admin.Delete("/faqs/:id", faqHandler.Delete)
}
