package middleware

import (
	"regexp"

	"github.com/gofiber/fiber/v2"
)

// ── scraped user-agent patterns ──────────────────────────────────────────

var scraperPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)python-requests`),
	regexp.MustCompile(`(?i)aiohttp`),
	regexp.MustCompile(`(?i)scrapy`),
	regexp.MustCompile(`(?i)curl\b`),
	regexp.MustCompile(`(?i)wget\b`),
	regexp.MustCompile(`(?i)libcurl`),
	regexp.MustCompile(`(?i)okhttp`),
	regexp.MustCompile(`(?i)httpx`),
	regexp.MustCompile(`(?i)HttpClient`),
	regexp.MustCompile(`(?i)PostmanRuntime`),
	regexp.MustCompile(`(?i)insomnia`),
	regexp.MustCompile(`(?i)Java/`),
	regexp.MustCompile(`(?i)ruby`),
	regexp.MustCompile(`(?i)faraday`),
	regexp.MustCompile(`(?i)(?:bot|spider|crawler)\b`),
}

// legitimateBotPatterns — bots we ALLOW despite matching the generic pattern above.
var legitimateBotPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)Googlebot`),
	regexp.MustCompile(`(?i)Bingbot`),
	regexp.MustCompile(`(?i)Slurp`), // Yahoo
	regexp.MustCompile(`(?i)DuckDuckBot`),
	regexp.MustCompile(`(?i)Baiduspider`),
	regexp.MustCompile(`(?i)YandexBot`),
	regexp.MustCompile(`(?i)Sogou`),
	regexp.MustCompile(`(?i)facebookexternalhit`),
	regexp.MustCompile(`(?i)Twitterbot`),
	regexp.MustCompile(`(?i)LinkedInBot`),
	regexp.MustCompile(`(?i)WhatsApp`),
	regexp.MustCompile(`(?i)TelegramBot`),
	regexp.MustCompile(`(?i)Discordbot`),
}

// isScraperUserAgent returns true when the User-Agent looks like an
// automated scraper / CLI tool rather than a real browser.
func isScraperUserAgent(ua string) bool {
	if ua == "" {
		return true // empty UA → always block
	}

	// First check for known-good bots so we don't block search engines.
	for _, pat := range legitimateBotPatterns {
		if pat.MatchString(ua) {
			return false
		}
	}

	for _, pat := range scraperPatterns {
		if pat.MatchString(ua) {
			return true
		}
	}

	return false
}

// ── middleware constructors ──────────────────────────────────────────────

// ScraperGuard blocks requests that appear to come from automated
// scrapers / bots (based on User-Agent) when they hit public content
// endpoints. Legitimate search-engine bots are allowed through.
func ScraperGuard() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ua := c.Get("User-Agent")
		if isScraperUserAgent(ua) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied: automated scraping is not allowed",
			})
		}
		return c.Next()
	}
}
