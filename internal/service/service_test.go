package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateLevel(t *testing.T) {
	tests := []struct {
		name       string
		totalPoints int
		want       int
	}{
		{"beginner 0 points", 0, 1},
		{"beginner 50 points", 50, 1},
		{"beginner 99 points", 99, 1},
		{"intermediate 100 points", 100, 2},
		{"intermediate 200 points", 200, 2},
		{"intermediate 499 points", 499, 2},
		{"advanced 500 points", 500, 3},
		{"advanced 750 points", 750, 3},
		{"advanced 999 points", 999, 3},
		{"expert 1000 points", 1000, 4},
		{"expert 1500 points", 1500, 4},
		{"expert 1999 points", 1999, 4},
		{"master 2000 points", 2000, 5},
		{"master 5000 points", 5000, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateLevel(tt.totalPoints)
			assert.Equal(t, tt.want, got, "CalculateLevel(%d)", tt.totalPoints)
		})
	}
}

func TestCalculateLevelEdgeCases(t *testing.T) {
	assert.Equal(t, 1, CalculateLevel(-100), "negative points should be level 1")
	assert.Equal(t, 1, CalculateLevel(-1), "negative points should be level 1")
}
