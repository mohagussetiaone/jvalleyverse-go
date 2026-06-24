package middleware

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"jvalleyverse/pkg/redis"

	"github.com/alicebob/miniredis/v2"
	"github.com/gofiber/fiber/v2"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ──── helpers ─────────────────────────────────────────────────────────────

// setupApp returns a Fiber app with IdempotencyMiddleware and a simple POST handler.
func setupApp() *fiber.App {
	app := fiber.New()
	app.Use(IdempotencyMiddleware())
	app.Post("/test", func(c *fiber.Ctx) error {
		return c.Status(201).JSON(fiber.Map{
			"id":   "order-123",
			"done": true,
		})
	})
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"status": "ok"})
	})
	return app
}

// setupAppWithErrorHandler returns an app whose POST handler returns 422.
func setupAppWithErrorHandler() *fiber.App {
	app := fiber.New()
	app.Use(IdempotencyMiddleware())
	app.Post("/test", func(c *fiber.Ctx) error {
		return c.Status(422).JSON(fiber.Map{"error": "unprocessable"})
	})
	return app
}

// setupAppWithMethod returns an app with the given method + handler.
func setupAppWithMethod(method, path string, handler fiber.Handler) *fiber.App {
	app := fiber.New()
	app.Use(IdempotencyMiddleware())
	app.Add(method, path, handler)
	return app
}

// withRedis sets up a miniredis instance and temporarily replaces the global
// redis.Client / redis.IsConnected. It returns the miniredis for assertions.
func withRedis(t *testing.T) *miniredis.Miniredis {
	t.Helper()

	mr, err := miniredis.Run()
	require.NoError(t, err)

	// Save originals
	origClient := redis.Client
	origConnected := redis.IsConnected

	// Replace with test client pointing to miniredis
	redis.Client = goredis.NewClient(&goredis.Options{
		Addr: mr.Addr(),
	})
	redis.IsConnected = true

	t.Cleanup(func() {
		redis.Client.Close()
		redis.Client = origClient
		redis.IsConnected = origConnected
		mr.Close()
	})

	return mr
}

// withoutRedis ensures Redis is unavailable for the duration of the test.
func withoutRedis(t *testing.T) {
	t.Helper()

	origClient := redis.Client
	origConnected := redis.IsConnected

	redis.Client = nil
	redis.IsConnected = false

	t.Cleanup(func() {
		redis.Client = origClient
		redis.IsConnected = origConnected
	})
}

// ──── tests ───────────────────────────────────────────────────────────────

func TestIdempotency_SkipGet(t *testing.T) {
	withoutRedis(t)
	app := setupApp()

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// X-Idempotency-Replayed should NEVER appear on GET
	assert.Empty(t, resp.Header.Get("X-Idempotency-Replayed"))
}

func TestIdempotency_SkipHead(t *testing.T) {
	withoutRedis(t)
	app := setupApp()

	req := httptest.NewRequest("HEAD", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestIdempotency_SkipOptions(t *testing.T) {
	withoutRedis(t)
	app := setupApp()

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	// OPTIONS should not be intercepted — Fiber returns 405 when no OPTIONS handler is registered
	assert.Equal(t, 405, resp.StatusCode, "OPTIONS should pass through middleware (not intercepted)")
}

func TestIdempotency_NoKey(t *testing.T) {
	withoutRedis(t)
	app := setupApp()

	req := httptest.NewRequest("POST", "/test", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
	assert.Empty(t, resp.Header.Get("X-Idempotency-Replayed"))
}

func TestIdempotency_InvalidUUID(t *testing.T) {
	withoutRedis(t)
	app := setupApp()

	tests := []struct {
		name string
		key  string
	}{
		{"not a UUID", "not-a-uuid"},
		{"missing dashes", "550e8400e29b41d4a716446655440000"},
		{"too short", "abc"},
		{"numeric", "12345"},
		{"partial uuid", "550e8400-e29b-41d4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/test", strings.NewReader(`{}`))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Idempotency-Key", tt.key)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, 400, resp.StatusCode)

			var body map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&body)
			assert.Contains(t, body["error"], "valid UUID")
		})
	}
}

func TestIdempotency_ValidUUIDNoRedis(t *testing.T) {
	withoutRedis(t)
	app := setupApp()

	req := httptest.NewRequest("POST", "/test", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", "550e8400-e29b-41d4-a716-446655440000")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
	assert.Empty(t, resp.Header.Get("X-Idempotency-Replayed"))
}



func TestIdempotency_CacheAndReplay(t *testing.T) {
	withRedis(t)

	app := setupApp()
	uuid := "550e8400-e29b-41d4-a716-446655440000"
	body := `{"hello":"world"}`
	reqBody := strings.NewReader(body)

	// ── 1st request: cache miss → processed normally ──
	req1 := httptest.NewRequest("POST", "/test", reqBody)
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Idempotency-Key", uuid)

	resp1, err := app.Test(req1)
	require.NoError(t, err)
	assert.Equal(t, 201, resp1.StatusCode)
	assert.Equal(t, "application/json", resp1.Header.Get("Content-Type"))

	var data1 map[string]interface{}
	json.NewDecoder(resp1.Body).Decode(&data1)
	assert.Equal(t, "order-123", data1["id"])

	// ── 2nd request: cache hit → replayed from cache ──
	req2 := httptest.NewRequest("POST", "/test", strings.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Idempotency-Key", uuid)

	resp2, err := app.Test(req2)
	require.NoError(t, err)
	assert.Equal(t, 201, resp2.StatusCode)
	assert.Equal(t, "application/json", resp2.Header.Get("Content-Type"),
		"replayed response should have Content-Type set")
	assert.Equal(t, "true", resp2.Header.Get("X-Idempotency-Replayed"),
		"replayed response should have X-Idempotency-Replayed header")

	var data2 map[string]interface{}
	json.NewDecoder(resp2.Body).Decode(&data2)
	assert.Equal(t, data1, data2, "replayed response should be identical to first response")
}

func TestIdempotency_DifferentKeys(t *testing.T) {
	withRedis(t)

	app := setupApp()

	keyA := "11111111-1111-4111-8111-111111111111"
	keyB := "22222222-2222-4222-8222-222222222222"

	// Request with key A
	reqA1 := httptest.NewRequest("POST", "/test", strings.NewReader(`{}`))
	reqA1.Header.Set("Content-Type", "application/json")
	reqA1.Header.Set("Idempotency-Key", keyA)
	respA1, err := app.Test(reqA1)
	require.NoError(t, err)
	assert.Equal(t, 201, respA1.StatusCode)

	// Request with key B
	reqB := httptest.NewRequest("POST", "/test", strings.NewReader(`{}`))
	reqB.Header.Set("Content-Type", "application/json")
	reqB.Header.Set("Idempotency-Key", keyB)
	respB, err := app.Test(reqB)
	require.NoError(t, err)
	assert.Equal(t, 201, respB.StatusCode)
	assert.Empty(t, respB.Header.Get("X-Idempotency-Replayed"), "key B should be processed fresh")

	// Replay key A → should get cached response (not fresh)
	reqA2 := httptest.NewRequest("POST", "/test", strings.NewReader(`{}`))
	reqA2.Header.Set("Content-Type", "application/json")
	reqA2.Header.Set("Idempotency-Key", keyA)
	respA2, err := app.Test(reqA2)
	require.NoError(t, err)
	assert.Equal(t, 201, respA2.StatusCode)
	assert.Equal(t, "true", respA2.Header.Get("X-Idempotency-Replayed"), "key A should replay")
}

func TestIdempotency_CacheOnly2xx(t *testing.T) {
	mr := withRedis(t)

	app := setupAppWithErrorHandler()
	uuid := "550e8400-e29b-41d4-a716-446655440000"

	// First request returns 422 (not cached)
	req1 := httptest.NewRequest("POST", "/test", strings.NewReader(`{}`))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Idempotency-Key", uuid)

	resp1, err := app.Test(req1)
	require.NoError(t, err)
	assert.Equal(t, 422, resp1.StatusCode)

	// Verify NOT cached in Redis
	_, err = mr.Get("idempotent:" + uuid)
	assert.Error(t, err, "4xx response should NOT be cached in Redis")

	// Second request: since first was not cached, handler runs again
	req2 := httptest.NewRequest("POST", "/test", strings.NewReader(`{}`))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Idempotency-Key", uuid)

	resp2, err := app.Test(req2)
	require.NoError(t, err)
	assert.Equal(t, 422, resp2.StatusCode)
	assert.Empty(t, resp2.Header.Get("X-Idempotency-Replayed"),
		"second request should NOT be replayed because 4xx was not cached")
}

func TestIdempotency_PutMethod(t *testing.T) {
	withRedis(t)

	app := setupAppWithMethod("PUT", "/test", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"updated": true})
	})

	req := httptest.NewRequest("PUT", "/test", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", "550e8400-e29b-41d4-a716-446655440000")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestIdempotency_DeleteMethod(t *testing.T) {
	withRedis(t)

	app := setupAppWithMethod("DELETE", "/test", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"deleted": true})
	})

	req := httptest.NewRequest("DELETE", "/test", strings.NewReader(``))
	req.Header.Set("Idempotency-Key", "550e8400-e29b-41d4-a716-446655440000")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestIdempotency_PatchMethod(t *testing.T) {
	withRedis(t)

	app := setupAppWithMethod("PATCH", "/test", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"patched": true})
	})

	req := httptest.NewRequest("PATCH", "/test", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", "550e8400-e29b-41d4-a716-446655440000")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestIdempotency_ConsecutiveCacheHits(t *testing.T) {
	withRedis(t)

	app := setupApp()
	uuid := "550e8400-e29b-41d4-a716-446655440000"

	// Request 3 times with same key — should all succeed
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("POST", "/test", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", uuid)

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)

		if i == 0 {
			assert.Empty(t, resp.Header.Get("X-Idempotency-Replayed"),
				"first request should NOT be replayed")
		} else {
			assert.Equal(t, "true", resp.Header.Get("X-Idempotency-Replayed"),
				"subsequent requests SHOULD be replayed")
		}
	}
}
