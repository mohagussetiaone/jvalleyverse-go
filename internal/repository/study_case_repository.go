package repository

import (
	"context"
	"jvalleyverse/internal/domain"

	"gorm.io/gorm"
)

type StudyCaseRepository struct {
	db *gorm.DB
}

func NewStudyCaseRepository(db *gorm.DB) *StudyCaseRepository {
	return &StudyCaseRepository{db: db}
}

type StudyCaseListFilter struct {
	CategoryID *string
}

func (r *StudyCaseRepository) Create(ctx context.Context, studyCase *domain.StudyCase) error {
	return r.db.WithContext(ctx).Create(studyCase).Error
}

func (r *StudyCaseRepository) FindByID(ctx context.Context, id string) (*domain.StudyCase, error) {
	studyCase := &domain.StudyCase{}
	if err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("User").
		Preload("Discussions", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(10)
		}).
		Preload("Discussions.User").
		Preload("Discussions.Replies", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Preload("Discussions.Replies.User").
		First(studyCase, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return studyCase, nil
}

func (r *StudyCaseRepository) ListAll(ctx context.Context, page, limit int, filter *StudyCaseListFilter) ([]domain.StudyCase, int64, error) {
	var studyCases []domain.StudyCase
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.StudyCase{})

	if filter != nil && filter.CategoryID != nil && *filter.CategoryID != "" {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.
		Preload("Category").
		Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&studyCases).Error; err != nil {
		return nil, 0, err
	}

	return studyCases, total, nil
}

func (r *StudyCaseRepository) ListByUserID(ctx context.Context, userID string, page, limit int) ([]domain.StudyCase, int64, error) {
	var studyCases []domain.StudyCase
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.StudyCase{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&studyCases).Error; err != nil {
		return nil, 0, err
	}

	return studyCases, total, nil
}

func (r *StudyCaseRepository) Update(ctx context.Context, studyCase *domain.StudyCase) error {
	return r.db.WithContext(ctx).Model(studyCase).Updates(studyCase).Error
}

func (r *StudyCaseRepository) DeleteByID(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.StudyCase{}).Error
}
