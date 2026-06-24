package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"valid email", "test@example.com", false},
		{"valid email with dots", "test.user@example.co.id", false},
		{"empty email", "", true},
		{"no @", "invalid", true},
		{"no domain", "test@", true},
		{"no tld", "test@example", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, "email", err.Field)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"valid password", "password123", false},
		{"min length", "abc123", false},
		{"empty password", "", true},
		{"too short", "abc12", true},
		{"too long", string(make([]byte, 73)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, "password", err.Field)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid name", "John Doe", false},
		{"single char", "J", true},
		{"empty", "", true},
		{"spaces only", "  ", true},
		{"very long", string(make([]byte, 101)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.input)
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, "name", err.Field)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestValidateRequired(t *testing.T) {
	assert.NotNil(t, ValidateRequired("title", "", "Title"))
	assert.Nil(t, ValidateRequired("title", "My Title", "Title"))
}

func TestValidateVisibility(t *testing.T) {
	assert.Nil(t, ValidateVisibility("public"))
	assert.Nil(t, ValidateVisibility("private"))
	assert.Nil(t, ValidateVisibility(""))
	assert.NotNil(t, ValidateVisibility("secret"))
	assert.NotNil(t, ValidateVisibility("hidden"))
}
