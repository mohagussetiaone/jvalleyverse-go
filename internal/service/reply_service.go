package service

import (
	"context"
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/dto"
	"jvalleyverse/internal/repository"
)

// IReplyService defines the business logic for managing replies within discussions
type IReplyService interface {
	CreateReply(ctx context.Context, userID, discussionID string, content string, parentID *string) (*domain.Reply, error)
	UpdateReply(ctx context.Context, replyID, userID string, content string) error
	DeleteReply(ctx context.Context, replyID, userID string, isAdmin bool) error
	LikeReply(ctx context.Context, userID, replyID string) error
	MarkBestReply(ctx context.Context, replyID, discussionID, userID string) error
	ListRepliesByUser(ctx context.Context, userID string, page, limit int) ([]dto.ReplyListItem, int64, error)
	ReactToReply(ctx context.Context, userID, replyID, emoji string) error
	UnreactFromReply(ctx context.Context, userID, replyID, emoji string) error
	GetReactionsForReplies(ctx context.Context, replyIDs []string, currentUserID string) (map[string][]dto.ReactionSummary, error)
}

type ReplyService struct {
	replyLikeRepo  *repository.ReplyLikeRepository
	replyRepo      *repository.ReplyRepository
	reactRepo      *repository.ReplyReactionRepository
	discussionRepo *repository.DiscussionRepository
	userService    IUserService
}

func NewReplyService(
	replyRepo *repository.ReplyRepository,
	reactRepo *repository.ReplyReactionRepository,
	replyLikeRepo *repository.ReplyLikeRepository,
	discussionRepo *repository.DiscussionRepository,
	userService IUserService,
) *ReplyService {
	return &ReplyService{
		replyRepo:      replyRepo,
		reactRepo:      reactRepo,
		discussionRepo: discussionRepo,
		replyLikeRepo:  replyLikeRepo,
		userService:    userService,
	}
}

// ListRepliesByUser returns paginated replies by a user
func (s *ReplyService) ListRepliesByUser(ctx context.Context, userID string, page, limit int) ([]dto.ReplyListItem, int64, error) {
	replies, total, err := s.replyRepo.ListByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.ReplyListItem, len(replies))
	for i, r := range replies {
		discussionTitle := ""
		if r.Discussion.ID != "" {
			discussionTitle = r.Discussion.Title
		}

		result[i] = dto.ReplyListItem{
			ID:              r.ID,
			Content:         r.Content,
			DiscussionID:    r.DiscussionID,
			DiscussionTitle: discussionTitle,
			ParentID:        r.ParentID,
			LikesCount:      r.LikesCount,
			IsMarkedBest:    r.IsMarkedBest,
			CreatedAt:       r.CreatedAt,
		}
	}

	return result, total, nil
}

// CreateReply creates new reply (can be nested)
func (s *ReplyService) CreateReply(ctx context.Context, userID, discussionID string, content string, parentID *string) (*domain.Reply, error) {
	reply := &domain.Reply{
		UserID:       userID,
		DiscussionID: discussionID,
		Content:      content,
		ParentID:     parentID,
		LikesCount:   0,
		IsMarkedBest: false,
	}

	if err := s.replyRepo.Create(ctx, reply); err != nil {
		return nil, err
	}

	// Award points for replying
	s.userService.AddPoints(ctx, userID, "create_reply", 5, map[string]interface{}{
		"discussion_id": discussionID,
		"reply_id":      reply.ID,
	})

	// Look up discussion for notification targets
	discussion, err := s.discussionRepo.FindByID(ctx, discussionID)

	// Self-notification as activity history
	if err == nil {
		if notifSvc := GetNotificationService(); notifSvc != nil {
			notifSvc.CreateNotification(ctx, userID, "reply_created",
				"Balasan Baru Dikirim",
				"Anda membalas diskusi: "+discussion.Title,
				"/discussions/"+discussionID,
			)
		}
	}

	// Notify discussion owner (if replier is not the owner)
	if err == nil && discussion.UserID != userID {
		if notifSvc := GetNotificationService(); notifSvc != nil {
			notifSvc.CreateNotification(ctx, discussion.UserID, "new_reply",
				"Diskusi Anda Dibalas",
				"Seseorang membalas diskusi Anda: "+discussion.Title,
				"/discussions/"+discussionID,
			)
		}
	}

	// Notify parent reply owner for nested replies
	if parentID != nil && *parentID != "" && discussion != nil {
		parentReply, err := s.replyRepo.FindByID(ctx, *parentID)
		if err == nil && parentReply.UserID != userID && parentReply.UserID != discussion.UserID {
			if notifSvc := GetNotificationService(); notifSvc != nil {
				notifSvc.CreateNotification(ctx, parentReply.UserID, "nested_reply",
					"Balasan Anda Dibalas",
					"Seseorang membalas komentar Anda dalam diskusi: "+discussion.Title,
					"/discussions/"+discussionID,
				)
			}
		}
	}

	return reply, nil
}

// UpdateReply updates reply (owner only)
func (s *ReplyService) UpdateReply(ctx context.Context, replyID, userID string, content string) error {
	reply, err := s.replyRepo.FindByID(ctx, replyID)
	if err != nil {
		return err
	}

	if reply.UserID != userID {
		return domain.ErrForbidden
	}

	reply.Content = content
	return s.replyRepo.Update(ctx, reply)
}

// DeleteReply deletes reply and cascade nested replies (owner or admin)
func (s *ReplyService) DeleteReply(ctx context.Context, replyID, userID string, isAdmin bool) error {
	reply, err := s.replyRepo.FindByID(ctx, replyID)
	if err != nil {
		return err
	}

	if reply.UserID != userID && !isAdmin {
		return domain.ErrForbidden
	}

	return s.replyRepo.DeleteByID(ctx, replyID)
}

// LikeReply adds like to reply and awards points to creator
func (s *ReplyService) LikeReply(ctx context.Context, userID, replyID string) error {
	reply, err := s.replyRepo.FindByID(ctx, replyID)
	if err != nil {
		return err
	}

	if err := s.replyRepo.IncrementLikes(ctx, replyID); err != nil {
		return err
	}

	// Award points to reply creator if not self-like
	if reply.UserID != userID {
		s.userService.AddPoints(ctx, reply.UserID, "receive_reply_like", 3, map[string]interface{}{
			"reply_id":  replyID,
			"from_user": userID,
		})
		// Notify reply creator
		if notifSvc := GetNotificationService(); notifSvc != nil {
			notifSvc.CreateNotification(ctx, reply.UserID, "reply_like",
				"Balasan Anda Disukai",
				"Seseorang menyukai balasan Anda dalam diskusi.",
				"/discussions/"+reply.DiscussionID,
			)
		}
	}

	return nil
}

// allowedEmojis is the set of supported reaction emojis
var allowedEmojis = map[string]bool{
	"\U0001f44d": true, // 👍
	"\u2764\ufe0f": true, // ❤️
	"\U0001f602": true, // 😂
	"\U0001f62e": true, // 😮
	"\U0001f622": true, // 😢
	"\U0001f621": true, // 😡
}

// ReactToReply adds an emoji reaction (validates emoji is allowed)
func (s *ReplyService) ReactToReply(ctx context.Context, userID, replyID, emoji string) error {
	if !allowedEmojis[emoji] {
		return domain.ErrInvalidInput
	}
	// Verify reply exists
	if _, err := s.replyRepo.FindByID(ctx, replyID); err != nil {
		return domain.ErrNotFound
	}
	return s.reactRepo.React(ctx, userID, replyID, emoji)
}

// UnreactFromReply removes an emoji reaction
func (s *ReplyService) UnreactFromReply(ctx context.Context, userID, replyID, emoji string) error {
	if !allowedEmojis[emoji] {
		return domain.ErrInvalidInput
	}
	// Verify reply exists
	if _, err := s.replyRepo.FindByID(ctx, replyID); err != nil {
		return domain.ErrNotFound
	}
	return s.reactRepo.Unreact(ctx, userID, replyID, emoji)
}

// GetReactionsForReplies returns reaction summaries for a set of reply IDs
func (s *ReplyService) GetReactionsForReplies(ctx context.Context, replyIDs []string, currentUserID string) (map[string][]dto.ReactionSummary, error) {
	if len(replyIDs) == 0 {
		return nil, nil
	}

	// Get all reactions for these replies
	reactions, err := s.reactRepo.GetAllReactionsByReplyIDs(ctx, replyIDs)
	if err != nil {
		return nil, err
	}

	// Get current user's reactions
	userReactions, err := s.reactRepo.GetUserReactionsByReplyIDs(ctx, currentUserID, replyIDs)
	if err != nil {
		return nil, err
	}

	// Build user react set for quick lookup: map[replyID]map[emoji]bool
	userReactMap := make(map[string]map[string]bool)
	for _, ur := range userReactions {
		if userReactMap[ur.ReplyID] == nil {
			userReactMap[ur.ReplyID] = make(map[string]bool)
		}
		userReactMap[ur.ReplyID][ur.Emoji] = true
	}

	// Build reaction summaries per reply
	// map[replyID]map[emoji]count
	reactCounts := make(map[string]map[string]int)
	for _, r := range reactions {
		if reactCounts[r.ReplyID] == nil {
			reactCounts[r.ReplyID] = make(map[string]int)
		}
		reactCounts[r.ReplyID][r.Emoji]++
	}

	result := make(map[string][]dto.ReactionSummary)
	emojiOrder := []string{"\U0001f44d", "\u2764\ufe0f", "\U0001f602", "\U0001f62e", "\U0001f622", "\U0001f621"}

	for replyID, emojis := range reactCounts {
		summaries := make([]dto.ReactionSummary, 0)
		for _, emoji := range emojiOrder {
			if count, ok := emojis[emoji]; ok {
				summaries = append(summaries, dto.ReactionSummary{
					Emoji:   emoji,
					Count:   count,
					Reacted: userReactMap[replyID][emoji],
				})
			}
		}
		// Also include any emojis not in the predefined order
		for emoji, count := range emojis {
			found := false
			for _, e := range emojiOrder {
				if e == emoji {
					found = true
					break
				}
			}
			if !found {
				summaries = append(summaries, dto.ReactionSummary{
					Emoji:   emoji,
					Count:   count,
					Reacted: userReactMap[replyID][emoji],
				})
			}
		}
		result[replyID] = summaries
	}

	return result, nil
}

// MarkBestReply marks reply as best answer (discussion owner only)
func (s *ReplyService) MarkBestReply(ctx context.Context, replyID, discussionID, userID string) error {
	discussion, err := s.discussionRepo.FindByID(ctx, discussionID)
	if err != nil {
		return err
	}

	if discussion.UserID != userID {
		return domain.ErrForbidden
	}

	reply, err := s.replyRepo.FindByID(ctx, replyID)
	if err != nil {
		return err
	}

	reply.IsMarkedBest = true
	if err := s.replyRepo.Update(ctx, reply); err != nil {
		return err
	}

	// Award points to reply creator
	s.userService.AddPoints(ctx, reply.UserID, "best_answer", 25, map[string]interface{}{
		"reply_id":      replyID,
		"discussion_id": discussionID,
	})

	// Notify reply creator
	if notifSvc := GetNotificationService(); notifSvc != nil {
		notifSvc.CreateNotification(ctx, reply.UserID, "best_answer",
			"Jawaban Anda Terpilih sebagai Terbaik",
			"Jawaban Anda dipilih sebagai jawaban terbaik dalam diskusi: "+discussion.Title,
			"/discussions/"+discussionID,
		)
	}

	return nil
}
