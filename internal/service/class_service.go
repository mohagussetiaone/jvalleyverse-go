package service

import (
	"encoding/json"
	"errors"
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/repository"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ClassService struct {
	classRepo    *repository.ClassRepository
	detailRepo   *repository.ClassDetailRepository
	progressRepo *repository.ClassProgressRepository
	certRepo     *repository.CertificateRepository
	userService  *UserService
	projectRepo  *repository.ProjectRepository
}

func NewClassService() *ClassService {
	return &ClassService{
		classRepo:    repository.NewClassRepository(),
		detailRepo:   repository.NewClassDetailRepository(),
		progressRepo: repository.NewClassProgressRepository(),
		certRepo:     repository.NewCertificateRepository(),
		userService:  NewUserService(),
		projectRepo:  repository.NewProjectRepository(),
	}
}

// GetClassBySlug returns class with its details and user progress
func (s *ClassService) GetClassBySlug(projectID uint, slug string, userID uint) (map[string]interface{}, error) {
	class, err := s.classRepo.FindBySlug(projectID, slug)
	if err != nil {
		return nil, err
	}

	// Fetch user progress
	progress, err := s.progressRepo.FindByUserAndClass(userID, class.ID)
	if err != nil {
		// Not started yet
		progress = &domain.ClassProgress{
			Status:             "not_started",
			ProgressPercentage: 0,
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
func (s *ClassService) StartClass(userID, classID uint) (*domain.ClassProgress, error) {
	progress, err := s.progressRepo.FindByUserAndClass(userID, classID)
	if err == nil {
		// Already started or exists
		if progress.Status == "not_started" {
			progress.Status = "started"
			now := time.Now()
			progress.StartedAt = &now
			err = s.progressRepo.Update(progress)
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
	err = s.progressRepo.Create(progress)
	return progress, err
}

// UpdateProgress updates user progress percentage
func (s *ClassService) UpdateProgress(userID, classID uint, percentage int, notes string) (*domain.ClassProgress, error) {
	progress, err := s.progressRepo.FindByUserAndClass(userID, classID)
	if err != nil {
		return nil, errors.New("class not started")
	}

	progress.ProgressPercentage = percentage
	progress.Notes = notes
	if percentage > 0 && percentage < 100 {
		progress.Status = "in_progress"
	}

	err = s.progressRepo.Update(progress)
	return progress, err
}

// CompleteClass marks class as completed and issues certificate
func (s *ClassService) CompleteClass(userID, classID uint) (map[string]interface{}, error) {
	progress, err := s.progressRepo.FindByUserAndClass(userID, classID)
	if err != nil {
		return nil, errors.New("class not started")
	}

	if progress.Status == "completed" {
		return nil, errors.New("class already completed")
	}

	now := time.Now()
	progress.Status = "completed"
	progress.ProgressPercentage = 100
	progress.CompletedAt = &now

	if err := s.progressRepo.Update(progress); err != nil {
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
	if err := s.certRepo.Create(cert); err != nil {
		return nil, err
	}

	// Award points (50 for completion as per previous logic)
	s.userService.AddPoints(userID, "complete_class", 50, map[string]interface{}{
		"class_id":  classID,
		"cert_code": code,
	})

	// Get next class info
	class, _ := s.classRepo.FindByID(classID)

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

func (s *ClassService) AdminCreateClass(adminID uint, input domain.Class) (*domain.Class, error) {
	input.AdminID = adminID
	if err := s.classRepo.Create(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *ClassService) AdminUpdateClass(id uint, input domain.Class) error {
	class, err := s.classRepo.FindByID(id)
	if err != nil {
		return err
	}
	// Update fields
	class.Title = input.Title
	class.Description = input.Description
	class.Slug = input.Slug
	class.Difficulty = input.Difficulty
	class.Duration = input.Duration
	class.OrderIndex = input.OrderIndex
	class.SequenceNum = input.SequenceNum
	class.NextClassID = input.NextClassID

	return s.classRepo.Update(class)
}

func (s *ClassService) AdminDeleteClass(id uint) error {
	return s.classRepo.DeleteByID(id)
}

func (s *ClassService) AdminCreateClassDetail(classID uint, about, rules string, tools, media, resources interface{}) (*domain.ClassDetail, error) {
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

	if err := s.detailRepo.Create(detail); err != nil {
		return nil, err
	}
	return detail, nil
}
