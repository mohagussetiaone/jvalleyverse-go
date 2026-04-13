package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
    Port        string
    JWTSecret   string
    JWTExpiry   time.Duration
    DBHost      string
    DBPort      string
    DBUser      string
    DBPassword  string
    DBName      string
    DBSSLMode   string
    RedisHost   string
    RedisPass   string
    RedisDB     int
    AdminEmail  string
    AdminPass   string
}

var AppConfig *Config

func LoadConfig() {
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using system env")
    }

    expiry, _ := time.ParseDuration(os.Getenv("JWT_EXPIRY"))
    if expiry == 0 {
        expiry = 24 * time.Hour
    }

    AppConfig = &Config{
        Port:       getEnv("PORT", "3000"),
        JWTSecret:  getEnv("JWT_SECRET", "defaultsecret"),
        JWTExpiry:  expiry,
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", ""),
        DBName:     getEnv("DB_NAME", "jvalleyverse"),
        DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
        RedisHost:  getEnv("REDIS_HOST", "localhost:6379"),
        RedisPass:  getEnv("REDIS_PASSWORD", ""),
        RedisDB:    getEnvInt("REDIS_DB", 0),
        AdminEmail: getEnv("ADMIN_EMAIL", "admin@example.com"),
        AdminPass:  getEnv("ADMIN_PASSWORD", "admin123"),
    }
}

func getEnv(key, fallback string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return fallback
}

func getEnvInt(key string, fallback int) int {
    // implementasi sederhana
    return fallback
}