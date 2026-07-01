# JValleyverse — Flow & API Reference

## Proyek Overview

Learning platform berbasis web untuk belajar programming. Dibangun dengan **Go (Fiber v2)**, **GORM ORM**, **PostgreSQL**, dan **Redis**.

```
Go 1.25 + Fiber v2 + GORM + PostgreSQL + Redis
```

---

## Domain Model

### Entity Relationship

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│ Category │─1─>│  Course  │─1─>│ Section  │─1─>│  Lesson  │
└──────────┘     └──────────┘     └──────────┘     └──────────┘
                    │                                    │
                    ├── Mentor (User, optional)          ├── LessonDetail (1:1)
                    ├── LearningObjectives (JSON[])      ├── LessonProgress (M:N User)
                    └── Review (M:N User)                ├── Certificate (M:N User)
                                                         └── Review (M:N User)

┌──────────┐     ┌──────────┐     ┌──────────┐
│  User    │─1─>│ Discussion    │─1─>│  Reply   │
│ (admin)  │     └──────────┘     └──────────┘
│ (mentor) │
│          │─1─>│ Showcase │─1─>│ ShowcaseLike   │
│          │     │          │─1─>│ ShowcaseComment │
│          │     └──────────┘     └──────────┘
│          │
│          │──>│ Review   │ (rating 1-5, message)
│          │──>│ CommunityPoint │ (activity log + points)
│          │──>│ RefreshToken   │
│          │──>│ AdminAuditLog  │ (admin actions)
└──────────┘
```

### Schema Quick Reference

| Table              | Key Fields                                                                         |
| ------------------ | ---------------------------------------------------------------------------------- |
| users              | email (unique), name, role (admin/user), points, level                             |
| categories         | name (unique), slug (unique)                                                       |
| courses            | title, category_id, admin_id, mentor_id, learning_objectives (JSON)                |
| sections           | title, course_id, order_index                                                      |
| lessons            | title, slug, course_id, section_id, admin_id, difficulty, next_lesson_id           |
| lesson_details     | lesson_id (unique), about, rules, tools (JSON), resources (JSON)                   |
| lesson_progresses  | user_id, lesson_id, status, progress_percentage                                    |
| course_enrollments | user_id, course_id (unique composite)                                              |
| learning_streaks   | user_id (unique), streak_count, longest_streak, last_activity_date                     |
| certificates       | user_id, lesson_id, unique_code (unique), verification_url, qr_code_url             |
| discussions        | title, content, user_id, lesson_id (nullable), category_id                         |
| replies            | content, user_id, discussion_id, parent_id (nullable)                              |
| showcases          | title, user_id, category_id, media_urls (JSONB)                                    |
| showcase_likes     | user_id + showcase_id (composite PK)                                               |
| showcase_comments  | content, user_id, showcase_id, parent_id                                           |
| community_points   | user_id, activity_type, points_earned, metadata (JSONB)                            |
| user_levels        | level, min_points, max_points, badge_name                                          |
| refresh_tokens     | user_id, token (unique), expires_at                                                |
| admin_audit_logs   | admin_id, action, resource_type, resource_id                                       |
| reviews            | user_id, course_id (nullable), lesson_id (nullable), rating (1-5), message         |
| study_cases        | name, description, img_url, youtube_url, tags (JSONB), user_id                     |
| notifications      | user_id, type, title, message, is_read, link, metadata (JSONB)                     |
| blogs              | title, slug, content, tags (JSONB), status (draft/published), user_id, category_id |

---

## API Route Table

### New Public Routes (Portfolio + Certificate Verification)

```
GET    /api/users/:id/portfolio           — Public user portfolio (aggregates certs, showcases, study cases)
GET    /api/certificates/:code/verify     — Verify certificate (public, no auth)
```

### New Protected Routes

```
GET    /api/users/me/streak               — My learning streak (JWT)
```

### Public (34 routes, some with optional JWT)

Catatan: Route yang ditandai `(Opt.JWT)` menggunakan middleware `OptionalJWTAuth` — response akan menyertakan `is_enrolled` jika user membawa token.

```
POST   /api/auth/register              — Register user (Email)
POST   /api/auth/login                 — Login (Email)
POST   /api/auth/google                — Google One Tap Login (ID token)
POST   /api/auth/refresh               — Refresh token
POST   /api/auth/logout                — Logout (JWT)
GET    /api/leaderboard                — Top users by points
GET    /api/mentors                    — List mentors (paginated)
GET    /api/showcases                  — Paginated showcases
GET    /api/showcases/:id              — Showcase detail
GET    /api/categories                 — All categories
GET    /api/categories/:slug           — Category by slug with courses
GET    /api/categories/:category_id/courses (Opt.JWT) — Courses in category (+ is_enrolled)
GET    /api/courses (Opt.JWT)          — Public courses (+ is_enrolled jika login) (?category_id=&min_price=&max_price=)
GET    /api/courses/:course_id (Opt.JWT) — Course with sections & lessons (+ is_enrolled)
GET    /api/courses/:course_id/sections — List sections
GET    /api/courses/:course_id/sections/:section_id — Section detail
GET    /api/courses/:course_id/reviews  — Reviews for a course (?page=&limit=)
GET    /api/lessons/:id                — Lesson detail
GET    /api/lessons/:id/reviews        — Reviews for a lesson (?page=&limit=)
GET    /api/courses/:course_id/lessons — Lessons in course
GET    /api/courses/:course_id/sections/:section_id/lessons — Lessons in section
GET    /api/courses/:course_id/lessons/:slug — Lesson by slug
GET    /api/blogs                      — Published blogs (?search=&category_id=&tag=)
GET    /api/blogs/:id                  — Blog detail
GET    /api/discussions                — List discussions (Opt.JWT, ?lesson_id=&study_case_id=)
GET    /api/discussions/:id            — Discussion with replies (Opt.JWT)
GET    /api/study-cases                — Paginated study cases
GET    /api/study-cases/:id            — Study case detail
GET    /api/health                     — Health check
GET    /api/health/detailed            — Detailed health + metrics
GET    /api/system/status              — System operational status (no auth)
GET    /api/users/:id                  — Public profile
GET    /api/company                    — Company profile (no auth)
```

### User (JWT required — no XSRF needed)

```
GET    /api/users/me                   — My profile
PUT    /api/users/me                   — Update profile
POST   /api/users/me/change-password   — Change password
POST   /api/users/me/avatar            — Update profile picture (multipart upload to MinIO)
GET    /api/users/me/activity          — My activity log
GET    /api/users/me/dashboard         — Dashboard widgets
```

### Safe — /api (JWT + Idempotency, no XSRF)

Endpoint-endpoint ini **tidak berbahaya** dan tidak memerlukan XSRF token. Cukup JWT.

```
POST   /api/courses/:id/enroll           — Enroll ke course
PUT    /api/courses/:id/last-lesson       — Update lesson terakhir dipelajari
GET    /api/users/me/courses              — Course yang sudah dienroll
POST   /api/lessons/:id/start             — Start lesson
PUT    /api/lessons/:id/progress          — Update progress
POST   /api/lessons/:id/complete          — Complete lesson
GET    /api/users/me/certificates         — My certificates
GET    /api/users/me/certificates/:code   — Certificate by code
GET    /api/certificates/:code/verify     — Verify certificate (public, no auth)
GET    /api/users/me/notifications        — My notifications
GET    /api/users/me/notifications/count  — Unread notification count
PUT    /api/users/me/notifications/:id/read — Mark notif as read
PUT    /api/users/me/notifications/read-all — Mark all as read
DELETE /api/users/me/notifications/:id    — Delete notification
GET    /api/notifications/stream          — SSE real-time notification stream
GET    /api/users/me/discussions          — My discussions
GET    /api/users/me/replies              — My replies
GET    /api/users/me/study-cases          — My study cases
GET    /api/users/me/showcases            — My showcases
GET    /api/users/me/blogs                — My blogs (?status=draft|published)
GET    /api/levels                        — Level info
GET    /api/users/:id/points              — User points & rank
POST   /api/upload                        — File upload (multipart)
```

### Dangerous — /api (JWT + XSRF + Idempotency)

Endpoint-endpoint ini **mengubah konten** dan memerlukan XSRF token tambahan.

```
POST   /api/showcases                     — Create showcase
PUT    /api/showcases/:id                 — Update showcase
DELETE /api/showcases/:id                 — Delete showcase
POST   /api/showcases/:id/like            — Like showcase
DELETE /api/showcases/:id/like            — Unlike showcase
POST   /api/discussions                   — Create discussion
PUT    /api/discussions/:id               — Update discussion
DELETE /api/discussions/:id               — Delete discussion
POST   /api/discussions/:id/close         — Close discussion
POST   /api/discussions/:id/replies       — Create reply
PUT    /api/replies/:id                   — Update reply
DELETE /api/replies/:id                   — Delete reply
POST   /api/replies/:id/like              — Like reply
POST   /api/replies/:id/best              — Mark as best answer
POST   /api/reviews                       — Create review
PUT    /api/reviews/:id                   — Update review
DELETE /api/reviews/:id                   — Delete review
```

### Admin — /api/admin (JWT + XSRF + Idempotency + role=admin)

```
GET    /api/admin/dashboard                     — Admin dashboard
GET    /api/admin/users                         — All users (paginated)
GET    /api/admin/blogs                         — All blogs (semua status, paginated)
POST   /api/admin/blogs                         — Create blog
PUT    /api/admin/blogs/:id                     — Update any blog
DELETE /api/admin/blogs/:id                     — Delete any blog
POST   /api/admin/courses                       — Create course
PUT    /api/admin/courses/:id                   — Update course
DELETE /api/admin/courses/:id                   — Delete course
POST   /api/admin/courses/:course_id/sections    — Create section
PUT    /api/admin/sections/:section_id           — Update section
DELETE /api/admin/sections/:section_id           — Delete section
POST   /api/admin/lessons                       — Create lesson
POST   /api/admin/lessons/:id/details           — Add lesson detail
PUT    /api/admin/lessons/:id                   — Update lesson
DELETE /api/admin/lessons/:id                   — Delete lesson
POST   /api/admin/study-cases                   — Create study case
PUT    /api/admin/study-cases/:id               — Update study case
DELETE /api/admin/study-cases/:id               — Delete study case
POST   /api/admin/categories                    — Create category
GET    /api/admin/categories                    — All categories (admin view)
PUT    /api/admin/categories/:id                — Update category
DELETE /api/admin/categories/:id                — Delete category
```

### Company Profile (Public + Admin Update)

```
GET  /api/company (public)            — Get company profile (no auth)

Admin:
  PUT  /api/admin/company             — Update company profile
```

### FAQ (Public + Admin CRUD)

```
GET  /api/faqs (public)              — List active FAQs (?page=&limit=, no auth)

Admin:
  GET    /api/admin/faqs              — List all FAQs (paginated)
  GET    /api/admin/faqs/:id          — Get single FAQ
  POST   /api/admin/faqs              — Create FAQ
  PUT    /api/admin/faqs/:id          — Update FAQ
  DELETE /api/admin/faqs/:id          — Delete FAQ
```

**Total: 105 endpoint definitions** (+ change-password, avatar, 5 FAQ routes, company public + admin)

### Notification Flow

Notifikasi dikirim **otomatis** oleh backend setiap ada event yang membutuhkan pemberitahuan. Berikut daftar lengkap semua tipe notifikasi:

#### Daftar Tipe Notifikasi

| Tipe                 | Pemicu                                           | Diterima Oleh         | Link                       | File Service          |
| -------------------- | ------------------------------------------------ | --------------------- | -------------------------- | --------------------- |
| `new_reply`          | Reply baru di diskusi                            | Owner diskusi         | `/discussions/:id`         | `reply_service.go`    |
| `nested_reply`       | Balasan nested ke reply                          | Parent reply owner    | `/discussions/:id`         | `reply_service.go`    |
| `reply_like`         | Reply seseorang di-like                          | Creator reply         | `/discussions/:id`         | `reply_service.go`    |
| `best_answer`        | Reply ditandai sebagai jawaban terbaik           | Creator reply         | `/discussions/:id`         | `reply_service.go`    |
| `showcase_like`      | Showcase di-like                                 | Owner showcase        | `/showcases/:id`           | `service.go`          |
| `course_enrollment`  | User baru mendaftar ke kursus                    | Admin course          | `/courses/:id`             | `course_service.go`   |
| `enrollment_success` | Pendaftaran kursus berhasil                      | User yang mendaftar   | `/courses/:id`             | `course_service.go`   |
| `new_review`         | Review baru untuk kursus                         | Admin course          | `/courses/:id`             | `review_service.go`   |
| `lesson_completed`   | Pelajaran selesai + sertifikat didapat           | User yang belajar     | `/courses/:id/lessons/:slug` | `lesson_service.go`  |
| `level_up`           | Level naik (dengan badge dari user_levels)       | User yang naik level  | `/users/:id/points`        | `service.go`          |
| `blog_published`     | Blog diterbitkan                                 | Author blog           | `/blogs/:id`               | `blog_service.go`     |
| `discussion_created` | Diskusi baru dibuat (terkait lesson)             | Creator diskusi       | `/discussions/:id`         | `discussion_service.go` |

#### Anti Self-Notifikasi

Sistem **tidak** mengirim notifikasi jika aktor == penerima untuk mencegah spam notifikasi ke diri sendiri:

- Reply baru: dilewati jika replier == discussion owner
- Nested reply: dilewati jika replier == parent reply owner
- Like reply: dilewati jika liker == reply creator (self-like)
- Like showcase: dilewati jika liker == showcase owner (self-like)
- Review baru: dilewati jika reviewer == course admin (self-review)

#### Badge Level Up

Saat level naik, notifikasi `level_up` menyertakan badge dari tabel `user_levels`:

```json
{
  "type": "level_up",
  "title": "Level Naik! 🏆",
  "message": "Selamat! Anda naik ke level Expert ⭐ — Badge: Expert",
  "link": "/users/user123/points"
}
```

Badge diambil dari database `user_levels.badge_name` dan `user_levels.badge_icon`. Jika tidak ada data, fallback ke hardcoded (🌱 Beginner, 🌿 Intermediate, 🌳 Advanced, ⭐ Expert, 👑 Master).

---

## DTO (Data Transfer Object) Pattern

Semua response API menggunakan **DTO struct bukan raw domain model** untuk:

- Mengurangi ukuran payload (hanya field yang FE butuhkan)
- Menghilangkan field duplikat (misal `course_id` di section dalam response course detail)
- Null safety: field yang tidak ada data return `null`, bukan object kosong
- Backward compatibility dengan json tag yang konsisten

### Naming Convention

| Suffix      | Contoh                                                                                                                                                                       | Penggunaan                                               |
| ----------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------- |
| `*Brief`    | `CategoryBrief`, `UserBrief`, `LessonBrief`, `SectionBrief`                                                                                                                  | Nested object — hanya field minimal: id, name, slug, dll |
| `*Item`     | `CourseListItem`, `ReviewItem`, `DiscussionListItem`, `ShowcaseListItem`, `CertificateItem`, `LeaderboardItem`, `ActivityItem`, `MentorItem`, `UserListItem`, `BlogListItem` | Item dalam list/array response                           |
| `*Detail`   | `SectionDetail`, `CourseDetailWithSections`, `DiscussionDetail`, `ShowcaseDetail`, `StudyCaseDetail`, `BlogDetail`                                                           | Response detail endpoint tunggal                         |
| `*Response` | `LessonDetailResponse`                                                                                                                                                       | Response kompleks dengan multiple nested objects         |

### Key Response Changes

| Endpoint                               | Before (full domain)                                                                                                                                                        | After (DTO)                                                             |
| -------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------- |
| `GET /api/courses`                     | `category: {id,name,slug,desc,created_at,updated_at,deleted_at,...}`, `mentor: {id,email,password,name,avatar,bio,role,points,total_points,level,is_active,created_at,...}` | `category: {id,name,slug}`, `mentor: {id,name,avatar,role}` atau `null` |
| `GET /api/courses/:course_id`          | Sections contain `course_id` (duplicate) + full `Course` object nested                                                                                                      | Sections tanpa `course_id` duplikat, tanpa nested Course                |
| `GET /api/categories`                  | Full Category with all timestamps + relations                                                                                                                               | `[{id, name, slug}]` — hanya esensial                                   |
| `GET /api/categories/:slug`            | Courses in category with full User objects                                                                                                                                  | Courses in category with `CourseListItem` (brief)                       |
| `GET /api/sections`                    | Full Section with Course relationship                                                                                                                                       | `SectionDetail` — tanpa nested Course                                   |
| `GET /api/lessons/:id`                 | Raw `map[string]interface{}`                                                                                                                                                | `LessonDetailResponse` — typed struct                                   |
| All discussion/reply/cert/gamification | `[]map[string]interface{}`                                                                                                                                                  | Typed DTO structs                                                       |

### Daftar Lengkap DTO

**`CategoryBrief`** — nested di course, showcase, dll

```json
{ "id": "...", "name": "Backend", "slug": "backend" }
```

**`UserBrief`** — nested di course, showcase, study case, discussion

```json
{ "id": "...", "name": "Budi", "avatar": "...", "role": "mentor" }
```

**`LessonBrief`** — nested di section

```json
{ "id": "...", "title": "Pengenalan Go", "slug": "pengenalan-go", "difficulty": "beginner", "duration": 45, "order_index": 1, "video_url": "..." }
```

**`SectionBrief`** — nested di course detail

```json
{ "id": "...", "title": "Module 1", "description": "...",
  "order_index": 1, "lessons": [LessonBrief, ...] }
```

**`SectionDetail`** — response endpoint section detail

```json
{ "id": "...", "title": "...", "description": "...",
  "course_id": "...", "order_index": 1, "lessons": [LessonBrief, ...] }
```

**`CourseListItem`** — response list course

```json
{ "id": "...", "title": "Go REST API", "description": "...",
  "thumbnail": "...", "category": CategoryBrief, "admin_name": "Admin",
  "mentor": UserBrief|null, "hours": 10, "section_count": 3,
  "is_enrolled": true, "created_at": "..." }
```

**`CourseDetailWithSections`** — response course detail

```json
{ "id": "...", "title": "...", "description": "...",
  "thumbnail": "...", "category": CategoryBrief, "admin_id": "...",
  "admin_name": "Admin", "mentor": UserBrief|null, "hours": 10,
  "total_duration_hours": 3, "visibility": "public",
  "sections": [SectionBrief, ...], "is_enrolled": true, "created_at": "..." }
```

**`CategoryWithCourses`** — response category detail

```json
{ "id": "...", "name": "Backend", "slug": "backend",
  "description": "...", "courses": [CourseListItem, ...] }
```

**`LessonDetailResponse`** — response lesson detail

```json
{ "lesson": LessonBrief, "details": LessonDetail,
  "progress": LessonProgress|null, "next_lesson": LessonBrief|null,
  "section": SectionBrief|null, "course": CourseListItem|null }
```

**`ReviewItem`** — response review

```json
{ "id": "...", "user_id": "...", "user_name": "Budi", "user_avatar": "...", "course_id": "...", "lesson_id": "...", "rating": 5, "message": "Mantap!", "created_at": "..." }
```

**`DiscussionListItem`** — list discussion

```json
{ "id": "...", "title": "...", "user": { "id": "...", "name": "Budi", "avatar": "...", "role": "user" }, "lesson_id": null, "study_case_id": null, "status": "open", "view_count": 42, "created_at": "..." }
```

**`DiscussionDetail`** — discussion dengan replies

```json
{
  "id": "...",
  "title": "...",
  "content": "...",
  "user": { "id": "...", "name": "Budi", "avatar": "...", "role": "user" },
  "status": "open",
  "view_count": 42,
  "created_at": "...",
  "replies": [{ "id": "...", "content": "...", "user": { "id": "...", "name": "Nama", "avatar": "...", "role": "user" }, "likes": 5, "is_best": false, "created_at": "..." }]
}
```

**`ReplyListItem`** — list reply milik user

```json
{ "id": "...", "content": "...", "discussion_id": "...", "discussion_title": "Judul Diskusi", "parent_id": null, "likes_count": 3, "is_marked_best": false, "created_at": "..." }
```

**`ShowcaseListItem`** — list showcase

```json
{ "id": "...", "title": "...", "description": "...",
  "media_urls": ["..."], "likes_count": 10, "views_count": 50,
  "visibility": "public", "user": UserBrief,
  "category": CategoryBrief, "created_at": "..." }
```

**`ShowcaseDetail`** — showcase detail

```json
{ "id": "...", "title": "...", "description": "...",
  "media_urls": ["..."], "likes_count": 10, "views_count": 50,
  "user": UserBrief, "category": CategoryBrief,
  "is_liked_by_me": true, "created_at": "..." }
```

**`CertificateItem`** — list certificate

```json
{ "id": "...", "unique_code": "CERT-abc12345", "issued_at": "...", "lesson_id": "...", "lesson_name": "Nama Lesson", "user_name": "Budi", "achievement": { "type": "certificate" }, "verification_url": "https://jvalleyverse.com/certificates/verify/CERT-abc12345", "qr_code_url": "https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=..." }
```

**`UserListItem`** — admin: list all users

```json
{ "id": "...", "email": "user@mail.com", "name": "Budi", "avatar": "...", "role": "user", "level": 3, "total_points": 500, "is_active": true, "created_at": "..." }
```

**`MentorItem`** — list mentors

```json
{ "id": "...", "name": "Budi Mentor", "avatar": "...", "bio": "...", "level": 5, "total_points": 2500 }
```

**`LeaderboardItem`** — leaderboard

```json
{ "rank": 1, "user_id": "...", "name": "Budi", "avatar": "...", "total_points": 2500, "level": 5 }
```

**`ActivityItem`** — activity log

```json
{ "id": "...", "activity": "complete_lesson", "points": 50, "timestamp": "..." }
```

**`UserStats`** — user points/level

```json
{ "user_id": "...", "name": "Budi", "total_points": 500,
  "current_level": 3, "recent_activity": [ActivityItem, ...] }
```

**`LevelInfo`** — level definitions

```json
{ "name": "Beginner", "threshold": 0, "color": "#6366f1", "description": "Just starting your journey" }
```

**`BlogListItem`** — list blog

```json
{
  "id": "...",
  "title": "Belajar Go",
  "slug": "belajar-go",
  "description": "...",
  "cover_img_url": "https://cdn...",
  "tags": ["golang", "backend"],
  "status": "published",
  "author": { "id": "...", "name": "Admin", "avatar": "..." },
  "category": { "id": "...", "name": "Backend", "slug": "backend" },
  "created_at": "2026-06-15T12:00:00Z"
}
```

**`BlogDetail`** — blog detail

```json
{
  "id": "...",
  "title": "Belajar Go",
  "slug": "belajar-go",
  "description": "...",
  "content": "# Full markdown...",
  "cover_img_url": "https://cdn...",
  "tags": ["golang"],
  "status": "published",
  "author": { "id": "...", "name": "Admin", "avatar": "..." },
  "category": { "id": "...", "name": "Backend", "slug": "backend" },
  "created_at": "2026-06-15T12:00:00Z",
  "updated_at": "2026-06-15T14:00:00Z"
}
```

### Real-time Notifications (SSE)

**Endpoint:** `GET /api/notifications/stream` (JWT required)

Notifikasi real-time dikirim melalui **Server-Sent Events (SSE)**. Setelah user login dan membuka koneksi SSE, setiap notifikasi baru yang dibuat oleh sistem akan langsung dikirim ke user tanpa perlu polling.

**Cara kerja:**

- Client membuka koneksi `GET /api/notifications/stream` dengan JWT
- Server menjaga koneksi tetap terbuka (long-lived HTTP connection)
- Setiap ada notifikasi baru, server push event ke client
- Heartbeat setiap 30 detik untuk menjaga koneksi tetap hidup
- Client disconnect akan auto-cleanup koneksi di server

**SSE Event Format:**

```
event: connected
data: {"status":"connected","user_id":"..."}

event: notification
data: {"type":"new_reply","title":"Reply Baru","message":"...","link":"/discussions/...","payload":{"id":"...","created_at":"...","is_read":false}}

: heartbeat
```

**Implementasi Teknis:**

- File: `internal/service/notification_hub.go` — Singleton hub manage koneksi per user
- File: `internal/handler/sse_handler.go` — SSE endpoint handler
- Hub menyimpan map `userID → {connID → channel}` dengan `sync.RWMutex`
- `CreateNotification()` setelah simpan ke DB, publish ke hub
- Detail teknis: lihat `sse_notification.md`

**Frontend JavaScript:**

```js
const evtSource = new EventSource("/api/notifications/stream", {
  headers: { Authorization: "Bearer " + token },
});

evtSource.addEventListener("notification", (e) => {
  const notif = JSON.parse(e.data);
  // tampilkan toast, update badge, dll
});

evtSource.onerror = () => {
  // auto-reconnect by default
};
```

### Dashboard Widget

```
GET /api/users/me/dashboard
Response: {
  "courses_in_progress": 2,    // LessonProgress status "in_progress"
  "courses_completed": 3,      // LessonProgress status "completed"
  "courses_dropped": 1,        // LessonProgress status "started"
  "unread_notifications": 5    // Notification is_read = false
}
```

### Course Enrollment

**Model:** `course_enrollments` — tabel baru untuk tracking enrollment user ke course.

| Field           | Tipe   | Keterangan                                |
| --------------- | ------ | ----------------------------------------- |
| id              | string | Primary key (CUID)                        |
| user_id         | string | Foreign key ke users (unique composite)   |
| course_id       | string | Foreign key ke courses (unique composite) |
| last_lesson_id  | string | nullable - lesson terakhir dipelajari     |
| original_price  | float  | harga asli sebelum diskon                 |
| discount_amount | float  | jumlah diskon                             |
| discount_code   | string | kode diskon                               |

**Endpoint Enrollment:**

- `POST /api/courses/:id/enroll` — Enroll user ke course (JWT required)
- `GET /api/users/me/courses` — Daftar course yang sudah dienroll (paginated, includes last_lesson_id)
- `PUT /api/courses/:id/last-lesson` — Update lesson terakhir dipelajari (body: { lesson_id: "..." })

**Field `is_enrolled` di response course:**

- `GET /api/courses` — Setiap course memiliki field `is_enrolled` **hanya jika** user login (bawa JWT)
- `GET /api/courses/:course_id` — Response menyertakan `is_enrolled` **hanya jika** user login
- `GET /api/categories/:category_id/courses` — Setiap course memiliki `is_enrolled` **hanya jika** user login

Jika user **tidak login** (tanpa JWT), field `is_enrolled` **tidak muncul sama sekali** di response.

### Course Update

Course model kini memiliki field **`hours`** (int) — durasi total course yang diisi admin saat create.

- `GET /api/courses/:course_id` juga mengembalikan `total_duration_hours` (otomatis dari jumlah durasi semua lesson / 60)

---

## Idempotency Key (Safe Retry)

Semua mutation endpoint (POST/PUT/DELETE) di protected `/api` group dan auth endpoints support **Idempotency-Key** untuk mencegah duplikasi akibat network retry.

### Cara Kerja

1. Client generate UUID v4 (`550e8400-e29b-41d4-a716-446655440000`)
2. Kirim header `Idempotency-Key: <uuid>` di request POST/PUT/DELETE
3. **Request pertama:** Server proses normally, cache response di Redis (TTL 24 jam)
4. **Retry (request kedua+):** Server deteksi key sudah ada di Redis → return response yang sama **tanpa memproses ulang**
5. Jika Redis down → middleware silent pass-through, API tetap jalan normal

### Response Headers

| Header                   | Value  | Keterangan                                                         |
| ------------------------ | ------ | ------------------------------------------------------------------ |
| `X-Idempotency-Replayed` | `true` | Hanya muncul jika response berasal dari cache (bukan request asli) |

### Contoh Request

```bash
# Generate UUID
$ python -c "import uuid; print(uuid.uuid4())"
550e8400-e29b-41d4-a716-446655440000

# Request pertama (proses normal)
curl -X POST http://localhost:3000/api/showcases \
  -H "Authorization: Bearer <token>" \
  -H "Idempotency-Key: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{"title":"My Project","category_id":"..."}'

# Retry (dapat response dari cache, X-Idempotency-Replayed: true)
curl -X POST http://localhost:3000/api/showcases \
  -H "Authorization: Bearer <token>" \
  -H "Idempotency-Key: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{"title":"My Project","category_id":"..."}'
```

### Error Response (Invalid Key)

```json
{
  "error": "Idempotency-Key must be a valid UUID v4 (e.g. 550e8400-e29b-41d4-a716-446655440000)"
}
```

### Endpoints dengan Idempotency

Semua **POST/PUT/DELETE** di bawah endpoint berikut:

**Auth:**

- `POST /api/auth/register`
- `POST /api/auth/login`
- `POST /api/auth/refresh`

**Protected (/api):**

- Semua mutation: create/update/delete showcases, discussions, replies, reviews, enrollments, certificates, progress tracking

**Admin (/api/admin):**

- Semua mutation: create/update/delete courses, sections, lessons, study cases, categories

### Notes

- Idempotency-Key **tidak wajib** — hanya untuk client yang ingin safe retry
- Key **wajib UUID v4** — format ketat, bukan string bebas
- Cache hanya untuk **2xx response** — error 4xx/5xx tidak di-cache
- TTL cache **24 jam** — cocok untuk retry dalam timeframe yang masuk akal

---

## Security & Rate Limiting

### Rate Limit Tiers

API memiliki tiga tier rate limiting berdasarkan jenis endpoint dan kerentanan:

| Tier        | Limit              | Target                 | Endpoint                                                         |
| ----------- | ------------------ | ---------------------- | ---------------------------------------------------------------- |
| **Global**  | **200 req/min/IP** | General browsing       | Semua route (baseline)                                           |
| **Content** | **60 req/min/IP**  | Anti-scraping          | Public GET: courses, lessons, showcases, categories, study-cases |
| **Auth**    | **10 req/min/IP**  | Brute force protection | `POST /api/auth/login`, `POST /api/auth/register`                |

Jika limit terlampaui, server mengembalikan:

```json
{
  "error": "Too many requests",
  "retry_after": 60
}
```

### Prometheus Metrics & Monitoring

**Endpoint:** `GET /metrics` (public, no auth)

Aplikasi mengekspos metrik Prometheus melalui endpoint `/metrics` menggunakan package `fiberprometheus`. Metrik ini di-scrape oleh Prometheus server setiap 15 detik.

**Metrik yang tersedia:**
- `http_requests_total` — Total request per method, status, path
- `http_request_duration_seconds` — Histogram durasi request
- `http_requests_in_progress` — Request aktif saat ini
- `go_goroutines` — Jumlah goroutine
- `go_memstats_heap_alloc_bytes` — Alokasi heap
- `go_gc_duration_seconds` — Durasi garbage collection
- `process_start_time_seconds` — Waktu mulai aplikasi

**Monitoring Stack (Grafana + Prometheus + Node Exporter):**

| Komponen   | Endpoint              | Port   |
| ---------- | --------------------- | ------ |
| Prometheus | http://localhost:9090 | `:9090` |
| Grafana    | http://localhost:3000 | `:3000` |
| Node Exp.  | http://localhost:9100 | `:9100` |

Konfigurasi dan dashboard tersedia di:
- `deploy/monitoring/prometheus.yml` — Konfigurasi scrape targets
- `deploy/monitoring/grafana-dashboard.json` — Dashboard JSON (import via Grafana UI)
- `scripts/setup-monitoring.sh` — Script instalasi otomatis untuk VPS

### Anti-Scraping (ScraperGuard)

Semua public content endpoint (`/api/courses/*`, `/api/lessons/*`, `/api/showcases/*`, `/api/categories/*`, `/api/study-cases/*`) dilindungi oleh **ScraperGuard** — middleware yang memblokir request berdasarkan User-Agent.

**Diblokir (403 Forbidden):**

- Empty User-Agent
- CLI tools: `curl`, `wget`, `libcurl`
- Python scrapers: `python-requests`, `aiohttp`, `scrapy`, `httpx`
- API clients: `PostmanRuntime`, `insomnia`, `HttpClient`, `okhttp`
- Bahasa pemrograman: `Java/`, `ruby`, `faraday`
- Generic bots: `bot`, `spider`, `crawler` (kecuali search engine)

**Diizinkan:**

- Search engine: `Googlebot`, `Bingbot`, `Slurp` (Yahoo), `DuckDuckBot`, `Baiduspider`, `YandexBot`, `Sogou`
- Social media: `facebookexternalhit`, `Twitterbot`, `LinkedInBot`, `WhatsApp`, `TelegramBot`, `Discordbot`

**Response ScraperGuard:**

```json
{
  "error": "Access denied: automated scraping is not allowed"
}
```

### Security Headers

Semua response dilengkapi security headers berikut:

| Header                      | Value                                      | Fungsi                                    |
| --------------------------- | ------------------------------------------ | ----------------------------------------- |
| `X-Content-Type-Options`    | `nosniff`                                  | Cegah MIME type sniffing                  |
| `X-Frame-Options`           | `DENY`                                     | Cegah clickjacking (tidak bisa di-iframe) |
| `X-XSS-Protection`          | `1; mode=block`                            | XSS filter (browser lama)                 |
| `Referrer-Policy`           | `strict-origin-when-cross-origin`          | Batasi informasi referer                  |
| `Permissions-Policy`        | `geolocation=(), microphone=(), camera=()` | Nonaktifkan API sensitif                  |
| `Strict-Transport-Security` | `max-age=31536000; includeSubDomains`      | Paksa HTTPS (setelah deploy)              |

### XSRF Protection

**XSRF hanya diterapkan ke endpoint berbahaya** (bukan semua endpoint).

| Group        | Middleware                       | Endpoints                                                                 |
| ------------ | -------------------------------- | ------------------------------------------------------------------------- |
| **Safe**     | JWT + Idempotency                | Enrollment, progress, notifications, upload, certificates, my-items, dll |
| **Dangerous**| JWT + XSRF + Idempotency         | Showcases (CRUD/like), Discussions (CRUD/close), Replies, Reviews         |
| **Admin**    | JWT + XSRF + Idempotency + admin | Admin CRUD: courses, sections, lessons, study-cases, categories, blogs    |

Endpoint **Safe** (termasuk `POST /api/courses/:id/enroll`, `POST /api/upload`) **tidak perlu XSRF token** — cukup JWT.

### Strategi: Double-Submit Cookie + Origin/Referer Fallback

Endpoint **Dangerous** dan **Admin** dilindungi oleh **2 lapis validasi berurutan**:

1. **Lapis 1 — Cookie Match:** Header `X-XSRF-TOKEN` harus sama dengan cookie `XSRF-TOKEN`. Jika cocok → izinkan.
2. **Lapis 2 — Origin/Referer Fallback:** Jika cookie check gagal, periksa header `Origin`. Jika tidak ada, pakai `Referer`. Jika origin termasuk dalam daftar CORS yang diizinkan → izinkan.
3. **Blokir (403):** Jika kedua lapis gagal.

Ini memastikan SPA client yang tidak bisa membaca XSRF cookie tetap aman karena browser tidak mengizinkan JavaScript dari domain lain memalsukan Origin header.

**Urutan prioritas:** `X-XSRF-TOKEN` header == cookie > Origin match > Referer match > 403 Forbidden.

**Cookie Attributes:**
- `SameSite: "None"` — mengizinkan pengiriman cookie dari frontend ke API pada cross-origin POST request
- `Secure: true` — cookie hanya dikirim via HTTPS (wajib untuk SameSite=None)
- `HTTPOnly: false` — JavaScript frontend dapat membaca cookie

```
Cookie: XSRF-TOKEN=abc123
Header: X-XSRF-TOKEN: abc123

# Atau tanpa cookie — cukup Origin header:
Origin: https://jvalleyverse.web.id
```

### CORS

Origin yang diizinkan dikonfigurasi via environment variable `CORS_ORIGINS`. Default:

```
http://localhost:3000, http://localhost:5173, https://jvalleyverse.vercel.app
```

### Idempotency Key

(lihat section di atas)

## Seed Data

### Default Credentials

| Email                    | Password   | Role   |
| ------------------------ | ---------- | ------ |
| admin@jvalleyverse.com   | Admin@123  | admin  |
| mentor1@jvalleyverse.com | Mentor@123 | mentor |
| mentor2@jvalleyverse.com | Mentor@123 | mentor |
| budi@example.com         | User@123   | user   |
| siti@example.com         | User@123   | user   |
| andi@example.com         | User@123   | user   |
| dewi@example.com         | User@123   | user   |

### Sample Learning Paths

**Go REST API Course** (3 lessons):

1. Pengenalan Go & Setup Environment → beginner, 45min
2. Fiber Framework & Routing → beginner, 60min
3. GORM & Database Integration → intermediate, 90min

**React Dashboard Course** (2 lessons):

1. Setup React + TypeScript Project → beginner, 30min
2. Membuat Komponen Dashboard → intermediate, 75min

**Flutter E-Commerce Course** (1 lesson):

1. Flutter Fundamentals & Dart Basics → beginner, 60min

---

## File Upload (MinIO)

Semua upload file (thumbnail course/lesson, avatar user, media showcase, img study case) dilakukan via **endpoint upload terpusat**, bukan langsung di handler masing-masing.

### Arsitektur

```
[Frontend]                          [Backend (Go/Fiber)]              [MinIO S3]
    |                                     |                              |
    | 1. POST /api/upload (multipart)    |                              |
    |    (file + folder)                  |                              |
    |------------------------------------>|                              |
    |                                     | 2. minioClient.PutObject()   |
    |                                     |----------------------------->|
    |                                     |                              |
    |  3. { url: "https://cdn.../uuid" }  |                              |
    |<------------------------------------|                              |
    |                                     |                              |
    | 4. POST /api/admin/courses          |  (JSON biasa, hanya URL)    |
    |    { thumbnail: "<cdn-url>" }       |                              |
    |------------------------------------>|                              |
```

### Endpoint Upload

```
POST /api/upload (JWT + Idempotency — no XSRF)
```

**Request (multipart/form-data):**

| Field    | Type   | Required | Keterangan                                                                 |
| -------- | ------ | -------- | -------------------------------------------------------------------------- |
| `file`   | File   | ✅       | File yang diupload (max 10MB)                                              |
| `folder` | String | ✅       | Folder tujuan: `courses`, `lessons`, `avatars`, `showcases`, `study-cases` |

**Validasi folder:** hanya alfanumerik + `-` + `_` (no path traversal, no nested folders)

**Validasi ekstensi:**

```
Images: jpg, jpeg, png, gif, webp, svg
Video:  mp4
Docs:   pdf
Archive: zip
```

**Response (201 Created):**

```json
{
  "url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/courses/abc123-def456.jpg",
  "object_name": "courses/abc123-def456.jpg",
  "size": 245760,
  "content_type": "image/jpeg"
}
```

**Error Responses:**

| Status | Error                     | Penyebab                        |
| ------ | ------------------------- | ------------------------------- |
| 400    | `"file is required"`      | Tidak ada file di field `file`  |
| 400    | `"folder is required"`    | Field `folder` kosong           |
| 400    | `"invalid folder name"`   | Folder mengandung `..` atau `/` |
| 400    | `"file type not allowed"` | Ekstensi tidak diizinkan        |
| 400    | `"file too large"`        | Melebihi 10MB                   |
| 500    | `"failed to upload file"` | Gagal upload ke MinIO           |

### Folder Convention

| Folder        | Entity    | Field Target    |
| ------------- | --------- | --------------- |
| `courses`     | Course    | `thumbnail`     |
| `lessons`     | Lesson    | `thumbnail`     |
| `avatars`     | User      | `avatar`        |
| `showcases`   | Showcase  | `media_urls[]`  |
| `blogs`       | Blog      | `cover_img_url` |
| `study-cases` | StudyCase | `img_url`       |

### MinIO Client

**File:** `internal/minio/client.go`

| Function                                                    | Deskripsi                                                |
| ----------------------------------------------------------- | -------------------------------------------------------- |
| `ConnectMinio()`                                            | Init client + auto-create bucket                         |
| `UploadFile(ctx, file, folder, filename)`                   | Upload file ke folder/, return object_name               |
| `DeleteFile(ctx, objectName)`                               | Hapus file dari bucket                                   |
| `GeneratePresignedUploadURL(ctx, folder, filename, expiry)` | Generate presigned URL untuk upload dari client langsung |
| `IsAvailable()`                                             | Cek apakah MinIO client siap dipakai                     |

### Contoh Frontend

```js
const formData = new FormData();
formData.append("file", fileInput.files[0]);
formData.append("folder", "courses");

fetch("/api/upload", {
  method: "POST",
  headers: {
    Authorization: "Bearer <token>",
    "X-XSRF-TOKEN": "<xsrf-token>",
    "Idempotency-Key": crypto.randomUUID(),
  },
  body: formData,
})
  .then((res) => res.json())
  .then((data) => {
    // data.url = "https://cdn.../courses/uuid.jpg"
    // Simpan URL ini ke field thumbnail/avatar/media_urls
  });
```

---

## User Level Definitions

| Level | Name         | Min Points | Max Points | Badge |
| ----- | ------------ | ---------- | ---------- | ----- |
| 1     | Beginner     | 0          | 99         | 🌱    |
| 2     | Intermediate | 100        | 499        | 🌿    |
| 3     | Advanced     | 500        | 999        | 🌳    |
| 4     | Expert       | 1000       | 1999       | ⭐    |
| 5     | Master       | 2000       | ∞          | 👑    |

---

## Environmental Config

| Variable             | Default                                                                     | Description                                 |
| -------------------- | --------------------------------------------------------------------------- | ------------------------------------------- |
| PORT                 | 3000                                                                        | HTTP listen port                            |
| JWT_SECRET           | supersecretkey                                                              | JWT signing key                             |
| JWT_EXPIRY           | 24h                                                                         | Token expiry                                |
| DB_HOST              | localhost                                                                   | PostgreSQL host                             |
| DB_PORT              | 5432                                                                        | PostgreSQL port                             |
| DB_USER              | postgres                                                                    | DB user                                     |
| DB_PASSWORD          | root                                                                        | DB password                                 |
| DB_NAME              | jvalleyverse                                                                | DB name                                     |
| REDIS_HOST           | localhost:6379                                                              | Redis host:port                             |
| CORS_ORIGINS         | http://localhost:3000,http://localhost:5173,https://jvalleyverse.vercel.app | Allowed origins                             |
| ADMIN_EMAIL          | admin@jvalleyverse.com                                                      | Default admin                               |
| ADMIN_PASSWORD       | admin123                                                                    | Default admin pass                          |
| GOOGLE_CLIENT_ID     | —                                                                           | Google OAuth client ID (for Google One Tap) |
| **MINIO_ENDPOINT**   | `minio.mohagussetiaone.my.id`                                               | MinIO server endpoint                       |
| **MINIO_ACCESS_KEY** | (required)                                                                  | MinIO access key                            |
| **MINIO_SECRET_KEY** | (required)                                                                  | MinIO secret key                            |
| **MINIO_BUCKET**     | `jvalleyverse`                                                              | MinIO bucket name                           |
| **MINIO_CDN_URL**    | `https://cdn.mohagussetiaone.my.id`                                         | Public CDN base URL                         |
| **MINIO_USE_SSL**    | `true`                                                                      | Use HTTPS for MinIO connection              |

## DB Reset

```powershell
# 1. Drop schema + recreate
psql -h localhost -U postgres -d jvalleyverse -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public; GRANT ALL ON SCHEMA public TO postgres; GRANT ALL ON SCHEMA public TO public;"

# 2. Run seed (migrate + seed data)
go run ./cmd/seed/main.go

# 3. Start server
go run ./cmd/api/main.go
```
