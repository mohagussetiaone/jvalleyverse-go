package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/dto"
	"jvalleyverse/internal/repository"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// IUserService defines the business logic for user management and gamification
type IUserService interface {
	GetProfile(ctx context.Context, userID string) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID string, name, bio, avatar string) error
	ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error
	AddPoints(ctx context.Context, userID string, category string, points int, metadata map[string]interface{}) error
	GetUserActivityLog(ctx context.Context, userID string, page, limit int) ([]dto.ActivityItem, int64, error)
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	ListAllUsers(ctx context.Context, page, limit int) ([]dto.UserListItem, int64, error)
	ListMentors(ctx context.Context, page, limit int) ([]dto.MentorItem, int64, error)
	GenerateRefreshToken(ctx context.Context, userID string) (*domain.RefreshToken, error)
	ValidateRefreshToken(ctx context.Context, token string) (*domain.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	RevokeAllUserTokens(ctx context.Context, userID string) error
}

type UserService struct {
	userRepo    *repository.UserRepository
	pointRepo   *repository.CommunityPointRepository
	levelRepo   *repository.UserLevelRepository
	refreshRepo *repository.RefreshTokenRepository
}

func NewUserService(userRepo *repository.UserRepository, pointRepo *repository.CommunityPointRepository, levelRepo *repository.UserLevelRepository, refreshRepo *repository.RefreshTokenRepository) *UserService {
	return &UserService{
		userRepo:    userRepo,
		pointRepo:   pointRepo,
		levelRepo:   levelRepo,
		refreshRepo: refreshRepo,
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

	if name != "" {
		user.Name = name
	}
	if bio != "" {
		user.Bio = bio
	}
	if avatar != "" {
		user.Avatar = avatar
	}

	return s.userRepo.Update(ctx, user)
}

// ChangePassword changes user password (validates current password first)
func (s *UserService) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	if newPassword == "" {
		return domain.ErrInvalidInput
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// User registered via Google may have empty password
	if user.Password != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
			return domain.ErrUnauthorized
		}
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(ctx, userID, string(hashed))
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

	newLevel := CalculateLevel(user.TotalPoints)
	oldLevel := user.Level

	// Only update if level changed
	if oldLevel != newLevel {
		if err := s.userRepo.UpdateLevel(ctx, userID, newLevel); err != nil {
			return err
		}
		// Notify user about level up with badge info
		if notifSvc := GetNotificationService(); notifSvc != nil {
			// Try to get badge from user_levels table
			badgeName := ""
			badgeIcon := ""
			levelDef, err := s.levelRepo.FindByLevel(ctx, newLevel)
			if err == nil && levelDef != nil {
				badgeName = levelDef.BadgeName
				badgeIcon = levelDef.BadgeIcon
			}
			// Fallback to hardcoded names
			if badgeName == "" {
				levelNames := map[int]string{1: "Beginner", 2: "Intermediate", 3: "Advanced", 4: "Expert", 5: "Master"}
				badgeName = levelNames[newLevel]
				if badgeName == "" {
					badgeName = fmt.Sprintf("Level %d", newLevel)
				}
			}
			if badgeIcon == "" {
				levelIcons := map[int]string{1: "🌱", 2: "🌿", 3: "🌳", 4: "⭐", 5: "👑"}
				badgeIcon = levelIcons[newLevel]
			}

			message := fmt.Sprintf("Selamat! Level Anda naik ke %s %s — Badge: %s", badgeName, badgeIcon, badgeName)
			if badgeIcon == "" {
				message = fmt.Sprintf("Selamat! Level Anda naik ke %s!", badgeName)
			}

			notifSvc.CreateNotification(ctx, userID, "level_up",
				"Level Naik! 🏆",
				message,
				"/users/"+userID+"/points",
			)
		}
	}

	return nil
}

// CalculateLevel determines user level based on total points
func CalculateLevel(totalPoints int) int {
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
func (s *UserService) GetUserActivityLog(ctx context.Context, userID string, page, limit int) ([]dto.ActivityItem, int64, error) {
	points, total, err := s.pointRepo.ListByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.ActivityItem, len(points))
	for i, p := range points {
		result[i] = dto.ActivityItem{
			ID:        p.ID,
			Activity:  p.ActivityType,
			Points:    p.PointsEarned,
			Timestamp: p.CreatedAt,
		}
	}
	return result, total, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, user *domain.User) error {
	return s.userRepo.Create(ctx, user)
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}

// GenerateRefreshToken creates a new refresh token for a user
func (s *UserService) GenerateRefreshToken(ctx context.Context, userID string) (*domain.RefreshToken, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	token := hex.EncodeToString(bytes)

	rt := &domain.RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	if err := s.refreshRepo.Create(ctx, rt); err != nil {
		return nil, err
	}
	return rt, nil
}

// ValidateRefreshToken validates and returns a refresh token
func (s *UserService) ValidateRefreshToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	return s.refreshRepo.FindByToken(ctx, token)
}

// RevokeRefreshToken revokes a specific refresh token
func (s *UserService) RevokeRefreshToken(ctx context.Context, token string) error {
	return s.refreshRepo.RevokeByToken(ctx, token)
}

// RevokeAllUserTokens revokes all refresh tokens for a user
func (s *UserService) RevokeAllUserTokens(ctx context.Context, userID string) error {
	return s.refreshRepo.RevokeByUserID(ctx, userID)
}

// ListMentors returns paginated list of users with role "mentor"
func (s *UserService) ListMentors(ctx context.Context, page, limit int) ([]dto.MentorItem, int64, error) {
	users, total, err := s.userRepo.ListByRole(ctx, "mentor", page, limit)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.MentorItem, len(users))
	for i, user := range users {
		result[i] = dto.MentorItem{
			ID:          user.ID,
			Name:        user.Name,
			Avatar:      user.Avatar,
			Bio:         user.Bio,
			Level:       user.Level,
			TotalPoints: user.TotalPoints,
		}
	}

	return result, total, nil
}

// ListAllUsers returns paginated list of all users (admin only)
func (s *UserService) ListAllUsers(ctx context.Context, page, limit int) ([]dto.UserListItem, int64, error) {
	users, total, err := s.userRepo.ListAll(ctx, page, limit)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.UserListItem, len(users))
	for i, user := range users {
		result[i] = dto.UserListItem{
			ID:          user.ID,
			Email:       user.Email,
			Name:        user.Name,
			Avatar:      user.Avatar,
			Role:        user.Role,
			Level:       user.Level,
			TotalPoints: user.TotalPoints,
			IsActive:    user.IsActive,
			CreatedAt:   user.CreatedAt,
		}
	}

	return result, total, nil
}

// IShowcaseService defines the business logic for managing user showcases
type IShowcaseService interface {
	CreateShowcase(ctx context.Context, userID string, title string, description string, mediaURLs []string, categoryID string, visibility string) (*domain.Showcase, error)
	ListShowcases(ctx context.Context, page, limit int, categoryID string, sort string) ([]dto.ShowcaseListItem, int64, error)
	ListMyShowcases(ctx context.Context, userID string, page, limit int) ([]dto.ShowcaseListItem, int64, error)
	GetShowcaseByID(ctx context.Context, showcaseID string) (*dto.ShowcaseDetail, error)
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

	// Notify creator as activity history
	if notifSvc := GetNotificationService(); notifSvc != nil {
		notifSvc.CreateNotification(ctx, userID, "showcase_created",
			"Showcase Baru Dibuat",
			"Anda membuat showcase: "+title,
			"/showcases/"+showcase.ID,
		)
	}

	return showcase, nil
}

// ListShowcases returns paginated list of public showcases
func (s *ShowcaseService) ListShowcases(ctx context.Context, page, limit int, categoryID string, sort string) ([]dto.ShowcaseListItem, int64, error) {
	var catPtr *string
	if categoryID != "" {
		catPtr = &categoryID
	}
	showcases, total, err := s.showcaseRepo.ListAll(ctx, page, limit, catPtr)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.ShowcaseListItem, len(showcases))
	for i, sc := range showcases {
		result[i] = dto.ToShowcaseListItem(sc)
	}
	return result, total, nil
}

// GetShowcaseByID returns a showcase by ID with user and category preloaded
func (s *ShowcaseService) GetShowcaseByID(ctx context.Context, showcaseID string) (*dto.ShowcaseDetail, error) {
	sc, err := s.showcaseRepo.FindByID(ctx, showcaseID)
	if err != nil {
		return nil, err
	}
	return dto.ToShowcaseDetail(sc), nil
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

// ListMyShowcases returns paginated showcases by user
func (s *ShowcaseService) ListMyShowcases(ctx context.Context, userID string, page, limit int) ([]dto.ShowcaseListItem, int64, error) {
	showcases, total, err := s.showcaseRepo.ListByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.ShowcaseListItem, len(showcases))
	for i, sc := range showcases {
		result[i] = dto.ToShowcaseListItem(sc)
	}

	return result, total, nil
}

// LikeShowcase adds like to showcase, awards points, and creates notification
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

	// Award points to showcase creator + send notification
	showcase, err := s.showcaseRepo.FindByID(ctx, showcaseID)
	if err == nil && showcase.UserID != userID {
		s.userService.AddPoints(ctx, showcase.UserID, "receive_like", 5, map[string]interface{}{
			"showcase_id": showcaseID,
			"from_user":   userID,
		})
		// Create notification for showcase owner
		if notifSvc != nil {
			notifSvc.CreateNotification(ctx, showcase.UserID, "showcase_like",
				"Showcase Anda Mendapat Like",
				"Seseorang menyukai showcase Anda: "+showcase.Title,
				"/showcases/"+showcaseID,
			)
		}
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
	lessonSvc       ILessonService
	blogSvc         IBlogService
	discussionSvc   IDiscussionService
	replySvc        IReplyService
	gamificationSvc IGamificationService
	studyCaseSvc    IStudyCaseService
	courseSvc       ICourseService
	categorySvc     ICategoryService
	sectionSvc      ISectionService
	reviewSvc       IReviewService
	faqSvc          IFaqService
	companySvc      ICompanyService
	auditSvc        *AuditService
	notifSvc        INotificationService
	dashboardSvc    IDashboardService
	streakSvc       *StreakService
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
	lessonRepo := repository.NewLessonRepository(db)
	classDetailRepo := repository.NewLessonDetailRepository(db)
	classProgressRepo := repository.NewLessonProgressRepository(db)
	discussionRepo := repository.NewDiscussionRepository(db)
	replyRepo := repository.NewReplyRepository(db)
	reactRepo := repository.NewReplyReactionRepository(db)
	replyLikeRepo := repository.NewReplyLikeRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	sectionRepo := repository.NewSectionRepository(db)
	enrollRepo := repository.NewEnrollmentRepository(db)
	blogRepo := repository.NewBlogRepository(db)

	// Initialize services with repositories
	refreshRepo := repository.NewRefreshTokenRepository(db)
	userSvc = NewUserService(userRepo, pointRepo, levelRepo, refreshRepo)
	showcaseSvc = NewShowcaseService(showcaseRepo, likeRepo, userSvc)
	certificateSvc = NewCertificateService(certificateRepo, lessonRepo, userRepo)

	// Learning Streak repository (needed by LessonService and DashboardService)
	streakRepo := repository.NewLearningStreakRepository(db)

	lessonSvc = NewLessonService(lessonRepo, classDetailRepo, classProgressRepo, certificateRepo, userSvc, courseRepo, enrollRepo, sectionRepo, streakRepo)
	discussionSvc = NewDiscussionService(discussionRepo, replyRepo, userRepo)
	replySvc = NewReplyService(replyRepo, reactRepo, replyLikeRepo, discussionRepo, userSvc)
	gamificationSvc = NewGamificationService(pointRepo, levelRepo, userRepo, showcaseRepo, userSvc)
	courseSvc = NewCourseService(courseRepo, lessonRepo, userRepo, enrollRepo)
	sectionSvc = NewSectionService(sectionRepo, courseRepo, enrollRepo) // StudyCase service
	studyCaseRepo := repository.NewStudyCaseRepository(db)
	studyCaseSvc = NewStudyCaseService(studyCaseRepo)

	blogSvc = NewBlogService(blogRepo)
	// Review service
	reviewRepo := repository.NewReviewRepository(db)
	reviewSvc = NewReviewService(reviewRepo, courseRepo)

	// Category service
	categoryRepo := repository.NewCategoryRepository(db)
	categorySvc = NewCategoryService(categoryRepo, enrollRepo)

	// Audit service
	auditSvc = NewAuditService(repository.NewAdminAuditLogRepository(db))

	// Notification service
	notifRepo := repository.NewNotificationRepository(db)
	notifSvc = NewNotificationService(notifRepo)

	// Streak service
	streakSvc = NewStreakService(streakRepo)

	// Dashboard service
	dashboardSvc = NewDashboardService(lessonRepo, classProgressRepo, notifSvc, streakRepo)

	// FAQ service
	faqRepo := repository.NewFAQRepository(db)
	faqSvc = NewFaqService(faqRepo)

	// Company service
	companyRepo := repository.NewCompanyRepository(db)
	companySvc = NewCompanyService(companyRepo)
}

func GetUserService() IUserService {
	return userSvc
}

func GetBlogService() IBlogService {
	return blogSvc
}

func GetShowcaseService() IShowcaseService {
	return showcaseSvc
}

func GetCertificateService() ICertificateService {
	return certificateSvc
}

func GetLessonService() ILessonService {
	return lessonSvc
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

func GetCourseService() ICourseService {
	return courseSvc
}

func GetCategoryService() ICategoryService {
	return categorySvc
}

func GetSectionService() ISectionService {
	return sectionSvc
}

func GetStudyCaseService() IStudyCaseService {
	return studyCaseSvc
}

func GetReviewService() IReviewService {
	return reviewSvc
}

func GetAuditService() *AuditService {
	return auditSvc
}

func GetNotificationService() INotificationService {
	return notifSvc
}

func GetDashboardService() IDashboardService {
	return dashboardSvc
}

func GetStreakService() *StreakService {
	return streakSvc
}

func GetFaqService() IFaqService {
	return faqSvc
}

func GetCompanyService() ICompanyService {
	return companySvc
}
