package repository

import (
	"context"
	"jvalleyverse/internal/domain"
	"time"

	"gorm.io/gorm"
)

type LearningStreakRepository struct {
	db *gorm.DB
}

func NewLearningStreakRepository(db *gorm.DB) *LearningStreakRepository {
	return &LearningStreakRepository{db: db}
}

// FindByUserID retrieves streak record for a user
func (r *LearningStreakRepository) FindByUserID(ctx context.Context, userID string) (*domain.LearningStreak, error) {
	streak := &domain.LearningStreak{}
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(streak).Error; err != nil {
		return nil, err
	}
	return streak, nil
}

// Upsert creates or updates streak record when user completes a lesson
func (r *LearningStreakRepository) Upsert(ctx context.Context, streak *domain.LearningStreak) error {
	// Use OnConflict for PostgreSQL: insert or update
	return r.db.WithContext(ctx).Where("user_id = ?", streak.UserID).Assign(*streak).FirstOrCreate(streak).Error
}

// GetTopStreaks returns top N users by streak count
func (r *LearningStreakRepository) GetTopStreaks(ctx context.Context, limit int) ([]domain.LearningStreak, error) {
	var streaks []domain.LearningStreak
	err := r.db.WithContext(ctx).
		Preload("User").
		Order("streak_count DESC").
		Limit(limit).
		Find(&streaks).Error
	return streaks, err
}

// GetActiveUserIDs returns user IDs who had activity on a specific date
func (r *LearningStreakRepository) GetActiveUserIDs(ctx context.Context, date time.Time) ([]string, error) {
	var userIDs []string
	err := r.db.WithContext(ctx).
		Model(&domain.LearningStreak{}).
		Where("DATE(last_activity_date) = DATE(?)", date).
		Pluck("user_id", &userIDs).Error
	return userIDs, err
}
