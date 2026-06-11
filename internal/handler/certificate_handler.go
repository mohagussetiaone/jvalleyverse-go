package handler

import (
	"jvalleyverse/internal/service"

	"github.com/gofiber/fiber/v2"
)

type CertificateHandler struct {
	certificateSvc service.ICertificateService
}

func NewCertificateHandler() *CertificateHandler {
	return &CertificateHandler{certificateSvc: service.GetCertificateService()}
}

// ListCertificates returns user's certificates
func (h *CertificateHandler) ListCertificates(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	certs, err := h.certificateSvc.ListUserCertificates(c.UserContext(), userID, page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": certs,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": len(certs),
		},
	})
}

// GetCertificate returns specific certificate by unique code
func (h *CertificateHandler) GetCertificate(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	role, _ := c.Locals("role").(string)
	code := c.Params("code")
	if code == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid certificate code"})
	}

	cert, err := h.certificateSvc.GetCertificateByCode(c.UserContext(), code, userID, role)
	if err != nil {
		return c.Status(mapServiceErrorToStatus(err)).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(cert)
}
