package repository

import (
	"context"
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

type SectionRepository struct {
	db *gorm.DB
}

func NewSectionRepository(db *gorm.DB) *SectionRepository {
	return &SectionRepository{db: db}
}

func (r *SectionRepository) Create(ctx context.Context, section *domain.Section) error {
	return r.db.WithContext(ctx).Create(section).Error
}

func (r *SectionRepository) FindByID(ctx context.Context, sectionID string) (*domain.Section, error) {
	section := &domain.Section{}
	err := r.db.WithContext(ctx).
		Where("id = ?", sectionID).
		Preload("Course").
		First(section).Error
	return section, err
}

func (r *SectionRepository) FindByIDWithLessons(ctx context.Context, sectionID string) (*domain.Section, error) {
	section := &domain.Section{}
	err := r.db.WithContext(ctx).
		Where("id = ?", sectionID).
		Preload("Course").
		Preload("Lessons", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC")
		}).
		First(section).Error
	return section, err
}

func (r *SectionRepository) ListByCourseID(ctx context.Context, courseID string) ([]domain.Section, error) {
	var sections []domain.Section
	err := r.db.WithContext(ctx).
		Where("course_id = ?", courseID).
		Order("order_index ASC").
		Find(&sections).Error
	return sections, err
}

func (r *SectionRepository) ListByCourseIDWithLessons(ctx context.Context, courseID string) ([]domain.Section, error) {
	var sections []domain.Section
	err := r.db.WithContext(ctx).
		Where("course_id = ?", courseID).
		Order("order_index ASC").
		Preload("Lessons", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC")
		}).
		Find(&sections).Error
	return sections, err
}

func (r *SectionRepository) Update(ctx context.Context, section *domain.Section) error {
	return r.db.WithContext(ctx).Model(section).Updates(section).Error
}

func (r *SectionRepository) DeleteByID(ctx context.Context, sectionID string) error {
	return r.db.WithContext(ctx).Where("id = ?", sectionID).Delete(&domain.Section{}).Error
}

func (r *SectionRepository) ExistsInCourse(ctx context.Context, sectionID, courseID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Section{}).
		Where("id = ? AND course_id = ?", sectionID, courseID).
		Count(&count).Error
	return count > 0, err
}
