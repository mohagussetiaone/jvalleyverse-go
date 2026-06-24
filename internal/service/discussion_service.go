package service

import (
	"context"
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/dto"
	"jvalleyverse/internal/repository"
)

// IDiscussionService defines the business logic for managing discussion threads
type IDiscussionService interface {
	CreateDiscussion(ctx context.Context, userID string, title, content string, lessonID *string, studyCaseID *string, categoryID string) (*domain.Discussion, error)
	ListDiscussions(ctx context.Context, page, limit int, lessonID *string, studyCaseID *string, status *string) ([]dto.DiscussionListItem, int64, error)
	ListUserDiscussions(ctx context.Context, userID string, page, limit int) ([]dto.DiscussionListItem, int64, error)
	GetDiscussionWithReplies(ctx context.Context, discussionID string) (*dto.DiscussionDetail, error)
	UpdateDiscussion(ctx context.Context, discussionID, userID string, title, content string) error
	CloseDiscussion(ctx context.Context, discussionID, userID string) error
	DeleteDiscussion(ctx context.Context, discussionID, userID string, isAdmin bool) error
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

func (s *DiscussionService) CreateDiscussion(ctx context.Context, userID string, title, content string, lessonID *string, studyCaseID *string, categoryID string) (*domain.Discussion, error) {
	discussion := &domain.Discussion{
		UserID:      userID,
		LessonID:    lessonID,
		StudyCaseID: studyCaseID,
		CategoryID:  categoryID,
		Title:       title,
		Content:     content,
		Status:      "open",
		ViewsCount:  0,
	}

	if err := s.discussionRepo.Create(ctx, discussion); err != nil {
		return nil, err
	}

	return discussion, nil
}

func (s *DiscussionService) ListDiscussions(ctx context.Context, page, limit int, lessonID *string, studyCaseID *string, status *string) ([]dto.DiscussionListItem, int64, error) {
	var discussions []domain.Discussion
	var total int64
	var err error

	if studyCaseID != nil && *studyCaseID != "" {
		discussions, total, err = s.discussionRepo.ListByStudyCaseID(ctx, *studyCaseID, page, limit)
	} else if lessonID != nil && *lessonID != "" {
		discussions, total, err = s.discussionRepo.ListByLessonID(ctx, *lessonID, page, limit)
	} else {
		discussions, total, err = s.discussionRepo.ListAll(ctx, page, limit)
	}

	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.DiscussionListItem, len(discussions))
	for i, d := range discussions {
		result[i] = dto.ToDiscussionListItem(d)
	}

	return result, total, nil
}

// GetDiscussionWithReplies returns discussion with all threaded replies
func (s *DiscussionService) GetDiscussionWithReplies(ctx context.Context, discussionID string) (*dto.DiscussionDetail, error) {
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

	replyData := make([]dto.ReplyInDiscussion, len(replies))
	for i, reply := range replies {
		replyData[i] = dto.ToReplyInDiscussion(reply)
	}

	return &dto.DiscussionDetail{
		ID:         discussion.ID,
		Title:      discussion.Title,
		Content:    discussion.Content,
		User:       dto.ToUserBrief(discussion.User),
		Status:     discussion.Status,
		ViewsCount: discussion.ViewsCount,
		CreatedAt:  discussion.CreatedAt,
		Replies:    replyData,
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

func (s *DiscussionService) ListUserDiscussions(ctx context.Context, userID string, page, limit int) ([]dto.DiscussionListItem, int64, error) {
	discussions, total, err := s.discussionRepo.ListByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.DiscussionListItem, len(discussions))
	for i, d := range discussions {
		result[i] = dto.ToDiscussionListItem(d)
	}

	return result, total, nil
}

// DeleteDiscussion deletes a discussion (owner or admin)
func (s *DiscussionService) DeleteDiscussion(ctx context.Context, discussionID, userID string, isAdmin bool) error {
	discussion, err := s.discussionRepo.FindByID(ctx, discussionID)
	if err != nil {
		return err
	}

	if discussion.UserID != userID && !isAdmin {
		return domain.ErrForbidden
	}

	return s.discussionRepo.DeleteByID(ctx, discussionID)
}
