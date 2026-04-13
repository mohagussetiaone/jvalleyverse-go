package middleware

import (
	"log"
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
	return cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-XSRF-TOKEN",
		AllowCredentials: true,
	})
}

// ==================== RATE LIMITER ====================
func RateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	})
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

		// ===== SAFE CAST =====
		userIDFloat, ok := claims["user_id"].(float64)
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

		// ===== DEBUG LOG =====
		log.Println("USER ID:", userIDFloat)
		log.Println("ROLE:", roleStr)

		c.Locals("userID", uint(userIDFloat))
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

		log.Println("CHECK ROLE:", userRole)

		if userRole != requiredRole && userRole != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied",
			})
		}

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