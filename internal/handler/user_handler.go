package handler

import (
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

// GetActivityLog returns user activity log
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
