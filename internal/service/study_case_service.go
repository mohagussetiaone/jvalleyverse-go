package service

import (
	"context"
	"encoding/json"
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/dto"
	"jvalleyverse/internal/repository"

	"gorm.io/datatypes"
)

type StudyCaseListFilter struct {
	CategoryID *string
}

type IStudyCaseService interface {
	CreateStudyCase(ctx context.Context, userID, name, description, imgURL, youtubeURL string, categoryID *string, tags []string) (*domain.StudyCase, error)
	GetStudyCaseByID(ctx context.Context, id string) (*dto.StudyCaseDetail, error)
	ListStudyCases(ctx context.Context, page, limit int, filter *StudyCaseListFilter) ([]dto.StudyCaseListItem, int64, error)
	ListStudyCasesByUser(ctx context.Context, userID string, page, limit int) ([]dto.StudyCaseListItem, int64, error)
	UpdateStudyCase(ctx context.Context, id string, name, description, imgURL, youtubeURL string, categoryID *string, tags []string) (*domain.StudyCase, error)
	DeleteStudyCase(ctx context.Context, id string) error
}

type StudyCaseService struct {
	studyCaseRepo *repository.StudyCaseRepository
}

func NewStudyCaseService(studyCaseRepo *repository.StudyCaseRepository) *StudyCaseService {
	return &StudyCaseService{studyCaseRepo: studyCaseRepo}
}

func toStudyCaseFilter(f *StudyCaseListFilter) *repository.StudyCaseListFilter {
	if f == nil {
		return nil
	}
	return &repository.StudyCaseListFilter{
		CategoryID: f.CategoryID,
	}
}

func (s *StudyCaseService) CreateStudyCase(ctx context.Context, userID, name, description, imgURL, youtubeURL string, categoryID *string, tags []string) (*domain.StudyCase, error) {
	if name == "" {
		return nil, domain.ErrInvalidInput
	}

	tagsJSON, _ := json.Marshal(tags)

	studyCase := &domain.StudyCase{
		Name:        name,
		Description: description,
		ImgURL:      imgURL,
		YoutubeURL:  youtubeURL,
		CategoryID:  categoryID,
		Tags:        datatypes.JSON(tagsJSON),
		UserID:      userID,
	}

	if err := s.studyCaseRepo.Create(ctx, studyCase); err != nil {
		return nil, err
	}

	// Notify creator as activity history
	if notifSvc := GetNotificationService(); notifSvc != nil {
		notifSvc.CreateNotification(ctx, userID, "study_case_created",
			"Study Case Baru Dibuat",
			"Anda membuat study case: "+name,
			"/study-case/"+studyCase.ID,
		)
	}

	return studyCase, nil
}

func (s *StudyCaseService) GetStudyCaseByID(ctx context.Context, id string) (*dto.StudyCaseDetail, error) {
	studyCase, err := s.studyCaseRepo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.ErrStudyCaseNotFound
	}

	var tags []string
	if studyCase.Tags != nil {
		json.Unmarshal(studyCase.Tags, &tags)
	}

	var category dto.CategoryBrief
	if studyCase.CategoryID != nil {
		category = dto.ToCategoryBrief(studyCase.Category)
	}

	discussions := make([]dto.DiscussionBrief, len(studyCase.Discussions))
	for i, d := range studyCase.Discussions {
		discussions[i] = dto.ToDiscussionBrief(d)
	}

	return &dto.StudyCaseDetail{
		ID:          studyCase.ID,
		Name:        studyCase.Name,
		Description: studyCase.Description,
		ImgURL:      studyCase.ImgURL,
		YoutubeURL:  studyCase.YoutubeURL,
		Tags:        tags,
		Category:    category,
		User:        dto.ToUserBrief(studyCase.User),
		CreatedAt:   studyCase.CreatedAt,
		Discussions: discussions,
	}, nil
}

func (s *StudyCaseService) ListStudyCasesByUser(ctx context.Context, userID string, page, limit int) ([]dto.StudyCaseListItem, int64, error) {
	studyCases, total, err := s.studyCaseRepo.ListByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.StudyCaseListItem, len(studyCases))
	for i, sc := range studyCases {
		result[i] = dto.ToStudyCaseListItem(sc)
	}

	return result, total, nil
}

func (s *StudyCaseService) ListStudyCases(ctx context.Context, page, limit int, filter *StudyCaseListFilter) ([]dto.StudyCaseListItem, int64, error) {
	studyCases, total, err := s.studyCaseRepo.ListAll(ctx, page, limit, toStudyCaseFilter(filter))
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.StudyCaseListItem, len(studyCases))
	for i, sc := range studyCases {
		result[i] = dto.ToStudyCaseListItem(sc)
	}

	return result, total, nil
}

func (s *StudyCaseService) UpdateStudyCase(ctx context.Context, id string, name, description, imgURL, youtubeURL string, categoryID *string, tags []string) (*domain.StudyCase, error) {
	studyCase, err := s.studyCaseRepo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.ErrStudyCaseNotFound
	}

	if name != "" {
		studyCase.Name = name
	}
	if description != "" {
		studyCase.Description = description
	}
	if imgURL != "" {
		studyCase.ImgURL = imgURL
	}
	if youtubeURL != "" {
		studyCase.YoutubeURL = youtubeURL
	}
	if categoryID != nil {
		studyCase.CategoryID = categoryID
	}
	if len(tags) > 0 {
		tagsJSON, _ := json.Marshal(tags)
		studyCase.Tags = datatypes.JSON(tagsJSON)
	}

	if err := s.studyCaseRepo.Update(ctx, studyCase); err != nil {
		return nil, err
	}
	return studyCase, nil
}

func (s *StudyCaseService) DeleteStudyCase(ctx context.Context, id string) error {
	_, err := s.studyCaseRepo.FindByID(ctx, id)
	if err != nil {
		return domain.ErrStudyCaseNotFound
	}
	return s.studyCaseRepo.DeleteByID(ctx, id)
}
