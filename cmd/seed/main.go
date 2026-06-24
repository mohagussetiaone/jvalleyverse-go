package main

import (
	"fmt"
	"jvalleyverse/internal/service"
	"jvalleyverse/pkg/config"
	"jvalleyverse/pkg/database"
	"log"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadConfig() {
    err := godotenv.Load(".env")
    if err != nil {
        log.Println("No .env file found, using system env")
    }
}

func main() {
    // Load configuration
    config.LoadConfig()

    // Drop old study_cases table before migration to avoid NOT NULL constraint issues
    cfg := config.AppConfig
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
        cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)
    if tempDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}); err == nil {
        tempDB.Exec("DROP TABLE IF EXISTS study_cases CASCADE")
        tempDB.Exec("UPDATE discussions SET study_case_id = NULL")
        fmt.Println("  → Dropped old study_cases table")
    }

    // Connect to the database
    database.ConnectDB()
    // Initialize services (repositories)
    service.InitServices(database.DB)

    // Run seeding
    if err := service.SeedInitialData(database.DB); err != nil {
        log.Fatalf("Seeding failed: %v", err)
    }
    fmt.Println("✅ Seeding completed successfully")
}
