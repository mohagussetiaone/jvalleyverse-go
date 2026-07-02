package service

import (
	"context"

	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/dto"
	"jvalleyverse/internal/repository"

	"gorm.io/datatypes"
)

type CourseListFilter struct {
	CategoryID *string
	MinPrice   *float64
	MaxPrice   *float64
}

type ICourseService interface {
	CreateCourse(ctx context.Context, adminID string, title, description, thumbnail string, categoryID string, mentorID string, price float64, hours int, tools datatypes.JSON, learningObjectives datatypes.JSON) (*domain.Course, error)
	UpdateCourse(ctx context.Context, courseID, adminID string, title, description string, price float64, visibility string, tools datatypes.JSON, learningObjectives datatypes.JSON) error
	DeleteCourse(ctx context.Context, courseID, adminID string) error
	ListPublicCourses(ctx context.Context, page, limit int, filter *CourseListFilter) ([]dto.CourseListItem, int64, error)
	ListPublicCoursesWithEnrollment(ctx context.Context, userID string, page, limit int, filter *CourseListFilter) ([]dto.CourseListItem, int64, error)
	EnrollCourse(ctx context.Context, userID, courseID string) error
	IsEnrolled(ctx context.Context, userID, courseID string) (bool, error)
	ListEnrolledCourses(ctx context.Context, userID string, page, limit int) ([]dto.EnrolledCourseItem, int64, error)
	SetLastLesson(ctx context.Context, userID, courseID, lessonID string) error
}

type CourseService struct {
	courseRepo *repository.CourseRepository
	lessonRepo *repository.LessonRepository
	userRepo   *repository.UserRepository
	enrollRepo *repository.EnrollmentRepository
}

func NewCourseService(
	courseRepo *repository.CourseRepository,
	lessonRepo *repository.LessonRepository,
	userRepo *repository.UserRepository,
	enrollRepo *repository.EnrollmentRepository,
) *CourseService {
	return &CourseService{
		courseRepo: courseRepo,
		lessonRepo: lessonRepo,
		userRepo:   userRepo,
		enrollRepo: enrollRepo,
	}
}

func (s *CourseService) CreateCourse(ctx context.Context, adminID string, title, description, thumbnail string, categoryID string, mentorID string, price float64, hours int, tools datatypes.JSON, learningObjectives datatypes.JSON) (*domain.Course, error) {
	if title == "" || categoryID == "" {
		return nil, domain.ErrInvalidInput
	}

	if price < 0 {
		return nil, domain.ErrInvalidInput
	}

	course := &domain.Course{
		Title:              title,
		Description:        description,
		Thumbnail:          thumbnail,
		AdminID:            adminID,
		MentorID:           mentorID,
		CategoryID:         categoryID,
		Price:              price,
		Hours:              hours,
		Visibility:         "public",
		Tools:              tools,
		LearningObjectives: learningObjectives,
	}

	if err := s.courseRepo.Create(ctx, course); err != nil {
		return nil, err
	}

	return course, nil
}

func toCourseFilter(f *CourseListFilter) *repository.CourseListFilter {
	if f == nil {
		return nil
	}
	return &repository.CourseListFilter{
		CategoryID: f.CategoryID,
		MinPrice:   f.MinPrice,
		MaxPrice:   f.MaxPrice,
	}
}

func (s *CourseService) ListPublicCourses(
	ctx context.Context,
	page, limit int,
	filter *CourseListFilter,
) ([]dto.CourseListItem, int64, error) {

	courses, total, err := s.courseRepo.ListPublic(ctx, page, limit, toCourseFilter(filter))
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.CourseListItem, len(courses))
	for i, p := range courses {
		result[i] = dto.CourseToListItem(p)
	}

	return result, total, nil
}

func (s *CourseService) ListPublicCoursesWithEnrollment(
	ctx context.Context,
	userID string,
	page, limit int,
	filter *CourseListFilter,
) ([]dto.CourseListItem, int64, error) {

	courses, total, err := s.courseRepo.ListPublic(ctx, page, limit, toCourseFilter(filter))
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.CourseListItem, len(courses))
	for i, p := range courses {
		item := dto.CourseToListItem(p)
		enrolled, _ := s.enrollRepo.Exists(ctx, userID, p.ID)
		item.IsEnrolled = enrolled
		result[i] = item
	}

	return result, total, nil
}

func (s *CourseService) UpdateCourse(ctx context.Context, courseID, adminID string, title, description string, price float64, visibility string, tools datatypes.JSON, learningObjectives datatypes.JSON) error {
	course, err := s.courseRepo.FindByID(ctx, courseID)
	if err != nil {
		return err
	}

	if course.AdminID != adminID {
		return domain.ErrForbidden
	}

	if title != "" {
		course.Title = title
	}
	if description != "" {
		course.Description = description
	}
	if price >= 0 {
		course.Price = price
	}
	if visibility != "" {
		course.Visibility = visibility
	}
	if len(tools) > 0 {
		course.Tools = tools
	}
	if len(learningObjectives) > 0 {
		course.LearningObjectives = learningObjectives
	}

	return s.courseRepo.Update(ctx, course)
}

func (s *CourseService) DeleteCourse(ctx context.Context, courseID, adminID string) error {
	course, err := s.courseRepo.FindByID(ctx, courseID)
	if err != nil {
		return err
	}

	if course.AdminID != adminID {
		return domain.ErrForbidden
	}

	return s.courseRepo.DeleteByID(ctx, courseID)
}

func (s *CourseService) EnrollCourse(ctx context.Context, userID, courseID string) error {
	// Check course exists
	_, err := s.courseRepo.FindByID(ctx, courseID)
	if err != nil {
		return domain.ErrCourseNotFound
	}

	// Check if already enrolled
	exists, err := s.enrollRepo.Exists(ctx, userID, courseID)
	if err != nil {
		return err
	}
	if exists {
		return nil // Already enrolled, no-op
	}

	enrollment := &domain.CourseEnrollment{
		UserID:   userID,
		CourseID: courseID,
	}

	if err := s.enrollRepo.Create(ctx, enrollment); err != nil {
		return err
	}

	// Notify course admin about new enrollment
	course, err := s.courseRepo.FindByID(ctx, courseID)
	if err == nil {
		notifSvc := GetNotificationService()
		if notifSvc != nil {
			// Notify admin
			notifSvc.CreateNotification(ctx, course.AdminID, "course_enrollment",
				"Pendaftar Baru",
				"Seseorang telah mendaftar di kursus: "+course.Title,
				"/courses/"+courseID,
			)
			// Notify enrolled user
			notifSvc.CreateNotification(ctx, userID, "enrollment_success",
				"Pendaftaran Berhasil",
				"Selamat! Anda berhasil mendaftar di kursus: "+course.Title,
				"/courses/"+courseID,
			)
		}
	}

	return nil
}

func (s *CourseService) SetLastLesson(ctx context.Context, userID, courseID, lessonID string) error {
	return s.enrollRepo.UpdateLastLesson(ctx, userID, courseID, lessonID)
}

func (s *CourseService) IsEnrolled(ctx context.Context, userID, courseID string) (bool, error) {
	if userID == "" {
		return false, nil
	}
	return s.enrollRepo.Exists(ctx, userID, courseID)
}

func (s *CourseService) ListEnrolledCourses(
	ctx context.Context,
	userID string,
	page,
	limit int,
) ([]dto.EnrolledCourseItem, int64, error) {

	enrollments, total, err := s.enrollRepo.ListByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.EnrolledCourseItem, len(enrollments))
	for i, e := range enrollments {
		item := dto.CourseToListItem(e.Course)
		result[i] = dto.EnrolledCourseItem{
			CourseListItem: item,
			EnrolledAt:     e.CreatedAt,
			LastLessonID:   e.LastLessonID,
		}
	}

	return result, total, nil
}
