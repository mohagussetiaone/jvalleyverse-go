package repository

import (
	"context"

	"jvalleyverse/internal/domain"
	"gorm.io/gorm"
)

// ReplyLikeRepository handles reply like data access
type ReplyLikeRepository struct {
	db *gorm.DB
}

func NewReplyLikeRepository(db *gorm.DB) *ReplyLikeRepository {
	return &ReplyLikeRepository{db: db}
}

// Exists checks if a user has liked a reply
func (r *ReplyLikeRepository) Exists(ctx context.Context, userID, replyID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.ReplyLike{}).
		Where("user_id = ? AND reply_id = ?", userID, replyID).
		Count(&count).Error
	return count > 0, err
}

// Create adds a like record
func (r *ReplyLikeRepository) Create(ctx context.Context, userID, replyID string) error {
	like := &domain.ReplyLike{
		UserID:  userID,
		ReplyID: replyID,
	}
	return r.db.WithContext(ctx).Create(like).Error
}

// Delete removes a like record
// FindLikedReplyIDs returns reply IDs that a user has liked
func (r *ReplyLikeRepository) FindLikedReplyIDs(ctx context.Context, userID string, replyIDs []string) (map[string]bool, error) {
	if len(replyIDs) == 0 {
		return make(map[string]bool), nil
	}
	var likedReplyIDs []string
	err := r.db.WithContext(ctx).Model(&domain.ReplyLike{}).
		Where("user_id = ? AND reply_id IN ?", userID, replyIDs).
		Pluck("reply_id", &likedReplyIDs).Error
	if err != nil {
		return nil, err
	}
	result := make(map[string]bool, len(likedReplyIDs))
	for _, id := range likedReplyIDs {
		result[id] = true
	}
	return result, nil
}

func (r *ReplyLikeRepository) Delete(ctx context.Context, userID, replyID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND reply_id = ?", userID, replyID).
		Delete(&domain.ReplyLike{}).Error
}
