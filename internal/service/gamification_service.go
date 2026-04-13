package service

import (
	"jvalleyverse/internal/repository"
)

type GamificationService struct {
	pointRepo  *repository.CommunityPointRepository
	levelRepo  *repository.UserLevelRepository
	userRepo   *repository.UserRepository
	showcaseRepo *repository.ShowcaseRepository
}

func NewGamificationService() *GamificationService {
	return &GamificationService{
		pointRepo:    repository.NewCommunityPointRepository(),
		levelRepo:    repository.NewUserLevelRepository(),
		userRepo:     repository.NewUserRepository(),
		showcaseRepo: repository.NewShowcaseRepository(),
	}
}

// AwardPoints awards points to user for activity
func (s *GamificationService) AwardPoints(userID uint, activityType string, points int, metadata map[string]interface{}) error {
	userSvc := NewUserService()
	return userSvc.AddPoints(userID, activityType, points, metadata)
}

// GetLeaderboard returns top users by points
func (s *GamificationService) GetLeaderboard(limit int) ([]map[string]interface{}, error) {
	// Top users overall - TODO: Query directly for better performance
	// For now, using repository methods
	result := []map[string]interface{}{}

	// Would normally do a SQL query like:
	// SELECT user_id, SUM(points) as total_points FROM community_points GROUP BY user_id ORDER BY total_points DESC LIMIT ?
	
	// For MVP, return empty and implement database query optimization later
	return result, nil
}

// GetUserActivityLog returns user's activity history
func (s *GamificationService) GetUserActivityLog(userID uint, page, limit int) ([]map[string]interface{}, error) {
	activities, _, err := s.pointRepo.ListByUserID(userID, page, limit)
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
func (s *GamificationService) GetUserStats(userID uint) (map[string]interface{}, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	level, _ := s.levelRepo.FindByUserID(userID)
	activityLog, _ := s.GetUserActivityLog(userID, 1, 10)

	return map[string]interface{}{
		"user_id":       user.ID,
		"name":          user.Name,
		"total_points":  user.TotalPoints,
		"current_level": level.Level,
		"recent_activity": activityLog,
	}, nil
}

// Points Distribution System:
// - Create showcase: 10 pts
// - Receive showcase like: 5 pts
// - Create discussion: 5 pts
// - Create reply: 5 pts
// - Receive reply like: 3 pts
// - Best answer selected: 25 pts
// - Complete class: 50 pts
//
// Level Thresholds:
// - Beginner: 0-99 pts
// - Intermediate: 100-499 pts
// - Advanced: 500-999 pts
// - Expert: 1000-1999 pts
// - Master: 2000+ pts
// - Discussion reply: 2 pts
// - Class completed: 100 pts
// - Certificate issued: 50 pts

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
