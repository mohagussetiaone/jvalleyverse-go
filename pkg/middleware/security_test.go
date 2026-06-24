package middleware

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ──── helpers ─────────────────────────────────────────────────────────────

// setupScraperApp returns a Fiber app with ScraperGuard on content routes
// and a health route without ScraperGuard.
func setupScraperApp() *fiber.App {
	app := fiber.New()

	// Content endpoint with ScraperGuard
	app.Get("/api/courses", ScraperGuard(), func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"data": []string{"course1"}})
	})
	app.Get("/api/lessons/:id", ScraperGuard(), func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"title": "Go Basics"})
	})
	app.Get("/api/showcases", ScraperGuard(), func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"data": []string{"showcase1"}})
	})

	// Health endpoint without ScraperGuard (should always work)
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"status": "ok"})
	})

	return app
}

// TestScraperGuard_EmptyUserAgent blocks empty UA.
func TestScraperGuard_EmptyUserAgent(t *testing.T) {
	app := setupScraperApp()

	req := httptest.NewRequest("GET", "/api/courses", nil)
	// Deliberately no User-Agent set

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 403, resp.StatusCode)

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Contains(t, body["error"], "scraping")
}

// TestScraperGuard_BlocksKnownScrapers tests that known scraper
// User-Agents are blocked with 403.
func TestScraperGuard_BlocksKnownScrapers(t *testing.T) {
	app := setupScraperApp()

	scrapers := []struct {
		name string
		ua   string
	}{
		{"python-requests", "python-requests/2.31.0"},
		{"aiohttp", "aiohttp/3.9.0"},
		{"scrapy", "Scrapy/2.11.0"},
		{"curl", "curl/8.4.0"},
		{"wget", "Wget/1.21.4"},
		{"libcurl", "libcurl/7.88.1"},
		{"okhttp", "okhttp/4.12.0"},
		{"httpx", "httpx/1.1.0"},
		{"HttpClient", "Apache-HttpClient/4.5.14"},
		{"PostmanRuntime", "PostmanRuntime/7.36.0"},
		{"insomnia", "insomnia/2023.5.8"},
		{"Java", "Java/17.0.9"},
		{"ruby", "Ruby/3.2.2"},
		{"faraday", "Faraday/v2.7.11"},
		{"generic bot", "Mozilla/5.0 bot v1.0"},
		{"generic spider", "MySpider/1.0"},
		{"generic crawler", "Crawler/2.0"},
	}

	for _, s := range scrapers {
		t.Run(s.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/courses", nil)
			req.Header.Set("User-Agent", s.ua)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, 403, resp.StatusCode,
				"expected 403 for User-Agent: %s", s.ua)

			var body map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&body)
			assert.Contains(t, body["error"], "scraping")
		})
	}
}

// TestScraperGuard_AllowsBrowsers tests that real browser User-Agents
// are allowed through.
func TestScraperGuard_AllowsBrowsers(t *testing.T) {
	app := setupScraperApp()

	browsers := []struct {
		name string
		ua   string
	}{
		{"Chrome Windows", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
		{"Firefox", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0"},
		{"Safari macOS", "Mozilla/5.0 (Macintosh; Intel Mac OS X 14_2) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Safari/605.1.15"},
		{"Edge", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0"},
		{"Mobile Chrome", "Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36"},
		{"Mobile Safari", "Mozilla/5.0 (iPhone; CPU iPhone OS 17_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Mobile/15E148 Safari/604.1"},
	}

	for _, b := range browsers {
		t.Run(b.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/courses", nil)
			req.Header.Set("User-Agent", b.ua)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode,
				"expected 200 for browser User-Agent: %s", b.ua)

			var body map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&body)
			assert.NotEmpty(t, body["data"])
		})
	}
}

// TestScraperGuard_AllowsSearchEngineBots verifies that legitimate
// search-engine crawlers are NOT blocked.
func TestScraperGuard_AllowsSearchEngineBots(t *testing.T) {
	app := setupScraperApp()

	bots := []struct {
		name string
		ua   string
	}{
		{"Googlebot", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"},
		{"Bingbot", "Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)"},
		{"Yahoo Slurp", "Mozilla/5.0 (compatible; Yahoo! Slurp; http://help.yahoo.com/help/us/ysearch/slurp)"},
		{"DuckDuckBot", "DuckDuckBot-Https/1.1; (+https://duckduckgo.com/duckduckbot)"},
		{"Baiduspider", "Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)"},
		{"YandexBot", "Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)"},
		{"Sogou", "Sogou web spider/4.0 (+http://www.sogou.com/docs/help/webmasters.htm#07)"},
		{"Facebook", "facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)"},
		{"Twitterbot", "Twitterbot/1.0"},
		{"LinkedInBot", "LinkedInBot/1.0 (compatible; Mozilla/5.0; +http://www.linkedin.com)"},
		{"WhatsApp", "WhatsApp/2.23.23.81"},
		{"TelegramBot", "TelegramBot (like TwitterBot)"},
		{"Discordbot", "Discordbot/2.0 (+https://discord.com)"},
	}

	for _, b := range bots {
		t.Run(b.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/courses", nil)
			req.Header.Set("User-Agent", b.ua)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode,
				"expected 200 for bot User-Agent: %s", b.ua)
		})
	}
}

// TestScraperGuard_CaseInsensitivity verifies that scraper detection is
// case-insensitive.
func TestScraperGuard_CaseInsensitivity(t *testing.T) {
	app := setupScraperApp()

	tests := []struct {
		name string
		ua   string
	}{
		{"uppercase CURL", "CURL/8.4.0"},
		{"mixed Python-Requests", "Python-Requests/2.31.0"},
		{"uppercase POSTMAN", "POSTMANRUNTIME/7.36.0"},
		{"mixed Scrapy", "ScRaPy/2.11.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/courses", nil)
			req.Header.Set("User-Agent", tt.ua)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, 403, resp.StatusCode,
				"expected 403 for User-Agent: %s", tt.ua)
		})
	}
}

// TestScraperGuard_HealthEndpointExempt verifies that health endpoints
// (without ScraperGuard) work even with scraper User-Agents.
func TestScraperGuard_HealthEndpointExempt(t *testing.T) {
	app := setupScraperApp()

	req := httptest.NewRequest("GET", "/api/health", nil)
	req.Header.Set("User-Agent", "curl/8.4.0")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode,
		"health endpoint should work even with scraper UA")

	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "ok", body["status"])
}

// TestScraperGuard_LessonEndpoint also uses ScraperGuard.
func TestScraperGuard_LessonEndpoint(t *testing.T) {
	app := setupScraperApp()

	// Browser → allowed
	req := httptest.NewRequest("GET", "/api/lessons/123", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Scraper → blocked
	req2 := httptest.NewRequest("GET", "/api/lessons/123", nil)
	req2.Header.Set("User-Agent", "python-requests/2.31.0")
	resp2, err := app.Test(req2)
	require.NoError(t, err)
	assert.Equal(t, 403, resp2.StatusCode)
}

// TestScraperGuard_ShowcaseEndpoint also uses ScraperGuard.
func TestScraperGuard_ShowcaseEndpoint(t *testing.T) {
	app := setupScraperApp()

	// Browser → allowed
	req := httptest.NewRequest("GET", "/api/showcases", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Firefox/121")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Scraper → blocked
	req2 := httptest.NewRequest("GET", "/api/showcases", nil)
	req2.Header.Set("User-Agent", "Wget/1.21.4")
	resp2, err := app.Test(req2)
	require.NoError(t, err)
	assert.Equal(t, 403, resp2.StatusCode)
}
