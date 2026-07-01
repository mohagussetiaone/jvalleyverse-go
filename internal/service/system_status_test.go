package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateSummary(t *testing.T) {
	tests := []struct {
		name              string
		statuses          []string
		wantTotal         int
		wantOperational   int
		wantDegraded      int
		wantDown          int
	}{
		{
			name:              "all operational",
			statuses:          []string{"operational", "operational", "operational"},
			wantTotal:         3,
			wantOperational:   3,
			wantDegraded:      0,
			wantDown:          0,
		},
		{
			name:              "one degraded",
			statuses:          []string{"operational", "degraded", "operational"},
			wantTotal:         3,
			wantOperational:   2,
			wantDegraded:      1,
			wantDown:          0,
		},
		{
			name:              "one down",
			statuses:          []string{"operational", "operational", "down"},
			wantTotal:         3,
			wantOperational:   2,
			wantDegraded:      0,
			wantDown:          1,
		},
		{
			name:              "all down",
			statuses:          []string{"down", "down", "down"},
			wantTotal:         3,
			wantOperational:   0,
			wantDegraded:      0,
			wantDown:          3,
		},
		{
			name:              "mixed statuses",
			statuses:          []string{"operational", "degraded", "down"},
			wantTotal:         3,
			wantOperational:   1,
			wantDegraded:      1,
			wantDown:          1,
		},
		{
			name:              "empty services",
			statuses:          []string{},
			wantTotal:         0,
			wantOperational:   0,
			wantDegraded:      0,
			wantDown:          0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := StatusSummary{}
			for _, status := range tt.statuses {
				health := ServiceHealth{Name: "test", Status: status}
				updateSummary(&summary, health)
			}
			assert.Equal(t, tt.wantTotal, summary.Total)
			assert.Equal(t, tt.wantOperational, summary.Operational)
			assert.Equal(t, tt.wantDegraded, summary.Degraded)
			assert.Equal(t, tt.wantDown, summary.Down)
		})
	}
}

func TestGetSystemStatus_Structure(t *testing.T) {
	// Test that GetSystemStatus returns a properly structured response
	// without actually connecting to DB/Redis/MinIO (they're nil in tests)
	status := GetSystemStatus(context.Background())

	assert.NotNil(t, status)
	assert.NotEmpty(t, status.Uptime)
	assert.Equal(t, "1.0.0", status.Version)
	assert.NotEmpty(t, status.Environment)
	assert.NotEmpty(t, status.Timestamp)
	assert.NotNil(t, status.Services)
	assert.NotNil(t, status.Summary)

	// Should have 3 services: database, redis, minio
	assert.Equal(t, 3, len(status.Services))
	assert.Equal(t, 3, status.Summary.Total)

	// In test environment, all services should be down/degraded
	// since DB, Redis, MinIO are not connected
	assert.GreaterOrEqual(t, status.Summary.Down+status.Summary.Degraded, 1)
}

func TestGetSystemStatus_Environment(t *testing.T) {
	// Default environment should be "development" when no ENV set
	status := GetSystemStatus(context.Background())
	assert.Equal(t, "development", status.Environment)
}

func TestGetSystemStatus_ServiceNames(t *testing.T) {
	status := GetSystemStatus(context.Background())
	names := make([]string, len(status.Services))
	for i, svc := range status.Services {
		names[i] = svc.Name
	}
	assert.Contains(t, names, "database")
	assert.Contains(t, names, "redis")
	assert.Contains(t, names, "minio")
}

func TestGetSystemStatus_ServiceHealth(t *testing.T) {
	status := GetSystemStatus(context.Background())

	for _, svc := range status.Services {
		assert.NotEmpty(t, svc.Name, "service name should not be empty")
		assert.Contains(t, []string{"operational", "degraded", "down"}, svc.Status,
			"service %s has invalid status: %s", svc.Name, svc.Status)
		assert.NotEmpty(t, svc.Latency, "service %s should have latency", svc.Name)
	}
}

func TestServiceHealthStruct(t *testing.T) {
	health := ServiceHealth{
		Name:    "test-service",
		Status:  "operational",
		Message: "All good",
		Latency: "5ms",
	}

	assert.Equal(t, "test-service", health.Name)
	assert.Equal(t, "operational", health.Status)
	assert.Equal(t, "All good", health.Message)
	assert.Equal(t, "5ms", health.Latency)
}

func TestServiceHealth_EmptyIsAllowed(t *testing.T) {
	// ServiceHealth with empty Status is valid (though we validate against enum)
	health := ServiceHealth{}
	assert.Empty(t, health.Name)
	assert.Empty(t, health.Status)
	assert.Empty(t, health.Message)
	assert.Empty(t, health.Latency)
}

func TestSystemStatusStruct(t *testing.T) {
	status := &SystemStatus{
		Status:      "all_operational",
		Uptime:      "10h",
		Version:     "1.0.0",
		Environment: "production",
		Services:    []ServiceHealth{},
		Summary: StatusSummary{
			Total:       0,
			Operational: 0,
			Degraded:    0,
			Down:        0,
		},
	}

	assert.Equal(t, "all_operational", status.Status)
	assert.Equal(t, "10h", status.Uptime)
	assert.Equal(t, "production", status.Environment)
}
