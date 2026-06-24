package domain

import "errors"

var (
	ErrNotFound        = errors.New("resource not found")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrInternal        = errors.New("internal server error")
	ErrInvalidInput    = errors.New("invalid input")
	ErrUserNotFound    = errors.New("user not found")
	ErrCourseNotFound  = errors.New("course not found")
	ErrLessonNotFound    = errors.New("lesson not found")
	ErrStudyCaseNotFound = errors.New("study case not found")
	ErrEmailExists       = errors.New("email already exists")
)
