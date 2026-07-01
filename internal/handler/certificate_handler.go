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
	page := clampPage(c.QueryInt("page", DefaultPage))
	limit := clampLimit(c.QueryInt("limit", DefaultLimit), DefaultLimit)

	certs, total, err := h.certificateSvc.ListUserCertificates(c.UserContext(), userID, page, limit)
	if err != nil {
		return safeError(c, 500, err)
	}

	return c.JSON(fiber.Map{
		"data": certs,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetCertificate returns specific certificate by unique code (owner/admin only)
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

// VerifyCertificate returns public certificate verification info (GET /api/certificates/:code/verify)
// No auth required — anyone can verify a certificate
func (h *CertificateHandler) VerifyCertificate(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid certificate code"})
	}

	cert, err := h.certificateSvc.VerifyCertificateByCode(c.UserContext(), code)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Certificate not found"})
	}

	return c.JSON(cert)
}
