package handler

import (
	"context"
	"jvalleyverse/internal/service"
)

var auditSvc *service.AuditService

func InitAuditService(svc *service.AuditService) {
	auditSvc = svc
}

func logAdminAction(ctx context.Context, adminID, action, resourceType, resourceID, details string) {
	if auditSvc == nil {
		return
	}
	auditSvc.Log(ctx, adminID, action, resourceType, resourceID, details)
}
