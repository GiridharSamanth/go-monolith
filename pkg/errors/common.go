package errors

import "fmt"

// ErrKind represents different types of errors that can occur
type ErrKind int

const (
	ErrKindNotFound ErrKind = iota
	ErrKindValidation
	ErrKindTransient
	ErrKindUnexpected
	ErrKindDatabase
	ErrKindHTTP
	ErrKindSession
)

// BaseError represents a common error type that can be used across the application
type BaseError struct {
	Kind    ErrKind
	Message string
	Err     error
}

func (e *BaseError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *BaseError) Unwrap() error {
	return e.Err
}

// Common error constructors
func NewNotFoundError(resource string, id string) error {
	return &BaseError{
		Kind:    ErrKindNotFound,
		Message: fmt.Sprintf("%s not found: %s", resource, id),
	}
}

func NewTransientError(err error) error {
	return &BaseError{
		Kind:    ErrKindTransient,
		Message: "temporary error occurred",
		Err:     err,
	}
}

func NewValidationError(msg string) error {
	return &BaseError{
		Kind:    ErrKindValidation,
		Message: msg,
	}
}

func NewUnexpectedError(err error) error {
	return &BaseError{
		Kind:    ErrKindUnexpected,
		Message: "unexpected error occurred",
		Err:     err,
	}
}
