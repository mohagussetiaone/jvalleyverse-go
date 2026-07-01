package service

import (
	"context"
	"jvalleyverse/internal/dto"
	"jvalleyverse/internal/repository"
)

type StreakService struct {
	streakRepo *repository.LearningStreakRepository
}

func NewStreakService(streakRepo *repository.LearningStreakRepository) *StreakService {
	return &StreakService{streakRepo: streakRepo}
}

// GetUserStreak returns the user's current streak info
func (s *StreakService) GetUserStreak(ctx context.Context, userID string) (*dto.StreakItem, error) {
	streak, err := s.streakRepo.FindByUserID(ctx, userID)
	if err != nil {
		// No streak record yet, return zero values
		return &dto.StreakItem{
			StreakCount:   0,
			LongestStreak: 0,
		}, nil
	}

	return &dto.StreakItem{
		StreakCount:      streak.StreakCount,
		LongestStreak:    streak.LongestStreak,
		LastActivityDate: streak.LastActivityDate,
	}, nil
}
