package domain

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ============================================================================
// CORE USER & AUTHENTICATION MODELS
// ============================================================================

// User represents a system user with roles and gamification tracking
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	Password     string         `json:"-"` // Never expose in JSON
	Name         string         `gorm:"not null" json:"name"`
	Avatar       string         `json:"avatar"`
	Bio          string         `json:"bio"`
	Role         string         `gorm:"default:'user';type:userrole" json:"role"` // admin || user
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	Points       int            `gorm:"default:0" json:"points"`       // Current points
	TotalPoints  int            `gorm:"default:0" json:"total_points"` // Lifetime points
	Level        int            `gorm:"default:1" json:"level"`        // 1-5

	// Relationships
	Projects          []Project          `gorm:"foreignKey:AdminID" json:"-"`
	Certificates      []Certificate      `gorm:"foreignKey:UserID" json:"-"`
	Discussions       []Discussion       `gorm:"foreignKey:UserID" json:"-"`
	Replies           []Reply            `gorm:"foreignKey:UserID" json:"-"`
	Showcases         []Showcase         `gorm:"foreignKey:UserID" json:"-"`
	ShowcaseLikes     []ShowcaseLike     `gorm:"foreignKey:UserID" json:"-"`
	CommunityPoints   []CommunityPoint   `gorm:"foreignKey:UserID" json:"-"`
	ShowcaseComments  []ShowcaseComment  `gorm:"foreignKey:UserID" json:"-"`
}

// ============================================================================
// CATEGORY MODEL (Shared across Project, Class, Showcase, Discussion)
// ============================================================================

// Category represents content categories
type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"uniqueIndex;not null" json:"name"`
	Slug        string         `gorm:"uniqueIndex;not null" json:"slug"`
	Description string         `json:"description"`

	// Relationships (reverse)
	Projects    []Project    `gorm:"foreignKey:CategoryID" json:"-"`
	Showcases   []Showcase   `gorm:"foreignKey:CategoryID" json:"-"`
	Discussions []Discussion `gorm:"foreignKey:CategoryID" json:"-"`
	
}

// ============================================================================
// PROJECT & CLASS MODELS (Admin-managed Learning Content)
// ============================================================================

// Project represents an admin-created learning project containing classes
type Project struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `gorm:"not null;index" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Thumbnail   string         `json:"thumbnail"` // Image URL
	CategoryID  uint           `gorm:"not null;index" json:"category_id"`
	Category    Category       `gorm:"foreignKey:CategoryID" json:"category"`
	AdminID     uint           `gorm:"not null;index" json:"admin_id"` // Only admin can create
	Admin       User           `gorm:"foreignKey:AdminID" json:"admin"`
	Visibility  string         `gorm:"default:'public'" json:"visibility"`

	// Relationships
	Classes []Class `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE" json:"classes,omitempty"`
}

// Class represents a learning class/module under a project
type Class struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `gorm:"not null;index" json:"title"`
	Slug        string         `gorm:"not null" json:"slug"` // ← NEW: URL slug
	Description string         `gorm:"type:text" json:"description"`
	Thumbnail   string         `json:"thumbnail"`
	ProjectID   uint           `gorm:"not null;index" json:"project_id"`
	Project     Project        `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	AdminID     uint           `gorm:"not null;index" json:"admin_id"` // Created by admin
	Admin       User           `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
	Difficulty  string         `gorm:"default:'beginner'" json:"difficulty"`
	Duration    int            `json:"duration"` // In minutes
	OrderIndex  int            `gorm:"default:0" json:"order_index"`
	SequenceNum int            `json:"sequence_number"`
	IsFirst     bool           `json:"is_first"`

	// Progression
	NextClassID *uint          `json:"next_class_id"` // ← NEW: Link to next class
	NextClass   *Class         `gorm:"foreignKey:NextClassID" json:"next_class,omitempty"`

	Visibility  string         `gorm:"default:'public'" json:"visibility"`

	// Relationships
	Details      *ClassDetail    `gorm:"foreignKey:ClassID" json:"details,omitempty"`
	Progress     []ClassProgress `gorm:"foreignKey:ClassID" json:"-"`
	Certificates []Certificate   `gorm:"foreignKey:ClassID;constraint:OnDelete:CASCADE" json:"certificates,omitempty"`
	Discussions  []Discussion    `gorm:"foreignKey:ClassID" json:"discussions,omitempty"`
}

// ClassDetail represents detailed content for a class
type ClassDetail struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ClassID       uint           `gorm:"uniqueIndex;not null" json:"class_id"`
	Class         Class          `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	About         string         `gorm:"type:text" json:"about"`
	Rules         string         `gorm:"type:text" json:"rules"`
	Tools         datatypes.JSON `gorm:"type:json" json:"tools"`          // ["tool1", "tool2"]
	ResourceMedia datatypes.JSON `gorm:"type:json" json:"resource_media"`  // { "videos": [...], "documents": [...], "images": [...] }
	Resources     datatypes.JSON `gorm:"type:json" json:"resources"`       // [ { "type": "pdf", "title": "...", "url": "..." } ]
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// ResourceMedia structure helper for JSON parsing
type ResourceMedia struct {
	Videos    []string `json:"videos"`
	Documents []string `json:"documents"`
	Images    []string `json:"images"`
}

// Resource structure helper for JSON parsing
type Resource struct {
	Type  string `json:"type"`  // pdf | video | link | document
	Title string `json:"title"`
	URL   string `json:"url"`
}

// ClassProgress tracks user learning progress in a class
type ClassProgress struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	UserID             uint           `gorm:"not null;index" json:"user_id"`
	User               User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ClassID            uint           `gorm:"not null;index" json:"class_id"`
	Class              Class          `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	Status             string         `gorm:"default:'not_started'" json:"status"` // not_started | started | in_progress | completed
	StartedAt          *time.Time     `json:"started_at"`
	CompletedAt        *time.Time     `json:"completed_at"`
	ProgressPercentage int            `json:"progress_percentage"` // 0-100
	Notes              string         `json:"notes"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

// ============================================================================
// CERTIFICATE MODEL (User-specific, Private)
// ============================================================================

// Certificate represents user achievement/completion of a class
// IMPORTANT: Only accessible to the user who owns it + admin
type Certificate struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ClassID     uint           `gorm:"not null;index" json:"class_id"`
	Class       Class          `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	UniqueCode  string         `gorm:"uniqueIndex;not null" json:"unique_code"` // UUID or slug
	BadgeURL    string         `json:"badge_url"`
	IssuedAt    time.Time      `json:"issued_at"`
	ExpiresAt   *time.Time     `json:"expires_at"` // Optional expiration
}

// ============================================================================
// DISCUSSION & REPLY MODELS (Community Discussion)
// ============================================================================

// Discussion represents a discussion thread (usually in a class context)
type Discussion struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `gorm:"not null;index" json:"title"`
	Content     string         `gorm:"type:text;not null" json:"content"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ClassID     *uint          `gorm:"index" json:"class_id"`              // Optional - can be standalone
	Class       *Class         `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	CategoryID  uint           `gorm:"index" json:"category_id"`           // For filtering
	Category    Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	ViewsCount  int            `gorm:"default:0" json:"views_count"`
	Status      string         `gorm:"default:'open';type:discussionstatus" json:"status"` // Can close discussion
	IsPinned    bool           `gorm:"default:false" json:"is_pinned"`

	// Relationships
	Replies []Reply `gorm:"foreignKey:DiscussionID;constraint:OnDelete:CASCADE" json:"replies,omitempty"`
}

// Reply represents a comment/reply in a discussion
// Can be nested (parent_id points to another reply for threaded discussions)
type Reply struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Content       string         `gorm:"type:text;not null" json:"content"`
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	User          User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	DiscussionID  uint           `gorm:"not null;index" json:"discussion_id"`
	Discussion    Discussion     `gorm:"foreignKey:DiscussionID" json:"discussion,omitempty"`
	ParentID      *uint          `gorm:"index" json:"parent_id"`                   // For nested replies
	Parent        *Reply         `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	LikesCount    int            `gorm:"default:0" json:"likes_count"`
	IsMarkedBest  bool           `gorm:"default:false" json:"is_marked_best"` // Discussion owner can mark best answer

	// Relationships
	ChildReplies []Reply `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"child_replies,omitempty"`
}

// ============================================================================
// SHOWCASE MODELS (User Portfolio & Performance)
// ============================================================================

// Showcase represents user portfolio item (projects completed, work display)
type Showcase struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Title        string         `gorm:"not null;index" json:"title"`
	Description  string         `gorm:"type:text" json:"description"`
	MediaURLs    datatypes.JSON `gorm:"type:jsonb" json:"media_urls"` // JSON array of image/video URLs
	UserID       uint           `gorm:"not null;index" json:"user_id"`
	User         User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CategoryID   uint           `gorm:"not null;index" json:"category_id"`
	Category     Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Status       string         `gorm:"default:'published';type:showcasestatus" json:"status"`
	Visibility   string         `gorm:"default:'public';type:showcasevisibility" json:"visibility"`
	LikesCount   int            `gorm:"default:0" json:"likes_count"`
	ViewsCount   int            `gorm:"default:0" json:"views_count"`

	// Relationships
	Likes     []ShowcaseLike    `gorm:"foreignKey:ShowcaseID;constraint:OnDelete:CASCADE" json:"likes,omitempty"`
	Comments  []ShowcaseComment `gorm:"foreignKey:ShowcaseID;constraint:OnDelete:CASCADE" json:"comments,omitempty"`
}

// ShowcaseLike represents a like on a showcase (composite key: user_id + showcase_id)
type ShowcaseLike struct {
	UserID     uint      `gorm:"primaryKey" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ShowcaseID uint      `gorm:"primaryKey" json:"showcase_id"`
	Showcase   Showcase  `gorm:"foreignKey:ShowcaseID" json:"showcase,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// ShowcaseComment represents comments on showcase items
type ShowcaseComment struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Content      string         `gorm:"type:text;not null" json:"content"`
	UserID       uint           `gorm:"not null;index" json:"user_id"`
	User         User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ShowcaseID   uint           `gorm:"not null;index" json:"showcase_id"`
	Showcase     Showcase       `gorm:"foreignKey:ShowcaseID" json:"showcase,omitempty"`
	ParentID     *uint          `gorm:"index" json:"parent_id"` // For nested comments
	Parent       *ShowcaseComment `gorm:"foreignKey:ParentID" json:"parent,omitempty"`

	// Relationships
	ChildComments []ShowcaseComment `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"child_comments,omitempty"`
}

// ============================================================================
// GAMIFICATION MODELS (Points & Levels)
// ============================================================================

// CommunityPoint represents activity log and points transaction
type CommunityPoint struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	User          User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ActivityType  string         `gorm:"not null;index;type:pointactivitytype" json:"activity_type"`
	PointsEarned  int            `gorm:"not null" json:"points_earned"`
	PointsAfter   int            `gorm:"not null" json:"points_after"` // Total after this activity
	LevelAfter    int            `json:"level_after"`                   // In case of level up
	Metadata      datatypes.JSON `gorm:"type:jsonb" json:"metadata"`    // {object_id, object_type, etc}
	Description   string         `json:"description"`                   // Human readable
}

// UserLevel represents level configuration and requirements
type UserLevel struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Level       int       `gorm:"uniqueIndex;not null" json:"level"` // 1, 2, 3, 4, 5
	MinPoints   int       `gorm:"not null;uniqueIndex" json:"min_points"`
	MaxPoints   int       `gorm:"not null" json:"max_points"`
	BadgeName   string    `json:"badge_name"`
	BadgeIcon   string    `json:"badge_icon"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ============================================================================
// HELPER METHODS
// ============================================================================

// TableName specifies table name for ShowcaseLike (composite key table)
func (ShowcaseLike) TableName() string {
	return "showcase_likes"
}

// String implements Stringer for UserLevel
func (ul UserLevel) String() string {
	return ul.BadgeName
}

// ============================================================================
// VALIDATION & PERMISSIONS
// ============================================================================

// IsUserOwnerOfCertificate checks if user owns this certificate (privacy check)
func (c *Certificate) IsUserOwnerOfCertificate(userID uint) bool {
	return c.UserID == userID
}

// IsUserOwnerOfShowcase checks if user owns this showcase
func (s *Showcase) IsUserOwnerOfShowcase(userID uint) bool {
	return s.UserID == userID
}

// CanUserEditShowcase checks if user can edit showcase (owner or admin)
func (s *Showcase) CanUserEditShowcase(user *User) bool {
	return s.UserID == user.ID || user.Role == "admin"
}

// CanUserCreateProject checks if user can create project (admin only)
func (p *Project) CanUserCreateProject(user *User) bool {
	return user.Role == "admin"
}

// ============================================================================
// AUTO MIGRATION
// ============================================================================

// AutoMigrate runs all database migrations
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Category{},
		&Project{},
		&Class{},
		&ClassDetail{},
		&ClassProgress{},
		&Certificate{},
		&Discussion{},
		&Reply{},
		&Showcase{},
		&ShowcaseLike{},
		&ShowcaseComment{},
		&CommunityPoint{},
		&UserLevel{},
	)
}