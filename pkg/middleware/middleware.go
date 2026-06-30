package middleware

import (
	"strings"
	"time"

	"jvalleyverse/pkg/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/golang-jwt/jwt/v5"
)

// ==================== CORS ====================
func SetupCORS() fiber.Handler {
	origins := config.AppConfig.CORSOrigins
	if origins == "" {
		origins = "http://localhost:3000,http://localhost:5173,http://localhost:5174,https://jvalleyverse.web.id"
	}
	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-XSRF-TOKEN",
		AllowCredentials: true,
	})
}

// ── rate limit tiers (requests per IP per minute) ─────────────────────

const (
	RateLimitGlobal  = 200 // General browsing — looser than before
	RateLimitAuth    = 10  // Login / register — brute force protection
	RateLimitContent = 60  // Public content endpoints — anti-scraping
)

// ==================== RATE LIMITER (global) ====================
func RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        RateLimitGlobal,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	})
}

// ==================== AUTH RATE LIMITER (stricter for login/register) ====================
func AuthRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        RateLimitAuth,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	})
}

// ==================== CONTENT RATE LIMITER (anti-scraping) ====================
// Applied to public content endpoints (courses, lessons, showcases).
// Lower limit than global because these are the most-targeted endpoints for scrapers.
func ContentRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        RateLimitContent,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		SkipSuccessfulRequests: false,
	})
}

// ==================== OPTIONAL JWT AUTH (does not reject if missing) ====================
func OptionalJWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			// No token, continue without user context
			return c.Next()
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Next()
		}

		tokenString := parts[1]

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Next()
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return c.Next()
		}

		roleStr, ok := claims["role"].(string)
		if !ok {
			return c.Next()
		}

		c.Locals("userID", userIDStr)
		c.Locals("role", strings.ToLower(strings.TrimSpace(roleStr)))

		return c.Next()
	}
}

// ==================== JWT AUTH ====================
func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization format",
			})
		}

		tokenString := parts[1]

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// ===== SAFE CAST - user_id should be string for CUID =====
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user_id in token",
			})
		}

		roleStr, ok := claims["role"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid role in token",
			})
		}

		// ===== NORMALIZE =====
		roleStr = strings.ToLower(strings.TrimSpace(roleStr))

		c.Locals("userID", userIDStr)
		c.Locals("role", roleStr)

		return c.Next()
	}
}

// ==================== ROLE GUARD ====================
func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleVal := c.Locals("role")

		if roleVal == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Role not found",
			})
		}

		userRole, ok := roleVal.(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Invalid role type",
			})
		}

		userRole = strings.ToLower(strings.TrimSpace(userRole))
		requiredRole = strings.ToLower(requiredRole)

		if userRole != requiredRole && userRole != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied",
			})
		}

		return c.Next()
	}
}

// SecurityHeaders adds security-related HTTP headers
func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		return c.Next()
	}
}

// XSRF Protection sederhana (cek header X-XSRF-TOKEN)
// Untuk production lebih baik gunakan double submit cookie pattern
func XSRFProtection() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Abaikan untuk method GET/HEAD/OPTIONS
		if c.Method() == "GET" || c.Method() == "HEAD" || c.Method() == "OPTIONS" {
			return c.Next()
		}
		token := c.Get("X-XSRF-TOKEN")
		// Bandingkan dengan cookie "XSRF-TOKEN" (harus dikirim client)
		cookieToken := c.Cookies("XSRF-TOKEN")
		if token == "" || cookieToken == "" || token != cookieToken {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "XSRF token invalid"})
		}
		return c.Next()
	}
}
