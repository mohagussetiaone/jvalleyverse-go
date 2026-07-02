package handler

import (
	"context"
	"jvalleyverse/internal/dto"
	"jvalleyverse/internal/minio"
	"jvalleyverse/internal/service"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userSvc       service.IUserService
	dashboardSvc  service.IDashboardService
	certificateSvc service.ICertificateService
	showcaseSvc   service.IShowcaseService
	studyCaseSvc  service.IStudyCaseService
	streakSvc     *service.StreakService
}

func NewUserHandler(userSvc service.IUserService, dashboardSvc service.IDashboardService) *UserHandler {
	return &UserHandler{
		userSvc:       userSvc,
		dashboardSvc:  dashboardSvc,
		certificateSvc: service.GetCertificateService(),
		showcaseSvc:   service.GetShowcaseService(),
		studyCaseSvc:  service.GetStudyCaseService(),
		streakSvc:     service.GetStreakService(),
	}
}

// ── Activity history helpers ──

func (h *UserHandler) fetchPointActivity(ctx context.Context, userID string, page, limit int) []dto.ActivityHistoryItem {
	items, _, err := h.userSvc.GetUserActivityLog(ctx, userID, page, limit)
	if err != nil || len(items) == 0 {
		return nil
	}
	result := make([]dto.ActivityHistoryItem, len(items))
	for i, item := range items {
		result[i] = dto.ActivityHistoryItem{
			ID:          item.ID,
			Type:        "point",
			Title:       item.Activity,
			Description: item.Activity,
			CreatedAt:   item.Timestamp,
		}
	}
	return result
}

func (h *UserHandler) fetchDiscussionActivity(ctx context.Context, userID string, page, limit int) []dto.ActivityHistoryItem {
	discussions, _, err := service.GetDiscussionService().ListUserDiscussions(ctx, userID, page, limit)
	if err != nil || len(discussions) == 0 {
		return nil
	}
	result := make([]dto.ActivityHistoryItem, len(discussions))
	for i, d := range discussions {
		result[i] = dto.ActivityHistoryItem{
			ID:          d.ID,
			Type:        "discussion",
			Title:       d.Title,
			Description: "Membuat diskusi",
			Link:        "/discussions/" + d.ID,
			CreatedAt:   d.CreatedAt,
		}
	}
	return result
}

func (h *UserHandler) fetchReplyActivity(ctx context.Context, userID string, page, limit int) []dto.ActivityHistoryItem {
	replies, _, err := service.GetReplyService().ListRepliesByUser(ctx, userID, page, limit)
	if err != nil || len(replies) == 0 {
		return nil
	}
	result := make([]dto.ActivityHistoryItem, len(replies))
	for i, r := range replies {
		result[i] = dto.ActivityHistoryItem{
			ID:          r.ID,
			Type:        "reply",
			Title:       r.DiscussionTitle,
			Description: r.Content,
			Link:        "/discussions/" + r.DiscussionID,
			CreatedAt:   r.CreatedAt,
		}
	}
	return result
}

func (h *UserHandler) fetchCertificateActivity(ctx context.Context, userID string, page, limit int) []dto.ActivityHistoryItem {
	certs, _, err := h.certificateSvc.ListUserCertificates(ctx, userID, page, limit)
	if err != nil || len(certs) == 0 {
		return nil
	}
	result := make([]dto.ActivityHistoryItem, len(certs))
	for i, c := range certs {
		result[i] = dto.ActivityHistoryItem{
			ID:          c.ID,
			Type:        "certificate",
			Title:       c.LessonName,
			Description: "Sertifikat diperoleh",
			Link:        "/certificates/" + c.UniqueCode + "/verify",
			CreatedAt:   c.IssuedAt,
		}
	}
	return result
}

func (h *UserHandler) fetchShowcaseActivity(ctx context.Context, userID string, page, limit int) []dto.ActivityHistoryItem {
	showcases, _, err := h.showcaseSvc.ListMyShowcases(ctx, userID, page, limit)
	if err != nil || len(showcases) == 0 {
		return nil
	}
	result := make([]dto.ActivityHistoryItem, len(showcases))
	for i, s := range showcases {
		result[i] = dto.ActivityHistoryItem{
			ID:          s.ID,
			Type:        "showcase",
			Title:       s.Title,
			Description: s.Description,
			Link:        "/showcases/" + s.ID,
			CreatedAt:   s.CreatedAt,
		}
	}
	return result
}

func (h *UserHandler) fetchStudyCaseActivity(ctx context.Context, userID string, page, limit int) []dto.ActivityHistoryItem {
	studyCases, _, err := h.studyCaseSvc.ListStudyCasesByUser(ctx, userID, page, limit)
	if err != nil || len(studyCases) == 0 {
		return nil
	}
	result := make([]dto.ActivityHistoryItem, len(studyCases))
	for i, sc := range studyCases {
		result[i] = dto.ActivityHistoryItem{
			ID:          sc.ID,
			Type:        "study_case",
			Title:       sc.Name,
			Description: sc.Description,
			Link:        "/study-cases/" + sc.ID,
			CreatedAt:   sc.CreatedAt,
		}
	}
	return result
}

func (h *UserHandler) fetchCourseActivity(ctx context.Context, userID string, page, limit int) []dto.ActivityHistoryItem {
	courses, _, err := service.GetCourseService().ListEnrolledCourses(ctx, userID, page, limit)
	if err != nil || len(courses) == 0 {
		return nil
	}
	result := make([]dto.ActivityHistoryItem, len(courses))
	for i, c := range courses {
		result[i] = dto.ActivityHistoryItem{
			ID:          c.ID,
			Type:        "course_enrollment",
			Title:       c.Title,
			Description: "Mendaftar kursus",
			Link:        "/courses/" + c.ID,
			CreatedAt:   c.EnrolledAt,
		}
	}
	return result
}

// GetProfile returns current user profile
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	user, err := h.userSvc.GetProfile(c.UserContext(), userID)
	if err != nil {
		return safeError(c, mapServiceErrorToStatus(err), err)
	}

	return c.JSON(user)
}

// UpdateProfile updates current user profile
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var input struct {
		Name   string `json:"name"`
		Bio    string `json:"bio"`
		Avatar string `json:"avatar"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.userSvc.UpdateProfile(c.UserContext(), userID, input.Name, input.Bio, input.Avatar); err != nil {
		return safeError(c, mapServiceErrorToStatus(err), err)
	}

	return c.JSON(fiber.Map{"message": "Profile updated"})
}

// ChangePassword changes current user password (must be logged in)
func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var input struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if input.NewPassword == "" {
		return c.Status(400).JSON(fiber.Map{"error": "new_password is required"})
	}

	if err := h.userSvc.ChangePassword(c.UserContext(), userID, input.CurrentPassword, input.NewPassword); err != nil {
		return safeError(c, mapServiceErrorToStatus(err), err)
	}

	return c.JSON(fiber.Map{"message": "Password changed successfully"})
}

// UpdateProfilePicture uploads avatar to MinIO and saves URL to DB
func (h *UserHandler) UpdateProfilePicture(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	if !minio.IsAvailable() {
		return c.Status(503).JSON(fiber.Map{
			"error": "Avatar upload is not available (MinIO not configured)",
		})
	}

	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "avatar file is required"})
	}

	// Validate file size (max 10MB)
	if fileHeader.Size > 10<<20 {
		return c.Status(400).JSON(fiber.Map{"error": "file too large, maximum size is 10 MB"})
	}

	// Upload to MinIO
	result, err := minio.DefaultClient.UploadFile(c.UserContext(), fileHeader, "avatars")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "file upload failed"})
	}

	// Update user avatar in DB
	if err := h.userSvc.UpdateProfile(c.UserContext(), userID, "", "", result.URL); err != nil {
		return safeError(c, mapServiceErrorToStatus(err), err)
	}

	return c.JSON(fiber.Map{
		"message":     "Avatar updated",
		"url":         result.URL,
		"object_name": result.ObjectName,
	})
}

// GetDashboard returns dashboard widgets and stats
func (h *UserHandler) GetDashboard(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	dashboard, err := h.dashboardSvc.GetDashboard(c.UserContext(), userID)
	if err != nil {
		return safeError(c, 500, err)
	}

	return c.JSON(dashboard)
}

// GetPublicProfile returns public user profile
func (h *UserHandler) GetPublicProfile(c *fiber.Ctx) error {
	userID := c.Params("id") // Direct string CUID
	if userID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	user, err := h.userSvc.GetProfile(c.UserContext(), userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Filter public info only or use a DTO
	return c.JSON(fiber.Map{
		"id":     user.ID,
		"name":   user.Name,
		"avatar": user.Avatar,
		"level":  user.Level,
		"points": user.Points,
	})
}

// GetActivityHistory returns unified activity history with optional type filter
func (h *UserHandler) GetActivityHistory(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	page := clampPage(c.QueryInt("page", DefaultPage))
	limit := clampLimit(c.QueryInt("limit", DefaultLimit), DefaultLimit)
	activityType := c.Query("type", "")

	var all []dto.ActivityHistoryItem

	switch activityType {
	case "", "all":
		// Fetch all types
		all = append(all, h.fetchPointActivity(c.UserContext(), userID, page, limit)...)
		all = append(all, h.fetchDiscussionActivity(c.UserContext(), userID, page, limit)...)
		all = append(all, h.fetchReplyActivity(c.UserContext(), userID, page, limit)...)
		all = append(all, h.fetchCertificateActivity(c.UserContext(), userID, page, limit)...)
		all = append(all, h.fetchShowcaseActivity(c.UserContext(), userID, page, limit)...)
		all = append(all, h.fetchStudyCaseActivity(c.UserContext(), userID, page, limit)...)
		all = append(all, h.fetchCourseActivity(c.UserContext(), userID, page, limit)...)

	case "discussion":
		all = h.fetchDiscussionActivity(c.UserContext(), userID, page, limit)
	case "reply", "comment":
		all = h.fetchReplyActivity(c.UserContext(), userID, page, limit)
	case "certificate":
		all = h.fetchCertificateActivity(c.UserContext(), userID, page, limit)
	case "showcase":
		all = h.fetchShowcaseActivity(c.UserContext(), userID, page, limit)
	case "study_case", "studycase":
		all = h.fetchStudyCaseActivity(c.UserContext(), userID, page, limit)
	case "course", "enrollment":
		all = h.fetchCourseActivity(c.UserContext(), userID, page, limit)
	case "point", "activity":
		all = h.fetchPointActivity(c.UserContext(), userID, page, limit)
	default:
		return c.Status(400).JSON(fiber.Map{"error": "Invalid type. Valid: all, discussion, reply, certificate, showcase, study_case, course, point"})
	}

	// Sort by created_at DESC
	for i := 0; i < len(all); i++ {
		for j := i + 1; j < len(all); j++ {
			if all[j].CreatedAt.After(all[i].CreatedAt) {
				all[i], all[j] = all[j], all[i]
			}
		}
	}

	// Paginate
	total := len(all)
	offset := (page - 1) * limit
	if offset >= len(all) {
		offset = len(all)
	}
	end := offset + limit
	if end > len(all) {
		end = len(all)
	}
	all = all[offset:end]
	if all == nil {
		all = []dto.ActivityHistoryItem{}
	}

	return c.JSON(fiber.Map{
		"data": all,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetActivityLog returns user activity log (point-based, legacy)
func (h *UserHandler) GetActivityLog(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	page := clampPage(c.QueryInt("page", DefaultPage))
	limit := clampLimit(c.QueryInt("limit", DefaultLimit), DefaultLimit)

	logs, total, err := h.userSvc.GetUserActivityLog(c.UserContext(), userID, page, limit)
	if err != nil {
		return safeError(c, mapServiceErrorToStatus(err), err)
	}

	return c.JSON(fiber.Map{
		"data": logs,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// ListMentors returns paginated list of mentors
func (h *UserHandler) ListMentors(c *fiber.Ctx) error {
	page := clampPage(c.QueryInt("page", DefaultPage))
	limit := clampLimit(c.QueryInt("limit", DefaultLimit), DefaultLimit)

	mentors, total, err := h.userSvc.ListMentors(c.UserContext(), page, limit)
	if err != nil {
		return safeError(c, 500, err)
	}

	return c.JSON(fiber.Map{
		"data": mentors,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetPortfolio returns public portfolio for a user (GET /api/users/:id/portfolio)
func (h *UserHandler) GetPortfolio(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	user, err := h.userSvc.GetProfile(c.UserContext(), userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Gather portfolio items
	var items []dto.PortfolioItem

	// Certificates
	certs, certTotal, _ := h.certificateSvc.ListUserCertificates(c.UserContext(), userID, 1, 100)
	for _, cert := range certs {
		items = append(items, dto.PortfolioItem{
			ID:          cert.ID,
			Title:       cert.LessonName,
			Description: "Certificate of completion",
			Type:        "certificate",
			URL:         cert.VerificationURL,
			CreatedAt:   cert.IssuedAt,
			Tags:        []string{"certificate"},
		})
	}

	// Showcases
	showcases, showcaseTotal, _ := h.showcaseSvc.ListMyShowcases(c.UserContext(), userID, 1, 100)
	for _, sc := range showcases {
		images := make([]string, 0)
		if len(sc.MediaURLs) > 0 {
			images = sc.MediaURLs
		}
		items = append(items, dto.PortfolioItem{
			ID:          sc.ID,
			Title:       sc.Title,
			Description: sc.Description,
			Type:        "showcase",
			ImageURL:    images[0],
			CreatedAt:   sc.CreatedAt,
			Tags:        []string{"showcase"},
		})
	}

	// Study Cases
	studyCases, studyCaseTotal, _ := h.studyCaseSvc.ListStudyCasesByUser(c.UserContext(), userID, 1, 100)
	for _, sc := range studyCases {
		items = append(items, dto.PortfolioItem{
			ID:          sc.ID,
			Title:       sc.Name,
			Description: sc.Description,
			Type:        "study_case",
			ImageURL:    sc.ImgURL,
			Tags:        sc.Tags,
			CreatedAt:   sc.CreatedAt,
		})
	}

	return c.JSON(dto.PortfolioResponse{
		User: dto.ToUserBrief(*user),
		TotalPoints: user.TotalPoints,
		Level:       user.Level,
		Items:       items,
		CertCount:     int(certTotal),
		ShowcaseCount: int(showcaseTotal),
		StudyCaseCount: int(studyCaseTotal),
	})
}

// GetMyStreak returns current user's learning streak (GET /api/users/me/streak)
func (h *UserHandler) GetMyStreak(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	streak, err := h.streakSvc.GetUserStreak(c.UserContext(), userID)
	if err != nil {
		return safeError(c, 500, err)
	}

	return c.JSON(streak)
}

// GetAllUsers returns paginated list of all users (admin only)
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	page := clampPage(c.QueryInt("page", DefaultPage))
	limit := clampLimit(c.QueryInt("limit", DefaultLimit), DefaultLimit)

	users, total, err := h.userSvc.ListAllUsers(c.UserContext(), page, limit)
	if err != nil {
		return safeError(c, 500, err)
	}

	return c.JSON(fiber.Map{
		"data": users,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}
