package handler

import (
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/service"

	"github.com/gofiber/fiber/v2"
)

type StudyCaseHandler struct {
	studyCaseSvc service.IStudyCaseService
}

func NewStudyCaseHandler(studyCaseSvc service.IStudyCaseService) *StudyCaseHandler {
	return &StudyCaseHandler{studyCaseSvc: studyCaseSvc}
}

// ListMyStudyCases returns current user's study cases (GET /api/users/me/study-cases)
func (h *StudyCaseHandler) ListMyStudyCases(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	page := clampPage(c.QueryInt("page", DefaultPage))
	limit := clampLimit(c.QueryInt("limit", DefaultLimit), DefaultLimit)

	data, total, err := h.studyCaseSvc.ListStudyCasesByUser(c.UserContext(), userID, page, limit)
	if err != nil {
		return safeError(c, 500, err)
	}

	return c.JSON(fiber.Map{
		"data": data,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// ListStudyCases godoc
// @Summary      List study cases
// @Description  Get paginated list of study cases with optional filter by category_id
// @Tags         Study Cases
// @Param        page         query int     false  "Page number (default: 1)"
// @Param        limit        query int     false  "Items per page (default: 20)"
// @Param        category_id  query string  false  "Filter by category ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /study-cases [get]
func (h *StudyCaseHandler) ListStudyCases(c *fiber.Ctx) error {
	page := clampPage(c.QueryInt("page", DefaultPage))
	limit := clampLimit(c.QueryInt("limit", DefaultLimit), DefaultLimit)

	categoryID := c.Query("category_id", "")
	var filter *service.StudyCaseListFilter
	if categoryID != "" {
		filter = &service.StudyCaseListFilter{CategoryID: &categoryID}
	}

	data, total, err := h.studyCaseSvc.ListStudyCases(c.UserContext(), page, limit, filter)
	if err != nil {
		return safeError(c, 500, err)
	}

	return c.JSON(fiber.Map{
		"data": data,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetStudyCase returns a study case by ID (GET /api/study-cases/:id)
func (h *StudyCaseHandler) GetStudyCase(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Study case ID is required"})
	}

	data, err := h.studyCaseSvc.GetStudyCaseByID(c.UserContext(), id)
	if err != nil {
		if err == domain.ErrStudyCaseNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "Study case not found"})
		}
		return safeError(c, 500, err)
	}

	return c.JSON(data)
}

// CreateStudyCase creates a new study case (POST /api/admin/study-cases)
func (h *StudyCaseHandler) CreateStudyCase(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var input struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		ImgURL      string   `json:"img_url"`
		YoutubeURL  string   `json:"youtube_url"`
		CategoryID  string   `json:"category_id"`
		Tags        []string `json:"tags"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if input.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name is required"})
	}

	var categoryID *string
	if input.CategoryID != "" {
		categoryID = &input.CategoryID
	}

	studyCase, err := h.studyCaseSvc.CreateStudyCase(c.UserContext(), userID, input.Name, input.Description, input.ImgURL, input.YoutubeURL, categoryID, input.Tags)
	if err != nil {
		return safeError(c, mapServiceErrorToStatus(err), err)
	}

	logAdminAction(c.UserContext(), userID, "create", "study_case", studyCase.ID, "Name: "+input.Name)
	return c.Status(201).JSON(studyCase)
}

// UpdateStudyCase updates a study case (PUT /api/admin/study-cases/:id)
func (h *StudyCaseHandler) UpdateStudyCase(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Study case ID is required"})
	}

	var input struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		ImgURL      string   `json:"img_url"`
		YoutubeURL  string   `json:"youtube_url"`
		CategoryID  string   `json:"category_id"`
		Tags        []string `json:"tags"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var categoryID *string
	if input.CategoryID != "" {
		categoryID = &input.CategoryID
	}

	studyCase, err := h.studyCaseSvc.UpdateStudyCase(c.UserContext(), id, input.Name, input.Description, input.ImgURL, input.YoutubeURL, categoryID, input.Tags)
	if err != nil {
		if err == domain.ErrStudyCaseNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "Study case not found"})
		}
		return safeError(c, mapServiceErrorToStatus(err), err)
	}

	logAdminAction(c.UserContext(), c.Locals("userID").(string), "update", "study_case", id, "Name: "+input.Name)
	return c.JSON(studyCase)
}

// DeleteStudyCase deletes a study case (DELETE /api/admin/study-cases/:id)
func (h *StudyCaseHandler) DeleteStudyCase(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Study case ID is required"})
	}

	if err := h.studyCaseSvc.DeleteStudyCase(c.UserContext(), id); err != nil {
		if err == domain.ErrStudyCaseNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "Study case not found"})
		}
		return safeError(c, mapServiceErrorToStatus(err), err)
	}

	logAdminAction(c.UserContext(), c.Locals("userID").(string), "delete", "study_case", id, "")
	return c.JSON(fiber.Map{"message": "Study case deleted"})
}
