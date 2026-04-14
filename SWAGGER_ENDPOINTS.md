# JValleyVerse API - Daftar Lengkap Endpoint

## 📋 Ringkas

Dokumentasi Swagger telah dilengkapi dengan **semua endpoint** dari aplikasi. Total: **40+ endpoints**

---

## 🔓 PUBLIC ENDPOINTS (Tanpa Autentikasi)

### Authentication

- `POST /auth/register` - Register user baru
- `POST /auth/login` - Login dan dapatkan JWT token

### Showcases (Public)

- `GET /leaderboard` - Lihat leaderboard top users
- `GET /showcases` - List semua showcases (dengan filter)
- `GET /showcases/{id}` - Detail showcase

### Classes (Public)

- `GET /classes/{id}` - Detail class
- `GET /projects/{id}/classes` - List classes dalam project

### Health

- `GET /health` - Server health check

---

## 🔒 PROTECTED ENDPOINTS (Memerlukan JWT Token)

### User Management

- `GET /users/me` - Profil user saat ini
- `PUT /users/me` - Update profil user
- `GET /users/{id}` - Lihat profil user lain (public)
- `GET /users/me/activity` - Activity log user

### Showcases (User)

- `POST /showcases` - Buat showcase baru
- `PUT /showcases/{id}` - Update showcase (owner only)
- `DELETE /showcases/{id}` - Hapus showcase (owner only)
- `POST /showcases/{id}/like` - Like showcase
- `DELETE /showcases/{id}/like` - Unlike showcase

### Certificates

- `GET /certificates` - List sertifikat user
- `GET /certificates/{code}` - View sertifikat (owner only)

### Discussions

- `POST /discussions` - Buat diskusi baru
- `GET /discussions` - List diskusi (dengan filter)
- `GET /discussions/{id}` - Detail diskusi dengan replies
- `PUT /discussions/{id}` - Update diskusi (owner only)
- `DELETE /discussions/{id}` - Hapus diskusi (owner only)

### Replies

- `POST /discussions/{id}/replies` - Balas diskusi
- `PUT /replies/{id}` - Update reply (owner only)
- `DELETE /replies/{id}` - Hapus reply (owner only)
- `POST /replies/{id}/replies` - Nested reply (balas reply)

### Classes (User)

- `POST /classes/{id}/complete` - Mark class sebagai completed

### Gamification

- `GET /levels` - Info level progression
- `GET /users/{id}/points` - Lihat points dan level user

---

## 👨‍💼 ADMIN ENDPOINTS (Memerlukan JWT + Role Admin)

### Dashboard

- `GET /admin/dashboard` - Admin dashboard

### Projects (Admin)

- `POST /admin/projects` - Buat project baru
- `GET /admin/projects` - List semua projects
- `PUT /admin/projects/{id}` - Update project
- `DELETE /admin/projects/{id}` - Hapus project

### Classes (Admin)

- `POST /admin/classes` - Buat class baru
- `PUT /admin/classes/{id}` - Update class
- `DELETE /admin/classes/{id}` - Hapus class

---

## 📊 Struktur Response

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
  "id": 1,
  "email": "user@example.com",
  "name": "John Doe",
  "avatar": "https://...",
  "bio": "...",
  "role": "user",
  "points": 100,
  "level": 1,
  "created_at": "2024-01-01T00:00:00Z"
}
```

### Pagination Response

```json
{
  "data": [...],
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
- ✓ Request/Response schemas lengkap
- ✓ Error responses terdefinisi
- ✓ Security schemes terdeklarasi
