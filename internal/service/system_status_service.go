package service

import (
	"context"
	"os"
	"time"

	"jvalleyverse/internal/minio"
	"jvalleyverse/pkg/database"
	"jvalleyverse/pkg/redis"
)

type ServiceHealth struct {
	Name    string `json:"name"`
	Status  string `json:"status"` // "operational", "degraded", "down"
	Message string `json:"message,omitempty"`
	Latency string `json:"latency,omitempty"`
}

type SystemStatus struct {
	Status      string          `json:"status"` // "all_operational", "partial_outage", "major_outage"
	Uptime      string          `json:"uptime"`
	Version     string          `json:"version"`
	Environment string          `json:"environment"`
	Timestamp   time.Time       `json:"timestamp"`
	Services    []ServiceHealth `json:"services"`
	Summary     StatusSummary   `json:"summary"`
}

type StatusSummary struct {
	Total     int `json:"total"`
	Operational int `json:"operational"`
	Degraded  int `json:"degraded"`
	Down      int `json:"down"`
}

var startTime time.Time

func init() {
	startTime = time.Now()
}

// GetSystemStatus checks all service dependencies and returns comprehensive status
func GetSystemStatus(ctx context.Context) *SystemStatus {
	services := make([]ServiceHealth, 0)
	status := "all_operational"
	summary := StatusSummary{}

	// 1. Database
	dbHealth := checkDatabase(ctx)
	services = append(services, dbHealth)
	updateSummary(&summary, dbHealth)

	// 2. Redis
	redisHealth := checkRedis(ctx)
	services = append(services, redisHealth)
	updateSummary(&summary, redisHealth)

	// 3. MinIO
	minioHealth := checkMinIO(ctx)
	services = append(services, minioHealth)
	updateSummary(&summary, minioHealth)

	// Determine overall status
	if summary.Down > 0 {
		status = "major_outage"
	} else if summary.Degraded > 0 {
		status = "partial_outage"
	}

	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = os.Getenv("ENV")
	}
	if env == "" {
		env = "development"
	}

	uptime := time.Since(startTime).Round(time.Second).String()

	return &SystemStatus{
		Status:      status,
		Uptime:      uptime,
		Version:     "1.0.0",
		Environment: env,
		Timestamp:   time.Now().UTC(),
		Services:    services,
		Summary:     summary,
	}
}

func checkDatabase(ctx context.Context) ServiceHealth {
	start := time.Now()
	if database.DB == nil {
		return ServiceHealth{
			Name:    "database",
			Status:  "down",
			Message: "Database not connected",
			Latency: time.Since(start).String(),
		}
	}

	sqlDB, err := database.DB.DB()
	if err != nil {
		return ServiceHealth{
			Name:    "database",
			Status:  "down",
			Message: err.Error(),
			Latency: time.Since(start).String(),
		}
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return ServiceHealth{
			Name:    "database",
			Status:  "down",
			Message: err.Error(),
			Latency: time.Since(start).String(),
		}
	}

	return ServiceHealth{
		Name:    "database",
		Status:  "operational",
		Message: "PostgreSQL connected",
		Latency: time.Since(start).String(),
	}
}

func checkRedis(ctx context.Context) ServiceHealth {
	start := time.Now()
	if !redis.IsAvailable() {
		return ServiceHealth{
			Name:    "redis",
			Status:  "degraded",
			Message: "Redis not available (caching disabled)",
			Latency: time.Since(start).String(),
		}
	}

	if err := redis.Client.Ping(ctx).Err(); err != nil {
		return ServiceHealth{
			Name:    "redis",
			Status:  "down",
			Message: err.Error(),
			Latency: time.Since(start).String(),
		}
	}

	return ServiceHealth{
		Name:    "redis",
		Status:  "operational",
		Message: "Redis connected",
		Latency: time.Since(start).String(),
	}
}

func checkMinIO(ctx context.Context) ServiceHealth {
	start := time.Now()
	if !minio.IsAvailable() {
		return ServiceHealth{
			Name:    "minio",
			Status:  "degraded",
			Message: "MinIO not configured (file upload disabled)",
			Latency: time.Since(start).String(),
		}
	}

	return ServiceHealth{
		Name:    "minio",
		Status:  "operational",
		Message: "MinIO connected",
		Latency: time.Since(start).String(),
	}
}

func updateSummary(summary *StatusSummary, health ServiceHealth) {
	summary.Total++
	switch health.Status {
	case "operational":
		summary.Operational++
	case "degraded":
		summary.Degraded++
	case "down":
		summary.Down++
	}
}
