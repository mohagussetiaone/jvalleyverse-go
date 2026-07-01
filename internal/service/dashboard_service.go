package service

import (
	"context"
	"jvalleyverse/internal/repository"
)

// IDashboardService defines dashboard business logic
type IDashboardService interface {
	GetDashboard(ctx context.Context, userID string) (map[string]interface{}, error)
	GetStreak(ctx context.Context, userID string) (int, error)
}

type DashboardService struct {
	lessonRepo      *repository.LessonRepository
	progressRepo    *repository.LessonProgressRepository
	notifSvc        INotificationService
	streakRepo      *repository.LearningStreakRepository
}

func NewDashboardService(
	lessonRepo *repository.LessonRepository,
	progressRepo *repository.LessonProgressRepository,
	notifSvc INotificationService,
	streakRepo *repository.LearningStreakRepository,
) *DashboardService {
	return &DashboardService{
		lessonRepo:   lessonRepo,
		progressRepo: progressRepo,
		notifSvc:     notifSvc,
		streakRepo:   streakRepo,
	}
}

func (s *DashboardService) GetDashboard(ctx context.Context, userID string) (map[string]interface{}, error) {
	progresses, err := s.progressRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	lessonsInProgress := 0
	lessonsCompleted := 0
	lessonsStarted := 0

	for _, p := range progresses {
		switch p.Status {
		case "in_progress":
			lessonsInProgress++
		case "completed":
			lessonsCompleted++
		case "started":
			lessonsStarted++
		}
	}

	unreadNotifs, _ := s.notifSvc.CountUnread(ctx, userID)

	// Get streak data
	streakCount := 0
	streak, err := s.streakRepo.FindByUserID(ctx, userID)
	if err == nil {
		streakCount = streak.StreakCount
	}

	return map[string]interface{}{
		"courses_in_progress":  lessonsInProgress,
		"courses_completed":    lessonsCompleted,
		"courses_dropped":      lessonsStarted,
		"unread_notifications": unreadNotifs,
		"streak_count":         streakCount,
	}, nil
}

func (s *DashboardService) GetStreak(ctx context.Context, userID string) (int, error) {
	streak, err := s.streakRepo.FindByUserID(ctx, userID)
	if err != nil {
		return 0, nil // Return 0 if no streak record exists yet
	}
	return streak.StreakCount, nil
}
