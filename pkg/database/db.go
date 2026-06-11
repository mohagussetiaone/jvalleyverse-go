package database

import (
	"errors"
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

	// Backfill legacy classes before AutoMigrate enforces phase_id NOT NULL.
	if err := migrateLegacyClassPhases(DB); err != nil {
		panic("Legacy phase migration failed: " + err.Error())
	}

	// Auto migrate schema
	if err := domain.AutoMigrate(DB); err != nil {
		panic("Migration failed: " + err.Error())
	}
}

func migrateLegacyClassPhases(db *gorm.DB) error {
	if !db.Migrator().HasTable(&domain.Project{}) || !db.Migrator().HasTable(&domain.Class{}) {
		return nil
	}

	if err := db.AutoMigrate(&domain.Phase{}); err != nil {
		return err
	}

	if !db.Migrator().HasColumn(&domain.Class{}, "phase_id") {
		if err := db.Exec(`ALTER TABLE "classes" ADD COLUMN "phase_id" text`).Error; err != nil {
			return err
		}
	}

	type legacyProjectRow struct {
		ProjectID string
	}

	var rows []legacyProjectRow
	if err := db.Raw(`
		SELECT DISTINCT project_id
		FROM classes
		WHERE deleted_at IS NULL
		  AND project_id IS NOT NULL
		  AND (phase_id IS NULL OR phase_id = '')
	`).Scan(&rows).Error; err != nil {
		return err
	}

	for _, row := range rows {
		if row.ProjectID == "" {
			continue
		}

		if err := db.Transaction(func(tx *gorm.DB) error {
			var phase domain.Phase
			err := tx.Where("project_id = ?", row.ProjectID).
				Order("order_index ASC").
				First(&phase).Error
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}

				phase = domain.Phase{
					ProjectID:   row.ProjectID,
					Title:       "General",
					Description: "Auto-generated default phase for legacy classes",
					OrderIndex:  0,
				}
				if err := tx.Create(&phase).Error; err != nil {
					return err
				}
			}

			return tx.Model(&domain.Class{}).
				Where("project_id = ? AND (phase_id IS NULL OR phase_id = '')", row.ProjectID).
				Update("phase_id", phase.ID).Error
		}); err != nil {
			return err
		}
	}

	return nil
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
