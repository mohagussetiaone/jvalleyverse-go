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
	GetPublicClassByID(ctx context.Context, classID string) (map[string]interface{}, error)
	GetClassBySlug(ctx context.Context, projectID string, slug string, userID string) (map[string]interface{}, error)
	StartClass(ctx context.Context, userID, classID string) (*domain.ClassProgress, error)
	UpdateProgress(ctx context.Context, userID, classID string, percentage int, notes string) (*domain.ClassProgress, error)
	CompleteClass(ctx context.Context, userID, classID string) (map[string]interface{}, error)
	AdminCreateClass(ctx context.Context, adminID string, input domain.Class) (*domain.Class, error)
	AdminUpdateClass(ctx context.Context, adminID, id string, input domain.Class) (*domain.Class, error)
	AdminDeleteClass(ctx context.Context, adminID, id string) error
	AdminCreateClassDetail(ctx context.Context, classID string, about, rules string, tools, media, resources interface{}) (*domain.ClassDetail, error)
	ListClassesByProject(ctx context.Context, projectID string, limit, offset int) ([]domain.Class, int64, error)
	ListClassesByPhase(ctx context.Context, phaseID string) ([]domain.Class, int64, error)
}

type ClassService struct {
	classRepo    *repository.ClassRepository
	detailRepo   *repository.ClassDetailRepository
	progressRepo *repository.ClassProgressRepository
	certRepo     *repository.CertificateRepository
	userService  IUserService
	projectRepo  *repository.ProjectRepository
	phaseRepo    *repository.PhaseRepository
}

func NewClassService(
	classRepo *repository.ClassRepository,
	detailRepo *repository.ClassDetailRepository,
	progressRepo *repository.ClassProgressRepository,
	certRepo *repository.CertificateRepository,
	userService IUserService,
	projectRepo *repository.ProjectRepository,
	phaseRepo *repository.PhaseRepository,
) *ClassService {
	return &ClassService{
		classRepo:    classRepo,
		detailRepo:   detailRepo,
		progressRepo: progressRepo,
		certRepo:     certRepo,
		userService:  userService,
		projectRepo:  projectRepo,
		phaseRepo:    phaseRepo,
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
		"phase":      class.Phase,
		"project":    class.Project,
	}, nil
}

// GetPublicClassByID returns public class detail
func (s *ClassService) GetPublicClassByID(
	ctx context.Context,
	classID string,
) (map[string]interface{}, error) {

	class, err := s.classRepo.FindPublicByID(ctx, classID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"class":      class,
		"details":    class.Details,
		"project":    class.Project,
		"phase":      class.Phase,
		"next_class": class.NextClass,
	}, nil
}

// StartClass initializes user progress for a class
func (s *ClassService) StartClass(ctx context.Context, userID, classID string) (*domain.ClassProgress, error) {
	if _, err := s.classRepo.FindByID(ctx, classID); err != nil {
		return nil, domain.ErrClassNotFound
	}

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
	if percentage < 0 || percentage > 100 {
		return nil, domain.ErrInvalidInput
	}

	progress, err := s.progressRepo.FindByUserAndClass(ctx, userID, classID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	progress.ProgressPercentage = percentage
	progress.Notes = notes
	switch {
	case percentage == 0:
		progress.Status = "started"
	case percentage > 0 && percentage < 100:
		progress.Status = "in_progress"
	case percentage == 100:
		progress.Status = "completed"
		now := time.Now()
		progress.CompletedAt = &now
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
	if err := s.userService.AddPoints(ctx, userID, "complete_class", 50, map[string]interface{}{
		"class_id":  classID,
		"cert_code": code,
	}); err != nil {
		return nil, err
	}

	// Get next class info
	class, err := s.classRepo.FindByID(ctx, classID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message":     "Class completed!",
		"certificate": cert,
		"achievement": map[string]interface{}{
			"type":        "certificate",
			"title":       class.Title,
			"unique_code": cert.UniqueCode,
			"issued_at":   cert.IssuedAt,
		},
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
	project, err := s.projectRepo.FindByID(ctx, input.ProjectID)
	if err != nil {
		return nil, domain.ErrProjectNotFound
	}
	if project.AdminID != adminID {
		return nil, domain.ErrForbidden
	}

	if input.PhaseID == "" {
		return nil, domain.ErrInvalidInput
	}

	exists, err := s.phaseRepo.ExistsInProject(ctx, input.PhaseID, input.ProjectID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrInvalidInput
	}

	input.AdminID = adminID
	if err := s.classRepo.Create(ctx, &input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *ClassService) AdminUpdateClass(ctx context.Context, adminID, id string, input domain.Class) (*domain.Class, error) {
	class, err := s.classRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if class.AdminID != adminID {
		return nil, domain.ErrForbidden
	}

	if input.ProjectID != "" && input.ProjectID != class.ProjectID {
		project, err := s.projectRepo.FindByID(ctx, input.ProjectID)
		if err != nil {
			return nil, domain.ErrProjectNotFound
		}
		if project.AdminID != adminID {
			return nil, domain.ErrForbidden
		}
		class.ProjectID = input.ProjectID
	}

	if input.PhaseID != "" {
		exists, err := s.phaseRepo.ExistsInProject(ctx, input.PhaseID, class.ProjectID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, domain.ErrInvalidInput
		}
		class.PhaseID = input.PhaseID
	}

	if input.Title != "" {
		class.Title = input.Title
	}
	if input.Description != "" {
		class.Description = input.Description
	}
	if input.Slug != "" {
		class.Slug = input.Slug
	}
	if input.Thumbnail != "" {
		class.Thumbnail = input.Thumbnail
	}
	if input.Difficulty != "" {
		class.Difficulty = input.Difficulty
	}
	if input.Duration > 0 {
		class.Duration = input.Duration
	}
	class.OrderIndex = input.OrderIndex
	class.SequenceNum = input.SequenceNum
	class.IsFirst = input.IsFirst
	class.NextClassID = input.NextClassID
	if input.Visibility != "" {
		class.Visibility = input.Visibility
	}

	if err := s.classRepo.Update(ctx, class); err != nil {
		return nil, err
	}

	updated, err := s.classRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *ClassService) AdminDeleteClass(ctx context.Context, adminID, id string) error {
	class, err := s.classRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if class.AdminID != adminID {
		return domain.ErrForbidden
	}

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

func (s *ClassService) ListClassesByProject(ctx context.Context, projectID string, limit, offset int) ([]domain.Class, int64, error) {
	return s.classRepo.ListByProjectID(ctx, projectID, limit, offset)
}

func (s *ClassService) ListClassesByPhase(ctx context.Context, phaseID string) ([]domain.Class, int64, error) {
	return s.classRepo.ListByPhaseID(ctx, phaseID)
}
