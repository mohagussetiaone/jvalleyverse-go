package service

import (
	"context"
	"encoding/json"
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/repository"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// IUserService defines the business logic for user management and gamification
type IUserService interface {
	GetProfile(ctx context.Context, userID string) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID string, name, bio, avatar string) error
	AddPoints(ctx context.Context, userID string, category string, points int, metadata map[string]interface{}) error
	GetUserActivityLog(ctx context.Context, userID string, page, limit int) ([]map[string]interface{}, error)
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	ListAllUsers(ctx context.Context, page, limit int) ([]map[string]interface{}, int64, error)
}

type UserService struct {
	userRepo  *repository.UserRepository
	pointRepo *repository.CommunityPointRepository
	levelRepo *repository.UserLevelRepository
}

func NewUserService(userRepo *repository.UserRepository, pointRepo *repository.CommunityPointRepository, levelRepo *repository.UserLevelRepository) *UserService {
	return &UserService{
		userRepo:  userRepo,
		pointRepo: pointRepo,
		levelRepo: levelRepo,
	}
}

// GetProfile returns user profile
func (s *UserService) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}

// UpdateProfile updates user profile
func (s *UserService) UpdateProfile(ctx context.Context, userID string, name, bio, avatar string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Name = name
	user.Bio = bio
	user.Avatar = avatar

	return s.userRepo.Update(ctx, user)
}

// AddPoints adds points to user and updates level if threshold reached
func (s *UserService) AddPoints(ctx context.Context, userID string, category string, points int, metadata map[string]interface{}) error {
	// Add points to user
	if err := s.userRepo.UpdatePoints(ctx, userID, points); err != nil {
		return err
	}

	// Record transaction
	meta, _ := json.Marshal(metadata)
	point := &domain.CommunityPoint{
		UserID:       userID,
		PointsEarned: points,
		ActivityType: category,
		Metadata:     datatypes.JSON(meta),
	}
	if err := s.pointRepo.Create(ctx, point); err != nil {
		return err
	}

	// Recalculate level
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	newLevel := calculateLevel(user.TotalPoints)
	oldLevel := user.Level

	// Only update if level changed
	if oldLevel != newLevel {
		if err := s.userRepo.UpdateLevel(ctx, userID, newLevel); err != nil {
			return err
		}
	}

	return nil
}

// calculateLevel determines user level based on total points
func calculateLevel(totalPoints int) int {
	if totalPoints < 100 {
		return 1
	} else if totalPoints < 500 {
		return 2
	} else if totalPoints < 1000 {
		return 3
	} else if totalPoints < 2000 {
		return 4
	}
	return 5
}

// GetUserActivityLog returns user's activity points history
func (s *UserService) GetUserActivityLog(ctx context.Context, userID string, page, limit int) ([]map[string]interface{}, error) {
	points, _, err := s.pointRepo.ListByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(points))
	for i, p := range points {
		result[i] = map[string]interface{}{
			"id":        p.ID,
			"activity":  p.ActivityType,
			"points":    p.PointsEarned,
			"timestamp": p.CreatedAt,
		}
	}
	return result, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, user *domain.User) error {
	return s.userRepo.Create(ctx, user)
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}

// ListAllUsers returns paginated list of all users (admin only)
func (s *UserService) ListAllUsers(ctx context.Context, page, limit int) ([]map[string]interface{}, int64, error) {
	users, total, err := s.userRepo.ListAll(ctx, page, limit)
	if err != nil {
		return nil, 0, err
	}

	result := make([]map[string]interface{}, len(users))
	for i, user := range users {
		result[i] = map[string]interface{}{
			"id":           user.ID,
			"email":        user.Email,
			"name":         user.Name,
			"avatar":       user.Avatar,
			"role":         user.Role,
			"level":        user.Level,
			"total_points": user.TotalPoints,
			"is_active":    user.IsActive,
			"created_at":   user.CreatedAt,
		}
	}

	return result, total, nil
}

// IShowcaseService defines the business logic for managing user showcases
type IShowcaseService interface {
	CreateShowcase(ctx context.Context, userID string, title string, description string, mediaURLs []string, categoryID string, visibility string) (*domain.Showcase, error)
	ListShowcases(ctx context.Context, page, limit int, categoryID string, sort string) ([]domain.Showcase, int64, error)
	GetShowcaseByID(ctx context.Context, showcaseID string) (*domain.Showcase, error)
	UpdateShowcase(ctx context.Context, showcaseID string, userID string, title string, description string, visibility string) (*domain.Showcase, error)
	DeleteShowcase(ctx context.Context, showcaseID string, userID string) error
	LikeShowcase(ctx context.Context, userID, showcaseID string) error
	UnlikeShowcase(ctx context.Context, userID, showcaseID string) error
}

type ShowcaseService struct {
	showcaseRepo *repository.ShowcaseRepository
	likeRepo     *repository.ShowcaseLikeRepository
	userService  IUserService
}

func NewShowcaseService(
	showcaseRepo *repository.ShowcaseRepository,
	likeRepo *repository.ShowcaseLikeRepository,
	userService IUserService,
) *ShowcaseService {
	return &ShowcaseService{
		showcaseRepo: showcaseRepo,
		likeRepo:     likeRepo,
		userService:  userService,
	}
}

// CreateShowcase creates new showcase and awards points
func (s *ShowcaseService) CreateShowcase(ctx context.Context, userID string, title string, description string, mediaURLs []string, categoryID string, visibility string) (*domain.Showcase, error) {
	mediaJSON, _ := json.Marshal(mediaURLs)

	if visibility == "" {
		visibility = "public"
	}

	showcase := &domain.Showcase{
		UserID:      userID,
		Title:       title,
		Description: description,
		MediaURLs:   datatypes.JSON(mediaJSON),
		CategoryID:  categoryID,
		Visibility:  visibility,
		LikesCount:  0,
	}

	if err := s.showcaseRepo.Create(ctx, showcase); err != nil {
		return nil, err
	}

	// Award points for creating showcase
	s.userService.AddPoints(ctx, userID, "create_showcase", 10, map[string]interface{}{
		"showcase_id": showcase.ID,
		"title":       title,
	})

	return showcase, nil
}

// ListShowcases returns paginated list of public showcases
func (s *ShowcaseService) ListShowcases(ctx context.Context, page, limit int, categoryID string, sort string) ([]domain.Showcase, int64, error) {
	var catPtr *string
	if categoryID != "" {
		catPtr = &categoryID
	}
	return s.showcaseRepo.ListAll(ctx, page, limit, catPtr)
}

// GetShowcaseByID returns a showcase by ID with user and category preloaded
func (s *ShowcaseService) GetShowcaseByID(ctx context.Context, showcaseID string) (*domain.Showcase, error) {
	return s.showcaseRepo.FindByID(ctx, showcaseID)
}

// UpdateShowcase updates title, description, visibility — only owner can do it
func (s *ShowcaseService) UpdateShowcase(ctx context.Context, showcaseID string, userID string, title string, description string, visibility string) (*domain.Showcase, error) {
	showcase, err := s.showcaseRepo.FindByID(ctx, showcaseID)
	if err != nil {
		return nil, err
	}
	if showcase.UserID != userID {
		return nil, domain.ErrForbidden
	}
	if title != "" {
		showcase.Title = title
	}
	if description != "" {
		showcase.Description = description
	}
	if visibility != "" {
		showcase.Visibility = visibility
	}
	if err := s.showcaseRepo.Update(ctx, showcase); err != nil {
		return nil, err
	}
	return showcase, nil
}

// DeleteShowcase deletes a showcase — only owner can do it
func (s *ShowcaseService) DeleteShowcase(ctx context.Context, showcaseID string, userID string) error {
	showcase, err := s.showcaseRepo.FindByID(ctx, showcaseID)
	if err != nil {
		return err
	}
	if showcase.UserID != userID {
		return domain.ErrForbidden
	}
	return s.showcaseRepo.Delete(ctx, showcaseID)
}

// LikeShowcase adds like to showcase and awards points to creator
func (s *ShowcaseService) LikeShowcase(ctx context.Context, userID, showcaseID string) error {
	// Check if already liked
	exists, err := s.likeRepo.Exists(ctx, userID, showcaseID)
	if err != nil {
		return err
	}
	if exists {
		return nil // Already liked, silent success
	}

	// Add like
	like := &domain.ShowcaseLike{
		UserID:     userID,
		ShowcaseID: showcaseID,
	}
	if err := s.likeRepo.Create(ctx, like); err != nil {
		return err
	}

	// Increment likes count
	if err := s.showcaseRepo.IncrementLikes(ctx, showcaseID); err != nil {
		return err
	}

	// Award points to showcase creator
	showcase, err := s.showcaseRepo.FindByID(ctx, showcaseID)
	if err == nil && showcase.UserID != userID {
		s.userService.AddPoints(ctx, showcase.UserID, "receive_like", 5, map[string]interface{}{
			"showcase_id": showcaseID,
			"from_user":   userID,
		})
	}

	return nil
}

// UnlikeShowcase removes like from showcase
func (s *ShowcaseService) UnlikeShowcase(ctx context.Context, userID, showcaseID string) error {
	exists, err := s.likeRepo.Exists(ctx, userID, showcaseID)
	if err != nil {
		return err
	}
	if !exists {
		return nil // Not liked, silent success
	}

	if err := s.likeRepo.DeleteByUserShowcase(ctx, userID, showcaseID); err != nil {
		return err
	}

	return s.showcaseRepo.DecrementLikes(ctx, showcaseID)
}

// ============================================================================
// GLOBAL SERVICE INSTANCES
// ============================================================================

var (
	userSvc         IUserService
	showcaseSvc     IShowcaseService
	certificateSvc  ICertificateService
	classSvc        IClassService
	discussionSvc   IDiscussionService
	replySvc        IReplyService
	gamificationSvc IGamificationService
	projectSvc      IProjectService
	categorySvc     ICategoryService
	phaseSvc        IPhaseService
)

// InitServices initializes all global service instances
func InitServices(db *gorm.DB) {
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	pointRepo := repository.NewCommunityPointRepository(db)
	levelRepo := repository.NewUserLevelRepository(db)
	showcaseRepo := repository.NewShowcaseRepository(db)
	likeRepo := repository.NewShowcaseLikeRepository(db)
	certificateRepo := repository.NewCertificateRepository(db)
	classRepo := repository.NewClassRepository(db)
	classDetailRepo := repository.NewClassDetailRepository(db)
	classProgressRepo := repository.NewClassProgressRepository(db)
	discussionRepo := repository.NewDiscussionRepository(db)
	replyRepo := repository.NewReplyRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	phaseRepo := repository.NewPhaseRepository(db)

	// Initialize services with repositories
	userSvc = NewUserService(userRepo, pointRepo, levelRepo)
	showcaseSvc = NewShowcaseService(showcaseRepo, likeRepo, userSvc)
	certificateSvc = NewCertificateService(certificateRepo, classRepo, userRepo)
	classSvc = NewClassService(classRepo, classDetailRepo, classProgressRepo, certificateRepo, userSvc, projectRepo, phaseRepo)
	discussionSvc = NewDiscussionService(discussionRepo, replyRepo, userRepo)
	replySvc = NewReplyService(replyRepo, discussionRepo, userSvc)
	gamificationSvc = NewGamificationService(pointRepo, levelRepo, userRepo, showcaseRepo, userSvc)
	projectSvc = NewProjectService(projectRepo, classRepo, userRepo)
	phaseSvc = NewPhaseService(phaseRepo, projectRepo)

	// Category service
	categoryRepo := repository.NewCategoryRepository(db)
	categorySvc = NewCategoryService(categoryRepo)
}

// Getter functions to access global services
func GetUserService() IUserService {
	return userSvc
}

func GetShowcaseService() IShowcaseService {
	return showcaseSvc
}

func GetCertificateService() ICertificateService {
	return certificateSvc
}

func GetClassService() IClassService {
	return classSvc
}

func GetDiscussionService() IDiscussionService {
	return discussionSvc
}

func GetReplyService() IReplyService {
	return replySvc
}

func GetGamificationService() IGamificationService {
	return gamificationSvc
}

func GetProjectService() IProjectService {
	return projectSvc
}

func GetCategoryService() ICategoryService {
	return categorySvc
}

func GetPhaseService() IPhaseService {
	return phaseSvc
}
