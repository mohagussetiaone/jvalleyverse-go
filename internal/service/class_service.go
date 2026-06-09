package service

import (
	"context"
	"encoding/json"
	"errors"
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// IClassService defines the business logic for managing learning classes and progress
type IClassService interface {
	GetClassBySlug(ctx context.Context, projectID string, slug string, userID string) (map[string]interface{}, error)
	StartClass(ctx context.Context, userID, classID string) (*domain.ClassProgress, error)
	UpdateProgress(ctx context.Context, userID, classID string, percentage int, notes string) (*domain.ClassProgress, error)
	CompleteClass(ctx context.Context, userID, classID string) (map[string]interface{}, error)
	AdminCreateClass(ctx context.Context, adminID string, input domain.Class) (*domain.Class, error)
	AdminUpdateClass(ctx context.Context, id string, input domain.Class) error
	AdminDeleteClass(ctx context.Context, id string) error
	AdminCreateClassDetail(ctx context.Context, classID string, about, rules string, tools, media, resources interface{}) (*domain.ClassDetail, error)
}

type ClassService struct {
	classRepo    *repository.ClassRepository
	detailRepo   *repository.ClassDetailRepository
	progressRepo *repository.ClassProgressRepository
	certRepo     *repository.CertificateRepository
	userService  IUserService
	projectRepo  *repository.ProjectRepository
}

func NewClassService(
	classRepo *repository.ClassRepository,
	detailRepo *repository.ClassDetailRepository,
	progressRepo *repository.ClassProgressRepository,
	certRepo *repository.CertificateRepository,
	userService IUserService,
	projectRepo *repository.ProjectRepository,
) *ClassService {
	return &ClassService{
		classRepo:    classRepo,
		detailRepo:   detailRepo,
		progressRepo: progressRepo,
		certRepo:     certRepo,
		userService:  userService,
		projectRepo:  projectRepo,
	}
}

// GetClassBySlug returns class with its details and user progress
func (s *ClassService) GetClassBySlug(ctx context.Context, projectID string, slug string, userID string) (map[string]interface{}, error) {
	class, err := s.classRepo.FindBySlug(ctx, projectID, slug)
	if err != nil {
		return nil, err
	}

	// Fetch user progress (only when logged in)
	var progress *domain.ClassProgress
	if userID != "" {
		progress, err = s.progressRepo.FindByUserAndClass(ctx, userID, class.ID)
		if err != nil {
			// Not started yet
			progress = &domain.ClassProgress{
				Status:             "not_started",
				ProgressPercentage: 0,
			}
		}
	}

	return map[string]interface{}{
		"class":      class,
		"details":    class.Details,
		"progress":   progress,
		"next_class": class.NextClass,
	}, nil
}

// StartClass initializes user progress for a class
func (s *ClassService) StartClass(ctx context.Context, userID, classID string) (*domain.ClassProgress, error) {
	progress, err := s.progressRepo.FindByUserAndClass(ctx, userID, classID)
	if err == nil {
		// Already started or exists
		if progress.Status == "not_started" {
			progress.Status = "started"
			now := time.Now()
			progress.StartedAt = &now
			err = s.progressRepo.Update(ctx, progress)
			return progress, err
		}
		return progress, nil
	}

	// Create new progress record
	now := time.Now()
	progress = &domain.ClassProgress{
		UserID:             userID,
		ClassID:            classID,
		Status:             "started",
		StartedAt:          &now,
		ProgressPercentage: 0,
	}
	err = s.progressRepo.Create(ctx, progress)
	return progress, err
}

// UpdateProgress updates user progress percentage
func (s *ClassService) UpdateProgress(ctx context.Context, userID, classID string, percentage int, notes string) (*domain.ClassProgress, error) {
	progress, err := s.progressRepo.FindByUserAndClass(ctx, userID, classID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	progress.ProgressPercentage = percentage
	progress.Notes = notes
	if percentage > 0 && percentage < 100 {
		progress.Status = "in_progress"
	}

	err = s.progressRepo.Update(ctx, progress)
	return progress, err
}

// CompleteClass marks class as completed and issues certificate
func (s *ClassService) CompleteClass(ctx context.Context, userID, classID string) (map[string]interface{}, error) {
	progress, err := s.progressRepo.FindByUserAndClass(ctx, userID, classID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	if progress.Status == "completed" {
		return nil, errors.New("class already completed")
	}

	now := time.Now()
	progress.Status = "completed"
	progress.ProgressPercentage = 100
	progress.CompletedAt = &now

	if err := s.progressRepo.Update(ctx, progress); err != nil {
		return nil, err
	}

	// Issue certificate
	code := generateUniqueCertificateCode()
	cert := &domain.Certificate{
		UserID:     userID,
		ClassID:    classID,
		UniqueCode: code,
		IssuedAt:   now,
	}
	if err := s.certRepo.Create(ctx, cert); err != nil {
		return nil, err
	}

	// Award points (50 for completion)
	s.userService.AddPoints(ctx, userID, "complete_class", 50, map[string]interface{}{
		"class_id":  classID,
		"cert_code": code,
	})

	// Get next class info
	class, _ := s.classRepo.FindByID(ctx, classID)

	return map[string]interface{}{
		"message":        "Class completed!",
		"certificate":    cert,
		"progress":       progress,
		"next_class":     class.NextClass,
		"points_awarded": 50,
	}, nil
}

// generateUniqueCertificateCode generates unique certificate code
func generateUniqueCertificateCode() string {
	return "CERT-" + uuid.New().String()[:8]
}

// Admin Methods

func (s *ClassService) AdminCreateClass(ctx context.Context, adminID string, input domain.Class) (*domain.Class, error) {
	input.AdminID = adminID
	if err := s.classRepo.Create(ctx, &input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *ClassService) AdminUpdateClass(ctx context.Context, id string, input domain.Class) error {
	class, err := s.classRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	class.Title = input.Title
	class.Description = input.Description
	class.Slug = input.Slug
	class.Difficulty = input.Difficulty
	class.Duration = input.Duration
	class.OrderIndex = input.OrderIndex
	class.SequenceNum = input.SequenceNum
	class.NextClassID = input.NextClassID

	return s.classRepo.Update(ctx, class)
}

func (s *ClassService) AdminDeleteClass(ctx context.Context, id string) error {
	return s.classRepo.DeleteByID(ctx, id)
}

func (s *ClassService) AdminCreateClassDetail(ctx context.Context, classID string, about, rules string, tools, media, resources interface{}) (*domain.ClassDetail, error) {
	toolsJSON, _ := json.Marshal(tools)
	mediaJSON, _ := json.Marshal(media)
	resourcesJSON, _ := json.Marshal(resources)

	detail := &domain.ClassDetail{
		ClassID:       classID,
		About:         about,
		Rules:         rules,
		Tools:         datatypes.JSON(toolsJSON),
		ResourceMedia: datatypes.JSON(mediaJSON),
		Resources:     datatypes.JSON(resourcesJSON),
	}

	if err := s.detailRepo.Create(ctx, detail); err != nil {
		return nil, err
	}
	return detail, nil
}
