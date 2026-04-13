# 🔄 STRUCTURE UPDATE - NEW Category/Project/Class Hierarchy

## Summary of Changes

User requested updated hierarchy structure to better organize learning content:

### **OLD Structure** ❌

```
Project
  └─ Class
      └─ Certificate
```

### **NEW Structure** ✅

```
Category (HEAD - Top Level)
  └─ Project
      └─ Class (with slug, next_class_id)
          ├─ ClassDetail (about, rules, tools, resource_media)
          ├─ ClassProgress (not_started → started → in_progress → completed)
          └─ Certificate
```

---

## 📋 What Changed

### 1. **Category Becomes the "Head"**

- Category is now top-level grouping (was just a tag)
- Projects belong to Categories
- Classes inherit category from Project (for filtering)

### 2. **Class Structure Enhanced**

- ✅ Added `slug` field for URL-friendly routing
- ✅ Added `next_class_id` for linear progression through classes
- ✅ Added `sequence_number` to order classes in project
- ✅ Added `is_first` boolean to mark first class
- ✅ Moved content details to separate `ClassDetail` table

### 3. **New Tables Created**

- ✅ `class_details` - Store: about, rules, tools, resource_media, resources
- ✅ `class_progress` - Track: user progress status, percentage, start/complete dates

### 4. **Progress Tracking Added**

- States: `not_started → started → in_progress → completed`
- Tracks: `started_at`, `completed_at`, `progress_percentage`
- Enables: User can see where they are in each class

---

## 📁 NEW Entity Relationships

```
User (admin)
  ├─→ creates many → Category
  ├─→ creates many → Project
  ├─→ creates many → Class
  └─→ progress in → ClassProgress

Category
  ├─→ contains many → Project
  └─→ categorizes → Class

Project
  └─→ contains many → Class

Class
  ├─→ has one → ClassDetail (resources, rules, tools)
  ├─→ has many → ClassProgress (multiple users can progress)
  ├─→ linked to → NextClass (next_class_id)
  └─→ issued → Certificate (on completion)

ClassProgress
  ├─→ belongs to → User
  ├─→ belongs to → Class
  └─→ tracks: status, percentage, dates

Certificate
  ├─→ belongs to → User
  ├─→ belongs to → Class
  └─→ unique: one per user per class
```

---

## 🗄️ Database Changes Required

### New Tables

#### `class_details`

```sql
CREATE TABLE class_details (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    class_id BIGINT NOT NULL UNIQUE,
    about TEXT,
    rules LONGTEXT,
    tools JSON,  -- ["tool1", "tool2"]
    resource_media JSON,  -- {videos, documents, images}
    resources JSON,  -- [{type, title, url}]
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE CASCADE
);
```

#### `class_progress`

```sql
CREATE TABLE class_progress (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    class_id BIGINT NOT NULL,
    status VARCHAR(50) DEFAULT 'not_started',  -- not_started|started|in_progress|completed
    started_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    progress_percentage INT DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    UNIQUE KEY unique_user_class (user_id, class_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE CASCADE
);
```

### Updated Tables

#### `classes` - Add new columns

```sql
ALTER TABLE classes ADD COLUMN slug VARCHAR(255) NOT NULL;
ALTER TABLE classes ADD COLUMN next_class_id BIGINT NULL;
ALTER TABLE classes ADD COLUMN sequence_number INT;
ALTER TABLE classes ADD COLUMN is_first BOOLEAN DEFAULT FALSE;

-- Add unique constraint on slug per project
ALTER TABLE classes ADD UNIQUE KEY unique_slug_per_project (project_id, slug);

-- Add foreign key for next_class_id
ALTER TABLE classes ADD FOREIGN KEY (next_class_id) REFERENCES classes(id);
```

#### `categories` - Already exists (good!)

- Already has: id, name, slug, description, icon, color
- No changes needed

---

## 🔀 API Changes

### New Endpoints

```
Category Management
POST   /api/admin/categories
GET    /api/categories
GET    /api/categories/{slug}
PUT    /api/admin/categories/{id}
DELETE /api/admin/categories/{id}

Updated Class Endpoints (use slug)
GET    /api/projects/{id}/classes/{slug}      ← NEW: use slug instead of ID
POST   /api/classes/{id}/start                ← NEW: start class
PUT    /api/classes/{id}/progress             ← NEW: update progress
POST   /api/classes/{id}/complete             ← UPDATED: now creates certificate & awards points

Progress Endpoints
GET    /api/users/me/progress                 ← NEW: all user progress
GET    /api/classes/{id}/progress             ← NEW: specific class progress
GET    /api/users/me/certificates             ← NEW: user certificates
```

### Response Structure Changes

**Old Response** (Class only):

```json
{
  "id": 1,
  "title": "HTML Basics",
  "content": "..."
}
```

**New Response** (Class detail with slug & progression):

```json
{
  "class": {
    "id": 1,
    "title": "HTML Basics",
    "slug": "html-basics",
    "sequence_number": 1,
    "next_class_id": 2
  },
  "details": {
    "about": "...",
    "rules": "...",
    "tools": ["VS Code"],
    "resource_media": {...}
  },
  "progress": {
    "status": "not_started",
    "progress_percentage": 0,
    "started_at": null
  },
  "next_class": {
    "id": 2,
    "slug": "functions",
    "title": "Functions"
  }
}
```

---

## 🔧 Code Changes Needed

### 1. Update Domain Models

- ✅ `Class` model: Add slug, next_class_id, sequence_number, is_first fields
- ✅ **Create** `ClassDetail` struct
- ✅ **Create** `ClassProgress` struct

### 2. Create Repositories

- ✅ `ClassDetailRepository` - CRUD for class details
- ✅ `ClassProgressRepository` - Get/Create/Update user progress

### 3. Update Services

- ✅ `ClassService` - Add methods for progression
- ✅ **Create** `ProgressService` - Handle start/progress/complete
- ✅ `CertificateService` - Link to ClassProgress

### 4. Update Handlers

- ✅ `ClassHandler` - Add to `/classes/{id}/start`, `/classes/{id}/progress`
- ✅ `ProgressHandler` - Handle progress tracking
- ✅ Update response structures

### 5. Update Routes

```go
// New routes
api.Post("/classes/:id/start", classHandler.StartClass)
api.Put("/classes/:id/progress", classHandler.UpdateProgress)
api.Post("/classes/:id/complete", classHandler.CompleteClass)
api.Get("/users/me/progress", progressHandler.GetMyProgress)

// Updated routes (now uses slug)
api.Get("/projects/:id/classes/:slug", classHandler.GetBySlug)
```

---

## 📊 Migration Guide

### For Existing Data

If you already have data in the old structure, here's how to migrate:

```sql
-- 1. Add new columns to classes
ALTER TABLE classes ADD COLUMN slug VARCHAR(255);
ALTER TABLE classes ADD COLUMN sequence_number INT;
ALTER TABLE classes ADD COLUMN is_first BOOLEAN DEFAULT FALSE;
ALTER TABLE classes ADD COLUMN next_class_id BIGINT;

-- 2. Generate slugs (convert title to slug)
UPDATE classes SET slug = CONCAT(
  LOWER(REPLACE(REPLACE(REPLACE(title, ' ', '-'), '.', '-'), '&', 'and')),
  '-', id
);

-- 3. Add sequence numbers
SET @project_id = 0;
SET @seq = 0;
UPDATE classes
SET sequence_number = (
  CASE
    WHEN @project_id = project_id THEN (@seq := @seq + 1)
    ELSE (@project_id := project_id, @seq := 1)
  END
)
ORDER BY project_id, id;

-- 4. Mark first classes
UPDATE classes SET is_first = TRUE WHERE sequence_number = 1;

-- 5. Link next classes
UPDATE classes c1
SET next_class_id = (
  SELECT c2.id FROM classes c2
  WHERE c2.project_id = c1.project_id
  AND c2.sequence_number = c1.sequence_number + 1
  LIMIT 1
);

-- 6. Create class_details from classes
INSERT INTO class_details (class_id, about, created_at, updated_at)
SELECT id, description, NOW(), NOW() FROM classes;

-- 7. Create class_progress for existing certificates
INSERT INTO class_progress (user_id, class_id, status, completed_at, progress_percentage, created_at, updated_at)
SELECT user_id, class_id, 'completed', issued_at, 100, issued_at, NOW() FROM certificates;
```

---

## 🎯 Example: Setting Up New Structure

### Admin Creates Complete Course:

```go
// 1. Create Category
category := &domain.Category{
    Name: "Web Development",
    Slug: "web-dev",
}
categoryRepo.Create(category)

// 2. Create Project in Category
project := &domain.Project{
    Title:      "JavaScript Mastery",
    CategoryID: category.ID,
    AdminID:    adminID,
}
projectRepo.Create(project)

// 3. Create Classes with Slug & Sequence
class1 := &domain.Class{
    ProjectID:     project.ID,
    Title:         "JS Basics",
    Slug:          "js-basics",        // Unique per project
    SequenceNum:   1,
    IsFirst:       true,
    Difficulty:    "beginner",
    Duration:      120,
}
classRepo.Create(class1)

class2 := &domain.Class{
    ProjectID:     project.ID,
    Title:         "Functions",
    Slug:          "functions",        // Unique per project
    SequenceNum:   2,
    NextClassID:   &class1.ID,         // Link to next
    Difficulty:    "beginner",
    Duration:      90,
}
classRepo.Create(class2)

// Update class1 to link to class2
class1.NextClassID = &class2.ID
classRepo.Update(class1)

// 4. Create Class Details
detail1 := &domain.ClassDetail{
    ClassID: class1.ID,
    About:   "Learn JavaScript fundamentals...",
    Rules:   "1. Install Node.js\n2. Use VS Code",
    Tools:   []string{"VS Code", "Node.js"},
    ResourceMedia: map[string][]string{
        "videos":    {"url1", "url2"},
        "documents": {"syllabus.pdf"},
    },
}
detailRepo.Create(detail1)
```

### User Learns:

```go
// 1. User clicks START on class
progress := &domain.ClassProgress{
    UserID:  userID,
    ClassID: class1.ID,
    Status:  "started",
    StartedAt: time.Now(),
}
progressRepo.Create(progress)

// 2. User studies & updates progress
progress.ProgressPercentage = 50
progress.Status = "in_progress"
progressRepo.Update(progress)

// 3. User completes class
progress.ProgressPercentage = 100
progress.Status = "completed"
progress.CompletedAt = time.Now().Ptr()
progressRepo.Update(progress)

// Create certificate
cert := &domain.Certificate{
    UserID:  userID,
    ClassID: class1.ID,
    Code:    generateCode(),
    IssuedAt: time.Now(),
}
certRepo.Create(cert)

// Award points & check level up
user.TotalPoints += 100
user.Points += 100
if user.Points >= pointsPerLevel {
    user.Points = 0
    user.Level++
}
userRepo.Update(user)

// Next class info returned in response
// User sees: "Congrats! Next: Functions" → link to /project/1/classes/functions
```

---

## ✅ Checklist for Implementation

### Database

- [ ] Create `class_details` table
- [ ] Create `class_progress` table
- [ ] Alter `classes` table (add slug, next_class_id, sequence_number, is_first)
- [ ] Create indexes
- [ ] Migrate existing data (if any)

### Models (Domain)

- [ ] Update `Class` struct
- [ ] Create `ClassDetail` struct
- [ ] Create `ClassProgress` struct
- [ ] Add relationships

### Repositories

- [ ] Create `ClassDetailRepository`
- [ ] Create `ClassProgressRepository`
- [ ] Update `ClassRepository` (add slug queries)

### Services

- [ ] Update `ClassService` (add slug methods)
- [ ] Create `ProgressService` (handle start/progress/complete)
- [ ] Update `CertificateService` (integrate with progress)

### Handlers

- [ ] Update `ClassHandler` (add start/progress/complete methods)
- [ ] Create `ProgressHandler` (if needed)
- [ ] Update response structures

### Routes

- [ ] Add new progress routes
- [ ] Update class routes (use slug)

### API Tests

- [ ] Test category endpoints
- [ ] Test class by slug
- [ ] Test start class
- [ ] Test progress update
- [ ] Test completion
- [ ] Test certificate generation

---

## 📚 Documentation Files

Two comprehensive guides have been created:

1. **PROJECT_CLASSES_FLOW_NEW.md** (Detailed)
   - Complete SQL schemas
   - Go models
   - Detailed flow diagrams
   - Example journey
   - State machines

2. **PROJECT_CLASSES_QUICK_REF_NEW.md** (Quick Reference)
   - Hierarchy diagram
   - Table overview
   - API endpoints summary
   - Example payloads
   - Common scenarios

---

## 🚀 Next Steps

1. Review the new structure in the documentation
2. Update database schema with new tables
3. Create new domain models
4. Implement repositories & services
5. Update handlers & routes
6. Test all endpoints
7. Migrate existing data (if applicable)
