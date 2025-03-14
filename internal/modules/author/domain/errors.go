package domain

import (
	"go-monolith/pkg/errors"
)

// AuthorError represents author-specific domain errors
type AuthorError struct {
	errors.BaseError
}

func NewAuthorError(message string) error {
	return &AuthorError{
		BaseError: errors.BaseError{
			Kind:    errors.ErrKindValidation,
			Message: message,
		},
	}
}

// Domain-specific error constructors
func NewInvalidFirstNameError() error {
	return NewAuthorError("first name must be at least 3 characters long")
}

func NewInvalidLastNameError() error {
	return NewAuthorError("last name cannot be empty")
}

func NewInvalidProfileImageError() error {
	return NewAuthorError("invalid profile image URL format")
}

func NewInvalidSlugError() error {
	return NewAuthorError("slug must be at least 8 characters long")
}

func NewAuthorNotFoundError(id string) error {
	return errors.NewNotFoundError("author", id)
}

func NewAuthorAlreadyExistsError(slug string) error {
	return NewAuthorError("author with slug '" + slug + "' already exists")
}
