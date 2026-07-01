package handler

import (
	"context"
	"encoding/json"
	"time"

	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/dto"
	"jvalleyverse/internal/service"
	"jvalleyverse/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
)

// ──────────────────────────────────────────────
// ALL SERVICE MOCK INTERFACES
// ──────────────────────────────────────────────

// mockUserService implements service.IUserService
type mockUserService struct {
	users        map[string]*domain.User
	refreshTokens map[string]*domain.RefreshToken
	activityLog  []dto.ActivityItem
	mentors      []dto.MentorItem
	allUsers     []dto.UserListItem
}

func newMockUserService() *mockUserService {
	return &mockUserService{
		users:         make(map[string]*domain.User),
		refreshTokens: make(map[string]*domain.RefreshToken),
	}
}

func (m *mockUserService) GetProfile(_ context.Context, userID string) (*domain.User, error) {
	user, ok := m.users[userID]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (m *mockUserService) UpdateProfile(_ context.Context, userID string, name, bio, avatar string) error {
	user, ok := m.users[userID]
	if !ok {
		return domain.ErrUserNotFound
	}
	if name != "" {
		user.Name = name
	}
	if bio != "" {
		user.Bio = bio
	}
	if avatar != "" {
		user.Avatar = avatar
	}
	return nil
}

func (m *mockUserService) AddPoints(_ context.Context, userID string, category string, points int, metadata map[string]interface{}) error {
	user, ok := m.users[userID]
	if !ok {
		return nil
	}
	user.TotalPoints += points
	user.Points += points
	return nil
}

func (m *mockUserService) GetUserActivityLog(_ context.Context, _ string, page, limit int) ([]dto.ActivityItem, int64, error) {
	return m.activityLog, int64(len(m.activityLog)), nil
}

func (m *mockUserService) CreateUser(_ context.Context, user *domain.User) error {
	if _, ok := m.users[user.Email]; ok {
		return domain.ErrEmailExists
	}
	user.ID = "mock-id-" + user.Email
	m.users[user.Email] = user
	m.users[user.ID] = user
	return nil
}

func (m *mockUserService) GetUserByEmail(_ context.Context, email string) (*domain.User, error) {
	user, ok := m.users[email]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (m *mockUserService) ListAllUsers(_ context.Context, page, limit int) ([]dto.UserListItem, int64, error) {
	return m.allUsers, int64(len(m.allUsers)), nil
}

func (m *mockUserService) GenerateRefreshToken(_ context.Context, userID string) (*domain.RefreshToken, error) {
	token := &domain.RefreshToken{Token: "mock-refresh-token-" + userID, UserID: userID}
	m.refreshTokens[token.Token] = token
	return token, nil
}

func (m *mockUserService) ValidateRefreshToken(_ context.Context, token string) (*domain.RefreshToken, error) {
	rt, ok := m.refreshTokens[token]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return rt, nil
}

func (m *mockUserService) RevokeRefreshToken(_ context.Context, token string) error {
	delete(m.refreshTokens, token)
	return nil
}

func (m *mockUserService) RevokeAllUserTokens(_ context.Context, userID string) error {
	for k, v := range m.refreshTokens {
		if v.UserID == userID {
			delete(m.refreshTokens, k)
		}
	}
	return nil
}

func (m *mockUserService) ChangePassword(_ context.Context, userID, currentPassword, newPassword string) error {
	user, ok := m.users[userID]
	if !ok {
		return domain.ErrUserNotFound
	}
	// Google users with no password can set a new one without current password
	if user.Password != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
			return domain.ErrUnauthorized
		}
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	user.Password = string(hashed)
	return nil
}

func (m *mockUserService) ListMentors(_ context.Context, page, limit int) ([]dto.MentorItem, int64, error) {
	return m.mentors, int64(len(m.mentors)), nil
}

// addTestUser is a helper to quickly create a test user with hashed password
func (m *mockUserService) addTestUser(id, email, password, name, role string) *domain.User {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &domain.User{
		ID:       id,
		Email:    email,
		Password: string(hashed),
		Name:     name,
		Role:     role,
		Avatar:   "https://example.com/avatar.jpg",
		Bio:      "Test bio",
		Level:    1,
		Points:   0,
	}
	m.users[email] = user
	m.users[id] = user
	return user
}

// ──────────────────────────────────────────────
// mockDashboardService implements service.IDashboardService
// ──────────────────────────────────────────────

type mockDashboardService struct{}

func newMockDashboardService() *mockDashboardService {
	return &mockDashboardService{}
}

func (m *mockDashboardService) GetDashboard(_ context.Context, _ string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"courses_in_progress":  2,
		"courses_completed":    3,
		"courses_dropped":      1,
		"unread_notifications": 5,
		"streak_count":         3,
	}, nil
}

func (m *mockDashboardService) GetStreak(_ context.Context, userID string) (int, error) {
	return 3, nil
}

// ──────────────────────────────────────────────
// mockShowcaseService implements service.IShowcaseService
// ──────────────────────────────────────────────

type mockShowcaseService struct {
	showcases map[string]*domain.Showcase
	likes     map[string]bool // "userID:showcaseID" -> true
}

func newMockShowcaseService() *mockShowcaseService {
	return &mockShowcaseService{
		showcases: make(map[string]*domain.Showcase),
		likes:     make(map[string]bool),
	}
}

func (m *mockShowcaseService) CreateShowcase(_ context.Context, userID, title, description string, mediaURLs []string, categoryID, visibility string) (*domain.Showcase, error) {
	if visibility == "" {
		visibility = "public"
	}
	mediaJSON, _ := json.Marshal(mediaURLs)
	sc := &domain.Showcase{
		ID:          "mock-showcase-" + title,
		Title:       title,
		Description: description,
		MediaURLs:   datatypes.JSON(mediaJSON),
		UserID:      userID,
		CategoryID:  categoryID,
		Visibility:  visibility,
		Status:      "published",
		LikesCount:  0,
		ViewsCount:  0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	m.showcases[sc.ID] = sc
	return sc, nil
}

func (m *mockShowcaseService) ListShowcases(_ context.Context, page, limit int, categoryID, sort string) ([]dto.ShowcaseListItem, int64, error) {
	return []dto.ShowcaseListItem{}, 0, nil
}

func (m *mockShowcaseService) ListMyShowcases(_ context.Context, userID string, page, limit int) ([]dto.ShowcaseListItem, int64, error) {
	return []dto.ShowcaseListItem{}, 0, nil
}

func (m *mockShowcaseService) GetShowcaseByID(_ context.Context, showcaseID string) (*dto.ShowcaseDetail, error) {
	sc, ok := m.showcases[showcaseID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	mediaURLs := make([]string, 0)
	if sc.MediaURLs != nil {
		json.Unmarshal(sc.MediaURLs, &mediaURLs)
	}
	return &dto.ShowcaseDetail{
		ID:          sc.ID,
		Title:       sc.Title,
		Description: sc.Description,
		MediaURLs:   mediaURLs,
		LikesCount:  sc.LikesCount,
		ViewsCount:  sc.ViewsCount,
		CreatedAt:   sc.CreatedAt,
	}, nil
}

func (m *mockShowcaseService) UpdateShowcase(_ context.Context, showcaseID, userID, title, description, visibility string) (*domain.Showcase, error) {
	sc, ok := m.showcases[showcaseID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	if sc.UserID != userID {
		return nil, domain.ErrForbidden
	}
	if title != "" {
		sc.Title = title
	}
	if description != "" {
		sc.Description = description
	}
	if visibility != "" {
		sc.Visibility = visibility
	}
	return sc, nil
}

func (m *mockShowcaseService) DeleteShowcase(_ context.Context, showcaseID, userID string) error {
	sc, ok := m.showcases[showcaseID]
	if !ok {
		return domain.ErrNotFound
	}
	if sc.UserID != userID {
		return domain.ErrForbidden
	}
	delete(m.showcases, showcaseID)
	return nil
}

func (m *mockShowcaseService) LikeShowcase(_ context.Context, userID, showcaseID string) error {
	key := userID + ":" + showcaseID
	if m.likes[key] {
		return nil
	}
	m.likes[key] = true
	if sc, ok := m.showcases[showcaseID]; ok {
		sc.LikesCount++
	}
	return nil
}

func (m *mockShowcaseService) UnlikeShowcase(_ context.Context, userID, showcaseID string) error {
	key := userID + ":" + showcaseID
	if !m.likes[key] {
		return nil
	}
	delete(m.likes, key)
	if sc, ok := m.showcases[showcaseID]; ok {
		sc.LikesCount--
	}
	return nil
}

// ──────────────────────────────────────────────
// mockCategoryService implements service.ICategoryService
// ──────────────────────────────────────────────

type mockCategoryService struct {
	categories map[string]*domain.Category
}

func newMockCategoryService() *mockCategoryService {
	return &mockCategoryService{
		categories: make(map[string]*domain.Category),
	}
}

func (m *mockCategoryService) addTestCategory(id, name, slug string) {
	m.categories[id] = &domain.Category{
		ID:   id,
		Name: name,
		Slug: slug,
	}
	m.categories[slug] = m.categories[id]
}

func (m *mockCategoryService) ListCategories(_ context.Context) ([]dto.CategoryBrief, error) {
	var result []dto.CategoryBrief
	seen := make(map[string]bool)
	for _, c := range m.categories {
		if !seen[c.ID] {
			result = append(result, dto.CategoryBrief{ID: c.ID, Name: c.Name, Slug: c.Slug})
			seen[c.ID] = true
		}
	}
	return result, nil
}

func (m *mockCategoryService) GetCategoryBySlug(_ context.Context, slug string) (*dto.CategoryWithCourses, error) {
	cat, ok := m.categories[slug]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return &dto.CategoryWithCourses{
		ID:          cat.ID,
		Name:        cat.Name,
		Slug:        cat.Slug,
		Description: cat.Description,
		Courses:     []dto.CourseListItem{},
	}, nil
}

func (m *mockCategoryService) GetCategoryByID(_ context.Context, id string) (*domain.Category, error) {
	cat, ok := m.categories[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return cat, nil
}

func (m *mockCategoryService) ListCoursesByCategoryID(_ context.Context, _ string) ([]domain.Course, error) {
	return []domain.Course{}, nil
}

func (m *mockCategoryService) ListCoursesByCategoryIDWithEnrollment(_ context.Context, _, _ string) ([]dto.CourseListItem, error) {
	return []dto.CourseListItem{}, nil
}

func (m *mockCategoryService) CreateCategory(_ context.Context, name, slug, description string) (*domain.Category, error) {
	cat := &domain.Category{
		ID:          "mock-cat-" + slug,
		Name:        name,
		Slug:        slug,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	m.categories[cat.ID] = cat
	m.categories[slug] = cat
	return cat, nil
}

func (m *mockCategoryService) UpdateCategory(_ context.Context, id, name, slug, description string) (*domain.Category, error) {
	cat, ok := m.categories[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	if name != "" {
		cat.Name = name
	}
	if slug != "" {
		cat.Slug = slug
	}
	if description != "" {
		cat.Description = description
	}
	return cat, nil
}

func (m *mockCategoryService) DeleteCategory(_ context.Context, id string) error {
	if _, ok := m.categories[id]; !ok {
		return domain.ErrNotFound
	}
	cat := m.categories[id]
	delete(m.categories, cat.Slug)
	delete(m.categories, id)
	return nil
}

// ──────────────────────────────────────────────
// mockCourseService implements service.ICourseService
// ──────────────────────────────────────────────

type mockCourseService struct {
	courses      map[string]*domain.Course
	enrollments  map[string]bool // "userID:courseID" -> bool
}

func newMockCourseService() *mockCourseService {
	return &mockCourseService{
		courses:     make(map[string]*domain.Course),
		enrollments: make(map[string]bool),
	}
}

func (m *mockCourseService) CreateCourse(_ context.Context, adminID, title, desc, thumbnail, categoryID, mentorID string, price float64, hours int, objectives datatypes.JSON) (*domain.Course, error) {
	if title == "" || categoryID == "" {
		return nil, domain.ErrInvalidInput
	}
	course := &domain.Course{
		ID:          "mock-course-" + title,
		Title:       title,
		Description: desc,
		Thumbnail:   thumbnail,
		AdminID:     adminID,
		MentorID:    mentorID,
		CategoryID:  categoryID,
		Price:       price,
		Hours:       hours,
		Visibility:  "public",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	m.courses[course.ID] = course
	return course, nil
}

func (m *mockCourseService) UpdateCourse(_ context.Context, courseID, adminID, title, desc string, price float64, visibility string, _ datatypes.JSON) error {
	course, ok := m.courses[courseID]
	if !ok {
		return domain.ErrCourseNotFound
	}
	if course.AdminID != adminID {
		return domain.ErrForbidden
	}
	if title != "" {
		course.Title = title
	}
	if desc != "" {
		course.Description = desc
	}
	if price >= 0 {
		course.Price = price
	}
	if visibility != "" {
		course.Visibility = visibility
	}
	return nil
}

func (m *mockCourseService) DeleteCourse(_ context.Context, courseID, adminID string) error {
	course, ok := m.courses[courseID]
	if !ok {
		return domain.ErrCourseNotFound
	}
	if course.AdminID != adminID {
		return domain.ErrForbidden
	}
	delete(m.courses, courseID)
	return nil
}

func (m *mockCourseService) ListPublicCourses(_ context.Context, page, limit int, _ *service.CourseListFilter) ([]dto.CourseListItem, int64, error) {
	return []dto.CourseListItem{}, 0, nil
}

func (m *mockCourseService) ListPublicCoursesWithEnrollment(_ context.Context, _ string, page, limit int, _ *service.CourseListFilter) ([]dto.CourseListItem, int64, error) {
	return []dto.CourseListItem{}, 0, nil
}

func (m *mockCourseService) EnrollCourse(_ context.Context, userID, courseID string) error {
	if _, ok := m.courses[courseID]; !ok {
		return domain.ErrCourseNotFound
	}
	m.enrollments[userID+":"+courseID] = true
	return nil
}

func (m *mockCourseService) IsEnrolled(_ context.Context, userID, courseID string) (bool, error) {
	return m.enrollments[userID+":"+courseID], nil
}

func (m *mockCourseService) ListEnrolledCourses(_ context.Context, userID string, page, limit int) ([]dto.EnrolledCourseItem, int64, error) {
	return []dto.EnrolledCourseItem{}, 0, nil
}

func (m *mockCourseService) SetLastLesson(_ context.Context, userID, courseID, lessonID string) error {
	return nil
}

// ──────────────────────────────────────────────
// mockSectionService implements service.ISectionService
// ──────────────────────────────────────────────

type mockSectionService struct {
	sections map[string]*domain.Section
}

func newMockSectionService() *mockSectionService {
	return &mockSectionService{
		sections: make(map[string]*domain.Section),
	}
}

func (m *mockSectionService) CreateSection(_ context.Context, adminID, courseID, title, desc string, orderIndex int) (*domain.Section, error) {
	sec := &domain.Section{
		ID:          "mock-sec-" + title,
		Title:       title,
		Description: desc,
		CourseID:    courseID,
		OrderIndex:  orderIndex,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	m.sections[sec.ID] = sec
	return sec, nil
}

func (m *mockSectionService) UpdateSection(_ context.Context, adminID, sectionID, title, desc string, orderIndex int) (*domain.Section, error) {
	sec, ok := m.sections[sectionID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	if title != "" {
		sec.Title = title
	}
	if desc != "" {
		sec.Description = desc
	}
	if orderIndex >= 0 {
		sec.OrderIndex = orderIndex
	}
	return sec, nil
}

func (m *mockSectionService) DeleteSection(_ context.Context, adminID, sectionID string) error {
	if _, ok := m.sections[sectionID]; !ok {
		return domain.ErrNotFound
	}
	delete(m.sections, sectionID)
	return nil
}

func (m *mockSectionService) GetSection(_ context.Context, sectionID string) (*dto.SectionDetail, error) {
	sec, ok := m.sections[sectionID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return &dto.SectionDetail{
		ID:         sec.ID,
		Title:      sec.Title,
		CourseID:   sec.CourseID,
		OrderIndex: sec.OrderIndex,
		Lessons:    []dto.LessonBrief{},
	}, nil
}

func (m *mockSectionService) ListSectionsByCourse(_ context.Context, courseID string) ([]dto.SectionDetail, error) {
	return []dto.SectionDetail{}, nil
}

func (m *mockSectionService) GetCourseWithSections(_ context.Context, courseID, userID string) (*dto.CourseDetailWithSections, error) {
	course, ok := m.getCourseByID(courseID)
	if !ok {
		return nil, domain.ErrCourseNotFound
	}
	return &dto.CourseDetailWithSections{
		ID:          course.ID,
		Title:       course.Title,
		Description: course.Description,
		Category:    dto.CategoryBrief{ID: course.CategoryID},
		AdminName:   "Admin",
		Sections:    []dto.SectionBrief{},
		CreatedAt:   course.CreatedAt,
	}, nil
}

func (m *mockSectionService) getCourseByID(id string) (*domain.Course, bool) {
	c, ok := m.sections[id]
	if ok {
		return &domain.Course{
			ID:          c.ID,
			Title:       c.Title,
			Description: c.Description,
			CategoryID:  c.CourseID,
			CreatedAt:   c.CreatedAt,
		}, true
	}
	return nil, false
}

// Helper for tests to register a course that sections can belong to
func (m *mockCourseService) addTestCourse(id, title, adminID, categoryID string) {
	m.courses[id] = &domain.Course{
		ID:         id,
		Title:      title,
		AdminID:    adminID,
		CategoryID: categoryID,
		Visibility: "public",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// ──────────────────────────────────────────────
// mockLessonService implements service.ILessonService
// ──────────────────────────────────────────────

type mockLessonService struct {
	lessons   map[string]*domain.Lesson
	progress  map[string]*domain.LessonProgress // "userID:lessonID"
}

func newMockLessonService() *mockLessonService {
	return &mockLessonService{
		lessons:  make(map[string]*domain.Lesson),
		progress: make(map[string]*domain.LessonProgress),
	}
}

func (m *mockLessonService) addTestLesson(id, title, slug, courseID, sectionID, adminID string) {
	m.lessons[id] = &domain.Lesson{
		ID:         id,
		Title:      title,
		Slug:       slug,
		CourseID:   courseID,
		SectionID:  sectionID,
		AdminID:    adminID,
		Difficulty: "beginner",
		Duration:   45,
		Visibility: "public",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func (m *mockLessonService) GetPublicLessonByID(_ context.Context, lessonID string) (*dto.LessonDetailResponse, error) {
	lesson, ok := m.lessons[lessonID]
	if !ok {
		return nil, domain.ErrLessonNotFound
	}
	return &dto.LessonDetailResponse{
		Lesson: dto.ToLessonBrief(*lesson),
		Details: &domain.LessonDetail{
			LessonID: lessonID,
			About:    "Test about",
		},
		Section: &dto.SectionBrief{
			ID: lesson.SectionID,
		},
	}, nil
}

func (m *mockLessonService) GetLessonBySlug(_ context.Context, courseID, slug, userID string) (*dto.LessonDetailResponse, error) {
	for _, l := range m.lessons {
		if l.Slug == slug && l.CourseID == courseID {
			resp := &dto.LessonDetailResponse{
				Lesson: dto.ToLessonBrief(*l),
			}
			if userID != "" {
				if p, ok := m.progress[userID+":"+l.ID]; ok {
					resp.Progress = p
				}
			}
			return resp, nil
		}
	}
	return nil, domain.ErrLessonNotFound
}

func (m *mockLessonService) StartLesson(_ context.Context, userID, lessonID string) (*domain.LessonProgress, error) {
	if _, ok := m.lessons[lessonID]; !ok {
		return nil, domain.ErrLessonNotFound
	}
	now := time.Now()
	p := &domain.LessonProgress{
		ID:                 "mock-progress-" + userID + "-" + lessonID,
		UserID:             userID,
		LessonID:           lessonID,
		Status:             "started",
		StartedAt:          &now,
		ProgressPercentage: 0,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	m.progress[userID+":"+lessonID] = p
	return p, nil
}

func (m *mockLessonService) UpdateProgress(_ context.Context, userID, lessonID string, percentage int, notes string) (*domain.LessonProgress, error) {
	if percentage < 0 || percentage > 100 {
		return nil, domain.ErrInvalidInput
	}
	p, ok := m.progress[userID+":"+lessonID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	p.ProgressPercentage = percentage
	p.Notes = notes
	switch {
	case percentage == 0:
		p.Status = "started"
	case percentage > 0 && percentage < 100:
		p.Status = "in_progress"
	case percentage == 100:
		p.Status = "completed"
		now := time.Now()
		p.CompletedAt = &now
	}
	return p, nil
}

func (m *mockLessonService) CompleteLesson(_ context.Context, userID, lessonID string) (map[string]interface{}, error) {
	lesson, ok := m.lessons[lessonID]
	if !ok {
		return nil, domain.ErrLessonNotFound
	}
	p, ok := m.progress[userID+":"+lessonID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	if p.Status == "completed" {
		return nil, domain.ErrInvalidInput // "lesson already completed"
	}
	now := time.Now()
	p.Status = "completed"
	p.ProgressPercentage = 100
	p.CompletedAt = &now

	return map[string]interface{}{
		"message": "Lesson completed!",
		"certificate": &domain.Certificate{
			ID:         "mock-cert-" + lessonID,
			UserID:     userID,
			LessonID:   lessonID,
			UniqueCode: "CERT-mock1234",
			IssuedAt:   now,
		},
		"achievement": map[string]interface{}{
			"type":        "certificate",
			"title":       lesson.Title,
			"unique_code": "CERT-mock1234",
		},
		"progress":       p,
		"points_awarded": 50,
		"next_lesson":    nil,
	}, nil
}

func (m *mockLessonService) AdminCreateLesson(_ context.Context, adminID string, input domain.Lesson) (*domain.Lesson, error) {
	input.ID = "mock-lesson-" + input.Slug
	input.AdminID = adminID
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()
	m.lessons[input.ID] = &input
	return &input, nil
}

func (m *mockLessonService) AdminUpdateLesson(_ context.Context, adminID, id string, input domain.Lesson) (*domain.Lesson, error) {
	lesson, ok := m.lessons[id]
	if !ok {
		return nil, domain.ErrLessonNotFound
	}
	if lesson.AdminID != adminID {
		return nil, domain.ErrForbidden
	}
	if input.Title != "" {
		lesson.Title = input.Title
	}
	if input.Slug != "" {
		lesson.Slug = input.Slug
	}
	if input.Description != "" {
		lesson.Description = input.Description
	}
	return lesson, nil
}

func (m *mockLessonService) AdminDeleteLesson(_ context.Context, adminID, id string) error {
	lesson, ok := m.lessons[id]
	if !ok {
		return domain.ErrLessonNotFound
	}
	if lesson.AdminID != adminID {
		return domain.ErrForbidden
	}
	delete(m.lessons, id)
	return nil
}

func (m *mockLessonService) AdminCreateLessonDetail(_ context.Context, lessonID, about, rules string, tools, media, resources interface{}) (*domain.LessonDetail, error) {
	if _, ok := m.lessons[lessonID]; !ok {
		return nil, domain.ErrLessonNotFound
	}
	toolsJSON, _ := json.Marshal(tools)
	mediaJSON, _ := json.Marshal(media)
	resJSON, _ := json.Marshal(resources)
	return &domain.LessonDetail{
		ID:            "mock-detail-" + lessonID,
		LessonID:      lessonID,
		About:         about,
		Rules:         rules,
		Tools:         datatypes.JSON(toolsJSON),
		ResourceMedia: datatypes.JSON(mediaJSON),
		Resources:     datatypes.JSON(resJSON),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

func (m *mockLessonService) ListLessonsByCourse(_ context.Context, courseID string, limit, offset int) ([]domain.Lesson, int64, error) {
	return []domain.Lesson{}, 0, nil
}

func (m *mockLessonService) ListLessonsBySection(_ context.Context, sectionID string) ([]domain.Lesson, int64, error) {
	return []domain.Lesson{}, 0, nil
}

// ──────────────────────────────────────────────
// mockDiscussionService implements service.IDiscussionService
// ──────────────────────────────────────────────

type mockDiscussionService struct {
	discussions map[string]*domain.Discussion
}

func newMockDiscussionService() *mockDiscussionService {
	return &mockDiscussionService{
		discussions: make(map[string]*domain.Discussion),
	}
}

func (m *mockDiscussionService) CreateDiscussion(_ context.Context, userID, title, content string, lessonID, studyCaseID *string, categoryID string) (*domain.Discussion, error) {
	d := &domain.Discussion{
		ID:          "mock-disc-" + title,
		Title:       title,
		Content:     content,
		UserID:      userID,
		LessonID:    lessonID,
		StudyCaseID: studyCaseID,
		CategoryID:  categoryID,
		Status:      "open",
		ViewsCount:  0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	m.discussions[d.ID] = d
	return d, nil
}

func (m *mockDiscussionService) ListDiscussions(_ context.Context, page, limit int, lessonID, studyCaseID, status *string) ([]dto.DiscussionListItem, int64, error) {
	return []dto.DiscussionListItem{}, 0, nil
}

func (m *mockDiscussionService) ListUserDiscussions(_ context.Context, userID string, page, limit int) ([]dto.DiscussionListItem, int64, error) {
	return []dto.DiscussionListItem{}, 0, nil
}

func (m *mockDiscussionService) GetDiscussionWithReplies(_ context.Context, discussionID string) (*dto.DiscussionDetail, error) {
	d, ok := m.discussions[discussionID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return &dto.DiscussionDetail{
		ID:         d.ID,
		Title:      d.Title,
		Content:    d.Content,
		Status:     d.Status,
		ViewsCount: d.ViewsCount,
		CreatedAt:  d.CreatedAt,
		Replies:    []dto.ReplyInDiscussion{},
	}, nil
}

func (m *mockDiscussionService) UpdateDiscussion(_ context.Context, discussionID, userID, title, content string) error {
	d, ok := m.discussions[discussionID]
	if !ok {
		return domain.ErrNotFound
	}
	if d.UserID != userID {
		return domain.ErrForbidden
	}
	d.Title = title
	d.Content = content
	return nil
}

func (m *mockDiscussionService) CloseDiscussion(_ context.Context, discussionID, userID string) error {
	d, ok := m.discussions[discussionID]
	if !ok {
		return domain.ErrNotFound
	}
	if d.UserID != userID {
		return domain.ErrForbidden
	}
	d.Status = "closed"
	return nil
}

func (m *mockDiscussionService) DeleteDiscussion(_ context.Context, discussionID, userID string, isAdmin bool) error {
	d, ok := m.discussions[discussionID]
	if !ok {
		return domain.ErrNotFound
	}
	if d.UserID != userID && !isAdmin {
		return domain.ErrForbidden
	}
	delete(m.discussions, discussionID)
	return nil
}

// ──────────────────────────────────────────────
// mockReplyService implements service.IReplyService
// ──────────────────────────────────────────────

type mockReplyService struct {
	replies map[string]*domain.Reply
}

func newMockReplyService() *mockReplyService {
	return &mockReplyService{
		replies: make(map[string]*domain.Reply),
	}
}

func (m *mockReplyService) CreateReply(_ context.Context, userID, discussionID, content string, parentID *string) (*domain.Reply, error) {
	idContent := content
	if len(idContent) > 10 {
		idContent = idContent[:10]
	}
	r := &domain.Reply{
		ID:           "mock-reply-" + idContent,
		Content:      content,
		UserID:       userID,
		DiscussionID: discussionID,
		ParentID:     parentID,
		LikesCount:   0,
		IsMarkedBest: false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	m.replies[r.ID] = r
	return r, nil
}

func (m *mockReplyService) UpdateReply(_ context.Context, replyID, userID, content string) error {
	r, ok := m.replies[replyID]
	if !ok {
		return domain.ErrNotFound
	}
	if r.UserID != userID {
		return domain.ErrForbidden
	}
	r.Content = content
	return nil
}

func (m *mockReplyService) DeleteReply(_ context.Context, replyID, userID string, isAdmin bool) error {
	r, ok := m.replies[replyID]
	if !ok {
		return domain.ErrNotFound
	}
	if r.UserID != userID && !isAdmin {
		return domain.ErrForbidden
	}
	delete(m.replies, replyID)
	return nil
}

func (m *mockReplyService) LikeReply(_ context.Context, userID, replyID string) error {
	r, ok := m.replies[replyID]
	if !ok {
		return domain.ErrNotFound
	}
	r.LikesCount++
	return nil
}

func (m *mockReplyService) MarkBestReply(_ context.Context, replyID, discussionID, userID string) error {
	r, ok := m.replies[replyID]
	if !ok {
		return domain.ErrNotFound
	}
	r.IsMarkedBest = true
	return nil
}

func (m *mockReplyService) ListRepliesByUser(_ context.Context, userID string, page, limit int) ([]dto.ReplyListItem, int64, error) {
	return []dto.ReplyListItem{}, 0, nil
}

// ──────────────────────────────────────────────
// mockReviewService implements service.IReviewService
// ──────────────────────────────────────────────

type mockReviewService struct {
	reviews map[string]*domain.Review
}

func newMockReviewService() *mockReviewService {
	return &mockReviewService{
		reviews: make(map[string]*domain.Review),
	}
}

func (m *mockReviewService) CreateReview(_ context.Context, userID, courseID, lessonID string, rating int, message string) (*domain.Review, error) {
	if rating < 1 || rating > 5 || message == "" {
		return nil, domain.ErrInvalidInput
	}
	id := userID
	if len(id) > 8 {
		id = id[:8]
	}
	r := &domain.Review{
		ID:       "mock-review-" + id,
		UserID:   userID,
		CourseID: courseID,
		LessonID: lessonID,
		Rating:   rating,
		Message:  message,
		User: domain.User{
			ID:     userID,
			Name:   "Test User",
			Avatar: "https://example.com/avatar.jpg",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.reviews[r.ID] = r
	return r, nil
}

func (m *mockReviewService) UpdateReview(_ context.Context, reviewID, userID string, rating int, message string) (*domain.Review, error) {
	r, ok := m.reviews[reviewID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	if r.UserID != userID {
		return nil, domain.ErrForbidden
	}
	if rating >= 1 && rating <= 5 {
		r.Rating = rating
	}
	if message != "" {
		r.Message = message
	}
	r.UpdatedAt = time.Now()
	return r, nil
}

func (m *mockReviewService) DeleteReview(_ context.Context, reviewID, userID string) error {
	r, ok := m.reviews[reviewID]
	if !ok {
		return domain.ErrNotFound
	}
	if r.UserID != userID {
		return domain.ErrForbidden
	}
	delete(m.reviews, reviewID)
	return nil
}

func (m *mockReviewService) ListCourseReviews(_ context.Context, courseID string, page, limit int) ([]dto.ReviewItem, int64, error) {
	return []dto.ReviewItem{}, 0, nil
}

func (m *mockReviewService) ListLessonReviews(_ context.Context, lessonID string, page, limit int) ([]dto.ReviewItem, int64, error) {
	return []dto.ReviewItem{}, 0, nil
}

// ──────────────────────────────────────────────
// mockCertificateService implements service.ICertificateService
// ──────────────────────────────────────────────

type mockCertificateService struct {
	certificates map[string]*domain.Certificate
}

func newMockCertificateService() *mockCertificateService {
	return &mockCertificateService{
		certificates: make(map[string]*domain.Certificate),
	}
}

func (m *mockCertificateService) IssueCertificate(_ context.Context, userID, lessonID, code string) (*domain.Certificate, error) {
	cert := &domain.Certificate{
		ID:         "mock-cert-" + code,
		UserID:     userID,
		LessonID:   lessonID,
		UniqueCode: code,
		IssuedAt:   time.Now(),
	}
	m.certificates[cert.ID] = cert
	return cert, nil
}

func (m *mockCertificateService) GetCertificate(_ context.Context, certID string) (*dto.CertificateItem, error) {
	cert, ok := m.certificates[certID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return &dto.CertificateItem{
		ID:          cert.ID,
		UniqueCode:  cert.UniqueCode,
		IssuedAt:    cert.IssuedAt,
		LessonName:  "Test Lesson",
		UserName:    "Test User",
		Achievement: &dto.AchievementInfo{Type: "certificate"},
	}, nil
}

func (m *mockCertificateService) GetCertificateByCode(_ context.Context, code, requesterID, requesterRole string) (*dto.CertificateItem, error) {
	for _, cert := range m.certificates {
		if cert.UniqueCode == code {
			if cert.UserID != requesterID && requesterRole != "admin" {
				return nil, domain.ErrForbidden
			}
			return &dto.CertificateItem{
				ID:          cert.ID,
				UniqueCode:  cert.UniqueCode,
				IssuedAt:    cert.IssuedAt,
				UserID:      cert.UserID,
				LessonID:    cert.LessonID,
				LessonName:  "Test Lesson",
				UserName:    "Test User",
				Achievement: &dto.AchievementInfo{Type: "certificate", Title: "Test Lesson", UniqueCode: cert.UniqueCode},
			}, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockCertificateService) VerifyCertificateByCode(_ context.Context, code string) (*dto.CertificateItem, error) {
	for _, cert := range m.certificates {
		if cert.UniqueCode == code {
			return &dto.CertificateItem{
				ID:          cert.ID,
				UniqueCode:  cert.UniqueCode,
				IssuedAt:    cert.IssuedAt,
				UserID:      cert.UserID,
				LessonID:    cert.LessonID,
				LessonName:  "Test Lesson",
				UserName:    "Test User",
				Achievement: &dto.AchievementInfo{Type: "certificate", Title: "Test Lesson", UniqueCode: cert.UniqueCode},
			}, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockCertificateService) ListUserCertificates(_ context.Context, userID string, page, limit int) ([]dto.CertificateItem, int64, error) {
	var items []dto.CertificateItem
	for _, cert := range m.certificates {
		if cert.UserID == userID {
			items = append(items, dto.CertificateItem{
				ID:          cert.ID,
				UniqueCode:  cert.UniqueCode,
				IssuedAt:    cert.IssuedAt,
				LessonName:  "Test Lesson",
				Achievement: &dto.AchievementInfo{Type: "certificate"},
			})
		}
	}
	return items, int64(len(items)), nil
}

// ──────────────────────────────────────────────
// mockGamificationService implements service.IGamificationService
// ──────────────────────────────────────────────

type mockGamificationService struct {
	userSvc *mockUserService
}

func newMockGamificationService(userSvc *mockUserService) *mockGamificationService {
	return &mockGamificationService{userSvc: userSvc}
}

func (m *mockGamificationService) AwardPoints(_ context.Context, userID, activityType string, points int, metadata map[string]interface{}) error {
	return m.userSvc.AddPoints(nil, userID, activityType, points, metadata)
}

func (m *mockGamificationService) GetLeaderboard(_ context.Context, limit int) ([]dto.LeaderboardItem, error) {
	return []dto.LeaderboardItem{
		{Rank: 1, UserID: "admin-id", Name: "Admin", TotalPoints: 5000, Level: 5},
		{Rank: 2, UserID: "user-id", Name: "User", TotalPoints: 2500, Level: 4},
	}, nil
}

func (m *mockGamificationService) GetUserActivityLog(_ context.Context, userID string, page, limit int) ([]dto.ActivityItem, error) {
	return []dto.ActivityItem{}, nil
}

func (m *mockGamificationService) GetLevelInfo() []dto.LevelInfo {
	return []dto.LevelInfo{
		{Name: "Beginner", Threshold: 0, Color: "#6366f1", Description: "Just starting your journey"},
		{Name: "Intermediate", Threshold: 100, Color: "#8b5cf6", Description: "Building momentum"},
		{Name: "Advanced", Threshold: 500, Color: "#d946ef", Description: "Getting serious"},
		{Name: "Expert", Threshold: 1000, Color: "#ec4899", Description: "Mastering skills"},
		{Name: "Master", Threshold: 2000, Color: "#f43f5e", Description: "Peak achievement"},
	}
}

func (m *mockGamificationService) GetUserStats(_ context.Context, userID string) (*dto.UserStats, error) {
	user, ok := m.userSvc.users[userID]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return &dto.UserStats{
		UserID:       user.ID,
		Name:         user.Name,
		TotalPoints:  user.TotalPoints,
		CurrentLevel: user.Level,
	}, nil
}

// ──────────────────────────────────────────────
// mockStudyCaseService implements service.IStudyCaseService
// ──────────────────────────────────────────────

type mockStudyCaseService struct {
	studyCases map[string]*domain.StudyCase
}

func newMockStudyCaseService() *mockStudyCaseService {
	return &mockStudyCaseService{
		studyCases: make(map[string]*domain.StudyCase),
	}
}

func (m *mockStudyCaseService) CreateStudyCase(_ context.Context, userID, name, desc, imgURL, youtubeURL string, categoryID *string, tags []string) (*domain.StudyCase, error) {
	if name == "" {
		return nil, domain.ErrInvalidInput
	}
	tagsJSON, _ := json.Marshal(tags)
	sc := &domain.StudyCase{
		ID:          "mock-sc-" + name,
		Name:        name,
		Description: desc,
		ImgURL:      imgURL,
		YoutubeURL:  youtubeURL,
		CategoryID:  categoryID,
		Tags:        datatypes.JSON(tagsJSON),
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	m.studyCases[sc.ID] = sc
	return sc, nil
}

func (m *mockStudyCaseService) GetStudyCaseByID(_ context.Context, id string) (*dto.StudyCaseDetail, error) {
	sc, ok := m.studyCases[id]
	if !ok {
		return nil, domain.ErrStudyCaseNotFound
	}
	return &dto.StudyCaseDetail{
		ID:          sc.ID,
		Name:        sc.Name,
		Description: sc.Description,
		ImgURL:      sc.ImgURL,
		YoutubeURL:  sc.YoutubeURL,
		CreatedAt:   sc.CreatedAt,
	}, nil
}

func (m *mockStudyCaseService) ListStudyCases(_ context.Context, page, limit int, filter *service.StudyCaseListFilter) ([]dto.StudyCaseListItem, int64, error) {
	return []dto.StudyCaseListItem{}, 0, nil
}

func (m *mockStudyCaseService) ListStudyCasesByUser(_ context.Context, userID string, page, limit int) ([]dto.StudyCaseListItem, int64, error) {
	return []dto.StudyCaseListItem{}, 0, nil
}

func (m *mockStudyCaseService) UpdateStudyCase(_ context.Context, id, name, desc, imgURL, youtubeURL string, categoryID *string, tags []string) (*domain.StudyCase, error) {
	sc, ok := m.studyCases[id]
	if !ok {
		return nil, domain.ErrStudyCaseNotFound
	}
	if name != "" {
		sc.Name = name
	}
	if desc != "" {
		sc.Description = desc
	}
	if imgURL != "" {
		sc.ImgURL = imgURL
	}
	if youtubeURL != "" {
		sc.YoutubeURL = youtubeURL
	}
	if len(tags) > 0 {
		tagsJSON, _ := json.Marshal(tags)
		sc.Tags = datatypes.JSON(tagsJSON)
	}
	return sc, nil
}

func (m *mockStudyCaseService) DeleteStudyCase(_ context.Context, id string) error {
	if _, ok := m.studyCases[id]; !ok {
		return domain.ErrStudyCaseNotFound
	}
	delete(m.studyCases, id)
	return nil
}

// ──────────────────────────────────────────────
// mockBlogService implements service.IBlogService
// ──────────────────────────────────────────────

type mockBlogService struct {
	blogs map[string]*domain.Blog
}

func newMockBlogService() *mockBlogService {
	return &mockBlogService{
		blogs: make(map[string]*domain.Blog),
	}
}

func (m *mockBlogService) CreateBlog(_ context.Context, userID string, req service.CreateBlogRequest) (*dto.BlogDetail, error) {
	if req.Title == "" {
		return nil, domain.ErrInvalidInput
	}
	status := req.Status
	if status == "" {
		status = "draft"
	}
	tagsJSON, _ := json.Marshal(req.Tags)
	blog := &domain.Blog{
		ID:          "mock-blog-" + req.Title,
		Title:       req.Title,
		Slug:        req.Title,
		Description: req.Description,
		Content:     req.Content,
		CoverImgURL: req.CoverImgURL,
		Status:      status,
		UserID:      userID,
		CategoryID:  req.CategoryID,
		Tags:        tagsJSON,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	m.blogs[blog.ID] = blog
	return &dto.BlogDetail{
		ID:          blog.ID,
		Title:       blog.Title,
		Slug:        blog.Slug,
		Description: blog.Description,
		Content:     blog.Content,
		CoverImgURL: blog.CoverImgURL,
		Status:      blog.Status,
		Author:      dto.UserBrief{ID: userID, Name: "Admin", Avatar: ""},
		Category:    dto.CategoryBrief{ID: req.CategoryID},
		CreatedAt:   blog.CreatedAt,
		UpdatedAt:   blog.UpdatedAt,
	}, nil
}

func (m *mockBlogService) ListBlogs(_ context.Context, page, limit int, search, categoryID, tag string) ([]dto.BlogListItem, *dto.Pagination, error) {
	var items []dto.BlogListItem
	for _, b := range m.blogs {
		if b.Status == "published" {
			items = append(items, dto.BlogListItem{
				ID:     b.ID,
				Title:  b.Title,
				Slug:   b.Slug,
				Status: b.Status,
			})
		}
	}
	return items, &dto.Pagination{Page: page, Limit: limit, Total: int64(len(items))}, nil
}

func (m *mockBlogService) ListMyBlogs(_ context.Context, userID string, page, limit int, status string) ([]dto.BlogListItem, *dto.Pagination, error) {
	var items []dto.BlogListItem
	for _, b := range m.blogs {
		if b.UserID == userID {
			if status == "" || b.Status == status {
				items = append(items, dto.BlogListItem{
					ID:     b.ID,
					Title:  b.Title,
					Slug:   b.Slug,
					Status: b.Status,
				})
			}
		}
	}
	return items, &dto.Pagination{Page: page, Limit: limit, Total: int64(len(items))}, nil
}

func (m *mockBlogService) GetBlogByID(_ context.Context, id string) (*dto.BlogDetail, error) {
	b, ok := m.blogs[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return &dto.BlogDetail{
		ID:          b.ID,
		Title:       b.Title,
		Slug:        b.Slug,
		Description: b.Description,
		Content:     b.Content,
		CoverImgURL: b.CoverImgURL,
		Status:      b.Status,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}, nil
}

func (m *mockBlogService) UpdateBlog(_ context.Context, blogID, userID string, req service.UpdateBlogRequest) error {
	b, ok := m.blogs[blogID]
	if !ok {
		return domain.ErrNotFound
	}
	if b.UserID != userID {
		return domain.ErrForbidden
	}
	if req.Title != nil {
		b.Title = *req.Title
	}
	return nil
}

func (m *mockBlogService) AdminUpdateBlog(_ context.Context, blogID string, req service.UpdateBlogRequest) error {
	b, ok := m.blogs[blogID]
	if !ok {
		return domain.ErrNotFound
	}
	if req.Title != nil {
		b.Title = *req.Title
	}
	return nil
}

func (m *mockBlogService) DeleteBlog(_ context.Context, blogID, userID string) error {
	b, ok := m.blogs[blogID]
	if !ok {
		return domain.ErrNotFound
	}
	if b.UserID != userID {
		return domain.ErrForbidden
	}
	delete(m.blogs, blogID)
	return nil
}

func (m *mockBlogService) AdminDeleteBlog(_ context.Context, blogID string) error {
	if _, ok := m.blogs[blogID]; !ok {
		return domain.ErrNotFound
	}
	delete(m.blogs, blogID)
	return nil
}

// ──────────────────────────────────────────────
// mockNotificationService implements service.INotificationService
// ──────────────────────────────────────────────

type mockNotificationService struct {
	notifications []map[string]interface{}
	unreadCount   int64
}

func newMockNotificationService() *mockNotificationService {
	return &mockNotificationService{
		unreadCount: 5,
	}
}

func (m *mockNotificationService) CreateNotification(_ context.Context, userID, nType, title, message, link string) error {
	return nil
}

func (m *mockNotificationService) ListNotifications(_ context.Context, _ string, page, limit int) ([]map[string]interface{}, int64, error) {
	return m.notifications, int64(len(m.notifications)), nil
}

func (m *mockNotificationService) CountUnread(_ context.Context, _ string) (int64, error) {
	return m.unreadCount, nil
}

func (m *mockNotificationService) MarkAsRead(_ context.Context, notificationID, userID string) error {
	return nil
}

func (m *mockNotificationService) MarkAllAsRead(_ context.Context, userID string) error {
	return nil
}

func (m *mockNotificationService) DeleteNotification(_ context.Context, notificationID, userID string) error {
	return nil
}

// ──────────────────────────────────────────────
// mockFaqService implements service.IFaqService
// ──────────────────────────────────────────────

type mockFaqService struct {
	faqs map[string]*dto.FAQItem
}

func newMockFaqService() *mockFaqService {
	return &mockFaqService{
		faqs: make(map[string]*dto.FAQItem),
	}
}

func (m *mockFaqService) addTestFAQ(id, question, answer, category string, orderIndex int, isActive bool) {
	m.faqs[id] = &dto.FAQItem{
		ID:         id,
		Question:   question,
		Answer:     answer,
		Category:   category,
		OrderIndex: orderIndex,
		IsActive:   isActive,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func (m *mockFaqService) CreateFAQ(_ context.Context, question, answer, category string, orderIndex int) (*dto.FAQItem, error) {
	if question == "" || answer == "" {
		return nil, domain.ErrInvalidInput
	}
	if category == "" {
		category = "general"
	}
	item := &dto.FAQItem{
		ID:         "mock-faq-" + question,
		Question:   question,
		Answer:     answer,
		Category:   category,
		OrderIndex: orderIndex,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	m.faqs[item.ID] = item
	return item, nil
}

func (m *mockFaqService) UpdateFAQ(_ context.Context, id, question, answer, category string, orderIndex int, isActive bool) (*dto.FAQItem, error) {
	item, ok := m.faqs[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	if question != "" {
		item.Question = question
	}
	if answer != "" {
		item.Answer = answer
	}
	if category != "" {
		item.Category = category
	}
	item.OrderIndex = orderIndex
	item.IsActive = isActive
	return item, nil
}

func (m *mockFaqService) DeleteFAQ(_ context.Context, id string) error {
	if _, ok := m.faqs[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.faqs, id)
	return nil
}

func (m *mockFaqService) ListAllFAQs(_ context.Context, page, limit int) ([]dto.FAQItem, int64, error) {
	items := make([]dto.FAQItem, 0, len(m.faqs))
	for _, f := range m.faqs {
		items = append(items, *f)
	}
	return items, int64(len(items)), nil
}

func (m *mockFaqService) ListPublicFAQs(_ context.Context, page, limit int) ([]dto.FAQItem, int64, error) {
	items := make([]dto.FAQItem, 0)
	for _, f := range m.faqs {
		if f.IsActive {
			items = append(items, *f)
		}
	}
	return items, int64(len(items)), nil
}

func (m *mockFaqService) GetFAQByID(_ context.Context, id string) (*dto.FAQItem, error) {
	item, ok := m.faqs[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return item, nil
}

// ──────────────────────────────────────────────
// mockCompanyService implements service.ICompanyService
// ──────────────────────────────────────────────

type mockCompanyService struct {
	company *dto.CompanyItem
}

func newMockCompanyService() *mockCompanyService {
	return &mockCompanyService{
		company: &dto.CompanyItem{
			ID:        "mock-company-1",
			BrandName: "JValleyVerse",
			Tagline:   "Learn, Build, Grow Together",
			Email:     "hello@jvalleyverse.com",
			Domain:    "https://jvalleyverse.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

func (m *mockCompanyService) GetCompany(_ context.Context) (*dto.CompanyItem, error) {
	return m.company, nil
}

func (m *mockCompanyService) UpdateCompany(_ context.Context, input domain.Company) (*dto.CompanyItem, error) {
	if input.BrandName != "" {
		m.company.BrandName = input.BrandName
	}
	if input.Tagline != "" {
		m.company.Tagline = input.Tagline
	}
	if input.Email != "" {
		m.company.Email = input.Email
	}
	if input.Domain != "" {
		m.company.Domain = input.Domain
	}
	m.company.UpdatedAt = time.Now()
	return m.company, nil
}

// ──────────────────────────────────────────────
// TOKEN HELPERS for test setup
// ──────────────────────────────────────────────

func generateTestToken(userID, role string) string {
	token, _ := utils.GenerateJWT(userID, role)
	return token
}

func generateTestXSRF() string {
	return utils.GenerateXSRFToken()
}

// setupProtectedApp creates a Fiber app with JWT auth, XSRF protection, and idempotency
func setupProtectedApp(userID, role string) *fiber.App {
	app := setupTestApp()
	token := generateTestToken(userID, role)
	xsrf := generateTestXSRF()

	// Add JWT middleware manually
	app.Use(func(c *fiber.Ctx) error {
		c.Request().Header.Set("Authorization", "Bearer "+token)
		c.Request().Header.Set("X-XSRF-TOKEN", xsrf)
		c.Request().Header.Set("Cookie", "XSRF-TOKEN="+xsrf)
		c.Locals("userID", userID)
		c.Locals("role", role)
		return c.Next()
	})

	return app
}

// setupAdminApp creates a Fiber app with admin JWT auth
func setupAdminApp(adminID string) *fiber.App {
	return setupProtectedApp(adminID, "admin")
}

// ──────────────────────────────────────────────
// Compile-time interface checks
// Ensure our mocks properly implement the service interfaces
// ──────────────────────────────────────────────

var _ interface {
	GetProfile(context.Context, string) (*domain.User, error)
	UpdateProfile(context.Context, string, string, string, string) error
	ChangePassword(context.Context, string, string, string) error
	AddPoints(context.Context, string, string, int, map[string]interface{}) error
	GetUserActivityLog(context.Context, string, int, int) ([]dto.ActivityItem, int64, error)
	CreateUser(context.Context, *domain.User) error
	GetUserByEmail(context.Context, string) (*domain.User, error)
	ListAllUsers(context.Context, int, int) ([]dto.UserListItem, int64, error)
	ListMentors(context.Context, int, int) ([]dto.MentorItem, int64, error)
	GenerateRefreshToken(context.Context, string) (*domain.RefreshToken, error)
	ValidateRefreshToken(context.Context, string) (*domain.RefreshToken, error)
	RevokeRefreshToken(context.Context, string) error
	RevokeAllUserTokens(context.Context, string) error
} = (*mockUserService)(nil)
