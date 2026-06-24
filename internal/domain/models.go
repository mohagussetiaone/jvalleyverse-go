package domain

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Email       string         `gorm:"uniqueIndex;not null" json:"email"`
	Password    string         `json:"-"` // Never expose in JSON
	Name        string         `gorm:"not null" json:"name"`
	Avatar      string         `json:"avatar"`
	Bio         string         `json:"bio"`
	Role        string         `gorm:"default:'user';type:userrole" json:"role"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	Points      int            `gorm:"default:0" json:"points"`
	TotalPoints int            `gorm:"default:0" json:"total_points"`
	Level       int            `gorm:"default:1" json:"level"`

	Courses          []Course          `gorm:"foreignKey:AdminID" json:"-"`
	Certificates     []Certificate     `gorm:"foreignKey:UserID" json:"-"`
	Discussions      []Discussion      `gorm:"foreignKey:UserID" json:"-"`
	Replies          []Reply           `gorm:"foreignKey:UserID" json:"-"`
	Showcases        []Showcase        `gorm:"foreignKey:UserID" json:"-"`
	ShowcaseLikes    []ShowcaseLike    `gorm:"foreignKey:UserID" json:"-"`
	CommunityPoints  []CommunityPoint  `gorm:"foreignKey:UserID" json:"-"`
	MentorCourses    []Course          `gorm:"foreignKey:MentorID" json:"-"`
	Reviews          []Review          `gorm:"foreignKey:UserID" json:"-"`
	ShowcaseComments []ShowcaseComment `gorm:"foreignKey:UserID" json:"-"`
}

type Category struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"uniqueIndex;not null" json:"name"`
	Slug        string         `gorm:"uniqueIndex;not null" json:"slug"`
	Description string         `json:"description"`

	Courses     []Course     `gorm:"foreignKey:CategoryID" json:"-"`
	Showcases   []Showcase   `gorm:"foreignKey:CategoryID" json:"-"`
	Discussions []Discussion `gorm:"foreignKey:CategoryID" json:"-"`
}

// Course represents an admin-created learning course containing sections
type Course struct {
	ID                 string         `gorm:"primaryKey" json:"id"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
	Title              string         `gorm:"not null;index" json:"title"`
	Description        string         `gorm:"type:text" json:"description"`
	Thumbnail          string         `json:"thumbnail"`
	CategoryID         string         `gorm:"not null;index" json:"category_id"`
	Category           Category       `gorm:"foreignKey:CategoryID" json:"category"`
	AdminID            string         `gorm:"not null;index" json:"admin_id"`
	Admin              User           `gorm:"foreignKey:AdminID" json:"admin"`
	MentorID           string         `gorm:"index" json:"mentor_id"`
	Mentor             User           `gorm:"foreignKey:MentorID" json:"mentor,omitempty"`
	Visibility         string         `gorm:"default:'public'" json:"visibility"`
	Price              float64        `gorm:"default:0" json:"price"`
	Hours              int            `json:"hours"`
	LearningObjectives datatypes.JSON `gorm:"type:json" json:"learning_objectives"`

	Sections []Section `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"sections,omitempty"`
	Reviews  []Review  `gorm:"foreignKey:CourseID" json:"-"`
}

// Section represents a learning section under a course, containing lessons
type Section struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `gorm:"not null;index" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	CourseID    string         `gorm:"not null;index" json:"course_id"`
	Course      Course         `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	OrderIndex  int            `gorm:"default:0" json:"order_index"`

	Lessons []Lesson `gorm:"foreignKey:SectionID;constraint:OnDelete:CASCADE" json:"lessons,omitempty"`
}

// Lesson represents a learning lesson under a section
type Lesson struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `gorm:"not null;index" json:"title"`
	Slug        string         `gorm:"not null" json:"slug"`
	Description string         `gorm:"type:text" json:"description"`
	Thumbnail   string         `json:"thumbnail"`
	CourseID    string         `gorm:"not null;index" json:"course_id"`
	Course      Course         `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	SectionID   string         `gorm:"not null;index" json:"section_id"`
	Section     Section        `gorm:"foreignKey:SectionID" json:"section,omitempty"`
	AdminID     string         `gorm:"not null;index" json:"admin_id"`
	Admin       User           `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
	Difficulty  string         `gorm:"default:'beginner'" json:"difficulty"`
	Duration    int            `json:"duration"`
	OrderIndex  int            `gorm:"default:0" json:"order_index"`
	SequenceNum int            `json:"sequence_number"`
	IsFirst     bool           `json:"is_first"`

	VideoURL string `json:"video_url"`

	NextLessonID *string `json:"next_lesson_id"`
	NextLesson   *Lesson `gorm:"foreignKey:NextLessonID" json:"next_lesson,omitempty"`

	Visibility string `gorm:"default:'public'" json:"visibility"`

	Details      *LessonDetail    `gorm:"foreignKey:LessonID" json:"details,omitempty"`
	Progress     []LessonProgress `gorm:"foreignKey:LessonID" json:"-"`
	Certificates []Certificate    `gorm:"foreignKey:LessonID;constraint:OnDelete:CASCADE" json:"certificates,omitempty"`
	Discussions  []Discussion     `gorm:"foreignKey:LessonID" json:"discussions,omitempty"`
	Reviews      []Review         `gorm:"foreignKey:LessonID" json:"-"`
}

// LessonDetail represents detailed content for a lesson
type LessonDetail struct {
	ID            string         `gorm:"primaryKey" json:"id"`
	LessonID      string         `gorm:"uniqueIndex;not null" json:"lesson_id"`
	Lesson        Lesson         `gorm:"foreignKey:LessonID" json:"lesson,omitempty"`
	About         string         `gorm:"type:text" json:"about"`
	Rules         string         `gorm:"type:text" json:"rules"`
	Tools         datatypes.JSON `gorm:"type:json" json:"tools"`
	ResourceMedia datatypes.JSON `gorm:"type:json" json:"resource_media"`
	Resources     datatypes.JSON `gorm:"type:json" json:"resources"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type ResourceMedia struct {
	Videos    []string `json:"videos"`
	Documents []string `json:"documents"`
	Images    []string `json:"images"`
}

type Resource struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

// LessonProgress tracks user learning progress in a lesson
type LessonProgress struct {
	ID                 string     `gorm:"primaryKey" json:"id"`
	UserID             string     `gorm:"not null;index" json:"user_id"`
	User               User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	LessonID           string     `gorm:"not null;index" json:"lesson_id"`
	Lesson             Lesson     `gorm:"foreignKey:LessonID" json:"lesson,omitempty"`
	Status             string     `gorm:"default:'not_started'" json:"status"`
	StartedAt          *time.Time `json:"started_at"`
	CompletedAt        *time.Time `json:"completed_at"`
	ProgressPercentage int        `json:"progress_percentage"`
	Notes              string     `json:"notes"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type Certificate struct {
	ID         string         `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	UserID     string         `gorm:"not null;index" json:"user_id"`
	User       User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	LessonID   string         `gorm:"not null;index" json:"lesson_id"`
	Lesson     Lesson         `gorm:"foreignKey:LessonID" json:"lesson,omitempty"`
	UniqueCode string         `gorm:"uniqueIndex;not null" json:"unique_code"`
	BadgeURL   string         `json:"badge_url"`
	IssuedAt   time.Time      `json:"issued_at"`
	ExpiresAt  *time.Time     `json:"expires_at"`
}

type Discussion struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `gorm:"not null;index" json:"title"`
	Content     string         `gorm:"type:text;not null" json:"content"`
	UserID      string         `gorm:"not null;index" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	LessonID    *string        `gorm:"index" json:"lesson_id"`
	Lesson      *Lesson        `gorm:"foreignKey:LessonID" json:"lesson,omitempty"`
	StudyCaseID *string        `gorm:"index" json:"study_case_id"`
	StudyCase   *StudyCase     `gorm:"foreignKey:StudyCaseID" json:"study_case,omitempty"`
	CategoryID  string         `gorm:"index" json:"category_id"`
	Category    Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	ViewsCount  int            `gorm:"default:0" json:"views_count"`
	Status      string         `gorm:"default:'open';type:discussionstatus" json:"status"`
	IsPinned    bool           `gorm:"default:false" json:"is_pinned"`

	Replies []Reply `gorm:"foreignKey:DiscussionID;constraint:OnDelete:CASCADE" json:"replies,omitempty"`
}

type Reply struct {
	ID           string         `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Content      string         `gorm:"type:text;not null" json:"content"`
	UserID       string         `gorm:"not null;index" json:"user_id"`
	User         User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	DiscussionID string         `gorm:"not null;index" json:"discussion_id"`
	Discussion   Discussion     `gorm:"foreignKey:DiscussionID" json:"discussion,omitempty"`
	ParentID     *string        `gorm:"index" json:"parent_id"`
	Parent       *Reply         `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	LikesCount   int            `gorm:"default:0" json:"likes_count"`
	IsMarkedBest bool           `gorm:"default:false" json:"is_marked_best"`

	ChildReplies []Reply `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"child_replies,omitempty"`
}

type Showcase struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `gorm:"not null;index" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	MediaURLs   datatypes.JSON `gorm:"type:jsonb" json:"media_urls"`
	UserID      string         `gorm:"not null;index" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CategoryID  string         `gorm:"not null;index" json:"category_id"`
	Category    Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Status      string         `gorm:"default:'published';type:showcasestatus" json:"status"`
	Visibility  string         `gorm:"default:'public';type:showcasevisibility" json:"visibility"`
	LikesCount  int            `gorm:"default:0" json:"likes_count"`
	ViewsCount  int            `gorm:"default:0" json:"views_count"`

	Likes    []ShowcaseLike    `gorm:"foreignKey:ShowcaseID;constraint:OnDelete:CASCADE" json:"likes,omitempty"`
	Comments []ShowcaseComment `gorm:"foreignKey:ShowcaseID;constraint:OnDelete:CASCADE" json:"comments,omitempty"`
}

type ShowcaseLike struct {
	UserID     string    `gorm:"primaryKey" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ShowcaseID string    `gorm:"primaryKey" json:"showcase_id"`
	Showcase   Showcase  `gorm:"foreignKey:ShowcaseID" json:"showcase,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type ShowcaseComment struct {
	ID         string           `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
	DeletedAt  gorm.DeletedAt   `gorm:"index" json:"-"`
	Content    string           `gorm:"type:text;not null" json:"content"`
	UserID     string           `gorm:"not null;index" json:"user_id"`
	User       User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ShowcaseID string           `gorm:"not null;index" json:"showcase_id"`
	Showcase   Showcase         `gorm:"foreignKey:ShowcaseID" json:"showcase,omitempty"`
	ParentID   *string          `gorm:"index" json:"parent_id"`
	Parent     *ShowcaseComment `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
}

type Blog struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `gorm:"not null;index" json:"title"`
	Slug        string         `gorm:"uniqueIndex" json:"slug"`
	Description string         `gorm:"type:text" json:"description"`
	Content     string         `gorm:"type:text" json:"content"`
	CoverImgURL string         `json:"cover_img_url"`
	Tags        datatypes.JSON `gorm:"type:jsonb" json:"tags"`
	Status        string         `gorm:"default:draft" json:"status"`
	UserID      string         `gorm:"not null;index" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"author,omitempty"`
	CategoryID  string         `gorm:"not null;index" json:"category_id"`
	Category    Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (Blog) TableName() string {
	return "blogs"
}

type Review struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	UserID    string         `gorm:"not null;index" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CourseID  string         `gorm:"index" json:"course_id"`
	Course    Course         `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	LessonID  string         `gorm:"index" json:"lesson_id"`
	Lesson    Lesson         `gorm:"foreignKey:LessonID" json:"lesson,omitempty"`
	Rating    int            `gorm:"not null" json:"rating"`
	Message   string         `gorm:"type:text" json:"message"`
}

type CommunityPoint struct {
	ID           string         `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	UserID       string         `gorm:"not null;index" json:"user_id"`
	User         User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ActivityType string         `gorm:"not null;index;type:pointactivitytype" json:"activity_type"`
	PointsEarned int            `gorm:"not null" json:"points_earned"`
	PointsAfter  int            `gorm:"not null" json:"points_after"`
	LevelAfter   int            `json:"level_after"`
	Metadata     datatypes.JSON `gorm:"type:jsonb" json:"metadata"`
	Description  string         `json:"description"`
}

type RefreshToken struct {
	ID        string     `gorm:"primaryKey" json:"id"`
	UserID    string     `gorm:"not null;index" json:"user_id"`
	User      User       `gorm:"foreignKey:UserID" json:"-"`
	Token     string     `gorm:"uniqueIndex;not null" json:"-"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	RevokedAt *time.Time `json:"-"`
	CreatedAt time.Time  `json:"created_at"`
}

type AdminAuditLog struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	AdminID      string    `gorm:"not null;index" json:"admin_id"`
	Admin        User      `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
	Action       string    `gorm:"not null;index" json:"action"`
	ResourceType string    `gorm:"not null;index" json:"resource_type"`
	ResourceID   string    `gorm:"index" json:"resource_id"`
	Details      string    `gorm:"type:text" json:"details,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserLevel struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Level       int       `gorm:"uniqueIndex;not null" json:"level"`
	MinPoints   int       `gorm:"not null;uniqueIndex" json:"min_points"`
	MaxPoints   int       `gorm:"not null" json:"max_points"`
	BadgeName   string    `json:"badge_name"`
	BadgeIcon   string    `json:"badge_icon"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// StudyCase represents a case study project
type StudyCase struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"not null;index" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	ImgURL      string         `json:"img_url"`
	Tags        datatypes.JSON `gorm:"type:jsonb" json:"tags"`
	YoutubeURL  string         `json:"youtube_url"`
	CategoryID  *string        `gorm:"index" json:"category_id"`
	Category    Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	UserID      string         `gorm:"not null;index" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`

	Discussions []Discussion `gorm:"foreignKey:StudyCaseID" json:"discussions,omitempty"`
}

// CourseEnrollment tracks user enrollment in a course
type CourseEnrollment struct {
	ID             string         `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	UserID         string         `gorm:"not null;uniqueIndex:idx_user_course" json:"user_id"`
	User           User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CourseID       string         `gorm:"not null;uniqueIndex:idx_user_course" json:"course_id"`
	Course         Course         `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	LastLessonID   *string        `json:"last_lesson_id"`
	OriginalPrice  float64        `gorm:"default:0" json:"original_price"`
	DiscountAmount float64        `gorm:"default:0" json:"discount_amount"`
	DiscountCode   string         `json:"discount_code"`
}

// Notification represents a user notification
type Notification struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	UserID    string         `gorm:"not null;index" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Type      string         `gorm:"not null;index" json:"type"`
	Title     string         `gorm:"not null" json:"title"`
	Message   string         `gorm:"type:text" json:"message"`
	IsRead    bool           `gorm:"default:false" json:"is_read"`
	Link      string         `json:"link"`
	Metadata  datatypes.JSON `gorm:"type:jsonb" json:"metadata,omitempty"`
}

// Helper methods

func (ShowcaseLike) TableName() string {
	return "showcase_likes"
}

func (ul UserLevel) String() string {
	return ul.BadgeName
}

func (c *Certificate) IsUserOwnerOfCertificate(userID string) bool {
	return c.UserID == userID
}

func (s *Showcase) IsUserOwnerOfShowcase(userID string) bool {
	return s.UserID == userID
}

func (s *Showcase) CanUserEditShowcase(user *User) bool {
	return s.UserID == user.ID || user.Role == "admin"
}

func (p *Course) CanUserCreateCourse(user *User) bool {
	return user.Role == "admin"
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Blog{},
		&User{},
		&Category{},
		&Course{},
		&Section{},
		&StudyCase{},
		&Lesson{},
		&LessonDetail{},
		&LessonProgress{},
		&Certificate{},
		&Discussion{},
		&Reply{},
		&Showcase{},
		&ShowcaseLike{},
		&ShowcaseComment{},
		&CommunityPoint{},
		&UserLevel{},
		&Review{},
		&RefreshToken{},
		&AdminAuditLog{},
		&CourseEnrollment{},
		&Notification{},
	)
}
