package main

import (
	"fmt"
	"jvalleyverse/internal/service"
	"jvalleyverse/pkg/config"
	"jvalleyverse/pkg/database"
	"log"
)

func main() {
    // Load configuration
    config.LoadConfig()
    // Connect to the database
    database.ConnectDB()
    // Initialize services (repositories)
    service.InitServices(database.DB)

    // Run seeding
    if err := service.SeedInitialData(); err != nil {
        log.Fatalf("Seeding failed: %v", err)
    }
    fmt.Println("✅ Seeding completed successfully")
}
