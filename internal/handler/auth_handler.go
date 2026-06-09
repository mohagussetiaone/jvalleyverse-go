package handler

import (
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/service"
	"jvalleyverse/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userSvc service.IUserService
}

func NewAuthHandler(userSvc service.IUserService) *AuthHandler {
	return &AuthHandler{userSvc: userSvc}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	user := &domain.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashed),
		Role:     "user",
		Points:   0,
		Level:    1,
	}

	if err := h.userSvc.CreateUser(c.UserContext(), user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Email already exists or internal error"})
	}

	return c.Status(210).JSON(fiber.Map{"message": "User created"})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	user, err := h.userSvc.GetUserByEmail(c.UserContext(), input.Email)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not generate token"})
	}

	// Set XSRF cookie
	xsrfToken := utils.GenerateXSRFToken()
	c.Cookie(&fiber.Cookie{
		Name:     "XSRF-TOKEN",
		Value:    xsrfToken,
		HTTPOnly: false,
		Secure:   false, // set true jika https
		SameSite: "Strict",
	})

	return c.JSON(fiber.Map{"token": token, "xsrf_token": xsrfToken})
}
