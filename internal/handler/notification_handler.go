package handler

import (
	"jvalleyverse/internal/service"

	"github.com/gofiber/fiber/v2"
)

type NotificationHandler struct {
	notifSvc service.INotificationService
}

func NewNotificationHandler(notifSvc service.INotificationService) *NotificationHandler {
	return &NotificationHandler{notifSvc: notifSvc}
}

// ListNotifications returns user's notifications
func (h *NotificationHandler) ListNotifications(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	page := clampPage(c.QueryInt("page", DefaultPage))
	limit := clampLimit(c.QueryInt("limit", DefaultLimit), DefaultLimit)

	notifs, total, err := h.notifSvc.ListNotifications(c.UserContext(), userID, page, limit)
	if err != nil {
		return safeError(c, 500, err)
	}

	unread, _ := h.notifSvc.CountUnread(c.UserContext(), userID)

	return c.JSON(fiber.Map{
		"data": notifs,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
		"unread_count": unread,
	})
}

// CountUnread returns unread notifications count
func (h *NotificationHandler) CountUnread(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	count, err := h.notifSvc.CountUnread(c.UserContext(), userID)
	if err != nil {
		return safeError(c, 500, err)
	}

	return c.JSON(fiber.Map{"unread_count": count})
}

// MarkAsRead marks a single notification as read
func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	notifID := c.Params("id")
	if notifID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid notification ID"})
	}

	if err := h.notifSvc.MarkAsRead(c.UserContext(), notifID, userID); err != nil {
		return safeError(c, 500, err)
	}

	return c.JSON(fiber.Map{"message": "Notification marked as read"})
}

// MarkAllAsRead marks all user notifications as read
func (h *NotificationHandler) MarkAllAsRead(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	if err := h.notifSvc.MarkAllAsRead(c.UserContext(), userID); err != nil {
		return safeError(c, 500, err)
	}

	return c.JSON(fiber.Map{"message": "All notifications marked as read"})
}

// DeleteNotification deletes a notification
func (h *NotificationHandler) DeleteNotification(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	notifID := c.Params("id")
	if notifID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid notification ID"})
	}

	if err := h.notifSvc.DeleteNotification(c.UserContext(), notifID, userID); err != nil {
		return safeError(c, 500, err)
	}

	return c.JSON(fiber.Map{"message": "Notification deleted"})
}
