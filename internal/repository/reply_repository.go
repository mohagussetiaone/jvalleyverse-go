package repository

import (
	"context"
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

// ReplyRepository handles reply data access
type ReplyRepository struct {
	db *gorm.DB
}

func NewReplyRepository(db *gorm.DB) *ReplyRepository {
	return &ReplyRepository{db: db}
}

// Create creates new reply
func (r *ReplyRepository) Create(ctx context.Context, reply *domain.Reply) error {
	return r.db.WithContext(ctx).Create(reply).Error
}

// FindByID finds reply with user info
func (r *ReplyRepository) FindByID(ctx context.Context, replyID string) (*domain.Reply, error) {
	reply := &domain.Reply{}
	if err := r.db.WithContext(ctx).Where("id = ?", replyID).Preload("User").First(reply).Error; err != nil {
		return nil, err
	}
	return reply, nil
}

// ListByDiscussionID lists direct replies (parent_id IS NULL) for a discussion
func (r *ReplyRepository) ListByDiscussionID(ctx context.Context, discussionID string) ([]domain.Reply, error) {
	var replies []domain.Reply
	err := r.db.WithContext(ctx).
		Where("discussion_id = ? AND parent_id IS NULL", discussionID).
		Preload("User").
		Order("created_at ASC").
		Find(&replies).Error
	return replies, err
}

// ListNestedByParentID lists nested replies under a parent reply
func (r *ReplyRepository) ListNestedByParentID(ctx context.Context, parentID string) ([]domain.Reply, error) {
	var replies []domain.Reply
	err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).
		Preload("User").
		Order("created_at ASC").
		Find(&replies).Error
	return replies, err
}

// Update updates reply
func (r *ReplyRepository) Update(ctx context.Context, reply *domain.Reply) error {
	return r.db.WithContext(ctx).Model(reply).Updates(reply).Error
}

// DeleteByID deletes reply and cascade nested replies
func (r *ReplyRepository) DeleteByID(ctx context.Context, replyID string) error {
	return r.db.WithContext(ctx).Where("id = ?", replyID).Delete(&domain.Reply{}).Error
}

// IncrementLikes increments like count on reply
func (r *ReplyRepository) IncrementLikes(ctx context.Context, replyID string) error {
	return r.db.WithContext(ctx).Model(&domain.Reply{}).Where("id = ?", replyID).
		Update("likes_count", gorm.Expr("likes_count + 1")).Error
}
