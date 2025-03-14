package errors

// SessionError represents session-specific errors
type SessionError struct {
	BaseError
}

func NewSessionError(message string, err error) error {
	return &SessionError{
		BaseError: BaseError{
			Kind:    ErrKindSession,
			Message: message,
			Err:     err,
		},
	}
}

func NewSessionExpiredError() error {
	return NewSessionError("session has expired", nil)
}

func NewInvalidSessionError() error {
	return NewSessionError("invalid session", nil)
}

func NewSessionCreationError(err error) error {
	return NewSessionError("failed to create session", err)
}

func NewSessionDeletionError(err error) error {
	return NewSessionError("failed to delete session", err)
}
