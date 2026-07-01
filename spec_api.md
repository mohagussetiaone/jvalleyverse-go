# API Specification — JValleyverse

## General

| Attribute    | Value                         |
| ------------ | ----------------------------- |
| Base URL     | `http://localhost:3000`       |
| Auth         | JWT Bearer Token (24h expiry) |
| ID Format    | CUID (string)                 |
| Timestamps   | ISO 8601 (UTC)                |
| Rate Limit   | Global 200 req/min/IP         |
| Content-Type | `application/json`            |
| File Upload  | `multipart/form-data`         |

## Response Codes

| Code | Meaning                                             |
| ---- | --------------------------------------------------- |
| 200  | OK                                                  |
| 201  | Created                                             |
| 204  | No Content (DELETE success)                         |
| 400  | Bad Request — invalid input                         |
| 401  | Unauthorized — no/invalid token                     |
| 403  | Forbidden — not owner / not admin / scraper blocked |
| 404  | Not Found                                           |
| 409  | Conflict — already exists                           |
| 429  | Too Many Requests — rate limit exceeded             |
| 500  | Internal Server Error                               |

## Security Layers

| Layer                  | Protection             | Detail                                                                                               |
| ---------------------- | ---------------------- | ---------------------------------------------------------------------------------------------------- |
| **Rate Limit Global**  | 200 req/min/IP         | Semua route                                                                                          |
| **Rate Limit Content** | 60 req/min/IP          | Public GET: courses, lessons, showcases, categories, study-cases                                     |
| **Rate Limit Auth**    | 10 req/min/IP          | Login & register                                                                                     |
| **Anti-Scraping**      | User-Agent block       | ScraperGuard: blokir curl, python-requests, Postman, wget, dll                                       |
| **Security Headers**   | 6 headers              | X-Content-Type-Options, X-Frame-Options, X-XSS-Protection, Referrer-Policy, Permissions-Policy, HSTS |
| **XSRF Protection**    | Double submit cookie + Origin/Referer fallback | Cookie `XSRF-TOKEN` (SameSite=None, Secure=true, HttpOnly=false) + Header `X-XSRF-TOKEN`. Jika cookie check gagal, fallback ke Origin/Referer header (hanya untuk grup Dangerous & Admin)                     |
| **CORS**               | Origin whitelist       | Dikonfigurasi via `CORS_ORIGINS`                                                                     |
| **Idempotency**        | Idempotency-Key (UUID) | Safe retry untuk semua POST/PUT/DELETE                                                               |

---

## PUBLIC ROUTES (No Auth)

Public content GET endpoints dilindungi **ScraperGuard** + **ContentRateLimiter (60/min/IP)**.
Route bertanda `(Opt.JWT)` menyertakan `is_enrolled` jika user login.

### Auth

| Method | Path               | Middleware                   |
| ------ | ------------------ | ---------------------------- |
| POST   | /api/auth/register | AuthRateLimiter, Idempotency |
| POST   | /api/auth/login    | AuthRateLimiter, Idempotency |
| POST   | /api/auth/google   | AuthRateLimiter, Idempotency |
| POST   | /api/auth/refresh  | Idempotency                  |
| POST   | /api/auth/logout   | JWTAuth                      |

#### POST /api/auth/register

```json
// Request
{
  "name": "Test User",
  "email": "test@example.com",
  "password": "password123"
}
// Response 201
{
  "message": "User created"
}
// Response 400
{
  "errors": [
    {"field": "email", "message": "Email is not valid"},
    {"field": "password", "message": "Password must be at least 8 characters"}
  ]
}
// Response 409
{
  "error": "Email already registered"
}
```

#### POST /api/auth/login

```json
// Request
{
  "email": "admin@jvalleyverse.com",
  "password": "Admin@123"
}
// Response 200
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiY21xZXIxOGtxMDAwaGdjdWpnOGx5aGh1eCIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTc1MDI4MTYwMH0.abc123",
  "refresh_token": "a1b2c3d4e5f6...32bytehex",
  "expires_in": 86400,
  "xsrf_token": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "user": {
    "id": "cmqer18kq000hgcujg8lyhhux",
    "name": "JValley Admin",
    "email": "admin@jvalleyverse.com",
    "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/uuid.jpg",
    "role": "admin"
  }
}
// Response 401
{
  "error": "Invalid credentials"
}
```

#### POST /api/auth/google

```json
// Request
{
  "token": "<google_id_token_jwt>"
}
// Response 200
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiY21xZXIxOGtxMDAwaGdjdWpnOGx5aGh1eCIsInJvbGUiOiJ1c2VyIiwiZXhwIjoxNzUwMjgxNjAwfQ.abc123",
  "refresh_token": "a1b2c3d4e5f6...32bytehex",
  "expires_in": 86400,
  "xsrf_token": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "user": {
    "id": "cmqer18kq000hgcujg8lyhhux",
    "name": "Google User",
    "email": "user@gmail.com",
    "avatar": "https://lh3.googleusercontent.com/a/photo.jpg",
    "role": "user"
  }
}
// Response 401
{
  "error": "Invalid Google token"
}
```

#### POST /api/auth/refresh

```json
// Request
{
  "refresh_token": "a1b2c3d4e5f6...32bytehex"
}
// Response 200
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.new_token",
  "expires_in": 86400
}
// Response 401
{
  "error": "Invalid or expired refresh token"
}
```

#### POST /api/auth/logout

```json
// Headers
Authorization: Bearer <access_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Request (optional — revokes specific token if provided, otherwise all)
{
  "refresh_token": "a1b2c3d4e5f6...32bytehex"
}
// Response 200
{
  "message": "Logged out successfully"
}
```

### Categories

| Method | Path                                 | Middleware                         |
| ------ | ------------------------------------ | ---------------------------------- |
| GET    | /api/categories                      | ScraperGuard + ContentRL           |
| GET    | /api/categories/:slug                | ScraperGuard + ContentRL           |
| GET    | /api/categories/:category_id/courses | ScraperGuard + ContentRL + Opt.JWT |

#### GET /api/categories

```json
// Response 200
[
  {
    "id": "cmqer18kq000hgcujg8lyhhux",
    "name": "Backend Development",
    "slug": "backend-development"
  },
  {
    "id": "cmqer18kq001hgcujg8lyhhvy",
    "name": "Frontend Development",
    "slug": "frontend-development"
  }
]
```

#### GET /api/categories/:slug

```json
// Response 200
{
  "id": "cmqer18kq000hgcujg8lyhhux",
  "name": "Backend Development",
  "slug": "backend-development",
  "description": "Belajar backend development dengan berbagai teknologi",
  "courses": [
    {
      "id": "cmqer18kq002hgcujg8lyhhwz",
      "title": "Membangun REST API dengan Go & Fiber",
      "description": "Pelajari cara membuat RESTful API yang scalable menggunakan Go dan Fiber framework",
      "thumbnail": "https://cdn.mohagussetiaone.my.id/jvalleyverse/courses/thumb.jpg",
      "price": 0,
      "category": {
        "id": "cmqer18kq000hgcujg8lyhhux",
        "name": "Backend Development",
        "slug": "backend-development"
      },
      "admin_name": "JValley Admin",
      "mentor": {
        "id": "cmqer18kq003hgcujg8lyhhxa",
        "name": "Budi Mentor",
        "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/mentor.jpg",
        "role": "mentor"
      },
      "hours": 40,
      "visibility": "public",
      "section_count": 2,
      "is_enrolled": false,
      "created_at": "2026-06-15T12:04:44.714Z"
    }
  ]
}
// Response 404
{
  "error": "Category not found"
}
```

#### GET /api/categories/:category_id/courses

```json
// Response 200 — returns CourseListItem array directly (no wrapper)
[
  {
    "id": "cmqer18kq002hgcujg8lyhhwz",
    "title": "Membangun REST API dengan Go & Fiber",
    "description": "Pelajari cara membuat RESTful API yang scalable menggunakan Go dan Fiber framework",
    "thumbnail": "https://cdn.mohagussetiaone.my.id/jvalleyverse/courses/thumb.jpg",
    "price": 0,
    "category": {
      "id": "cmqer18kq000hgcujg8lyhhux",
      "name": "Backend Development",
      "slug": "backend-development"
    },
    "admin_name": "JValley Admin",
    "mentor": {
        "id": "cmqer18kq003hgcujg8lyhhxa",
        "name": "Budi Mentor",
        "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/mentor.jpg",
        "role": "mentor"
      },
    "hours": 40,
    "visibility": "public",
    "section_count": 2,
    "is_enrolled": true,
    "created_at": "2026-06-15T12:04:44.714Z"
  }
]
```

### Courses & Sections & Lessons

| Method | Path                                                 | Middleware                         |
| ------ | ---------------------------------------------------- | ---------------------------------- |
| GET    | /api/courses                                         | ScraperGuard + ContentRL + Opt.JWT |
| GET    | /api/courses/:course_id                              | ScraperGuard + ContentRL + Opt.JWT |
| GET    | /api/courses/:course_id/sections                     | ScraperGuard + ContentRL           |
| GET    | /api/courses/:course_id/sections/:section_id         | ScraperGuard + ContentRL           |
| GET    | /api/courses/:course_id/reviews                      | ScraperGuard + ContentRL           |
| GET    | /api/lessons/:id                                     | ScraperGuard + ContentRL           |
| GET    | /api/lessons/:id/reviews                             | ScraperGuard + ContentRL           |
| GET    | /api/courses/:course_id/lessons                      | ScraperGuard + ContentRL           |
| GET    | /api/courses/:course_id/sections/:section_id/lessons | ScraperGuard + ContentRL           |
| GET    | /api/courses/:course_id/lessons/:slug                | ScraperGuard + ContentRL           |

#### GET /api/courses

```
// Query: ?page=1&limit=10&category_id=cmq...&min_price=0&max_price=100000
//   category_id — Filter by category ID (optional)
//   min_price   — Filter courses with price >= this value (optional)
//   max_price   — Filter courses with price <= this value (optional)
```

```json
// Query: ?page=1&limit=10
// Response 200
{
  "data": [
    {
      "id": "cmqer18kq000hgcujg8lyhhux",
      "title": "Membangun REST API dengan Go & Fiber",
      "description": "Pelajari cara membuat RESTful API yang scalable menggunakan Go dan Fiber framework",
      "thumbnail": "https://cdn.mohagussetiaone.my.id/jvalleyverse/courses/thumb.jpg",
      "price": 0,
      "category": {
        "id": "cmqer18kq001hgcujg8lyhhvy",
        "name": "Backend Development",
        "slug": "backend-development"
      },
      "admin_name": "JValley Admin",
      "mentor": {
        "id": "cmqer18kq003hgcujg8lyhhxa",
        "name": "Budi Mentor",
        "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/mentor.jpg",
        "role": "mentor"
      },
      "hours": 40,
      "visibility": "public",
      "section_count": 2,
      "is_enrolled": true,
      "created_at": "2026-06-15T12:04:44.714Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 5
  }
}
```

#### GET /api/courses/:course_id

```json
// Response 200
{
  "id": "cmqer18kq000hgcujg8lyhhux",
  "title": "Membangun REST API dengan Go & Fiber",
  "description": "Pelajari cara membuat RESTful API yang scalable menggunakan Go dan Fiber framework",
  "thumbnail": "https://cdn.mohagussetiaone.my.id/jvalleyverse/courses/thumb.jpg",
  "price": 0,
  "category": {
    "id": "cmqer18kq001hgcujg8lyhhvy",
    "name": "Backend Development",
    "slug": "backend-development"
  },
  "admin_id": "cmqer18kq002hgcujg8lyhhwz",
  "admin_name": "JValley Admin",
  "mentor": {
        "id": "cmqer18kq003hgcujg8lyhhxa",
        "name": "Budi Mentor",
        "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/mentor.jpg",
        "role": "mentor"
      },
  "hours": 40,
  "total_duration_hours": 10,
  "visibility": "public",
  "is_enrolled": true,
  "created_at": "2026-06-15T12:04:44.714Z",
  "sections": [
    {
      "id": "cmqer18kq004hgcujg8lyhhyb",
      "title": "Introduction to Go",
      "description": "Dasar-dasar bahasa pemrograman Go",
      "order_index": 1,
      "lessons": [
        {
          "id": "cmqer18kq005hgcujg8lyhhyc",
          "title": "Go Basics",
          "slug": "go-basics",
          "difficulty": "beginner",
          "duration": 45,
          "order_index": 1,
          "video_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/lessons/video.mp4"
        }
      ]
    },
    {
      "id": "cmqer18kq006hgcujg8lyhhyd",
      "title": "Building REST API with Fiber",
      "description": "Implementasi REST API menggunakan Fiber",
      "order_index": 2,
      "lessons": [
        {
          "id": "cmqer18kq007hgcujg8lyhhyf",
          "title": "Fiber Framework",
          "slug": "fiber-framework",
          "difficulty": "intermediate",
          "duration": 60,
          "order_index": 1,
          "video_url": ""
        }
      ]
    }
  ]
}
// Response 404
{
  "error": "Course not found"
}
```

#### GET /api/courses/:course_id/sections

```json
// Response 200
{
  "data": [
    {
      "id": "cmqer18kq004hgcujg8lyhhyb",
      "title": "Introduction to Go",
      "description": "Dasar-dasar bahasa pemrograman Go",
      "course_id": "cmqer18kq000hgcujg8lyhhux",
      "order_index": 1,
      "lessons": [
        {
          "id": "cmqer18kq005hgcujg8lyhhyc",
          "title": "Go Basics",
          "slug": "go-basics",
          "difficulty": "beginner",
          "duration": 45,
          "order_index": 1,
          "video_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/lessons/video.mp4"
        }
      ]
    }
  ]
}
```

#### GET /api/courses/:course_id/sections/:section_id

```json
// Response 200
{
  "id": "cmqer18kq004hgcujg8lyhhyb",
  "title": "Introduction to Go",
  "description": "Dasar-dasar bahasa pemrograman Go",
  "course_id": "cmqer18kq000hgcujg8lyhhux",
  "order_index": 1,
  "lessons": [
    {
      "id": "cmqer18kq005hgcujg8lyhhyc",
      "title": "Go Basics",
      "slug": "go-basics",
      "difficulty": "beginner",
      "duration": 45,
      "order_index": 1,
      "video_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/lessons/video.mp4"
    }
  ]
}
// Response 404
{
  "error": "Section not found"
}
```

#### GET /api/courses/:course_id/reviews

```
// Query: ?page=1&limit=20
```

```json
// Response 200
{
  "data": [
    {
      "id": "cmqer18kq008hgcujg8lyhhyg",
      "user_id": "cmqer18kq009hgcujg8lyhhyh",
      "user_name": "John Doe",
      "user_avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/user.jpg",
      "course_id": "cmqer18kq000hgcujg8lyhhux",
      "lesson_id": "",
      "rating": 5,
      "message": "Kursusnya sangat bagus dan mudah dipahami",
      "created_at": "2026-06-15T12:04:44.714Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 3
  }
}
```

#### GET /api/lessons/:id

```json
// Response 200
{
  "lesson": {
    "id": "cmqer18kq005hgcujg8lyhhyc",
    "title": "Go Basics",
    "slug": "go-basics",
    "difficulty": "beginner",
    "duration": 45,
    "order_index": 1,
    "video_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/lessons/video.mp4"
  },
  "details": {
    "id": "cmqer18kq010hgcujg8lyhhyi",
    "lesson_id": "cmqer18kq005hgcujg8lyhhyc",
    "about": "Pelajari fundamentals Go: variables, types, control flow",
    "rules": "1. Selesaikan semua latihan\n2. Submit code kamu\n3. Diskusi jika ada kendala",
    "tools": [
      "Go 1.22+",
      "VS Code",
      "Terminal"
    ],
    "resource_media": {
      "videos": ["https://cdn.mohagussetiaone.my.id/jvalleyverse/lessons/intro.mp4"],
      "documents": ["https://cdn.mohagussetiaone.my.id/jvalleyverse/lessons/guide.pdf"],
      "images": ["https://cdn.mohagussetiaone.my.id/jvalleyverse/lessons/diagram.png"]
    },
    "resources": [
      {"type": "article", "title": "Go Tour", "url": "https://go.dev/tour"},
      {"type": "github", "title": "Source Code", "url": "https://github.com/example/go-basics"}
    ],
    "created_at": "2026-06-15T12:04:44.714Z",
    "updated_at": "2026-06-15T12:04:44.714Z"
  },
  "progress": null,
  "next_lesson": {
    "id": "cmqer18kq007hgcujg8lyhhyf",
    "title": "Fiber Framework",
    "slug": "fiber-framework",
    "difficulty": "intermediate",
    "duration": 60,
    "order_index": 2,
    "video_url": ""
  },
  "section": {
    "id": "cmqer18kq004hgcujg8lyhhyb",
    "title": "Introduction to Go",
    "description": "Dasar-dasar bahasa pemrograman Go",
    "order_index": 1,
    "lessons": null
  },
  "course": {
    "id": "cmqer18kq000hgcujg8lyhhux",
    "title": "Membangun REST API dengan Go & Fiber",
    "description": "Pelajari cara membuat RESTful API yang scalable menggunakan Go dan Fiber framework",
    "thumbnail": "https://cdn.mohagussetiaone.my.id/jvalleyverse/courses/thumb.jpg",
    "price": 0,
    "category": {
      "id": "cmqer18kq001hgcujg8lyhhvy",
      "name": "Backend Development",
      "slug": "backend-development"
    },
    "admin_name": "JValley Admin",
    "mentor": {
        "id": "cmqer18kq003hgcujg8lyhhxa",
        "name": "Budi Mentor",
        "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/mentor.jpg",
        "role": "mentor"
      },
    "hours": 40,
    "visibility": "public",
    "section_count": 2,
    "created_at": "2026-06-15T12:04:44.714Z"
  }
}
// Response 404
{
  "error": "Lesson not found"
}
```

#### GET /api/lessons/:id/reviews

```
// Query: ?page=1&limit=20
```

```json
// Response 200
{
  "data": [
    {
      "id": "cmqer18kq008hgcujg8lyhhyg",
      "user_id": "cmqer18kq009hgcujg8lyhhyh",
      "user_name": "John Doe",
      "user_avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/user.jpg",
      "course_id": "",
      "lesson_id": "cmqer18kq005hgcujg8lyhhyc",
      "rating": 4,
      "message": "Materinya bagus, tapi butuh lebih banyak contoh",
      "created_at": "2026-06-15T12:04:44.714Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 2
  }
}
```

#### GET /api/courses/:course_id/lessons

```json
// Query: ?page=1&limit=20
// Response 200
{
  "data": [
    {
      "id": "cmqer18kq005hgcujg8lyhhyc",
      "title": "Go Basics",
      "slug": "go-basics",
      "difficulty": "beginner",
      "duration": 45,
      "order_index": 1,
      "video_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/lessons/video.mp4"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 10
  }
}
```

#### GET /api/courses/:course_id/sections/:section_id/lessons

```json
// Response 200
{
  "data": [
    {
      "id": "cmqer18kq005hgcujg8lyhhyc",
      "title": "Go Basics",
      "slug": "go-basics",
      "difficulty": "beginner",
      "duration": 45,
      "order_index": 1,
      "video_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/lessons/video.mp4"
    }
  ],
  "total": 5
}
```

#### GET /api/courses/:course_id/lessons/:slug

```json
// Response 200 — same structure as GET /api/lessons/:id
{
  "lesson": {
    "id": "cmqer18kq005hgcujg8lyhhyc",
    "title": "Go Basics",
    "slug": "go-basics",
    "difficulty": "beginner",
    "duration": 45,
    "order_index": 1,
    "video_url": ""
  },
  "details": {
    "id": "cmqer18kq010hgcujg8lyhhyi",
    "lesson_id": "cmqer18kq005hgcujg8lyhhyc",
    "about": "Pelajari fundamentals Go: variables, types, control flow",
    "rules": "1. Selesaikan semua latihan\n2. Submit code kamu\n3. Diskusi jika ada kendala",
    "tools": ["Go 1.22+", "VS Code", "Terminal"],
    "resource_media": {
      "videos": ["https://cdn.mohagussetiaone.my.id/jvalleyverse/lessons/intro.mp4"],
      "documents": ["https://cdn.mohagussetiaone.my.id/jvalleyverse/lessons/guide.pdf"],
      "images": ["https://cdn.mohagussetiaone.my.id/jvalleyverse/lessons/diagram.png"]
    },
    "resources": [
      {"type": "article", "title": "Go Tour", "url": "https://go.dev/tour"},
      {"type": "github", "title": "Source Code", "url": "https://github.com/example/go-basics"}
    ],
    "created_at": "2026-06-15T12:04:44.714Z",
    "updated_at": "2026-06-15T12:04:44.714Z"
  },
  "progress": null,
  "next_lesson": null,
  "section": {
    "id": "cmqer18kq004hgcujg8lyhhyb",
    "title": "Introduction to Go",
    "description": "Dasar-dasar bahasa pemrograman Go",
    "order_index": 1,
    "lessons": null
  },
  "course": {
    "id": "cmqer18kq000hgcujg8lyhhux",
    "title": "Membangun REST API dengan Go & Fiber",
    "description": "Pelajari cara membuat RESTful API yang scalable menggunakan Go dan Fiber framework",
    "thumbnail": "https://cdn.mohagussetiaone.my.id/jvalleyverse/courses/thumb.jpg",
    "price": 0,
    "category": {
      "id": "cmqer18kq001hgcujg8lyhhvy",
      "name": "Backend Development",
      "slug": "backend-development"
    },
    "admin_name": "JValley Admin",      "mentor": null,
      "hours": 40,
    "visibility": "public",
    "section_count": 2,
    "created_at": "2026-06-15T12:04:44.714Z"
  }
}
// Response 404
{
  "error": "Lesson not found"
}
```

### Blogs

| Method | Path           | Middleware               |
| ------ | -------------- | ------------------------ |
| GET    | /api/blogs     | ScraperGuard + ContentRL |
| GET    | /api/blogs/:id | ScraperGuard + ContentRL |

#### GET /api/blogs

```json
// Query: ?page=1&limit=10&search=golang&category_id=cmq...&tag=golang
// Response 200
{
  "data": [
    {
      "id": "cmqer18kq011hgcujg8lyhhyj",
      "title": "Belajar Go untuk Pemula",
      "slug": "belajar-go-untuk-pemula",
      "description": "Panduan lengkap belajar Go dari dasar hingga mahir untuk pemula",
      "cover_img_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/blogs/cover.jpg",
      "tags": ["golang", "backend"],
      "status": "published",
      "author": {
        "id": "cmqer18kq012hgcujg8lyhhyk",
        "name": "JValley Admin",
        "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/admin.jpg"
      },
      "category": {
        "id": "cmqer18kq001hgcujg8lyhhvy",
        "name": "Backend Development",
        "slug": "backend-development"
      },
      "created_at": "2026-06-15T12:00:00.000Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 3
  }
}
```

#### GET /api/blogs/:id

```json
// Response 200
{
  "id": "cmqer18kq011hgcujg8lyhhyj",
  "title": "Belajar Go untuk Pemula",
  "slug": "belajar-go-untuk-pemula",
  "description": "Panduan lengkap belajar Go dari dasar hingga mahir untuk pemula",
  "content": "# Belajar Go untuk Pemula\n\n## Introduction\nGo adalah bahasa pemrograman yang dikembangkan oleh Google...\n\n## Setup\n1. Download Go dari golang.org\n2. Install dan setup GOPATH\n3. Tulis program pertama kamu...",
  "cover_img_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/blogs/cover.jpg",
  "tags": ["golang", "backend"],
  "status": "published",
  "author": {
    "id": "cmqer18kq012hgcujg8lyhhyk",
    "name": "JValley Admin",
    "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/admin.jpg"
  },
  "category": {
    "id": "cmqer18kq001hgcujg8lyhhvy",
    "name": "Backend Development",
    "slug": "backend-development"
  },
  "created_at": "2026-06-15T12:00:00.000Z",
  "updated_at": "2026-06-15T14:00:00.000Z"
}
// Response 404
{
  "error": "Not found"
}
```

### Discussions (Read-only Public)

| Method | Path                 | Middleware                         |
| ------ | -------------------- | ---------------------------------- |
| GET    | /api/discussions     | ScraperGuard + ContentRL + Opt.JWT |
| GET    | /api/discussions/:id | ScraperGuard + ContentRL + Opt.JWT |

#### GET /api/discussions

```json
// Query: ?page=1&limit=20&lesson_id=cmq...&study_case_id=cmq...
// Response 200
{
  "data": [
    {
      "id": "cmqer18kq013hgcujg8lyhhyl",
      "title": "Perbedaan goroutine dan thread biasa?",
      "content": "Saya baru belajar Go, apakah ada yang bisa jelaskan perbedaan goroutine dengan thread OS?",
      "user": {
        "id": "cmqer18kq009hgcujg8lyhhyh",
        "name": "John Doe",
        "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/user.jpg",
        "role": "user"
      },
      "lesson_id": "cmqer18kq005hgcujg8lyhhyc",
      "study_case_id": null,
      "status": "open",
      "view_count": 15,
      "created_at": "2026-06-15T12:04:44.714Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 50
  }
}
```

#### GET /api/discussions/:id

```json
// Response 200
{
  "id": "cmqer18kq013hgcujg8lyhhyl",
  "title": "Perbedaan goroutine dan thread biasa?",
  "content": "Saya baru belajar Go, apakah ada yang bisa jelaskan perbedaan goroutine dengan thread OS?",
  "user": {
    "id": "cmqer18kq009hgcujg8lyhhyh",
    "name": "John Doe",
    "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/user.jpg",
    "role": "user"
  },
  "status": "open",
  "view_count": 15,
  "created_at": "2026-06-15T12:04:44.714Z",
  "replies": [
    {
      "id": "cmqer18kq014hgcujg8lyhhym",
      "content": "Goroutine lebih ringan dari thread OS karena menggunakan scheduler sendiri yang berjalan di userspace",
      "user": {
        "id": "cmqer18kq012hgcujg8lyhhyk",
        "name": "JValley Admin",
        "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/admin.jpg",
        "role": "admin"
      },
      "likes": 5,
      "is_best": true,
      "created_at": "2026-06-15T13:00:00.000Z"
    },
    {
      "id": "cmqer18kq015hgcujg8lyhhyn",
      "content": "Selain itu, goroutine punya stack yang dinamis, awalnya hanya 2KB",
      "user": {
        "id": "cmqer18kq009hgcujg8lyhhyh",
        "name": "John Doe",
        "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/user.jpg",
        "role": "user"
      },
      "likes": 2,
      "is_best": false,
      "created_at": "2026-06-15T13:30:00.000Z"
    }
  ]
}
// Response 404
{
  "error": "Discussion not found"
}
```

### Study Cases

| Method | Path                 | Middleware               |
| ------ | -------------------- | ------------------------ |
| GET    | /api/study-cases     | ScraperGuard + ContentRL |
| GET    | /api/study-cases/:id | ScraperGuard + ContentRL |

#### GET /api/study-cases

```
// Query: ?page=1&limit=20&category_id=cmq...
//   category_id — Filter by category ID (optional)
```

```json
// Query: ?page=1&limit=20
// Response 200
{
  "data": [
    {
      "id": "cmqer18kq016hgcujg8lyhhyo",
      "name": "Belajar Go Basic",
      "description": "Studi kasus fundamental Go untuk pemula",
      "img_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/study-cases/img.jpg",
      "youtube_url": "https://youtube.com/watch?v=abc123",
      "tags": ["golang", "beginner"],
      "category": {
        "id": "cmqer18kq001hgcujg8lyhhvy",
        "name": "Backend Development",
        "slug": "backend-development"
      },
      "user": {
        "id": "cmqer18kq012hgcujg8lyhhyk",
        "name": "JValley Admin",
        "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/admin.jpg"
      },
      "created_at": "2026-06-15T12:04:44.714Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 10
  }
}
```

#### GET /api/study-cases/:id

```json
// Response 200
{
  "id": "cmqer18kq016hgcujg8lyhhyo",
  "name": "Belajar Go Basic",
  "description": "Studi kasus fundamental Go untuk pemula",
  "img_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/study-cases/img.jpg",
  "youtube_url": "https://youtube.com/watch?v=abc123",
  "tags": ["golang", "beginner"],
  "user": {
    "id": "cmqer18kq012hgcujg8lyhhyk",
    "name": "JValley Admin",
    "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/admin.jpg",
    "role": "admin"
  },
  "discussions": [
    {
      "id": "cmqer18kq017hgcujg8lyhhyp",
      "title": "Diskusi tentang Go",
      "user_id": "cmqer18kq009hgcujg8lyhhyh",
      "user_name": "Budi",
      "status": "open",
      "views_count": 10,
      "created_at": "2026-06-15T12:04:44.714Z"
    }
  ],
  "created_at": "2026-06-15T12:04:44.714Z"
}
// Response 404
{
  "error": "Study case not found"
}
```

### Company Profile (Public)

| Method | Path          | Middleware |
| ------ | ------------- | ---------- |
| GET    | /api/company  | —          |

#### GET /api/company

```json
// Response 200
{
  "id": "cmqer18kq027hgcujg8lyhhzc",
  "created_at": "2026-06-18T10:00:00.000Z",
  "updated_at": "2026-06-18T10:00:00.000Z",
  "brand_name": "JValleyVerse",
  "tagline": "Learn, Build, Grow Together",
  "vision": "Menjadi platform edukasi teknologi terdepan di Indonesia yang mencetak talenta digital berkualitas dan siap bersaing di era global.",
  "mission": "Menyediakan materi pembelajaran berkualitas tinggi yang mudah diakses\nMembangun komunitas belajar yang kolaboratif dan suportif\nMenjembatani kesenjangan antara pendidikan formal dan kebutuhan industri\nMemberikan pengalaman belajar interaktif dengan gamifikasi dan sertifikasi",
  "logo_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/logo.png",
  "domain": "https://jvalleyverse.com",
  "email": "hello@jvalleyverse.com",
  "facebook": "https://facebook.com/jvalleyverse",
  "instagram": "https://instagram.com/jvalleyverse",
  "twitter": "https://x.com/jvalleyverse",
  "tiktok": "https://tiktok.com/@jvalleyverse",
  "youtube": "https://youtube.com/@jvalleyverse",
  "linkedin": "https://linkedin.com/company/jvalleyverse",
  "whatsapp": "https://wa.me/6281234567890",
  "address": "Jakarta, Indonesia",
  "phone": "+62 812-3456-7890"
}
```

### FAQs (Public)

| Method | Path       | Middleware |
| ------ | ---------- | ---------- |
| GET    | /api/faqs  | —          |

#### GET /api/faqs

```
// Query: ?page=1&limit=20
```

```json
// Response 200
{
  "data": [
    {
      "id": "cmqer18kq026hgcujg8lyhhzb",
      "question": "Apa itu JValleyverse?",
      "answer": "JValleyverse adalah platform belajar coding online dengan sistem gamifikasi, sertifikat, dan showcase.",
      "category": "general",
      "order_index": 1,
      "is_active": true,
      "created_at": "2026-06-18T10:00:00.000Z",
      "updated_at": "2026-06-18T10:00:00.000Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 5
  }
}
```

### Other Public

| Method | Path                 | Middleware               |
| ------ | -------------------- | ------------------------ |
| GET    | /api/leaderboard     | —                        |
| GET    | /api/mentors         | —                        |
| GET    | /api/showcases       | ScraperGuard + ContentRL |
| GET    | /api/showcases/:id   | ScraperGuard + ContentRL |
| GET    | /api/health          | — (no guard)             |
| GET    | /api/health/detailed | — (no guard)             |
| GET    | /api/system/status   | — (no guard)             |
| GET    | /api/users/:id       | —                        |
| GET    | /api/faqs            | — (no guard)             |
| GET    | /api/company         | — (no guard)             |

#### GET /api/leaderboard

```json
// Query: ?limit=10
// Response 200
{
  "data": [
    {
      "rank": 1,
      "user_id": "cmqer18kq012hgcujg8lyhhyk",
      "name": "JValley Admin",
      "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/admin.jpg",
      "total_points": 5000,
      "level": 5
    },
    {
      "rank": 2,
      "user_id": "cmqer18kq009hgcujg8lyhhyh",
      "name": "John Doe",
      "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/user.jpg",
      "total_points": 2500,
      "level": 5
    }
  ]
}
```

#### GET /api/mentors

```json
// Query: ?page=1&limit=20
// Response 200
{
  "data": [
    {
      "id": "cmqer18kq003hgcujg8lyhhxa",
      "name": "Budi Mentor",
      "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/mentor.jpg",
      "bio": "Full-stack developer dengan 10 tahun pengalaman",
      "level": 5,
      "total_points": 3500
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 3
  }
}
```

#### GET /api/showcases

```json
// Query: ?page=1&limit=20&category_id=cmq...&sort=newest
// Response 200
{
  "data": [
    {
      "id": "cmqer18kq018hgcujg8lyhhyq",
      "title": "My Portfolio Website",
      "description": "Built with React & Next.js",
      "media_urls": ["https://cdn.mohagussetiaone.my.id/jvalleyverse/showcases/img1.jpg"],
      "likes_count": 24,
      "views_count": 156,
      "visibility": "public",
      "user": {
        "id": "cmqer18kq009hgcujg8lyhhyh",
        "name": "John Doe",
        "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/user.jpg"
      },
      "category": {
        "id": "cmqer18kq019hgcujg8lyhhyv",
        "name": "Frontend Development",
        "slug": "frontend-development"
      },
      "created_at": "2026-06-15T12:04:44.714Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 50
  }
}
```

#### GET /api/showcases/:id

```json
// Response 200
{
  "id": "cmqer18kq018hgcujg8lyhhyq",
  "title": "My Portfolio Website",
  "description": "Built with React & Next.js",
  "media_urls": ["https://cdn.mohagussetiaone.my.id/jvalleyverse/showcases/img1.jpg"],
  "likes_count": 24,
  "views_count": 157,
  "user": {
    "id": "cmqer18kq009hgcujg8lyhhyh",
    "name": "John Doe",
    "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/user.jpg",
    "role": "user"
  },
  "category": {
    "id": "cmqer18kq019hgcujg8lyhhyv",
    "name": "Frontend Development",
    "slug": "frontend-development"
  },
  "is_liked_by_me": false,
  "created_at": "2026-06-15T12:04:44.714Z"
}
// Response 404
{
  "error": "Showcase not found"
}
```

#### GET /api/health

```json
// Response 200
{
  "status": "ok",
  "timestamp": "2026-06-18T10:00:00.000Z",
  "version": "1.0.0"
}
```

#### GET /api/health/detailed

```json
// Response 200 (local)
{
  "status": "ok",
  "environment": "local",
  "timestamp": "2026-06-18T10:00:00.000Z",
  "version": "1.0.0"
}
// Response 200 (production)
{
  "status": "ok",
  "environment": "production",
  "hostname": "vps-12345",
  "timestamp": "2026-06-18T10:00:00.000Z",
  "version": "1.0.0",
  "system": {
    "os": "linux",
    "arch": "amd64"
  },
  "memory": {
    "allocated_mb": 45.23,
    "total_allocated_mb": 120.45,
    "system_mb": 67.89,
    "num_gc": 152
  },
  "goroutines": 12
}
```

#### GET /api/users/:id

```json
// Response 200
{
  "id": "cmqer18kq009hgcujg8lyhhyh",
  "name": "John Doe",
  "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/user.jpg",
  "level": 3,
  "points": 850
}
// Response 404
{
  "error": "User not found"
}
```

#### GET /api/users/:id/portfolio

```
// No auth required (public)
// Aggregates certificates, showcases, and study cases for a shareable portfolio
```

```json
// Response 200
{
  "user": {
    "id": "cmqer18kq009hgcujg8lyhhyh",
    "name": "John Doe",
    "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/user.jpg"
  },
  "total_points": 850,
  "level": 3,
  "items": [
    {
      "id": "cmqer18kq023hgcujg8lyhhyz",
      "title": "Go Basics",
      "description": "Certificate of completion",
      "type": "certificate",
      "url": "https://jvalleyverse.com/certificates/verify/CERT-abc12345",
      "created_at": "2026-06-15T12:00:00.000Z",
      "tags": ["certificate"]
    },
    {
      "id": "cmqer18kq018hgcujg8lyhhyq",
      "title": "My Portfolio Website",
      "description": "Built with React & Next.js",
      "type": "showcase",
      "image_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/showcases/img1.jpg",
      "created_at": "2026-06-15T12:04:44.714Z",
      "tags": ["showcase"]
    },
    {
      "id": "cmqer18kq016hgcujg8lyhhyo",
      "title": "Belajar Go Basic",
      "description": "Studi kasus fundamental Go",
      "type": "study_case",
      "image_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/study-cases/img.jpg",
      "created_at": "2026-06-15T12:04:44.714Z",
      "tags": ["golang", "beginner"]
    }
  ],
  "cert_count": 5,
  "showcase_count": 3,
  "study_case_count": 2
}
// Response 404
{
  "error": "User not found"
}
```

#### GET /api/certificates/:code/verify

```
// No auth required (public)
// Verify a certificate by its unique code
```

```json
// Response 200
{
  "id": "cmqer18kq023hgcujg8lyhhyz",
  "unique_code": "CERT-abc12345",
  "issued_at": "2026-06-15T12:00:00.000Z",
  "user_id": "cmqer18kq009hgcujg8lyhhyh",
  "lesson_id": "cmqer18kq005hgcujg8lyhhyc",
  "lesson_name": "Go Basics",
  "user_name": "John Doe",
  "verification_url": "https://jvalleyverse.com/certificates/verify/CERT-abc12345",
  "qr_code_url": "https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=https%3A%2F%2Fjvalleyverse.com%2Fcertificates%2Fverify%2FCERT-abc12345",
  "achievement": {
    "type": "certificate",
    "title": "Go Basics",
    "unique_code": "CERT-abc12345"
  }
}
// Response 404
{
  "error": "Certificate not found"
}
```

#### GET /api/system/status

```
// No auth required (public)
// Returns real-time operational status of all system dependencies
```

```json
// Response 200
{
  "status": "all_operational",
  "uptime": "72h15m30s",
  "version": "1.0.0",
  "environment": "production",
  "timestamp": "2026-07-01T10:00:00.000Z",
  "services": [
    {
      "name": "database",
      "status": "operational",
      "message": "PostgreSQL connected",
      "latency": "2.5ms"
    },
    {
      "name": "redis",
      "status": "operational",
      "message": "Redis connected",
      "latency": "1.2ms"
    },
    {
      "name": "minio",
      "status": "operational",
      "message": "MinIO connected",
      "latency": "0.5ms"
    }
  ],
  "summary": {
    "total": 3,
    "operational": 3,
    "degraded": 0,
    "down": 0
  }
}
```

---

## PROTECTED ROUTES (JWT)

### User Profile & Dashboard (JWT only — no XSRF)

| Method | Path                          | Middleware |
| ------ | ----------------------------- | ---------- |
| GET    | /api/users/me                 | JWTAuth    |
| PUT    | /api/users/me                 | JWTAuth    |
| POST   | /api/users/me/change-password | JWTAuth    |
| POST   | /api/users/me/avatar          | JWTAuth    |
| GET    | /api/users/me/activity        | JWTAuth    |
| GET    | /api/users/me/dashboard       | JWTAuth    |

#### GET /api/users/me

```json
// Headers
Authorization: Bearer <access_token>

// Response 200
{
  "id": "cmqer18kq012hgcujg8lyhhyk",
  "created_at": "2026-06-10T08:00:00.000Z",
  "updated_at": "2026-06-18T10:00:00.000Z",
  "email": "admin@jvalleyverse.com",
  "name": "JValley Admin",
  "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/admin.jpg",
  "bio": "Full-stack developer dan content creator",
  "role": "admin",
  "is_active": true,
  "points": 85,
  "total_points": 5000,
  "level": 5
}
// Response 401
{
  "error": "Missing or invalid JWT token"
}
```

#### PUT /api/users/me

```json
// Headers
Authorization: Bearer <access_token>

// Request (all fields optional)
{
  "name": "JValley Admin Updated",
  "bio": "Senior Full-stack Developer & Mentor",
  "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/new.jpg"
}
// Response 200
{
  "message": "Profile updated"
}
// Response 400
{
  "error": "Invalid input"
}
```

#### GET /api/users/me/activity

```json
// Headers
Authorization: Bearer <access_token>

// Query: ?page=1&limit=20
// Response 200
{
  "data": [
    {
      "id": "cmqer18kq020hgcujg8lyhhyw",
      "activity": "create_showcase",
      "points": 10,
      "timestamp": "2026-06-15T12:04:44.714Z"
    },
    {
      "id": "cmqer18kq021hgcujg8lyhhyx",
      "activity": "lesson_completed",
      "points": 50,
      "timestamp": "2026-06-15T13:00:00.000Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 87
  }
}
```

#### GET /api/users/me/dashboard

```json
// Headers
Authorization: Bearer <access_token>

// Response 200
{
  "courses_in_progress": 2,
  "courses_completed": 3,
  "courses_dropped": 1,
  "unread_notifications": 5,
  "streak_count": 5
}
```

#### GET /api/users/me/streak

| Method | Path                | Middleware |
| ------ | ------------------- | ---------- |
| GET    | /api/users/me/streak | JWTAuth    |

```json
// Headers
Authorization: Bearer <access_token>

// Response 200
{
  "streak_count": 5,
  "longest_streak": 12,
  "last_activity_date": "2026-06-20T10:00:00.000Z"
}
// Response 401
{
  "error": "Missing or invalid JWT token"
}
```

### Safe Group (JWT + Idempotency — no XSRF)

Group `/api` — Endpoint **aman** yang tidak perlu XSRF. Cukup JWT + Idempotency-Key.

#### Change Password

| Method | Path                          | Middleware      |
| ------ | ----------------------------- | --------------- |
| POST   | /api/users/me/change-password | JWTAuth         |

```json
// Request
{
  "current_password": "OldPass123",
  "new_password": "NewPass456!"
}
// Response 200
{
  "message": "Password changed successfully"
}
// Response 401
{
  "error": "unauthorized"
}
```

#### Update Avatar (multipart upload to MinIO)

| Method | Path                   | Middleware |
| ------ | ---------------------- | ---------- |
| POST   | /api/users/me/avatar   | JWTAuth    |

**Request:** `multipart/form-data` dengan field `avatar` (file)

```json
// Response 200
{
  "message": "Avatar updated",
  "url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/uuid.jpg",
  "object_name": "avatars/uuid.jpg"
}
// Response 503
{
  "error": "Avatar upload is not available (MinIO not configured)"
}
```

#### Learning & Progress

| Method | Path                      | Middleware        |
| ------ | ------------------------- | ----------------- |
| POST   | /api/lessons/:id/start    | JWT + Idempotency |
| PUT    | /api/lessons/:id/progress | JWT + Idempotency |
| POST   | /api/lessons/:id/complete | JWT + Idempotency |

#### POST /api/lessons/:id/start

```json
// Headers
Authorization: Bearer <access_token>
Idempotency-Key: <uuid>

// Response 200
{
  "message": "Lesson started!",
  "progress": {
    "id": "cmqer18kq022hgcujg8lyhhyy",
    "user_id": "cmqer18kq009hgcujg8lyhhyh",
    "lesson_id": "cmqer18kq005hgcujg8lyhhyc",
    "status": "started",
    "started_at": "2026-06-15T12:00:00.000Z",
    "progress_percentage": 0,
    "notes": "",
    "created_at": "2026-06-15T12:00:00.000Z",
    "updated_at": "2026-06-15T12:00:00.000Z"
  }
}
```

#### PUT /api/lessons/:id/progress

```json
// Headers
Authorization: Bearer <access_token>
Idempotency-Key: <uuid>

// Request
{
  "progress_percentage": 50,
  "notes": "Saya sudah sampai di bagian middleware"
}
// Response 200 — returns updated progress object
{
  "id": "cmqer18kq022hgcujg8lyhhyy",
  "user_id": "cmqer18kq009hgcujg8lyhhyh",
  "lesson_id": "cmqer18kq005hgcujg8lyhhyc",
  "status": "started",
  "started_at": "2026-06-15T12:00:00.000Z",
  "completed_at": null,
  "progress_percentage": 50,
  "notes": "Saya sudah sampai di bagian middleware",
  "created_at": "2026-06-15T12:00:00.000Z",
  "updated_at": "2026-06-15T12:30:00.000Z"
}
```

#### POST /api/lessons/:id/complete

```json
// Headers
Authorization: Bearer <access_token>
Idempotency-Key: <uuid>

// Response 200
{
  "message": "Lesson completed!",
  "certificate": {
    "id": "cmqer18kq023hgcujg8lyhhyz",
    "unique_code": "CERT-abc12345",
    "issued_at": "2026-06-15T12:00:00.000Z",
    "user_id": "cmqer18kq009hgcujg8lyhhyh",
    "lesson_id": "cmqer18kq005hgcujg8lyhhyc",
    "lesson_name": "Go Basics",
    "user_name": "John Doe",
    "achievement": {
      "type": "certificate"
    }
  },
  "progress": {
    "id": "cmqer18kq022hgcujg8lyhhyy",
    "user_id": "cmqer18kq009hgcujg8lyhhyh",
    "lesson_id": "cmqer18kq005hgcujg8lyhhyc",
    "status": "completed",
    "started_at": "2026-06-15T12:00:00.000Z",
    "completed_at": "2026-06-15T12:00:00.000Z",
    "progress_percentage": 100,
    "notes": "",
    "created_at": "2026-06-15T12:00:00.000Z",
    "updated_at": "2026-06-15T12:00:00.000Z"
  },
  "next_lesson": {
    "id": "cmqer18kq007hgcujg8lyhhyf",
    "title": "Fiber Framework & Routing"
  },
  "points_awarded": 50
}
```

### Certificates

| Method | Path                             | Middleware                |
| ------ | -------------------------------- | ------------------------- |
| GET    | /api/users/me/certificates       | JWT + Idempotency (no XSRF) |
| GET    | /api/users/me/certificates/:code | JWT + Idempotency (no XSRF) |

#### GET /api/users/me/certificates

```json
// Headers
Authorization: Bearer <access_token>
Idempotency-Key: <uuid>

// Query: ?page=1&limit=20
// Response 200
{
  "data": [
    {
      "id": "cmqer18kq023hgcujg8lyhhyz",
      "unique_code": "CERT-abc12345",
      "issued_at": "2026-06-15T12:00:00.000Z",
      "user_id": "cmqer18kq009hgcujg8lyhhyh",
      "lesson_id": "cmqer18kq005hgcujg8lyhhyc",
      "lesson_name": "Go Basics",
      "user_name": "John Doe",
      "achievement": {
        "type": "certificate"
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 5
  }
}
```

#### GET /api/users/me/certificates/:code

```json
// Headers
Authorization: Bearer <access_token>
Idempotency-Key: <uuid>

// Response 200
{
  "id": "cmqer18kq023hgcujg8lyhhyz",
  "unique_code": "CERT-abc12345",
  "issued_at": "2026-06-15T12:00:00.000Z",
  "user_id": "cmqer18kq009hgcujg8lyhhyh",
  "lesson_id": "cmqer18kq005hgcujg8lyhhyc",
  "lesson_name": "Go Basics",
  "user_name": "John Doe",
  "achievement": {
    "type": "certificate",
    "title": "Go Basics",
    "unique_code": "CERT-abc12345"
  }
}
// Response 403
{
  "error": "Forbidden: not certificate owner"
}
// Response 404
{
  "error": "Certificate not found"
}
```

### Showcases

| Method | Path                    | Middleware                |
| ------ | ----------------------- | ------------------------- |
| POST   | /api/showcases          | JWT+XSRF+Idempotency      |
| PUT    | /api/showcases/:id      | JWT+XSRF+Idempotency      |
| DELETE | /api/showcases/:id      | JWT+XSRF+Idempotency      |
| POST   | /api/showcases/:id/like | JWT+XSRF+Idempotency      |
| DELETE | /api/showcases/:id/like | JWT+XSRF+Idempotency      |
| GET    | /api/users/me/showcases | JWT + Idempotency (no XSRF) |

#### POST /api/showcases

```json
// Headers
Authorization: Bearer <access_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Request
{
  "title": "My Portfolio Website",
  "description": "Built with React & Next.js",
  "media_urls": ["https://cdn.mohagussetiaone.my.id/jvalleyverse/showcases/img1.jpg"],
  "category_id": "cmqer18kq019hgcujg8lyhhyv",
  "visibility": "public"
}
// Response 201
{
  "id": "cmqer18kq018hgcujg8lyhhyq",
  "created_at": "2026-06-18T10:00:00.000Z",
  "updated_at": "2026-06-18T10:00:00.000Z",
  "title": "My Portfolio Website",
  "description": "Built with React & Next.js",
  "media_urls": ["https://cdn.mohagussetiaone.my.id/jvalleyverse/showcases/img1.jpg"],
  "user_id": "cmqer18kq009hgcujg8lyhhyh",
  "category_id": "cmqer18kq019hgcujg8lyhhyv",
  "status": "published",
  "visibility": "public",
  "likes_count": 0,
  "views_count": 0
}
```

#### PUT /api/showcases/:id

```json
// Headers
Authorization: Bearer <access_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Request (all fields optional)
{
  "title": "Updated Portfolio Title",
  "description": "Updated description with new features",
  "visibility": "public"
}
// Response 200 — returns updated showcase object
{
  "id": "cmqer18kq018hgcujg8lyhhyq",
  "created_at": "2026-06-15T12:00:00.000Z",
  "updated_at": "2026-06-18T10:00:00.000Z",
  "title": "Updated Portfolio Title",
  "description": "Updated description with new features",
  "media_urls": ["https://cdn.mohagussetiaone.my.id/jvalleyverse/showcases/img1.jpg"],
  "user_id": "cmqer18kq009hgcujg8lyhhyh",
  "category_id": "cmqer18kq019hgcujg8lyhhyv",
  "status": "published",
  "visibility": "public",
  "likes_count": 5,
  "views_count": 50
}
// Response 403
{
  "error": "You do not own this showcase"
}
```

#### DELETE /api/showcases/:id

```json
// Headers
Authorization: Bearer <access_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Response 200
{
  "message": "Showcase deleted"
}
// Response 403
{
  "error": "You do not own this showcase"
}
```

#### POST /api/showcases/:id/like

```json
// Headers
Authorization: Bearer <access_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Response 200
{
  "message": "Liked successfully"
}
```

#### DELETE /api/showcases/:id/like

```json
// Headers
Authorization: Bearer <access_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Response 200
{
  "message": "Showcase unliked"
}
```

#### GET /api/users/me/showcases

```json
// Headers
Authorization: Bearer <access_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Query: ?page=1&limit=20
// Response 200
{
  "data": [
    {
      "id": "cmqer18kq018hgcujg8lyhhyq",
      "title": "My Portfolio Website",
      "description": "Built with React & Next.js",
      "media_urls": ["https://cdn.mohagussetiaone.my.id/jvalleyverse/showcases/img1.jpg"],
      "likes_count": 24,
      "views_count": 156,
      "visibility": "public",
      "user": {
        "id": "cmqer18kq009hgcujg8lyhhyh",
        "name": "John Doe",
        "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/user.jpg"
      },
      "category": {
        "id": "cmqer18kq019hgcujg8lyhhyv",
        "name": "Frontend Development",
        "slug": "frontend-development"
      },
      "created_at": "2026-06-15T12:04:44.714Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 8
  }
}
```

### Discussions & Replies

| Method | Path                         | Middleware           |
| ------ | ---------------------------- | -------------------- |
| POST   | /api/discussions             | JWT+XSRF+Idempotency |
| PUT    | /api/discussions/:id         | JWT+XSRF+Idempotency |
| DELETE | /api/discussions/:id         | JWT+XSRF+Idempotency |
| POST   | /api/discussions/:id/close   | JWT+XSRF+Idempotency |
| POST   | /api/discussions/:id/replies | JWT+XSRF+Idempotency |
| PUT    | /api/replies/:id             | JWT+XSRF+Idempotency |
| DELETE | /api/replies/:id             | JWT+XSRF+Idempotency |
| POST   | /api/replies/:id/like        | JWT+XSRF+Idempotency |
| POST   | /api/replies/:id/best        | JWT+XSRF+Idempotency |

#### POST /api/discussions

```json
// Headers
Authorization: Bearer <access_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Request
{
  "title": "Perbedaan goroutine dan thread biasa?",
  "content": "Saya baru belajar Go, apakah ada yang bisa jelaskan perbedaan goroutine dengan thread OS?",
  "lesson_id": "cmqer18kq005hgcujg8lyhhyc",
  "study_case_id": null,
  "category_id": "cmqer18kq001hgcujg8lyhhvy"
}
// Response 201 — returns created discussion object
{
  "id": "cmqer18kq013hgcujg8lyhhyl",
  "created_at": "2026-06-18T10:00:00.000Z",
  "updated_at": "2026-06-18T10:00:00.000Z",
  "title": "Perbedaan goroutine dan thread biasa?",
  "content": "Saya baru belajar Go, apakah ada yang bisa jelaskan perbedaan goroutine dengan thread OS?",
  "user_id": "cmqer18kq009hgcujg8lyhhyh",
  "lesson_id": "cmqer18kq005hgcujg8lyhhyc",
  "study_case_id": null,
  "category_id": "cmqer18kq001hgcujg8lyhhvy",
  "views_count": 0,
  "status": "open",
  "is_pinned": false
}
```

#### PUT /api/discussions/:id

```json
// Headers
Authorization: Bearer <access_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Request
{
  "title": "Updated Goroutine Question",
  "content": "Updated content with more details..."
}
// Response 200
{
  "message": "Discussion updated"
}
// Response 403
{
  "error": "You do not own this discussion"
}
```

#### DELETE /api/discussions/:id

```json
// Headers
Authorization: Bearer <access_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Response 200
{
  "message": "Discussion deleted"
}
// Response 403
{
  "error": "You do not own this discussion"
}
```

#### POST /api/discussions/:id/close

```json
// Headers
Authorization: Bearer <access_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Response 200
{
  "message": "Discussion closed"
}
// Response 403
{
  "error": "You do not own this discussion"
}
```

#### POST /api/discussions/:id/replies

```json
// Headers
Authorization: Bearer <access_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Request (parent_id optional — for nested replies)
{
  "content": "Goroutine lebih ringan dari thread OS karena menggunakan scheduler sendiri yang berjalan di userspace",
  "parent_id": null
}
// Response 201 — returns created reply object
{
  "id": "cmqer18kq014hgcujg8lyhhym",
  "created_at": "2026-06-18T10:00:00.000Z",
  "updated_at": "2026-06-18T10:00:00.000Z",
  "content": "Goroutine lebih ringan dari thread OS karena menggunakan scheduler sendiri yang berjalan di userspace",
  "user_id": "cmqer18kq012hgcujg8lyhhyk",
  "discussion_id": "cmqer18kq013hgcujg8lyhhyl",
  "parent_id": null,
  "likes_count": 0,
  "is_marked_best": false
}
```

#### PUT /api/replies/:id

```json
// Headers
Authorization: Bearer <access_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Request
{
  "content": "Updated reply content with better explanation"
}
// Response 200
{
  "message": "Reply updated"
}
// Response 403
{
  "error": "You do not own this reply"
}
```

#### DELETE /api/replies/:id

```json
// Headers
Authorization: Bearer <access_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Response 200
{
  "message": "Reply deleted"
}
// Response 403
{
  "error": "You do not own this reply"
}
```

#### POST /api/replies/:id/like

```json
// Headers
Authorization: Bearer <access_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Response 200
{
  "message": "Reply liked"
}
// Response 404
{
  "error": "Reply not found"
}
```

#### POST /api/replies/:id/best

```json
// Headers
Authorization: Bearer <access_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Request
{
  "discussion_id": "cmqer18kq013hgcujg8lyhhyl"
}
// Response 200
{
  "message": "Reply marked as best answer"
}
// Response 403
{
  "error": "Only discussion owner can mark best answer"
}
```

### Reviews

| Method | Path             | Middleware           |
| ------ | ---------------- | -------------------- |
| POST   | /api/reviews     | JWT+XSRF+Idempotency |
| PUT    | /api/reviews/:id | JWT+XSRF+Idempotency |
| DELETE | /api/reviews/:id | JWT+XSRF+Idempotency |

#### POST /api/reviews

```json
// Headers
Authorization: Bearer <access_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Request
{
  "course_id": "cmqer18kq000hgcujg8lyhhux",
  "lesson_id": "",
  "rating": 5,
  "message": "Kursusnya sangat bagus dan mudah dipahami"
}
// Response 201 — returns created domain.Review (without preloaded user)
{
  "id": "cmqer18kq008hgcujg8lyhhyg",
  "created_at": "2026-06-18T10:00:00.000Z",
  "updated_at": "2026-06-18T10:00:00.000Z",
  "user_id": "cmqer18kq009hgcujg8lyhhyh",
  "course_id": "cmqer18kq000hgcujg8lyhhux",
  "lesson_id": "",
  "rating": 5,
  "message": "Kursusnya sangat bagus dan mudah dipahami"
}
// Response 400
{
  "error": "rating (1-5) and message are required"
}
```

#### PUT /api/reviews/:id

```json
// Headers
Authorization: Bearer <access_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Request
{
  "rating": 4,
  "message": "Updated review: materinya bagus tapi butuh lebih banyak latihan"
}
// Response 200 — returns updated domain.Review (without preloaded user)
{
  "id": "cmqer18kq008hgcujg8lyhhyg",
  "created_at": "2026-06-15T12:00:00.000Z",
  "updated_at": "2026-06-18T10:00:00.000Z",
  "user_id": "cmqer18kq009hgcujg8lyhhyh",
  "course_id": "cmqer18kq000hgcujg8lyhhux",
  "lesson_id": "",
  "rating": 4,
  "message": "Updated review: materinya bagus tapi butuh lebih banyak latihan"
}
```

#### DELETE /api/reviews/:id

```json
// Headers
Authorization: Bearer <access_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Response 200
{
  "message": "Review deleted"
}
```

### Enrollment & My Courses

| Method | Path                           | Middleware                |
| ------ | ------------------------------ | ------------------------- |
| POST   | /api/courses/:id/enroll        | JWT + Idempotency (no XSRF) |
| GET    | /api/users/me/courses          | JWT + Idempotency (no XSRF) |
| PUT    | /api/courses/:id/last-lesson   | JWT + Idempotency (no XSRF) |

#### POST /api/courses/:id/enroll

```json
// Headers
Authorization: Bearer <access_token>
Idempotency-Key: <uuid>

// Response 201
{
  "message": "Successfully enrolled in course"
}
```

#### GET /api/users/me/courses

```json
// Headers
Authorization: Bearer <access_token>
Idempotency-Key: <uuid>

// Query: ?page=1&limit=10
// Response 200
{
  "data": [
    {
      "id": "cmqer18kq000hgcujg8lyhhux",
      "title": "Membangun REST API dengan Go & Fiber",
      "description": "Pelajari cara membuat RESTful API yang scalable menggunakan Go dan Fiber framework",
      "thumbnail": "https://cdn.mohagussetiaone.my.id/jvalleyverse/courses/thumb.jpg",
      "price": 0,
      "category": {
        "id": "cmqer18kq001hgcujg8lyhhvy",
        "name": "Backend Development",
        "slug": "backend-development"
      },
      "admin_name": "JValley Admin",
      "mentor": {
        "id": "cmqer18kq003hgcujg8lyhhxa",
        "name": "Budi Mentor",
        "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/mentor.jpg",
        "role": "mentor"
      },
      "hours": 40,
      "visibility": "public",
      "section_count": 2,
      "is_enrolled": true,
      "created_at": "2026-06-15T12:04:44.714Z",
      "enrolled_at": "2026-06-16T08:00:00.000Z",
      "last_lesson_id": "cmqer18kq010hgcujg8lyiiab"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 3
  }
}
```

#### PUT /api/courses/:id/last-lesson

```json
// Request
{
  "lesson_id": "cmqer18kq010hgcujg8lyiiab"
}

// Response 200
{
  "message": "Last lesson updated"
}
```

### My Items

| Method | Path                      | Middleware                |
| ------ | ------------------------- | ------------------------- |
| GET    | /api/users/me/study-cases | JWT + Idempotency (no XSRF) |
| GET    | /api/users/me/blogs       | JWT + Idempotency (no XSRF) |
| GET    | /api/users/me/replies     | JWT + Idempotency (no XSRF) |
| GET    | /api/users/me/discussions | JWT + Idempotency (no XSRF) |

#### GET /api/users/me/study-cases

```json
// Headers
Authorization: Bearer <access_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Query: ?page=1&limit=20
// Response 200
{
  "data": [
    {
      "id": "cmqer18kq016hgcujg8lyhhyo",
      "name": "Belajar Go Basic",
      "description": "Studi kasus fundamental Go untuk pemula",
      "img_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/study-cases/img.jpg",
      "youtube_url": "https://youtube.com/watch?v=abc123",
      "tags": ["golang", "beginner"],
      "user": {
        "id": "cmqer18kq012hgcujg8lyhhyk",
        "name": "JValley Admin",
        "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/admin.jpg"
      },
      "created_at": "2026-06-15T12:04:44.714Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 3
  }
}
```

#### GET /api/users/me/blogs

```json
// Headers
Authorization: Bearer <access_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Query: ?page=1&limit=10&status=draft
// Response 200
{
  "data": [
    {
      "id": "cmqer18kq011hgcujg8lyhhyj",
      "title": "Belajar Go untuk Pemula",
      "slug": "belajar-go-untuk-pemula",
      "description": "Panduan lengkap belajar Go dari dasar hingga mahir untuk pemula",
      "cover_img_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/blogs/cover.jpg",
      "tags": ["golang", "backend"],
      "status": "published",
      "author": {
        "id": "cmqer18kq012hgcujg8lyhhyk",
        "name": "JValley Admin",
        "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/admin.jpg"
      },
      "category": {
        "id": "cmqer18kq001hgcujg8lyhhvy",
        "name": "Backend Development",
        "slug": "backend-development"
      },
      "created_at": "2026-06-15T12:00:00.000Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 5
  }
}
```

#### GET /api/users/me/replies

```json
// Headers
Authorization: Bearer <access_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Query: ?page=1&limit=20
// Response 200
{
  "data": [
    {
      "id": "cmqer18kq014hgcujg8lyhhym",
      "content": "Goroutine lebih ringan dari thread OS karena menggunakan scheduler sendiri yang berjalan di userspace",
      "discussion_id": "cmqer18kq013hgcujg8lyhhyl",
      "discussion_title": "Perbedaan goroutine dan thread biasa?",
      "parent_id": null,
      "likes_count": 5,
      "is_marked_best": true,
      "created_at": "2026-06-15T13:00:00.000Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 15
  }
}
```

### File Upload

| Method | Path        | Middleware                |
| ------ | ----------- | ------------------------- |
| POST   | /api/upload | JWT + Idempotency (no XSRF) |

**Request:** `multipart/form-data`

- `file` — File (max 10MB, ekstensi: jpg,jpeg,png,gif,webp,svg,mp4,pdf,zip)
- `folder` — String: `courses`, `lessons`, `avatars`, `showcases`, `study-cases`, `blogs`

```javascript
const formData = new FormData();
formData.append("file", fileInput.files[0]);
formData.append("folder", "courses");

fetch("/api/upload", {
  method: "POST",
  headers: {
    Authorization: "Bearer <token>",
    "Idempotency-Key": crypto.randomUUID(),
  },
  body: formData,
});
```

```json
// Response 201
{
  "url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/courses/a1b2c3d4-e5f6-7890-abcd-ef1234567890.jpg",
  "object_name": "courses/a1b2c3d4-e5f6-7890-abcd-ef1234567890.jpg",
  "size": 245760,
  "content_type": "image/jpeg"
}
// Response 400
{
  "error": "file too large, maximum size is 10 MB"
}
// Response 503
{
  "error": "File upload is not available (MinIO not configured)"
}
```

### Notifications

| Method | Path                                 | Middleware                |
| ------ | ------------------------------------ | ------------------------- |
| GET    | /api/notifications/stream            | JWT + Idempotency (no XSRF) |
| GET    | /api/users/me/notifications          | JWT + Idempotency (no XSRF) |
| GET    | /api/users/me/notifications/count    | JWT + Idempotency (no XSRF) |
| PUT    | /api/users/me/notifications/:id/read | JWT + Idempotency (no XSRF) |
| PUT    | /api/users/me/notifications/read-all | JWT + Idempotency (no XSRF) |
| DELETE | /api/users/me/notifications/:id      | JWT + Idempotency (no XSRF) |

#### Tipe Notifikasi (Otomatis Oleh Sistem)

Sistem membuat notifikasi secara otomatis untuk 12 event berbeda:

| Tipe                 | Pemicu                                           | Diterima Oleh         |
| -------------------- | ------------------------------------------------ | --------------------- |
| `new_reply`          | Reply baru di diskusi                            | Owner diskusi         |
| `nested_reply`       | Balasan nested ke reply                          | Parent reply owner    |
| `reply_like`         | Reply seseorang di-like                          | Creator reply         |
| `best_answer`        | Reply ditandai sebagai jawaban terbaik           | Creator reply         |
| `showcase_like`      | Showcase di-like                                 | Owner showcase        |
| `course_enrollment`  | User baru mendaftar ke kursus                    | Admin course          |
| `enrollment_success` | Pendaftaran kursus berhasil                      | User yang mendaftar   |
| `new_review`         | Review baru untuk kursus                         | Admin course          |
| `lesson_completed`   | Pelajaran selesai + sertifikat didapat           | User yang belajar     |
| `level_up`           | Level naik (dengan badge dari user_levels)       | User yang naik level  |
| `blog_published`     | Blog diterbitkan                                 | Author blog           |
| `discussion_created` | Diskusi baru dibuat (terkait lesson)             | Creator diskusi       |

**Anti Self-Notifikasi:** Sistem tidak mengirim notifikasi jika aktor == penerima (self-reply, self-like, self-review dilewati).

#### GET /api/notifications/stream (SSE)

```text
event: connected
data: {"status":"connected","user_id":"cmqer18kq009hgcujg8lyhhyh"}

event: notification
data: {"type":"showcase_like","title":"Showcase Anda Mendapat Like","message":"Seseorang menyukai showcase Anda: My Portfolio Website","link":"/showcases/cmqer18kq018hgcujg8lyhhyq","payload":{"id":"cmqer18kq024hgcujg8lyhhz","created_at":"2026-06-18T10:00:00.000Z","is_read":false}}

: heartbeat
```

#### GET /api/users/me/notifications

```json
// Headers
Authorization: Bearer <access_token>
Idempotency-Key: <uuid>

// Query: ?page=1&limit=20
// Response 200
{
  "data": [
    {
      "id": "cmqer18kq024hgcujg8lyhhz",
      "type": "showcase_like",
      "title": "Showcase Anda Mendapat Like",
      "message": "Seseorang menyukai showcase Anda: My Portfolio Website",
      "is_read": false,
      "link": "/showcases/cmqer18kq018hgcujg8lyhhyq",
      "created_at": "2026-06-18T10:00:00.000Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 10
  },
  "unread_count": 3
}
```

#### GET /api/users/me/notifications/count

```json
// Headers
Authorization: Bearer <access_token>
Idempotency-Key: <uuid>

// Response 200
{
  "unread_count": 3
}
```

#### PUT /api/users/me/notifications/:id/read

```json
// Headers
Authorization: Bearer <access_token>
Idempotency-Key: <uuid>

// Response 200
{
  "message": "Notification marked as read"
}
```

#### PUT /api/users/me/notifications/read-all

```json
// Headers
Authorization: Bearer <access_token>
Idempotency-Key: <uuid>

// Response 200
{
  "message": "All notifications marked as read"
}
```

#### DELETE /api/users/me/notifications/:id

```json
// Headers
Authorization: Bearer <access_token>
Idempotency-Key: <uuid>

// Response 200
{
  "message": "Notification deleted"
}
```

### Gamification

| Method | Path                  | Middleware                |
| ------ | --------------------- | ------------------------- |
| GET    | /api/levels           | JWT + Idempotency (no XSRF) |
| GET    | /api/users/:id/points | JWT + Idempotency (no XSRF) |

#### GET /api/levels

```json
// Headers
Authorization: Bearer <access_token>
Idempotency-Key: <uuid>

// Response 200
{
  "data": [
    {
      "name": "Beginner",
      "threshold": 0,
      "color": "#6366f1",
      "description": "Just starting your journey"
    },
    {
      "name": "Intermediate",
      "threshold": 100,
      "color": "#8b5cf6",
      "description": "Building momentum"
    },
    {
      "name": "Advanced",
      "threshold": 500,
      "color": "#d946ef",
      "description": "Getting serious"
    },
    {
      "name": "Expert",
      "threshold": 1000,
      "color": "#ec4899",
      "description": "Mastering skills"
    },
    {
      "name": "Master",
      "threshold": 2000,
      "color": "#f43f5e",
      "description": "Peak achievement"
    }
  ]
}
```

#### GET /api/users/:id/points

```json
// Headers
Authorization: Bearer <access_token>
Idempotency-Key: <uuid>

// Response 200
{
  "user_id": "cmqer18kq009hgcujg8lyhhyh",
  "name": "John Doe",
  "total_points": 850,
  "current_level": 3,
  "recent_activity": [
    {
      "id": "cmqer18kq020hgcujg8lyhhyw",
      "activity": "create_showcase",
      "points": 10,
      "timestamp": "2026-06-15T12:04:44.714Z"
    }
  ]
}
// Response 404
{
  "error": "User not found"
}
```

---

## ADMIN ROUTES (JWT + RequireRole("admin"))

Group `/api/admin` — Semua route di bawah ini perlu **JWT + XSRF + Idempotency + role=admin**.

### Dashboard

| Method | Path                 |
| ------ | -------------------- |
| GET    | /api/admin/dashboard |

#### GET /api/admin/dashboard

```json
// Headers
Authorization: Bearer <admin_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Response 200
{
  "message": "Welcome admin"
}
```

### Users

| Method | Path             |
| ------ | ---------------- |
| GET    | /api/admin/users |

#### GET /api/admin/users

```json
// Headers
Authorization: Bearer <admin_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Query: ?page=1&limit=20
// Response 200
{
  "data": [
    {
      "id": "cmqer18kq012hgcujg8lyhhyk",
      "email": "admin@jvalleyverse.com",
      "name": "JValley Admin",
      "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/admin.jpg",
      "role": "admin",
      "level": 5,
      "total_points": 5000,
      "is_active": true,
      "created_at": "2026-01-15T08:00:00.000Z"
    },
    {
      "id": "cmqer18kq009hgcujg8lyhhyh",
      "email": "john@example.com",
      "name": "John Doe",
      "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/user.jpg",
      "role": "user",
      "level": 2,
      "total_points": 450,
      "is_active": true,
      "created_at": "2026-01-20T10:30:00.000Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 156
  }
}
```

### Blogs

| Method | Path                 |
| ------ | -------------------- |
| GET    | /api/admin/blogs     |
| POST   | /api/admin/blogs     |
| PUT    | /api/admin/blogs/:id |
| DELETE | /api/admin/blogs/:id |

#### GET /api/admin/blogs

```json
// Headers
Authorization: Bearer <admin_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Query: ?page=1&limit=10&search=golang&category_id=cmq...&tag=golang
// Response 200 — same structure as GET /api/blogs but includes all statuses (not only published)
{
  "data": [
    {
      "id": "cmqer18kq011hgcujg8lyhhyj",
      "title": "Belajar Go untuk Pemula",
      "slug": "belajar-go-untuk-pemula",
      "description": "Panduan lengkap belajar Go dari dasar hingga mahir untuk pemula",
      "cover_img_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/blogs/cover.jpg",
      "tags": ["golang", "backend"],
      "status": "draft",
      "author": {
        "id": "cmqer18kq012hgcujg8lyhhyk",
        "name": "JValley Admin",
        "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/admin.jpg"
      },
      "category": {
        "id": "cmqer18kq001hgcujg8lyhhvy",
        "name": "Backend Development",
        "slug": "backend-development"
      },
      "created_at": "2026-06-15T12:00:00.000Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 8
  }
}
```

#### POST /api/admin/blogs

```json
// Headers
Authorization: Bearer <admin_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Request
{
  "title": "Belajar Go untuk Pemula",
  "description": "Panduan lengkap belajar Go dari dasar hingga mahir untuk pemula",
  "content": "# Full markdown content here...\n\n## Introduction\n...",
  "cover_img_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/blogs/cover.jpg",
  "tags": ["golang", "backend"],
  "status": "draft",
  "category_id": "cmqer18kq001hgcujg8lyhhvy"
}
// Response 201 — returns BlogDetail
{
  "id": "cmqer18kq011hgcujg8lyhhyj",
  "title": "Belajar Go untuk Pemula",
  "slug": "belajar-go-untuk-pemula",
  "description": "Panduan lengkap belajar Go dari dasar hingga mahir untuk pemula",
  "content": "# Full markdown content here...\n\n## Introduction\n...",
  "cover_img_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/blogs/cover.jpg",
  "tags": ["golang", "backend"],
  "status": "draft",
  "author": {
    "id": "cmqer18kq012hgcujg8lyhhyk",
    "name": "JValley Admin",
    "avatar": "https://cdn.mohagussetiaone.my.id/jvalleyverse/avatars/admin.jpg"
  },
  "category": {
    "id": "cmqer18kq001hgcujg8lyhhvy",
    "name": "Backend Development",
    "slug": "backend-development"
  },
  "created_at": "2026-06-18T10:00:00.000Z",
  "updated_at": "2026-06-18T10:00:00.000Z"
}
```

#### PUT /api/admin/blogs/:id

```json
// Headers
Authorization: Bearer <admin_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Request (all fields optional)
{
  "title": "Updated Blog Title",
  "description": "Updated description",
  "content": "Updated full markdown content...",
  "cover_img_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/blogs/new-cover.jpg",
  "tags": ["golang", "backend", "tutorial"],
  "status": "published",
  "category_id": "cmqer18kq001hgcujg8lyhhvy"
}
// Response 200
{
  "message": "Blog updated successfully"
}
```

#### DELETE /api/admin/blogs/:id

```json
// Headers
Authorization: Bearer <admin_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Response 200
{
  "message": "Blog deleted successfully"
}
```

### Courses

| Method | Path                   |
| ------ | ---------------------- |
| POST   | /api/admin/courses     |
| PUT    | /api/admin/courses/:id |
| DELETE | /api/admin/courses/:id |

#### POST /api/admin/courses

```json
// Headers
Authorization: Bearer <admin_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Request
{
  "title": "Membangun REST API dengan Go & Fiber",
  "description": "Pelajari cara membuat RESTful API yang scalable menggunakan Go dan Fiber framework",
  "thumbnail": "https://cdn.mohagussetiaone.my.id/jvalleyverse/courses/thumb.jpg",
  "category_id": "cmqer18kq001hgcujg8lyhhvy",
  "mentor_id": "cmqer18kq003hgcujg8lyhhxa",
  "price": 0,
  "hours": 40,
  "learning_objectives": ["Memahami dasar Go", "Membuat REST API", "Integrasi database"]
}
// Response 201 — returns created course object (associations not preloaded after Create)
{
  "id": "cmqer18kq000hgcujg8lyhhux",
  "created_at": "2026-06-18T10:00:00.000Z",
  "updated_at": "2026-06-18T10:00:00.000Z",
  "title": "Membangun REST API dengan Go & Fiber",
  "description": "Pelajari cara membuat RESTful API yang scalable menggunakan Go dan Fiber framework",
  "thumbnail": "https://cdn.mohagussetiaone.my.id/jvalleyverse/courses/thumb.jpg",
  "category_id": "cmqer18kq001hgcujg8lyhhvy",
  "category": {
    "id": "",
    "name": "",
    "slug": ""
  },
  "admin_id": "cmqer18kq012hgcujg8lyhhyk",
  "admin": {
    "id": "",
    "created_at": "0001-01-01T00:00:00.000Z",
    "updated_at": "0001-01-01T00:00:00.000Z",
    "email": "",
    "name": "",
    "avatar": "",
    "bio": "",
    "role": "",
    "is_active": false,
    "points": 0,
    "total_points": 0,
    "level": 0
  },
  "mentor_id": "cmqer18kq003hgcujg8lyhhxa",
  "mentor": null,
  "visibility": "public",
  "price": 0,
  "hours": 40,
  "learning_objectives": ["Memahami dasar Go", "Membuat REST API", "Integrasi database"]
}
```

#### PUT /api/admin/courses/:id

```json
// Headers
Authorization: Bearer <admin_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Request
{
  "title": "Updated Course Title",
  "description": "Updated description",
  "price": 99000,
  "visibility": "public",
  "learning_objectives": ["Objective 1", "Objective 2"]
}
// Response 200
{
  "message": "Course updated"
}
```

#### DELETE /api/admin/courses/:id

```json
// Headers
Authorization: Bearer <admin_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Response 200
{
  "message": "Course deleted"
}
```

### Sections

| Method | Path                                   |
| ------ | -------------------------------------- |
| POST   | /api/admin/courses/:course_id/sections |
| PUT    | /api/admin/sections/:section_id        |
| DELETE | /api/admin/sections/:section_id        |

#### POST /api/admin/courses/:course_id/sections

```json
// Headers
Authorization: Bearer <admin_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Request
{
  "title": "Introduction to Go",
  "description": "Dasar-dasar bahasa pemrograman Go",
  "order_index": 1
}
// Response 201 — returns created section object
{
  "id": "cmqer18kq004hgcujg8lyhhyb",
  "created_at": "2026-06-18T10:00:00.000Z",
  "updated_at": "2026-06-18T10:00:00.000Z",
  "title": "Introduction to Go",
  "description": "Dasar-dasar bahasa pemrograman Go",
  "course_id": "cmqer18kq000hgcujg8lyhhux",
  "order_index": 1
}
```

#### PUT /api/admin/sections/:section_id

```json
// Headers
Authorization: Bearer <admin_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Request
{
  "title": "Updated Section Title",
  "description": "Updated description",
  "order_index": 2
}
// Response 200 — returns updated section object
{
  "id": "cmqer18kq004hgcujg8lyhhyb",
  "created_at": "2026-06-15T12:00:00.000Z",
  "updated_at": "2026-06-18T10:00:00.000Z",
  "title": "Updated Section Title",
  "description": "Updated description",
  "course_id": "cmqer18kq000hgcujg8lyhhux",
  "order_index": 2
}
```

#### DELETE /api/admin/sections/:section_id

```json
// Headers
Authorization: Bearer <admin_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Response 200
{
  "message": "Section deleted"
}
```

### Lessons

| Method | Path                           |
| ------ | ------------------------------ |
| POST   | /api/admin/lessons             |
| POST   | /api/admin/lessons/:id/details |
| PUT    | /api/admin/lessons/:id         |
| DELETE | /api/admin/lessons/:id         |

#### POST /api/admin/lessons

```json
// Headers
Authorization: Bearer <admin_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Request
{
  "course_id": "cmqer18kq000hgcujg8lyhhux",
  "section_id": "cmqer18kq004hgcujg8lyhhyb",
  "title": "Go Basics",
  "slug": "go-basics",
  "description": "Pelajari fundamentals Go: variables, types, control flow",
  "thumbnail": "https://cdn.mohagussetiaone.my.id/jvalleyverse/lessons/thumb.jpg",
  "difficulty": "beginner",
  "duration": 45,
  "order_index": 1,
  "video_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/lessons/video.mp4",
  "visibility": "public"
}
// Response 201 — returns created lesson object
{
  "id": "cmqer18kq005hgcujg8lyhhyc",
  "created_at": "2026-06-18T10:00:00.000Z",
  "updated_at": "2026-06-18T10:00:00.000Z",
  "course_id": "cmqer18kq000hgcujg8lyhhux",
  "section_id": "cmqer18kq004hgcujg8lyhhyb",
  "title": "Go Basics",
  "slug": "go-basics",
  "description": "Pelajari fundamentals Go: variables, types, control flow",
  "thumbnail": "https://cdn.mohagussetiaone.my.id/jvalleyverse/lessons/thumb.jpg",
  "admin_id": "cmqer18kq012hgcujg8lyhhyk",
  "difficulty": "beginner",
  "duration": 45,
  "order_index": 1,
  "video_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/lessons/video.mp4",
  "visibility": "public"
}
```

#### POST /api/admin/lessons/:id/details

```json
// Headers
Authorization: Bearer <admin_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Request
{
  "about": "Pelajari fundamentals Go: variables, types, control flow",
  "rules": "1. Selesaikan semua latihan\n2. Submit code kamu\n3. Diskusi jika ada kendala",
  "tools": ["Go 1.22+", "VS Code", "Terminal"],
  "resource_media": {
    "videos": ["https://cdn.mohagussetiaone.my.id/jvalleyverse/lessons/intro.mp4"],
    "documents": ["https://cdn.mohagussetiaone.my.id/jvalleyverse/lessons/guide.pdf"],
    "images": ["https://cdn.mohagussetiaone.my.id/jvalleyverse/lessons/diagram.png"]
  },
  "resources": [
    {"type": "article", "title": "Go Tour", "url": "https://go.dev/tour"},
    {"type": "github", "title": "Source Code", "url": "https://github.com/example/go-basics"}
  ]
}
// Response 201 — returns created LessonDetail
{
  "id": "cmqer18kq010hgcujg8lyhhyi",
  "lesson_id": "cmqer18kq005hgcujg8lyhhyc",
  "about": "Pelajari fundamentals Go: variables, types, control flow",
  "rules": "1. Selesaikan semua latihan\n2. Submit code kamu\n3. Diskusi jika ada kendala",
  "tools": ["Go 1.22+", "VS Code", "Terminal"],
  "resource_media": {
    "videos": ["https://cdn.mohagussetiaone.my.id/jvalleyverse/lessons/intro.mp4"],
    "documents": ["https://cdn.mohagussetiaone.my.id/jvalleyverse/lessons/guide.pdf"],
    "images": ["https://cdn.mohagussetiaone.my.id/jvalleyverse/lessons/diagram.png"]
  },
  "resources": [
    {"type": "article", "title": "Go Tour", "url": "https://go.dev/tour"},
    {"type": "github", "title": "Source Code", "url": "https://github.com/example/go-basics"}
  ],
  "created_at": "2026-06-18T10:00:00.000Z",
  "updated_at": "2026-06-18T10:00:00.000Z"
}
```

#### PUT /api/admin/lessons/:id

```json
// Headers
Authorization: Bearer <admin_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Request (all fields optional)
{
  "title": "Updated Lesson Title",
  "slug": "updated-lesson-title",
  "description": "Updated description with more details",
  "difficulty": "intermediate",
  "duration": 60,
  "video_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/lessons/updated.mp4",
  "visibility": "public"
}
// Response 200 — returns updated lesson object
{
  "id": "cmqer18kq005hgcujg8lyhhyc",
  "created_at": "2026-06-15T12:00:00.000Z",
  "updated_at": "2026-06-18T10:00:00.000Z",
  "course_id": "cmqer18kq000hgcujg8lyhhux",
  "section_id": "cmqer18kq004hgcujg8lyhhyb",
  "title": "Updated Lesson Title",
  "slug": "updated-lesson-title",
  "description": "Updated description with more details",
  "thumbnail": "https://cdn.mohagussetiaone.my.id/jvalleyverse/lessons/thumb.jpg",
  "admin_id": "cmqer18kq012hgcujg8lyhhyk",
  "difficulty": "intermediate",
  "duration": 60,
  "order_index": 1,
  "video_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/lessons/updated.mp4",
  "visibility": "public"
}
```

#### DELETE /api/admin/lessons/:id

```json
// Headers
Authorization: Bearer <admin_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Response 200
{
  "message": "Lesson deleted"
}
```

### Study Cases

| Method | Path                       |
| ------ | -------------------------- |
| POST   | /api/admin/study-cases     |
| PUT    | /api/admin/study-cases/:id |
| DELETE | /api/admin/study-cases/:id |

#### POST /api/admin/study-cases

```json
// Headers
Authorization: Bearer <admin_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Request
{
  "name": "Belajar Go Basic",
  "description": "Studi kasus fundamental Go untuk pemula",
  "img_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/study-cases/img.jpg",
  "youtube_url": "https://youtube.com/watch?v=abc123",
  "tags": ["golang", "beginner"]
}
// Response 201 — returns created study case object
{
  "id": "cmqer18kq016hgcujg8lyhhyo",
  "created_at": "2026-06-18T10:00:00.000Z",
  "updated_at": "2026-06-18T10:00:00.000Z",
  "name": "Belajar Go Basic",
  "description": "Studi kasus fundamental Go untuk pemula",
  "img_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/study-cases/img.jpg",
  "tags": ["golang", "beginner"],
  "youtube_url": "https://youtube.com/watch?v=abc123",
  "user_id": "cmqer18kq012hgcujg8lyhhyk"
}
```

#### PUT /api/admin/study-cases/:id

```json
// Headers
Authorization: Bearer <admin_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Request (all fields optional)
{
  "name": "Updated Study Case",
  "description": "Updated description",
  "img_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/study-cases/new-img.jpg",
  "youtube_url": "https://youtube.com/watch?v=def456",
  "tags": ["golang", "intermediate"]
}
// Response 200 — returns updated study case object
{
  "id": "cmqer18kq016hgcujg8lyhhyo",
  "created_at": "2026-06-15T12:00:00.000Z",
  "updated_at": "2026-06-18T10:00:00.000Z",
  "name": "Updated Study Case",
  "description": "Updated description",
  "img_url": "https://cdn.mohagussetiaone.my.id/jvalleyverse/study-cases/new-img.jpg",
  "tags": ["golang", "intermediate"],
  "youtube_url": "https://youtube.com/watch?v=def456",
  "user_id": "cmqer18kq012hgcujg8lyhhyk"
}
// Response 404
{
  "error": "Study case not found"
}
```

#### DELETE /api/admin/study-cases/:id

```json
// Headers
Authorization: Bearer <admin_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Response 200
{
  "message": "Study case deleted"
}
// Response 404
{
  "error": "Study case not found"
}
```

### FAQs

| Method | Path                |
| ------ | ------------------- |
| GET    | /api/admin/faqs     |
| GET    | /api/admin/faqs/:id |
| POST   | /api/admin/faqs     |
| PUT    | /api/admin/faqs/:id |
| DELETE | /api/admin/faqs/:id |

#### GET /api/admin/faqs

```json
// Headers
Authorization: Bearer <admin_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Query: ?page=1&limit=20
// Response 200
{
  "data": [
    {
      "id": "cmqer18kq026hgcujg8lyhhzb",
      "question": "Apa itu JValleyverse?",
      "answer": "JValleyverse adalah platform belajar coding online dengan sistem gamifikasi...",
      "category": "general",
      "order_index": 1,
      "is_active": true,
      "created_at": "2026-06-18T10:00:00.000Z",
      "updated_at": "2026-06-18T10:00:00.000Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 5
  }
}
```

#### POST /api/admin/faqs

```json
// Headers
Authorization: Bearer <admin_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Request
{
  "question": "Bagaimana cara mendaftar?",
  "answer": "Klik tombol Daftar di pojok kanan atas, isi email dan password, lalu verifikasi.",
  "category": "account",
  "order_index": 2
}
// Response 201
{
  "id": "cmqer18kq026hgcujg8lyhhzb",
  "question": "Bagaimana cara mendaftar?",
  "answer": "Klik tombol Daftar di pojok kanan atas, isi email dan password, lalu verifikasi.",
  "category": "account",
  "order_index": 2,
  "is_active": true,
  "created_at": "2026-06-18T10:00:00.000Z",
  "updated_at": "2026-06-18T10:00:00.000Z"
}
```

#### PUT /api/admin/faqs/:id

```json
// Headers
Authorization: Bearer <admin_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Request (all fields optional)
{
  "question": "Updated question",
  "answer": "Updated answer",
  "category": "general",
  "order_index": 1,
  "is_active": false
}
// Response 200 — returns updated FAQItem
{
  "id": "cmqer18kq026hgcujg8lyhhzb",
  "question": "Updated question",
  "answer": "Updated answer",
  "category": "general",
  "order_index": 1,
  "is_active": false,
  "created_at": "2026-06-18T10:00:00.000Z",
  "updated_at": "2026-06-18T10:00:00.000Z"
}
```

#### DELETE /api/admin/faqs/:id

```json
// Headers
Authorization: Bearer <admin_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Response 200
{
  "message": "FAQ deleted"
}
```

### Categories

| Method | Path                      |
| ------ | ------------------------- |
| POST   | /api/admin/categories     |
| GET    | /api/admin/categories     |
| PUT    | /api/admin/categories/:id |
| DELETE | /api/admin/categories/:id |

#### POST /api/admin/categories

```json
// Headers
Authorization: Bearer <admin_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Request
{
  "name": "Mobile Development",
  "slug": "mobile-development",
  "description": "Belajar membuat aplikasi mobile dengan Flutter dan React Native"
}
// Response 201 — returns created category object
{
  "id": "cmqer18kq025hgcujg8lyhhza",
  "created_at": "2026-06-18T10:00:00.000Z",
  "updated_at": "2026-06-18T10:00:00.000Z",
  "name": "Mobile Development",
  "slug": "mobile-development",
  "description": "Belajar membuat aplikasi mobile dengan Flutter dan React Native"
}
```

#### GET /api/admin/categories

Same as `GET /api/categories` (public) — returns all categories.

```json
// Response 200
[
  {
    "id": "cmqer18kq000hgcujg8lyhhux",
    "name": "Backend Development",
    "slug": "backend-development"
  }
]
```

#### PUT /api/admin/categories/:id

```json
// Headers
Authorization: Bearer <admin_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Request (all fields optional)
{
  "name": "Updated Category",
  "slug": "updated-category",
  "description": "Updated description"
}
// Response 200 — returns updated category object
{
  "id": "cmqer18kq000hgcujg8lyhhux",
  "created_at": "2026-06-15T12:04:44.714Z",
  "updated_at": "2026-06-18T10:00:00.000Z",
  "name": "Updated Category",
  "slug": "updated-category",
  "description": "Updated description"
}
```

#### DELETE /api/admin/categories/:id

```json
// Headers
Authorization: Bearer <admin_token>
X-XSRF-TOKEN: <xsrf_token>
Idempotency-Key: <uuid>

// Response 200
{
  "message": "Category deleted"
}
```

---

## Idempotency-Key (Safe Retry)

Semua mutation endpoint (POST/PUT/DELETE) mendukung `Idempotency-Key: <UUIDv4>` header untuk safe retry.

```bash
# Generate UUID
$ python -c "import uuid; print(uuid.uuid4())"
550e8400-e29b-41d4-a716-446655440000

# Request pertama — diproses normal
curl -X POST http://localhost:3000/api/showcases \
  -H "Authorization: Bearer <token>" \
  -H "Idempotency-Key: 550e8400-e29b-41d4-a716-446655440000" \
  -H "Content-Type: application/json" \
  -d '{"title":"My Project","category_id":"cmqer18kq019hgcujg8lyhhyv"}'

# Retry — response dari cache (X-Idempotency-Replayed: true)
```

---

## Rate Limiting

| Tier        | Limit          | Target        | Endpoints                                                        |
| ----------- | -------------- | ------------- | ---------------------------------------------------------------- |
| **Global**  | 200 req/min/IP | Baseline      | Semua route                                                      |
| **Content** | 60 req/min/IP  | Anti-scraping | Public GET: courses, lessons, showcases, categories, study-cases |
| **Auth**    | 10 req/min/IP  | Brute force   | POST /api/auth/login, /api/auth/register                         |

```json
// Response 429 Too Many Requests
{ "error": "Too many requests", "retry_after": 60 }
```

---

## Anti-Scraping (ScraperGuard)

Semua public content endpoint diproteksi **ScraperGuard** — 403 jika User-Agent terdeteksi sebagai scraper:

**Diblokir:** Empty UA, curl, wget, python-requests, aiohttp, scrapy, httpx, PostmanRuntime, insomnia, HttpClient, okhttp, Java/, ruby, faraday, generic bot/spider/crawler

**Diizinkan:** Googlebot, Bingbot, Yahoo Slurp, DuckDuckBot, Baiduspider, YandexBot, Sogou, Facebook, Twitter, LinkedIn, WhatsApp, Telegram, Discord

```json
// Response 403
{ "error": "Access denied: automated scraping is not allowed" }
```

---

## MinIO File Upload

Semua file (thumbnail, avatar, showcase media, study case image) diupload via endpoint terpusat `POST /api/upload`. Handler create/update hanya menerima URL CDN string.

**File structure di MinIO:**

```
{BUCKET}/{folder}/{uuid}.{ext}
→ https://cdn.mohagussetiaone.my.id/jvalleyverse/courses/a1b2c3d4-e5f6-7890-abcd-ef1234567890.jpg
```

**Folder convention:**
| Folder | Entity | Field |
|--------|--------|-------|
| courses | Course | thumbnail |
| lessons | Lesson | thumbnail |
| avatars | User | avatar |
| showcases | Showcase | media_urls[] |
| blogs | Blog | cover_img_url |
| study-cases | StudyCase | img_url |
