package repository

import (
	"context"
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

type CourseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) Create(ctx context.Context, course *domain.Course) error {
	return r.db.WithContext(ctx).Create(course).Error
}

func (r *CourseRepository) FindByID(ctx context.Context, courseID string) (*domain.Course, error) {
	course := &domain.Course{}
	if err := r.db.WithContext(ctx).Where("id = ?", courseID).
		Preload("Admin").
		Preload("Category").
		Preload("Mentor").
		First(course).Error; err != nil {
		return nil, err
	}
	return course, nil
}

func (r *CourseRepository) FindByIDWithSections(ctx context.Context, courseID string) (*domain.Course, error) {
	course := &domain.Course{}
	if err := r.db.WithContext(ctx).Where("id = ?", courseID).
		Preload("Admin").
		Preload("Category").
		Preload("Mentor").
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC")
		}).
		Preload("Sections.Lessons", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC")
		}).
		Preload("Sections.Lessons.Details").
		Preload("Sections.Lessons.NextLesson").
		Preload("Reviews", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(10)
		}).
		Preload("Reviews.User").
		First(course).Error; err != nil {
		return nil, err
	}
	return course, nil
}

type CourseListFilter struct {
	CategoryID *string
	MinPrice   *float64
	MaxPrice   *float64
}

func (r *CourseRepository) ListPublic(
	ctx context.Context,
	page, limit int,
	filter *CourseListFilter,
) ([]domain.Course, int64, error) {

	var courses []domain.Course
	var total int64

	query := r.db.WithContext(ctx).
		Model(&domain.Course{}).
		Where("visibility = ?", "public")

	if filter != nil {
		if filter.CategoryID != nil && *filter.CategoryID != "" {
			query = query.Where("category_id = ?", *filter.CategoryID)
		}
		if filter.MinPrice != nil {
			query = query.Where("price >= ?", *filter.MinPrice)
		}
		if filter.MaxPrice != nil {
			query = query.Where("price <= ?", *filter.MaxPrice)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit

	if err := query.
		Preload("Admin").
		Preload("Category").
		Preload("Mentor").
		Preload("Sections").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}

func (r *CourseRepository) ListByAdminID(ctx context.Context, adminID string, page, limit int) ([]domain.Course, int64, error) {
	var courses []domain.Course
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Course{}).Where("admin_id = ?", adminID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).
		Where("admin_id = ?", adminID).
		Preload("Category").
		Preload("Mentor").
		Preload("Sections").
		Offset(offset).
		Limit(limit).
		Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}

func (r *CourseRepository) ListByMentorID(ctx context.Context, mentorID string, page, limit int) ([]domain.Course, int64, error) {
	var courses []domain.Course
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Course{}).Where("mentor_id = ?", mentorID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).
		Where("mentor_id = ?", mentorID).
		Preload("Category").
		Preload("Admin").
		Preload("Mentor").
		Preload("Sections").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&courses).Error; err != nil {
		return nil, 0, err
	}

	return courses, total, nil
}

func (r *CourseRepository) Update(ctx context.Context, course *domain.Course) error {
	return r.db.WithContext(ctx).Model(course).Updates(course).Error
}

func (r *CourseRepository) DeleteByID(ctx context.Context, courseID string) error {
	return r.db.WithContext(ctx).Where("id = ?", courseID).Delete(&domain.Course{}).Error
}
