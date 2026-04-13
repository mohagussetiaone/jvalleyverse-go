package repository

import (
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

// ReplyRepository handles reply data access
type ReplyRepository struct {
	db *gorm.DB
}

func NewReplyRepository() *ReplyRepository {
	return &ReplyRepository{db: db}
}

// Create creates new reply
func (r *ReplyRepository) Create(reply *domain.Reply) error {
	return r.db.Create(reply).Error
}

// FindByID finds reply with nested replies
func (r *ReplyRepository) FindByID(replyID uint) (*domain.Reply, error) {
	reply := &domain.Reply{}
	if err := r.db.Preload("User").First(reply, replyID).Error; err != nil {
		return nil, err
	}
	return reply, nil
}

// ListByDiscussionID lists direct replies (parent_id IS NULL) for a discussion with nested replies
func (r *ReplyRepository) ListByDiscussionID(discussionID uint) ([]domain.Reply, error) {
	var replies []domain.Reply
	// Get direct replies
	err := r.db.Where("discussion_id = ? AND (parent_id = 0 OR parent_id IS NULL)", discussionID).
		Preload("User").
		Order("created_at ASC").
		Find(&replies).Error
	return replies, err
}

// ListNestedByParentID lists nested replies under a parent reply
func (r *ReplyRepository) ListNestedByParentID(parentID uint) ([]domain.Reply, error) {
	var replies []domain.Reply
	err := r.db.Where("parent_id = ?", parentID).
		Preload("User").
		Order("created_at ASC").
		Find(&replies).Error
	return replies, err
}

// Update updates reply
func (r *ReplyRepository) Update(reply *domain.Reply) error {
	return r.db.Model(reply).Updates(reply).Error
}

// DeleteByID deletes reply and cascade nested replies
func (r *ReplyRepository) DeleteByID(replyID uint) error {
	return r.db.Delete(&domain.Reply{}, replyID).Error
}

// IncrementLikes increments like count on reply
func (r *ReplyRepository) IncrementLikes(replyID uint) error {
	return r.db.Model(&domain.Reply{}).Where("id = ?", replyID).
		Update("likes_count", gorm.Expr("likes_count + 1")).Error
}
