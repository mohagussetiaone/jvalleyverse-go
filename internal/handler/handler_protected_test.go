package handler

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ──────────────────────────────────────────────
// SHOWCASE CRUD FLOW TESTS (flow.md section 7)
// ──────────────────────────────────────────────

func TestShowcaseHandler_CreateShowcase(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	mockSvc := newMockShowcaseService()
	handler := NewShowcaseHandler(mockSvc)
	app.Post("/api/showcases", handler.Create)

	body := `{"title":"My Project","description":"A cool project","media_urls":["https://example.com/img.jpg"],"category_id":"cat1","visibility":"public"}`
	req := httptest.NewRequest("POST", "/api/showcases", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "My Project", result["title"])
	assert.Equal(t, "user1", result["user_id"])
}

func TestShowcaseHandler_CreateShowcase_MissingTitle(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	handler := NewShowcaseHandler(newMockShowcaseService())
	app.Post("/api/showcases", handler.Create)

	body := `{"description":"No title","category_id":"cat1"}`
	req := httptest.NewRequest("POST", "/api/showcases", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestShowcaseHandler_UpdateShowcase(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	mockSvc := newMockShowcaseService()
	mockSvc.CreateShowcase(nil, "user1", "Original", "Original desc", []string{"img.jpg"}, "cat1", "public")
	handler := NewShowcaseHandler(mockSvc)
	app.Put("/api/showcases/:id", handler.Update)

	body := `{"title":"Updated Title","description":"Updated desc","visibility":"public"}`
	req := httptest.NewRequest("PUT", "/api/showcases/mock-showcase-Original", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Updated Title", result["title"])
}

func TestShowcaseHandler_UpdateShowcase_Forbidden(t *testing.T) {
	app := setupProtectedApp("user2", "user") // different user
	mockSvc := newMockShowcaseService()
	mockSvc.CreateShowcase(nil, "user1", "Original", "Original desc", []string{"img.jpg"}, "cat1", "public")
	handler := NewShowcaseHandler(mockSvc)
	app.Put("/api/showcases/:id", handler.Update)

	body := `{"title":"Hacked!","visibility":"public"}`
	req := httptest.NewRequest("PUT", "/api/showcases/mock-showcase-Original", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 403, resp.StatusCode)
}

func TestShowcaseHandler_DeleteShowcase(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	mockSvc := newMockShowcaseService()
	mockSvc.CreateShowcase(nil, "user1", "To-Delete", "Desc", []string{}, "cat1", "public")
	handler := NewShowcaseHandler(mockSvc)
	app.Delete("/api/showcases/:id", handler.Delete)

	req := httptest.NewRequest("DELETE", "/api/showcases/mock-showcase-To-Delete", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestShowcaseHandler_DeleteShowcase_Forbidden(t *testing.T) {
	app := setupProtectedApp("user2", "user")
	mockSvc := newMockShowcaseService()
	mockSvc.CreateShowcase(nil, "user1", "To-Delete", "Desc", []string{}, "cat1", "public")
	handler := NewShowcaseHandler(mockSvc)
	app.Delete("/api/showcases/:id", handler.Delete)

	req := httptest.NewRequest("DELETE", "/api/showcases/mock-showcase-To-Delete", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 403, resp.StatusCode)
}

func TestShowcaseHandler_LikeShowcase(t *testing.T) {
	app := setupProtectedApp("user2", "user")
	mockSvc := newMockShowcaseService()
	mockSvc.CreateShowcase(nil, "user1", "Liked-Showcase", "Desc", []string{}, "cat1", "public")
	handler := NewShowcaseHandler(mockSvc)
	app.Post("/api/showcases/:id/like", handler.Like)

	req := httptest.NewRequest("POST", "/api/showcases/mock-showcase-Liked-Showcase/like", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Liked successfully", result["message"])
}

func TestShowcaseHandler_UnlikeShowcase(t *testing.T) {
	app := setupProtectedApp("user2", "user")
	mockSvc := newMockShowcaseService()
	mockSvc.CreateShowcase(nil, "user1", "Unliked-Showcase", "Desc", []string{}, "cat1", "public")
	mockSvc.LikeShowcase(nil, "user2", "mock-showcase-Unliked-Showcase")
	handler := NewShowcaseHandler(mockSvc)
	app.Delete("/api/showcases/:id/like", handler.Unlike)

	req := httptest.NewRequest("DELETE", "/api/showcases/mock-showcase-Unliked-Showcase/like", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestShowcaseHandler_ListMyShowcases(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	handler := NewShowcaseHandler(newMockShowcaseService())
	app.Get("/api/users/me/showcases", handler.ListMyShowcases)

	req := httptest.NewRequest("GET", "/api/users/me/showcases", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// ──────────────────────────────────────────────
// DISCUSSION CRUD FLOW TESTS (flow.md section 6)
// ──────────────────────────────────────────────

func TestDiscussionHandler_CreateDiscussion(t *testing.T) {
	app := setupProtectedApp("user1", "user")

	app.Post("/api/discussions", func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(string)
		var input struct {
			Title      string  `json:"title"`
			Content    string  `json:"content"`
			LessonID   *string `json:"lesson_id"`
			CategoryID string  `json:"category_id"`
		}
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
		}
		mock := newMockDiscussionService()
		d, err := mock.CreateDiscussion(c.UserContext(), userID, input.Title, input.Content, input.LessonID, nil, input.CategoryID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(201).JSON(d)
	})

	body := `{"title":"My Question","content":"Need help with Go","lesson_id":"lesson1","category_id":"cat1"}`
	req := httptest.NewRequest("POST", "/api/discussions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestDiscussionHandler_UpdateDiscussion(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	mockSvc := newMockDiscussionService()
	mockSvc.CreateDiscussion(nil, "user1", "Original", "Original content", nil, nil, "cat1")

	app.Put("/api/discussions/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		userID := c.Locals("userID").(string)
		var input struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}
		c.BodyParser(&input)
		if err := mockSvc.UpdateDiscussion(c.UserContext(), id, userID, input.Title, input.Content); err != nil {
			return c.Status(403).JSON(fiber.Map{"error": "You do not own this discussion"})
		}
		return c.JSON(fiber.Map{"message": "Discussion updated"})
	})

	body := `{"title":"Updated","content":"Updated content"}`
	req := httptest.NewRequest("PUT", "/api/discussions/mock-disc-Original", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestDiscussionHandler_DeleteDiscussion(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	mockSvc := newMockDiscussionService()
	mockSvc.CreateDiscussion(nil, "user1", "To-Delete", "Content", nil, nil, "cat1")

	app.Delete("/api/discussions/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		userID := c.Locals("userID").(string)
		role := c.Locals("role").(string)
		isAdmin := role == "admin"
		if err := mockSvc.DeleteDiscussion(c.UserContext(), id, userID, isAdmin); err != nil {
			return c.Status(403).JSON(fiber.Map{"error": "You do not own this discussion"})
		}
		return c.JSON(fiber.Map{"message": "Discussion deleted"})
	})

	req := httptest.NewRequest("DELETE", "/api/discussions/mock-disc-To-Delete", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestDiscussionHandler_CloseDiscussion(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	mockSvc := newMockDiscussionService()
	mockSvc.CreateDiscussion(nil, "user1", "To-Close", "Content", nil, nil, "cat1")

	app.Post("/api/discussions/:id/close", func(c *fiber.Ctx) error {
		id := c.Params("id")
		userID := c.Locals("userID").(string)
		if err := mockSvc.CloseDiscussion(c.UserContext(), id, userID); err != nil {
			return c.Status(403).JSON(fiber.Map{"error": "You do not own this discussion"})
		}
		return c.JSON(fiber.Map{"message": "Discussion closed"})
	})

	req := httptest.NewRequest("POST", "/api/discussions/mock-disc-To-Close/close", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestDiscussionHandler_ListMyDiscussions(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	mockSvc := newMockDiscussionService()

	app.Get("/api/users/me/discussions", func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(string)
		page := c.QueryInt("page", 1)
		limit := c.QueryInt("limit", 20)
		data, total, _ := mockSvc.ListUserDiscussions(c.UserContext(), userID, page, limit)
		return c.JSON(fiber.Map{
			"data": data,
			"pagination": fiber.Map{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		})
	})

	req := httptest.NewRequest("GET", "/api/users/me/discussions", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// ──────────────────────────────────────────────
// REPLY FLOW TESTS (flow.md section 6)
// ──────────────────────────────────────────────

func TestReplyHandler_CreateReply(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	mockSvc := newMockReplyService()

	app.Post("/api/discussions/:id/replies", func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(string)
		discussionID := c.Params("id")
		var input struct {
			Content  string  `json:"content"`
			ParentID *string `json:"parent_id"`
		}
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
		}
		reply, err := mockSvc.CreateReply(c.UserContext(), userID, discussionID, input.Content, input.ParentID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(201).JSON(reply)
	})

	body := `{"content":"Great question! Let me help."}`
	req := httptest.NewRequest("POST", "/api/discussions/disc1/replies", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestReplyHandler_UpdateReply(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	mockSvc := newMockReplyService()
	mockSvc.CreateReply(nil, "user1", "disc1", "Original-rp", nil)

	app.Put("/api/replies/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		userID := c.Locals("userID").(string)
		var input struct {
			Content string `json:"content"`
		}
		c.BodyParser(&input)
		if err := mockSvc.UpdateReply(c.UserContext(), id, userID, input.Content); err != nil {
			return c.Status(403).JSON(fiber.Map{"error": "You do not own this reply"})
		}
		return c.JSON(fiber.Map{"message": "Reply updated"})
	})

	body := `{"content":"Updated reply with better explanation"}`
	req := httptest.NewRequest("PUT", "/api/replies/mock-reply-Original-r", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestReplyHandler_DeleteReply(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	mockSvc := newMockReplyService()
	mockSvc.CreateReply(nil, "user1", "disc1", "To-delete-x", nil)

	app.Delete("/api/replies/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		userID := c.Locals("userID").(string)
		role := c.Locals("role").(string)
		if err := mockSvc.DeleteReply(c.UserContext(), id, userID, role == "admin"); err != nil {
			return c.Status(403).JSON(fiber.Map{"error": "You do not own this reply"})
		}
		return c.JSON(fiber.Map{"message": "Reply deleted"})
	})

	req := httptest.NewRequest("DELETE", "/api/replies/mock-reply-To-delete-", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestReplyHandler_LikeReply(t *testing.T) {
	app := setupProtectedApp("user2", "user")
	mockSvc := newMockReplyService()
	mockSvc.CreateReply(nil, "user1", "disc1", "Reply-to-l", nil)

	app.Post("/api/replies/:id/like", func(c *fiber.Ctx) error {
		id := c.Params("id")
		userID := c.Locals("userID").(string)
		if err := mockSvc.LikeReply(c.UserContext(), userID, id); err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Reply not found"})
		}
		return c.JSON(fiber.Map{"message": "Reply liked"})
	})

	req := httptest.NewRequest("POST", "/api/replies/mock-reply-Reply-to-l/like", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestReplyHandler_MarkBestReply(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	mockSvc := newMockReplyService()
	mockSvc.CreateReply(nil, "user2", "disc1", "Best-answer", nil)

	app.Post("/api/replies/:id/best", func(c *fiber.Ctx) error {
		id := c.Params("id")
		userID := c.Locals("userID").(string)
		var input struct {
			DiscussionID string `json:"discussion_id"`
		}
		c.BodyParser(&input)
		if err := mockSvc.MarkBestReply(c.UserContext(), id, input.DiscussionID, userID); err != nil {
			return c.Status(403).JSON(fiber.Map{"error": "Only discussion owner can mark best answer"})
		}
		return c.JSON(fiber.Map{"message": "Reply marked as best answer"})
	})

	body := `{"discussion_id":"disc1"}`
	req := httptest.NewRequest("POST", "/api/replies/mock-reply-Best-answe/best", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestReplyHandler_GetMyReplies(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	mockSvc := newMockReplyService()

	app.Get("/api/users/me/replies", func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(string)
		page := c.QueryInt("page", 1)
		limit := c.QueryInt("limit", 20)
		data, total, _ := mockSvc.ListRepliesByUser(c.UserContext(), userID, page, limit)
		return c.JSON(fiber.Map{
			"data": data,
			"pagination": fiber.Map{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		})
	})

	req := httptest.NewRequest("GET", "/api/users/me/replies", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// ──────────────────────────────────────────────
// LEARNING PROGRESS FLOW TESTS (flow.md section 4)
// ──────────────────────────────────────────────

func TestLessonHandler_StartLesson(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	mockSvc := newMockLessonService()
	mockSvc.addTestLesson("lesson1", "Go Basics", "go-basics", "course1", "section1", "admin1")
	handler := NewLessonHandler(mockSvc)
	app.Post("/api/lessons/:id/start", handler.StartLesson)

	req := httptest.NewRequest("POST", "/api/lessons/lesson1/start", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Lesson started!", result["message"])
	progress := result["progress"].(map[string]interface{})
	assert.Equal(t, "started", progress["status"])
}

func TestLessonHandler_StartLesson_NotFound(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	handler := NewLessonHandler(newMockLessonService())
	app.Post("/api/lessons/:id/start", handler.StartLesson)

	req := httptest.NewRequest("POST", "/api/lessons/nonexistent/start", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestLessonHandler_UpdateProgress(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	mockSvc := newMockLessonService()
	mockSvc.addTestLesson("lesson1", "Go Basics", "go-basics", "course1", "section1", "admin1")
	mockSvc.StartLesson(nil, "user1", "lesson1")
	handler := NewLessonHandler(mockSvc)
	app.Put("/api/lessons/:id/progress", handler.UpdateProgress)

	body := `{"progress_percentage":50,"notes":"Learning middleware"}`
	req := httptest.NewRequest("PUT", "/api/lessons/lesson1/progress", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, float64(50), result["progress_percentage"])
	assert.Equal(t, "in_progress", result["status"])
}

func TestLessonHandler_CompleteLesson(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	mockSvc := newMockLessonService()
	mockSvc.addTestLesson("lesson1", "Go Basics", "go-basics", "course1", "section1", "admin1")
	mockSvc.StartLesson(nil, "user1", "lesson1")
	handler := NewLessonHandler(mockSvc)
	app.Post("/api/lessons/:id/complete", handler.CompleteLesson)

	req := httptest.NewRequest("POST", "/api/lessons/lesson1/complete", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Lesson completed!", result["message"])
	assert.Contains(t, result, "certificate")
	assert.Contains(t, result, "achievement")
	assert.Equal(t, float64(50), result["points_awarded"])
}

func TestLessonHandler_CompleteLesson_NotStarted(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	mockSvc := newMockLessonService()
	mockSvc.addTestLesson("lesson1", "Go Basics", "go-basics", "course1", "section1", "admin1")
	handler := NewLessonHandler(mockSvc)
	app.Post("/api/lessons/:id/complete", handler.CompleteLesson)

	req := httptest.NewRequest("POST", "/api/lessons/lesson1/complete", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode) // Not started, so not found
}

// ──────────────────────────────────────────────
// ENROLLMENT FLOW TESTS (flow.md section 4)
// ──────────────────────────────────────────────

func TestCourseHandler_EnrollCourse(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	mockSvc := newMockCourseService()
	mockSvc.addTestCourse("course1", "Test Course", "admin1", "cat1")
	handler := NewCourseHandler(mockSvc)
	app.Post("/api/courses/:id/enroll", handler.EnrollCourse)

	req := httptest.NewRequest("POST", "/api/courses/course1/enroll", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Successfully enrolled in course", result["message"])
}

func TestCourseHandler_EnrollCourse_NotFound(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	handler := NewCourseHandler(newMockCourseService())
	app.Post("/api/courses/:id/enroll", handler.EnrollCourse)

	req := httptest.NewRequest("POST", "/api/courses/nonexistent/enroll", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestCourseHandler_ListEnrolledCourses(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	handler := NewCourseHandler(newMockCourseService())
	app.Get("/api/users/me/courses", handler.ListEnrolledCourses)

	req := httptest.NewRequest("GET", "/api/users/me/courses", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestCourseHandler_SetLastLesson(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	handler := NewCourseHandler(newMockCourseService())
	app.Put("/api/courses/:id/last-lesson", handler.SetLastLesson)

	body := `{"lesson_id":"lesson1"}`
	req := httptest.NewRequest("PUT", "/api/courses/course1/last-lesson", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// ──────────────────────────────────────────────
// CERTIFICATE FLOW TESTS (flow.md section 5)
// ──────────────────────────────────────────────

func TestCertificateHandler_ListCertificates(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	mockSvc := newMockCertificateService()

	app.Get("/api/users/me/certificates", func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(string)
		page := c.QueryInt("page", 1)
		limit := c.QueryInt("limit", 20)
		data, total, _ := mockSvc.ListUserCertificates(c.UserContext(), userID, page, limit)
		return c.JSON(fiber.Map{
			"data": data,
			"pagination": fiber.Map{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		})
	})

	req := httptest.NewRequest("GET", "/api/users/me/certificates", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestCertificateHandler_GetCertificateByCode(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	mockSvc := newMockCertificateService()
	mockSvc.IssueCertificate(nil, "user1", "lesson1", "CERT-test1234")

	app.Get("/api/users/me/certificates/:code", func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(string)
		role := c.Locals("role").(string)
		code := c.Params("code")
		cert, err := mockSvc.GetCertificateByCode(c.UserContext(), code, userID, role)
		if err != nil {
			return c.Status(403).JSON(fiber.Map{"error": "Forbidden: not certificate owner"})
		}
		return c.JSON(cert)
	})

	req := httptest.NewRequest("GET", "/api/users/me/certificates/CERT-test1234", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestCertificateHandler_GetCertificateForbidden(t *testing.T) {
	app := setupProtectedApp("user2", "user")
	mockSvc := newMockCertificateService()
	mockSvc.IssueCertificate(nil, "user1", "lesson1", "CERT-user1only")

	app.Get("/api/users/me/certificates/:code", func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(string)
		role := c.Locals("role").(string)
		code := c.Params("code")
		_, err := mockSvc.GetCertificateByCode(c.UserContext(), code, userID, role)
		if err != nil {
			return c.Status(403).JSON(fiber.Map{"error": "Forbidden: not certificate owner"})
		}
		return c.JSON(fiber.Map{})
	})

	req := httptest.NewRequest("GET", "/api/users/me/certificates/CERT-user1only", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 403, resp.StatusCode)
}

// ──────────────────────────────────────────────
// REVIEW CRUD FLOW TESTS (flow.md section 3)
// ──────────────────────────────────────────────

func TestReviewHandler_CreateReview(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	mockSvc := newMockReviewService()
	handler := NewReviewHandler(mockSvc)
	app.Post("/api/reviews", handler.CreateReview)

	body := `{"course_id":"course1","lesson_id":"","rating":5,"message":"Excellent course!"}`
	req := httptest.NewRequest("POST", "/api/reviews", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, float64(5), result["rating"])
}

func TestReviewHandler_CreateReview_Invalid(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	handler := NewReviewHandler(newMockReviewService())
	app.Post("/api/reviews", handler.CreateReview)

	body := `{"course_id":"course1","rating":0,"message":""}`
	req := httptest.NewRequest("POST", "/api/reviews", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestReviewHandler_UpdateReview(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	mockSvc := newMockReviewService()
	mockSvc.CreateReview(nil, "user1", "course1", "", 5, "Good course")
	created, _ := mockSvc.CreateReview(nil, "user1", "course1", "", 5, "Good course")
	handler := NewReviewHandler(mockSvc)
	app.Put("/api/reviews/:id", handler.UpdateReview)

	body := `{"rating":4,"message":"Updated review"}`
	req := httptest.NewRequest("PUT", "/api/reviews/"+created.ID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestReviewHandler_DeleteReview(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	mockSvc := newMockReviewService()
	created, _ := mockSvc.CreateReview(nil, "user1", "course1", "", 5, "Good")
	handler := NewReviewHandler(mockSvc)
	app.Delete("/api/reviews/:id", handler.DeleteReview)

	req := httptest.NewRequest("DELETE", "/api/reviews/"+created.ID, nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// ──────────────────────────────────────────────
// NOTIFICATION FLOW TESTS (flow.md section 8)
// ──────────────────────────────────────────────

func TestNotificationHandler_ListNotifications(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	handler := NewNotificationHandler(newMockNotificationService())
	app.Get("/api/users/me/notifications", handler.ListNotifications)

	req := httptest.NewRequest("GET", "/api/users/me/notifications", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestNotificationHandler_CountUnread(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	handler := NewNotificationHandler(newMockNotificationService())
	app.Get("/api/users/me/notifications/count", handler.CountUnread)

	req := httptest.NewRequest("GET", "/api/users/me/notifications/count", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, float64(5), result["unread_count"])
}

func TestNotificationHandler_MarkAsRead(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	handler := NewNotificationHandler(newMockNotificationService())
	app.Put("/api/users/me/notifications/:id/read", handler.MarkAsRead)

	req := httptest.NewRequest("PUT", "/api/users/me/notifications/notif1/read", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestNotificationHandler_MarkAllAsRead(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	handler := NewNotificationHandler(newMockNotificationService())
	app.Put("/api/users/me/notifications/read-all", handler.MarkAllAsRead)

	req := httptest.NewRequest("PUT", "/api/users/me/notifications/read-all", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestNotificationHandler_DeleteNotification(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	handler := NewNotificationHandler(newMockNotificationService())
	app.Delete("/api/users/me/notifications/:id", handler.DeleteNotification)

	req := httptest.NewRequest("DELETE", "/api/users/me/notifications/notif1", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// ──────────────────────────────────────────────
// GAMIFICATION FLOW TESTS (flow.md section 9)
// ──────────────────────────────────────────────

func TestGamificationHandler_GetLevels(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	mockSvc := newMockGamificationService(newMockUserService())

	app.Get("/api/levels", func(c *fiber.Ctx) error {
		levels := mockSvc.GetLevelInfo()
		return c.JSON(fiber.Map{"data": levels})
	})

	req := httptest.NewRequest("GET", "/api/levels", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].([]interface{})
	assert.Equal(t, 5, len(data))
	first := data[0].(map[string]interface{})
	assert.Equal(t, "Beginner", first["name"])
}

func TestGamificationHandler_GetUserPoints(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	mockUserSvc := newMockUserService()
	mockUserSvc.addTestUser("user1", "user1@test.com", "pass", "Test User", "user")
	mockSvc := newMockGamificationService(mockUserSvc)

	app.Get("/api/users/:id/points", func(c *fiber.Ctx) error {
		userID := c.Params("id")
		stats, err := mockSvc.GetUserStats(c.UserContext(), userID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "User not found"})
		}
		return c.JSON(stats)
	})

	req := httptest.NewRequest("GET", "/api/users/user1/points", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Test User", result["name"])
}

func TestGamificationHandler_GetUserPoints_NotFound(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	mockSvc := newMockGamificationService(newMockUserService())

	app.Get("/api/users/:id/points", func(c *fiber.Ctx) error {
		userID := c.Params("id")
		_, err := mockSvc.GetUserStats(c.UserContext(), userID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "User not found"})
		}
		return c.JSON(fiber.Map{})
	})

	req := httptest.NewRequest("GET", "/api/users/nonexistent/points", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

// ──────────────────────────────────────────────
// MY ITEMS TESTS
// ──────────────────────────────────────────────

func TestBlogHandler_ListMyBlogs(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	mockSvc := newMockBlogService()
	handler := NewBlogHandler(mockSvc)
	app.Get("/api/users/me/blogs", handler.ListMyBlogs)

	req := httptest.NewRequest("GET", "/api/users/me/blogs", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestStudyCaseHandler_ListMyStudyCases(t *testing.T) {
	app := setupProtectedApp("user1", "user")
	handler := NewStudyCaseHandler(newMockStudyCaseService())
	app.Get("/api/users/me/study-cases", handler.ListMyStudyCases)

	req := httptest.NewRequest("GET", "/api/users/me/study-cases", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
