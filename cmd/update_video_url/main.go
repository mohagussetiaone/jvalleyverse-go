package main

import (
	"fmt"
	"jvalleyverse/pkg/config"
	"jvalleyverse/pkg/database"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load config
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found, using system env")
	}
	config.LoadConfig()

	// Connect to database
	database.ConnectDB()
	db := database.DB

	// Video URL mappings by slug
	type lessonVideo struct {
		slug    string
		videoURL string
	}

	lessons := []lessonVideo{
		{slug: "pengenalan-go-setup-environment", videoURL: "https://www.youtube.com/watch?v=Q0sJQtgj1gY"},
		{slug: "fiber-framework-routing", videoURL: "https://www.youtube.com/watch?v=9qB-KfOI31E"},
		{slug: "gorm-database-integration", videoURL: "https://www.youtube.com/watch?v=Z0h0yRKKjQg"},
		{slug: "setup-react-typescript-project", videoURL: "https://www.youtube.com/watch?v=8a9hsZPmXms"},
		{slug: "membuat-komponen-dashboard", videoURL: "https://www.youtube.com/watch?v=6aBd6Zw3k8A"},
		{slug: "flutter-fundamentals-dart-basics", videoURL: "https://www.youtube.com/watch?v=CD1Y2xDOqC0"},
	}

	updated := 0
	for _, l := range lessons {
		result := db.Model(&struct{}{}).Table("lessons").
			Where("slug = ?", l.slug).
			Update("video_url", l.videoURL)
		if result.Error != nil {
			fmt.Printf("  ❌ Error updating %s: %v\n", l.slug, result.Error)
			continue
		}
		if result.RowsAffected > 0 {
			fmt.Printf("  ✅ Updated %s → %s\n", l.slug, l.videoURL)
			updated++
		} else {
			fmt.Printf("  ⚠️  No lesson found with slug: %s\n", l.slug)
		}
	}

	fmt.Printf("\n✅ Done! %d lessons updated with video_url.\n", updated)
}
