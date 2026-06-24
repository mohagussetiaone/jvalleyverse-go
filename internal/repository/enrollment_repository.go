package repository

import (
	"context"
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

type EnrollmentRepository struct {
	db *gorm.DB
}

func NewEnrollmentRepository(db *gorm.DB) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}

func (r *EnrollmentRepository) Create(ctx context.Context, enrollment *domain.CourseEnrollment) error {
	return r.db.WithContext(ctx).Create(enrollment).Error
}

func (r *EnrollmentRepository) Exists(ctx context.Context, userID, courseID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.CourseEnrollment{}).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Count(&count).Error
	return count > 0, err
}

func (r *EnrollmentRepository) ListByUserID(ctx context.Context, userID string, page, limit int) ([]domain.CourseEnrollment, int64, error) {
	var enrollments []domain.CourseEnrollment
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.CourseEnrollment{}).
		Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	if err := query.
		Preload("Course").
		Preload("Course.Category").
		Preload("Course.Admin").
		Preload("Course.Mentor").
		Preload("Course.Sections").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&enrollments).Error; err != nil {
		return nil, 0, err
	}

	return enrollments, total, nil
}

func (r *EnrollmentRepository) DeleteByUserAndCourse(ctx context.Context, userID, courseID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Delete(&domain.CourseEnrollment{}).Error
}

func (r *EnrollmentRepository) FindByUserAndCourse(ctx context.Context, userID, courseID string) (*domain.CourseEnrollment, error) {
	var enrollment domain.CourseEnrollment
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		First(&enrollment).Error
	if err != nil {
		return nil, err
	}
	return &enrollment, nil
}

func (r *EnrollmentRepository) UpdateLastLesson(ctx context.Context, userID, courseID string, lessonID string) error {
	return r.db.WithContext(ctx).
		Model(&domain.CourseEnrollment{}).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Update("last_lesson_id", lessonID).Error
}
