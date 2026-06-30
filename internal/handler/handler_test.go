package handler

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/dto"
	"jvalleyverse/pkg/config"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	os.Setenv("JWT_SECRET", "test-secret-for-testing")
	os.Setenv("JWT_EXPIRY", "24h")
	config.LoadConfig()
}

func setupTestApp() *fiber.App {
	app := fiber.New()
	return app
}

// ──────────────────────────────────────────────
// AUTH FLOW TESTS (flow.md section 1)
// ──────────────────────────────────────────────

func TestAuthHandler_Register_Success(t *testing.T) {
	app := setupTestApp()
	mockSvc := newMockUserService()
	handler := NewAuthHandler(mockSvc)
	app.Post("/api/auth/register", handler.Register)

	body := `{"name":"Test User","email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "User created", result["message"])
}

func TestAuthHandler_Register_DuplicateEmail(t *testing.T) {
	app := setupTestApp()
	mockSvc := newMockUserService()
	mockSvc.addTestUser("existing-id", "existing@example.com", "password123", "Existing", "user")
	handler := NewAuthHandler(mockSvc)
	app.Post("/api/auth/register", handler.Register)

	body := `{"name":"Test User","email":"existing@example.com","password":"password123"}`
	req := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 409, resp.StatusCode)
}

func TestAuthHandler_Register_InvalidInput(t *testing.T) {
	app := setupTestApp()
	handler := NewAuthHandler(newMockUserService())
	app.Post("/api/auth/register", handler.Register)

	// Empty body
	req := httptest.NewRequest("POST", "/api/auth/register", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestAuthHandler_Register_ValidationErrors(t *testing.T) {
	app := setupTestApp()
	handler := NewAuthHandler(newMockUserService())
	app.Post("/api/auth/register", handler.Register)

	body := `{"name":"A","email":"invalid","password":"12"}`
	req := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Contains(t, result, "errors")
}

func TestAuthHandler_Login_Success(t *testing.T) {
	app := setupTestApp()
	mockSvc := newMockUserService()
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockSvc.users["test@example.com"] = &domain.User{
		ID:       "test-id",
		Email:    "test@example.com",
		Password: string(hashed),
		Name:     "Test User",
		Role:     "user",
		Avatar:   "https://example.com/avatar.jpg",
	}

	handler := NewAuthHandler(mockSvc)
	app.Post("/api/auth/login", handler.Login)

	body := `{"email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(t, result["access_token"])
	assert.NotEmpty(t, result["refresh_token"])
	assert.NotEmpty(t, result["xsrf_token"])
	assert.NotEmpty(t, result["expires_in"])
	assert.NotEmpty(t, result["user"])
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	app := setupTestApp()
	mockSvc := newMockUserService()
	handler := NewAuthHandler(mockSvc)
	app.Post("/api/auth/login", handler.Login)

	body := `{"email":"nonexistent@example.com","password":"wrongpass"}`
	req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestAuthHandler_Login_WrongPassword(t *testing.T) {
	app := setupTestApp()
	mockSvc := newMockUserService()
	hashed, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)
	mockSvc.users["user@test.com"] = &domain.User{
		ID:       "user-id",
		Email:    "user@test.com",
		Password: string(hashed),
		Name:     "Test User",
		Role:     "user",
	}
	handler := NewAuthHandler(mockSvc)
	app.Post("/api/auth/login", handler.Login)

	body := `{"email":"user@test.com","password":"wrongpass"}`
	req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestAuthHandler_Refresh_Success(t *testing.T) {
	app := setupTestApp()
	mockSvc := newMockUserService()
	mockSvc.addTestUser("test-id", "test@test.com", "pass", "Test", "user")
	mockSvc.GenerateRefreshToken(nil, "test-id")
	handler := NewAuthHandler(mockSvc)
	app.Post("/api/auth/refresh", handler.Refresh)

	body := `{"refresh_token":"mock-refresh-token-test-id"}`
	req := httptest.NewRequest("POST", "/api/auth/refresh", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.NotEmpty(t, result["access_token"])
	assert.NotEmpty(t, result["xsrf_token"])
	assert.NotEmpty(t, result["expires_in"])
}

func TestAuthHandler_Refresh_InvalidToken(t *testing.T) {
	app := setupTestApp()
	mockSvc := newMockUserService()
	handler := NewAuthHandler(mockSvc)
	app.Post("/api/auth/refresh", handler.Refresh)

	body := `{"refresh_token":"invalid-token"}`
	req := httptest.NewRequest("POST", "/api/auth/refresh", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
}

func TestAuthHandler_Logout(t *testing.T) {
	mockSvc := newMockUserService()
	mockSvc.addTestUser("test-id", "test@test.com", "pass", "Test", "user")
	handler := NewAuthHandler(mockSvc)

	app := setupTestApp()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", "test-id")
		return c.Next()
	})
	app.Post("/api/auth/logout", handler.Logout)

	req := httptest.NewRequest("POST", "/api/auth/logout", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Contains(t, result["message"], "Logged out")
}

func TestAuthHandler_GoogleLogin_MissingToken(t *testing.T) {
	app := setupTestApp()
	mockSvc := newMockUserService()
	handler := NewAuthHandler(mockSvc)
	app.Post("/api/auth/google", handler.GoogleLogin)

	body := `{}`
	req := httptest.NewRequest("POST", "/api/auth/google", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

// ──────────────────────────────────────────────
// USER FLOW TESTS (flow.md section 8)
// ──────────────────────────────────────────────

func TestUserHandler_GetProfile(t *testing.T) {
	app := setupProtectedApp("test-id", "user")
	mockSvc := newMockUserService()
	mockSvc.addTestUser("test-id", "user@test.com", "pass", "Test User", "user")
	handler := NewUserHandler(mockSvc, newMockDashboardService())
	app.Get("/api/users/me", handler.GetProfile)

	req := httptest.NewRequest("GET", "/api/users/me", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Test User", result["name"])
	assert.Equal(t, "user@test.com", result["email"])
	assert.Equal(t, "user", result["role"])
}

func TestUserHandler_GetProfile_WithProtectedApp(t *testing.T) {
	app := setupProtectedApp("test-id", "user")
	mockSvc := newMockUserService()
	mockSvc.addTestUser("test-id", "user@test.com", "pass", "Test User", "user")
	handler := NewUserHandler(mockSvc, newMockDashboardService())
	app.Get("/api/users/me", handler.GetProfile)

	req := httptest.NewRequest("GET", "/api/users/me", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Test User", result["name"])
	assert.Equal(t, "user@test.com", result["email"])
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	app := setupProtectedApp("test-id", "user")
	mockSvc := newMockUserService()
	mockSvc.addTestUser("test-id", "user@test.com", "pass", "Old Name", "user")
	handler := NewUserHandler(mockSvc, newMockDashboardService())
	app.Put("/api/users/me", handler.UpdateProfile)

	body := `{"name":"New Name","bio":"Updated bio","avatar":"https://example.com/new-avatar.jpg"}`
	req := httptest.NewRequest("PUT", "/api/users/me", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify profile was updated
	assert.Equal(t, "New Name", mockSvc.users["test-id"].Name)
	assert.Equal(t, "Updated bio", mockSvc.users["test-id"].Bio)
}

func TestUserHandler_GetDashboard(t *testing.T) {
	app := setupProtectedApp("test-id", "user")
	handler := NewUserHandler(newMockUserService(), newMockDashboardService())
	app.Get("/api/users/me/dashboard", handler.GetDashboard)

	req := httptest.NewRequest("GET", "/api/users/me/dashboard", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, float64(2), result["courses_in_progress"])
	assert.Equal(t, float64(3), result["courses_completed"])
	assert.Equal(t, float64(5), result["unread_notifications"])
}

func TestUserHandler_GetActivityLog(t *testing.T) {
	app := setupProtectedApp("test-id", "user")
	mockSvc := newMockUserService()
	mockSvc.activityLog = []dto.ActivityItem{
		{ID: "act1", Activity: "complete_lesson", Points: 50, Timestamp: time.Now()},
		{ID: "act2", Activity: "create_showcase", Points: 10, Timestamp: time.Now()},
	}
	handler := NewUserHandler(mockSvc, newMockDashboardService())
	app.Get("/api/users/me/activity", handler.GetActivityLog)

	req := httptest.NewRequest("GET", "/api/users/me/activity", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].([]interface{})
	assert.Equal(t, 2, len(data))
	pagination := result["pagination"].(map[string]interface{})
	assert.Equal(t, float64(2), pagination["total"])
}

func TestUserHandler_GetPublicProfile(t *testing.T) {
	app := setupTestApp()
	mockSvc := newMockUserService()
	mockSvc.addTestUser("public-id", "public@test.com", "pass", "Public User", "user")
	handler := NewUserHandler(mockSvc, newMockDashboardService())
	app.Get("/api/users/:id", handler.GetPublicProfile)

	req := httptest.NewRequest("GET", "/api/users/public-id", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Public User", result["name"])
	assert.NotContains(t, result, "email") // public profile should not expose email
}

func TestUserHandler_GetPublicProfile_NotFound(t *testing.T) {
	app := setupTestApp()
	handler := NewUserHandler(newMockUserService(), newMockDashboardService())
	app.Get("/api/users/:id", handler.GetPublicProfile)

	req := httptest.NewRequest("GET", "/api/users/nonexistent-id", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestUserHandler_ListMentors(t *testing.T) {
	app := setupTestApp()
	mockSvc := newMockUserService()
	mockSvc.mentors = []dto.MentorItem{
		{ID: "mentor1", Name: "Mentor One", Level: 5, TotalPoints: 2500},
		{ID: "mentor2", Name: "Mentor Two", Level: 4, TotalPoints: 1500},
	}
	handler := NewUserHandler(mockSvc, newMockDashboardService())
	app.Get("/api/mentors", handler.ListMentors)

	req := httptest.NewRequest("GET", "/api/mentors", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// ──────────────────────────────────────────────
// HEALTH CHECK TESTS
// ──────────────────────────────────────────────

func TestHealthHandler_Health(t *testing.T) {
	app := setupTestApp()
	handler := NewHealthHandler()
	app.Get("/api/health", handler.Health)

	req := httptest.NewRequest("GET", "/api/health", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "ok", result["status"])
	assert.Equal(t, "1.0.0", result["version"])
	assert.NotEmpty(t, result["timestamp"])
}

func TestHealthHandler_HealthDetailed_Local(t *testing.T) {
	app := setupTestApp()
	handler := NewHealthHandler()
	app.Get("/api/health/detailed", handler.HealthDetailed)

	req := httptest.NewRequest("GET", "/api/health/detailed", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "ok", result["status"])
	assert.Equal(t, "local", result["environment"])
}

// ──────────────────────────────────────────────
// ERROR MAPPING TESTS
// ──────────────────────────────────────────────

func TestMapServiceErrorToStatus(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"ErrInvalidInput", domain.ErrInvalidInput, 400},
		{"ErrUnauthorized", domain.ErrUnauthorized, 401},
		{"ErrForbidden", domain.ErrForbidden, 403},
		{"ErrNotFound", domain.ErrNotFound, 404},
		{"ErrCourseNotFound", domain.ErrCourseNotFound, 404},
		{"ErrLessonNotFound", domain.ErrLessonNotFound, 404},
		{"ErrUserNotFound", domain.ErrUserNotFound, 404},
		{"ErrStudyCaseNotFound", domain.ErrStudyCaseNotFound, 404},
		{"ErrEmailExists", domain.ErrEmailExists, 409},
		{"unknown error", assert.AnError, 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapServiceErrorToStatus(tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSafeError(t *testing.T) {
	app := fiber.New()
	app.Get("/test-safe-error", func(c *fiber.Ctx) error {
		return safeError(c, 400, domain.ErrInvalidInput)
	})

	req := httptest.NewRequest("GET", "/test-safe-error", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}
