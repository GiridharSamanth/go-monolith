package errors

import "fmt"

// DatabaseError represents database-specific errors
type DatabaseError struct {
	BaseError
	Operation string
}

func NewDatabaseError(operation string, err error) error {
	return &DatabaseError{
		BaseError: BaseError{
			Kind:    ErrKindDatabase,
			Message: fmt.Sprintf("database error during %s", operation),
			Err:     err,
		},
		Operation: operation,
	}
}

func NewConnectionError(err error) error {
	return &DatabaseError{
		BaseError: BaseError{
			Kind:    ErrKindDatabase,
			Message: "failed to connect to database",
			Err:     err,
		},
		Operation: "connect",
	}
}

func NewQueryError(operation string, err error) error {
	return &DatabaseError{
		BaseError: BaseError{
			Kind:    ErrKindDatabase,
			Message: fmt.Sprintf("query error during %s", operation),
			Err:     err,
		},
		Operation: operation,
	}
}

func NewTransactionError(err error) error {
	return &DatabaseError{
		BaseError: BaseError{
			Kind:    ErrKindDatabase,
			Message: "transaction error",
			Err:     err,
		},
		Operation: "transaction",
	}
}
