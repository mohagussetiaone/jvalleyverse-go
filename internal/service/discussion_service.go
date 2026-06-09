package service

import (
	"context"
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/repository"
)

// IDiscussionService defines the business logic for managing discussion threads
type IDiscussionService interface {
	CreateDiscussion(ctx context.Context, userID string, title, content string, classID *string, categoryID string) (*domain.Discussion, error)
	ListDiscussions(ctx context.Context, page, limit int, classID *string, status *string) ([]map[string]interface{}, error)
	GetDiscussionWithReplies(ctx context.Context, discussionID string) (map[string]interface{}, error)
	UpdateDiscussion(ctx context.Context, discussionID, userID string, title, content string) error
	CloseDiscussion(ctx context.Context, discussionID, userID string) error
}

type DiscussionService struct {
	discussionRepo *repository.DiscussionRepository
	replyRepo      *repository.ReplyRepository
	userRepo       *repository.UserRepository
}

func NewDiscussionService(
	discussionRepo *repository.DiscussionRepository,
	replyRepo *repository.ReplyRepository,
	userRepo *repository.UserRepository,
) *DiscussionService {
	return &DiscussionService{
		discussionRepo: discussionRepo,
		replyRepo:      replyRepo,
		userRepo:       userRepo,
	}
}

// CreateDiscussion creates new discussion thread
func (s *DiscussionService) CreateDiscussion(ctx context.Context, userID string, title, content string, classID *string, categoryID string) (*domain.Discussion, error) {
	discussion := &domain.Discussion{
		UserID:     userID,
		ClassID:    classID,
		CategoryID: categoryID,
		Title:      title,
		Content:    content,
		Status:     "open",
		ViewsCount: 0,
	}

	if err := s.discussionRepo.Create(ctx, discussion); err != nil {
		return nil, err
	}

	return discussion, nil
}

// ListDiscussions returns paginated discussions
func (s *DiscussionService) ListDiscussions(ctx context.Context, page, limit int, classID *string, status *string) ([]map[string]interface{}, error) {
	var discussions []domain.Discussion
	var err error

	if classID != nil {
		discussions, _, err = s.discussionRepo.ListByClassID(ctx, *classID, page, limit)
	} else {
		discussions = []domain.Discussion{}
	}

	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(discussions))
	for i, d := range discussions {
		result[i] = map[string]interface{}{
			"id":         d.ID,
			"title":      d.Title,
			"content":    d.Content,
			"user_id":    d.UserID,
			"user_name":  d.User.Name,
			"class_id":   d.ClassID,
			"status":     d.Status,
			"view_count": d.ViewsCount,
			"created_at": d.CreatedAt,
		}
	}

	return result, nil
}

// GetDiscussionWithReplies returns discussion with all threaded replies
func (s *DiscussionService) GetDiscussionWithReplies(ctx context.Context, discussionID string) (map[string]interface{}, error) {
	discussion, err := s.discussionRepo.FindByID(ctx, discussionID)
	if err != nil {
		return nil, err
	}

	// Increment view count
	discussion.ViewsCount++
	s.discussionRepo.Update(ctx, discussion)

	// Get all replies
	replies, err := s.replyRepo.ListByDiscussionID(ctx, discussionID)
	if err != nil {
		return nil, err
	}

	replyData := make([]map[string]interface{}, len(replies))
	for i, reply := range replies {
		replyData[i] = map[string]interface{}{
			"id":         reply.ID,
			"content":    reply.Content,
			"user_id":    reply.UserID,
			"user_name":  reply.User.Name,
			"likes":      reply.LikesCount,
			"is_best":    reply.IsMarkedBest,
			"created_at": reply.CreatedAt,
		}
	}

	return map[string]interface{}{
		"id":         discussion.ID,
		"title":      discussion.Title,
		"content":    discussion.Content,
		"user_id":    discussion.UserID,
		"user_name":  discussion.User.Name,
		"status":     discussion.Status,
		"view_count": discussion.ViewsCount,
		"created_at": discussion.CreatedAt,
		"replies":    replyData,
	}, nil
}

// UpdateDiscussion updates discussion (owner only)
func (s *DiscussionService) UpdateDiscussion(ctx context.Context, discussionID, userID string, title, content string) error {
	discussion, err := s.discussionRepo.FindByID(ctx, discussionID)
	if err != nil {
		return err
	}

	// Only owner can update
	if discussion.UserID != userID {
		return domain.ErrForbidden
	}

	discussion.Title = title
	discussion.Content = content
	return s.discussionRepo.Update(ctx, discussion)
}

// CloseDiscussion closes a discussion (owner or admin)
func (s *DiscussionService) CloseDiscussion(ctx context.Context, discussionID, userID string) error {
	discussion, err := s.discussionRepo.FindByID(ctx, discussionID)
	if err != nil {
		return err
	}

	if discussion.UserID != userID {
		return domain.ErrForbidden
	}

	discussion.Status = "closed"
	return s.discussionRepo.Update(ctx, discussion)
}
