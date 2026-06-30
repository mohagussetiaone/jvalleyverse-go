package handler

import (
	"context"
	"time"

	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/service"
	"jvalleyverse/pkg/config"
	"jvalleyverse/pkg/utils"
	"jvalleyverse/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"
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

	if errs := validateRegisterInput(input); len(errs) > 0 {
		return c.Status(400).JSON(fiber.Map{"errors": errs})
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
		return safeError(c, mapServiceErrorToStatus(err), err)
	}

	return c.Status(201).JSON(fiber.Map{"message": "User created"})
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

	return h.generateAuthResponse(c, user)
}

// GoogleLogin handles Google One Tap / Sign-In with ID token.
// The frontend sends the credential JWT from Google Identity Services.
func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	var input struct {
		Token string `json:"token"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if input.Token == "" {
		return c.Status(400).JSON(fiber.Map{"error": "token is required"})
	}

	clientID := config.AppConfig.GoogleClientID
	if clientID == "" {
		return c.Status(500).JSON(fiber.Map{"error": "Google login not configured"})
	}

	// Verify the Google ID token
	payload, err := idtoken.Validate(context.Background(), input.Token, clientID)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid Google token"})
	}

	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)

	if email == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Email not provided by Google"})
	}

	// Find existing user by email, or create a new one
	user, err := h.userSvc.GetUserByEmail(c.UserContext(), email)
	if err != nil {
		// User not found — create a new one with Google data
		user = &domain.User{
			Name:     name,
			Email:    email,
			Avatar:   picture,
			Role:     "user",
			Password: "", // OAuth users have no password
			Points:   0,
			Level:    1,
			IsActive: true,
		}
		if err := h.userSvc.CreateUser(c.UserContext(), user); err != nil {
			return safeError(c, mapServiceErrorToStatus(err), err)
		}
	} else {
		// Update avatar/name from Google if they've changed
		if picture != "" && picture != user.Avatar {
			user.Avatar = picture
		}
		if name != "" && name != user.Name {
			user.Name = name
		}
		// Re-fetch after potential update
		user, _ = h.userSvc.GetUserByEmail(c.UserContext(), email)
	}

	return h.generateAuthResponse(c, user)
}

// generateAuthResponse creates JWT, refresh token, XSRF token and returns them.
// Used by both Login and GoogleLogin to avoid duplication.
func (h *AuthHandler) generateAuthResponse(c *fiber.Ctx, user *domain.User) error {
	accessToken, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not generate token"})
	}

	refreshToken, err := h.userSvc.GenerateRefreshToken(c.UserContext(), user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not generate refresh token"})
	}

	xsrfToken := utils.GenerateXSRFToken()
	c.Cookie(&fiber.Cookie{
		Name:     "XSRF-TOKEN",
		Value:    xsrfToken,
		HTTPOnly: false,
		Secure:   false,
		SameSite: "Lax",
		MaxAge:   int(config.AppConfig.JWTExpiry.Seconds()),
	})

	return c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken.Token,
		"expires_in":    int(config.AppConfig.JWTExpiry.Seconds()),
		"xsrf_token":    xsrfToken,
		"user": fiber.Map{
			"id":     user.ID,
			"name":   user.Name,
			"email":  user.Email,
			"avatar": user.Avatar,
			"role":   user.Role,
		},
	})
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if input.RefreshToken == "" {
		return c.Status(400).JSON(fiber.Map{"error": "refresh_token is required"})
	}

	tokenData, err := h.userSvc.ValidateRefreshToken(c.UserContext(), input.RefreshToken)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired refresh token"})
	}

	newAccessToken, err := utils.GenerateJWT(tokenData.UserID, "")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not generate token"})
	}

	// Get user role for the new token
	user, err := h.userSvc.GetProfile(c.UserContext(), tokenData.UserID)
	if err == nil {
		newAccessToken, _ = utils.GenerateJWT(user.ID, user.Role)
	}

	// Generate new XSRF token on refresh so cookie stays in sync with session
	xsrfToken := utils.GenerateXSRFToken()
	c.Cookie(&fiber.Cookie{
		Name:     "XSRF-TOKEN",
		Value:    xsrfToken,
		HTTPOnly: false,
		Secure:   false,
		SameSite: "Lax",
		MaxAge:   int(config.AppConfig.JWTExpiry.Seconds()),
	})

	return c.JSON(fiber.Map{
		"access_token": newAccessToken,
		"expires_in":   int(config.AppConfig.JWTExpiry.Seconds()),
		"xsrf_token":   xsrfToken,
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var input struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BodyParser(&input); err == nil && input.RefreshToken != "" {
		h.userSvc.RevokeRefreshToken(c.UserContext(), input.RefreshToken)
	} else {
		h.userSvc.RevokeAllUserTokens(c.UserContext(), userID)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "XSRF-TOKEN",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: false,
	})

	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}

func validateRegisterInput(input struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}) []validator.ValidationError {
	var errs []validator.ValidationError

	if err := validator.ValidateName(input.Name); err != nil {
		errs = append(errs, *err)
	}
	if err := validator.ValidateEmail(input.Email); err != nil {
		errs = append(errs, *err)
	}
	if err := validator.ValidatePassword(input.Password); err != nil {
		errs = append(errs, *err)
	}

	return errs
}
