package service

import (
	"context"
	"jvalleyverse/internal/dto"
	"jvalleyverse/internal/repository"
)

// IGamificationService defines the business logic for the point and level system
type IGamificationService interface {
	AwardPoints(ctx context.Context, userID string, activityType string, points int, metadata map[string]interface{}) error
	GetLeaderboard(ctx context.Context, limit int) ([]dto.LeaderboardItem, error)
	GetUserActivityLog(ctx context.Context, userID string, page, limit int) ([]dto.ActivityItem, error)
	GetLevelInfo() []dto.LevelInfo
	GetUserStats(ctx context.Context, userID string) (*dto.UserStats, error)
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
func (s *GamificationService) GetLeaderboard(ctx context.Context, limit int) ([]dto.LeaderboardItem, error) {
	if limit <= 0 {
		limit = 10
	}

	users, err := s.userRepo.GetTopByPoints(ctx, limit)
	if err != nil {
		return nil, err
	}

	result := make([]dto.LeaderboardItem, len(users))
	for i, u := range users {
		result[i] = dto.LeaderboardItem{
			Rank:        i + 1,
			UserID:      u.ID,
			Name:        u.Name,
			Avatar:      u.Avatar,
			TotalPoints: u.TotalPoints,
			Level:       u.Level,
		}
	}

	return result, nil
}

// GetUserActivityLog returns user's activity history
func (s *GamificationService) GetUserActivityLog(ctx context.Context, userID string, page, limit int) ([]dto.ActivityItem, error) {
	activities, _, err := s.pointRepo.ListByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	result := make([]dto.ActivityItem, len(activities))
	for i, activity := range activities {
		result[i] = dto.ActivityItem{
			ID:        activity.ID,
			Activity:  activity.ActivityType,
			Points:    activity.PointsEarned,
			Timestamp: activity.CreatedAt,
		}
	}

	return result, nil
}

// GetLevelInfo returns all level requirements
func (s *GamificationService) GetLevelInfo() []dto.LevelInfo {
	return []dto.LevelInfo{
		{Name: "Beginner", Threshold: 0, Color: "#6366f1", Description: "Just starting your journey"},
		{Name: "Intermediate", Threshold: 100, Color: "#8b5cf6", Description: "Building momentum"},
		{Name: "Advanced", Threshold: 500, Color: "#d946ef", Description: "Getting serious"},
		{Name: "Expert", Threshold: 1000, Color: "#ec4899", Description: "Mastering skills"},
		{Name: "Master", Threshold: 2000, Color: "#f43f5e", Description: "Peak achievement"},
	}
}

// GetUserStats returns comprehensive user statistics
func (s *GamificationService) GetUserStats(ctx context.Context, userID string) (*dto.UserStats, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	activityLog, _ := s.GetUserActivityLog(ctx, userID, 1, 10)

	return &dto.UserStats{
		UserID:         user.ID,
		Name:           user.Name,
		TotalPoints:    user.TotalPoints,
		CurrentLevel:   user.Level,
		RecentActivity: activityLog,
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
