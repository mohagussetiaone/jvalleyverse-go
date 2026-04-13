package service

import (
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/repository"
)

type ReplyService struct {
	replyRepo      *repository.ReplyRepository
	discussionRepo *repository.DiscussionRepository
	userService    *UserService
}

func NewReplyService() *ReplyService {
	return &ReplyService{
		replyRepo:      repository.NewReplyRepository(),
		discussionRepo: repository.NewDiscussionRepository(),
		userService:    NewUserService(),
	}
}

// CreateReply creates new reply (can be nested)
func (s *ReplyService) CreateReply(userID, discussionID uint, content string, parentID *uint) (*domain.Reply, error) {
	reply := &domain.Reply{
		UserID:       userID,
		DiscussionID: discussionID,
		Content:      content,
		ParentID:     parentID,
		LikesCount:   0,
		IsMarkedBest: false,
	}

	if err := s.replyRepo.Create(reply); err != nil {
		return nil, err
	}

	// Award points for replying
	s.userService.AddPoints(userID, "create_reply", 5, map[string]interface{}{
		"discussion_id": discussionID,
		"reply_id":      reply.ID,
	})

	return reply, nil
}

// UpdateReply updates reply (owner only)
func (s *ReplyService) UpdateReply(replyID, userID uint, content string) error {
	reply, err := s.replyRepo.FindByID(replyID)
	if err != nil {
		return err
	}

	// Only owner can update
	if reply.UserID != userID {
		return nil
	}

	reply.Content = content
	return s.replyRepo.Update(reply)
}

// DeleteReply deletes reply and cascade nested replies (owner or admin)
func (s *ReplyService) DeleteReply(replyID, userID uint, isAdmin bool) error {
	reply, err := s.replyRepo.FindByID(replyID)
	if err != nil {
		return err
	}

	// Only owner or admin can delete
	if reply.UserID != userID && !isAdmin {
		return nil
	}

	return s.replyRepo.DeleteByID(replyID)
}

// LikeReply adds like to reply and awards points to creator
func (s *ReplyService) LikeReply(userID, replyID uint) error {
	reply, err := s.replyRepo.FindByID(replyID)
	if err != nil {
		return err
	}

	// Increment likes
	if err := s.replyRepo.IncrementLikes(replyID); err != nil {
		return err
	}

	// Award points to reply creator if not self-like
	if reply.UserID != userID {
		s.userService.AddPoints(reply.UserID, "receive_reply_like", 3, map[string]interface{}{
			"reply_id": replyID,
			"from_user": userID,
		})
	}

	return nil
}

// MarkBestReply marks reply as best answer (discussion owner only)
func (s *ReplyService) MarkBestReply(replyID, discussionID, userID uint) error {
	discussion, err := s.discussionRepo.FindByID(discussionID)
	if err != nil {
		return err
	}

	// Only discussion owner can mark best
	if discussion.UserID != userID {
		return nil
	}

	reply, err := s.replyRepo.FindByID(replyID)
	if err != nil {
		return err
	}

	reply.IsMarkedBest = true
	if err := s.replyRepo.Update(reply); err != nil {
		return err
	}

	// Award points to reply creator
	s.userService.AddPoints(reply.UserID, "best_answer", 25, map[string]interface{}{
		"reply_id":       replyID,
		"discussion_id":  discussionID,
	})

	return nil
}
