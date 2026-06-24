package service

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

// SSEEvent represents a notification event sent via SSE
type SSEEvent struct {
	Type    string      `json:"type"`
	Title   string      `json:"title"`
	Message string      `json:"message"`
	Link    string      `json:"link,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

// NotificationHub manages SSE connections per user
type NotificationHub struct {
	mu          sync.RWMutex
	subscribers map[string]map[string]chan []byte // userID → {connID → channel}
}

var (
	hub     *NotificationHub
	hubOnce sync.Once
)

// GetNotificationHub returns the singleton hub instance
func GetNotificationHub() *NotificationHub {
	hubOnce.Do(func() {
		hub = &NotificationHub{
			subscribers: make(map[string]map[string]chan []byte),
		}
	})
	return hub
}

// Subscribe adds a new SSE connection for a user, returns connection ID and channel
func (h *NotificationHub) Subscribe(userID string) (string, chan []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	connID := generateConnID()
	ch := make(chan []byte, 10) // buffered to avoid blocking

	if h.subscribers[userID] == nil {
		h.subscribers[userID] = make(map[string]chan []byte)
	}
	h.subscribers[userID][connID] = ch

	log.Printf("[SSE] User %s connected (total: %d connections)", userID, h.countAll())
	return connID, ch
}

// Unsubscribe removes an SSE connection for a user
func (h *NotificationHub) Unsubscribe(userID, connID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if userConns, ok := h.subscribers[userID]; ok {
		if ch, exists := userConns[connID]; exists {
			close(ch)
			delete(userConns, connID)
		}
		if len(userConns) == 0 {
			delete(h.subscribers, userID)
		}
	}
	log.Printf("[SSE] User %s disconnected (total: %d connections)", userID, h.countAll())
}

// Publish sends an event to all SSE connections for a user
func (h *NotificationHub) Publish(userID string, event SSEEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("[SSE] Failed to marshal event: %v", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	userConns, ok := h.subscribers[userID]
	if !ok {
		return
	}

	for connID, ch := range userConns {
		select {
		case ch <- data:
		default:
			log.Printf("[SSE] Channel full for connection %s, dropping event", connID)
		}
	}
}

// PublishJSON sends raw JSON data to all SSE connections for a user
func (h *NotificationHub) PublishJSON(userID string, data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	userConns, ok := h.subscribers[userID]
	if !ok {
		return
	}

	for connID, ch := range userConns {
		select {
		case ch <- data:
		default:
			log.Printf("[SSE] Channel full for connection %s, dropping event", connID)
		}
	}
}

func (h *NotificationHub) countAll() int {
	count := 0
	for _, conns := range h.subscribers {
		count += len(conns)
	}
	return count
}

var connIDCounter int
var connIDMu sync.Mutex

func generateConnID() string {
	connIDMu.Lock()
	defer connIDMu.Unlock()
	connIDCounter++
	return fmt.Sprintf("conn-%d", connIDCounter)
}
