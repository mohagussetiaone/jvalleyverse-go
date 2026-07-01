package service

import (
	"encoding/json"
	"fmt"
	"jvalleyverse/internal/domain"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// SeedInitialData populates the database with initial data for all tables.
// Uses FirstOrCreate to make seeding idempotent (safe to run multiple times).
func SeedInitialData(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database connection not initialized")
	}

	fmt.Println("🌱 Starting database seeding...")

	// =========================================================================
	// 1. SEED USERS
	// =========================================================================
	fmt.Println("  → Seeding users...")
	users, err := seedUsers(db)
	if err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	// =========================================================================
	// 2. SEED USER LEVELS
	// =========================================================================
	fmt.Println("  → Seeding user levels...")
	if err := seedUserLevels(db); err != nil {
		return fmt.Errorf("failed to seed user levels: %w", err)
	}

	// =========================================================================
	// 3. SEED CATEGORIES
	// =========================================================================
	fmt.Println("  → Seeding categories...")
	categories, err := seedCategories(db)
	if err != nil {
		return fmt.Errorf("failed to seed categories: %w", err)
	}

	// =========================================================================
	// 4. SEED COURSES
	// =========================================================================
	fmt.Println("  → Seeding courses...")
	courses, err := seedCourses(db, users["admin"], users["mentor1"], users["mentor2"], categories)
	if err != nil {
		return fmt.Errorf("failed to seed courses: %w", err)
	}

	// =========================================================================
	// 5. SEED SECTIONS
	// =========================================================================
	fmt.Println("  → Seeding sections...")
	sections, err := seedSections(db, courses)
	if err != nil {
		return fmt.Errorf("failed to seed sections: %w", err)
	}

	// =========================================================================
	// 6. SEED LESSONS + LESSON DETAILS
	// =========================================================================
	fmt.Println("  → Seeding lessons & lesson details...")
	lessons, err := seedLessons(db, users["admin"], courses, sections)
	if err != nil {
		return fmt.Errorf("failed to seed lessons: %w", err)
	}

	// =========================================================================
	// 6. SEED LESSON PROGRESS
	// =========================================================================
	fmt.Println("  → Seeding lesson progress...")
	if err := seedClassProgress(db, users, lessons); err != nil {
		return fmt.Errorf("failed to seed lesson progress: %w", err)
	}

	// =========================================================================
	// 7. SEED DISCUSSIONS & REPLIES
	// =========================================================================
	fmt.Println("  → Seeding discussions & replies...")
	if err := seedDiscussions(db, users, lessons, categories); err != nil {
		return fmt.Errorf("failed to seed discussions: %w", err)
	}

	// =========================================================================
	// 8. SEED SHOWCASES
	// =========================================================================
	fmt.Println("  → Seeding showcases...")
	if err := seedShowcases(db, users, categories); err != nil {
		return fmt.Errorf("failed to seed showcases: %w", err)
	}

	// =========================================================================
	// 9. SEED STUDY CASES
	// =========================================================================
	fmt.Println("  → Seeding study cases...")
	studyCases, err := seedStudyCases(db, users)
	if err != nil {
		return fmt.Errorf("failed to seed study cases: %w", err)
	}

	// =========================================================================
	// 10. SEED STUDY CASE DISCUSSIONS
	// =========================================================================
	fmt.Println("  → Seeding study case discussions...")
	if err := seedStudyCaseDiscussions(db, users, studyCases, categories); err != nil {
		return fmt.Errorf("failed to seed study case discussions: %w", err)
	}

	// =========================================================================
	// 11. SEED REVIEWS
	// =========================================================================
	fmt.Println("  → Seeding reviews...")
	if err := seedReviews(db, users, courses, lessons); err != nil {
		return fmt.Errorf("failed to seed reviews: %w", err)
	}

	// =========================================================================
	// 12. SEED BLOGS
	// =========================================================================
	fmt.Println("  → Seeding blogs...")
	if err := seedBlogs(db, users, categories); err != nil {
		return fmt.Errorf("failed to seed blogs: %w", err)
	}
	// =========================================================================
	// 13. SEED COMPANY PROFILE
	// =========================================================================
	fmt.Println("  → Seeding company profile...")
	if err := seedCompany(db); err != nil {
		return fmt.Errorf("failed to seed company: %w", err)
	}
	if err := seedFAQs(db); err != nil {
		return fmt.Errorf("failed to seed faqs: %w", err)
	}

	fmt.Println("🌱 Seeding completed!")
	return nil
}

// ============================================================================
// 1. SEED USERS
// ============================================================================

func seedUsers(db *gorm.DB) (map[string]*domain.User, error) {
	hashPassword := func(pw string) string {
		hash, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
		return string(hash)
	}

	usersData := []domain.User{
		{
			Email:       "admin@jvalleyverse.com",
			Password:    hashPassword("Admin@123"),
			Name:        "JValley Admin",
			Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=admin",
			Bio:         "Platform administrator & lead instructor at JValleyverse",
			Role:        "admin",
			IsActive:    true,
			Points:      500,
			TotalPoints: 500,
			Level:       3,
		},
		{
			Email:       "budi@example.com",
			Password:    hashPassword("User@123"),
			Name:        "Budi Santoso",
			Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=budi",
			Bio:         "Full-stack developer yang sedang belajar Go & React",
			Role:        "user",
			IsActive:    true,
			Points:      150,
			TotalPoints: 150,
			Level:       2,
		},
		{
			Email:       "siti@example.com",
			Password:    hashPassword("User@123"),
			Name:        "Siti Rahayu",
			Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=siti",
			Bio:         "UI/UX Designer belajar frontend development",
			Role:        "user",
			IsActive:    true,
			Points:      75,
			TotalPoints: 75,
			Level:       1,
		},
		{
			Email:       "andi@example.com",
			Password:    hashPassword("User@123"),
			Name:        "Andi Pratama",
			Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=andi",
			Bio:         "Backend engineer, suka ngoprek database & microservices",
			Role:        "user",
			IsActive:    true,
			Points:      320,
			TotalPoints: 320,
			Level:       2,
		},
		{
			Email:       "dewi@example.com",
			Password:    hashPassword("User@123"),
			Name:        "Dewi Lestari",
			Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=dewi",
			Bio:         "Fresh graduate, antusias belajar programming",
			Role:        "user",
			IsActive:    true,
			Points:      30,
			TotalPoints: 30,
			Level:       1,
		},
		{
			Email:       "mentor1@jvalleyverse.com",
			Password:    hashPassword("Mentor@123"),
			Name:        "Eko Prasetyo",
			Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=eko",
			Bio:         "Senior Backend Engineer & Go enthusiast. 10+ years experience in building scalable systems.",
			Role:        "mentor",
			IsActive:    true,
			Points:      2500,
			TotalPoints: 2500,
			Level:       5,
		},
		{
			Email:       "mentor2@jvalleyverse.com",
			Password:    hashPassword("Mentor@123"),
			Name:        "Fitri Handayani",
			Avatar:      "https://api.dicebear.com/7.x/avataaars/svg?seed=fitri",
			Bio:         "Frontend Architect & UI/UX specialist. Expert in React, Vue, and modern CSS.",
			Role:        "mentor",
			IsActive:    true,
			Points:      2000,
			TotalPoints: 2000,
			Level:       4,
		},
	}

	result := make(map[string]*domain.User)
	keys := []string{"admin", "budi", "siti", "andi", "dewi", "mentor1", "mentor2"}

	for i, userData := range usersData {
		user := userData // create a copy
		if err := db.Where("email = ?", user.Email).FirstOrCreate(&user).Error; err != nil {
			return nil, err
		}
		result[keys[i]] = &user
		fmt.Printf("    ✓ User: %s (ID: %s, Role: %s)\n", user.Name, user.ID, user.Role)
	}

	return result, nil
}

// ============================================================================
// 2. SEED USER LEVELS
// ============================================================================

func seedUserLevels(db *gorm.DB) error {
	levels := []domain.UserLevel{
		{Level: 1, MinPoints: 0, MaxPoints: 99, BadgeName: "Beginner", BadgeIcon: "🌱", Description: "Baru memulai perjalanan belajar"},
		{Level: 2, MinPoints: 100, MaxPoints: 499, BadgeName: "Intermediate", BadgeIcon: "🌿", Description: "Mulai menguasai dasar-dasar"},
		{Level: 3, MinPoints: 500, MaxPoints: 999, BadgeName: "Advanced", BadgeIcon: "🌳", Description: "Sudah mahir dan berpengalaman"},
		{Level: 4, MinPoints: 1000, MaxPoints: 1999, BadgeName: "Expert", BadgeIcon: "⭐", Description: "Ahli dan menjadi panutan"},
		{Level: 5, MinPoints: 2000, MaxPoints: 99999, BadgeName: "Master", BadgeIcon: "👑", Description: "Penguasa ilmu tertinggi"},
	}

	for _, level := range levels {
		l := level
		if err := db.Where("level = ?", l.Level).FirstOrCreate(&l).Error; err != nil {
			return err
		}
		fmt.Printf("    ✓ Level %d: %s %s (%d-%d pts)\n", l.Level, l.BadgeIcon, l.BadgeName, l.MinPoints, l.MaxPoints)
	}

	return nil
}

// ============================================================================
// 3. SEED CATEGORIES
// ============================================================================

func seedCategories(db *gorm.DB) (map[string]*domain.Category, error) {
	categoriesData := []domain.Category{
		{Name: "Web Development", Slug: "web-development", Description: "Belajar membangun website modern dari frontend hingga backend"},
		{Name: "Mobile Development", Slug: "mobile-development", Description: "Belajar membuat aplikasi mobile Android & iOS"},
		{Name: "Backend Development", Slug: "backend-development", Description: "Belajar server-side programming, API, dan database"},
		{Name: "Frontend Development", Slug: "frontend-development", Description: "Belajar HTML, CSS, JavaScript dan framework modern"},
		{Name: "DevOps & Cloud", Slug: "devops-cloud", Description: "Belajar deployment, CI/CD, Docker, dan cloud computing"},
		{Name: "Data Science", Slug: "data-science", Description: "Belajar analisis data, machine learning, dan AI"},
		{Name: "UI/UX Design", Slug: "ui-ux-design", Description: "Belajar desain antarmuka dan pengalaman pengguna"},
	}

	result := make(map[string]*domain.Category)

	for _, catData := range categoriesData {
		cat := catData
		if err := db.Where("slug = ?", cat.Slug).FirstOrCreate(&cat).Error; err != nil {
			return nil, err
		}
		result[cat.Slug] = &cat
		fmt.Printf("    ✓ Category: %s (ID: %s)\n", cat.Name, cat.ID)
	}

	return result, nil
}

// ============================================================================
// 4. SEED COURSES
// ============================================================================

func seedCourses(db *gorm.DB, admin, mentor1, mentor2 *domain.User, categories map[string]*domain.Category) (map[string]*domain.Course, error) {
	coursesData := []struct {
		key         string
		categoryKey string
		mentor      *domain.User
		course      domain.Course
	}{
		{
			key:         "go-rest-api",
			categoryKey: "backend-development",
			mentor:      mentor1,
			course: domain.Course{
				Title:       "Membangun REST API dengan Go & Fiber",
				Description: "Pelajari cara membuat RESTful API yang scalable menggunakan Go, Fiber framework, GORM, dan PostgreSQL. Dari setup project hingga deployment.",
				Thumbnail:   "https://images.unsplash.com/photo-1555066931-4365d14bab8c?w=800",
				Visibility:  "public",
				LearningObjectives: mustJSON([]string{
					"Setting up the Go development environment",
					"Understanding HTTP fundamentals & REST principles",
					"Building RESTful endpoints with Fiber",
					"Database integration with GORM & PostgreSQL",
					"Implementing authentication & middleware",
					"Deploying Go API to production",
				}),
			},
		},
		{
			key:         "react-dashboard",
			categoryKey: "frontend-development",
			mentor:      mentor2,
			course: domain.Course{
				Title:       "Build Modern Dashboard dengan React & TypeScript",
				Description: "Membangun dashboard interaktif menggunakan React, TypeScript, Tailwind CSS, dan Chart.js. Lengkap dengan autentikasi dan state management.",
				Thumbnail:   "https://images.unsplash.com/photo-1551288049-bebda4e38f71?w=800",
				Visibility:  "public",
				LearningObjectives: mustJSON([]string{
					"Setting up React with TypeScript & Vite",
					"Building reusable UI components with Tailwind",
					"State management with React Context & Redux",
					"Data visualization with Chart.js",
					"Implementing authentication & routing",
					"Building responsive dashboards",
				}),
			},
		},
		{
			key:         "flutter-ecommerce",
			categoryKey: "mobile-development",
			mentor:      mentor2,
			course: domain.Course{
				Title:       "Flutter E-Commerce App dari Nol",
				Description: "Membuat aplikasi e-commerce lengkap dengan Flutter. Mulai dari UI design, state management dengan BLoC, integrasi API, hingga payment gateway.",
				Thumbnail:   "https://images.unsplash.com/photo-1512941937669-90a1b58e7e9c?w=800",
				Visibility:  "public",
				LearningObjectives: mustJSON([]string{
					"Setting up Flutter development environment",
					"Building beautiful UI with Flutter widgets",
					"State management with BLoC pattern",
					"Integrating REST API & payment gateway",
					"Implementing authentication & user management",
					"Publishing app to Play Store & App Store",
				}),
			},
		},
		{
			key:         "fullstack-nextjs",
			categoryKey: "web-development",
			mentor:      mentor2,
			course: domain.Course{
				Title:       "Fullstack Web App dengan Next.js 14",
				Description: "Panduan lengkap membangun aplikasi web fullstack dengan Next.js 14, Server Components, Prisma ORM, dan Vercel deployment.",
				Thumbnail:   "https://images.unsplash.com/photo-1460925895917-afdab827c52f?w=800",
				Visibility:  "public",
				LearningObjectives: mustJSON([]string{
					"Understanding Next.js 14 App Router & Server Components",
					"Building API routes & database models with Prisma",
					"Implementing authentication with NextAuth.js",
					"Building server-rendered & static pages",
					"Deploying to Vercel with CI/CD",
					"Optimizing performance & SEO",
				}),
			},
		},
		{
			key:         "docker-kubernetes",
			categoryKey: "devops-cloud",
			mentor:      mentor1,
			course: domain.Course{
				Title:       "Docker & Kubernetes untuk Developer",
				Description: "Menguasai containerization dan orchestration. Dari Docker basics hingga Kubernetes cluster management di production.",
				Thumbnail:   "https://images.unsplash.com/photo-1667372393119-3d4c48d07fc9?w=800",
				Visibility:  "public",
				LearningObjectives: mustJSON([]string{
					"Understanding containerization concepts",
					"Building & optimizing Docker images",
					"Orchestrating containers with Kubernetes",
					"Managing deployments, services & ingress",
					"Implementing CI/CD pipelines",
					"Monitoring & logging in production",
				}),
			},
		},
	}

	result := make(map[string]*domain.Course)

	for _, cd := range coursesData {
		course := cd.course
		course.AdminID = admin.ID
		course.MentorID = cd.mentor.ID
		course.CategoryID = categories[cd.categoryKey].ID

		if err := db.Where("title = ? AND admin_id = ?", course.Title, course.AdminID).FirstOrCreate(&course).Error; err != nil {
			return nil, err
		}

		// Force update mentor_id in case course existed before mentor support was added
		if cd.mentor != nil {
			if err := db.Model(&domain.Course{}).Where("id = ?", course.ID).Update("mentor_id", cd.mentor.ID).Error; err != nil {
				fmt.Printf("    ⚠️  Failed to update mentor_id for %s: %v\n", course.Title, err)
			}
		}

		result[cd.key] = &course
		fmt.Printf("    ✓ Course: %s (ID: %s, Mentor: %s)\n", course.Title, course.ID, cd.mentor.Name)
	}

	return result, nil
}

// ============================================================================
// 5. SEED SECTIONS
// ============================================================================

func seedSections(db *gorm.DB, courses map[string]*domain.Course) (map[string]*domain.Section, error) {
	sectionsData := []struct {
		key       string
		courseKey string
		section   domain.Section
	}{
		{
			key:       "go-intro-section",
			courseKey: "go-rest-api",
			section: domain.Section{
				Title:       "Introduction to Go",
				Description: "Basic introduction to Go programming language",
				OrderIndex:  1,
			},
		},
		{
			key:       "go-fiber-section",
			courseKey: "go-rest-api",
			section: domain.Section{
				Title:       "Building REST API with Fiber",
				Description: "Learn Fiber framework and REST API development",
				OrderIndex:  2,
			},
		},
		{
			key:       "react-setup-section",
			courseKey: "react-dashboard",
			section: domain.Section{
				Title:       "React Fundamentals",
				Description: "Learn React basics and setup",
				OrderIndex:  1,
			},
		},
		{
			key:       "react-components-section",
			courseKey: "react-dashboard",
			section: domain.Section{
				Title:       "Dashboard Components",
				Description: "Building reusable dashboard UI components",
				OrderIndex:  2,
			},
		},
		{
			key:       "flutter-section",
			courseKey: "flutter-ecommerce",
			section: domain.Section{
				Title:       "Flutter Fundamentals",
				Description: "Introduction to Flutter and Dart",
				OrderIndex:  1,
			},
		},
	}

	result := make(map[string]*domain.Section)

	for _, sd := range sectionsData {
		section := sd.section
		section.CourseID = courses[sd.courseKey].ID
		if err := db.Where("title = ? AND course_id = ?", section.Title, section.CourseID).FirstOrCreate(&section).Error; err != nil {
			return nil, err
		}
		result[sd.key] = &section
		fmt.Printf("    ✓ Section: %s (ID: %s, Course: %s)\n", section.Title, section.ID, sd.courseKey)
	}

	return result, nil
}

// ============================================================================
// 6. SEED LESSONS + LESSON DETAILS
// ============================================================================

func seedLessons(db *gorm.DB, admin *domain.User, courses map[string]*domain.Course, sections map[string]*domain.Section) (map[string]*domain.Lesson, error) {
	lessonsData := []struct {
		key        string
		courseKey  string
		sectionKey string
		lesson     domain.Lesson
		detail     domain.LessonDetail
	}{
		{
			key:        "go-intro",
			courseKey:  "go-rest-api",
			sectionKey: "go-intro-section",
			lesson: domain.Lesson{
				Title:       "Pengenalan Go & Setup Environment",
				Slug:        "pengenalan-go-setup-environment",
				Description: "Mengenal bahasa Go, instalasi tools, dan setup project pertama",
				Thumbnail:   "https://images.unsplash.com/photo-1515879218367-8466d910auj8?w=800",
				VideoURL:    "https://www.youtube.com/watch?v=Q0sJQtgj1gY",
				Difficulty:  "beginner",
				Duration:    45,
				OrderIndex:  1,
				SequenceNum: 1,
				IsFirst:     true,
				Visibility:  "public",
			},
			detail: domain.LessonDetail{
				About: "Di kelas ini kamu akan belajar dasar-dasar bahasa Go (Golang), mulai dari instalasi, konfigurasi environment, hingga membuat program pertama.",
				Rules: "1. Pastikan sudah menginstall Go versi terbaru\n2. Gunakan code editor (VS Code direkomendasikan)\n3. Selesaikan semua latihan sebelum lanjut ke kelas berikutnya",
				Tools: mustJSON([]string{"Go 1.21+", "VS Code", "Git", "Terminal/Command Prompt"}),
				ResourceMedia: mustJSON(map[string][]string{
					"videos":    {"https://youtube.com/watch?v=example1"},
					"documents": {"https://go.dev/doc/tutorial/getting-started"},
					"images":    {"https://go.dev/images/go-logo-blue.svg"},
				}),
				Resources: mustJSON([]map[string]string{
					{"type": "link", "title": "Dokumentasi Resmi Go", "url": "https://go.dev/doc/"},
					{"type": "link", "title": "Go Playground", "url": "https://go.dev/play/"},
					{"type": "pdf", "title": "Go Cheatsheet", "url": "https://example.com/go-cheatsheet.pdf"},
				}),
			},
		},
		{
			key:        "go-fiber-basics",
			courseKey:  "go-rest-api",
			sectionKey: "go-fiber-section",
			lesson: domain.Lesson{
				Title:       "Fiber Framework & Routing",
				Slug:        "fiber-framework-routing",
				Description: "Belajar Fiber framework, routing, middleware, dan request handling",
				Thumbnail:   "https://images.unsplash.com/photo-1516116216624-53e697fedbea?w=800",
				VideoURL:    "https://www.youtube.com/watch?v=9qB-KfOI31E",
				Difficulty:  "beginner",
				Duration:    60,
				OrderIndex:  2,
				SequenceNum: 2,
				IsFirst:     false,
				Visibility:  "public",
			},
			detail: domain.LessonDetail{
				About: "Fiber adalah web framework Go yang terinspirasi dari Express.js. Setup Fiber, membuat route, menggunakan middleware, dan handling request/response.",
				Rules: "1. Sudah menyelesaikan kelas Pengenalan Go\n2. Praktikkan setiap contoh kode\n3. Buat minimal 5 endpoint berbeda sebagai latihan",
				Tools: mustJSON([]string{"Go 1.21+", "Fiber v2", "Postman/Insomnia", "VS Code"}),
				ResourceMedia: mustJSON(map[string][]string{
					"videos":    {"https://youtube.com/watch?v=example2"},
					"documents": {"https://docs.gofiber.io/"},
					"images":    {},
				}),
				Resources: mustJSON([]map[string]string{
					{"type": "link", "title": "Fiber Documentation", "url": "https://docs.gofiber.io/"},
					{"type": "link", "title": "Fiber GitHub", "url": "https://github.com/gofiber/fiber"},
				}),
			},
		},
		{
			key:        "go-gorm-db",
			courseKey:  "go-rest-api",
			sectionKey: "go-fiber-section",
			lesson: domain.Lesson{
				Title:       "GORM & Database Integration",
				Slug:        "gorm-database-integration",
				Description: "Menghubungkan Go dengan PostgreSQL menggunakan GORM ORM",
				Thumbnail:   "https://images.unsplash.com/photo-1544383835-bda2bc66a55d?w=800",
				VideoURL:    "https://www.youtube.com/watch?v=Z0h0yRKKjQg",
				Difficulty:  "intermediate",
				Duration:    90,
				OrderIndex:  3,
				SequenceNum: 3,
				IsFirst:     false,
				Visibility:  "public",
			},
			detail: domain.LessonDetail{
				About: "GORM adalah ORM paling populer untuk Go. Koneksi database, model definition, migration, CRUD operations, relationships, dan query optimization.",
				Rules: "1. Install PostgreSQL di local atau gunakan Docker\n2. Sudah menyelesaikan kelas Fiber Framework\n3. Pahami konsep dasar SQL sebelum memulai",
				Tools: mustJSON([]string{"Go 1.21+", "PostgreSQL 15+", "GORM v2", "pgAdmin/DBeaver", "Docker (optional)"}),
				ResourceMedia: mustJSON(map[string][]string{
					"videos":    {"https://youtube.com/watch?v=example3"},
					"documents": {"https://gorm.io/docs/"},
					"images":    {},
				}),
				Resources: mustJSON([]map[string]string{
					{"type": "link", "title": "GORM Documentation", "url": "https://gorm.io/docs/"},
					{"type": "link", "title": "PostgreSQL Tutorial", "url": "https://www.postgresql.org/docs/current/tutorial.html"},
				}),
			},
		},
		{
			key:        "react-setup",
			courseKey:  "react-dashboard",
			sectionKey: "react-setup-section",
			lesson: domain.Lesson{
				Title:       "Setup React + TypeScript Project",
				Slug:        "setup-react-typescript-project",
				Description: "Setup project React dengan Vite, TypeScript, dan Tailwind CSS",
				Thumbnail:   "https://images.unsplash.com/photo-1633356122102-3fe601e05bd2?w=800",
				VideoURL:    "https://www.youtube.com/watch?v=8a9hsZPmXms",
				Difficulty:  "beginner",
				Duration:    30,
				OrderIndex:  1,
				SequenceNum: 1,
				IsFirst:     true,
				Visibility:  "public",
			},
			detail: domain.LessonDetail{
				About: "Memulai project React modern menggunakan Vite sebagai build tool, TypeScript untuk type safety, dan Tailwind CSS untuk styling yang cepat.",
				Rules: "1. Install Node.js versi 18+\n2. Familiar dengan HTML, CSS, JavaScript\n3. Gunakan npm atau yarn sebagai package manager",
				Tools: mustJSON([]string{"Node.js 18+", "VS Code", "Chrome DevTools", "npm/yarn"}),
				ResourceMedia: mustJSON(map[string][]string{
					"videos":    {"https://youtube.com/watch?v=example4"},
					"documents": {"https://react.dev/learn"},
					"images":    {},
				}),
				Resources: mustJSON([]map[string]string{
					{"type": "link", "title": "React Documentation", "url": "https://react.dev/"},
					{"type": "link", "title": "Vite Guide", "url": "https://vitejs.dev/guide/"},
				}),
			},
		},
		{
			key:        "react-components",
			courseKey:  "react-dashboard",
			sectionKey: "react-components-section",
			lesson: domain.Lesson{
				Title:       "Membuat Komponen Dashboard",
				Slug:        "membuat-komponen-dashboard",
				Description: "Membangun komponen-komponen UI dashboard: Sidebar, Navbar, Cards, Charts",
				Thumbnail:   "https://images.unsplash.com/photo-1551288049-bebda4e38f71?w=800",
				VideoURL:    "https://www.youtube.com/watch?v=6aBd6Zw3k8A",
				Difficulty:  "intermediate",
				Duration:    75,
				OrderIndex:  2,
				SequenceNum: 2,
				IsFirst:     false,
				Visibility:  "public",
			},
			detail: domain.LessonDetail{
				About: "Membuat komponen-komponen UI reusable untuk dashboard. Layout system, sidebar navigation, stats cards, hingga data table dengan sorting dan filtering.",
				Rules: "1. Sudah menyelesaikan kelas Setup React\n2. Jangan copy-paste, ketik sendiri untuk memahami\n3. Buat variasi komponen sebagai latihan",
				Tools: mustJSON([]string{"React 18+", "TypeScript", "Tailwind CSS", "Recharts/Chart.js"}),
				ResourceMedia: mustJSON(map[string][]string{
					"videos":    {"https://youtube.com/watch?v=example5"},
					"documents": {},
					"images":    {},
				}),
				Resources: mustJSON([]map[string]string{
					{"type": "link", "title": "Recharts", "url": "https://recharts.org/"},
					{"type": "link", "title": "Tailwind Components", "url": "https://tailwindui.com/components"},
				}),
			},
		},
		{
			key:        "flutter-intro",
			courseKey:  "flutter-ecommerce",
			sectionKey: "flutter-section",
			lesson: domain.Lesson{
				Title:       "Flutter Fundamentals & Dart Basics",
				Slug:        "flutter-fundamentals-dart-basics",
				Description: "Mengenal Flutter framework dan dasar-dasar bahasa Dart",
				Thumbnail:   "https://images.unsplash.com/photo-1617040619263-41c5a9ca7521?w=800",
				VideoURL:    "https://www.youtube.com/watch?v=CD1Y2xDOqC0",
				Difficulty:  "beginner",
				Duration:    60,
				OrderIndex:  1,
				SequenceNum: 1,
				IsFirst:     true,
				Visibility:  "public",
			},
			detail: domain.LessonDetail{
				About: "Pengenalan Flutter dan Dart untuk pemula. Widget tree, layout system, state management dasar, dan navigasi antar halaman.",
				Rules: "1. Install Flutter SDK & Android Studio/VS Code\n2. Setup emulator atau gunakan physical device\n3. Selesaikan Flutter doctor tanpa error",
				Tools: mustJSON([]string{"Flutter SDK 3.x", "Dart", "Android Studio / VS Code", "Android Emulator / iOS Simulator"}),
				ResourceMedia: mustJSON(map[string][]string{
					"videos":    {"https://youtube.com/watch?v=example6"},
					"documents": {"https://flutter.dev/docs"},
					"images":    {},
				}),
				Resources: mustJSON([]map[string]string{
					{"type": "link", "title": "Flutter Documentation", "url": "https://flutter.dev/docs"},
					{"type": "link", "title": "Dart Language Tour", "url": "https://dart.dev/language"},
				}),
			},
		},
	}

	result := make(map[string]*domain.Lesson)

	for _, ld := range lessonsData {
		lesson := ld.lesson
		lesson.AdminID = admin.ID
		lesson.CourseID = courses[ld.courseKey].ID
		lesson.SectionID = sections[ld.sectionKey].ID

		videoURL := lesson.VideoURL

		if err := db.Where("slug = ? AND course_id = ?", lesson.Slug, lesson.CourseID).FirstOrCreate(&lesson).Error; err != nil {
			return nil, err
		}

		// Force update video_url regardless (FirstOrCreate doesn't update existing records)
		if videoURL != "" {
			if err := db.Model(&domain.Lesson{}).Where("id = ?", lesson.ID).Update("video_url", videoURL).Error; err != nil {
				fmt.Printf("    ⚠️  Failed to update video_url for %s: %v\n", lesson.Slug, err)
			}
		}

		result[ld.key] = &lesson

		detail := ld.detail
		detail.LessonID = lesson.ID
		if err := db.Where("lesson_id = ?", detail.LessonID).FirstOrCreate(&detail).Error; err != nil {
			return nil, err
		}

		fmt.Printf("    ✓ Lesson: %s (ID: %s, Course: %s)\n", lesson.Title, lesson.ID, ld.courseKey)
	}

	linkLessonProgression(db, result, "go-intro", "go-fiber-basics")
	linkLessonProgression(db, result, "go-fiber-basics", "go-gorm-db")
	linkLessonProgression(db, result, "react-setup", "react-components")

	return result, nil
}

func linkLessonProgression(db *gorm.DB, lessons map[string]*domain.Lesson, fromKey, toKey string) {
	from := lessons[fromKey]
	to := lessons[toKey]
	if from != nil && to != nil {
		db.Model(&domain.Lesson{}).Where("id = ?", from.ID).Update("next_lesson_id", to.ID)
	}
}

// ============================================================================
// 6. SEED LESSON PROGRESS
// ============================================================================

func seedClassProgress(db *gorm.DB, users map[string]*domain.User, lessons map[string]*domain.Lesson) error {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	lastWeek := now.Add(-7 * 24 * time.Hour)

	progressData := []domain.LessonProgress{
		{
			UserID:             users["budi"].ID,
			LessonID:           lessons["go-intro"].ID,
			Status:             "completed",
			StartedAt:          &lastWeek,
			CompletedAt:        &yesterday,
			ProgressPercentage: 100,
			Notes:              "Kelas yang sangat bagus untuk pemula Go!",
		},
		{
			UserID:             users["budi"].ID,
			LessonID:           lessons["go-fiber-basics"].ID,
			Status:             "in_progress",
			StartedAt:          &yesterday,
			ProgressPercentage: 45,
			Notes:              "Sedang belajar middleware",
		},
		{
			UserID:             users["siti"].ID,
			LessonID:           lessons["react-setup"].ID,
			Status:             "started",
			StartedAt:          &now,
			ProgressPercentage: 10,
		},
		{
			UserID:             users["andi"].ID,
			LessonID:           lessons["go-intro"].ID,
			Status:             "completed",
			StartedAt:          &lastWeek,
			CompletedAt:        &lastWeek,
			ProgressPercentage: 100,
		},
		{
			UserID:             users["andi"].ID,
			LessonID:           lessons["go-fiber-basics"].ID,
			Status:             "completed",
			StartedAt:          &lastWeek,
			CompletedAt:        &yesterday,
			ProgressPercentage: 100,
		},
		{
			UserID:             users["andi"].ID,
			LessonID:           lessons["go-gorm-db"].ID,
			Status:             "in_progress",
			StartedAt:          &now,
			ProgressPercentage: 30,
			Notes:              "Relationships di GORM lumayan tricky",
		},
		{
			UserID:             users["dewi"].ID,
			LessonID:           lessons["flutter-intro"].ID,
			Status:             "started",
			StartedAt:          &now,
			ProgressPercentage: 5,
		},
	}

	for _, progress := range progressData {
		p := progress
		if err := db.Where("user_id = ? AND lesson_id = ?", p.UserID, p.LessonID).FirstOrCreate(&p).Error; err != nil {
			return err
		}
	}
	fmt.Printf("    ✓ Created %d lesson progress records\n", len(progressData))

	return nil
}

// ============================================================================
// 7. SEED DISCUSSIONS & REPLIES
// ============================================================================

func seedDiscussions(db *gorm.DB, users map[string]*domain.User, lessons map[string]*domain.Lesson, categories map[string]*domain.Category) error {
	lessonGoIntroID := lessons["go-intro"].ID
	lessonReactSetupID := lessons["react-setup"].ID

	// Discussion 1: Budi bertanya tentang Go
	disc1 := domain.Discussion{
		Title:      "Perbedaan goroutine dan thread biasa?",
		Content:    "Halo semua, saya baru belajar Go dan masih bingung bedanya goroutine dengan thread biasa di bahasa lain. Bisa tolong jelaskan dengan contoh sederhana?",
		UserID:     users["budi"].ID,
		LessonID:   &lessonGoIntroID,
		CategoryID: categories["backend-development"].ID,
		ViewsCount: 42,
		Status:     "open",
	}
	if err := db.Where("title = ? AND user_id = ?", disc1.Title, disc1.UserID).FirstOrCreate(&disc1).Error; err != nil {
		return err
	}

	// Replies for Discussion 1
	reply1 := domain.Reply{
		Content:      "Goroutine itu lightweight thread yang dikelola oleh Go runtime, bukan OS. Satu OS thread bisa menjalankan ribuan goroutine karena Go punya scheduler sendiri (M:N scheduling).",
		UserID:       users["andi"].ID,
		DiscussionID: disc1.ID,
		LikesCount:   5,
		IsMarkedBest: true,
	}
	if err := db.Where("discussion_id = ? AND user_id = ? AND is_marked_best = ?", reply1.DiscussionID, reply1.UserID, true).FirstOrCreate(&reply1).Error; err != nil {
		return err
	}

	reply2 := domain.Reply{
		Content:      "Tambahannya, goroutine juga lebih murah untuk context switching karena Go scheduler yang handle, bukan kernel. Makanya bisa jalankan puluhan ribu goroutine tanpa masalah.",
		UserID:       users["admin"].ID,
		DiscussionID: disc1.ID,
		LikesCount:   3,
	}
	if err := db.Where("discussion_id = ? AND user_id = ? AND content LIKE ?", reply2.DiscussionID, reply2.UserID, "%context switching%").FirstOrCreate(&reply2).Error; err != nil {
		return err
	}

	// Nested reply
	nestedReply := domain.Reply{
		Content:      "Terima kasih penjelasannya! Sekarang sudah lebih paham. Jadi intinya goroutine itu lebih efisien dari thread ya.",
		UserID:       users["budi"].ID,
		DiscussionID: disc1.ID,
		ParentID:     &reply1.ID,
	}
	if err := db.Where("discussion_id = ? AND user_id = ? AND parent_id = ?", nestedReply.DiscussionID, nestedReply.UserID, nestedReply.ParentID).FirstOrCreate(&nestedReply).Error; err != nil {
		return err
	}

	// Discussion 2: Siti bertanya tentang React
	disc2 := domain.Discussion{
		Title:      "Kapan pakai useState vs useReducer?",
		Content:    "Saya lagi belajar React hooks dan bingung kapan harus pakai useState dan kapan pakai useReducer. Ada rule of thumb-nya gak?",
		UserID:     users["siti"].ID,
		LessonID:   &lessonReactSetupID,
		CategoryID: categories["frontend-development"].ID,
		ViewsCount: 28,
		Status:     "open",
	}
	if err := db.Where("title = ? AND user_id = ?", disc2.Title, disc2.UserID).FirstOrCreate(&disc2).Error; err != nil {
		return err
	}

	reply3 := domain.Reply{
		Content:      "Rule of thumb: useState untuk state sederhana, useReducer untuk state kompleks dengan banyak field yang saling bergantung.",
		UserID:       users["andi"].ID,
		DiscussionID: disc2.ID,
		LikesCount:   4,
	}
	if err := db.Where("discussion_id = ? AND user_id = ?", reply3.DiscussionID, reply3.UserID).FirstOrCreate(&reply3).Error; err != nil {
		return err
	}

	// Discussion 3: General discussion
	disc3 := domain.Discussion{
		Title:      "Tips belajar programming untuk pemula",
		Content:    "Halo teman-teman! Share tips belajar programming: konsisten, langsung praktik, buat project sendiri, gabung komunitas, jangan takut error.",
		UserID:     users["admin"].ID,
		CategoryID: categories["web-development"].ID,
		ViewsCount: 156,
		Status:     "open",
		IsPinned:   true,
	}
	if err := db.Where("title = ? AND user_id = ?", disc3.Title, disc3.UserID).FirstOrCreate(&disc3).Error; err != nil {
		return err
	}

	fmt.Printf("    ✓ Created 3 discussions with replies\n")
	return nil
}

// ============================================================================
// 8. SEED SHOWCASES
// ============================================================================

func seedShowcases(db *gorm.DB, users map[string]*domain.User, categories map[string]*domain.Category) error {
	showcasesData := []struct {
		key     string
		userKey string
		catKey  string
		s       domain.Showcase
	}{
		{
			key:     "portfolio-budi",
			userKey: "budi",
			catKey:  "web-development",
			s: domain.Showcase{
				Title:       "Portfolio Website Personal",
				Description: "Website portfolio menggunakan Next.js 14 dan Tailwind CSS. Fitur: dark mode, animasi smooth, blog dengan MDX, dan SEO optimized.",
				MediaURLs:   mustJSON([]string{"https://images.unsplash.com/photo-1517180102446-f3ece451e9d8?w=800", "https://images.unsplash.com/photo-1460925895917-afdab827c52f?w=800"}),
				Status:      "published",
				Visibility:  "public",
				LikesCount:  12,
				ViewsCount:  89,
			},
		},
		{
			key:     "microservices-andi",
			userKey: "andi",
			catKey:  "backend-development",
			s: domain.Showcase{
				Title:       "REST API Microservices Go",
				Description: "Arsitektur microservices menggunakan Go, gRPC, dan RabbitMQ. Termasuk API gateway, service discovery, dan distributed tracing.",
				MediaURLs:   mustJSON([]string{"https://images.unsplash.com/photo-1558494949-ef010cbdcc31?w=800"}),
				Status:      "published",
				Visibility:  "public",
				LikesCount:  23,
				ViewsCount:  145,
			},
		},
		{
			key:     "uikit-siti",
			userKey: "siti",
			catKey:  "ui-ux-design",
			s: domain.Showcase{
				Title:       "Mobile App UI Kit Design",
				Description: "Koleksi UI Kit untuk aplikasi mobile modern. 50+ screen design untuk e-commerce, social media, dan productivity app.",
				MediaURLs:   mustJSON([]string{"https://images.unsplash.com/photo-1512941937669-90a1b58e7e9c?w=800"}),
				Status:      "published",
				Visibility:  "public",
				LikesCount:  18,
				ViewsCount:  112,
			},
		},
		{
			key:     "todo-flutter-dewi",
			userKey: "dewi",
			catKey:  "mobile-development",
			s: domain.Showcase{
				Title:       "Aplikasi Todo Sederhana Flutter",
				Description: "Aplikasi todo list pertama dengan Flutter. Fitur CRUD, local storage dengan Hive, dan UI yang clean.",
				MediaURLs:   mustJSON([]string{"https://images.unsplash.com/photo-1617040619263-41c5a9ca7521?w=800"}),
				Status:      "published",
				Visibility:  "public",
				LikesCount:  7,
				ViewsCount:  34,
			},
		},
	}

	// Create showcases
	showcaseMap := make(map[string]*domain.Showcase)
	for _, sd := range showcasesData {
		showcase := sd.s
		showcase.UserID = users[sd.userKey].ID
		showcase.CategoryID = categories[sd.catKey].ID

		if err := db.Where("title = ? AND user_id = ?", showcase.Title, showcase.UserID).FirstOrCreate(&showcase).Error; err != nil {
			return err
		}
		showcaseMap[sd.key] = &showcase
		fmt.Printf("    ✓ Showcase: %s (ID: %s)\n", showcase.Title, showcase.ID)
	}

	// Seed showcase likes using actual CUID IDs
	likes := []domain.ShowcaseLike{
		{UserID: users["andi"].ID, ShowcaseID: showcaseMap["portfolio-budi"].ID},
		{UserID: users["siti"].ID, ShowcaseID: showcaseMap["portfolio-budi"].ID},
		{UserID: users["admin"].ID, ShowcaseID: showcaseMap["microservices-andi"].ID},
		{UserID: users["budi"].ID, ShowcaseID: showcaseMap["microservices-andi"].ID},
		{UserID: users["dewi"].ID, ShowcaseID: showcaseMap["uikit-siti"].ID},
		{UserID: users["budi"].ID, ShowcaseID: showcaseMap["uikit-siti"].ID},
		{UserID: users["andi"].ID, ShowcaseID: showcaseMap["todo-flutter-dewi"].ID},
	}

	for _, like := range likes {
		l := like
		db.Where("user_id = ? AND showcase_id = ?", l.UserID, l.ShowcaseID).FirstOrCreate(&l)
	}

	// Seed showcase comments
	comments := []domain.ShowcaseComment{
		{
			Content:    "Keren banget portfolionya! Animasinya smooth.",
			UserID:     users["andi"].ID,
			ShowcaseID: showcaseMap["portfolio-budi"].ID,
		},
		{
			Content:    "Arsitektur microservicesnya rapi, boleh share repo-nya?",
			UserID:     users["budi"].ID,
			ShowcaseID: showcaseMap["microservices-andi"].ID,
		},
		{
			Content:    "UI Kit-nya cantik banget! Warna-warnanya harmonious.",
			UserID:     users["dewi"].ID,
			ShowcaseID: showcaseMap["uikit-siti"].ID,
		},
	}

	for _, comment := range comments {
		c := comment
		db.Where("content = ? AND user_id = ? AND showcase_id = ?", c.Content, c.UserID, c.ShowcaseID).FirstOrCreate(&c)
	}

	fmt.Printf("    ✓ Created showcase likes & comments\n")
	return nil
}

// ============================================================================
// 9. SEED REVIEWS
// ============================================================================

func seedStudyCases(db *gorm.DB, users map[string]*domain.User) (map[string]*domain.StudyCase, error) {
	studyCasesData := []struct {
		key       string
		creator   string
		studyCase domain.StudyCase
	}{
		{
			key:     "go-microservice",
			creator: "admin",
			studyCase: domain.StudyCase{
				Name:        "Microservices E-Commerce dengan Go",
				Description: "Studi kasus implementasi microservices untuk platform e-commerce menggunakan Go, gRPC, RabbitMQ, dan PostgreSQL. Membahas service discovery, distributed tracing, dan API gateway pattern.",
				ImgURL:      "https://images.unsplash.com/photo-1558494949-ef010cbdcc31?w=800",
				YoutubeURL:  "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
				Tags:        mustJSON([]string{"Go", "Microservices", "gRPC", "RabbitMQ"}),
			},
		},
		{
			key:     "react-ecommerce",
			creator: "admin",
			studyCase: domain.StudyCase{
				Name:        "Fullstack E-Commerce dengan React & Node.js",
				Description: "Studi kasus membangun platform e-commerce end-to-end dengan React, TypeScript, Node.js, Express, dan MongoDB. Meliputi payment integration, cart management, dan admin dashboard.",
				ImgURL:      "https://images.unsplash.com/photo-1556742049-0cfed4f6a45d?w=800",
				YoutubeURL:  "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
				Tags:        mustJSON([]string{"React", "TypeScript", "Node.js", "MongoDB"}),
			},
		},
	}

	result := make(map[string]*domain.StudyCase)

	for _, sd := range studyCasesData {
		studyCase := sd.studyCase
		studyCase.UserID = users[sd.creator].ID

		if err := db.Where("name = ?", studyCase.Name).FirstOrCreate(&studyCase).Error; err != nil {
			return nil, err
		}
		result[sd.key] = &studyCase
		fmt.Printf("    ✓ Study Case: %s (ID: %s)\n", studyCase.Name, studyCase.ID)
	}

	return result, nil
}

func seedStudyCaseDiscussions(db *gorm.DB, users map[string]*domain.User, studyCases map[string]*domain.StudyCase, categories map[string]*domain.Category) error {
	disc1 := domain.Discussion{
		Title:       "Best practices microservices dengan Go?",
		Content:     "Saya baru selesai mempelajari studi kasus microservices e-commerce. Ada yang punya pengalaman production-ready? Share dong best practices-nya!",
		UserID:      users["budi"].ID,
		StudyCaseID: &studyCases["go-microservice"].ID,
		CategoryID:  categories["backend-development"].ID,
		ViewsCount:  23,
		Status:      "open",
	}
	if err := db.Where("title = ? AND study_case_id = ?", disc1.Title, disc1.StudyCaseID).FirstOrCreate(&disc1).Error; err != nil {
		return err
	}

	disc2 := domain.Discussion{
		Title:       "Saran payment gateway untuk e-commerce React?",
		Content:     "Studi kasus e-commerce-nya keren! Ada saran payment gateway yang mudah diintegrasikan untuk skala kecil menengah?",
		UserID:      users["siti"].ID,
		StudyCaseID: &studyCases["react-ecommerce"].ID,
		CategoryID:  categories["frontend-development"].ID,
		ViewsCount:  15,
		Status:      "open",
	}
	if err := db.Where("title = ? AND study_case_id = ?", disc2.Title, disc2.StudyCaseID).FirstOrCreate(&disc2).Error; err != nil {
		return err
	}

	fmt.Printf("    ✓ Created 2 study case discussions\n")
	return nil
}

func seedReviews(db *gorm.DB, users map[string]*domain.User, courses map[string]*domain.Course, lessons map[string]*domain.Lesson) error {
	courseTitle := make(map[string]string)
	for key, c := range courses {
		courseTitle[c.ID] = key
	}
	lessonTitle := make(map[string]string)
	for key, l := range lessons {
		lessonTitle[l.ID] = key
	}
	userName := make(map[string]string)
	for key, u := range users {
		userName[u.ID] = key
	}

	reviewsData := []domain.Review{
		{
			UserID:   users["budi"].ID,
			CourseID: courses["go-rest-api"].ID,
			LessonID: lessons["go-intro"].ID,
			Rating:   5,
			Message:  "Kelas pengenalan Go yang sangat jelas! Cocok untuk pemula seperti saya.",
		},
		{
			UserID:   users["budi"].ID,
			CourseID: courses["go-rest-api"].ID,
			Rating:   4,
			Message:  "Kursus Go yang komprehensif. Penjelasan dari dasar sampai deployment. Recommended!",
		},
		{
			UserID:   users["andi"].ID,
			CourseID: courses["go-rest-api"].ID,
			LessonID: lessons["go-fiber-basics"].ID,
			Rating:   5,
			Message:  "Materi Fiber dan routing-nya lengkap banget, langsung bisa dipraktikkan.",
		},
		{
			UserID:   users["andi"].ID,
			CourseID: courses["go-rest-api"].ID,
			Rating:   5,
			Message:  "Best Go course ever! Materi GORM-nya bikin paham database integration.",
		},
		{
			UserID:   users["siti"].ID,
			CourseID: courses["react-dashboard"].ID,
			LessonID: lessons["react-setup"].ID,
			Rating:   4,
			Message:  "Setup React dengan Vite + TypeScript dijelaskan step by step. Mudah diikuti.",
		},
		{
			UserID:   users["siti"].ID,
			CourseID: courses["react-dashboard"].ID,
			Rating:   4,
			Message:  "Dashboard React-nya keren! Banyak komponen UI yang reusable.",
		},
		{
			UserID:   users["dewi"].ID,
			CourseID: courses["flutter-ecommerce"].ID,
			LessonID: lessons["flutter-intro"].ID,
			Rating:   5,
			Message:  "Senang banget akhirnya paham Flutter & Dart. Tutorialnya ramah pemula.",
		},
	}

	for _, review := range reviewsData {
		r := review
		query := db.Where("user_id = ?", r.UserID)
		if r.LessonID != "" {
			query = query.Where("lesson_id = ?", r.LessonID)
		} else {
			query = query.Where("course_id = ?", r.CourseID)
		}
		if err := query.FirstOrCreate(&r).Error; err != nil {
			return err
		}
		fmt.Printf("    ✓ Review: %s → %s (Rating: %d)\n", userName[r.UserID], func() string {
			if r.LessonID != "" {
				return "Lesson: " + lessonTitle[r.LessonID]
			}
			return "Course: " + courseTitle[r.CourseID]
		}(), r.Rating)
	}

	return nil
}

// ============================================================================
// UTILITY HELPERS
// ============================================================================

// mustJSON marshals value to datatypes.JSON, panics on error
func mustJSON(v interface{}) datatypes.JSON {
	data, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal JSON: %v", err))
	}
	return datatypes.JSON(data)
}

// ============================================================================
// 12. SEED BLOGS
// ============================================================================
func seedBlogs(db *gorm.DB, users map[string]*domain.User, categories map[string]*domain.Category) error {
	type blogData struct {
		key     string
		userKey string
		catKey  string
		blog    domain.Blog
	}

	blogsData := []blogData{
		{
			key:     "blog1",
			userKey: "admin",
			catKey:  "backend-development",
			blog: domain.Blog{
				Title:       "Getting Started with Go Programming",
				Slug:        "getting-started-with-go-programming",
				Description: "A comprehensive guide to start learning Go programming language from scratch.",
				Content:     "<h2>Introduction to Go</h2><p>Go is a statically typed, compiled programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson. It is syntactically similar to C, but with memory safety, garbage collection, structural typing, and CSP-style concurrency.</p><h2>Why Learn Go?</h2><p>Go is known for its simplicity, efficiency, and built-in concurrency features. It's widely used for building web servers, cloud services, and DevOps tools.</p><h2>Getting Started</h2><p>To get started with Go, you need to install the Go compiler and set up your development environment. The official documentation provides excellent resources for beginners.</p>",
				CoverImgURL: "https://images.unsplash.com/photo-1629654297299-c8506221ca97?w=800",
				Status:      "published",
			},
		},
		{
			key:     "blog2",
			userKey: "mentor1",
			catKey:  "frontend-development",
			blog: domain.Blog{
				Title:       "Modern React Patterns in 2024",
				Slug:        "modern-react-patterns-2024",
				Description: "Explore the latest React patterns and best practices for building modern web applications.",
				Content:     "<h2>React Modern Patterns</h2><p>React has evolved significantly over the years. With the introduction of hooks, server components, and new patterns, building React applications has become more efficient and enjoyable.</p><h2>Key Patterns</h2><p>Some of the most important patterns include custom hooks, compound components, render props, and the context API. Understanding these patterns will help you write more maintainable and reusable code.</p>",
				CoverImgURL: "https://images.unsplash.com/photo-1633356122102-3fe601e05bd2?w=800",
				Status:      "published",
			},
		},
		{
			key:     "blog3",
			userKey: "mentor2",
			catKey:  "devops-cloud",
			blog: domain.Blog{
				Title:       "Docker Containerization Best Practices",
				Slug:        "docker-containerization-best-practices",
				Description: "Learn the best practices for containerizing your applications with Docker.",
				Content:     "<h2>Docker Best Practices</h2><p>Docker has revolutionized the way we deploy and manage applications. However, to get the most out of Docker, you need to follow best practices for building efficient and secure containers.</p><h2>Key Practices</h2><p>Use multi-stage builds to reduce image size, implement health checks, use .dockerignore files, and avoid running containers as root. These practices will help you build production-ready containers.</p>",
				CoverImgURL: "https://images.unsplash.com/photo-1605745341112-85968b19335b?w=800",
				Status:      "published",
			},
		},
		{
			key:     "blog4",
			userKey: "admin",
			catKey:  "mobile-development",
			blog: domain.Blog{
				Title:       "Building Cross-Platform Apps with Flutter",
				Slug:        "building-cross-platform-apps-flutter",
				Description: "Discover how to build beautiful cross-platform mobile applications using Flutter.",
				Content:     "<h2>Flutter Framework</h2><p>Flutter is Google's UI toolkit for building natively compiled applications for mobile, web, and desktop from a single codebase. It uses the Dart programming language and provides a rich set of pre-built widgets.</p><h2>Getting Started with Flutter</h2><p>Flutter's hot reload feature makes development fast and enjoyable. You can see changes instantly without losing the state of your application.</p>",
				CoverImgURL: "https://images.unsplash.com/photo-1551650975-87deedd944c3?w=800",
				Status:      "published",
			},
		},
		{
			key:     "blog5",
			userKey: "mentor1",
			catKey:  "backend-development",
			blog: domain.Blog{
				Title:       "RESTful API Design Guidelines",
				Slug:        "restful-api-design-guidelines",
				Description: "Comprehensive guidelines for designing RESTful APIs that scale.",
				Content:     "<h2>API Design Principles</h2><p>Designing a good RESTful API requires careful planning and adherence to standards. This guide covers the essential principles of API design.</p><h2>Best Practices</h2><p>Use proper HTTP methods, implement consistent error handling, version your APIs, and provide comprehensive documentation. These practices ensure your API is easy to use and maintain.</p>",
				CoverImgURL: "https://images.unsplash.com/photo-1516259762381-22954d7d3ad2?w=800",
				Status:      "published",
			},
		},
		{
			key:     "blog6",
			userKey: "budi",
			catKey:  "frontend-development",
			blog: domain.Blog{
				Title:       "CSS Grid vs Flexbox: When to Use What",
				Slug:        "css-grid-vs-flexbox-when-to-use",
				Description: "Understanding the differences between CSS Grid and Flexbox for layout.",
				Content:     "<h2>CSS Layout Modules</h2><p>CSS Grid and Flexbox are two powerful layout modules in CSS. While they can sometimes be used interchangeably, each has its strengths and ideal use cases.</p><h2>When to Use Flexbox</h2><p>Flexbox is best for one-dimensional layouts - either a row OR a column. It's great for navigation bars, centering content, and distributing space among items in a container.</p><h2>When to Use Grid</h2><p>Grid is designed for two-dimensional layouts - rows AND columns simultaneously. Use it for complex page layouts, image galleries, and dashboard designs.</p>",
				CoverImgURL: "https://images.unsplash.com/photo-1507721999472-8ed4421c4af2?w=800",
				Status:      "published",
			},
		},
		{
			key:     "blog7",
			userKey: "admin",
			catKey:  "devops-cloud",
			blog: domain.Blog{
				Title:       "Kubernetes for Beginners",
				Slug:        "kubernetes-for-beginners",
				Description: "An introduction to Kubernetes container orchestration for beginners.",
				Content:     "<h2>What is Kubernetes?</h2><p>Kubernetes (K8s) is an open-source container orchestration platform that automates the deployment, scaling, and management of containerized applications.</p><h2>Key Concepts</h2><p>Understanding pods, services, deployments, and namespaces is essential for working with Kubernetes. This guide walks you through these concepts with practical examples.</p>",
				CoverImgURL: "https://images.unsplash.com/photo-1667372393119-3d4c48d07fc9?w=800",
				Status:      "draft",
			},
		},
	}

	for _, item := range blogsData {
		item.blog.UserID = users[item.userKey].ID
		item.blog.CategoryID = categories[item.catKey].ID

		tags := []string{"programming", "tutorial"}
		if item.key == "blog1" || item.key == "blog5" {
			tags = []string{"go", "backend", "programming"}
		} else if item.key == "blog2" || item.key == "blog6" {
			tags = []string{"react", "frontend", "javascript"}
		} else if item.key == "blog3" || item.key == "blog7" {
			tags = []string{"docker", "devops", "container"}
		} else if item.key == "blog4" {
			tags = []string{"flutter", "mobile", "dart"}
		}
		item.blog.Tags = mustJSON(tags)

		result := db.Where("title = ?", item.blog.Title).FirstOrCreate(&item.blog)
		if err := result.Error; err != nil {
			return fmt.Errorf("failed to seed blog %s: %w", item.key, err)
		}
		if result.RowsAffected > 0 {
			fmt.Printf("    ✓ Blog: %s\n", func() string {
				if len(item.blog.Title) > 40 {
					return item.blog.Title[:40]
				}
				return item.blog.Title
			}())
		}
	}

	fmt.Println("  ✓ Blogs seeded successfully!")
	return nil
}

// ============================================================================
// 13. SEED COMPANY PROFILE
// ============================================================================

func seedCompany(db *gorm.DB) error {
	var existing domain.Company
	result := db.First(&existing)
	if result.RowsAffected > 0 {
		fmt.Printf("    ✓ Company already exists: %s\n", existing.BrandName)
		return nil
	}

	company := domain.Company{
		BrandName: "JValleyVerse",
		Tagline:   "Learn, Build, Grow Together",
		Vision:    "Menjadi platform edukasi teknologi terdepan di Indonesia yang mencetak talenta digital berkualitas dan siap bersaing di era global.",
		Mission:   "Menyediakan materi pembelajaran berkualitas tinggi yang mudah diakses\nMembangun komunitas belajar yang kolaboratif dan suportif\nMenjembatani kesenjangan antara pendidikan formal dan kebutuhan industri\nMemberikan pengalaman belajar interaktif dengan gamifikasi dan sertifikasi",
		LogoURL:   "https://cdn.mohagussetiaone.my.id/jvalleyverse/logo.png",
		Domain:    "https://jvalleyverse.com",
		Email:     "hello@jvalleyverse.com",
		Facebook:  "https://facebook.com/jvalleyverse",
		Instagram: "https://instagram.com/jvalleyverse",
		Twitter:   "https://x.com/jvalleyverse",
		TikTok:    "https://tiktok.com/@jvalleyverse",
		Youtube:   "https://youtube.com/@jvalleyverse",
		LinkedIn:  "https://linkedin.com/company/jvalleyverse",
		WhatsApp:  "https://wa.me/6281234567890",
		Address:   "Jakarta, Indonesia",
		Phone:     "+62 812-3456-7890",
	}

	if err := db.Create(&company).Error; err != nil {
		return fmt.Errorf("failed to create company: %w", err)
	}

	fmt.Printf("    ✓ Company created: %s (ID: %s)\n", company.BrandName, company.ID)
	return nil
}

func seedFAQs(db *gorm.DB) error {
	var existing domain.FAQ
	if result := db.First(&existing); result.RowsAffected > 0 {
		fmt.Println("    ✓ FAQs already exist, skipping...")
		return nil
	}

	faqs := []domain.FAQ{
		{
			Question:   "Apa itu JValleyverse?",
			Answer:     "JValleyverse adalah platform belajar coding online gratis dengan sistem gamifikasi, sertifikat, dan komunitas yang saling mendukung. Kami menyediakan kursus berkualitas tinggi untuk membantu developer Indonesia berkembang.",
			Category:   "general",
			OrderIndex: 1,
			IsActive:   true,
		},
		{
			Question:   "Apakah kursus di JValleyverse benar-benar gratis?",
			Answer:     "Ya, semua course yang tersedia di JValleyverse dapat diakses secara gratis tanpa biaya. Kami ingin menciptakan akses belajar yang inklusif untuk semua kalangan.",
			Category:   "general",
			OrderIndex: 2,
			IsActive:   true,
		},
		{
			Question:   "Bagaimana cara mendaftar akun?",
			Answer:     "Klik tombol Daftar di pojok kanan atas halaman utama, isi email dan password Anda, lalu klik Daftar. Anda juga bisa mendaftar menggunakan akun Google untuk proses yang lebih cepat.",
			Category:   "account",
			OrderIndex: 1,
			IsActive:   true,
		},
		{
			Question:   "Bagaimana cara bergabung dengan komunitas?",
			Answer:     "Anda bisa bergabung dengan komunitas kami melalui tautan yang tersedia di halaman utama, atau mengikuti media sosial resmi JValleyverse untuk info terbaru dan undangan ke grup diskusi.",
			Category:   "general",
			OrderIndex: 3,
			IsActive:   true,
		},
		{
			Question:   "Siapa saja yang bisa berkontribusi?",
			Answer:     "Semua orang! Baik Anda developer, content creator, maupun enthusiast, Anda bisa berkontribusi dalam bentuk course, artikel, mentoring, atau sekadar berbagi pengalaman di komunitas.",
			Category:   "general",
			OrderIndex: 4,
			IsActive:   true,
		},
		{
			Question:   "Apakah saya bisa membuat course sendiri?",
			Answer:     "Tentu! Jika Anda ingin berbagi ilmu dan pengalaman, silakan hubungi tim kami melalui halaman Kontak atau DM sosial media kami untuk mulai berdiskusi soal pembuatan course.",
			Category:   "course",
			OrderIndex: 1,
			IsActive:   true,
		},
		{
			Question:   "Apakah ada event atau workshop?",
			Answer:     "Ya! Kami rutin mengadakan event seperti webinar, coding challenge, dan workshop online yang bisa diikuti secara gratis oleh anggota komunitas.",
			Category:   "general",
			OrderIndex: 5,
			IsActive:   true,
		},
		{
			Question:   "Bagaimana cara mendapatkan update terbaru?",
			Answer:     "Anda bisa berlangganan newsletter kami atau follow sosial media resmi JValleyverse untuk update course baru, event, dan konten komunitas lainnya.",
			Category:   "general",
			OrderIndex: 6,
			IsActive:   true,
		},
	}

	for _, faq := range faqs {
		if err := db.Create(&faq).Error; err != nil {
			return fmt.Errorf("failed to seed FAQ: %w", err)
		}
	}

	fmt.Printf("    ✓ %d FAQs seeded successfully\n", len(faqs))
	return nil
}
