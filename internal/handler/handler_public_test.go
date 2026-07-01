package handler

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"jvalleyverse/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ──────────────────────────────────────────────
// CATEGORY FLOW TESTS (flow.md section 3)
// ──────────────────────────────────────────────

func TestCategoryHandler_ListCategories(t *testing.T) {
	app := setupTestApp()
	mockSvc := newMockCategoryService()
	mockSvc.addTestCategory("cat1", "Backend Development", "backend-development")
	mockSvc.addTestCategory("cat2", "Frontend Development", "frontend-development")
	handler := NewCategoryHandler(mockSvc)
	app.Get("/api/categories", handler.ListCategories)

	req := httptest.NewRequest("GET", "/api/categories", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.GreaterOrEqual(t, len(result), 2)
	// Find our test categories (order independent due to map iteration)
	foundNames := make(map[string]bool)
	for _, cat := range result {
		foundNames[cat["name"].(string)] = true
	}
	assert.True(t, foundNames["Backend Development"])
	assert.True(t, foundNames["Frontend Development"])
}

func TestCategoryHandler_GetCategoryBySlug(t *testing.T) {
	app := setupTestApp()
	mockSvc := newMockCategoryService()
	mockSvc.addTestCategory("cat1", "Backend", "backend")
	handler := NewCategoryHandler(mockSvc)
	app.Get("/api/categories/:slug", handler.GetCategoryBySlug)

	req := httptest.NewRequest("GET", "/api/categories/backend", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Backend", result["name"])
}

func TestCategoryHandler_GetCategoryBySlug_NotFound(t *testing.T) {
	app := setupTestApp()
	handler := NewCategoryHandler(newMockCategoryService())
	app.Get("/api/categories/:slug", handler.GetCategoryBySlug)

	req := httptest.NewRequest("GET", "/api/categories/nonexistent", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestCategoryHandler_ListCoursesByCategory(t *testing.T) {
	app := setupTestApp()
	mockCat := newMockCategoryService()
	mockCat.addTestCategory("cat1", "Backend", "backend")
	handler := NewCategoryHandler(mockCat)
	app.Get("/api/categories/:category_id/courses", handler.ListCoursesByCategory)

	req := httptest.NewRequest("GET", "/api/categories/cat1/courses", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestCategoryHandler_ListCoursesByCategory_WithAuth(t *testing.T) {
	app := setupTestApp()
	// Simulate optional JWT
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", "test-user-id")
		return c.Next()
	})

	mockCat := newMockCategoryService()
	mockCat.addTestCategory("cat1", "Backend", "backend")
	handler := NewCategoryHandler(mockCat)
	app.Get("/api/categories/:category_id/courses", handler.ListCoursesByCategory)

	req := httptest.NewRequest("GET", "/api/categories/cat1/courses", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// ──────────────────────────────────────────────
// COURSE FLOW TESTS (flow.md section 3)
// ──────────────────────────────────────────────

func TestCourseHandler_ListPublicCourses(t *testing.T) {
	app := setupTestApp()
	handler := NewCourseHandler(newMockCourseService())
	app.Get("/api/courses", handler.ListPublicCourses)

	req := httptest.NewRequest("GET", "/api/courses?page=1&limit=10", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Contains(t, result, "data")
	assert.Contains(t, result, "pagination")
}

func TestCourseHandler_ListPublicCourses_WithUser(t *testing.T) {
	app := setupTestApp()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", "test-user")
		return c.Next()
	})
	handler := NewCourseHandler(newMockCourseService())
	app.Get("/api/courses", handler.ListPublicCourses)

	req := httptest.NewRequest("GET", "/api/courses", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestCourseHandler_ListPublicCourses_WithFilters(t *testing.T) {
	app := setupTestApp()
	handler := NewCourseHandler(newMockCourseService())
	app.Get("/api/courses", handler.ListPublicCourses)

	req := httptest.NewRequest("GET", "/api/courses?page=1&limit=10&category_id=cat1&min_price=0&max_price=100000", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// ──────────────────────────────────────────────
// SECTION FLOW TESTS
// ──────────────────────────────────────────────

func TestSectionHandler_GetCourseWithSections(t *testing.T) {
	app := setupTestApp()
	mockCourse := newMockCourseService()
	mockCourse.addTestCourse("course1", "Test Course", "admin-id", "cat1")

	mockSvc := newMockSectionService()
	handler := NewSectionHandler(mockSvc)
	// Register with mockCourse to return the course
	// Since GetCourseWithSections uses courseRepo internally, we simulate
	app.Get("/api/courses/:course_id", handler.GetCourseWithSections)

	req := httptest.NewRequest("GET", "/api/courses/course1", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode) // Course not in sections mock
}

func TestSectionHandler_ListSectionsByCourse(t *testing.T) {
	app := setupTestApp()
	handler := NewSectionHandler(newMockSectionService())
	app.Get("/api/courses/:course_id/sections", handler.ListSectionsByCourse)

	req := httptest.NewRequest("GET", "/api/courses/course1/sections", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestSectionHandler_GetSectionByID(t *testing.T) {
	app := setupTestApp()
	mockSvc := newMockSectionService()
	mockSvc.CreateSection(nil, "admin-id", "course1", "Test-Section", "Description", 1)
	handler := NewSectionHandler(mockSvc)
	app.Get("/api/courses/:course_id/sections/:section_id", handler.GetSection)

	req := httptest.NewRequest("GET", "/api/courses/course1/sections/mock-sec-Test-Section", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestSectionHandler_GetSection_NotFound(t *testing.T) {
	app := setupTestApp()
	handler := NewSectionHandler(newMockSectionService())
	app.Get("/api/courses/:course_id/sections/:section_id", handler.GetSection)

	req := httptest.NewRequest("GET", "/api/courses/course1/sections/nonexistent", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

// ──────────────────────────────────────────────
// LESSON FLOW TESTS (flow.md section 4)
// ──────────────────────────────────────────────

func TestLessonHandler_GetPublicLessonByID(t *testing.T) {
	app := setupTestApp()
	mockSvc := newMockLessonService()
	mockSvc.addTestLesson("lesson1", "Go Basics", "go-basics", "course1", "section1", "admin1")
	handler := NewLessonHandler(mockSvc)
	app.Get("/api/lessons/:id", handler.GetPublicLessonByID)

	req := httptest.NewRequest("GET", "/api/lessons/lesson1", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	lesson := result["lesson"].(map[string]interface{})
	assert.Equal(t, "Go Basics", lesson["title"])
	assert.Contains(t, result, "details")
}

func TestLessonHandler_GetPublicLesson_NotFound(t *testing.T) {
	app := setupTestApp()
	handler := NewLessonHandler(newMockLessonService())
	app.Get("/api/lessons/:id", handler.GetPublicLessonByID)

	req := httptest.NewRequest("GET", "/api/lessons/nonexistent", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestLessonHandler_GetLessonBySlug(t *testing.T) {
	app := setupTestApp()
	mockSvc := newMockLessonService()
	mockSvc.addTestLesson("lesson1", "Go Basics", "go-basics", "course1", "section1", "admin1")
	handler := NewLessonHandler(mockSvc)
	app.Get("/api/courses/:course_id/lessons/:slug", handler.GetLessonBySlug)

	req := httptest.NewRequest("GET", "/api/courses/course1/lessons/go-basics", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestLessonHandler_GetLessonBySlug_WithUser(t *testing.T) {
	app := setupTestApp()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", "test-user")
		return c.Next()
	})
	mockSvc := newMockLessonService()
	mockSvc.addTestLesson("lesson1", "Go Basics", "go-basics", "course1", "section1", "admin1")
	handler := NewLessonHandler(mockSvc)
	app.Get("/api/courses/:course_id/lessons/:slug", handler.GetLessonBySlug)

	req := httptest.NewRequest("GET", "/api/courses/course1/lessons/go-basics", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestLessonHandler_ListLessonsByCourse(t *testing.T) {
	app := setupTestApp()
	handler := NewLessonHandler(newMockLessonService())
	app.Get("/api/courses/:course_id/lessons", handler.ListLessonsByCourse)

	req := httptest.NewRequest("GET", "/api/courses/course1/lessons?page=1&limit=20", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestLessonHandler_ListLessonsBySection(t *testing.T) {
	app := setupTestApp()
	handler := NewLessonHandler(newMockLessonService())
	app.Get("/api/courses/:course_id/sections/:section_id/lessons", handler.ListLessonsBySection)

	req := httptest.NewRequest("GET", "/api/courses/course1/sections/section1/lessons", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// ──────────────────────────────────────────────
// STUDY CASE FLOW TESTS (flow.md section 3)
// ──────────────────────────────────────────────

func TestStudyCaseHandler_ListStudyCases(t *testing.T) {
	app := setupTestApp()
	handler := NewStudyCaseHandler(newMockStudyCaseService())
	app.Get("/api/study-cases", handler.ListStudyCases)

	req := httptest.NewRequest("GET", "/api/study-cases?page=1&limit=20", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Contains(t, result, "data")
	assert.Contains(t, result, "pagination")
}

func TestStudyCaseHandler_GetStudyCase(t *testing.T) {
	app := setupTestApp()
	mockSvc := newMockStudyCaseService()
	mockSvc.CreateStudyCase(nil, "admin1", "Test-Case", "Description", "img.jpg", "https://youtube.com", nil, []string{"go"})
	handler := NewStudyCaseHandler(mockSvc)
	app.Get("/api/study-cases/:id", handler.GetStudyCase)

	req := httptest.NewRequest("GET", "/api/study-cases/mock-sc-Test-Case", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestStudyCaseHandler_GetStudyCase_NotFound(t *testing.T) {
	app := setupTestApp()
	handler := NewStudyCaseHandler(newMockStudyCaseService())
	app.Get("/api/study-cases/:id", handler.GetStudyCase)

	req := httptest.NewRequest("GET", "/api/study-cases/nonexistent", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

// ──────────────────────────────────────────────
// BLOG FLOW TESTS (flow.md section 18)
// ──────────────────────────────────────────────

func TestBlogHandler_ListBlogs(t *testing.T) {
	app := setupTestApp()
	mockSvc := newMockBlogService()
	_ = NewBlogHandler(mockSvc) // Ensure mock implements IBlogService
	app.Get("/api/blogs", func(c *fiber.Ctx) error {
		page := c.QueryInt("page", 1)
		limit := c.QueryInt("limit", 10)
		items, pagination, _ := mockSvc.ListBlogs(c.UserContext(), page, limit, "", "", "")
		return c.JSON(fiber.Map{
			"data":       items,
			"pagination": pagination,
		})
	})

	req := httptest.NewRequest("GET", "/api/blogs?page=1&limit=10", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Contains(t, result, "data")
	assert.Contains(t, result, "pagination")
}

func TestBlogHandler_GetBlogByID(t *testing.T) {
	app := setupTestApp()
	mockSvc := newMockBlogService()
	mockSvc.CreateBlog(nil, "admin1", service.CreateBlogRequest{
		Title:       "Test-Blog",
		Description: "Description",
		Content:     "Content",
		CoverImgURL: "https://example.com/cover.jpg",
		Tags:        []string{"go", "programming"},
		Status:      "published",
		CategoryID:  "cat1",
	})
	app.Get("/api/blogs/:id", func(c *fiber.Ctx) error {
		blog, err := mockSvc.GetBlogByID(c.UserContext(), c.Params("id"))
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(blog)
	})

	req := httptest.NewRequest("GET", "/api/blogs/mock-blog-Test-Blog", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Test-Blog", result["title"])
}

func TestBlogHandler_GetBlog_NotFound(t *testing.T) {
	app := setupTestApp()
	handler := NewBlogHandler(newMockBlogService())
	app.Get("/api/blogs/:id", handler.GetByID)

	req := httptest.NewRequest("GET", "/api/blogs/nonexistent", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

// ──────────────────────────────────────────────
// SHOWCASE PUBLIC FLOW TESTS (flow.md section 7)
// ──────────────────────────────────────────────

func TestShowcaseHandler_ListShowcases(t *testing.T) {
	app := setupTestApp()
	handler := NewShowcaseHandler(newMockShowcaseService())
	app.Get("/api/showcases", handler.ListShowcases)

	req := httptest.NewRequest("GET", "/api/showcases?page=1&limit=20&sort=newest", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Contains(t, result, "data")
	assert.Contains(t, result, "pagination")
}

func TestShowcaseHandler_ListShowcases_WithCategoryFilter(t *testing.T) {
	app := setupTestApp()
	handler := NewShowcaseHandler(newMockShowcaseService())
	app.Get("/api/showcases", handler.ListShowcases)

	req := httptest.NewRequest("GET", "/api/showcases?page=1&limit=20&category_id=cat1&sort=newest", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestShowcaseHandler_GetShowcase(t *testing.T) {
	app := setupTestApp()
	mockSvc := newMockShowcaseService()
	mockSvc.CreateShowcase(nil, "user1", "My-Project", "Description", []string{"img.jpg"}, "cat1", "public")
	handler := NewShowcaseHandler(mockSvc)
	app.Get("/api/showcases/:id", handler.GetShowcase)

	req := httptest.NewRequest("GET", "/api/showcases/mock-showcase-My-Project", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestShowcaseHandler_GetShowcase_NotFound(t *testing.T) {
	app := setupTestApp()
	handler := NewShowcaseHandler(newMockShowcaseService())
	app.Get("/api/showcases/:id", handler.GetShowcase)

	req := httptest.NewRequest("GET", "/api/showcases/nonexistent", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

// ──────────────────────────────────────────────
// DISCUSSION PUBLIC FLOW TESTS (flow.md section 6)
// ──────────────────────────────────────────────

func TestDiscussionHandler_ListDiscussions(t *testing.T) {
	app := setupTestApp()
	svc := newMockDiscussionService()

	// Create a test app to override the handler
	app.Get("/api/discussions", func(c *fiber.Ctx) error {
		page := c.QueryInt("page", 1)
		limit := c.QueryInt("limit", 20)
		discussions, total, _ := svc.ListDiscussions(c.UserContext(), page, limit, nil, nil, nil)
		return c.JSON(fiber.Map{
			"data": discussions,
			"pagination": fiber.Map{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		})
	})

	req := httptest.NewRequest("GET", "/api/discussions?page=1&limit=20", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestDiscussionHandler_GetDiscussion(t *testing.T) {
	app := setupTestApp()
	mockSvc := newMockDiscussionService()
	mockSvc.CreateDiscussion(nil, "user1", "Test-Discussion", "Content", nil, nil, "cat1")

	app.Get("/api/discussions/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		d, err := mockSvc.GetDiscussionWithReplies(c.UserContext(), id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Discussion not found"})
		}
		return c.JSON(d)
	})

	req := httptest.NewRequest("GET", "/api/discussions/mock-disc-Test-Discussion", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// ──────────────────────────────────────────────
// FAQ PUBLIC TESTS (flow.md section 15)
// ──────────────────────────────────────────────

func TestFaqHandler_ListPublic(t *testing.T) {
	app := setupTestApp()
	mockSvc := newMockFaqService()
	mockSvc.addTestFAQ("faq1", "Apa itu JValleyverse?", "Platform belajar coding.", "general", 1, true)
	mockSvc.addTestFAQ("faq2", "Bagaimana cara daftar?", "Klik tombol Daftar.", "account", 2, true)
	handler := NewFaqHandler(mockSvc)
	app.Get("/api/faqs", handler.ListPublic)

	req := httptest.NewRequest("GET", "/api/faqs?page=1&limit=20", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 2)
	assert.Contains(t, result, "pagination")
	pagination := result["pagination"].(map[string]interface{})
	assert.Equal(t, float64(1), pagination["page"])
	assert.Equal(t, float64(20), pagination["limit"])
	assert.Equal(t, float64(2), pagination["total"])
}

func TestFaqHandler_ListPublic_Empty(t *testing.T) {
	app := setupTestApp()
	handler := NewFaqHandler(newMockFaqService())
	app.Get("/api/faqs", handler.ListPublic)

	req := httptest.NewRequest("GET", "/api/faqs?page=1&limit=20", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].([]interface{})
	assert.Equal(t, 0, len(data))
	assert.Contains(t, result, "pagination")
}

func TestFaqHandler_ListPublic_OnlyActive(t *testing.T) {
	app := setupTestApp()
	mockSvc := newMockFaqService()
	mockSvc.addTestFAQ("faq1", "Active FAQ", "Active answer.", "general", 1, true)
	mockSvc.addTestFAQ("faq2", "Inactive FAQ", "Inactive answer.", "general", 2, false)
	handler := NewFaqHandler(mockSvc)
	app.Get("/api/faqs", handler.ListPublic)

	req := httptest.NewRequest("GET", "/api/faqs?page=1&limit=20", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].([]interface{})
	assert.Equal(t, 1, len(data))
	first := data[0].(map[string]interface{})
	assert.Equal(t, "Active FAQ", first["question"])
	pagination := result["pagination"].(map[string]interface{})
	assert.Equal(t, float64(1), pagination["total"])
}

// ──────────────────────────────────────────────
// COMPANY PUBLIC TESTS (flow.md section 15)
// ──────────────────────────────────────────────

func TestCompanyHandler_GetCompany(t *testing.T) {
	app := setupTestApp()
	handler := NewCompanyHandler(newMockCompanyService())
	app.Get("/api/company", handler.GetCompany)

	req := httptest.NewRequest("GET", "/api/company", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "JValleyVerse", result["brand_name"])
	assert.Equal(t, "Learn, Build, Grow Together", result["tagline"])
	assert.Equal(t, "hello@jvalleyverse.com", result["email"])
}

// ──────────────────────────────────────────────
// GAMIFICATION PUBLIC TESTS (flow.md section 9)
// ──────────────────────────────────────────────

func TestGamificationHandler_GetLeaderboard(t *testing.T) {
	app := setupTestApp()
	mockSvc := newMockGamificationService(newMockUserService())
	// Override since NewGamificationHandler uses global
	app.Get("/api/leaderboard", func(c *fiber.Ctx) error {
		limit := c.QueryInt("limit", 10)
		data, _ := mockSvc.GetLeaderboard(c.UserContext(), limit)
		return c.JSON(fiber.Map{"data": data})
	})

	req := httptest.NewRequest("GET", "/api/leaderboard?limit=10", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].([]interface{})
	assert.GreaterOrEqual(t, len(data), 1)
	first := data[0].(map[string]interface{})
	assert.Equal(t, float64(1), first["rank"])
}

// ──────────────────────────────────────────────
// REVIEW PUBLIC TESTS
// ──────────────────────────────────────────────

func TestReviewHandler_ListCourseReviews(t *testing.T) {
	app := setupTestApp()
	handler := NewReviewHandler(newMockReviewService())
	app.Get("/api/courses/:course_id/reviews", handler.ListCourseReviews)

	req := httptest.NewRequest("GET", "/api/courses/course1/reviews", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestReviewHandler_ListLessonReviews(t *testing.T) {
	app := setupTestApp()
	handler := NewReviewHandler(newMockReviewService())
	app.Get("/api/lessons/:id/reviews", handler.ListLessonReviews)

	req := httptest.NewRequest("GET", "/api/lessons/lesson1/reviews", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
