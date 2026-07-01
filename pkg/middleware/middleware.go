package middleware

import (
	"net/url"
	"strings"
	"time"

	"jvalleyverse/pkg/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/golang-jwt/jwt/v5"
)

// allowedOriginsSet is a quick-lookup set of allowed origins (from CORS config).
// Rebuilt on every call so it picks up ENV changes without restart.
func getAllowedOrigins() []string {
	raw := config.AppConfig.CORSOrigins
	if raw == "" {
		raw = "http://localhost:3000,http://localhost:3001,http://localhost:5173,http://localhost:5174,https://jvalleyverse.web.id"
	}
	return strings.Split(raw, ",")
}

// originIsAllowed checks if the given Origin or Referer matches any configured CORS origin.
// Referer URLs may include paths (e.g. "http://localhost:3000/some-page"),
// so we parse them and compare only the origin (scheme + host + port).
func originIsAllowed(origin string) bool {
	if origin == "" {
		return false
	}

	// If the value looks like a full URL (contains a path), parse out just the origin.
	// Origin header never has a path, but Referer does.
	if strings.Contains(origin, "//") && strings.Count(origin, "/") >= 3 {
		if parsed, err := url.Parse(origin); err == nil {
			origin = parsed.Scheme + "://" + parsed.Host
		}
	}

	// Strip trailing slash for comparison
	origin = strings.TrimRight(origin, "/")

	for _, allowed := range getAllowedOrigins() {
		allowed = strings.TrimSpace(allowed)
		if strings.EqualFold(origin, allowed) {
			return true
		}
	}
	return false
}

// ==================== CORS ====================
func SetupCORS() fiber.Handler {
	origins := config.AppConfig.CORSOrigins
	if origins == "" {
		origins = "http://localhost:3000,http://localhost:3001,http://localhost:5173,http://localhost:5174,https://jvalleyverse.web.id"
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

// XSRF Protection with Origin/Referer fallback.
//
// Strategy:
//  1. Standard double-submit cookie check: X-XSRF-TOKEN header must match XSRF-TOKEN cookie.
//  2. If cookie check fails, fall back to verifying the Origin (or Referer) header
//     matches the configured CORS origins.  This is a well-known CSRF defence for SPAs
//     because the browser will not allow JavaScript on another origin to spoof this header.
//
// This means SPA clients that cannot reliably read the XSRF-TOKEN cookie will still get
// protection through the Same-Origin policy.
func XSRFProtection() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Abaikan untuk method GET/HEAD/OPTIONS
		if c.Method() == "GET" || c.Method() == "HEAD" || c.Method() == "OPTIONS" {
			return c.Next()
		}

		token := c.Get("X-XSRF-TOKEN")
		cookieToken := c.Cookies("XSRF-TOKEN")

		// ── Check 1: Double-submit cookie ──
		if token != "" && cookieToken != "" && token == cookieToken {
			return c.Next() // XSRF token valid
		}

		// ── Check 2: Origin / Referer fallback ──
		origin := c.Get("Origin")
		if origin == "" {
			origin = c.Get("Referer")
		}
		if origin != "" && originIsAllowed(origin) {
			// Origin matches configured CORS origins → allowed
			return c.Next()
		}

		// Both checks failed → reject
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "XSRF token invalid"})
	}
}
