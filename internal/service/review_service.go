package service

import (
	"context"
	"fmt"
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/dto"
	"jvalleyverse/internal/repository"
)

type IReviewService interface {
	CreateReview(ctx context.Context, userID, courseID, lessonID string, rating int, message string) (*domain.Review, error)
	UpdateReview(ctx context.Context, reviewID, userID string, rating int, message string) (*domain.Review, error)
	DeleteReview(ctx context.Context, reviewID, userID string) error
	ListCourseReviews(ctx context.Context, courseID string) ([]dto.ReviewItem, error)
	ListLessonReviews(ctx context.Context, lessonID string) ([]dto.ReviewItem, error)
}

type ReviewService struct {
	reviewRepo *repository.ReviewRepository
	courseRepo *repository.CourseRepository
}

func NewReviewService(reviewRepo *repository.ReviewRepository, courseRepo *repository.CourseRepository) *ReviewService {
	return &ReviewService{reviewRepo: reviewRepo, courseRepo: courseRepo}
}

func (s *ReviewService) CreateReview(ctx context.Context, userID, courseID, lessonID string, rating int, message string) (*domain.Review, error) {
	if rating < 1 || rating > 5 {
		return nil, domain.ErrInvalidInput
	}
	if message == "" {
		return nil, domain.ErrInvalidInput
	}

	review := &domain.Review{
		UserID:   userID,
		CourseID: courseID,
		LessonID: lessonID,
		Rating:   rating,
		Message:  message,
	}

	if err := s.reviewRepo.Create(ctx, review); err != nil {
		return nil, err
	}

	created, err := s.reviewRepo.FindByID(ctx, review.ID)
	if err == nil {
		// Notify course admin about new review
		if notifSvc := GetNotificationService(); notifSvc != nil && courseID != "" {
			course, courseErr := s.courseRepo.FindByID(ctx, courseID)
			if courseErr == nil && course.AdminID != userID {
				notifSvc.CreateNotification(ctx, course.AdminID, "new_review",
					"Review Baru",
					"Ada review baru untuk kursus: "+course.Title+" (rating: "+fmt.Sprintf("%d", rating)+"/5)",
					"/courses/"+courseID,
				)
			}
		}
	}

	return created, nil
}

func (s *ReviewService) UpdateReview(ctx context.Context, reviewID, userID string, rating int, message string) (*domain.Review, error) {
	review, err := s.reviewRepo.FindByID(ctx, reviewID)
	if err != nil {
		return nil, err
	}

	if review.UserID != userID {
		return nil, domain.ErrForbidden
	}

	if rating >= 1 && rating <= 5 {
		review.Rating = rating
	}
	if message != "" {
		review.Message = message
	}

	if err := s.reviewRepo.Update(ctx, review); err != nil {
		return nil, err
	}

	return s.reviewRepo.FindByID(ctx, review.ID)
}

func (s *ReviewService) DeleteReview(ctx context.Context, reviewID, userID string) error {
	review, err := s.reviewRepo.FindByID(ctx, reviewID)
	if err != nil {
		return err
	}

	if review.UserID != userID {
		return domain.ErrForbidden
	}

	return s.reviewRepo.DeleteByID(ctx, reviewID)
}

func (s *ReviewService) ListCourseReviews(ctx context.Context, courseID string) ([]dto.ReviewItem, error) {
	reviews, err := s.reviewRepo.ListByCourse(ctx, courseID)
	if err != nil {
		return nil, err
	}
	return dto.ToReviewItems(reviews), nil
}

func (s *ReviewService) ListLessonReviews(ctx context.Context, lessonID string) ([]dto.ReviewItem, error) {
	reviews, err := s.reviewRepo.ListByLesson(ctx, lessonID)
	if err != nil {
		return nil, err
	}
	return dto.ToReviewItems(reviews), nil
}
