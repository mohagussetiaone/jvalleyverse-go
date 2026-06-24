package handler

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/dto"
	"jvalleyverse/pkg/config"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// mockUserService implements service.IUserService for testing
type mockUserService struct {
	users map[string]*domain.User
}

func newMockUserService() *mockUserService {
	return &mockUserService{users: make(map[string]*domain.User)}
}

func (m *mockUserService) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	user, ok := m.users[userID]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (m *mockUserService) UpdateProfile(ctx context.Context, userID string, name, bio, avatar string) error {
	user, ok := m.users[userID]
	if !ok {
		return domain.ErrUserNotFound
	}
	user.Name = name
	user.Bio = bio
	user.Avatar = avatar
	return nil
}

func (m *mockUserService) AddPoints(ctx context.Context, userID string, category string, points int, metadata map[string]interface{}) error {
	return nil
}

func (m *mockUserService) GetUserActivityLog(ctx context.Context, userID string, page, limit int) ([]dto.ActivityItem, int64, error) {
	return []dto.ActivityItem{}, 0, nil
}

func (m *mockUserService) CreateUser(ctx context.Context, user *domain.User) error {
	if _, ok := m.users[user.Email]; ok {
		return domain.ErrEmailExists
	}
	m.users[user.Email] = user
	m.users[user.ID] = user
	return nil
}

func (m *mockUserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, ok := m.users[email]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (m *mockUserService) ListAllUsers(ctx context.Context, page, limit int) ([]dto.UserListItem, int64, error) {
	return []dto.UserListItem{}, 0, nil
}

func (m *mockUserService) GenerateRefreshToken(ctx context.Context, userID string) (*domain.RefreshToken, error) {
	return &domain.RefreshToken{Token: "mock-refresh-token", UserID: userID}, nil
}

func (m *mockUserService) ValidateRefreshToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	if token == "mock-refresh-token" {
		return &domain.RefreshToken{Token: token, UserID: "test-id"}, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockUserService) RevokeRefreshToken(ctx context.Context, token string) error {
	return nil
}

func (m *mockUserService) RevokeAllUserTokens(ctx context.Context, userID string) error {
	return nil
}

func (m *mockUserService) ListMentors(ctx context.Context, page, limit int) ([]dto.MentorItem, int64, error) {
	return []dto.MentorItem{}, 0, nil
}

func setupTestApp() *fiber.App {
	config.LoadConfig()
	app := fiber.New()
	return app
}

func init() {
	os.Setenv("JWT_SECRET", "test-secret-for-testing")
}

func TestAuthHandler_Register_Success(t *testing.T) {
	app := setupTestApp()
	handler := NewAuthHandler(newMockUserService())
	app.Post("/api/auth/register", handler.Register)

	body := `{"name":"Test User","email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestAuthHandler_Register_InvalidInput(t *testing.T) {
	app := setupTestApp()
	handler := NewAuthHandler(newMockUserService())
	app.Post("/api/auth/register", handler.Register)

	req := httptest.NewRequest("POST", "/api/auth/register", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
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
		{"unknown error", assert.AnError, 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapServiceErrorToStatus(tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}


