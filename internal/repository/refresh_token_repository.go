package repository

import (
	"context"
	"jvalleyverse/internal/domain"
	"time"

	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *RefreshTokenRepository) FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	var t domain.RefreshToken
	if err := r.db.WithContext(ctx).Where("token = ? AND revoked_at IS NULL AND expires_at > ?", token, time.Now()).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *RefreshTokenRepository) RevokeByUserID(ctx context.Context, userID string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&domain.RefreshToken{}).Where("user_id = ? AND revoked_at IS NULL", userID).Update("revoked_at", &now).Error
}

func (r *RefreshTokenRepository) RevokeByToken(ctx context.Context, token string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&domain.RefreshToken{}).Where("token = ?", token).Update("revoked_at", &now).Error
}

func (r *RefreshTokenRepository) CleanupExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&domain.RefreshToken{}).Error
}
