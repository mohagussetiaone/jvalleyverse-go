package handler

import (
	"bufio"
	"fmt"
	"jvalleyverse/internal/service"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// SSEHandler handles Server-Sent Events for real-time notifications
type SSEHandler struct {
	hub *service.NotificationHub
}

func NewSSEHandler() *SSEHandler {
	return &SSEHandler{
		hub: service.GetNotificationHub(),
	}
}

// StreamNotifications streams real-time notifications via SSE for the authenticated user
func (h *SSEHandler) StreamNotifications(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	// Set SSE headers
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("X-Accel-Buffering", "no")

	// Capture the Done channel while still in the handler goroutine
	// where the RequestCtx is guaranteed valid.
	// The SetBodyStreamWriter callback runs in a separate goroutine
	// AFTER the handler returns, at which point c.Context() may
	// have been recycled — calling Done() on it would panic.
	done := c.Context().Done()

	// Use SetBodyStreamWriter for streaming
	// Subscribe/Unsubscribe live INSIDE the callback because
	// SetBodyStreamWriter runs asynchronously — the handler returns
	// immediately and defer in the outer function would fire before
	// the stream writer goroutine starts, closing the channel early.
	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		connID, ch := h.hub.Subscribe(userID)
		defer h.hub.Unsubscribe(userID, connID)

		// Send initial connection event
		initData := fmt.Sprintf(`{"status":"connected","user_id":"%s"}`, userID)
		if _, err := fmt.Fprintf(w, "event: connected\ndata: %s\n\n", initData); err != nil {
			log.Printf("[SSE] Initial write error for user %s: %v", userID, err)
			return
		}
		if err := w.Flush(); err != nil {
			log.Printf("[SSE] Initial flush error for user %s: %v", userID, err)
			return
		}

		// Send heartbeat every 30 seconds to keep connection alive
		heartbeat := time.NewTicker(30 * time.Second)
		defer heartbeat.Stop()

		for {
			select {
			case <-done:
				log.Printf("[SSE] Connection closed for user %s (conn: %s)", userID, connID)
				return

			case data, ok := <-ch:
				if !ok {
					return
				}
				if _, err := fmt.Fprintf(w, "event: notification\ndata: %s\n\n", string(data)); err != nil {
					log.Printf("[SSE] Write error for user %s: %v", userID, err)
					return
				}
				if err := w.Flush(); err != nil {
					log.Printf("[SSE] Flush error for user %s: %v", userID, err)
					return
				}

			case <-heartbeat.C:
				if _, err := fmt.Fprintf(w, ": heartbeat\n\n"); err != nil {
					return
				}
				if err := w.Flush(); err != nil {
					return
				}
			}
		}
	})

	return nil
}
