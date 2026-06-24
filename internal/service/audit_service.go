package service

import (
	"context"
	"jvalleyverse/internal/domain"
	"jvalleyverse/internal/repository"
)

type AuditService struct {
	auditRepo *repository.AdminAuditLogRepository
}

func NewAuditService(auditRepo *repository.AdminAuditLogRepository) *AuditService {
	return &AuditService{auditRepo: auditRepo}
}

func (s *AuditService) Log(ctx context.Context, adminID, action, resourceType, resourceID, details string) error {
	log := &domain.AdminAuditLog{
		AdminID:      adminID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Details:      details,
	}
	return s.auditRepo.Create(ctx, log)
}

func (s *AuditService) List(ctx context.Context, page, limit int) ([]domain.AdminAuditLog, int64, error) {
	return s.auditRepo.List(ctx, page, limit)
}
