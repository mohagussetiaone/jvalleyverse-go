# 📚 JValleyVerse API Documentation

> **API Version**: 1.0.0  
> **Base URL**: `http://localhost:3000/api/v1` (development)  
> **Authorization**: Bearer JWT Token  
> **Content-Type**: `application/json`

---

## 📖 Table of Contents

1. [Authentication](#authentication)
2. [User Management](#user-management)
3. [Projects](#projects-admin)
4. [Classes](#classes)
5. [Discussions & Replies](#discussions--replies)
6. [Certificates](#certificates)
7. [Showcases](#showcases)
8. [Gamification](#gamification)
9. [Error Handling](#error-handling)
10. [Rate Limiting](#rate-limiting)

---

## 🔐 Authentication

### Register

**Endpoint**: `POST /auth/register`

Create a new user account.

**Request Body**:

```json
{
  "email": "user@example.com",
  "password": "secure_password",
  "name": "John Developer"
}
```

**Response** (201 Created):

```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Developer",
  "role": "user",
  "points": 0,
  "level": 1,
  "created_at": "2024-04-02T10:00:00Z"
}
```

**Error Responses**:

- `400 Bad Request` - Missing or invalid fields
- `409 Conflict` - Email already exists

---

### Login

**Endpoint**: `POST /auth/login`

Authenticate user and receive JWT token.

**Request Body**:

```json
{
  "email": "user@example.com",
  "password": "secure_password"
}
```

**Response** (200 OK):

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
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

**Error Responses**:

- `400 Bad Request` - Invalid input format
- `401 Unauthorized` - Wrong email or password

---

## 👤 User Management

### Get Current User Profile

**Endpoint**: `GET /users/me`

Get current authenticated user's profile.

**Headers**:

```
Authorization: Bearer {access_token}
```

**Response** (200 OK):

```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Developer",
  "avatar": "https://...",
  "bio": "Full-stack developer",
  "role": "user",
  "is_active": true,
  "points": 85,
  "total_points": 85,
  "level": 1,
  "created_at": "2024-04-02T10:00:00Z",
  "updated_at": "2024-04-02T12:30:00Z"
}
```

**Error Responses**:

- `401 Unauthorized` - Invalid or missing token

---

### Update User Profile

**Endpoint**: `PUT /users/me`

Update current user's profile information.

**Headers**:

```
Authorization: Bearer {access_token}
```

**Request Body** (all optional):

```json
{
  "name": "John Developer",
  "bio": "Senior Full-stack Developer",
  "avatar": "https://example.com/avatar.jpg"
}
```

**Response** (200 OK): Updated user object

---

### Get Public User Profile

**Endpoint**: `GET /users/{id}`

Get public profile of another user.

**Parameters**:

- `id` (path, required): User ID

**Response** (200 OK):

```json
{
  "id": 5,
  "name": "Sarah Expert",
  "avatar": "https://...",
  "bio": "Data Scientist",
  "level": 3,
  "points": 850,
  "showcase_count": 15,
  "certificate_count": 8
}
```

**Error Responses**:

- `404 Not Found` - User doesn't exist

---

### Get Leaderboard

**Endpoint**: `GET /leaderboard`

Get top users ranked by points.

**Query Parameters**:

- `page` (optional, default: 1): Page number
- `limit` (optional, default: 50, max: 100): Results per page
- `timeframe` (optional, default: all): `all`, `month`, `week`

**Response** (200 OK):

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

---

## 📌 Projects (Admin)

### Create Project

**Endpoint**: `POST /admin/projects`

Create new learning project (Admin only).

**Headers**:

```
Authorization: Bearer {admin_token}
Content-Type: application/json
```

**Request Body**:

```json
{
  "title": "Advanced Go Programming",
  "description": "Master backend development with Go",
  "thumbnail": "https://example.com/thumbnail.jpg",
  "category_id": 1,
  "visibility": "public"
}
```

**Response** (201 Created):

```json
{
  "id": 1,
  "title": "Advanced Go Programming",
  "description": "Master backend development with Go",
  "thumbnail": "https://...",
  "category_id": 1,
  "admin_id": 1,
  "visibility": "public",
  "created_at": "2024-04-02T10:00:00Z"
}
```

**Error Responses**:

- `403 Forbidden` - User is not admin
- `400 Bad Request` - Invalid input

---

### List Projects (Admin)

**Endpoint**: `GET /admin/projects`

List all projects with admin stats (Admin only).

**Headers**:

```
Authorization: Bearer {admin_token}
```

**Query Parameters**:

- `page` (optional, default: 1)
- `limit` (optional, default: 20)

**Response** (200 OK):

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
      "created_at": "2024-04-02T10:00:00Z"
    }
  ],
  "pagination": { "page": 1, "limit": 20, "total": 12 }
}
```

---

### Update Project

**Endpoint**: `PUT /admin/projects/{id}`

Update project details (Admin only).

**Request Body**:

```json
{
  "title": "Advanced Go (Updated)",
  "description": "...",
  "visibility": "public"
}
```

**Response** (200 OK): Updated project object

---

### Delete Project

**Endpoint**: `DELETE /admin/projects/{id}`

Delete project and cascade delete classes (Admin only).

**Response** (204 No Content)

---

## 📚 Classes

### Create Class

**Endpoint**: `POST /admin/classes`

Create class under project (Admin only).

**Request Body**:

```json
{
  "title": "Module 1: Goroutines & Concurrency",
  "description": "Learn how to use goroutines effectively",
  "content": "# Goroutines\n## Introduction\n...",
  "thumbnail": "https://...",
  "project_id": 1,
  "category_id": 1,
  "difficulty": "intermediate",
  "duration": 120,
  "order": 1
}
```

**Response** (201 Created):

```json
{
  "id": 10,
  "title": "Module 1: Goroutines & Concurrency",
  "project_id": 1,
  "category_id": 1,
  "admin_id": 1,
  "difficulty": "intermediate",
  "duration": 120,
  "order": 1,
  "created_at": "2024-04-02T10:00:00Z"
}
```

---

### Get Class Details

**Endpoint**: `GET /classes/{id}`

Get full class content including markdown.

**Response** (200 OK):

```json
{
  "id": 10,
  "title": "Module 1: Goroutines & Concurrency",
  "description": "Learn how to use goroutines effectively",
  "content": "# Goroutines\n## Introduction\n...",
  "thumbnail": "https://...",
  "project": { "id": 1, "title": "Advanced Go Programming" },
  "category": { "id": 1, "name": "Backend" },
  "admin": { "id": 1, "name": "Admin User" },
  "difficulty": "intermediate",
  "duration": 120
}
```

---

### List Project Classes

**Endpoint**: `GET /projects/{id}/classes`

Get all classes in a project.

**Response** (200 OK):

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

---

## 🏆 Certificates

### Get User Certificates

**Endpoint**: `GET /certificates`

Get current user's certificates (private).

**Headers**:

```
Authorization: Bearer {access_token}
```

**Response** (200 OK):

```json
{
  "data": [
    {
      "id": 42,
      "unique_code": "CERT-2024-ABC123XYZ",
      "class": {
        "id": 10,
        "title": "Module 1: Goroutines & Concurrency"
      },
      "badge_url": "https://...",
      "issued_at": "2024-04-02T10:00:00Z"
    }
  ]
}
```

---

### View Certificate

**Endpoint**: `GET /certificates/{code}`

View specific certificate (Owner + Admin only).

**Headers**:

```
Authorization: Bearer {access_token}
```

**Response** (200 OK):

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
  "issued_at": "2024-04-02T10:00:00Z",
  "expires_at": null
}
```

**Error Responses**:

- `403 Forbidden` - Not certificate owner
- `404 Not Found` - Certificate doesn't exist

---

## 💬 Discussions & Replies

### Create Discussion

**Endpoint**: `POST /discussions`

Create new discussion thread.

**Headers**:

```
Authorization: Bearer {access_token}
```

**Request Body**:

```json
{
  "title": "How to optimize database queries?",
  "content": "I have slow queries in production...",
  "class_id": 10,
  "category_id": 1
}
```

**Response** (201 Created):

```json
{
  "id": 100,
  "title": "How to optimize database queries?",
  "content": "I have slow queries in production...",
  "user": { "id": 5, "name": "John Dev" },
  "class_id": 10,
  "category_id": 1,
  "status": "open",
  "views_count": 0,
  "created_at": "2024-04-02T10:00:00Z"
}
```

---

### List Discussions

**Endpoint**: `GET /discussions`

List discussions with filtering.

**Query Parameters**:

- `page` (optional, default: 1)
- `class_id` (optional)
- `category_id` (optional)
- `status` (optional): `open`, `closed`

**Response** (200 OK):

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
      "created_at": "2024-04-02T10:00:00Z"
    }
  ],
  "pagination": { "page": 1, "limit": 20, "total": 156 }
}
```

---

### Get Discussion with Threaded Replies

**Endpoint**: `GET /discussions/{id}`

Get discussion with all replies organized in threads.

**Response** (200 OK):

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
      "created_at": "2024-04-02T10:30:00Z",
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

### Reply to Discussion

**Endpoint**: `POST /discussions/{id}/replies`

Post top-level reply to discussion.

**Request Body**:

```json
{
  "content": "Try adding indexes to foreign keys!"
}
```

**Response** (201 Created):

```json
{
  "id": 200,
  "content": "Try adding indexes to foreign keys!",
  "user": { "id": 10, "name": "Sarah" },
  "discussion_id": 100,
  "parent_id": null,
  "likes_count": 0,
  "created_at": "2024-04-02T10:30:00Z"
}
```

---

### Reply to a Reply (Nested)

**Endpoint**: `POST /replies/{id}/replies`

Post reply to another reply (threading).

**Request Body**:

```json
{
  "content": "Especially on foreign keys, yes!"
}
```

**Response** (201 Created):

```json
{
  "id": 201,
  "content": "Especially on foreign keys, yes!",
  "user": { "id": 11, "name": "Mike" },
  "discussion_id": 100,
  "parent_id": 200,
  "likes_count": 0,
  "created_at": "2024-04-02T10:35:00Z"
}
```

---

### Update Reply

**Endpoint**: `PUT /replies/{id}`

Edit reply (Owner only).

**Request Body**:

```json
{
  "content": "Updated content..."
}
```

**Response** (200 OK): Updated reply

---

### Delete Reply

**Endpoint**: `DELETE /replies/{id}`

Delete reply (Owner or Admin).

**Response** (204 No Content)

---

## 🎨 Showcases

### Create Showcase

**Endpoint**: `POST /showcases`

Create portfolio item.

**Headers**:

```
Authorization: Bearer {access_token}
```

**Request Body**:

```json
{
  "title": "My ML Model - 95% Accuracy",
  "description": "Trained CNN classifier on CIFAR-10",
  "media_urls": ["https://cdn.example.com/project1.jpg", "https://cdn.example.com/project2.jpg"],
  "category_id": 2,
  "visibility": "public"
}
```

**Response** (201 Created):

```json
{
  "id": 42,
  "title": "My ML Model - 95% Accuracy",
  "description": "Trained CNN classifier on CIFAR-10",
  "user": { "id": 5, "name": "John Dev" },
  "category_id": 2,
  "status": "published",
  "visibility": "public",
  "likes_count": 0,
  "views_count": 0,
  "created_at": "2024-04-02T10:00:00Z"
}
```

---

### List Showcases

**Endpoint**: `GET /showcases`

Browse public showcases (paginated feed).

**Query Parameters**:

- `page` (optional, default: 1)
- `category_id` (optional)
- `sort` (optional): `newest` (default), `trending`, `most_liked`
- `search` (optional): Search in title/description

**Response** (200 OK):

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
      "category": {
        "id": 2,
        "name": "Machine Learning",
        "color": "#FF6B6B"
      },
      "likes_count": 24,
      "views_count": 156,
      "is_liked_by_me": false,
      "created_at": "2024-04-02T10:00:00Z"
    }
  ],
  "pagination": { "page": 1, "limit": 20, "total": 543 }
}
```

---

### Get Showcase Detail

**Endpoint**: `GET /showcases/{id}`

Get showcase with comments and full details.

**Response** (200 OK):

```json
{
  "id": 42,
  "title": "My ML Model - 95% Accuracy",
  "description": "Trained CNN classifier on CIFAR-10",
  "media_urls": ["https://...", "https://..."],
  "user": {
    "id": 5,
    "name": "John Dev",
    "avatar": "https://...",
    "level": 3
  },
  "category": {
    "id": 2,
    "name": "Machine Learning",
    "color": "#FF6B6B"
  },
  "likes_count": 24,
  "views_count": 157,
  "is_liked_by_me": true,
  "comments": [
    {
      "id": 1001,
      "content": "Great work!",
      "user": { "id": 10, "name": "Sarah" },
      "parent_id": null,
      "created_at": "2024-04-02T11:00:00Z",
      "child_comments": []
    }
  ],
  "created_at": "2024-04-02T10:00:00Z"
}
```

---

### Like Showcase

**Endpoint**: `POST /showcases/{id}/like`

Like a showcase.

**Headers**:

```
Authorization: Bearer {access_token}
```

**Response** (200 OK):

```json
{
  "showcase_id": 42,
  "liked": true,
  "likes_count": 25
}
```

---

### Unlike Showcase

**Endpoint**: `DELETE /showcases/{id}/like`

Remove like from showcase.

**Headers**:

```
Authorization: Bearer {access_token}
```

**Response** (204 No Content)

---

### Update Showcase

**Endpoint**: `PUT /showcases/{id}`

Update showcase (Owner only).

**Request Body**:

```json
{
  "title": "Updated title",
  "description": "Updated description",
  "visibility": "private"
}
```

**Response** (200 OK): Updated showcase

---

### Delete Showcase

**Endpoint**: `DELETE /showcases/{id}`

Delete showcase (Owner or Admin).

**Response** (204 No Content)

---

## 🎮 Gamification

### Get User Points & Level

**Endpoint**: `GET /users/{id}/points`

Get user's points and level information (public).

**Response** (200 OK):

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

---

### Get User Activity Log

**Endpoint**: `GET /users/me/activity`

Get current user's activity history (points transactions).

**Headers**:

```
Authorization: Bearer {access_token}
```

**Query Parameters**:

- `page` (optional, default: 1)
- `limit` (optional, default: 50)

**Response** (200 OK):

```json
{
  "data": [
    {
      "id": 1,
      "activity_type": "showcase_liked",
      "points_earned": 5,
      "points_after": 255,
      "description": "Showcase 'My ML Model' was liked",
      "metadata": {
        "showcase_id": 42,
        "liker_name": "Sarah"
      },
      "created_at": "2024-04-02T14:30:00Z"
    },
    {
      "id": 2,
      "activity_type": "discussion_reply",
      "points_earned": 2,
      "points_after": 250,
      "description": "Replied to 'Database Optimization' discussion",
      "created_at": "2024-04-02T12:15:00Z"
    }
  ],
  "pagination": { "page": 1, "limit": 50, "total": 87 }
}
```

---

### Get Level Information

**Endpoint**: `GET /levels`

Get all level configuration and requirements.

**Response** (200 OK):

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
    },
    {
      "level": 3,
      "badge_name": "Contributor",
      "min_points": 500,
      "max_points": 999,
      "badge_icon": "https://...",
      "description": "Great contributions!"
    },
    {
      "level": 4,
      "badge_name": "Expert",
      "min_points": 1000,
      "max_points": 1999,
      "badge_icon": "https://...",
      "description": "You're an expert!"
    },
    {
      "level": 5,
      "badge_name": "Master",
      "min_points": 2000,
      "max_points": null,
      "badge_icon": "https://...",
      "description": "Master of the platform!"
    }
  ]
}
```

---

## ⚠️ Error Handling

All errors follow standard HTTP status codes and include error details.

### Error Response Format

```json
{
  "error": "Error message",
  "status": 400,
  "timestamp": "2024-04-02T10:00:00Z"
}
```

### Common HTTP Status Codes

| Code | Meaning              | Typical Cause                |
| ---- | -------------------- | ---------------------------- |
| 200  | OK                   | Successful GET or PUT        |
| 201  | Created              | Successful POST              |
| 204  | No Content           | Successful DELETE            |
| 400  | Bad Request          | Invalid input format         |
| 401  | Unauthorized         | Missing or invalid JWT token |
| 403  | Forbidden            | Insufficient permissions     |
| 404  | Not Found            | Resource doesn't exist       |
| 409  | Conflict             | Resource already exists      |
| 422  | Unprocessable Entity | Validation error             |
| 429  | Too Many Requests    | Rate limit exceeded          |
| 500  | Server Error         | Internal server error        |

---

## 🚦 Rate Limiting

API endpoints are rate limited to prevent abuse.

**Default Limits**:

- **Authenticated requests**: 100 requests per minute per user
- **Public requests**: 30 requests per minute per IP
- **Auth endpoints**: 5 requests per minute per IP (to prevent brute force)

**Rate Limit Headers**:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 99
X-RateLimit-Reset: 1617321600
```

When rate limit exceeded (429):

```json
{
  "error": "Too many requests",
  "retry_after": 60
}
```

---

## 🔒 Security Notes

### JWT Token

- Valid for **24 hours**
- Must be included in `Authorization: Bearer {token}` header
- Store securely on client (httpOnly cookie recommended)

### XSRF Protection

- XSRF token provided in login response
- Include in request headers for state-changing operations

### Certificate Privacy

- Certificates only accessible by owner + admin
- Verified by JWT user ID matching certificate user ID

### Admin Operations

- Project/Class creation restricted to admin role
- Verified at middleware level before handler

---

## 📝 Testing with cURL

### Register User

```bash
curl -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "John Dev"
  }'
```

### Login

```bash
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Get Profile (with token)

```bash
curl -X GET http://localhost:3000/api/v1/users/me \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Create Showcase

```bash
curl -X POST http://localhost:3000/api/v1/showcases \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My Project",
    "description": "Description here",
    "media_urls": ["https://example.com/img1.jpg"],
    "category_id": 1,
    "visibility": "public"
  }'
```

---

## 📚 Additional Resources

- **OpenAPI Spec**: See `openapi.json`
- **Postman Collection**: Import openapi.json into Postman
- **Implementation Guide**: See `README_IMPLEMENTATION.md`
- **Database Schema**: See `SCHEMA_DESIGN.md`

---

**Last Updated**: April 2, 2024  
**API Version**: 1.0.0
