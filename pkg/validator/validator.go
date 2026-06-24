package validator

import (
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidateEmail(email string) *ValidationError {
	if strings.TrimSpace(email) == "" {
		return &ValidationError{Field: "email", Message: "Email is required"}
	}
	if !emailRegex.MatchString(email) {
		return &ValidationError{Field: "email", Message: "Email format is invalid"}
	}
	return nil
}

func ValidatePassword(password string) *ValidationError {
	if password == "" {
		return &ValidationError{Field: "password", Message: "Password is required"}
	}
	if len(password) < 6 {
		return &ValidationError{Field: "password", Message: "Password must be at least 6 characters"}
	}
	if len(password) > 72 {
		return &ValidationError{Field: "password", Message: "Password must be at most 72 characters"}
	}
	return nil
}

func ValidateName(name string) *ValidationError {
	if strings.TrimSpace(name) == "" {
		return &ValidationError{Field: "name", Message: "Name is required"}
	}
	if len(name) < 2 {
		return &ValidationError{Field: "name", Message: "Name must be at least 2 characters"}
	}
	if len(name) > 100 {
		return &ValidationError{Field: "name", Message: "Name must be at most 100 characters"}
	}
	return nil
}

func ValidateRequired(field, value, fieldName string) *ValidationError {
	if strings.TrimSpace(value) == "" {
		return &ValidationError{Field: field, Message: fieldName + " is required"}
	}
	return nil
}

func ValidateVisibility(visibility string) *ValidationError {
	if visibility == "" {
		return nil
	}
	visibility = strings.ToLower(visibility)
	if visibility != "public" && visibility != "private" {
		return &ValidationError{Field: "visibility", Message: "Visibility must be 'public' or 'private'"}
	}
	return nil
}
