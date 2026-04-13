package service

import (
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/repository"
)

type DiscussionService struct {
	discussionRepo *repository.DiscussionRepository
	replyRepo      *repository.ReplyRepository
	userRepo       *repository.UserRepository
}

func NewDiscussionService() *DiscussionService {
	return &DiscussionService{
		discussionRepo: repository.NewDiscussionRepository(),
		replyRepo:      repository.NewReplyRepository(),
		userRepo:       repository.NewUserRepository(),
	}
}

// CreateDiscussion creates new discussion thread
func (s *DiscussionService) CreateDiscussion(userID uint, title, content string, classID uint) (*domain.Discussion, error) {
	discussion := &domain.Discussion{
		UserID:    userID,
		ClassID:   &classID,
		Title:     title,
		Content:   content,
		Status:    "open",
		ViewsCount: 0,
	}

	if err := s.discussionRepo.Create(discussion); err != nil {
		return nil, err
	}

	return discussion, nil
}

// ListDiscussions returns paginated discussions
func (s *DiscussionService) ListDiscussions(page, limit int, classID *uint, status *string) ([]map[string]interface{}, error) {
	var discussions []domain.Discussion
	var err error

	if classID != nil {
		discussions, _, err = s.discussionRepo.ListByClassID(*classID, page, limit)
	} else {
		// TODO: Implement general list all discussions
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
func (s *DiscussionService) GetDiscussionWithReplies(discussionID uint) (map[string]interface{}, error) {
	discussion, err := s.discussionRepo.FindByID(discussionID)
	if err != nil {
		return nil, err
	}

	// Increment view count
	discussion.ViewsCount++
	s.discussionRepo.Update(discussion)

	// Get all replies
	replies, err := s.replyRepo.ListByDiscussionID(discussionID)
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
func (s *DiscussionService) UpdateDiscussion(discussionID, userID uint, title, content string) error {
	discussion, err := s.discussionRepo.FindByID(discussionID)
	if err != nil {
		return err
	}

	// Only owner can update
	if discussion.UserID != userID {
		return nil // Silent fail or return error based on preference
	}

	discussion.Title = title
	discussion.Content = content
	return s.discussionRepo.Update(discussion)
}

// CloseDiscussion closes a discussion (owner or admin)
func (s *DiscussionService) CloseDiscussion(discussionID, userID uint) error {
	discussion, err := s.discussionRepo.FindByID(discussionID)
	if err != nil {
		return err
	}

	if discussion.UserID != userID {
		return nil // Only owner can close
	}

	discussion.Status = "closed"
	return s.discussionRepo.Update(discussion)
}
