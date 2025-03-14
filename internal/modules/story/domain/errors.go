package domain

import (
	"go-monolith/pkg/errors"
)

// StoryError represents story-specific domain errors
type StoryError struct {
	errors.BaseError
}

func NewStoryError(message string, err error) error {
	return &StoryError{
		BaseError: errors.BaseError{
			Kind:    errors.ErrKindValidation,
			Message: message,
			Err:     err,
		},
	}
}

// Domain-specific error constructors
func NewStoryNotFoundError(id string) error {
	return errors.NewNotFoundError("story", id)
}

func NewStoryTransientError(err error) error {
	return errors.NewTransientError(err)
}

func NewStoryValidationError(msg string) error {
	return errors.NewValidationError(msg)
}

func NewStoryUnexpectedError(err error) error {
	return errors.NewUnexpectedError(err)
}

// Additional story-specific validation errors
func NewInvalidTitleError() error {
	return NewStoryError("title must be at least 3 characters long", nil)
}

func NewInvalidContentError() error {
	return NewStoryError("content cannot be empty", nil)
}

func NewInvalidAuthorError() error {
	return NewStoryError("author is required", nil)
}

func NewInvalidStatusError() error {
	return NewStoryError("invalid story status", nil)
}
