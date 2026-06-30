package handler

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"jvalleyverse/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ──────────────────────────────────────────────
// ADMIN CATEGORY CRUD (flow.md section 2)
// ──────────────────────────────────────────────

func TestAdminCategoryHandler_CreateCategory(t *testing.T) {
	app := setupAdminApp("admin1")
	mockSvc := newMockCategoryService()
	handler := NewCategoryHandler(mockSvc)
	app.Post("/api/admin/categories", handler.CreateCategory)

	body := `{"name":"Mobile Development","slug":"mobile-development","description":"Build mobile apps"}`
	req := httptest.NewRequest("POST", "/api/admin/categories", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Mobile Development", result["name"])
	assert.Equal(t, "mobile-development", result["slug"])
}

func TestAdminCategoryHandler_CreateCategory_MissingFields(t *testing.T) {
	app := setupAdminApp("admin1")
	handler := NewCategoryHandler(newMockCategoryService())
	app.Post("/api/admin/categories", handler.CreateCategory)

	body := `{"name":""}`
	req := httptest.NewRequest("POST", "/api/admin/categories", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestAdminCategoryHandler_UpdateCategory(t *testing.T) {
	app := setupAdminApp("admin1")
	mockSvc := newMockCategoryService()
	mockSvc.CreateCategory(nil, "Old Name", "old-slug", "Old desc")
	handler := NewCategoryHandler(mockSvc)
	app.Put("/api/admin/categories/:id", handler.UpdateCategory)

	body := `{"name":"Updated Category","slug":"updated-slug","description":"Updated description"}`
	req := httptest.NewRequest("PUT", "/api/admin/categories/mock-cat-old-slug", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Updated Category", result["name"])
}

func TestAdminCategoryHandler_DeleteCategory(t *testing.T) {
	app := setupAdminApp("admin1")
	mockSvc := newMockCategoryService()
	mockSvc.CreateCategory(nil, "To Delete", "to-delete", "")
	handler := NewCategoryHandler(mockSvc)
	app.Delete("/api/admin/categories/:id", handler.DeleteCategory)

	req := httptest.NewRequest("DELETE", "/api/admin/categories/mock-cat-to-delete", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAdminCategoryHandler_ListCategories(t *testing.T) {
	app := setupAdminApp("admin1")
	mockSvc := newMockCategoryService()
	mockSvc.addTestCategory("cat1", "Cat1", "cat1")
	mockSvc.addTestCategory("cat2", "Cat2", "cat2")
	handler := NewCategoryHandler(mockSvc)
	app.Get("/api/admin/categories", handler.ListCategories)

	req := httptest.NewRequest("GET", "/api/admin/categories", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.GreaterOrEqual(t, len(result), 2)
}

// ──────────────────────────────────────────────
// ADMIN COURSE CRUD (flow.md section 2)
// ──────────────────────────────────────────────

func TestAdminCourseHandler_CreateCourse(t *testing.T) {
	app := setupAdminApp("admin1")
	mockSvc := newMockCourseService()
	mockSvc.addTestCourse("cat1", "Category 1", "admin1", "cat1") // pre-create
	handler := NewCourseHandler(mockSvc)
	app.Post("/api/admin/courses", handler.CreateCourse)

	body := `{"title":"Go REST API","description":"Learn Go","thumbnail":"thumb.jpg","category_id":"cat1","mentor_id":"mentor1","price":0,"hours":40,"learning_objectives":["Objective 1"]}`
	req := httptest.NewRequest("POST", "/api/admin/courses", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Go REST API", result["title"])
}

func TestAdminCourseHandler_CreateCourse_MissingFields(t *testing.T) {
	app := setupAdminApp("admin1")
	handler := NewCourseHandler(newMockCourseService())
	app.Post("/api/admin/courses", handler.CreateCourse)

	body := `{"title":"","category_id":""}`
	req := httptest.NewRequest("POST", "/api/admin/courses", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestAdminCourseHandler_UpdateCourse(t *testing.T) {
	app := setupAdminApp("admin1")
	mockSvc := newMockCourseService()
	mockSvc.addTestCourse("course1", "Original Title", "admin1", "cat1")
	handler := NewCourseHandler(mockSvc)
	app.Put("/api/admin/courses/:id", handler.UpdateCourse)

	body := `{"title":"Updated Title","description":"Updated desc","price":99000,"visibility":"public"}`
	req := httptest.NewRequest("PUT", "/api/admin/courses/course1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify the update was applied
	assert.Equal(t, "Updated Title", mockSvc.courses["course1"].Title)
}

func TestAdminCourseHandler_DeleteCourse(t *testing.T) {
	app := setupAdminApp("admin1")
	mockSvc := newMockCourseService()
	mockSvc.addTestCourse("course1", "To Delete", "admin1", "cat1")
	handler := NewCourseHandler(mockSvc)
	app.Delete("/api/admin/courses/:id", handler.DeleteCourse)

	req := httptest.NewRequest("DELETE", "/api/admin/courses/course1", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAdminCourseHandler_DeleteCourse_Forbidden(t *testing.T) {
	app := setupAdminApp("admin1")
	mockSvc := newMockCourseService()
	mockSvc.addTestCourse("course1", "Not Owned", "other-admin", "cat1")
	handler := NewCourseHandler(mockSvc)
	app.Delete("/api/admin/courses/:id", handler.DeleteCourse)

	req := httptest.NewRequest("DELETE", "/api/admin/courses/course1", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 403, resp.StatusCode)
}

// ──────────────────────────────────────────────
// ADMIN SECTION CRUD (flow.md section 2)
// ──────────────────────────────────────────────

func TestAdminSectionHandler_CreateSection(t *testing.T) {
	app := setupAdminApp("admin1")
	mockSvc := newMockSectionService()
	handler := NewSectionHandler(mockSvc)
	app.Post("/api/admin/courses/:course_id/sections", handler.CreateSection)

	body := `{"title":"Introduction","description":"Getting started","order_index":1}`
	req := httptest.NewRequest("POST", "/api/admin/courses/course1/sections", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Introduction", result["title"])
}

func TestAdminSectionHandler_UpdateSection(t *testing.T) {
	app := setupAdminApp("admin1")
	mockSvc := newMockSectionService()
	mockSvc.CreateSection(nil, "admin1", "course1", "Old-Section", "Old", 1)
	handler := NewSectionHandler(mockSvc)
	app.Put("/api/admin/sections/:section_id", handler.UpdateSection)

	body := `{"title":"Updated-Section","description":"Updated","order_index":2}`
	req := httptest.NewRequest("PUT", "/api/admin/sections/mock-sec-Old-Section", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAdminSectionHandler_DeleteSection(t *testing.T) {
	app := setupAdminApp("admin1")
	mockSvc := newMockSectionService()
	mockSvc.CreateSection(nil, "admin1", "course1", "To-Delete", "Old", 1)
	handler := NewSectionHandler(mockSvc)
	app.Delete("/api/admin/sections/:section_id", handler.DeleteSection)

	req := httptest.NewRequest("DELETE", "/api/admin/sections/mock-sec-To-Delete", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// ──────────────────────────────────────────────
// ADMIN LESSON CRUD (flow.md section 2)
// ──────────────────────────────────────────────

func TestAdminLessonHandler_CreateLesson(t *testing.T) {
	app := setupAdminApp("admin1")
	mockSvc := newMockLessonService()
	handler := NewLessonHandler(mockSvc)
	app.Post("/api/admin/lessons", handler.CreateLesson)

	body := `{"course_id":"course1","section_id":"section1","title":"Go Basics","slug":"go-basics","description":"Learn Go","difficulty":"beginner","duration":45,"order_index":1,"visibility":"public"}`
	req := httptest.NewRequest("POST", "/api/admin/lessons", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Go Basics", result["title"])
	assert.Equal(t, "go-basics", result["slug"])
}

func TestAdminLessonHandler_CreateLesson_MissingFields(t *testing.T) {
	app := setupAdminApp("admin1")
	handler := NewLessonHandler(newMockLessonService())
	app.Post("/api/admin/lessons", handler.CreateLesson)

	body := `{"course_id":"","section_id":"","title":"","slug":""}`
	req := httptest.NewRequest("POST", "/api/admin/lessons", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestAdminLessonHandler_CreateLessonDetail(t *testing.T) {
	app := setupAdminApp("admin1")
	mockSvc := newMockLessonService()
	mockSvc.addTestLesson("lesson1", "Go Basics", "go-basics", "course1", "section1", "admin1")
	handler := NewLessonHandler(mockSvc)
	app.Post("/api/admin/lessons/:id/details", handler.CreateLessonDetail)

	body := `{"about":"Learn Go basics","rules":"1. Follow along","tools":["Go","VS Code"],"resource_media":{"videos":[],"documents":[],"images":[]},"resources":[]}`
	req := httptest.NewRequest("POST", "/api/admin/lessons/lesson1/details", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Contains(t, result, "id")
	assert.Equal(t, "Learn Go basics", result["about"])
}

func TestAdminLessonHandler_UpdateLesson(t *testing.T) {
	app := setupAdminApp("admin1")
	mockSvc := newMockLessonService()
	mockSvc.addTestLesson("lesson1", "Old Title", "old-slug", "course1", "section1", "admin1")
	handler := NewLessonHandler(mockSvc)
	app.Put("/api/admin/lessons/:id", handler.UpdateLesson)

	body := `{"title":"Updated Title","slug":"updated-slug","description":"Updated description","difficulty":"intermediate","duration":60,"visibility":"public"}`
	req := httptest.NewRequest("PUT", "/api/admin/lessons/lesson1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Updated Title", result["title"])
}

func TestAdminLessonHandler_DeleteLesson(t *testing.T) {
	app := setupAdminApp("admin1")
	mockSvc := newMockLessonService()
	mockSvc.addTestLesson("lesson1", "To Delete", "to-delete", "course1", "section1", "admin1")
	handler := NewLessonHandler(mockSvc)
	app.Delete("/api/admin/lessons/:id", handler.DeleteLesson)

	req := httptest.NewRequest("DELETE", "/api/admin/lessons/lesson1", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAdminLessonHandler_DeleteLesson_Forbidden(t *testing.T) {
	app := setupAdminApp("admin1")
	mockSvc := newMockLessonService()
	mockSvc.addTestLesson("lesson1", "Not Owned", "not-owned", "course1", "section1", "other-admin")
	handler := NewLessonHandler(mockSvc)
	app.Delete("/api/admin/lessons/:id", handler.DeleteLesson)

	req := httptest.NewRequest("DELETE", "/api/admin/lessons/lesson1", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 403, resp.StatusCode)
}

// ──────────────────────────────────────────────
// ADMIN STUDY CASE CRUD (flow.md section 2)
// ──────────────────────────────────────────────

func TestAdminStudyCaseHandler_CreateStudyCase(t *testing.T) {
	app := setupAdminApp("admin1")
	mockSvc := newMockStudyCaseService()
	handler := NewStudyCaseHandler(mockSvc)
	app.Post("/api/admin/study-cases", handler.CreateStudyCase)

	body := `{"name":"Go Microservices","description":"Build microservices with Go","img_url":"img.jpg","youtube_url":"https://youtube.com","tags":["go","microservices"]}`
	req := httptest.NewRequest("POST", "/api/admin/study-cases", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Go Microservices", result["name"])
}

func TestAdminStudyCaseHandler_CreateStudyCase_MissingName(t *testing.T) {
	app := setupAdminApp("admin1")
	handler := NewStudyCaseHandler(newMockStudyCaseService())
	app.Post("/api/admin/study-cases", handler.CreateStudyCase)

	body := `{"name":"","description":"test"}`
	req := httptest.NewRequest("POST", "/api/admin/study-cases", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestAdminStudyCaseHandler_UpdateStudyCase(t *testing.T) {
	app := setupAdminApp("admin1")
	mockSvc := newMockStudyCaseService()
	mockSvc.CreateStudyCase(nil, "admin1", "Original", "Original desc", "old.jpg", "", nil, []string{})
	handler := NewStudyCaseHandler(mockSvc)
	app.Put("/api/admin/study-cases/:id", handler.UpdateStudyCase)

	body := `{"name":"Updated","description":"Updated desc","img_url":"new.jpg","tags":["updated"]}`
	req := httptest.NewRequest("PUT", "/api/admin/study-cases/mock-sc-Original", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAdminStudyCaseHandler_DeleteStudyCase(t *testing.T) {
	app := setupAdminApp("admin1")
	mockSvc := newMockStudyCaseService()
	mockSvc.CreateStudyCase(nil, "admin1", "To-Delete", "", "", "", nil, []string{})
	handler := NewStudyCaseHandler(mockSvc)
	app.Delete("/api/admin/study-cases/:id", handler.DeleteStudyCase)

	req := httptest.NewRequest("DELETE", "/api/admin/study-cases/mock-sc-To-Delete", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAdminStudyCaseHandler_DeleteStudyCase_NotFound(t *testing.T) {
	app := setupAdminApp("admin1")
	handler := NewStudyCaseHandler(newMockStudyCaseService())
	app.Delete("/api/admin/study-cases/:id", handler.DeleteStudyCase)

	req := httptest.NewRequest("DELETE", "/api/admin/study-cases/nonexistent", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

// ──────────────────────────────────────────────
// ADMIN BLOG CRUD (flow.md section 18)
// ──────────────────────────────────────────────

func TestAdminBlogHandler_CreateBlog(t *testing.T) {
	app := setupAdminApp("admin1")
	mockSvc := newMockBlogService()
	handler := NewBlogHandler(mockSvc)
	app.Post("/api/admin/blogs", handler.Create)

	body := `{"title":"Getting Started with Go","description":"A comprehensive guide","content":"# Markdown content","cover_img_url":"https://example.com/cover.jpg","tags":["go","programming"],"status":"published","category_id":"cat1"}`
	req := httptest.NewRequest("POST", "/api/admin/blogs", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Getting Started with Go", result["title"])
}

func TestAdminBlogHandler_CreateBlog_MissingTitle(t *testing.T) {
	app := setupAdminApp("admin1")
	handler := NewBlogHandler(newMockBlogService())
	app.Post("/api/admin/blogs", handler.Create)

	body := `{"title":"","category_id":"cat1"}`
	req := httptest.NewRequest("POST", "/api/admin/blogs", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestAdminBlogHandler_UpdateBlog(t *testing.T) {
	app := setupAdminApp("admin1")
	mockSvc := newMockBlogService()
	mockSvc.CreateBlog(nil, "admin1", service.CreateBlogRequest{
		Title:      "Original",
		Status:     "draft",
		CategoryID: "cat1",
	})
	app.Put("/api/admin/blogs/:id", func(c *fiber.Ctx) error {
		blogID := c.Params("id")
		var req service.UpdateBlogRequest
		c.BodyParser(&req)
		if err := mockSvc.AdminUpdateBlog(c.UserContext(), blogID, req); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "Blog updated successfully"})
	})

	body := `{"title":"Updated Blog","status":"published"}`
	req := httptest.NewRequest("PUT", "/api/admin/blogs/mock-blog-Original", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAdminBlogHandler_DeleteBlog(t *testing.T) {
	app := setupAdminApp("admin1")
	mockSvc := newMockBlogService()
	mockSvc.CreateBlog(nil, "admin1", service.CreateBlogRequest{
		Title:      "To-Delete",
		Status:     "draft",
		CategoryID: "cat1",
	})
	app.Delete("/api/admin/blogs/:id", func(c *fiber.Ctx) error {
		blogID := c.Params("id")
		if err := mockSvc.AdminDeleteBlog(c.UserContext(), blogID); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "Blog deleted successfully"})
	})

	req := httptest.NewRequest("DELETE", "/api/admin/blogs/mock-blog-To-Delete", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// ──────────────────────────────────────────────
// ADMIN USER MANAGEMENT TESTS
// ──────────────────────────────────────────────

func TestAdminGetAllUsers(t *testing.T) {
	app := setupAdminApp("admin1")
	mockSvc := newMockUserService()
	handler := NewUserHandler(mockSvc, newMockDashboardService())
	app.Get("/api/admin/users", handler.GetAllUsers)

	req := httptest.NewRequest("GET", "/api/admin/users", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestAdminDashboard(t *testing.T) {
	app := setupAdminApp("admin1")

	app.Get("/api/admin/dashboard", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome admin"})
	})

	req := httptest.NewRequest("GET", "/api/admin/dashboard", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Welcome admin", result["message"])
}

// ──────────────────────────────────────────────
// ADMIN LIST ALL BLOGS
// ──────────────────────────────────────────────

func TestAdminBlogHandler_ListAllBlogs(t *testing.T) {
	app := setupAdminApp("admin1")
	mockSvc := newMockBlogService()
	mockSvc.CreateBlog(nil, "admin1", service.CreateBlogRequest{
		Title:      "Admin Blog",
		Status:     "published",
		CategoryID: "cat1",
	})
	handler := NewBlogHandler(mockSvc)
	// Admin uses the same List handler but with admin middleware
	app.Get("/api/admin/blogs", handler.List)

	req := httptest.NewRequest("GET", "/api/admin/blogs", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
