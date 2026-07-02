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

	// Notify creator as activity history
	if notifSvc := GetNotificationService(); notifSvc != nil {
		notifSvc.CreateNotification(ctx, userID, "discussion_created",
			"Diskusi Baru Dibuat",
			"Anda membuat diskusi: "+title,
			"/discussions/"+discussion.ID,
		)
	}

	// Notify study case owner if discussion is related to a study case
	if studyCaseID != nil && *studyCaseID != "" {
		studyCase, err := s.discussionRepo.FindStudyCaseByID(ctx, *studyCaseID)
		if err == nil && studyCase.UserID != userID {
			if notifSvc := GetNotificationService(); notifSvc != nil {
				notifSvc.CreateNotification(ctx, studyCase.UserID, "discussion_on_study_case",
					"Diskusi Baru di Study Case",
					"Ada diskusi baru di study case '"+studyCase.Name+"': "+title,
					"/discussions/"+discussion.ID,
				)
			}
		}
	}

	// Notify lesson admin if discussion is related to a lesson
	if lessonID != nil && *lessonID != "" {
		lesson, err := s.discussionRepo.FindLessonByID(ctx, *lessonID)
		if err == nil && lesson.AdminID != userID {
			if notifSvc := GetNotificationService(); notifSvc != nil {
				notifSvc.CreateNotification(ctx, lesson.AdminID, "discussion_on_lesson",
					"Diskusi Baru di Lesson",
					"Ada diskusi baru di lesson '"+lesson.Title+"': "+title,
					"/discussions/"+discussion.ID,
				)
			}
		}
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

// GetDiscussionWithReplies returns discussion with all threaded replies (unlimited depth)
func (s *DiscussionService) GetDiscussionWithReplies(ctx context.Context, discussionID string) (*dto.DiscussionDetail, error) {
	discussion, err := s.discussionRepo.FindByID(ctx, discussionID)
	if err != nil {
		return nil, err
	}

	// Increment view count
	discussion.ViewsCount++
	s.discussionRepo.Update(ctx, discussion)

	// Get ALL replies (not just top-level)
	allReplies, err := s.replyRepo.ListAllByDiscussionID(ctx, discussionID)
	if err != nil {
		return nil, err
	}

	// Collect all reply IDs for batch reaction loading
	replyIDs := make([]string, len(allReplies))
	for i, r := range allReplies {
		replyIDs[i] = r.ID
	}

	// Convert to DTOs first
	replyMap := make(map[string]dto.ReplyInDiscussion, len(allReplies))
	for _, r := range allReplies {
		replyMap[r.ID] = dto.ToReplyInDiscussion(r)
	}

	// Fetch reactions for all replies (anonymous user — no "reacted" status)
	if svc := GetReplyService(); svc != nil {
		reactions, err := svc.GetReactionsForReplies(ctx, replyIDs, "")
		if err == nil {
			for replyID, rxn := range reactions {
				if rdto, ok := replyMap[replyID]; ok {
					rdto.Reactions = rxn
					replyMap[replyID] = rdto
				}
			}
		}
	}

	// Build tree: top-level replies (parent_id IS NULL)
	var topLevel []dto.ReplyInDiscussion
	for _, r := range allReplies {
		if r.ParentID == nil || *r.ParentID == "" {
			rdto := replyMap[r.ID]
			rdto.Children = buildReplyTree(r.ID, allReplies, replyMap)
			topLevel = append(topLevel, rdto)
		}
	}

	return &dto.DiscussionDetail{
		ID:         discussion.ID,
		Title:      discussion.Title,
		Content:    discussion.Content,
		User:       dto.ToUserBrief(discussion.User),
		Status:     discussion.Status,
		ViewsCount: discussion.ViewsCount,
		CreatedAt:  discussion.CreatedAt,
		Replies:    topLevel,
	}, nil
}

// buildReplyTree recursively builds nested reply tree for a parent reply
func buildReplyTree(parentID string, allReplies []domain.Reply, replyMap map[string]dto.ReplyInDiscussion) []dto.ReplyInDiscussion {
	var children []dto.ReplyInDiscussion
	for _, r := range allReplies {
		if r.ParentID != nil && *r.ParentID == parentID {
			rdto := replyMap[r.ID]
			rdto.Children = buildReplyTree(r.ID, allReplies, replyMap)
			children = append(children, rdto)
		}
	}
	return children
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
