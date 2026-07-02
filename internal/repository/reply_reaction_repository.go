package repository

import (
	"context"
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ReplyReactionRepository handles reply reaction data access
type ReplyReactionRepository struct {
	db *gorm.DB
}

func NewReplyReactionRepository(db *gorm.DB) *ReplyReactionRepository {
	return &ReplyReactionRepository{db: db}
}

// React creates a reaction (ignore if already exists due to composite PK)
func (r *ReplyReactionRepository) React(ctx context.Context, userID, replyID, emoji string) error {
	reaction := &domain.ReplyReaction{
		UserID:  userID,
		ReplyID: replyID,
		Emoji:   emoji,
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(reaction).Error
}

// Unreact removes a reaction
func (r *ReplyReactionRepository) Unreact(ctx context.Context, userID, replyID, emoji string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND reply_id = ? AND emoji = ?", userID, replyID, emoji).
		Delete(&domain.ReplyReaction{}).Error
}

// GetSummaryByReplyID returns emoji counts and which ones the user reacted with
type ReactionCount struct {
	Emoji string `json:"emoji"`
	Count int    `json:"count"`
}

func (r *ReplyReactionRepository) GetSummaryByReplyID(ctx context.Context, replyID string) ([]ReactionCount, error) {
	var counts []ReactionCount
	err := r.db.WithContext(ctx).
		Model(&domain.ReplyReaction{}).
		Select("emoji, COUNT(*) as count").
		Where("reply_id = ?", replyID).
		Group("emoji").
		Order("count DESC").
		Find(&counts).Error
	return counts, err
}

// GetAllReactionsByReplyIDs returns all reactions for multiple reply IDs
func (r *ReplyReactionRepository) GetAllReactionsByReplyIDs(ctx context.Context, replyIDs []string) ([]domain.ReplyReaction, error) {
	if len(replyIDs) == 0 {
		return nil, nil
	}
	var reactions []domain.ReplyReaction
	err := r.db.WithContext(ctx).
		Where("reply_id IN ?", replyIDs).
		Find(&reactions).Error
	return reactions, err
}

// GetUserReactionsByReplyIDs returns the user's reactions for a set of reply IDs
func (r *ReplyReactionRepository) GetUserReactionsByReplyIDs(ctx context.Context, userID string, replyIDs []string) ([]domain.ReplyReaction, error) {
	if len(replyIDs) == 0 {
		return nil, nil
	}
	var reactions []domain.ReplyReaction
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND reply_id IN ?", userID, replyIDs).
		Find(&reactions).Error
	return reactions, err
}
