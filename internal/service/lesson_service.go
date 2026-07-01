package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/dto"
	"jvalleyverse/internal/repository"
	"jvalleyverse/pkg/config"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ILessonService interface {
	GetPublicLessonByID(ctx context.Context, lessonID string) (*dto.LessonDetailResponse, error)
	GetLessonBySlug(ctx context.Context, courseID string, slug string, userID string) (*dto.LessonDetailResponse, error)
	StartLesson(ctx context.Context, userID, lessonID string) (*domain.LessonProgress, error)
	UpdateProgress(ctx context.Context, userID, lessonID string, percentage int, notes string) (*domain.LessonProgress, error)
	CompleteLesson(ctx context.Context, userID, lessonID string) (map[string]interface{}, error)
	AdminCreateLesson(ctx context.Context, adminID string, input domain.Lesson) (*domain.Lesson, error)
	AdminUpdateLesson(ctx context.Context, adminID, id string, input domain.Lesson) (*domain.Lesson, error)
	AdminDeleteLesson(ctx context.Context, adminID, id string) error
	AdminCreateLessonDetail(ctx context.Context, lessonID string, about, rules string, tools, media, resources interface{}) (*domain.LessonDetail, error)
	ListLessonsByCourse(ctx context.Context, courseID string, limit, offset int) ([]domain.Lesson, int64, error)
	ListLessonsBySection(ctx context.Context, sectionID string) ([]domain.Lesson, int64, error)
}

type LessonService struct {
	lessonRepo         *repository.LessonRepository
	lessonDetailRepo   *repository.LessonDetailRepository
	lessonProgressRepo *repository.LessonProgressRepository
	certRepo           *repository.CertificateRepository
	userService        IUserService
	courseRepo         *repository.CourseRepository
	enrollRepo         *repository.EnrollmentRepository
	sectionRepo        *repository.SectionRepository
	streakRepo         *repository.LearningStreakRepository
}

func NewLessonService(
	lessonRepo *repository.LessonRepository,
	lessonDetailRepo *repository.LessonDetailRepository,
	lessonProgressRepo *repository.LessonProgressRepository,
	certRepo *repository.CertificateRepository,
	userService IUserService,
	courseRepo *repository.CourseRepository,
	enrollRepo *repository.EnrollmentRepository,
	sectionRepo *repository.SectionRepository,
	streakRepo *repository.LearningStreakRepository,
) *LessonService {
	return &LessonService{
		lessonRepo:         lessonRepo,
		lessonDetailRepo:   lessonDetailRepo,
		lessonProgressRepo: lessonProgressRepo,
		certRepo:           certRepo,
		userService:        userService,
		courseRepo:         courseRepo,
		enrollRepo:         enrollRepo,
		sectionRepo:        sectionRepo,
		streakRepo:         streakRepo,
	}
}

func (s *LessonService) GetLessonBySlug(ctx context.Context, courseID string, slug string, userID string) (*dto.LessonDetailResponse, error) {
	lesson, err := s.lessonRepo.FindBySlug(ctx, courseID, slug)
	if err != nil {
		return nil, err
	}

	var progress *domain.LessonProgress
	if userID != "" {
		progress, err = s.lessonProgressRepo.FindByUserAndLesson(ctx, userID, lesson.ID)
		if err != nil {
			progress = &domain.LessonProgress{
				Status:             "not_started",
				ProgressPercentage: 0,
			}
		}
	}

	return &dto.LessonDetailResponse{
		Lesson:     dto.ToLessonBrief(*lesson),
		Details:    lesson.Details,
		Progress:   progress,
		NextLesson: dto.ToLessonBriefPtr(lesson.NextLesson),
		Section:    dto.ToSectionBriefPtr(&lesson.Section),
		Course:     dto.CourseToListItemPtr(&lesson.Course),
	}, nil
}

func (s *LessonService) GetPublicLessonByID(
	ctx context.Context,
	lessonID string,
) (*dto.LessonDetailResponse, error) {

	lesson, err := s.lessonRepo.FindPublicByID(ctx, lessonID)
	if err != nil {
		return nil, err
	}

	return &dto.LessonDetailResponse{
		Lesson:     dto.ToLessonBrief(*lesson),
		Details:    lesson.Details,
		NextLesson: dto.ToLessonBriefPtr(lesson.NextLesson),
		Section:    dto.ToSectionBriefPtr(&lesson.Section),
		Course:     dto.CourseToListItemPtr(&lesson.Course),
	}, nil
}

func (s *LessonService) StartLesson(ctx context.Context, userID, lessonID string) (*domain.LessonProgress, error) {
	if _, err := s.lessonRepo.FindByID(ctx, lessonID); err != nil {
		return nil, domain.ErrLessonNotFound
	}

	progress, err := s.lessonProgressRepo.FindByUserAndLesson(ctx, userID, lessonID)
	if err == nil {
		if progress.Status == "not_started" {
			progress.Status = "started"
			now := time.Now()
			progress.StartedAt = &now
			err = s.lessonProgressRepo.Update(ctx, progress)
			// Auto-track last lesson
			if err == nil {
				lesson, findErr := s.lessonRepo.FindByID(ctx, lessonID)
				if findErr == nil && lesson.CourseID != "" {
					s.enrollRepo.UpdateLastLesson(ctx, userID, lesson.CourseID, lessonID)
				}
			}
			return progress, err
		}
		return progress, nil
	}

	now := time.Now()
	progress = &domain.LessonProgress{
		UserID:             userID,
		LessonID:           lessonID,
		Status:             "started",
		StartedAt:          &now,
		ProgressPercentage: 0,
	}
	err = s.lessonProgressRepo.Create(ctx, progress)
	if err == nil {
		lesson, findErr := s.lessonRepo.FindByID(ctx, lessonID)
		if findErr == nil && lesson.CourseID != "" {
			s.enrollRepo.UpdateLastLesson(ctx, userID, lesson.CourseID, lessonID)
		}
	}
	return progress, err
}

func (s *LessonService) UpdateProgress(ctx context.Context, userID, lessonID string, percentage int, notes string) (*domain.LessonProgress, error) {
	if percentage < 0 || percentage > 100 {
		return nil, domain.ErrInvalidInput
	}

	progress, err := s.lessonProgressRepo.FindByUserAndLesson(ctx, userID, lessonID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	progress.ProgressPercentage = percentage
	progress.Notes = notes
	switch {
	case percentage == 0:
		progress.Status = "started"
	case percentage > 0 && percentage < 100:
		progress.Status = "in_progress"
	case percentage == 100:
		progress.Status = "completed"
		now := time.Now()
		progress.CompletedAt = &now
	}

	err = s.lessonProgressRepo.Update(ctx, progress)

	// Auto-track last lesson
	lesson, findErr := s.lessonRepo.FindByID(ctx, lessonID)
	if findErr == nil && lesson.CourseID != "" {
		s.enrollRepo.UpdateLastLesson(ctx, userID, lesson.CourseID, lessonID)
	}

	return progress, err
}

func (s *LessonService) CompleteLesson(ctx context.Context, userID, lessonID string) (map[string]interface{}, error) {
	progress, err := s.lessonProgressRepo.FindByUserAndLesson(ctx, userID, lessonID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	if progress.Status == "completed" {
		return nil, errors.New("lesson already completed")
	}

	now := time.Now()
	progress.Status = "completed"
	progress.ProgressPercentage = 100
	progress.CompletedAt = &now

	if err := s.lessonProgressRepo.Update(ctx, progress); err != nil {
		return nil, err
	}

	// Auto-track last lesson
	lesson, findErr := s.lessonRepo.FindByID(ctx, lessonID)
	if findErr == nil && lesson.CourseID != "" {
		s.enrollRepo.UpdateLastLesson(ctx, userID, lesson.CourseID, lessonID)
	}

	code := generateUniqueCertificateCode()
	baseURL := "https://jvalleyverse.com"
	if config.AppConfig != nil && config.AppConfig.FrontendURL != "" {
		baseURL = config.AppConfig.FrontendURL
	}
	verificationURL := fmt.Sprintf("%s/certificates/verify/%s", baseURL, code)
	qrCodeURL := fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=%s", verificationURL)
	cert := &domain.Certificate{
		UserID:          userID,
		LessonID:        lessonID,
		UniqueCode:      code,
		IssuedAt:        now,
		VerificationURL: verificationURL,
		QRCodeURL:       qrCodeURL,
	}
	if err := s.certRepo.Create(ctx, cert); err != nil {
		return nil, err
	}

	if err := s.userService.AddPoints(ctx, userID, "complete_lesson", 50, map[string]interface{}{
		"lesson_id": lessonID,
		"cert_code": code,
	}); err != nil {
		return nil, err
	}

	lesson, err = s.lessonRepo.FindByID(ctx, lessonID)
	if err != nil {
		return nil, err
	}

	// Notify user about lesson completion and certificate
	if notifSvc := GetNotificationService(); notifSvc != nil {
		notifSvc.CreateNotification(ctx, userID, "lesson_completed",
			"Pelajaran Selesai!",
			"Selamat! Anda telah menyelesaikan '"+lesson.Title+"' dan mendapat sertifikat.",
			"/courses/"+lesson.CourseID+"/lessons/"+lesson.Slug,
		)
	}

	// Update learning streak
	s.updateLearningStreak(ctx, userID, now)

	return map[string]interface{}{
		"message":     "Lesson completed!",
		"certificate": cert,
		"achievement": map[string]interface{}{
			"type":        "certificate",
			"title":       lesson.Title,
			"unique_code": cert.UniqueCode,
			"issued_at":   cert.IssuedAt,
		},
		"progress":       progress,
		"next_lesson":    lesson.NextLesson,
		"points_awarded": 50,
	}, nil
}

func (s *LessonService) updateLearningStreak(ctx context.Context, userID string, now time.Time) {
	streak, err := s.streakRepo.FindByUserID(ctx, userID)
	if err != nil {
		// No streak record yet, create one
		streak = &domain.LearningStreak{
			UserID:           userID,
			StreakCount:      1,
			LongestStreak:    1,
			LastActivityDate: now,
		}
		s.streakRepo.Upsert(ctx, streak)
		return
	}

	today := now.Truncate(24 * time.Hour)
	lastActivity := streak.LastActivityDate.Truncate(24 * time.Hour)

	if today.Equal(lastActivity) {
		// Same day activity, don't change streak
		return
	}

	yesterday := today.AddDate(0, 0, -1)
	if lastActivity.Equal(yesterday) {
		// Consecutive day, increment streak
		streak.StreakCount++
	} else {
		// Streak broken, reset to 1
		streak.StreakCount = 1
	}

	// Update longest streak if current streak is longer
	if streak.StreakCount > streak.LongestStreak {
		streak.LongestStreak = streak.StreakCount
	}

	streak.LastActivityDate = now
	s.streakRepo.Upsert(ctx, streak)
}

func generateUniqueCertificateCode() string {
	return "CERT-" + uuid.New().String()[:8]
}

func (s *LessonService) AdminCreateLesson(ctx context.Context, adminID string, input domain.Lesson) (*domain.Lesson, error) {
	course, err := s.courseRepo.FindByID(ctx, input.CourseID)
	if err != nil {
		return nil, domain.ErrCourseNotFound
	}
	if course.AdminID != adminID {
		return nil, domain.ErrForbidden
	}

	if input.SectionID == "" {
		return nil, domain.ErrInvalidInput
	}

	exists, err := s.sectionRepo.ExistsInCourse(ctx, input.SectionID, input.CourseID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrInvalidInput
	}

	input.AdminID = adminID
	if err := s.lessonRepo.Create(ctx, &input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *LessonService) AdminUpdateLesson(ctx context.Context, adminID, id string, input domain.Lesson) (*domain.Lesson, error) {
	lesson, err := s.lessonRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if lesson.AdminID != adminID {
		return nil, domain.ErrForbidden
	}

	if input.CourseID != "" && input.CourseID != lesson.CourseID {
		course, err := s.courseRepo.FindByID(ctx, input.CourseID)
		if err != nil {
			return nil, domain.ErrCourseNotFound
		}
		if course.AdminID != adminID {
			return nil, domain.ErrForbidden
		}
		lesson.CourseID = input.CourseID
	}

	if input.SectionID != "" {
		exists, err := s.sectionRepo.ExistsInCourse(ctx, input.SectionID, lesson.CourseID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, domain.ErrInvalidInput
		}
		lesson.SectionID = input.SectionID
	}

	if input.Title != "" {
		lesson.Title = input.Title
	}
	if input.Description != "" {
		lesson.Description = input.Description
	}
	if input.Slug != "" {
		lesson.Slug = input.Slug
	}
	if input.Thumbnail != "" {
		lesson.Thumbnail = input.Thumbnail
	}
	if input.VideoURL != "" {
		lesson.VideoURL = input.VideoURL
	}
	if input.Difficulty != "" {
		lesson.Difficulty = input.Difficulty
	}
	if input.Duration > 0 {
		lesson.Duration = input.Duration
	}
	lesson.OrderIndex = input.OrderIndex
	lesson.SequenceNum = input.SequenceNum
	lesson.IsFirst = input.IsFirst
	lesson.NextLessonID = input.NextLessonID
	if input.Visibility != "" {
		lesson.Visibility = input.Visibility
	}

	if err := s.lessonRepo.Update(ctx, lesson); err != nil {
		return nil, err
	}

	updated, err := s.lessonRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *LessonService) AdminDeleteLesson(ctx context.Context, adminID, id string) error {
	lesson, err := s.lessonRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if lesson.AdminID != adminID {
		return domain.ErrForbidden
	}

	return s.lessonRepo.DeleteByID(ctx, id)
}

func (s *LessonService) AdminCreateLessonDetail(ctx context.Context, lessonID string, about, rules string, tools, media, resources interface{}) (*domain.LessonDetail, error) {
	toolsJSON, _ := json.Marshal(tools)
	mediaJSON, _ := json.Marshal(media)
	resourcesJSON, _ := json.Marshal(resources)

	detail := &domain.LessonDetail{
		LessonID:      lessonID,
		About:         about,
		Rules:         rules,
		Tools:         datatypes.JSON(toolsJSON),
		ResourceMedia: datatypes.JSON(mediaJSON),
		Resources:     datatypes.JSON(resourcesJSON),
	}

	if err := s.lessonDetailRepo.Create(ctx, detail); err != nil {
		return nil, err
	}
	return detail, nil
}

func (s *LessonService) ListLessonsByCourse(ctx context.Context, courseID string, limit, offset int) ([]domain.Lesson, int64, error) {
	return s.lessonRepo.ListByCourseID(ctx, courseID, limit, offset)
}

func (s *LessonService) ListLessonsBySection(ctx context.Context, sectionID string) ([]domain.Lesson, int64, error) {
	return s.lessonRepo.ListBySectionID(ctx, sectionID)
}
