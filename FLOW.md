# API Documentation

REST API untuk platform belajar berbasis komunitas. Dibangun dengan Go + GORM.

---

## Daftar Isi

- [Base URL & Auth](#base-url--auth)
- [Role & Permission](#role--permission)
- [Admin Flow](#admin-flow)
  - [1. Category](#1-category)
  - [2. Project](#2-project)
  - [3. Class](#3-class)
  - [4. Class Detail](#4-class-detail)
  - [5. Certificate](#5-certificate)
  - [6. User Management](#6-user-management)
  - [7. Level Config](#7-level-config)
- [User Flow](#user-flow)
  - [1. Auth](#1-auth)
  - [2. Browse Konten](#2-browse-konten)
  - [3. Progress Belajar](#3-progress-belajar)
  - [4. Discussion](#4-discussion)
  - [5. Showcase](#5-showcase)
  - [6. Profil & Gamifikasi](#6-profil--gamifikasi)
- [Gamifikasi — Poin & Level](#gamifikasi--poin--level)
- [Error Response](#error-response)

---

## Base URL & Auth

```
Base URL: /api/v1
```

Semua endpoint yang membutuhkan autentikasi menggunakan **JWT Bearer Token**.

```
Authorization: Bearer <token>
```

Token didapat dari endpoint `POST /auth/login` atau `POST /auth/register`.

---

## Role & Permission

| Role     | Keterangan                  |
| -------- | --------------------------- |
| `public` | Tidak perlu token           |
| `user`   | User terdaftar, token valid |
| `admin`  | Role admin, token valid     |

---

## Admin Flow

> Urutan pembuatan konten **wajib diikuti**: Category → Project → Class → Class Detail.  
> Tidak bisa membuat Project tanpa Category, tidak bisa membuat Class tanpa Project.

---

### 1. Category

#### `POST /admin/categories`

Buat category baru.

**Auth:** `admin`

**Request Body:**

```json
{
  "name": "Web Development",
  "slug": "web-development",
  "description": "Kelas seputar pengembangan web",
  "icon": "code",
  "color": "#3B82F6"
}
```

**Response `201`:**

```json
{
  "id": 1,
  "name": "Web Development",
  "slug": "web-development",
  "description": "Kelas seputar pengembangan web",
  "icon": "code",
  "color": "#3B82F6",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

> `slug` harus unik dan lowercase. Digunakan sebagai filter di seluruh konten (project, showcase, discussion).

---

#### `GET /admin/categories`

List semua categories.

**Auth:** `admin`

**Response `200`:**

```json
{
  "data": [
    {
      "id": 1,
      "name": "Web Development",
      "slug": "web-development",
      "color": "#3B82F6"
    }
  ],
  "total": 1
}
```

---

#### `PUT /admin/categories/:id`

Update category.

**Auth:** `admin`

**Request Body:** _(field yang ingin diubah)_

```json
{
  "name": "Web Dev",
  "color": "#6366F1"
}
```

**Response `200`:** Objek category yang sudah diupdate.

---

#### `DELETE /admin/categories/:id`

Hapus category (soft delete).

**Auth:** `admin`

**Response `200`:**

```json
{
  "message": "category deleted"
}
```

---

### 2. Project

#### `POST /admin/projects`

Buat project baru.

**Auth:** `admin`

**Request Body:**

```json
{
  "title": "Fullstack React & Go",
  "description": "Belajar fullstack dari nol hingga deployment",
  "thumbnail": "https://cdn.example.com/thumbnail.jpg",
  "category_id": 1,
  "visibility": "public"
}
```

> `visibility`: `public` | `draft`

**Response `201`:**

```json
{
  "id": 1,
  "title": "Fullstack React & Go",
  "category_id": 1,
  "admin_id": 2,
  "visibility": "public",
  "created_at": "2025-01-01T00:00:00Z"
}
```

> `admin_id` otomatis diambil dari JWT — tidak perlu dikirim di body.

---

#### `GET /admin/projects`

List semua project.

**Auth:** `admin`

**Query Params:**
| Param | Tipe | Keterangan |
|-------|------|-----------|
| `category_id` | int | Filter by category |
| `visibility` | string | `public` / `draft` |
| `page` | int | Default: 1 |
| `limit` | int | Default: 10 |

**Response `200`:**

```json
{
  "data": [
    {
      "id": 1,
      "title": "Fullstack React & Go",
      "category": { "id": 1, "name": "Web Development" },
      "visibility": "public"
    }
  ],
  "total": 1,
  "page": 1,
  "limit": 10
}
```

---

#### `PUT /admin/projects/:id`

Update project.

**Auth:** `admin`

**Request Body:** _(field yang ingin diubah)_

```json
{
  "title": "Fullstack React & Go — Updated",
  "visibility": "draft"
}
```

---

#### `DELETE /admin/projects/:id`

Hapus project beserta semua class di dalamnya (cascade).

**Auth:** `admin`

---

### 3. Class

#### `POST /admin/projects/:id/classes`

Buat class baru di dalam project.

**Auth:** `admin`

**Request Body:**

```json
{
  "title": "Setup Environment",
  "slug": "setup-environment",
  "description": "Instalasi tools dasar yang diperlukan",
  "thumbnail": "https://cdn.example.com/class.jpg",
  "difficulty": "beginner",
  "duration": 30,
  "order_index": 1,
  "is_first": true,
  "next_class_id": null,
  "visibility": "public"
}
```

> `difficulty`: `beginner` | `intermediate` | `advanced`  
> `duration`: dalam menit  
> `next_class_id`: isi `null` dulu, update setelah semua class dibuat  
> `is_first`: hanya satu class per project yang boleh `true`

**Response `201`:**

```json
{
  "id": 10,
  "title": "Setup Environment",
  "slug": "setup-environment",
  "project_id": 1,
  "order_index": 1,
  "is_first": true,
  "next_class_id": null,
  "difficulty": "beginner",
  "duration": 30,
  "visibility": "public"
}
```

---

#### `PUT /admin/classes/:id`

Update class. Gunakan ini untuk **chain `next_class_id`** setelah semua class dibuat.

**Auth:** `admin`

**Request Body:**

```json
{
  "next_class_id": 11,
  "order_index": 1
}
```

**Response `200`:** Objek class yang sudah diupdate.

> **Perhatian:** Validasi di server wajib mencegah circular reference (Class A → B → A).

---

#### `GET /admin/projects/:id/classes`

List semua class dalam project, diurutkan by `order_index`.

**Auth:** `admin`

---

#### `DELETE /admin/classes/:id`

Hapus class beserta detail dan sertifikat terkait (cascade).

**Auth:** `admin`

---

### 4. Class Detail

#### `POST /admin/classes/:id/detail`

Buat detail konten class. Satu class hanya punya satu ClassDetail.

**Auth:** `admin`

**Request Body:**

```json
{
  "about": "Di kelas ini kamu akan belajar cara setup environment Go dan React dari awal.",
  "rules": "Wajib menyelesaikan semua resource dan submit project akhir.",
  "tools": ["VSCode", "Go 1.22", "Node 20", "Docker"],
  "resource_media": {
    "videos": ["https://youtube.com/watch?v=xxx"],
    "documents": ["https://cdn.example.com/materi.pdf"],
    "images": ["https://cdn.example.com/diagram.png"]
  },
  "resources": [
    {
      "type": "pdf",
      "title": "Slide Materi Setup",
      "url": "https://cdn.example.com/slide.pdf"
    },
    {
      "type": "link",
      "title": "Repo GitHub Starter",
      "url": "https://github.com/example/starter"
    }
  ]
}
```

> `resource_media`: untuk embed konten langsung (video player, pdf viewer).  
> `resources`: daftar unduhan / referensi eksternal.  
> `type` pada resources: `pdf` | `video` | `link` | `document`

**Response `201`:**

```json
{
  "id": 5,
  "class_id": 10,
  "about": "Di kelas ini kamu akan...",
  "tools": ["VSCode", "Go 1.22", "Node 20", "Docker"],
  "created_at": "2025-01-01T00:00:00Z"
}
```

---

#### `PUT /admin/classes/:id/detail`

Update detail class yang sudah ada.

**Auth:** `admin`

**Request Body:** _(field yang ingin diubah)_

---

### 5. Certificate

#### `POST /admin/certificates`

Terbitkan sertifikat ke user setelah menyelesaikan class.

**Auth:** `admin`

**Request Body:**

```json
{
  "user_id": 42,
  "class_id": 10,
  "badge_url": "https://cdn.example.com/badge-setup-env.png",
  "expires_at": null
}
```

> `expires_at`: opsional. Isi `null` jika sertifikat tidak kedaluwarsa.  
> `unique_code` di-generate otomatis di server (UUID).

**Response `201`:**

```json
{
  "id": 7,
  "user_id": 42,
  "class_id": 10,
  "unique_code": "CERT-abc123xyz",
  "badge_url": "https://cdn.example.com/badge-setup-env.png",
  "issued_at": "2025-01-01T00:00:00Z",
  "expires_at": null
}
```

---

#### `GET /admin/users/:id/certificates`

Lihat semua sertifikat milik user tertentu.

**Auth:** `admin`

---

### 6. User Management

#### `GET /admin/users`

List semua user.

**Auth:** `admin`

**Query Params:**
| Param | Tipe | Keterangan |
|-------|------|-----------|
| `search` | string | Cari by name / email |
| `is_active` | bool | Filter aktif/nonaktif |
| `level` | int | Filter by level (1–5) |
| `page` | int | Default: 1 |
| `limit` | int | Default: 20 |

**Response `200`:**

```json
{
  "data": [
    {
      "id": 42,
      "name": "Budi Santoso",
      "email": "budi@mail.com",
      "role": "user",
      "level": 2,
      "points": 450,
      "total_points": 450,
      "is_active": true
    }
  ],
  "total": 1
}
```

---

#### `PUT /admin/users/:id`

Update role atau status aktif user.

**Auth:** `admin`

**Request Body:**

```json
{
  "is_active": false,
  "role": "user"
}
```

> `role`: `admin` | `user`  
> Set `is_active: false` untuk suspend user sementara.

**Response `200`:** Objek user yang sudah diupdate.

---

### 7. Level Config

#### `POST /admin/levels`

Konfigurasi threshold poin untuk tiap level.

**Auth:** `admin`

**Request Body:**

```json
{
  "level": 2,
  "min_points": 100,
  "max_points": 499,
  "badge_name": "Pelajar Aktif",
  "badge_icon": "star",
  "description": "Sudah menyelesaikan beberapa kelas"
}
```

> `level`: 1–5, harus unik.  
> `min_points` harus unik per level dan tidak overlap antar level.

**Response `201`:**

```json
{
  "id": 2,
  "level": 2,
  "min_points": 100,
  "max_points": 499,
  "badge_name": "Pelajar Aktif",
  "badge_icon": "star"
}
```

---

#### `GET /admin/levels`

List semua konfigurasi level.

**Auth:** `admin`

---

## User Flow

---

### 1. Auth

#### `POST /auth/register`

Daftar akun baru.

**Auth:** `public`

**Request Body:**

```json
{
  "name": "Budi Santoso",
  "email": "budi@mail.com",
  "password": "min8karakter"
}
```

**Response `201`:**

```json
{
  "token": "eyJhbGc...",
  "user": {
    "id": 42,
    "name": "Budi Santoso",
    "email": "budi@mail.com",
    "level": 1,
    "points": 0,
    "is_active": true
  }
}
```

---

#### `POST /auth/login`

Login dan dapatkan token.

**Auth:** `public`

**Request Body:**

```json
{
  "email": "budi@mail.com",
  "password": "min8karakter"
}
```

**Response `200`:**

```json
{
  "token": "eyJhbGc...",
  "user": {
    "id": 42,
    "name": "Budi Santoso",
    "level": 2,
    "points": 450
  }
}
```

---

### 2. Browse Konten

#### `GET /projects`

Daftar semua project yang `visibility: public`.

**Auth:** `user`

**Query Params:**
| Param | Tipe | Keterangan |
|-------|------|-----------|
| `category_id` | int | Filter by category |
| `search` | string | Cari by title |
| `page` | int | Default: 1 |
| `limit` | int | Default: 10 |

**Response `200`:**

```json
{
  "data": [
    {
      "id": 1,
      "title": "Fullstack React & Go",
      "thumbnail": "https://cdn.example.com/thumbnail.jpg",
      "category": { "id": 1, "name": "Web Development" }
    }
  ],
  "total": 1
}
```

---

#### `GET /projects/:id/classes`

Daftar class dalam project, diurutkan by `order_index`.

**Auth:** `user`

**Response `200`:**

```json
{
  "data": [
    {
      "id": 10,
      "title": "Setup Environment",
      "slug": "setup-environment",
      "order_index": 1,
      "is_first": true,
      "difficulty": "beginner",
      "duration": 30
    },
    {
      "id": 11,
      "title": "Build REST API",
      "slug": "build-rest-api",
      "order_index": 2,
      "is_first": false,
      "difficulty": "intermediate",
      "duration": 60
    }
  ]
}
```

---

#### `GET /classes/:id`

Detail class beserta konten lengkap dan `next_class_id`.

**Auth:** `user`

**Response `200`:**

```json
{
  "id": 10,
  "title": "Setup Environment",
  "slug": "setup-environment",
  "difficulty": "beginner",
  "duration": 30,
  "next_class_id": 11,
  "details": {
    "about": "Di kelas ini kamu akan...",
    "rules": "Wajib submit project akhir.",
    "tools": ["VSCode", "Go 1.22"],
    "resource_media": {
      "videos": ["https://youtube.com/watch?v=xxx"],
      "documents": [],
      "images": []
    },
    "resources": [{ "type": "pdf", "title": "Slide Materi", "url": "https://cdn.example.com/slide.pdf" }]
  },
  "discussions": []
}
```

---

### 3. Progress Belajar

#### `POST /classes/:id/progress`

Mulai atau update progress belajar.

**Auth:** `user`

**Request Body:**

```json
{
  "status": "in_progress",
  "progress_percentage": 60,
  "notes": "Sudah sampai bagian instalasi Go"
}
```

> `status`: `not_started` | `started` | `in_progress` | `completed`  
> Saat `status: "completed"`, server otomatis:
>
> 1. Insert `CommunityPoint` dengan activity_type `class_completed`
> 2. Update `User.points` dan `User.total_points`
> 3. Cek threshold `UserLevel` — update `User.level` jika naik

**Response `200`:**

```json
{
  "id": 3,
  "user_id": 42,
  "class_id": 10,
  "status": "in_progress",
  "progress_percentage": 60,
  "started_at": "2025-01-01T08:00:00Z",
  "completed_at": null
}
```

---

#### `GET /me/progress`

Semua progress belajar user yang sedang login.

**Auth:** `user`

**Response `200`:**

```json
{
  "data": [
    {
      "class_id": 10,
      "class_title": "Setup Environment",
      "status": "completed",
      "progress_percentage": 100,
      "completed_at": "2025-01-01T10:00:00Z"
    },
    {
      "class_id": 11,
      "class_title": "Build REST API",
      "status": "in_progress",
      "progress_percentage": 40
    }
  ]
}
```

---

### 4. Discussion

#### `POST /discussions`

Buat thread diskusi baru.

**Auth:** `user`

**Request Body:**

```json
{
  "title": "Kenapa goroutine lebih ringan dari thread OS?",
  "content": "Saya penasaran soal perbedaan goroutine dan thread biasa...",
  "class_id": 10,
  "category_id": 1
}
```

> `class_id`: opsional — bisa membuat diskusi standalone tanpa class.

**Response `201`:**

```json
{
  "id": 5,
  "title": "Kenapa goroutine lebih ringan dari thread OS?",
  "user_id": 42,
  "class_id": 10,
  "category_id": 1,
  "status": "open",
  "is_pinned": false,
  "views_count": 0,
  "created_at": "2025-01-01T00:00:00Z"
}
```

---

#### `GET /discussions`

List semua diskusi.

**Auth:** `user`

**Query Params:**
| Param | Tipe | Keterangan |
|-------|------|-----------|
| `class_id` | int | Filter by class |
| `category_id` | int | Filter by category |
| `status` | string | `open` / `closed` |
| `page` | int | Default: 1 |
| `limit` | int | Default: 20 |

---

#### `GET /discussions/:id`

Detail diskusi beserta replies.

**Auth:** `user`

> Views count otomatis bertambah setiap kali endpoint ini dipanggil.

---

#### `POST /discussions/:id/replies`

Balas diskusi.

**Auth:** `user`

**Request Body:**

```json
{
  "content": "Goroutine dikelola oleh Go runtime, bukan OS, sehingga lebih ringan...",
  "parent_id": null
}
```

> `parent_id`: isi ID reply lain jika ingin nested reply (balas reply).

**Response `201`:**

```json
{
  "id": 20,
  "discussion_id": 5,
  "user_id": 42,
  "content": "Goroutine dikelola oleh Go runtime...",
  "parent_id": null,
  "likes_count": 0,
  "is_marked_best": false
}
```

---

#### `POST /replies/:id/mark-best`

Tandai reply sebagai jawaban terbaik.

**Auth:** `user` _(hanya owner discussion)_

**Response `200`:**

```json
{
  "id": 20,
  "is_marked_best": true
}
```

---

### 5. Showcase

#### `POST /showcases`

Upload karya portofolio.

**Auth:** `user`

**Request Body:**

```json
{
  "title": "Todo App with Go + React",
  "description": "Project akhir kelas Fullstack. Fitur: auth, CRUD, realtime update.",
  "media_urls": ["https://cdn.example.com/screenshot1.png", "https://cdn.example.com/screenshot2.png"],
  "category_id": 1,
  "visibility": "public"
}
```

> `media_urls`: array URL gambar/video yang sudah diupload ke storage terlebih dahulu.  
> `visibility`: `public` | `private`

**Response `201`:**

```json
{
  "id": 8,
  "title": "Todo App with Go + React",
  "user_id": 42,
  "category_id": 1,
  "status": "published",
  "visibility": "public",
  "likes_count": 0,
  "views_count": 0
}
```

---

#### `GET /showcases`

List semua showcase publik.

**Auth:** `user`

**Query Params:**
| Param | Tipe | Keterangan |
|-------|------|-----------|
| `category_id` | int | Filter by category |
| `user_id` | int | Filter by user |
| `page` | int | Default: 1 |
| `limit` | int | Default: 12 |

---

#### `POST /showcases/:id/like`

Toggle like pada showcase.

**Auth:** `user`

**Request Body:** _(tidak diperlukan)_

**Response `200`:**

```json
{
  "showcase_id": 8,
  "user_id": 42,
  "liked": true
}
```

> Menggunakan composite key `(user_id, showcase_id)` — otomatis toggle antara like dan unlike.

---

#### `POST /showcases/:id/comments`

Komentar pada showcase.

**Auth:** `user`

**Request Body:**

```json
{
  "content": "Keren banget projectnya! Stack-nya solid.",
  "parent_id": null
}
```

> `parent_id`: isi ID komentar lain untuk nested comment.

**Response `201`:**

```json
{
  "id": 15,
  "showcase_id": 8,
  "user_id": 42,
  "content": "Keren banget projectnya!",
  "parent_id": null,
  "created_at": "2025-01-01T00:00:00Z"
}
```

---

### 6. Profil & Gamifikasi

#### `GET /me`

Profil lengkap user yang sedang login.

**Auth:** `user`

**Response `200`:**

```json
{
  "id": 42,
  "name": "Budi Santoso",
  "email": "budi@mail.com",
  "avatar": "https://cdn.example.com/avatar.jpg",
  "bio": "Learning everyday",
  "role": "user",
  "level": 2,
  "points": 450,
  "total_points": 450,
  "is_active": true
}
```

---

#### `PUT /me`

Update profil sendiri.

**Auth:** `user`

**Request Body:**

```json
{
  "name": "Budi S.",
  "avatar": "https://cdn.example.com/new-avatar.jpg",
  "bio": "Fullstack developer in progress"
}
```

> User tidak bisa mengubah `email`, `role`, atau `level` sendiri.

---

#### `GET /me/points`

Riwayat transaksi poin user.

**Auth:** `user`

**Response `200`:**

```json
{
  "data": [
    {
      "id": 1,
      "activity_type": "class_completed",
      "points_earned": 100,
      "points_after": 450,
      "level_after": 2,
      "description": "Selesai kelas: Setup Environment",
      "metadata": {
        "class_id": 10,
        "class_title": "Setup Environment"
      },
      "created_at": "2025-01-01T10:00:00Z"
    }
  ]
}
```

---

#### `GET /me/certificates`

Sertifikat milik user yang sedang login.

**Auth:** `user`

**Response `200`:**

```json
{
  "data": [
    {
      "id": 7,
      "unique_code": "CERT-abc123xyz",
      "badge_url": "https://cdn.example.com/badge.png",
      "issued_at": "2025-01-01T00:00:00Z",
      "expires_at": null,
      "class": {
        "id": 10,
        "title": "Setup Environment",
        "project": { "id": 1, "title": "Fullstack React & Go" }
      }
    }
  ]
}
```

> **Private** — hanya mengembalikan sertifikat milik user yang login. Admin menggunakan `GET /admin/users/:id/certificates`.

---

#### `GET /levels`

List semua konfigurasi level (untuk ditampilkan di UI).

**Auth:** `user`

**Response `200`:**

```json
{
  "data": [
    { "level": 1, "min_points": 0, "max_points": 99, "badge_name": "Pemula", "badge_icon": "seedling" },
    { "level": 2, "min_points": 100, "max_points": 499, "badge_name": "Pelajar Aktif", "badge_icon": "star" },
    { "level": 3, "min_points": 500, "max_points": 1499, "badge_name": "Kontributor", "badge_icon": "fire" },
    { "level": 4, "min_points": 1500, "max_points": 4999, "badge_name": "Expert", "badge_icon": "trophy" },
    { "level": 5, "min_points": 5000, "max_points": 999999, "badge_name": "Master", "badge_icon": "crown" }
  ]
}
```

---

## Gamifikasi — Poin & Level

Sistem poin ditrigger otomatis di server. Client **tidak perlu** mengirim request poin secara manual.

| Activity Type        | Poin | Keterangan                 |
| -------------------- | ---- | -------------------------- |
| `class_completed`    | 100  | Menyelesaikan satu class   |
| `showcase_posted`    | 50   | Upload showcase baru       |
| `discussion_created` | 20   | Membuat thread diskusi     |
| `reply_posted`       | 10   | Membalas diskusi           |
| `reply_marked_best`  | 30   | Reply ditandai best answer |
| `showcase_liked`     | 5    | Mendapat like di showcase  |

**Level naik** diperiksa setiap kali poin bertambah. Jika `User.points >= UserLevel.min_points` untuk level berikutnya, `User.level` diupdate dan `CommunityPoint.level_after` mencatat level baru.

---

## Error Response

Semua error menggunakan format berikut:

```json
{
  "error": "pesan error yang human-readable",
  "code": "ERROR_CODE"
}
```

| HTTP Status | Keterangan                                             |
| ----------- | ------------------------------------------------------ |
| `400`       | Request tidak valid / field wajib kosong               |
| `401`       | Token tidak ada atau expired                           |
| `403`       | Tidak punya akses (role tidak sesuai)                  |
| `404`       | Resource tidak ditemukan                               |
| `409`       | Conflict — misal slug sudah dipakai, sudah pernah like |
| `500`       | Internal server error                                  |

**Contoh `401`:**

```json
{
  "error": "token tidak valid atau sudah expired",
  "code": "UNAUTHORIZED"
}
```

**Contoh `403`:**

```json
{
  "error": "hanya admin yang bisa mengakses endpoint ini",
  "code": "FORBIDDEN"
}
```
