package errors

import "fmt"

// HTTPError represents HTTP-specific errors
type HTTPError struct {
	BaseError
	StatusCode int
}

func NewHTTPError(statusCode int, message string) error {
	return &HTTPError{
		BaseError: BaseError{
			Kind:    ErrKindHTTP,
			Message: message,
		},
		StatusCode: statusCode,
	}
}

func NewBadRequestError(message string) error {
	return NewHTTPError(400, message)
}

func NewUnauthorizedError(message string) error {
	return NewHTTPError(401, message)
}

func NewForbiddenError(message string) error {
	return NewHTTPError(403, message)
}

func NewHTTPNotFoundError(resource string, id string) error {
	return NewHTTPError(404, fmt.Sprintf("%s not found: %s", resource, id))
}

func NewConflictError(message string) error {
	return NewHTTPError(409, message)
}

func NewInternalServerError(err error) error {
	return &HTTPError{
		BaseError: BaseError{
			Kind:    ErrKindHTTP,
			Message: "internal server error",
			Err:     err,
		},
		StatusCode: 500,
	}
}
