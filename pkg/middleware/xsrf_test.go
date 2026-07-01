package middleware

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"jvalleyverse/pkg/config"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ──── helpers ─────────────────────────────────────────────────────────────

// initTestConfig ensures config.AppConfig is initialized for tests.
func initTestConfig() {
	if config.AppConfig == nil {
		config.AppConfig = &config.Config{
			CORSOrigins: "http://localhost:3000,http://localhost:5173,http://localhost:5174,https://jvalleyverse.web.id",
		}
	}
}

// setupXSRFApp returns a Fiber app with XSRFProtection middleware.
func setupXSRFApp() *fiber.App {
	initTestConfig()
	app := fiber.New()
	app.Use(XSRFProtection())
	app.Post("/dangerous", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"status": "ok"})
	})
	app.Get("/dangerous", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"status": "ok"})
	})
	app.Get("/safe", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"status": "ok"})
	})
	return app
}

// setCORSOrigins temporarily overrides config for testing.
func setCORSOrigins(origins string) func() {
	initTestConfig()
	original := config.AppConfig.CORSOrigins
	config.AppConfig.CORSOrigins = origins
	return func() {
		config.AppConfig.CORSOrigins = original
	}
}

// ──── Tests: GET/HEAD/OPTIONS should always pass ──────────────────────────

func TestXSRF_SkipsSafeMethods(t *testing.T) {
	app := setupXSRFApp()

	tests := []struct {
		method string
		path   string
	}{
		{"GET", "/safe"},
		{"HEAD", "/safe"},
		{"GET", "/dangerous"},     // now has GET handler registered
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)
		})
	}
}

// ──── Tests: Cookie check (standard double-submit) ────────────────────────

func TestXSRF_CookieMatch_Allows(t *testing.T) {
	app := setupXSRFApp()

	req := httptest.NewRequest("POST", "/dangerous", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-XSRF-TOKEN", "my-token-123")
	req.Header.Set("Cookie", "XSRF-TOKEN=my-token-123")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestXSRF_CookieMismatch_Blocks(t *testing.T) {
	app := setupXSRFApp()

	req := httptest.NewRequest("POST", "/dangerous", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-XSRF-TOKEN", "token-a")
	req.Header.Set("Cookie", "XSRF-TOKEN=token-b")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 403, resp.StatusCode)

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "XSRF token invalid", body["error"])
}

func TestXSRF_EmptyCookie_Blocks(t *testing.T) {
	app := setupXSRFApp()

	req := httptest.NewRequest("POST", "/dangerous", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-XSRF-TOKEN", "some-token")
	req.Header.Set("Cookie", "")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 403, resp.StatusCode)
}

func TestXSRF_EmptyHeader_Blocks(t *testing.T) {
	app := setupXSRFApp()

	req := httptest.NewRequest("POST", "/dangerous", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "XSRF-TOKEN=some-token")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 403, resp.StatusCode)
}

// ──── Tests: Origin fallback ──────────────────────────────────────────────

func TestXSRF_OriginMatch_Allows(t *testing.T) {
	restore := setCORSOrigins("http://localhost:3000,https://jvalleyverse.web.id")
	defer restore()

	app := setupXSRFApp()

	req := httptest.NewRequest("POST", "/dangerous", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "http://localhost:3000")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestXSRF_OriginMismatch_Blocks(t *testing.T) {
	restore := setCORSOrigins("http://localhost:3000")
	defer restore()

	app := setupXSRFApp()

	req := httptest.NewRequest("POST", "/dangerous", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://evil-site.com")
	// No cookie match either

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 403, resp.StatusCode)
}

func TestXSRF_OriginWithTrailingSlash_Allows(t *testing.T) {
	restore := setCORSOrigins("http://localhost:5173")
	defer restore()

	app := setupXSRFApp()

	// SPA often sends Origin with trailing slash
	req := httptest.NewRequest("POST", "/dangerous", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "http://localhost:5173/")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestXSRF_OriginPortMatch_Allows(t *testing.T) {
	restore := setCORSOrigins("http://localhost:5173,http://localhost:3000")
	defer restore()

	app := setupXSRFApp()

	req := httptest.NewRequest("POST", "/dangerous", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "http://localhost:5173")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// ──── Tests: Referer fallback ─────────────────────────────────────────────

func TestXSRF_RefererMatch_Allows(t *testing.T) {
	restore := setCORSOrigins("http://localhost:3000,https://jvalleyverse.web.id")
	defer restore()

	app := setupXSRFApp()

	req := httptest.NewRequest("POST", "/dangerous", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	// No Origin header — should fall back to Referer
	req.Header.Set("Referer", "http://localhost:3000/some-page")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestXSRF_RefererUsedOnlyWhenNoOrigin(t *testing.T) {
	restore := setCORSOrigins("http://localhost:5173")
	defer restore()

	app := setupXSRFApp()

	// Origin is set but doesn't match → even if Referer matches, should block
	req := httptest.NewRequest("POST", "/dangerous", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://evil.com")
	req.Header.Set("Referer", "http://localhost:5173/page")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 403, resp.StatusCode)
}

// ──── Tests: Cookie takes priority over Origin ────────────────────────────

func TestXSRF_CookiePriority(t *testing.T) {
	restore := setCORSOrigins("http://localhost:5173")
	defer restore()

	app := setupXSRFApp()

	// Even if Origin is from evil site, valid cookie should pass
	req := httptest.NewRequest("POST", "/dangerous", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-XSRF-TOKEN", "valid-token")
	req.Header.Set("Cookie", "XSRF-TOKEN=valid-token")
	req.Header.Set("Origin", "https://evil-site.com")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// ──── Tests: originIsAllowed helper function ──────────────────────────────

func TestOriginIsAllowed_ExactMatch(t *testing.T) {
	restore := setCORSOrigins("http://localhost:3000")
	defer restore()

	assert.True(t, originIsAllowed("http://localhost:3000"))
}

func TestOriginIsAllowed_CaseInsensitive(t *testing.T) {
	restore := setCORSOrigins("http://localhost:3000")
	defer restore()

	assert.True(t, originIsAllowed("HTTP://LOCALHOST:3000"))
}

func TestOriginIsAllowed_EmptyOrigin(t *testing.T) {
	restore := setCORSOrigins("http://localhost:3000")
	defer restore()

	assert.False(t, originIsAllowed(""))
}

func TestOriginIsAllowed_NoMatch(t *testing.T) {
	restore := setCORSOrigins("http://localhost:3000")
	defer restore()

	assert.False(t, originIsAllowed("https://attacker.com"))
}

func TestOriginIsAllowed_EmptyConfigUseDefault(t *testing.T) {
	restore := setCORSOrigins("")
	defer restore()

	// Should fall back to default origins
	assert.True(t, originIsAllowed("http://localhost:3000"))
	assert.True(t, originIsAllowed("http://localhost:5174"))
	assert.True(t, originIsAllowed("https://jvalleyverse.web.id"))
	assert.False(t, originIsAllowed("https://evil.com"))
}

func TestOriginIsAllowed_MultipleOrigins(t *testing.T) {
	restore := setCORSOrigins("https://app.com,https://admin.com")
	defer restore()

	assert.True(t, originIsAllowed("https://app.com"))
	assert.True(t, originIsAllowed("https://admin.com"))
	assert.False(t, originIsAllowed("https://other.com"))
}

func TestOriginIsAllowed_TrailingSlash(t *testing.T) {
	restore := setCORSOrigins("http://localhost:3000")
	defer restore()

	assert.True(t, originIsAllowed("http://localhost:3000/"))
}

// ──── Tests: Integration — end-to-end flow ───────────────────────────────

func TestXSRF_Integration_NoTokenButValidOrigin(t *testing.T) {
	restore := setCORSOrigins("https://jvalleyverse.web.id")
	defer restore()

	app := fiber.New()
	app.Use(XSRFProtection())
	app.Post("/api/showcases", func(c *fiber.Ctx) error {
		return c.Status(201).JSON(fiber.Map{"id": "sc-123"})
	})

	// Simulate a real browser SPA request: no XSRF cookie/header but valid Origin
	req := httptest.NewRequest("POST", "/api/showcases", strings.NewReader(`{"title":"My Project"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://jvalleyverse.web.id")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "sc-123", body["id"])
}

func TestXSRF_Integration_NoOriginNoCookie(t *testing.T) {
	app := setupXSRFApp()

	// No Origin, no Referer, no XSRF cookies — should be blocked
	req := httptest.NewRequest("POST", "/dangerous", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 403, resp.StatusCode)
}
