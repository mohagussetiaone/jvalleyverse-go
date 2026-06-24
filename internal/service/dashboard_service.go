package service

import (
	"context"
	"jvalleyverse/internal/repository"
)

// IDashboardService defines dashboard business logic
type IDashboardService interface {
	GetDashboard(ctx context.Context, userID string) (map[string]interface{}, error)
}

type DashboardService struct {
	lessonRepo      *repository.LessonRepository
	progressRepo    *repository.LessonProgressRepository
	notifSvc        INotificationService
}

func NewDashboardService(
	lessonRepo *repository.LessonRepository,
	progressRepo *repository.LessonProgressRepository,
	notifSvc INotificationService,
) *DashboardService {
	return &DashboardService{
		lessonRepo:   lessonRepo,
		progressRepo: progressRepo,
		notifSvc:     notifSvc,
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

	return map[string]interface{}{
		"courses_in_progress":  lessonsInProgress,
		"courses_completed":    lessonsCompleted,
		"courses_dropped":      lessonsStarted,
		"unread_notifications": unreadNotifs,
	}, nil
}
