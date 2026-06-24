package repository

import (
	"context"
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

type LessonRepository struct {
	db *gorm.DB
}

func NewLessonRepository(db *gorm.DB) *LessonRepository {
	return &LessonRepository{db: db}
}

func (r *LessonRepository) FindPublicByID(ctx context.Context, lessonID string) (*domain.Lesson, error) {
	lesson := &domain.Lesson{}

	if err := r.db.WithContext(ctx).
		Where("id = ?", lessonID).
		Where("visibility = ?", "public").
		Preload("Course").
		Preload("Section").
		Preload("Admin").
		Preload("Details").
		Preload("NextLesson").
		First(lesson).Error; err != nil {
		return nil, err
	}

	return lesson, nil
}

func (r *LessonRepository) Create(ctx context.Context, lesson *domain.Lesson) error {
	return r.db.WithContext(ctx).Create(lesson).Error
}

func (r *LessonRepository) FindByID(ctx context.Context, lessonID string) (*domain.Lesson, error) {
	lesson := &domain.Lesson{}
	if err := r.db.WithContext(ctx).Where("id = ?", lessonID).
		Preload("Course").Preload("Admin").Preload("Details").First(lesson).Error; err != nil {
		return nil, err
	}
	return lesson, nil
}

func (r *LessonRepository) FindBySlug(ctx context.Context, courseID string, slug string) (*domain.Lesson, error) {
	lesson := &domain.Lesson{}
	if err := r.db.WithContext(ctx).Where("course_id = ? AND slug = ?", courseID, slug).
		Preload("Course").
		Preload("Admin").
		Preload("Details").
		Preload("NextLesson").
		First(lesson).Error; err != nil {
		return nil, err
	}
	return lesson, nil
}

func (r *LessonRepository) FindNextLesson(ctx context.Context, nextLessonID string) (*domain.Lesson, error) {
	lesson := &domain.Lesson{}
	if err := r.db.WithContext(ctx).Where("id = ?", nextLessonID).First(lesson).Error; err != nil {
		return nil, err
	}
	return lesson, nil
}

func (r *LessonRepository) ListByCourseID(ctx context.Context, courseID string, limit, offset int) ([]domain.Lesson, int64, error) {
	var lessons []domain.Lesson
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Lesson{}).Where("course_id = ?", courseID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Where("course_id = ?", courseID).
		Preload("Admin").
		Order("order_index ASC").
		Offset(offset).Limit(limit).Find(&lessons).Error; err != nil {
		return nil, 0, err
	}

	return lessons, total, nil
}

func (r *LessonRepository) ListBySectionID(ctx context.Context, sectionID string) ([]domain.Lesson, int64, error) {
	var lessons []domain.Lesson
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Lesson{}).Where("section_id = ?", sectionID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Where("section_id = ?", sectionID).
		Preload("Admin").
		Order("order_index ASC").
		Find(&lessons).Error; err != nil {
		return nil, 0, err
	}

	return lessons, total, nil
}

func (r *LessonRepository) ListAll(ctx context.Context, page, limit int, difficulty *string) ([]domain.Lesson, int64, error) {
	var lessons []domain.Lesson
	var total int64

	query := r.db.WithContext(ctx)
	if difficulty != nil {
		query = query.Where("difficulty = ?", *difficulty)
	}

	if err := query.Model(&domain.Lesson{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Preload("Course").Preload("Admin").Offset(offset).Limit(limit).Find(&lessons).Error; err != nil {
		return nil, 0, err
	}

	return lessons, total, nil
}

func (r *LessonRepository) Update(ctx context.Context, lesson *domain.Lesson) error {
	return r.db.WithContext(ctx).Model(lesson).Updates(lesson).Error
}

func (r *LessonRepository) DeleteByID(ctx context.Context, lessonID string) error {
	return r.db.WithContext(ctx).Where("id = ?", lessonID).Delete(&domain.Lesson{}).Error
}
