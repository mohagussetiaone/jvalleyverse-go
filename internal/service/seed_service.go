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
func SeedInitialData() error {
	db := getDB()
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
	// 4. SEED PROJECTS
	// =========================================================================
	fmt.Println("  → Seeding projects...")
	projects, err := seedProjects(db, users["admin"], categories)
	if err != nil {
		return fmt.Errorf("failed to seed projects: %w", err)
	}

	// =========================================================================
	// 5. SEED CLASSES + CLASS DETAILS
	// =========================================================================
	fmt.Println("  → Seeding classes & details...")
	classes, err := seedClasses(db, users["admin"], projects)
	if err != nil {
		return fmt.Errorf("failed to seed classes: %w", err)
	}

	// =========================================================================
	// 6. SEED CLASS PROGRESS
	// =========================================================================
	fmt.Println("  → Seeding class progress...")
	if err := seedClassProgress(db, users, classes); err != nil {
		return fmt.Errorf("failed to seed class progress: %w", err)
	}

	// =========================================================================
	// 7. SEED DISCUSSIONS & REPLIES
	// =========================================================================
	fmt.Println("  → Seeding discussions & replies...")
	if err := seedDiscussions(db, users, classes, categories); err != nil {
		return fmt.Errorf("failed to seed discussions: %w", err)
	}

	// =========================================================================
	// 8. SEED SHOWCASES
	// =========================================================================
	fmt.Println("  → Seeding showcases...")
	if err := seedShowcases(db, users, categories); err != nil {
		return fmt.Errorf("failed to seed showcases: %w", err)
	}

	fmt.Println("🌱 Seeding completed!")
	return nil
}

// ============================================================================
// HELPER: get the global DB from repository package
// ============================================================================

func getDB() *gorm.DB {
	// Access the DB via a new repository that exposes it
	repo := newSeedRepository()
	return repo.db
}

type seedRepository struct {
	db *gorm.DB
}

func newSeedRepository() *seedRepository {
	// We access the DB through the repository package's global var
	// Since InitServices has been called, repositories have the DB
	return &seedRepository{db: getSeedDB()}
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
	}

	result := make(map[string]*domain.User)
	keys := []string{"admin", "budi", "siti", "andi", "dewi"}

	for i, userData := range usersData {
		user := userData // create a copy
		if err := db.Where("email = ?", user.Email).FirstOrCreate(&user).Error; err != nil {
			return nil, err
		}
		result[keys[i]] = &user
		fmt.Printf("    ✓ User: %s (ID: %d, Role: %s)\n", user.Name, user.ID, user.Role)
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
		fmt.Printf("    ✓ Category: %s (ID: %d)\n", cat.Name, cat.ID)
	}

	return result, nil
}

// ============================================================================
// 4. SEED PROJECTS
// ============================================================================

func seedProjects(db *gorm.DB, admin *domain.User, categories map[string]*domain.Category) (map[string]*domain.Project, error) {
	projectsData := []struct {
		key         string
		categoryKey string
		project     domain.Project
	}{
		{
			key:         "go-rest-api",
			categoryKey: "backend-development",
			project: domain.Project{
				Title:       "Membangun REST API dengan Go & Fiber",
				Description: "Pelajari cara membuat RESTful API yang scalable menggunakan Go, Fiber framework, GORM, dan PostgreSQL. Dari setup project hingga deployment.",
				Thumbnail:   "https://images.unsplash.com/photo-1555066931-4365d14bab8c?w=800",
				Visibility:  "public",
			},
		},
		{
			key:         "react-dashboard",
			categoryKey: "frontend-development",
			project: domain.Project{
				Title:       "Build Modern Dashboard dengan React & TypeScript",
				Description: "Membangun dashboard interaktif menggunakan React, TypeScript, Tailwind CSS, dan Chart.js. Lengkap dengan autentikasi dan state management.",
				Thumbnail:   "https://images.unsplash.com/photo-1551288049-bebda4e38f71?w=800",
				Visibility:  "public",
			},
		},
		{
			key:         "flutter-ecommerce",
			categoryKey: "mobile-development",
			project: domain.Project{
				Title:       "Flutter E-Commerce App dari Nol",
				Description: "Membuat aplikasi e-commerce lengkap dengan Flutter. Mulai dari UI design, state management dengan BLoC, integrasi API, hingga payment gateway.",
				Thumbnail:   "https://images.unsplash.com/photo-1512941937669-90a1b58e7e9c?w=800",
				Visibility:  "public",
			},
		},
		{
			key:         "fullstack-nextjs",
			categoryKey: "web-development",
			project: domain.Project{
				Title:       "Fullstack Web App dengan Next.js 14",
				Description: "Panduan lengkap membangun aplikasi web fullstack dengan Next.js 14, Server Components, Prisma ORM, dan Vercel deployment.",
				Thumbnail:   "https://images.unsplash.com/photo-1460925895917-afdab827c52f?w=800",
				Visibility:  "public",
			},
		},
		{
			key:         "docker-kubernetes",
			categoryKey: "devops-cloud",
			project: domain.Project{
				Title:       "Docker & Kubernetes untuk Developer",
				Description: "Menguasai containerization dan orchestration. Dari Docker basics hingga Kubernetes cluster management di production.",
				Thumbnail:   "https://images.unsplash.com/photo-1667372393119-3d4c48d07fc9?w=800",
				Visibility:  "public",
			},
		},
	}

	result := make(map[string]*domain.Project)

	for _, pd := range projectsData {
		project := pd.project
		project.AdminID = admin.ID
		project.CategoryID = categories[pd.categoryKey].ID

		if err := db.Where("title = ? AND admin_id = ?", project.Title, project.AdminID).FirstOrCreate(&project).Error; err != nil {
			return nil, err
		}
		result[pd.key] = &project
		fmt.Printf("    ✓ Project: %s (ID: %d)\n", project.Title, project.ID)
	}

	return result, nil
}

// ============================================================================
// 5. SEED CLASSES + CLASS DETAILS
// ============================================================================

func seedClasses(db *gorm.DB, admin *domain.User, projects map[string]*domain.Project) (map[string]*domain.Class, error) {
	classesData := []struct {
		key        string
		projectKey string
		class      domain.Class
		detail     domain.ClassDetail
	}{
		// === GO REST API Project Classes ===
		{
			key:        "go-intro",
			projectKey: "go-rest-api",
			class: domain.Class{
				Title:       "Pengenalan Go & Setup Environment",
				Slug:        "pengenalan-go-setup-environment",
				Description: "Mengenal bahasa Go, instalasi tools, dan setup project pertama",
				Thumbnail:   "https://images.unsplash.com/photo-1515879218367-8466d910auj8?w=800",
				Difficulty:  "beginner",
				Duration:    45,
				OrderIndex:  1,
				SequenceNum: 1,
				IsFirst:     true,
				Visibility:  "public",
			},
			detail: domain.ClassDetail{
				About: "Di kelas ini kamu akan belajar dasar-dasar bahasa Go (Golang), mulai dari instalasi, konfigurasi environment, hingga membuat program pertama. Go adalah bahasa yang dikembangkan oleh Google, terkenal dengan performa tinggi dan kesederhanaan syntaxnya.",
				Rules: "1. Pastikan sudah menginstall Go versi terbaru\n2. Gunakan code editor (VS Code direkomendasikan)\n3. Selesaikan semua latihan sebelum lanjut ke kelas berikutnya\n4. Jangan ragu bertanya di forum diskusi",
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
			projectKey: "go-rest-api",
			class: domain.Class{
				Title:       "Fiber Framework & Routing",
				Slug:        "fiber-framework-routing",
				Description: "Belajar Fiber framework, routing, middleware, dan request handling",
				Thumbnail:   "https://images.unsplash.com/photo-1516116216624-53e697fedbea?w=800",
				Difficulty:  "beginner",
				Duration:    60,
				OrderIndex:  2,
				SequenceNum: 2,
				IsFirst:     false,
				Visibility:  "public",
			},
			detail: domain.ClassDetail{
				About: "Fiber adalah web framework Go yang terinspirasi dari Express.js. Di kelas ini kamu akan belajar setup Fiber, membuat route, menggunakan middleware, dan handling request/response.",
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
			projectKey: "go-rest-api",
			class: domain.Class{
				Title:       "GORM & Database Integration",
				Slug:        "gorm-database-integration",
				Description: "Menghubungkan Go dengan PostgreSQL menggunakan GORM ORM",
				Thumbnail:   "https://images.unsplash.com/photo-1544383835-bda2bc66a55d?w=800",
				Difficulty:  "intermediate",
				Duration:    90,
				OrderIndex:  3,
				SequenceNum: 3,
				IsFirst:     false,
				Visibility:  "public",
			},
			detail: domain.ClassDetail{
				About: "GORM adalah ORM paling populer untuk Go. Kamu akan belajar koneksi database, model definition, migration, CRUD operations, relationships, dan query optimization.",
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

		// === REACT DASHBOARD Project Classes ===
		{
			key:        "react-setup",
			projectKey: "react-dashboard",
			class: domain.Class{
				Title:       "Setup React + TypeScript Project",
				Slug:        "setup-react-typescript-project",
				Description: "Setup project React dengan Vite, TypeScript, dan Tailwind CSS",
				Thumbnail:   "https://images.unsplash.com/photo-1633356122102-3fe601e05bd2?w=800",
				Difficulty:  "beginner",
				Duration:    30,
				OrderIndex:  1,
				SequenceNum: 1,
				IsFirst:     true,
				Visibility:  "public",
			},
			detail: domain.ClassDetail{
				About: "Memulai project React modern menggunakan Vite sebagai build tool, TypeScript untuk type safety, dan Tailwind CSS untuk styling yang cepat dan konsisten.",
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
			projectKey: "react-dashboard",
			class: domain.Class{
				Title:       "Membuat Komponen Dashboard",
				Slug:        "membuat-komponen-dashboard",
				Description: "Membangun komponen-komponen UI dashboard: Sidebar, Navbar, Cards, Charts",
				Thumbnail:   "https://images.unsplash.com/photo-1551288049-bebda4e38f71?w=800",
				Difficulty:  "intermediate",
				Duration:    75,
				OrderIndex:  2,
				SequenceNum: 2,
				IsFirst:     false,
				Visibility:  "public",
			},
			detail: domain.ClassDetail{
				About: "Belajar membuat komponen-komponen UI yang reusable untuk dashboard. Mulai dari layout system, sidebar navigation, stats cards, hingga data table dengan sorting dan filtering.",
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

		// === FLUTTER E-COMMERCE Project Classes ===
		{
			key:        "flutter-intro",
			projectKey: "flutter-ecommerce",
			class: domain.Class{
				Title:       "Flutter Fundamentals & Dart Basics",
				Slug:        "flutter-fundamentals-dart-basics",
				Description: "Mengenal Flutter framework dan dasar-dasar bahasa Dart",
				Thumbnail:   "https://images.unsplash.com/photo-1617040619263-41c5a9ca7521?w=800",
				Difficulty:  "beginner",
				Duration:    60,
				OrderIndex:  1,
				SequenceNum: 1,
				IsFirst:     true,
				Visibility:  "public",
			},
			detail: domain.ClassDetail{
				About: "Pengenalan Flutter dan Dart untuk pemula. Belajar widget tree, layout system, state management dasar, dan navigasi antar halaman.",
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

	result := make(map[string]*domain.Class)

	for _, cd := range classesData {
		class := cd.class
		class.AdminID = admin.ID
		class.ProjectID = projects[cd.projectKey].ID

		if err := db.Where("slug = ? AND project_id = ?", class.Slug, class.ProjectID).FirstOrCreate(&class).Error; err != nil {
			return nil, err
		}
		result[cd.key] = &class

		// Seed ClassDetail
		detail := cd.detail
		detail.ClassID = class.ID
		if err := db.Where("class_id = ?", detail.ClassID).FirstOrCreate(&detail).Error; err != nil {
			return nil, err
		}

		fmt.Printf("    ✓ Class: %s (ID: %d, Project: %s)\n", class.Title, class.ID, cd.projectKey)
	}

	// Link classes with NextClassID for progression
	linkClassProgression(db, result, "go-intro", "go-fiber-basics")
	linkClassProgression(db, result, "go-fiber-basics", "go-gorm-db")
	linkClassProgression(db, result, "react-setup", "react-components")

	return result, nil
}

func linkClassProgression(db *gorm.DB, classes map[string]*domain.Class, fromKey, toKey string) {
	from := classes[fromKey]
	to := classes[toKey]
	if from != nil && to != nil {
		db.Model(&domain.Class{}).Where("id = ?", from.ID).Update("next_class_id", to.ID)
	}
}

// ============================================================================
// 6. SEED CLASS PROGRESS
// ============================================================================

func seedClassProgress(db *gorm.DB, users map[string]*domain.User, classes map[string]*domain.Class) error {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	lastWeek := now.Add(-7 * 24 * time.Hour)

	progressData := []domain.ClassProgress{
		// Budi sudah selesai kelas Go intro, sedang di Fiber
		{
			UserID:             users["budi"].ID,
			ClassID:            classes["go-intro"].ID,
			Status:             "completed",
			StartedAt:          &lastWeek,
			CompletedAt:        &yesterday,
			ProgressPercentage: 100,
			Notes:              "Kelas yang sangat bagus untuk pemula Go!",
		},
		{
			UserID:             users["budi"].ID,
			ClassID:            classes["go-fiber-basics"].ID,
			Status:             "in_progress",
			StartedAt:          &yesterday,
			ProgressPercentage: 45,
			Notes:              "Sedang belajar middleware",
		},
		// Siti baru mulai React
		{
			UserID:             users["siti"].ID,
			ClassID:            classes["react-setup"].ID,
			Status:             "started",
			StartedAt:          &now,
			ProgressPercentage: 10,
		},
		// Andi sudah advanced, selesai beberapa kelas
		{
			UserID:             users["andi"].ID,
			ClassID:            classes["go-intro"].ID,
			Status:             "completed",
			StartedAt:          &lastWeek,
			CompletedAt:        &lastWeek,
			ProgressPercentage: 100,
		},
		{
			UserID:             users["andi"].ID,
			ClassID:            classes["go-fiber-basics"].ID,
			Status:             "completed",
			StartedAt:          &lastWeek,
			CompletedAt:        &yesterday,
			ProgressPercentage: 100,
		},
		{
			UserID:             users["andi"].ID,
			ClassID:            classes["go-gorm-db"].ID,
			Status:             "in_progress",
			StartedAt:          &now,
			ProgressPercentage: 30,
			Notes:              "Relationships di GORM lumayan tricky",
		},
		// Dewi baru mulai Flutter
		{
			UserID:             users["dewi"].ID,
			ClassID:            classes["flutter-intro"].ID,
			Status:             "started",
			StartedAt:          &now,
			ProgressPercentage: 5,
		},
	}

	for _, progress := range progressData {
		p := progress
		if err := db.Where("user_id = ? AND class_id = ?", p.UserID, p.ClassID).FirstOrCreate(&p).Error; err != nil {
			return err
		}
	}
	fmt.Printf("    ✓ Created %d class progress records\n", len(progressData))

	return nil
}

// ============================================================================
// 7. SEED DISCUSSIONS & REPLIES
// ============================================================================

func seedDiscussions(db *gorm.DB, users map[string]*domain.User, classes map[string]*domain.Class, categories map[string]*domain.Category) error {
	// Discussion 1: Budi bertanya tentang Go
	disc1 := domain.Discussion{
		Title:      "Perbedaan goroutine dan thread biasa?",
		Content:    "Halo semua, saya baru belajar Go dan masih bingung bedanya goroutine dengan thread biasa di bahasa lain. Bisa tolong jelaskan dengan contoh sederhana?\n\nSaya sudah baca dokumentasi tapi masih kurang paham bagian scheduler-nya.",
		UserID:     users["budi"].ID,
		ClassID:    &classes["go-intro"].ID,
		CategoryID: categories["backend-development"].ID,
		ViewsCount: 42,
		Status:     "open",
	}
	if err := db.Where("title = ? AND user_id = ?", disc1.Title, disc1.UserID).FirstOrCreate(&disc1).Error; err != nil {
		return err
	}

	// Replies for Discussion 1
	reply1 := domain.Reply{
		Content:      "Goroutine itu lightweight thread yang dikelola oleh Go runtime, bukan OS. Satu OS thread bisa menjalankan ribuan goroutine karena Go punya scheduler sendiri (M:N scheduling).\n\nContoh sederhana:\n```go\ngo func() {\n    fmt.Println(\"Hello dari goroutine!\")\n}()\n```\n\nBeda dengan thread biasa yang berat (biasanya 1-8MB stack), goroutine cuma pakai ~2KB stack yang bisa grow/shrink otomatis.",
		UserID:       users["andi"].ID,
		DiscussionID: disc1.ID,
		LikesCount:   5,
		IsMarkedBest: true,
	}
	if err := db.Where("discussion_id = ? AND user_id = ? AND is_marked_best = ?", reply1.DiscussionID, reply1.UserID, true).FirstOrCreate(&reply1).Error; err != nil {
		return err
	}

	reply2 := domain.Reply{
		Content:      "Tambahannya, goroutine juga lebih murah untuk context switching karena Go scheduler yang handle, bukan kernel. Makanya bisa jalankan puluhan ribu goroutine tanpa masalah 🚀",
		UserID:       users["admin"].ID,
		DiscussionID: disc1.ID,
		LikesCount:   3,
	}
	if err := db.Where("discussion_id = ? AND user_id = ? AND content LIKE ?", reply2.DiscussionID, reply2.UserID, "%context switching%").FirstOrCreate(&reply2).Error; err != nil {
		return err
	}

	// Nested reply
	nestedReply := domain.Reply{
		Content:      "Terima kasih penjelasannya! Sekarang sudah lebih paham. Jadi intinya goroutine itu lebih efisien dari thread ya 👍",
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
		Content:    "Saya lagi belajar React hooks dan bingung kapan harus pakai useState dan kapan pakai useReducer. Ada rule of thumb-nya gak?\n\nTerutama untuk form yang kompleks, lebih baik pakai yang mana?",
		UserID:     users["siti"].ID,
		ClassID:    &classes["react-setup"].ID,
		CategoryID: categories["frontend-development"].ID,
		ViewsCount: 28,
		Status:     "open",
	}
	if err := db.Where("title = ? AND user_id = ?", disc2.Title, disc2.UserID).FirstOrCreate(&disc2).Error; err != nil {
		return err
	}

	reply3 := domain.Reply{
		Content:      "Rule of thumb saya:\n- **useState**: untuk state sederhana (toggle, single value, counter)\n- **useReducer**: untuk state kompleks (form dengan banyak field, state yang saling bergantung)\n\nUseReducer juga bagus kalau kamu perlu predictable state transitions. Mirip seperti Redux tapi built-in.",
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
		Content:    "Halo teman-teman! Saya mau share tips belajar programming berdasarkan pengalaman saya:\n\n1. Konsisten lebih penting dari intensitas\n2. Jangan cuma nonton tutorial, langsung praktik\n3. Buat project sendiri, sekecil apapun\n4. Gabung komunitas untuk diskusi\n5. Jangan takut error, itu bagian dari proses\n\nShare juga tips kalian dong!",
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
	showcasesData := []domain.Showcase{
		{
			Title:       "Portfolio Website Personal",
			Description: "Website portfolio menggunakan Next.js 14 dan Tailwind CSS. Fitur: dark mode, animasi smooth, blog dengan MDX, dan SEO optimized. Responsive untuk semua device.",
			MediaURLs:   mustJSON([]string{"https://images.unsplash.com/photo-1517180102446-f3ece451e9d8?w=800", "https://images.unsplash.com/photo-1460925895917-afdab827c52f?w=800"}),
			UserID:      users["budi"].ID,
			CategoryID:  categories["web-development"].ID,
			Status:      "published",
			Visibility:  "public",
			LikesCount:  12,
			ViewsCount:  89,
		},
		{
			Title:       "REST API Microservices Go",
			Description: "Arsitektur microservices menggunakan Go, gRPC, dan RabbitMQ. Termasuk API gateway, service discovery, dan distributed tracing dengan Jaeger.",
			MediaURLs:   mustJSON([]string{"https://images.unsplash.com/photo-1558494949-ef010cbdcc31?w=800"}),
			UserID:      users["andi"].ID,
			CategoryID:  categories["backend-development"].ID,
			Status:      "published",
			Visibility:  "public",
			LikesCount:  23,
			ViewsCount:  145,
		},
		{
			Title:       "Mobile App UI Kit Design",
			Description: "Koleksi UI Kit untuk aplikasi mobile modern. Termasuk 50+ screen design untuk e-commerce, social media, dan productivity app. Dibuat dengan Figma dan diimplementasi dengan Flutter.",
			MediaURLs:   mustJSON([]string{"https://images.unsplash.com/photo-1512941937669-90a1b58e7e9c?w=800", "https://images.unsplash.com/photo-1616469829581-73993eb86b02?w=800"}),
			UserID:      users["siti"].ID,
			CategoryID:  categories["ui-ux-design"].ID,
			Status:      "published",
			Visibility:  "public",
			LikesCount:  18,
			ViewsCount:  112,
		},
		{
			Title:       "Aplikasi Todo Sederhana Flutter",
			Description: "Aplikasi todo list pertama saya pakai Flutter! Fitur: CRUD task, reminder, dark mode. Masih belajar tapi seneng bisa bikin ini 😊",
			MediaURLs:   mustJSON([]string{"https://images.unsplash.com/photo-1611224923853-80b023f02d71?w=800"}),
			UserID:      users["dewi"].ID,
			CategoryID:  categories["mobile-development"].ID,
			Status:      "published",
			Visibility:  "public",
			LikesCount:  7,
			ViewsCount:  34,
		},
	}

	for _, showcaseData := range showcasesData {
		showcase := showcaseData
		if err := db.Where("title = ? AND user_id = ?", showcase.Title, showcase.UserID).FirstOrCreate(&showcase).Error; err != nil {
			return err
		}
		fmt.Printf("    ✓ Showcase: %s by User#%d\n", showcase.Title, showcase.UserID)
	}

	// Seed some showcase likes
	likes := []domain.ShowcaseLike{
		{UserID: users["andi"].ID, ShowcaseID: 1, CreatedAt: time.Now()},
		{UserID: users["siti"].ID, ShowcaseID: 1, CreatedAt: time.Now()},
		{UserID: users["budi"].ID, ShowcaseID: 2, CreatedAt: time.Now()},
		{UserID: users["dewi"].ID, ShowcaseID: 2, CreatedAt: time.Now()},
		{UserID: users["budi"].ID, ShowcaseID: 3, CreatedAt: time.Now()},
		{UserID: users["andi"].ID, ShowcaseID: 4, CreatedAt: time.Now()},
	}

	for _, like := range likes {
		l := like
		db.Where("user_id = ? AND showcase_id = ?", l.UserID, l.ShowcaseID).FirstOrCreate(&l)
	}

	// Seed showcase comments
	comments := []domain.ShowcaseComment{
		{
			Content:    "Keren banget portfolionya! Animasinya smooth 🔥",
			UserID:     users["andi"].ID,
			ShowcaseID: 1,
		},
		{
			Content:    "Arsitektur microservicesnya rapi, boleh share repo-nya?",
			UserID:     users["budi"].ID,
			ShowcaseID: 2,
		},
		{
			Content:    "UI Kit-nya cantik banget! Warna-warnanya harmonious 😍",
			UserID:     users["dewi"].ID,
			ShowcaseID: 3,
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
