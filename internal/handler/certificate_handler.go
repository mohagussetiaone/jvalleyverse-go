package handler

import (
	"jvalleyverse/internal/service"

	"github.com/gofiber/fiber/v2"
)

type CertificateHandler struct {
	certificateSvc *service.CertificateService
}

func NewCertificateHandler() *CertificateHandler {
	return &CertificateHandler{certificateSvc: service.NewCertificateService()}
}

// ListCertificates returns user's certificates
func (h *CertificateHandler) ListCertificates(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	_ = userID // TODO: Query certificates from service

	return c.JSON(fiber.Map{
		"data": []fiber.Map{},
		"pagination": fiber.Map{
			"page":  1,
			"limit": 20,
			"total": 0,
		},
	})
}

// GetCertificate returns specific certificate
func (h *CertificateHandler) GetCertificate(c *fiber.Ctx) error {
	code := c.Params("code")
	userID := c.Locals("userID").(uint)
	_ = userID // TODO: Verify certificate owner

	return c.JSON(fiber.Map{
		"id":           0,
		"unique_code":  code,
		"user_id":      userID,
		"class_id":     0,
		"badge_url":    "",
		"issued_at":    "",
		"issued_to":    "",
	})
}
