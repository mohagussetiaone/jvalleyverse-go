package repository

import "gorm.io/gorm"

// Global db instance shared by all repositories
var db *gorm.DB

// InitRepository initializes the database connection for all repositories
func InitRepository(database *gorm.DB) {
	db = database
}

// GetDB returns the global database instance
func GetDB() *gorm.DB {
	return db
}
