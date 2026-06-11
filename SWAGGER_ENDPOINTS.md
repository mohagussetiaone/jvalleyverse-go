# JValleyVerse API - Daftar Lengkap Endpoint

## 📋 Ringkas

Dokumentasi Swagger telah dilengkapi dengan **semua endpoint** dari aplikasi.

---

## 🔗 Relasi Project ↔ Class

```
Project (1) ──────────── (many) Class
   │                              │
   │ has many Classes             │ belongs to Project (via project_id)
   │                              │
   └── GET /api/projects/:project_id/classes          → List semua class dalam project
   └── GET /api/projects/:project_id/classes/:slug    → Detail class berdasarkan slug
```

**Aturan penting:**
- Setiap `Class` **wajib** memiliki `project_id` — class tidak bisa berdiri sendiri tanpa project
- Untuk mengakses detail class, selalu butuh `project_id` + `slug` class (bukan ID langsung)
- Admin membuat `Project` dulu, kemudian membuat `Class` dengan menyertakan `project_id`
- `ClassDetail` (konten lengkap: about, rules, tools, resources) merupakan entitas terpisah yang di-attach ke class via `POST /api/admin/classes/:id/details`

---

## 🔓 PUBLIC ENDPOINTS (Tanpa Autentikasi)

### Authentication

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| `POST` | `/api/auth/register` | Register user baru |
| `POST` | `/api/auth/login` | Login dan dapatkan JWT token |

### Showcases (Public)

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| `GET` | `/api/leaderboard` | Lihat leaderboard top users |
| `GET` | `/api/showcases` | List semua showcases (dengan filter) |
| `GET` | `/api/showcases/:id` | Detail showcase |

### Projects & Classes (Public)

> Class diakses melalui konteks project-nya — selalu ada `project_id` dalam path.

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| `GET` | `/api/projects/:project_id/classes` | List semua class dalam project tertentu |
| `GET` | `/api/projects/:project_id/classes/:slug` | Detail class berdasarkan slug |

### Categories (Public)

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| `GET` | `/api/categories` | List semua kategori |
| `GET` | `/api/categories/:slug` | Detail kategori berdasarkan slug |
| `GET` | `/api/categories/:category_id/projects` | List projects dalam kategori tertentu |

### Health

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| `GET` | `/api/health` | Server health check |
| `GET` | `/api/health/detailed` | Health check detail (DB, Redis, dll) |

---

## 🔒 PROTECTED ENDPOINTS (Memerlukan JWT Token)

### User Management

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| `GET` | `/api/users/me` | Profil user saat ini |
| `PUT` | `/api/users/me` | Update profil user |
| `GET` | `/api/users/:id` | Lihat profil user lain (public) |
| `GET` | `/api/users/me/activity` | Activity log user |

### Showcases (User)

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| `POST` | `/api/showcases` | Buat showcase baru |
| `PUT` | `/api/showcases/:id` | Update showcase (owner only) |
| `DELETE` | `/api/showcases/:id` | Hapus showcase (owner only) |
| `POST` | `/api/showcases/:id/like` | Like showcase |
| `DELETE` | `/api/showcases/:id/like` | Unlike showcase |

### Certificates

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| `GET` | `/api/certificates` | List sertifikat milik user sendiri |
| `GET` | `/api/certificates/:code` | View sertifikat (owner only) |

### Discussions

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| `POST` | `/api/discussions` | Buat diskusi baru (opsional: attach ke class via `class_id`) |
| `GET` | `/api/discussions` | List diskusi (dengan filter) |
| `GET` | `/api/discussions/:id` | Detail diskusi dengan replies |
| `PUT` | `/api/discussions/:id` | Update diskusi (owner only) |

### Replies

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| `POST` | `/api/discussions/:id/replies` | Balas diskusi |
| `PUT` | `/api/replies/:id` | Update reply (owner only) |
| `DELETE` | `/api/replies/:id` | Hapus reply (owner only) |

> Nested reply (balas reply) dilakukan via `POST /api/discussions/:id/replies` dengan menyertakan `parent_id` di body.

### Classes (User — Progress Tracking)

> Aksi progress class menggunakan `class_id` langsung (tidak perlu `project_id`), karena user sudah tahu class-nya.

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| `POST` | `/api/classes/:id/start` | Mulai / tandai class sebagai started |
| `PUT` | `/api/classes/:id/progress` | Update persentase progress class |
| `POST` | `/api/classes/:id/complete` | Mark class sebagai completed |

### Gamification

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| `GET` | `/api/levels` | Info level progression (level 1–5 + syarat poin) |
| `GET` | `/api/users/:id/points` | Lihat points dan level user |

---

## 👨‍💼 ADMIN ENDPOINTS (Memerlukan JWT + Role `admin`)

### Dashboard

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| `GET` | `/api/admin/dashboard` | Admin dashboard |

### User Management (Admin)

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| `GET` | `/api/admin/users` | List semua users dengan pagination |

### Projects (Admin)

> Admin harus membuat Project terlebih dahulu sebelum bisa membuat Class.

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| `POST` | `/api/admin/projects` | Buat project baru |
| `GET` | `/api/admin/projects` | List semua projects |
| `PUT` | `/api/admin/projects/:id` | Update project |
| `DELETE` | `/api/admin/projects/:id` | Hapus project (cascade hapus semua class-nya) |

### Classes (Admin)

> Class dibuat dengan menyertakan `project_id` untuk mengikat class ke project tertentu.

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| `POST` | `/api/admin/classes` | Buat class baru (wajib sertakan `project_id`) |
| `POST` | `/api/admin/classes/:id/details` | Tambah/set konten detail class (about, rules, tools, resources) |
| `PUT` | `/api/admin/classes/:id` | Update class |
| `DELETE` | `/api/admin/classes/:id` | Hapus class |

### Categories (Admin)

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| `POST` | `/api/admin/categories` | Buat kategori baru |
| `GET` | `/api/admin/categories` | List semua kategori (admin view) |
| `PUT` | `/api/admin/categories/:id` | Update kategori |
| `DELETE` | `/api/admin/categories/:id` | Hapus kategori |

---

## 📊 Alur Kerja Admin: Membuat Konten Pembelajaran

```
1. POST /api/admin/categories          → Buat kategori (misal: "Web Development")
2. POST /api/admin/projects            → Buat project, sertakan category_id
3. POST /api/admin/classes             → Buat class, sertakan project_id
4. POST /api/admin/classes/:id/details → Set konten detail class
```

## 📊 Alur Kerja User: Mengikuti Pembelajaran

```
1. GET  /api/categories/:slug/projects           → Temukan project berdasarkan kategori
2. GET  /api/projects/:project_id/classes        → Lihat daftar class dalam project
3. GET  /api/projects/:project_id/classes/:slug  → Baca detail class
4. POST /api/classes/:id/start                   → Mulai class
5. PUT  /api/classes/:id/progress                → Update progress
6. POST /api/classes/:id/complete                → Selesaikan class → otomatis dapat poin + sertifikat
```

---

## 📊 Struktur Response Umum

### Authentication Response

```json
{
  "token": "jwt_token_here",
  "xsrf_token": "xsrf_token_here"
}
```

### User Response

```json
{
  "id": "cuid_string",
  "email": "user@example.com",
  "name": "John Doe",
  "avatar": "https://...",
  "bio": "...",
  "role": "user",
  "points": 100,
  "total_points": 250,
  "level": 2,
  "created_at": "2024-01-01T00:00:00Z"
}
```

### Project Response

```json
{
  "id": "cuid_string",
  "title": "Belajar Go",
  "description": "...",
  "thumbnail": "https://...",
  "category_id": "cuid_string",
  "category": { "id": "...", "name": "Backend", "slug": "backend" },
  "admin_id": "cuid_string",
  "visibility": "public",
  "classes": []
}
```

### Class Response

```json
{
  "id": "cuid_string",
  "title": "Pengenalan Go",
  "slug": "pengenalan-go",
  "description": "...",
  "project_id": "cuid_string",
  "project": { "id": "...", "title": "Belajar Go" },
  "difficulty": "beginner",
  "duration": 60,
  "order_index": 1,
  "next_class_id": "cuid_string_or_null",
  "details": { "about": "...", "tools": [], "resources": [] }
}
```

### Pagination Response

```json
{
  "data": [],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100
  }
}
```

---

## 🔐 Security

Semua endpoint protected menggunakan JWT Bearer Token:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

---

## 🚀 Akses Dokumentasi Swagger

File dokumentasi: `openapi.json`

UI Swagger dapat diakses dari endpoint Swagger Docs:

- Local: `http://localhost:3000/api/docs`
- Production: `https://jvalleyverse.mohagussetiaone.my.id/api/docs`

---

## ✅ Status Dokumentasi

- ✓ Semua PUBLIC endpoints didokumentasikan
- ✓ Semua PROTECTED endpoints didokumentasikan
- ✓ Semua ADMIN endpoints didokumentasikan
- ✓ Relasi Project ↔ Class dijelaskan secara eksplisit
- ✓ Alur kerja admin dan user tersedia
- ✓ Request/Response schemas lengkap
- ✓ Error responses terdefinisi
- ✓ Security schemes terdeklarasi
