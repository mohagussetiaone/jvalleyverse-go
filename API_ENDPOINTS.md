# API Endpoints & Request/Response Examples

## 📋 Quick Reference

```
┌─────────────────────────────────────────────────────────────┐
│                    COMPLETE API MAP                         │
├─────────────────────────────────────────────────────────────┤
│ AUTH              │ USER              │ PROJECTS (ADMIN)    │
│ POST /login       │ GET /me           │ POST /admin/projects│
│ POST /register    │ PUT /me           │ GET /admin/projects │
│ POST /refresh     │ GET /leaderboard  │ PUT /admin/:id      │
│ POST /logout      │ GET /users/:id    │ DELETE /admin/:id   │
├─────────────────────────────────────────────────────────────┤
│ CLASSES           │ CERTIFICATES      │ DISCUSSIONS         │
│ GET /classes      │ POST /generate    │ POST /discussions   │
│ GET /classes/:id  │ GET /my-certs     │ GET /discussions    │
│ PUT /admin/:id    │ GET /:code        │ GET /:id            │
│                   │                   │ PUT /:id            │
├─────────────────────────────────────────────────────────────┤
│ REPLIES           │ SHOWCASES         │ GAMIFICATION        │
│ POST /replies     │ POST /showcases   │ GET /leaderboard    │
│ GET /replies      │ GET /showcases    │ GET /users/:id/pts  │
│ POST /:id/like    │ GET /:id          │ GET /levels         │
│ DELETE /:id       │ POST /:id/like    │ GET /my/activity    │
└─────────────────────────────────────────────────────────────┘
```

## 🔐 Authentication Endpoints

### POST /api/v1/auth/register

Create new account

**Request:**

```json
{
  "email": "user@example.com",
  "password": "secure_password",
  "name": "John Developer"
}
```

**Response (201):**

```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Developer",
  "role": "user",
  "points": 0,
  "level": 1,
  "created_at": "2024-01-20T10:00:00Z"
}
```

### POST /api/v1/auth/login

Authenticate user

**Request:**

```json
{
  "email": "user@example.com",
  "password": "secure_password"
}
```

**Response (200):**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Developer",
    "role": "user",
    "points": 50,
    "level": 1
  }
}
```

---

## 👤 User Endpoints

### GET /api/v1/users/me

Get current user profile

**Headers:** `Authorization: Bearer {access_token}`

**Response (200):**

```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Developer",
  "avatar": "https://...",
  "bio": "Full-stack developer",
  "role": "user",
  "points": 85,
  "total_points": 85,
  "level": 1,
  "is_active": true,
  "created_at": "2024-01-20T10:00:00Z"
}
```

### PUT /api/v1/users/me

Update profile

**Request:**

```json
{
  "name": "John Developer",
  "bio": "Senior Full-stack Developer",
  "avatar": "https://..."
}
```

**Response (200):** Updated user object

### GET /api/v1/leaderboard?page=1&limit=50

Get top users by points

**Query Params:**

- `page=1` (default: 1)
- `limit=50` (default: 50, max: 100)

**Response (200):**

```json
{
  "data": [
    {
      "rank": 1,
      "user_id": 5,
      "name": "Sarah Expert",
      "points": 2500,
      "level": 5,
      "badge": "master.png"
    },
    {
      "rank": 2,
      "user_id": 10,
      "name": "Mike Developer",
      "points": 2150,
      "level": 5,
      "badge": "expert.png"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 1523
  }
}
```

### GET /api/v1/users/:id

Get public user profile

**Response (200):**

```json
{
  "id": 5,
  "name": "Sarah Expert",
  "avatar": "https://...",
  "bio": "Data Scientist",
  "role": "user",
  "points": 2500,
  "level": 5,
  "showcase_count": 15,
  "certificate_count": 8
}
```

---

## 📌 Project Endpoints (Admin Only)

### POST /api/v1/admin/projects

Create new project

**Headers:** `Authorization: Bearer {admin_token}`

**Request:**

```json
{
  "title": "Advanced Go Programming",
  "description": "Master backend development with Go",
  "thumbnail": "https://...",
  "category_id": 1,
  "visibility": "public"
}
```

**Response (201):**

```json
{
  "id": 1,
  "title": "Advanced Go Programming",
  "description": "Master backend development with Go",
  "thumbnail": "https://...",
  "category_id": 1,
  "admin_id": 1,
  "visibility": "public",
  "created_at": "2024-01-20T10:00:00Z"
}
```

### GET /api/v1/admin/projects

List all projects (admin view with stats)

**Headers:** `Authorization: Bearer {admin_token}`

**Response (200):**

```json
{
  "data": [
    {
      "id": 1,
      "title": "Advanced Go Programming",
      "category": { "id": 1, "name": "Backend" },
      "admin": { "id": 1, "name": "Admin User" },
      "class_count": 5,
      "enrollment_count": 234,
      "created_at": "2024-01-20T10:00:00Z"
    }
  ],
  "pagination": { "page": 1, "limit": 20, "total": 12 }
}
```

### PUT /api/v1/admin/projects/:id

Update project

**Request:**

```json
{
  "title": "Advanced Go (Updated)",
  "description": "...",
  "visibility": "public"
}
```

### DELETE /api/v1/admin/projects/:id

Delete project (cascade delete classes)

**Response (204):** No content

---

## 📚 Class Endpoints

### POST /api/v1/admin/classes

Create class under project (Admin)

**Headers:** `Authorization: Bearer {admin_token}`

**Request:**

```json
{
  "title": "Module 1: Goroutines & Concurrency",
  "description": "Learn how to use goroutines",
  "content": "# Goroutines\n## Introduction\n...",
  "thumbnail": "https://...",
  "project_id": 1,
  "category_id": 1,
  "difficulty": "intermediate",
  "duration": 120,
  "order": 1
}
```

**Response (201):**

```json
{
  "id": 10,
  "title": "Module 1: Goroutines & Concurrency",
  "project_id": 1,
  "category_id": 1,
  "admin_id": 1,
  "difficulty": "intermediate",
  "order": 1,
  "created_at": "2024-01-20T10:00:00Z"
}
```

### GET /api/v1/projects/:id/classes

List classes in a project

**Response (200):**

```json
{
  "data": [
    {
      "id": 10,
      "title": "Module 1: Goroutines & Concurrency",
      "description": "Learn how to use goroutines",
      "difficulty": "intermediate",
      "duration": 120,
      "order": 1
    }
  ]
}
```

### GET /api/v1/classes/:id

Get full class content

**Response (200):**

```json
{
  "id": 10,
  "title": "Module 1: Goroutines & Concurrency",
  "content": "# Goroutines\n## Introduction\n...",
  "project": { "id": 1, "title": "Advanced Go Programming" },
  "admin": { "id": 1, "name": "Admin User" },
  "difficulty": "intermediate",
  "duration": 120
}
```

---

## 🏆 Certificate Endpoints (Private)

### POST /api/v1/certificates/generate

Issue certificate to user (System/Admin)

**Headers:** `Authorization: Bearer {admin_token}`

**Request:**

```json
{
  "user_id": 5,
  "class_id": 10
}
```

**Response (201):**

```json
{
  "id": 42,
  "user_id": 5,
  "class_id": 10,
  "unique_code": "CERT-2024-ABC123XYZ",
  "badge_url": "https://...",
  "issued_at": "2024-01-20T10:00:00Z",
  "expires_at": null
}
```

### GET /api/v1/certificates

Get current user's certificates

**Headers:** `Authorization: Bearer {user_token}`

**Response (200):**

```json
{
  "data": [
    {
      "id": 42,
      "unique_code": "CERT-2024-ABC123XYZ",
      "class": { "id": 10, "title": "Module 1: ..." },
      "issued_at": "2024-01-20T10:00:00Z"
    }
  ]
}
```

### GET /api/v1/certificates/:code

View specific certificate (Owner + Admin only)

**Headers:** `Authorization: Bearer {user_token}`

**Response (200):**

```json
{
  "id": 42,
  "user_id": 5,
  "class_id": 10,
  "unique_code": "CERT-2024-ABC123XYZ",
  "badge_url": "https://...",
  "class": {
    "id": 10,
    "title": "Module 1: Goroutines & Concurrency"
  },
  "issued_at": "2024-01-20T10:00:00Z"
}
```

**Response (403 - if not owner):**

```json
{
  "error": "Unauthorized: Certificate not accessible"
}
```

---

## 💬 Discussion Endpoints

### POST /api/v1/discussions

Create discussion

**Headers:** `Authorization: Bearer {user_token}`

**Request:**

```json
{
  "title": "How to optimize database queries?",
  "content": "I have slow queries in production...",
  "class_id": 10,
  "category_id": 1
}
```

**Response (201):**

```json
{
  "id": 100,
  "title": "How to optimize database queries?",
  "user": { "id": 5, "name": "John Dev" },
  "class_id": 10,
  "category_id": 1,
  "status": "open",
  "views_count": 0,
  "created_at": "2024-01-20T10:00:00Z"
}
```

### GET /api/v1/discussions?page=1&class_id=10

List discussions

**Query Params:**

- `page=1`
- `class_id=10` (optional)
- `category_id=1` (optional)
- `status=open` (optional)

**Response (200):**

```json
{
  "data": [
    {
      "id": 100,
      "title": "How to optimize database queries?",
      "user": { "id": 5, "name": "John Dev" },
      "replies_count": 3,
      "views_count": 42,
      "is_pinned": false,
      "created_at": "2024-01-20T10:00:00Z"
    }
  ],
  "pagination": { "page": 1, "limit": 20, "total": 156 }
}
```

### GET /api/v1/discussions/:id

Get discussion with replies (threaded)

**Response (200):**

```json
{
  "id": 100,
  "title": "How to optimize database queries?",
  "content": "I have slow queries in production...",
  "user": { "id": 5, "name": "John Dev", "level": 3 },
  "views_count": 42,
  "replies": [
    {
      "id": 200,
      "content": "Use indexes!",
      "user": { "id": 10, "name": "Sarah" },
      "parent_id": null,
      "likes_count": 5,
      "is_marked_best": true,
      "child_replies": [
        {
          "id": 201,
          "content": "Especially on foreign keys",
          "user": { "id": 11, "name": "Mike" },
          "parent_id": 200,
          "likes_count": 2
        }
      ]
    }
  ]
}
```

---

## ↩️ Reply Endpoints (Nested)

### POST /api/v1/discussions/:id/replies

Reply to discussion (top-level)

**Headers:** `Authorization: Bearer {user_token}`

**Request:**

```json
{
  "content": "Try adding indexes to foreign keys!"
}
```

**Response (201):**

```json
{
  "id": 200,
  "content": "Try adding indexes to foreign keys!",
  "user": { "id": 10, "name": "Sarah" },
  "discussion_id": 100,
  "parent_id": null,
  "likes_count": 0,
  "created_at": "2024-01-20T10:00:00Z"
}
```

### POST /api/v1/replies/:id/replies

Reply to a reply (nested/threaded)

**Request:**

```json
{
  "content": "Especially on foreign keys, yes!"
}
```

**Response (201):**

```json
{
  "id": 201,
  "content": "Especially on foreign keys, yes!",
  "user": { "id": 11, "name": "Mike" },
  "discussion_id": 100,
  "parent_id": 200,
  "likes_count": 0,
  "created_at": "2024-01-20T10:00:00Z"
}
```

### PUT /api/v1/replies/:id

Edit reply (owner only)

**Request:**

```json
{
  "content": "Updated content..."
}
```

### DELETE /api/v1/replies/:id

Delete reply (owner or admin)

**Response (204):** No content

---

## 🎨 Showcase Endpoints

### POST /api/v1/showcases

Create showcase

**Headers:** `Authorization: Bearer {user_token}`

**Request:**

```json
{
  "title": "My ML Model - 95% Accuracy",
  "description": "Trained CNN classifier on CIFAR-10",
  "media_urls": ["https://cdn.example.com/project1.jpg", "https://cdn.example.com/project2.jpg"],
  "category_id": 2,
  "visibility": "public"
}
```

**Response (201):**

```json
{
  "id": 42,
  "title": "My ML Model - 95% Accuracy",
  "user": { "id": 5, "name": "John Dev" },
  "category_id": 2,
  "status": "published",
  "visibility": "public",
  "likes_count": 0,
  "created_at": "2024-01-20T10:00:00Z"
}
```

### GET /api/v1/showcases?page=1&category_id=2&sort=newest

List showcases (public feed)

**Query Params:**

- `page=1`
- `category_id=2` (optional)
- `sort=newest|trending|most_liked` (default: newest)
- `search=keyword` (optional)

**Response (200):**

```json
{
  "data": [
    {
      "id": 42,
      "title": "My ML Model - 95% Accuracy",
      "description": "Trained CNN classifier on CIFAR-10",
      "media_urls": ["https://...", "https://..."],
      "user": {
        "id": 5,
        "name": "John Dev",
        "level": 3,
        "points": 850
      },
      "category": { "id": 2, "name": "Machine Learning" },
      "likes_count": 24,
      "views_count": 156,
      "created_at": "2024-01-20T10:00:00Z",
      "is_liked_by_me": false
    }
  ],
  "pagination": { "page": 1, "limit": 20, "total": 543 }
}
```

### GET /api/v1/showcases/:id

Get showcase detail with comments

**Response (200):**

```json
{
  "id": 42,
  "title": "My ML Model - 95% Accuracy",
  "description": "Trained CNN classifier on CIFAR-10",
  "media_urls": ["https://...", "https://..."],
  "user": { "id": 5, "name": "John Dev", "level": 3 },
  "category": { "id": 2, "name": "Machine Learning", "color": "#FF6B6B" },
  "likes_count": 24,
  "views_count": 157,
  "is_liked_by_me": true,
  "comments": [
    {
      "id": 1001,
      "content": "Great work!",
      "user": { "id": 10, "name": "Sarah" },
      "parent_id": null,
      "child_comments": []
    }
  ],
  "created_at": "2024-01-20T10:00:00Z"
}
```

### POST /api/v1/showcases/:id/like

Like showcase

**Headers:** `Authorization: Bearer {user_token}`

**Response (200):**

```json
{
  "showcase_id": 42,
  "liked": true,
  "likes_count": 25
}
```

### DELETE /api/v1/showcases/:id/like

Unlike showcase

**Response (200):**

```json
{
  "showcase_id": 42,
  "liked": false,
  "likes_count": 24
}
```

### PUT /api/v1/showcases/:id

Edit showcase (owner only)

**Request:**

```json
{
  "title": "Updated title",
  "description": "Updated description",
  "visibility": "private"
}
```

### DELETE /api/v1/showcases/:id

Delete showcase (owner or admin)

**Response (204):** No content

---

## 🎮 Gamification Endpoints

### GET /api/v1/leaderboard?page=1&timeframe=all

Get leaderboard

**Query Params:**

- `timeframe=all|month|week` (default: all)
- `page=1`

**Response (200):** [See leaderboard response above]

### GET /api/v1/users/:id/points

Get user's points and level (public)

**Response (200):**

```json
{
  "user_id": 5,
  "name": "John Dev",
  "points": 250,
  "level": 2,
  "level_info": {
    "level": 2,
    "badge_name": "Learner",
    "min_points": 200,
    "max_points": 499,
    "progress": 50
  },
  "rank": 42
}
```

### GET /api/v1/users/me/activity

Get current user's activity log

**Headers:** `Authorization: Bearer {user_token}`

**Response (200):**

```json
{
  "data": [
    {
      "id": 1,
      "activity_type": "showcase_liked",
      "points_earned": 5,
      "points_after": 255,
      "description": "Showcase 'My ML Model' was liked",
      "metadata": { "showcase_id": 42, "liker_name": "Sarah" },
      "created_at": "2024-01-20T14:30:00Z"
    },
    {
      "id": 2,
      "activity_type": "discussion_reply",
      "points_earned": 2,
      "points_after": 250,
      "description": "Replied to discussion 'Database Optimization'",
      "created_at": "2024-01-20T12:15:00Z"
    }
  ],
  "pagination": { "page": 1, "limit": 50, "total": 87 }
}
```

### GET /api/v1/levels

Get all level information

**Response (200):**

```json
{
  "data": [
    {
      "level": 1,
      "badge_name": "Beginner",
      "min_points": 0,
      "max_points": 199,
      "badge_icon": "https://...",
      "description": "Welcome to the community!"
    },
    {
      "level": 2,
      "badge_name": "Learner",
      "min_points": 200,
      "max_points": 499,
      "badge_icon": "https://...",
      "description": "You're making progress!"
    }
  ]
}
```

---

## 📌 Common Response Codes

| Code | Meaning                                 |
| ---- | --------------------------------------- |
| 200  | OK - Success                            |
| 201  | Created - Resource created              |
| 204  | No Content - Success, no response body  |
| 400  | Bad Request - Invalid input             |
| 401  | Unauthorized - No/invalid token         |
| 403  | Forbidden - No permission               |
| 404  | Not Found - Resource doesn't exist      |
| 409  | Conflict - Resource already exists      |
| 422  | Unprocessable Entity - Validation error |
| 429  | Too Many Requests - Rate limited        |
| 500  | Internal Server Error                   |

---

**All timestamps are in ISO 8601 format (UTC)**
