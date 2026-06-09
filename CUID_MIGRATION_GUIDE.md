# CUID Migration Guide (Production-Ready)

## Overview
Migrating from sequential `uint` IDs to distributed `string` CUIDs for all entities.

## Architecture Changes

### 1. Database Schema
**Before (uint auto-increment):**
```sql
users:
  id SERIAL PRIMARY KEY (1, 2, 3, ...)
  
classes:
  id SERIAL PRIMARY KEY (1, 2, 3, ...)
  user_id INT FOREIGN KEY
```

**After (CUID strings):**
```sql
users:
  id VARCHAR(25) PRIMARY KEY (chj6lj2mk0000nj2..., clh6lj2mk0001nj2...)
  
classes:
  id VARCHAR(25) PRIMARY KEY
  user_id VARCHAR(25) FOREIGN KEY
```

### 2. Code Layer Changes

#### JWT Claims
- **Before:** `{ "user_id": 1 (numeric) }`
- **After:** `{ "user_id": "chj6lj2mk0000nj2" (CUID string) }`
- **File:** `pkg/utils/jwt.go` ✅ DONE

#### Middleware
- **Before:** Extract `float64` and cast to `uint`
- **After:** Extract `string` directly  
- **File:** `pkg/middleware/middleware.go` ✅ DONE

#### URL Parameters
- **Before:** `/api/classes/:id` → `ParamsInt("id")` → `uint`
- **After:** `/api/classes/:id` → `Params("id")` → `string`
- **Files:** All handler files (IN PROGRESS)

#### Database Models
- **Before:** `ID uint`
- **After:** `ID string` with `BeforeCreate` hook for auto-generation
- **File:** `internal/domain/models.go` ✅ DONE
- **File:** `internal/domain/cuid_hooks.go` ✅ DONE

### 3. Data Migration Strategy

#### Option A: Zero-Downtime (Recommended for production)
1. **Phase 1:** Add new CUID columns alongside existing uint columns
2. **Phase 2:** Create migration script to populate CUID values
3. **Phase 3:** Update application to write to both columns
4. **Phase 4:** Migrate all existing users/data to new CUID format
5. **Phase 5:** Cutover to CUID-only queries
6. **Phase 6:** Drop uint columns

#### Option B: Drop & Recreate (Current approach - dev environment)
1. Delete all existing data
2. Drop old tables
3. Run migrations with new schema
4. Seed with new CUID-based data
5. Test thoroughly
6. Deploy to production with fresh data or with Phase 1-6 strategy

### 4. Files To Update

#### ✅ COMPLETE
- `internal/domain/models.go` - All ID fields converted to string
- `internal/domain/cuid_hooks.go` - GORM hooks for auto CUID generation  
- `pkg/utils/jwt.go` - JWT generation/parsing with string userID
- `pkg/middleware/middleware.go` - Middleware extracts string userID
- `go.mod` - Added github.com/lucsky/cuid

#### 🔄 IN PROGRESS
- `internal/handler/*.go` - Extract string IDs from context/params
- `internal/service/*.go` - Update all service methods for string IDs
- `internal/repository/*.go` - Update all CRUD operations for string IDs
- `cmd/seed/main.go` - Create new seed data with CUIDs

#### ⏳ TODO (Optional)
- Database migration files for production
- API documentation update
- Swagger spec update for IDs as strings

### 5. Pattern Examples

#### Handler - Before
```go
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
    userID := c.Locals("userID").(uint)  // ← uint
    user, err := h.userSvc.GetUser(c.UserContext(), userID)
    ...
}
```

#### Handler - After
```go
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
    userID := c.Locals("userID").(string)  // ← string CUID
    user, err := h.userSvc.GetUser(c.UserContext(), userID)
    ...
}
```

#### URL Params - Before
```go
classID, _ := c.ParamsInt("id")        // ← returns int
userID, _ := c.QueryInt("user_id")     // ← returns int
```

#### URL Params - After
```go
classID := c.Params("id")              // ← returns string CUID
userID := c.Query("user_id")           // ← returns string CUID
```

#### Repository - Before
```go
func (r *UserRepository) FindByID(id uint) (*domain.User, error) {
    var user domain.User
    return &user, r.db.First(&user, id).Error
}
```

#### Repository - After
```go
func (r *UserRepository) FindByID(id string) (*domain.User, error) {
    var user domain.User
    return &user, r.db.First(&user, id).Error
}
```

### 6. Testing Checklist

- [ ] Build succeeds: `go build ./cmd/api`
- [ ] Database migrations run: `migrate up`
- [ ] Seed data creates users with CUIDs
- [ ] Registration creates CUID-based users
- [ ] JWT tokens contain CUID in user_id claim
- [ ] Login returns JWT with CUID
- [ ] Protected routes parse JWT correctly
- [ ] User endpoints work with CUID in context
- [ ] Admin endpoints enforce RBAC correctly
- [ ] Query parameters accept CUIDs
- [ ] URL path parameters accept CUIDs
- [ ] All foreign key relationships work

### 7. Migration Commands (When Ready)

```bash
# 1. Create new migration
migrate create -ext sql -dir pkg/database/migrations -seq alter_ids_to_cuid

# 2. Manual rollback if needed
migrate down

# 3. Seed with new data
go run cmd/seed/main.go

# 4. Test endpoints
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","email":"test@example.com","password":"pass123"}'
```

## Next Steps

1. Update all handlers to extract string userID and IDs from params
2. Update all service interfaces and implementations
3. Update all repository methods
4. Update seed command
5. Test end-to-end
6. Create proper database migrations
