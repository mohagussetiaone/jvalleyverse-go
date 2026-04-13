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

    // Create ENUM types if they don't exist
    createEnumTypes()

    // Auto migrate schema
    if err := domain.AutoMigrate(DB); err != nil {
        panic("Migration failed: " + err.Error())
    }
}

// createEnumTypes creates PostgreSQL ENUM types for the application
func createEnumTypes() {
    enums := []struct {
        typeName string
        values   string
    }{
        {"userrole", "'admin','user'"},
        {"projectstatus", "'draft','published','archived'"},
        {"certificatestatus", "'issued','revoked','expired'"},
        {"discussionstatus", "'open','closed','pinned'"},
        {"showcaseStatus", "'draft','published','archived'"},
        {"showcasevisibility", "'public','private','draft'"},
        {"classDifficulty", "'beginner','intermediate','advanced'"},
        {"pointactivitytype", "'showcase_created','showcase_liked','discussion_created','discussion_reply','class_completed','certificate_issued'"},
    }

    for _, enum := range enums {
        // Check if enum exists, if not create it
        sql := fmt.Sprintf(`DO $$ BEGIN
            IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = '%s') THEN
                CREATE TYPE %s AS ENUM (%s);
            END IF;
        END $$;`, enum.typeName, enum.typeName, enum.values)
        
        if err := DB.Exec(sql).Error; err != nil {
            fmt.Printf("Warning: Could not create enum %s: %v\n", enum.typeName, err)
        }
    }
}