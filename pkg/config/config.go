package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	JWTSecret      string
	JWTExpiry      time.Duration
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	RedisHost      string
	RedisPass      string
	RedisDB        int
	AdminEmail     string
	AdminPass      string
	GoogleClientID string
	CORSOrigins    string
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioCDNURL    string
	MinioUseSSL    bool
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

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required. Set it in .env or environment variables.")
	}

	AppConfig = &Config{
		Port:           getEnv("PORT", "3000"),
		JWTSecret:      jwtSecret,
		JWTExpiry:      expiry,
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBName:         getEnv("DB_NAME", "jvalleyverse"),
		DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
		RedisHost:      getEnv("REDIS_HOST", "localhost:6379"),
		RedisPass:      getEnv("REDIS_PASSWORD", ""),
		RedisDB:        getEnvInt("REDIS_DB", 0),
		AdminEmail:     getEnv("ADMIN_EMAIL", "admin@example.com"),
		AdminPass:      getEnv("ADMIN_PASSWORD", "admin123"),
		GoogleClientID: getEnv("GOOGLE_CLIENT_ID", ""),
		CORSOrigins:    getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:5173,https://localhost:5174,https://jvalleyverse.web.id"),
		MinioEndpoint:  getEnv("MINIO_ENDPOINT", "minio.mohagussetiaone.my.id"),
		MinioAccessKey: getEnv("MINIO_ACCESS_KEY", ""),
		MinioSecretKey: getEnv("MINIO_SECRET_KEY", ""),
		MinioBucket:    getEnv("MINIO_BUCKET", "jvalleyverse"),
		MinioCDNURL:    getEnv("MINIO_CDN_URL", "https://cdn.mohagussetiaone.my.id"),
		MinioUseSSL:    getEnv("MINIO_USE_SSL", "true") == "true",
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	i, err := parseInt(val)
	if err != nil {
		return fallback
	}
	return i
}

func parseInt(s string) (int, error) {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("not a number: %s", s)
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}
