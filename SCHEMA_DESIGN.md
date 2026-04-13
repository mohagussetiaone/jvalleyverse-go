# JValleyVerse - Complete Database & API Schema Design

## 📋 Entity Relationship Diagram (ERD)

```
┌─────────────────────────────────────────────────────────────────────┐
│                         JVALLEYVERSE SYSTEM                         │
└─────────────────────────────────────────────────────────────────────┘

ENTITIES:
├── User (Core)
│   ├── id, email, password, name, avatar
│   ├── role (admin, user)
│   ├── points, level, totalPoints
│   └── timestamps
│
├── Category (Shared across features)
│   ├── id, name, slug, description
│   ├── icon, color
│   └── timestamps
│
├── PROJECT (Admin Management)
│   ├── id, title, description, thumbnail
│   ├── category_id (FK Category)
│   ├── admin_id (FK User - admin only)
│   ├── visibility (public, private)
│   └── timestamps
│
├── CLASS (Learning Path)
│   ├── id, title, description, thumbnail
│   ├── content (markdown)
│   ├── project_id (FK Project)
│   ├── category_id (FK Category)
│   ├── admin_id (FK User - created by admin)
│   ├── difficulty (beginner, intermediate, advanced)
│   └── timestamps
│
├── CERTIFICATE (User-specific, Private)
│   ├── id, unique_code, badge_url
│   ├── user_id (FK User - only owner can view)
│   ├── class_id (FK Class)
│   ├── issued_at, expires_at
│   └── timestamps
│
├── DISCUSSION (Class-related)
│   ├── id, title, content
│   ├── user_id (FK User - creator)
│   ├── class_id (FK Class)
│   ├── category_id (FK Category)
│   ├── views_count, status (open, closed)
│   └── timestamps
│
├── REPLY (Discussion Comments)
│   ├── id, content
│   ├── user_id (FK User)
│   ├── discussion_id (FK Discussion)
│   ├── parent_id (FK Reply - for nested replies)
│   ├── likes_count
│   └── timestamps
│
├── SHOWCASE (User Portfolio)
│   ├── id, title, description, media_urls[]
│   ├── user_id (FK User - creator/owner)
│   ├── category_id (FK Category)
│   ├── status (published, draft, archived)
│   ├── visibility (public, private, friends_only)
│   ├── likes_count, views_count
│   └── timestamps
│
├── SHOWCASE_LIKE (Composite)
│   ├── user_id (FK User - who liked)
│   ├── showcase_id (FK Showcase - what is liked)
│   └── created_at (timestamp)
│
├── SHOWCASE_COMMENT (Showcase Discussion)
│   ├── id, content
│   ├── user_id (FK User)
│   ├── showcase_id (FK Showcase)
│   ├── parent_id (FK ShowcaseComment - nested)
│   └── timestamps
│
├── COMMUNITY_POINT (Gamification - Activity Log)
│   ├── id, user_id (FK User)
│   ├── activity_type (certificate_issued, showcase_liked, discussion_reply, etc)
│   ├── points_earned, total_after
│   ├── metadata (related_object_id)
│   └── timestamps
│
└── USER_LEVEL (Points to Level Mapping)
    ├── level, min_points, max_points
    ├── badge_name, badge_icon
    └── description

RELATIONSHIPS:
User 1──∞ Class (admin_id)
User 1──∞ Project (admin_id)
User 1──∞ Certificate (owner - private)
User 1──∞ Discussion
User 1──∞ Reply
User 1──∞ Showcase
User 1──∞ ShowcaseLike
User 1──∞ CommunityPoint

Category 1──∞ Class
Category 1──∞ Project
Category 1──∞ Showcase
Category 1──∞ Discussion

Project 1──∞ Class

Class 1──∞ Certificate
Class 1──∞ Discussion

Discussion 1──∞ Reply
Reply 0...1──∞ Reply (self-referential)

Showcase 1──∞ ShowcaseLike
Showcase 1──∞ ShowcaseComment
ShowcaseComment 0...1──∞ ShowcaseComment (self-referential)

ShowcaseLike: Composite Primary Key (user_id, showcase_id)
```

## 🔐 Permission Matrix (Admin vs User vs Public)

| Feature                 | Admin  | User   | Public         | Notes                        |
| ----------------------- | ------ | ------ | -------------- | ---------------------------- |
| Create Project          | ✅     | ❌     | ❌             | Admin only feature           |
| Edit Project            | ✅ Own | ❌     | ❌             | Only admins can edit         |
| Delete Project          | ✅     | ❌     | ❌             | Admin only                   |
| Create Class            | ✅     | ❌     | ❌             | Under projects (admin)       |
| Edit Class              | ✅ Own | ❌     | ❌             | Creator admin only           |
| View Class              | ✅     | ✅     | ✅ (if public) | Public/private access        |
| Certificate Access      | ✅ Own | ✅ Own | ❌             | User-specific, JWT validated |
| Create Discussion       | ✅     | ✅     | ❌             | Logged-in only               |
| Reply Discussion        | ✅     | ✅     | ❌             | Logged-in only               |
| Create Showcase         | ✅     | ✅     | ❌             | Logged-in users              |
| Edit Showcase           | ✅ Own | ✅ Own | ❌             | Creator only                 |
| Delete Showcase         | ✅     | ✅ Own | ❌             | Creator/Admin                |
| Manage Showcase (Admin) | ✅     | ❌     | ❌             | Remove, archive, etc         |
| Like Showcase           | ✅     | ✅     | ❌             | Logged-in only               |
| View Points/Level       | ✅     | ✅ Own | ✅ Leaderboard | Public leaderboard           |

## 🎮 Gamification Rules

Points awarded for:

- Certificate Issued: **+50 pts**
- Showcase Posted: **+10 pts**
- Showcase Liked: **+5 pts** (receive per like)
- Discussion Created: **+5 pts**
- Reply to Discussion: **+2 pts**
- Reply Liked: **+1 pt**

Level Progression:

- Level 1: 0 - 199 pts (Beginner)
- Level 2: 200 - 499 pts (Learner)
- Level 3: 500 - 999 pts (Contributor)
- Level 4: 1000 - 1999 pts (Expert)
- Level 5: 2000+ pts (Master)

## 🔄 Data Flow Diagrams

### 1. Admin Project Creation Flow

```
Admin → Create Project + Category → DB (Project table)
                ↓
             Set Classes under Project → DB (Class table)
```

### 2. User Certificate Flow

```
User → Completes Class → System Validates → Create Certificate (PRIVATE)
                                                ↓
                                         Add +50 Points to User
                                                ↓
                                         Check Level Up
                                                ↓
                                         Log Activity (CommunityPoint)
```

### 3. User Showcase → Like Flow

```
User A → Create Showcase + Category → DB (Showcase - visibility: public)
                                               ↓
User B → View & Like → DB (ShowcaseLike - composite key)
                            ↓
                     Increase showcase.likes_count
                            ↓
                     User A receives +5 points (CommunityPoint)
                            ↓
                     Update User A level if needed
```

### 4. Discussion Reply Flow

```
User A → Create Discussion (in Class) → DB (Discussion)
                ↓
         Category tagged automatically from Class
                ↓
User B → Reply → DB (Reply with discussion_id)
           ↓
    User B gets +2 points
           ↓
User C → Reply to User B's Reply → DB (Reply with parent_id pointing to B's reply)
           ↓
    Nested reply display (limit depth to 2-3 levels)
```

## 📊 API Endpoints Overview

### Authentication

- POST /api/v1/auth/register
- POST /api/v1/auth/login
- POST /api/v1/auth/refresh
- POST /api/v1/auth/logout

### User Management

- GET /api/v1/users/me (profile)
- PUT /api/v1/users/me (update profile)
- GET /api/v1/users/:id (public profile)
- GET /api/v1/leaderboard (points ranking)

### Projects (Admin)

- POST /api/v1/admin/projects (create)
- GET /api/v1/projects (list - public)
- GET /api/v1/admin/projects (list - admin only)
- PUT /api/v1/admin/projects/:id (edit)
- DELETE /api/v1/admin/projects/:id (delete)

### Classes

- POST /api/v1/admin/classes (create under project)
- GET /api/v1/projects/:id/classes (list)
- GET /api/v1/classes/:id (view)
- PUT /api/v1/admin/classes/:id (edit)

### Certificates (Private)

- POST /api/v1/certificates/generate (issue when class completed)
- GET /api/v1/certificates (user's own certificates)
- GET /api/v1/certificates/:code (view own certificate - JWT validated)
- GET /api/v1/users/:id/certificates (public: show count only, private data hidden)

### Discussions

- POST /api/v1/discussions (create)
- GET /api/v1/classes/:id/discussions (list by class)
- GET /api/v1/discussions/:id (view single)
- PUT /api/v1/discussions/:id (edit)
- DELETE /api/v1/discussions/:id (delete)

### Replies

- POST /api/v1/discussions/:id/replies (create)
- GET /api/v1/discussions/:id/replies (list)
- POST /api/v1/replies/:id/replies (nested reply)
- PUT /api/v1/replies/:id (edit)
- DELETE /api/v1/replies/:id (delete)
- POST /api/v1/replies/:id/like (like reply)

### Showcase

- POST /api/v1/showcases (create)
- GET /api/v1/showcases (list public, paginated)
- GET /api/v1/showcases?category=tech (filter by category)
- GET /api/v1/showcases/:id (view)
- PUT /api/v1/showcases/:id (edit)
- DELETE /api/v1/showcases/:id (delete)
- POST /api/v1/showcases/:id/like (like)
- DELETE /api/v1/showcases/:id/like (unlike)
- GET /api/v1/showcases/:id/comments (list comments)
- POST /api/v1/showcases/:id/comments (create comment)

### Categories

- GET /api/v1/categories (list)
- POST /api/v1/admin/categories (create)
- PUT /api/v1/admin/categories/:id (edit)

### Gamification

- GET /api/v1/users/:id/points (public: leaderboard points)
- GET /api/v1/users/me/activity (user's own activity log)
- GET /api/v1/levels (level progression info)
- GET /api/v1/leaderboard (top users by points)
