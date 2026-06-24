package database

import (
	"fmt"
	"jvalleyverse/internal/domain"
	"jvalleyverse/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	cfg := config.AppConfig
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	createEnumTypes()

	// Clean up old study_cases table if it exists with incompatible schema
	// (e.g. previous model had MentorID, now uses UserID with NOT NULL)
	DB.Exec("UPDATE discussions SET study_case_id = NULL WHERE study_case_id IS NOT NULL")
	DB.Exec("DROP TABLE IF EXISTS study_cases CASCADE")

	if err := domain.AutoMigrate(DB); err != nil {
		panic("Migration failed: " + err.Error())
	}
}

func createEnumTypes() {
	enums := []struct {
		typeName string
		values   string
	}{
		{"userrole", "'admin','user','mentor'"},
		{"coursestatus", "'draft','published','archived'"},
		{"certificatestatus", "'issued','revoked','expired'"},
		{"discussionstatus", "'open','closed','pinned'"},
		{"showcaseStatus", "'draft','published','archived'"},
		{"showcasevisibility", "'public','private','draft'"},
		{"lessondifficulty", "'beginner','intermediate','advanced'"},
		{"pointactivitytype", "'showcase_created','showcase_liked','discussion_created','discussion_reply','lesson_completed','certificate_issued'"},
	}

	for _, enum := range enums {
		sql := fmt.Sprintf(`DO $$ BEGIN
            IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = '%s') THEN
                CREATE TYPE %s AS ENUM (%s);
            END IF;
        END $$;`, enum.typeName, enum.typeName, enum.values)

		if err := DB.Exec(sql).Error; err != nil {
			fmt.Printf("Warning: Could not create enum %s: %v\n", enum.typeName, err)
		}
	}

	// Add 'mentor' value to existing userrole enum (safe to run repeatedly)
	DB.Exec("ALTER TYPE userrole ADD VALUE IF NOT EXISTS 'mentor'")
}
