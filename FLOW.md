# Flow Final

Dokumen ini merangkum flow final backend yang sekarang aktif.

## Struktur Utama

```text
Category
└─ Project
   └─ Phase
      └─ Class
         ├─ Class Detail
         ├─ Class Progress
         └─ Certificate
```

## Peran

- `public`: bisa lihat konten publik
- `user`: bisa belajar, diskusi, showcase, lihat certificate sendiri
- `admin`: mengelola konten dan user

## 1. Auth

Flow:

1. User register
2. User login
3. Server memberi JWT
4. JWT dipakai untuk endpoint protected

Endpoint:

- `POST /api/auth/register`
- `POST /api/auth/login`

## 2. Admin Content Flow

Urutan wajib:

1. Buat `Category`
2. Buat `Project` di dalam category
3. Buat `Phase` di dalam project
4. Buat `Class` di dalam phase
5. Tambah `Class Detail`

Aturan:

- `class` wajib punya `project_id` dan `phase_id`
- admin hanya bisa ubah data miliknya sendiri
- class lama yang belum punya `phase` akan di-backfill ke default phase saat migrasi

Endpoint admin:

- `POST /api/admin/categories`
- `GET /api/admin/categories`
- `PUT /api/admin/categories/:id`
- `DELETE /api/admin/categories/:id`
- `POST /api/admin/projects`
- `GET /api/admin/projects`
- `PUT /api/admin/projects/:id`
- `DELETE /api/admin/projects/:id`
- `POST /api/admin/projects/:project_id/phases`
- `PUT /api/admin/phases/:phase_id`
- `DELETE /api/admin/phases/:phase_id`
- `POST /api/admin/classes`
- `POST /api/admin/classes/:id/details`
- `PUT /api/admin/classes/:id`
- `DELETE /api/admin/classes/:id`

## 3. Public Browse Flow

Flow:

1. Lihat daftar category
2. Lihat project di category
3. Buka project
4. Lihat phases
5. Lihat classes
6. Buka detail class

Endpoint public:

- `GET /api/categories`
- `GET /api/categories/:slug`
- `GET /api/categories/:category_id/projects`
- `GET /api/projects/:project_id`
- `GET /api/projects/:project_id/phases`
- `GET /api/projects/:project_id/phases/:phase_id`
- `GET /api/projects/:project_id/classes`
- `GET /api/projects/:project_id/phases/:phase_id/classes`
- `GET /api/projects/:project_id/classes/:slug`

## 4. Learning Flow

Flow user belajar:

1. User buka detail class
2. User `start class`
3. User `update progress`
4. User `complete class`
5. Server update progress
6. Server tambah poin
7. Server buat certificate

Status progress:

```text
not_started -> started -> in_progress -> completed
```

Endpoint:

- `POST /api/classes/:id/start`
- `PUT /api/classes/:id/progress`
- `POST /api/classes/:id/complete`

Hasil saat complete:

- progress jadi `completed`
- points bertambah
- response mengandung `certificate`
- response mengandung `achievement`
- response mengandung `next_class` bila ada

## 5. Certificate Flow

Flow:

1. Certificate dibuat saat class selesai
2. Owner bisa lihat daftar certificate miliknya
3. Certificate by code hanya bisa diakses owner atau admin

Endpoint:

- `GET /api/certificates`
- `GET /api/certificates/:code`

## 6. User Profile Flow

Flow:

1. User lihat profil sendiri
2. User update profil sendiri
3. User lihat activity log
4. Public bisa lihat profil ringkas user

Endpoint:

- `GET /api/users/me`
- `PUT /api/users/me`
- `GET /api/users/me/activity`
- `GET /api/users/:id`

## 7. Discussion Flow

Flow:

1. User buat discussion
2. Discussion bisa dikaitkan ke class
3. User lain balas discussion
4. Owner reply bisa edit atau hapus reply miliknya
5. Owner discussion bisa update thread miliknya

Endpoint:

- `POST /api/discussions`
- `GET /api/discussions`
- `GET /api/discussions/:id`
- `PUT /api/discussions/:id`
- `POST /api/discussions/:id/replies`
- `PUT /api/replies/:id`
- `DELETE /api/replies/:id`

## 8. Showcase Flow

Flow:

1. User buat showcase
2. Public bisa lihat showcase
3. Owner bisa update atau delete showcase
4. User lain bisa like atau unlike
5. Showcase owner bisa dapat poin dari interaksi

Endpoint:

- `GET /api/showcases`
- `GET /api/showcases/:id`
- `POST /api/showcases`
- `PUT /api/showcases/:id`
- `DELETE /api/showcases/:id`
- `POST /api/showcases/:id/like`
- `DELETE /api/showcases/:id/like`

## 9. Gamification Flow

Flow:

1. Aktivitas user menghasilkan poin
2. Poin tersimpan ke activity log
3. Level user dihitung dari total poin
4. Leaderboard menampilkan ranking user

Endpoint:

- `GET /api/leaderboard`
- `GET /api/levels`
- `GET /api/users/:id/points`

Contoh aktivitas yang memicu poin:

- menyelesaikan class
- membuat showcase
- mendapat like showcase
- membuat reply atau diskusi

## 10. Admin User Flow

Flow:

1. Admin buka dashboard
2. Admin lihat daftar user
3. Admin memantau konten dan struktur belajar

Endpoint:

- `GET /api/admin/dashboard`
- `GET /api/admin/users`

## Ringkasan Besar

```text
Admin:
Category -> Project -> Phase -> Class -> Class Detail

User:
Browse -> Open Class -> Start -> Progress -> Complete -> Certificate

Community:
Discussion + Reply + Showcase + Like

Gamification:
Activity -> Points -> Level -> Leaderboard
```
