package service

import (
	"context"
	"jvalleyverse/internal/repository"
)

// IGamificationService defines the business logic for the point and level system
type IGamificationService interface {
	AwardPoints(ctx context.Context, userID string, activityType string, points int, metadata map[string]interface{}) error
	GetLeaderboard(ctx context.Context, limit int) ([]map[string]interface{}, error)
	GetUserActivityLog(ctx context.Context, userID string, page, limit int) ([]map[string]interface{}, error)
	GetLevelInfo() []map[string]interface{}
	GetUserStats(ctx context.Context, userID string) (map[string]interface{}, error)
}

type GamificationService struct {
	pointRepo    *repository.CommunityPointRepository
	levelRepo    *repository.UserLevelRepository
	userRepo     *repository.UserRepository
	showcaseRepo *repository.ShowcaseRepository
	userService  IUserService
}

func NewGamificationService(
	pointRepo *repository.CommunityPointRepository,
	levelRepo *repository.UserLevelRepository,
	userRepo *repository.UserRepository,
	showcaseRepo *repository.ShowcaseRepository,
	userService IUserService,
) *GamificationService {
	return &GamificationService{
		pointRepo:    pointRepo,
		levelRepo:    levelRepo,
		userRepo:     userRepo,
		showcaseRepo: showcaseRepo,
		userService:  userService,
	}
}

// AwardPoints awards points to user for activity
func (s *GamificationService) AwardPoints(ctx context.Context, userID string, activityType string, points int, metadata map[string]interface{}) error {
	return s.userService.AddPoints(ctx, userID, activityType, points, metadata)
}

// GetLeaderboard returns top users by points
func (s *GamificationService) GetLeaderboard(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	if limit <= 0 {
		limit = 10
	}

	users, err := s.userRepo.GetTopByPoints(ctx, limit)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(users))
	for i, u := range users {
		result[i] = map[string]interface{}{
			"rank":         i + 1,
			"user_id":      u.ID,
			"name":         u.Name,
			"avatar":       u.Avatar,
			"total_points": u.TotalPoints,
			"level":        u.Level,
		}
	}

	return result, nil
}

// GetUserActivityLog returns user's activity history
func (s *GamificationService) GetUserActivityLog(ctx context.Context, userID string, page, limit int) ([]map[string]interface{}, error) {
	activities, _, err := s.pointRepo.ListByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(activities))
	for i, activity := range activities {
		result[i] = map[string]interface{}{
			"id":        activity.ID,
			"activity":  activity.ActivityType,
			"points":    activity.PointsEarned,
			"timestamp": activity.CreatedAt,
		}
	}

	return result, nil
}

// GetLevelInfo returns all level requirements
func (s *GamificationService) GetLevelInfo() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":        "Beginner",
			"threshold":   0,
			"color":       "#6366f1",
			"description": "Just starting your journey",
		},
		{
			"name":        "Intermediate",
			"threshold":   100,
			"color":       "#8b5cf6",
			"description": "Building momentum",
		},
		{
			"name":        "Advanced",
			"threshold":   500,
			"color":       "#d946ef",
			"description": "Getting serious",
		},
		{
			"name":        "Expert",
			"threshold":   1000,
			"color":       "#ec4899",
			"description": "Mastering skills",
		},
		{
			"name":        "Master",
			"threshold":   2000,
			"color":       "#f43f5e",
			"description": "Peak achievement",
		},
	}
}

// GetUserStats returns comprehensive user statistics
func (s *GamificationService) GetUserStats(ctx context.Context, userID string) (map[string]interface{}, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	activityLog, _ := s.GetUserActivityLog(ctx, userID, 1, 10)

	return map[string]interface{}{
		"user_id":         user.ID,
		"name":            user.Name,
		"total_points":    user.TotalPoints,
		"current_level":   user.Level,
		"recent_activity": activityLog,
	}, nil
}

func (s *GamificationService) calculateLevel(points int) int {
	if points >= 2000 {
		return 5
	} else if points >= 1000 {
		return 4
	} else if points >= 500 {
		return 3
	} else if points >= 200 {
		return 2
	}
	return 1
}
