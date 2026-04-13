package main

import (
	"fmt"
	"jvalleyverse/internal/service"
	"jvalleyverse/pkg/config"
	"jvalleyverse/pkg/database"
	"log"

	"github.com/joho/godotenv"
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
