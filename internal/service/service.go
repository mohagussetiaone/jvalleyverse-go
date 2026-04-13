package service

import (
	"encoding/json"
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/repository"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Initialize all repositories after database connection
func InitServices(db *gorm.DB) {
	repository.InitRepository(db)
}

type UserService struct {
	userRepo  *repository.UserRepository
	pointRepo *repository.CommunityPointRepository
	levelRepo *repository.UserLevelRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo:  repository.NewUserRepository(),
		pointRepo: repository.NewCommunityPointRepository(),
		levelRepo: repository.NewUserLevelRepository(),
	}
}

// AddPoints adds points to user and updates level if threshold reached
func (s *UserService) AddPoints(userID uint, category string, points int, metadata map[string]interface{}) error {
	// Add points to user
	if err := s.userRepo.UpdatePoints(userID, points); err != nil {
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
	if err := s.pointRepo.Create(point); err != nil {
		return err
	}

	// Recalculate level
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	newLevel := calculateLevel(user.TotalPoints)
	oldLevel := user.Level

	// Only update if level changed
	if oldLevel != newLevel {
		if err := s.userRepo.UpdateLevel(userID, newLevel); err != nil {
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
func (s *UserService) GetUserActivityLog(userID uint, page, limit int) ([]map[string]interface{}, error) {
	points, _, err := s.pointRepo.ListByUserID(userID, page, limit)
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

type ShowcaseService struct {
	showcaseRepo *repository.ShowcaseRepository
	likeRepo     *repository.ShowcaseLikeRepository
	userService  *UserService
}

func NewShowcaseService() *ShowcaseService {
	return &ShowcaseService{
		showcaseRepo: repository.NewShowcaseRepository(),
		likeRepo:     repository.NewShowcaseLikeRepository(),
		userService:  NewUserService(),
	}
}

// CreateShowcase creates new showcase and awards points
func (s *ShowcaseService) CreateShowcase(userID uint, title string, mediaURLs []string, categoryID uint) (*domain.Showcase, error) {
	mediaJSON, _ := json.Marshal(mediaURLs)

	showcase := &domain.Showcase{
		UserID:     userID,
		Title:      title,
		MediaURLs:  datatypes.JSON(mediaJSON),
		CategoryID: categoryID,
		LikesCount: 0,
	}

	if err := s.showcaseRepo.Create(showcase); err != nil {
		return nil, err
	}

	// Award points for creating showcase
	s.userService.AddPoints(userID, "create_showcase", 10, map[string]interface{}{
		"showcase_id": showcase.ID,
		"title":       title,
	})

	return showcase, nil
}

// LikeShowcase adds like to showcase and awards points to creator
func (s *ShowcaseService) LikeShowcase(userID, showcaseID uint) error {
	// Check if already liked
	exists, err := s.likeRepo.Exists(userID, showcaseID)
	if err != nil {
		return err
	}
	if exists {
		return nil // Already liked, silent success
	}

	// Add like
	like := &domain.ShowcaseLike{
		UserID:      userID,
		ShowcaseID:  showcaseID,
	}
	if err := s.likeRepo.Create(like); err != nil {
		return err
	}

	// Increment likes count
	if err := s.showcaseRepo.IncrementLikes(showcaseID); err != nil {
		return err
	}

	// Award points to showcase creator
	showcase, err := s.showcaseRepo.FindByID(showcaseID)
	if err == nil && showcase.UserID != userID {
		s.userService.AddPoints(showcase.UserID, "receive_like", 5, map[string]interface{}{
			"showcase_id": showcaseID,
			"from_user":   userID,
		})
	}

	return nil
}

// UnlikeShowcase removes like from showcase
func (s *ShowcaseService) UnlikeShowcase(userID, showcaseID uint) error {
	exists, err := s.likeRepo.Exists(userID, showcaseID)
	if err != nil {
		return err
	}
	if !exists {
		return nil // Not liked, silent success
	}

	if err := s.likeRepo.DeleteByUserShowcase(userID, showcaseID); err != nil {
		return err
	}

	return s.showcaseRepo.DecrementLikes(showcaseID)
}