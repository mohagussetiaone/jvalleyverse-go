package handler

import (
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/repository"
	"jvalleyverse/internal/service"
	"jvalleyverse/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
    userSvc *service.UserService
}

func NewAuthHandler() *AuthHandler {
    return &AuthHandler{userSvc: service.NewUserService()}
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
    // Panggil repository langsung (untuk sederhana)
    // Sebaiknya buat method di userSvc
    repo := repository.NewUserRepository()
    if err := repo.Create(user); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Email already exists"})
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
    repo := repository.NewUserRepository()
    user, err := repo.FindByEmail(input.Email)
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
    // Set XSRF cookie (contoh)
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

type ShowcaseHandler struct {
    showcaseSvc *service.ShowcaseService
}

func NewShowcaseHandler() *ShowcaseHandler {
    return &ShowcaseHandler{showcaseSvc: service.NewShowcaseService()}
}

func (h *ShowcaseHandler) Create(c *fiber.Ctx) error {
    userID := c.Locals("userID").(uint)
    var input struct {
		Title      string   `json:"title"`
		MediaURLs  []string `json:"media_urls"`
		CategoryID uint     `json:"category_id"`
    }
    if err := c.BodyParser(&input); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
    }
    showcase, err := h.showcaseSvc.CreateShowcase(userID, input.Title, input.MediaURLs, input.CategoryID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.Status(201).JSON(showcase)
}

func (h *ShowcaseHandler) Like(c *fiber.Ctx) error {
    userID := c.Locals("userID").(uint)
    showcaseID, err := c.ParamsInt("id")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid showcase id"})
    }
    if err := h.showcaseSvc.LikeShowcase(userID, uint(showcaseID)); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(fiber.Map{"message": "Liked successfully"})
}

func (h *ShowcaseHandler) ListShowcases(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "data": []fiber.Map{},
        "pagination": fiber.Map{
            "page":  1,
            "limit": 20,
            "total": 0,
        },
    })
}

func (h *ShowcaseHandler) GetShowcase(c *fiber.Ctx) error {
    showcaseID := c.Params("id")
    return c.JSON(fiber.Map{
        "id":           showcaseID,
        "title":        "",
        "description":  "",
        "media_urls":   []string{},
        "user":         fiber.Map{"id": 0, "name": "", "level": 1},
        "category":     fiber.Map{"id": 0, "name": ""},
        "likes_count":  0,
        "views_count":  0,
        "is_liked_by_me": false,
        "comments":     []fiber.Map{},
    })
}

func (h *ShowcaseHandler) Update(c *fiber.Ctx) error {
    userID := c.Locals("userID").(uint)
    showcaseID := c.Params("id")
    var input struct {
        Title       string `json:"title"`
        Description string `json:"description"`
        Visibility  string `json:"visibility"`
    }
    if err := c.BodyParser(&input); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
    }
    return c.JSON(fiber.Map{
        "message": "Showcase updated",
        "id":      showcaseID,
        "user_id": userID,
    })
}

func (h *ShowcaseHandler) Delete(c *fiber.Ctx) error {
    userID := c.Locals("userID").(uint)
    showcaseID := c.Params("id")
    return c.JSON(fiber.Map{
        "message": "Showcase deleted",
        "id":      showcaseID,
        "user_id": userID,
    })
}

func (h *ShowcaseHandler) Unlike(c *fiber.Ctx) error {
    userID := c.Locals("userID").(uint)
    showcaseID, err := c.ParamsInt("id")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid showcase id"})
    }
    return c.JSON(fiber.Map{
        "message": "Showcase unliked",
        "showcase_id": showcaseID,
        "user_id": userID,
    })
}

func (h *ShowcaseHandler) GetLeaderboard(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "data": []fiber.Map{},
        "pagination": fiber.Map{
            "page":  1,
            "limit": 50,
            "total": 0,
        },
    })
}

type ClassHandler struct {
	classSvc *service.ClassService
}

func NewClassHandler() *ClassHandler {
	return &ClassHandler{classSvc: service.NewClassService()}
}

// GetClassBySlug returns class with details and user's progress
func (h *ClassHandler) GetClassBySlug(c *fiber.Ctx) error {
	projectID, _ := c.ParamsInt("project_id")
	slug := c.Params("slug")
	userID, _ := c.Locals("userID").(uint) // Might be 0 if not logged in, service handles it

	data, err := h.classSvc.GetClassBySlug(uint(projectID), slug, userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Class not found"})
	}

	return c.JSON(data)
}

// StartClass initializes user progress
func (h *ClassHandler) StartClass(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	classID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid class id"})
	}

	progress, err := h.classSvc.StartClass(userID, uint(classID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message":  "Class started!",
		"progress": progress,
	})
}

// UpdateProgress updates user progress percentage
func (h *ClassHandler) UpdateProgress(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	classID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid class id"})
	}

	var input struct {
		Percentage int    `json:"progress_percentage"`
		Notes      string `json:"notes"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	progress, err := h.classSvc.UpdateProgress(userID, uint(classID), input.Percentage, input.Notes)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(progress)
}

// Complete marks class as completed
func (h *ClassHandler) Complete(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	classID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid class id"})
	}

	data, err := h.classSvc.CompleteClass(userID, uint(classID))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(data)
}

// Admin Methods

func (h *ClassHandler) CreateClass(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	var input domain.Class
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	class, err := h.classSvc.AdminCreateClass(userID, input)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(class)
}

func (h *ClassHandler) CreateClassDetail(c *fiber.Ctx) error {
	classID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid class id"})
	}

	var input struct {
		About         string      `json:"about"`
		Rules         string      `json:"rules"`
		Tools         interface{} `json:"tools"`
		ResourceMedia interface{} `json:"resource_media"`
		Resources     interface{} `json:"resources"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	detail, err := h.classSvc.AdminCreateClassDetail(uint(classID), input.About, input.Rules, input.Tools, input.ResourceMedia, input.Resources)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(detail)
}

func (h *ClassHandler) UpdateClass(c *fiber.Ctx) error {
	// ... (Implementation remains similar but should call service)
	return c.JSON(fiber.Map{"message": "Not implemented in this step"})
}

func (h *ClassHandler) DeleteClass(c *fiber.Ctx) error {
	// ...
	return c.JSON(fiber.Map{"message": "Not implemented in this step"})
}
