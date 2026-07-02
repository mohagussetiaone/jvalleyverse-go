package service

import (
	"context"
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/dto"
	"jvalleyverse/internal/repository"
)

type ISectionService interface {
	CreateSection(ctx context.Context, adminID, courseID string, title, description string, orderIndex int) (*domain.Section, error)
	UpdateSection(ctx context.Context, adminID, sectionID string, title, description string, orderIndex int) (*domain.Section, error)
	DeleteSection(ctx context.Context, adminID, sectionID string) error
	GetSection(ctx context.Context, sectionID string) (*dto.SectionDetail, error)
	ListSectionsByCourse(ctx context.Context, courseID string) ([]dto.SectionDetail, error)
	GetCourseWithSections(ctx context.Context, courseID string, userID string) (*dto.CourseDetailWithSections, error)
}

type SectionService struct {
	sectionRepo *repository.SectionRepository
	courseRepo  *repository.CourseRepository
	enrollRepo  *repository.EnrollmentRepository
}

func NewSectionService(
	sectionRepo *repository.SectionRepository,
	courseRepo *repository.CourseRepository,
	enrollRepo *repository.EnrollmentRepository,
) *SectionService {
	return &SectionService{
		sectionRepo: sectionRepo,
		courseRepo:  courseRepo,
		enrollRepo: enrollRepo,
	}
}

func (s *SectionService) CreateSection(ctx context.Context, adminID, courseID string, title, description string, orderIndex int) (*domain.Section, error) {
	course, err := s.courseRepo.FindByID(ctx, courseID)
	if err != nil {
		return nil, domain.ErrCourseNotFound
	}
	if course.AdminID != adminID {
		return nil, domain.ErrForbidden
	}

	section := &domain.Section{
		CourseID:    courseID,
		Title:       title,
		Description: description,
		OrderIndex:  orderIndex,
	}

	if err := s.sectionRepo.Create(ctx, section); err != nil {
		return nil, err
	}

	return section, nil
}

func (s *SectionService) UpdateSection(ctx context.Context, adminID, sectionID string, title, description string, orderIndex int) (*domain.Section, error) {
	section, err := s.sectionRepo.FindByID(ctx, sectionID)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	course, err := s.courseRepo.FindByID(ctx, section.CourseID)
	if err != nil {
		return nil, domain.ErrCourseNotFound
	}
	if course.AdminID != adminID {
		return nil, domain.ErrForbidden
	}

	if title != "" {
		section.Title = title
	}
	if description != "" {
		section.Description = description
	}
	if orderIndex >= 0 {
		section.OrderIndex = orderIndex
	}

	if err := s.sectionRepo.Update(ctx, section); err != nil {
		return nil, err
	}

	return section, nil
}

func (s *SectionService) DeleteSection(ctx context.Context, adminID, sectionID string) error {
	section, err := s.sectionRepo.FindByID(ctx, sectionID)
	if err != nil {
		return domain.ErrNotFound
	}
	course, err := s.courseRepo.FindByID(ctx, section.CourseID)
	if err != nil {
		return domain.ErrCourseNotFound
	}
	if course.AdminID != adminID {
		return domain.ErrForbidden
	}
	return s.sectionRepo.DeleteByID(ctx, sectionID)
}

func (s *SectionService) GetSection(ctx context.Context, sectionID string) (*dto.SectionDetail, error) {
	section, err := s.sectionRepo.FindByIDWithLessons(ctx, sectionID)
	if err != nil {
		return nil, err
	}
	return dto.ToSectionDetail(section), nil
}

func (s *SectionService) ListSectionsByCourse(ctx context.Context, courseID string) ([]dto.SectionDetail, error) {
	sections, err := s.sectionRepo.ListByCourseID(ctx, courseID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.SectionDetail, len(sections))
	for i, sec := range sections {
		result[i] = *dto.ToSectionDetail(&sec)
	}
	return result, nil
}

func (s *SectionService) GetCourseWithSections(ctx context.Context, courseID string, userID string) (*dto.CourseDetailWithSections, error) {
	course, err := s.courseRepo.FindByIDWithSections(ctx, courseID)
	if err != nil {
		return nil, err
	}

	totalDurationMinutes := 0
	for _, section := range course.Sections {
		for _, lesson := range section.Lessons {
			totalDurationMinutes += lesson.Duration
		}
	}

	sections := make([]dto.SectionBrief, len(course.Sections))
	for i, s := range course.Sections {
		sections[i] = dto.ToSectionBrief(s)
	}

	tools := dto.ParseTools(course.Tools)
	reviews := dto.ToReviewItems(course.Reviews)

	result := &dto.CourseDetailWithSections{
		ID:                 course.ID,
		Title:              course.Title,
		Description:        course.Description,
		Thumbnail:          course.Thumbnail,
		Price:              course.Price,
		Category:           dto.ToCategoryBrief(course.Category),
		AdminID:            course.AdminID,
		AdminName:          course.Admin.Name,
		Mentor:             dto.ToUserBriefPtr(&course.Mentor, course.MentorID),
		Tools:              tools,
		Hours:              course.Hours,
		TotalDurationHours: totalDurationMinutes / 60,
		Visibility:         course.Visibility,
		Sections:           sections,
		Reviews:            reviews,
		CreatedAt:          course.CreatedAt,
	}

	if userID != "" {
		enrolled, _ := s.enrollRepo.Exists(ctx, userID, courseID)
		result.IsEnrolled = enrolled
	}

	return result, nil
}
