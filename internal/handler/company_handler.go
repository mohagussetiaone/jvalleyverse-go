package handler

import (
	"github.com/gofiber/fiber/v2"

	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/service"
)

type CompanyHandler struct {
	companySvc service.ICompanyService
}

func NewCompanyHandler(companySvc service.ICompanyService) *CompanyHandler {
	return &CompanyHandler{companySvc: companySvc}
}

// GET /api/company — public, no auth required
func (h *CompanyHandler) GetCompany(c *fiber.Ctx) error {
	company, err := h.companySvc.GetCompany(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch company profile"})
	}
	return c.JSON(company)
}

// PUT /api/admin/company — admin only, update company profile
func (h *CompanyHandler) UpdateCompany(c *fiber.Ctx) error {
	var input domain.Company
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	company, err := h.companySvc.UpdateCompany(c.UserContext(), input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Log admin action
	userID := c.Locals("userID").(string)
	logAdminAction(c.UserContext(), userID, "update", "company", company.ID, "")

	return c.JSON(company)
}
