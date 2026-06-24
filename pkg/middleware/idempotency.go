package middleware

import (
	"encoding/json"
	"regexp"
	"time"

	"jvalleyverse/pkg/redis"

	"github.com/gofiber/fiber/v2"
)

// idempotentResponse stores the status code and body for replay.
type idempotentResponse struct {
	StatusCode int             `json:"status_code"`
	Data       json.RawMessage `json:"data"`
}

// uuidRegex validates UUID v4 format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx).
var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// cacheTTL is how long a completed idempotent response is kept.
const cacheTTL = 24 * time.Hour

// IdempotencyMiddleware ensures safe retries by deduplicating mutation requests.
//
// Clients MUST send the Idempotency-Key header with a UUID v4 value.
// On the first request the response is cached in Redis (TTL 24h).
// On retry the cached response is returned immediately without processing.
//
// If Redis is unavailable the middleware silently passes through,
// so the API degrades gracefully when Redis is down.
//
// Usage:
//
//	app.Post("/orders", middleware.IdempotencyMiddleware(), handler.PlaceOrder)
func IdempotencyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip read-only methods
		switch c.Method() {
		case fiber.MethodGet, fiber.MethodHead, fiber.MethodOptions:
			return c.Next()
		}

		key := c.Get("Idempotency-Key")
		if key == "" {
			return c.Next() // No idempotency requested
		}

		// Validate UUID format
		if !uuidRegex.MatchString(key) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Idempotency-Key must be a valid UUID v4 (e.g. 550e8400-e29b-41d4-a716-446655440000)",
			})
		}

		cacheKey := "idempotent:" + key

		// --- Cache hit → replay stored response ---
		if redis.IsAvailable() {
			cached, err := redis.Client.Get(c.Context(), cacheKey).Result()
			if err == nil && cached != "" {
				var resp idempotentResponse
				if json.Unmarshal([]byte(cached), &resp) == nil {
					c.Status(resp.StatusCode)
					c.Set("Content-Type", "application/json")
					c.Set("X-Idempotency-Replayed", "true")
					return c.Send(resp.Data)
				}
			}
		}

		// --- Cache miss → process normally ---
		if err := c.Next(); err != nil {
			return err
		}

		// --- Cache the response only for 2xx ---
		statusCode := c.Response().StatusCode()
		if statusCode >= 200 && statusCode < 300 && redis.IsAvailable() {
			body := c.Response().Body()
			resp := idempotentResponse{
				StatusCode: statusCode,
				Data:       body,
			}
			if data, err := json.Marshal(resp); err == nil {
				redis.Client.Set(c.Context(), cacheKey, data, cacheTTL)
			}
		}

		return nil
	}
}
