# SSE Real-time Notifications — Technical Documentation

## Overview

**SSE (Server-Sent Events)** adalah protokol HTTP long-lived connection yang memungkinkan server mengirim data ke client secara real-time. Berbeda dengan WebSocket yang bidirectional, SSE bersifat **one-way** (server → client), sangat cocok untuk notifikasi.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Go Fiber Server                          │
│                                                                 │
│  ┌──────────────────────┐            ┌──────────────────────┐  │
│  │ NotificationService   │            │   NotificationHub    │  │
│  │                      │            │   (Singleton)        │  │
│  │  CreateNotification() │───────────>│                      │  │
│  │  ┌─ simpan ke DB     │            │  subscribers:        │  │
│  │  └─ hub.Publish()    │            │  map[userID]         │  │
│  └──────────────────────┘            │    map[connID]chan   │  │
│                                       └──────────┬───────────┘  │
│                                                  │              │
│               ┌──────────────────────────────────┼──────┐       │
│               │              Subscribe            │      │       │
│               ▼                                   ▼      │       │
│  ┌─────────────────────┐          ┌─────────────────────┐ │       │
│  │  SSE Handler (User A)│          │  SSE Handler (User B)│ │       │
│  │  ch := hub.Subscribe │          │  ch := hub.Subscribe │ │       │
│  │  for data := range ch│          │  for data := range ch│ │       │
│  │    w.Fprintf(SSE)    │          │    w.Fprintf(SSE)    │ │       │
│  │    w.Flush()         │          │    w.Flush()         │ │       │
│  └─────────────────────┘          └─────────────────────┘ │       │
│           │                                │              │       │
└───────────┼────────────────────────────────┼──────────────┘       │
            │                                │                      │
            ▼                                ▼                      │
    EventSource (Browser)           EventSource (Browser)           │
                                                                     │
    Contoh: User "budi@example.com" punya 2 tab browser → 2 SSE     │
    koneksi. Saat notif masuk, kedua tab langsung dapat event.      │
└──────────────────────────────────────────────────────────────────┘
```

## File Structure

### 1. `internal/service/notification_hub.go`

**Fungsi:** Mengelola semua koneksi SSE per user.

**Key Components:**

- **`subscribers map[string]map[string]chan []byte`**
  - Level 1: key = `userID` (string)
  - Level 2: key = `connID` (string, format `"conn-1"`, `"conn-2"`, dll)
  - Value: channel `chan []byte` dengan buffer 10

- **`sync.RWMutex`** — thread-safe untuk concurrent read/write dari banyak goroutine

**Methods:**

```go
func (h *NotificationHub) Subscribe(userID string) (string, chan []byte)
  // - Buat channel baru (buffer 10)
  // - Simpan di map subscribers[userID][connID]
  // - Return connID + channel

func (h *NotificationHub) Unsubscribe(userID, connID string)
  // - Close channel
  // - Hapus dari map
  // - Cleanup jika user sudah tidak punya koneksi

func (h *NotificationHub) Publish(userID string, event SSEEvent)
  // - Marshal event ke JSON
  // - Iterasi semua channel milik user
  // - Kirim data via channel (non-blocking select)
```

**Kenapa buffer 10?**
- Mencegah blocking jika client lambat membaca
- Notifikasi jarang burst besar (>10 dalam waktu singkat)
- Jika channel penuh, event di-drop dengan log warning

### 2. `internal/handler/sse_handler.go`

**Endpoint:** `GET /api/notifications/stream` (JWT required)

**Alur:**

1. Extract `userID` dari JWT (`c.Locals("userID")`)
2. Set SSE headers:
   - `Content-Type: text/event-stream`
   - `Cache-Control: no-cache`
   - `Connection: keep-alive`
   - `X-Accel-Buffering: no` (untuk nginx)
3. Subscribe ke NotificationHub → dapat `connID` + `ch`
4. `defer hub.Unsubscribe(userID, connID)` — cleanup otomatis
5. Gunakan `c.Context().SetBodyStreamWriter(func(w *bufio.Writer) { ... })`
6. Dalam stream writer:
   - Kirim event `connected` sebagai konfirmasi
   - Loop `select`:
     - `case <-c.Context().Done()`: client disconnect → return
     - `case data := <-ch`: ada notif baru → tulis SSE event + flush
     - `case <-heartbeat.C`: tiap 30 detik → kirim `: heartbeat\n\n`
7. Jika koneksi putus (client tutup browser), Fiber otomatis cancel context

**Header Penting:**
```go
c.Set("X-Accel-Buffering", "no")  // Mencegah nginx buffering
```
Tanpa header ini, nginx/proxy bisa menahan response sampai buffer penuh, menyebabkan delay.

### 3. `internal/service/notification_service.go`

**Perubahan:** `CreateNotification()` sekarang:

```go
func (s *NotificationService) CreateNotification(...) error {
    // 1. Simpan notif ke database (seperti biasa)
    notif := &domain.Notification{...}
    s.notifRepo.Create(ctx, notif)

    // 2. Push real-time via SSE hub
    hub := GetNotificationHub()
    hub.Publish(userID, SSEEvent{
        Type:    nType,
        Title:   title,
        Message: message,
        Link:    link,
        Payload: map[string]interface{}{
            "id":         notif.ID,
            "created_at": notif.CreatedAt,
            "is_read":    false,
        },
    })
}
```

## SSE Protocol Format

```
event: connected
data: {"status":"connected","user_id":"abc123"}

event: notification
data: {"type":"new_reply","title":"Reply Baru","message":"Budi membalas diskusi Anda","link":"/discussions/xyz","payload":{"id":"notif123","created_at":"2024-01-20T10:00:00Z","is_read":false}}

: heartbeat
```

**Format Rules:**
- Setiap event dipisahkan oleh `\n\n` (double newline)
- `event: <nama_event>` — jenis event (optional)
- `data: <json>` — payload event
- `: <comment>` — komentar (digunakan untuk heartbeat)
- Baris yang dimulai `:` diabaikan browser

## Performance Analysis

### Memory per Connection

| Komponen | Ukuran |
|----------|--------|
| Goroutine | ~4 KB |
| Channel (buffer 10 × 256 bytes) | ~2.5 KB |
| Map entry | ~0.5 KB |
| **Total per koneksi** | **~7 KB** |

### Estimasi untuk VPS

| User Online | SSE Koneksi | Memory | CPU | Bandwidth |
|-------------|-------------|--------|-----|-----------|
| 100 | 100 | ~0.7 MB | idle | ~1 KB/notifikasi |
| 500 | 500 | ~3.5 MB | idle | ~1 KB/notifikasi |
| 1000 | 1000 | ~7 MB | idle | ~1 KB/notifikasi |
| 5000 | 5000 | ~35 MB | idle | ~1 KB/notifikasi |

**Catatan:** CPU hanya aktif saat ada notifikasi baru (write + flush). Jika rata-rata 1 notifikasi/detik, CPU usage tidak signifikan.

### Perbandingan dengan Polling

| Metrik | SSE | Polling (setiap 5 detik) |
|--------|-----|--------------------------|
| HTTP Request/jam | 1 (long-lived) | 720 per user |
| DB Query/jam | 0 (push dari memory) | 720 per user |
| Latensi notifikasi | <100ms (real-time) | ~5 detik (max) |
| Bandwidth idle | ~40 byte/30 detik (heartbeat) | ~2 KB/response |
| Resource VPS 500 user | ~3.5 MB + idle CPU | ~1000 req/menit ke DB |

## Client Implementation

### JavaScript (Browser)

```js
class NotificationStream {
  constructor(token) {
    this.token = token;
    this.eventSource = null;
    this.listeners = new Map();
  }

  connect() {
    this.eventSource = new EventSource('/api/notifications/stream', {
      headers: { Authorization: 'Bearer ' + this.token }
    });

    this.eventSource.addEventListener('connected', (e) => {
      console.log('SSE connected:', JSON.parse(e.data));
    });

    this.eventSource.addEventListener('notification', (e) => {
      const notif = JSON.parse(e.data);
      // Panggil callback yang terdaftar
      this.listeners.forEach(cb => cb(notif));
      // Contoh: showToast(notif.title, notif.message)
    });

    this.eventSource.onerror = (e) => {
      if (this.eventSource.readyState === EventSource.CLOSED) {
        console.log('SSE connection closed. Reconnecting...');
      }
      // EventSource auto-reconnect by default
    };
  }

  onNotification(callback) {
    const id = Date.now();
    this.listeners.set(id, callback);
    return () => this.listeners.delete(id);
  }

  disconnect() {
    if (this.eventSource) {
      this.eventSource.close();
      this.eventSource = null;
    }
  }
}

// Usage
const notifStream = new NotificationStream('jwt_token_here');
notifStream.connect();

notifStream.onNotification((notif) => {
  // Update badge count
  updateBadge();
  // Show toast
  showToast(notif.title, notif.message);
  // Play sound
  playNotificationSound();
});
```

### React

```tsx
import { useEffect, useState } from 'react';

function useNotificationStream(token: string | null) {
  const [notifications, setNotifications] = useState<any[]>([]);

  useEffect(() => {
    if (!token) return;

    const evtSource = new EventSource('/api/notifications/stream', {
      headers: { Authorization: 'Bearer ' + token }
    });

    evtSource.addEventListener('notification', (e) => {
      const notif = JSON.parse(e.data);
      setNotifications(prev => [notif, ...prev].slice(0, 50));
    });

    evtSource.onerror = () => {
      // Auto-reconnect by EventSource
    };

    return () => evtSource.close();
  }, [token]);

  return notifications;
}
```

## Error Handling

### Server Side

| Skenario | Handling |
|----------|----------|
| Client disconnect | `c.Context().Done()` → return → `defer Unsubscribe()` cleanup |
| Channel full (client lambat) | `select default` → drop event + log warning |
| Write error | Log error → return → cleanup |
| JWT invalid | Middleware return 401 sebelum handler |

### Client Side

| Skenario | Handling |
|----------|----------|
| Koneksi putus | EventSource **auto-reconnect** (built-in) |
| Server restart | Client reconnect, dapat event `connected` lagi |
| Unauthorized (401) | EventSource error → redirect ke login |
| Network timeout | Browser retry dengan interval exponential backoff |

## Potential Improvements

1. **Redis Pub/Sub** — Jika ada multiple server instances (horizontal scaling), hub harus pakai Redis Pub/Sub agar notifikasi dari server A bisa dikirim ke koneksi SSE di server B.

2. **Unread Count via SSE** — Kirim unread_count setiap ada notif baru, jadi frontend bisa update badge tanpa perlu REST call.

3. **Notification Persistence + SSE** — Saat user reconnect, kirim notifikasi yang belum dibaca sebagai initial state.

4. **Connection Health** — Track jumlah koneksi per user untuk monitoring.

## Debugging

### Test dengan curl
```bash
# Buka SSE stream
curl -N -H "Authorization: Bearer <token>" \
  http://localhost:3000/api/notifications/stream

# Output:
# event: connected
# data: {"status":"connected","user_id":"abc123"}
#
# : heartbeat
# : heartbeat
```

### Test dengan Postman
1. Method: GET
2. URL: `http://localhost:3000/api/notifications/stream`
3. Headers: `Authorization: Bearer <token>`
4. Klik "Send" → lihat response streaming di bagian bawah

### Log Server
```
[SSE] User abc123 connected (total: 5 connections)
[SSE] New notification → pushing to user abc123
[SSE] User abc123 disconnected (total: 4 connections)
```
