package handler

import (
	"errors"
	"jvalleyverse/internal/domain"
)

func mapServiceErrorToStatus(err error) int {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return 400
	case errors.Is(err, domain.ErrUnauthorized):
		return 401
	case errors.Is(err, domain.ErrForbidden):
		return 403
	case errors.Is(err, domain.ErrNotFound),
		errors.Is(err, domain.ErrProjectNotFound),
		errors.Is(err, domain.ErrClassNotFound),
		errors.Is(err, domain.ErrUserNotFound):
		return 404
	default:
		return 500
	}
}
