package service

import (
	"context"
	"sync"
	"testing"

	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/repository"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// ──────────────────────────────────────────────
// NOTIFICATION SPY
// ──────────────────────────────────────────────

type notifCall struct {
	UserID  string
	NType   string
	Title   string
	Message string
	Link    string
}

type notificationSpy struct {
	mu    sync.Mutex
	calls []notifCall
}

func (s *notificationSpy) CreateNotification(_ context.Context, userID, nType, title, message, link string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls = append(s.calls, notifCall{UserID: userID, NType: nType, Title: title, Message: message, Link: link})
	return nil
}

func (s *notificationSpy) ListNotifications(_ context.Context, _ string, _, _ int) ([]map[string]interface{}, int64, error) {
	return nil, 0, nil
}
func (s *notificationSpy) CountUnread(_ context.Context, _ string) (int64, error)  { return 0, nil }
func (s *notificationSpy) MarkAsRead(_ context.Context, _, _ string) error         { return nil }
func (s *notificationSpy) MarkAllAsRead(_ context.Context, _ string) error         { return nil }
func (s *notificationSpy) DeleteNotification(_ context.Context, _, _ string) error { return nil }

func (s *notificationSpy) PopCalls() []notifCall {
	s.mu.Lock()
	defer s.mu.Unlock()
	calls := s.calls
	s.calls = nil
	return calls
}

func (s *notificationSpy) CountByType(nType string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	count := 0
	for _, c := range s.calls {
		if c.NType == nType {
			count++
		}
	}
	return count
}

// ──────────────────────────────────────────────
// TEST HELPERS
// ──────────────────────────────────────────────

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(
		&domain.User{},
		&domain.Category{},
		&domain.Course{},
		&domain.Section{},
		&domain.Lesson{},
		&domain.LessonDetail{},
		&domain.LessonProgress{},
		&domain.CourseEnrollment{},
		&domain.Discussion{},
		&domain.Reply{},
		&domain.Review{},
		&domain.Blog{},
		&domain.Showcase{},
		&domain.ShowcaseLike{},
		&domain.CommunityPoint{},
		&domain.UserLevel{},
		&domain.Certificate{},
		&domain.Notification{},
		&domain.LearningStreak{},
	))
	return db
}

func setupNotificationSpy(t *testing.T) *notificationSpy {
	t.Helper()
	spy := &notificationSpy{}
	orig := notifSvc
	notifSvc = spy
	t.Cleanup(func() { notifSvc = orig })
	return spy
}

func seedUser(db *gorm.DB, id, email, name, role string) {
	db.Create(&domain.User{ID: id, Email: email, Name: name, Role: role})
}

func seedCategory(db *gorm.DB, id, name, slug string) {
	db.Create(&domain.Category{ID: id, Name: name, Slug: slug})
}

func makeUserSvc(db *gorm.DB) IUserService {
	return NewUserService(
		repository.NewUserRepository(db),
		repository.NewCommunityPointRepository(db),
		repository.NewUserLevelRepository(db),
		repository.NewRefreshTokenRepository(db),
	)
}

func makeReplySvc(db *gorm.DB, userSvc IUserService) *ReplyService {
	return NewReplyService(
		repository.NewReplyRepository(db),
		repository.NewReplyReactionRepository(db),
		repository.NewReplyLikeRepository(db),
		repository.NewDiscussionRepository(db),
		userSvc,
	)
}

// ──────────────────────────────────────────────
// 1. COURSE ENROLLMENT
// ──────────────────────────────────────────────

func TestNotif_EnrollCourse_SendsAdminAndUserNotifications(t *testing.T) {
	db := setupTestDB(t)
	spy := setupNotificationSpy(t)
	seedUser(db, "admin1", "a@t.com", "Admin", "admin")
	seedCategory(db, "cat1", "Backend", "backend")
	require.NoError(t, db.Create(&domain.Course{ID: "c1", Title: "Go API", AdminID: "admin1", CategoryID: "cat1", Visibility: "public"}).Error)

	svc := NewCourseService(repository.NewCourseRepository(db), repository.NewLessonRepository(db), repository.NewUserRepository(db), repository.NewEnrollmentRepository(db))

	require.NoError(t, svc.EnrollCourse(context.Background(), "user1", "c1"))
	calls := spy.PopCalls()
	require.Len(t, calls, 2)
	assert.Equal(t, "admin1", calls[0].UserID)
	assert.Equal(t, "course_enrollment", calls[0].NType)
	assert.Equal(t, "user1", calls[1].UserID)
	assert.Equal(t, "enrollment_success", calls[1].NType)
}

func TestNotif_EnrollCourse_Duplicate_NoNewNotification(t *testing.T) {
	db := setupTestDB(t)
	spy := setupNotificationSpy(t)
	seedUser(db, "admin1", "a@t.com", "Admin", "admin")
	seedCategory(db, "cat1", "B", "b")
	require.NoError(t, db.Create(&domain.Course{ID: "c1", Title: "Go API", AdminID: "admin1", CategoryID: "cat1", Visibility: "public"}).Error)

	svc := NewCourseService(repository.NewCourseRepository(db), repository.NewLessonRepository(db), repository.NewUserRepository(db), repository.NewEnrollmentRepository(db))

	require.NoError(t, svc.EnrollCourse(context.Background(), "user1", "c1"))
	spy.PopCalls() // clear
	require.NoError(t, svc.EnrollCourse(context.Background(), "user1", "c1"))
	assert.Len(t, spy.PopCalls(), 0, "no new notifications for duplicate enrollment")
}

// ──────────────────────────────────────────────
// 2. REPLY → NESTED REPLY
// ──────────────────────────────────────────────

func TestNotif_NestedReply_NotifiesParentOwner(t *testing.T) {
	db := setupTestDB(t)
	spy := setupNotificationSpy(t)
	seedUser(db, "u1", "u1@t.com", "U1", "user")
	seedUser(db, "u2", "u2@t.com", "U2", "user")
	seedCategory(db, "cat1", "G", "g")
	require.NoError(t, db.Create(&domain.Discussion{ID: "d1", Title: "T", Content: "C", UserID: "u1", CategoryID: "cat1", Status: "open"}).Error)
	require.NoError(t, db.Create(&domain.Reply{ID: "r1", Content: "Parent", UserID: "u2", DiscussionID: "d1"}).Error)

	svc := makeReplySvc(db, makeUserSvc(db))
	parentID := "r1"
	_, err := svc.CreateReply(context.Background(), "u1", "d1", "Nested reply comment", &parentID)
	require.NoError(t, err)

	// u1 replier == u1 discussion owner, so new_reply is skipped (no self-notification)
	// But nested_reply is sent to u2 (parent reply owner)
	calls := spy.PopCalls()
	countNested, countNew := 0, 0
	for _, c := range calls {
		switch c.NType {
		case "nested_reply":
			countNested++
			assert.Equal(t, "u2", c.UserID)
		case "new_reply":
			countNew++
		}
	}
	assert.Equal(t, 1, countNested, "parent reply owner should get nested_reply")
	assert.Equal(t, 0, countNew, "no new_reply when replier == discussion owner")
}

func TestNotif_NestedReply_ParentEqualsOwner_Dedup(t *testing.T) {
	db := setupTestDB(t)
	spy := setupNotificationSpy(t)
	seedUser(db, "u1", "u1@t.com", "U1", "user")
	seedUser(db, "u2", "u2@t.com", "U2", "user")
	seedCategory(db, "cat1", "G", "g")
	require.NoError(t, db.Create(&domain.Discussion{ID: "d1", Title: "T", Content: "C", UserID: "u1", CategoryID: "cat1", Status: "open"}).Error)
	// parent reply owner (u1) == discussion owner (u1) → should dedup nested_reply
	require.NoError(t, db.Create(&domain.Reply{ID: "r1", Content: "Parent", UserID: "u1", DiscussionID: "d1"}).Error)

	svc := makeReplySvc(db, makeUserSvc(db))
	parentID := "r1"
	_, err := svc.CreateReply(context.Background(), "u2", "d1", "Nested reply comment", &parentID)
	require.NoError(t, err)

	calls := spy.PopCalls()
	countNested, countNew := 0, 0
	for _, c := range calls {
		switch c.NType {
		case "nested_reply":
			countNested++
		case "new_reply":
			countNew++
		}
	}
	assert.Equal(t, 0, countNested, "no nested_reply when parent owner == discussion owner")
	assert.Equal(t, 1, countNew, "single new_reply to discussion owner")
}

// ──────────────────────────────────────────────
// 3. REVIEW → COURSE ADMIN
// ──────────────────────────────────────────────

func TestNotif_Review_NotifiesCourseAdmin(t *testing.T) {
	db := setupTestDB(t)
	spy := setupNotificationSpy(t)
	seedUser(db, "admin1", "a@t.com", "Admin", "admin")
	seedUser(db, "user1", "u@t.com", "User", "user")
	seedCategory(db, "cat1", "B", "b")
	require.NoError(t, db.Create(&domain.Course{ID: "c1", Title: "Go API", AdminID: "admin1", CategoryID: "cat1", Visibility: "public"}).Error)

	svc := NewReviewService(repository.NewReviewRepository(db), repository.NewCourseRepository(db))
	_, err := svc.CreateReview(context.Background(), "user1", "c1", "", 5, "Great course!")
	require.NoError(t, err)

	calls := spy.PopCalls()
	require.Len(t, calls, 1)
	assert.Equal(t, "admin1", calls[0].UserID)
	assert.Equal(t, "new_review", calls[0].NType)
}

func TestNotif_Review_SelfReview_NoNotification(t *testing.T) {
	db := setupTestDB(t)
	spy := setupNotificationSpy(t)
	seedUser(db, "admin1", "a@t.com", "Admin", "admin")
	seedCategory(db, "cat1", "B", "b")
	require.NoError(t, db.Create(&domain.Course{ID: "c1", Title: "Go API", AdminID: "admin1", CategoryID: "cat1", Visibility: "public"}).Error)

	svc := NewReviewService(repository.NewReviewRepository(db), repository.NewCourseRepository(db))
	_, err := svc.CreateReview(context.Background(), "admin1", "c1", "", 5, "My own course review")
	require.NoError(t, err)

	assert.Len(t, spy.PopCalls(), 0, "no notification for self-review")
}

// ──────────────────────────────────────────────
// 4. LESSON COMPLETE → CERTIFICATE
// ──────────────────────────────────────────────

func TestNotif_LessonComplete_SendsCertificateNotification(t *testing.T) {
	db := setupTestDB(t)
	spy := setupNotificationSpy(t)
	seedUser(db, "admin1", "a@t.com", "Admin", "admin")
	seedUser(db, "user1", "u@t.com", "User", "user")
	seedCategory(db, "cat1", "B", "b")
	require.NoError(t, db.Create(&domain.Course{ID: "c1", Title: "Go", AdminID: "admin1", CategoryID: "cat1", Visibility: "public"}).Error)
	require.NoError(t, db.Create(&domain.Section{ID: "s1", Title: "M1", CourseID: "c1", OrderIndex: 1}).Error)
	require.NoError(t, db.Create(&domain.Lesson{ID: "l1", Title: "Basics", Slug: "go-basics", CourseID: "c1", SectionID: "s1", AdminID: "admin1", Difficulty: "beginner", Duration: 45, Visibility: "public"}).Error)

	svc := NewLessonService(
		repository.NewLessonRepository(db),
		repository.NewLessonDetailRepository(db),
		repository.NewLessonProgressRepository(db),
		repository.NewCertificateRepository(db),
		makeUserSvc(db),
		repository.NewCourseRepository(db),
		repository.NewEnrollmentRepository(db),
		repository.NewSectionRepository(db),
		repository.NewLearningStreakRepository(db),
	)

	_, err := svc.StartLesson(context.Background(), "user1", "l1")
	require.NoError(t, err)
	spy.PopCalls() // clear start-lesson calls

	_, err = svc.CompleteLesson(context.Background(), "user1", "l1")
	require.NoError(t, err)

	calls := spy.PopCalls()
	found := false
	for _, c := range calls {
		if c.UserID == "user1" && c.NType == "lesson_completed" {
			found = true
			assert.Contains(t, c.Link, "/courses/c1/lessons/go-basics")
		}
	}
	assert.True(t, found, "user should receive lesson_completed notification")
}

// ──────────────────────────────────────────────
// 5. LEVEL UP
// ──────────────────────────────────────────────

func TestNotif_LevelUp_SendsNotification(t *testing.T) {
	db := setupTestDB(t)
	spy := setupNotificationSpy(t)
	require.NoError(t, db.Create(&domain.User{ID: "u1", Email: "u@t.com", Name: "U", Role: "user", Points: 95, TotalPoints: 95, Level: 1}).Error)

	svc := makeUserSvc(db)
	require.NoError(t, svc.AddPoints(context.Background(), "u1", "test", 10, nil))

	calls := spy.PopCalls()
	found := false
	for _, c := range calls {
		if c.UserID == "u1" && c.NType == "level_up" {
			found = true
			assert.Contains(t, c.Title, "Level Naik")
			// Verifikasi badge content dalam message
			// Karena user_levels tidak di-seed, fallback hardcoded: Level 2 = Intermediate 🌿
			assert.Contains(t, c.Message, "Badge", "level_up message should mention Badge")
			assert.Contains(t, c.Message, "Intermediate", "level_up fallback should be Intermediate")
		}
	}
	assert.True(t, found, "user should receive level_up when level changes")
}

func TestNotif_LevelUp_NoChange_NoNotification(t *testing.T) {
	db := setupTestDB(t)
	spy := setupNotificationSpy(t)
	require.NoError(t, db.Create(&domain.User{ID: "u1", Email: "u@t.com", Name: "U", Role: "user", Points: 50, TotalPoints: 50, Level: 1}).Error)

	svc := makeUserSvc(db)
	require.NoError(t, svc.AddPoints(context.Background(), "u1", "test", 10, nil)) // 60 → still level 1

	calls := spy.PopCalls()
	countLevelUp := 0
	for _, c := range calls {
		if c.NType == "level_up" {
			countLevelUp++
		}
	}
	assert.Equal(t, 0, countLevelUp, "no level_up without actual level change")
}

// ──────────────────────────────────────────────
// 6. BLOG PUBLISHED
// ──────────────────────────────────────────────

func TestNotif_BlogPublished_SendsNotification(t *testing.T) {
	db := setupTestDB(t)
	spy := setupNotificationSpy(t)
	seedUser(db, "a1", "a@t.com", "A", "user")
	seedCategory(db, "cat1", "T", "t")

	svc := NewBlogService(repository.NewBlogRepository(db))
	_, err := svc.CreateBlog(context.Background(), "a1", CreateBlogRequest{
		Title: "My Post", Description: "D", Content: "# C",
		Status: "published", CategoryID: "cat1",
	})
	require.NoError(t, err)

	calls := spy.PopCalls()
	countPublished := 0
	for _, c := range calls {
		if c.NType == "blog_published" {
			countPublished++
			assert.Equal(t, "a1", c.UserID)
		}
	}
	assert.Equal(t, 1, countPublished, "should have blog_published notification")
}

func TestNotif_BlogDraft_NoNotification(t *testing.T) {
	db := setupTestDB(t)
	spy := setupNotificationSpy(t)
	seedUser(db, "a1", "a@t.com", "A", "user")
	seedCategory(db, "cat1", "T", "t")

	svc := NewBlogService(repository.NewBlogRepository(db))
	_, err := svc.CreateBlog(context.Background(), "a1", CreateBlogRequest{
		Title: "Draft", Description: "D", Content: "# D",
		Status: "draft", CategoryID: "cat1",
	})
	require.NoError(t, err)

	assert.Equal(t, 0, spy.CountByType("blog_published"), "no notification for draft blog")
}

// ──────────────────────────────────────────────
// 7. DISCUSSION CREATED
// ──────────────────────────────────────────────

func TestNotif_DiscussionCreated_WithLesson_SendsConfirmation(t *testing.T) {
	db := setupTestDB(t)
	spy := setupNotificationSpy(t)
	seedUser(db, "u1", "u@t.com", "U", "user")
	seedCategory(db, "cat1", "G", "g")

	svc := NewDiscussionService(repository.NewDiscussionRepository(db), repository.NewReplyRepository(db), repository.NewUserRepository(db))
	lessonID := "lesson1"
	_, err := svc.CreateDiscussion(context.Background(), "u1", "Q?", "Help", &lessonID, nil, "cat1")
	require.NoError(t, err)

	assert.Equal(t, 1, spy.CountByType("discussion_created"), "confirmation notification when lessonID present")
}

func TestNotif_DiscussionCreated_NoLesson_HasSelfNotification(t *testing.T) {
	db := setupTestDB(t)
	spy := setupNotificationSpy(t)
	seedUser(db, "u1", "u@t.com", "U", "user")
	seedCategory(db, "cat1", "G", "g")

	svc := NewDiscussionService(repository.NewDiscussionRepository(db), repository.NewReplyRepository(db), repository.NewUserRepository(db))
	_, err := svc.CreateDiscussion(context.Background(), "u1", "Q?", "Help", nil, nil, "cat1")
	require.NoError(t, err)

	// Self-notification (activity history) is always sent regardless of lessonID
	assert.Equal(t, 1, spy.CountByType("discussion_created"), "self-notification should be sent as activity history")
}

// ──────────────────────────────────────────────
// 8. EXISTING NOTIFICATIONS
// ──────────────────────────────────────────────

func TestNotif_NewReply_NotifiesDiscussionOwner(t *testing.T) {
	db := setupTestDB(t)
	spy := setupNotificationSpy(t)
	seedUser(db, "o1", "o@t.com", "O", "user")
	seedUser(db, "r1", "r@t.com", "R", "user")
	seedCategory(db, "cat1", "G", "g")
	require.NoError(t, db.Create(&domain.Discussion{ID: "d1", Title: "T", Content: "C", UserID: "o1", CategoryID: "cat1", Status: "open"}).Error)

	svc := makeReplySvc(db, makeUserSvc(db))
	_, err := svc.CreateReply(context.Background(), "r1", "d1", "Reply text here", nil)
	require.NoError(t, err)

	assert.Equal(t, 1, spy.CountByType("new_reply"))
}

func TestNotif_ReplyLike_NotifiesReplyCreator(t *testing.T) {
	db := setupTestDB(t)
	spy := setupNotificationSpy(t)
	seedUser(db, "u1", "u@t.com", "U1", "user")
	seedUser(db, "u2", "u2@t.com", "U2", "user")
	seedCategory(db, "cat1", "G", "g")
	require.NoError(t, db.Create(&domain.Discussion{ID: "d1", Title: "T", Content: "C", UserID: "u1", CategoryID: "cat1", Status: "open"}).Error)
	require.NoError(t, db.Create(&domain.Reply{ID: "r1", Content: "Helpful", UserID: "u1", DiscussionID: "d1"}).Error)

	svc := makeReplySvc(db, makeUserSvc(db))
	require.NoError(t, svc.LikeReply(context.Background(), "u2", "r1"))

	assert.Equal(t, 1, spy.CountByType("reply_like"))
}

func TestNotif_ReplyLike_SelfLike_NoNotification(t *testing.T) {
	db := setupTestDB(t)
	spy := setupNotificationSpy(t)
	seedUser(db, "u1", "u@t.com", "U1", "user")
	seedCategory(db, "cat1", "G", "g")
	require.NoError(t, db.Create(&domain.Discussion{ID: "d1", Title: "T", Content: "C", UserID: "u1", CategoryID: "cat1", Status: "open"}).Error)
	require.NoError(t, db.Create(&domain.Reply{ID: "r1", Content: "My reply", UserID: "u1", DiscussionID: "d1"}).Error)

	svc := makeReplySvc(db, makeUserSvc(db))
	require.NoError(t, svc.LikeReply(context.Background(), "u1", "r1")) // self-like

	assert.Equal(t, 0, spy.CountByType("reply_like"), "no notification for self-like")
}

func TestNotif_BestAnswer_NotifiesReplyCreator(t *testing.T) {
	db := setupTestDB(t)
	spy := setupNotificationSpy(t)
	seedUser(db, "o1", "o@t.com", "O", "user")
	seedUser(db, "r1", "r@t.com", "R", "user")
	seedCategory(db, "cat1", "G", "g")
	require.NoError(t, db.Create(&domain.Discussion{ID: "d1", Title: "T", Content: "C", UserID: "o1", CategoryID: "cat1", Status: "open"}).Error)
	require.NoError(t, db.Create(&domain.Reply{ID: "r1", Content: "Best answer", UserID: "r1", DiscussionID: "d1"}).Error)

	svc := makeReplySvc(db, makeUserSvc(db))
	require.NoError(t, svc.MarkBestReply(context.Background(), "r1", "d1", "o1"))

	assert.Equal(t, 1, spy.CountByType("best_answer"))
}

func TestNotif_ShowcaseLike_NotifiesOwner(t *testing.T) {
	db := setupTestDB(t)
	spy := setupNotificationSpy(t)
	seedUser(db, "o1", "o@t.com", "O", "user")
	seedUser(db, "l1", "l@t.com", "L", "user")
	seedCategory(db, "cat1", "S", "s")
	require.NoError(t, db.Create(&domain.Showcase{ID: "s1", Title: "Project", UserID: "o1", CategoryID: "cat1", Visibility: "public", Status: "published"}).Error)

	svc := NewShowcaseService(repository.NewShowcaseRepository(db), repository.NewShowcaseLikeRepository(db), makeUserSvc(db))
	require.NoError(t, svc.LikeShowcase(context.Background(), "l1", "s1"))

	assert.Equal(t, 1, spy.CountByType("showcase_like"))
}

// ──────────────────────────────────────────────
// 9. NOTIFICATION SERVICE CRUD
// ──────────────────────────────────────────────

func TestNotif_CreateAndListCRUD(t *testing.T) {
	db := setupTestDB(t)

	svc := NewNotificationService(repository.NewNotificationRepository(db))

	// Create
	err := svc.CreateNotification(context.Background(), "u1", "test_type", "Test Title", "Test Message", "/link")
	require.NoError(t, err)

	// List
	notifs, total, err := svc.ListNotifications(context.Background(), "u1", 1, 10)
	require.NoError(t, err)
	require.Equal(t, int64(1), total)
	assert.Equal(t, "test_type", notifs[0]["type"])
	assert.Equal(t, false, notifs[0]["is_read"])

	// Count unread
	cnt, _ := svc.CountUnread(context.Background(), "u1")
	assert.Equal(t, int64(1), cnt)

	// Mark as read
	notifID := notifs[0]["id"].(string)
	require.NoError(t, svc.MarkAsRead(context.Background(), notifID, "u1"))
	cnt, _ = svc.CountUnread(context.Background(), "u1")
	assert.Equal(t, int64(0), cnt)

	// Mark all as read
	_ = svc.CreateNotification(context.Background(), "u1", "t2", "T2", "M2", "/l2")
	require.NoError(t, svc.MarkAllAsRead(context.Background(), "u1"))
	cnt, _ = svc.CountUnread(context.Background(), "u1")
	assert.Equal(t, int64(0), cnt)

	// Delete
	notifs, _, _ = svc.ListNotifications(context.Background(), "u1", 1, 10)
	if len(notifs) > 0 {
		require.NoError(t, svc.DeleteNotification(context.Background(), notifs[0]["id"].(string), "u1"))
	}
}
