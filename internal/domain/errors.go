package domain

import (
	"errors"
	"fmt"
)

// Sentinel errors for common cases
var (
	// ErrNotFound is returned when a resource is not found
	ErrNotFound = errors.New("resource not found")

	// ErrInvalidInput is returned when input validation fails
	ErrInvalidInput = errors.New("invalid input")

	// ErrInvalidStatusTransition is returned when a status transition is not allowed
	ErrInvalidStatusTransition = errors.New("invalid status transition")

	// ErrAlreadyExists is returned when a resource already exists
	ErrAlreadyExists = errors.New("resource already exists")

	// ErrDriverNotAvailable is returned when a driver is not available
	ErrDriverNotAvailable = errors.New("driver not available")

	// ErrInternal is returned for internal server errors
	ErrInternal = errors.New("internal server error")

	// ErrConflict is returned when there's a conflict with current state
	ErrConflict = errors.New("conflict with current state")

	// ErrTimeout is returned when operation times out
	ErrTimeout = errors.New("operation timeout")
)

// DomainError represents a domain-specific error with context
type DomainError struct {
	Op      string // operation that failed
	Code    string // error code for client identification
	Message string // human-readable message
	Err     error  // underlying error
}

// Error implements the error interface
func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Op, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Op, e.Message)
}

// Unwrap implements the unwrap interface for errors.Is and errors.As
func (e *DomainError) Unwrap() error {
	return e.Err
}

// NewDomainError creates a new domain error
func NewDomainError(op, code, message string, err error) *DomainError {
	return &DomainError{
		Op:      op,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
	Err     error
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("validation error on %s: %s: %v", e.Field, e.Message, e.Err)
	}
	return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

// Unwrap implements the unwrap interface
func (e *ValidationError) Unwrap() error {
	return e.Err
}

// Is checks if the error is ErrInvalidInput
func (e *ValidationError) Is(target error) bool {
	return target == ErrInvalidInput
}

// NotFoundError represents a resource not found error
type NotFoundError struct {
	Resource string
	ID       string
	Err      error
}

// Error implements the error interface
func (e *NotFoundError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s not found with id %s: %v", e.Resource, e.ID, e.Err)
	}
	return fmt.Sprintf("%s not found with id %s", e.Resource, e.ID)
}

// Unwrap implements the unwrap interface
func (e *NotFoundError) Unwrap() error {
	return e.Err
}

// Is checks if the error is ErrNotFound
func (e *NotFoundError) Is(target error) bool {
	return target == ErrNotFound
}

// ConflictError represents a state conflict error
type ConflictError struct {
	Resource     string
	CurrentState string
	RequestedOp  string
	Message      string
	Err          error
}

// Error implements the error interface
func (e *ConflictError) Error() string {
	msg := fmt.Sprintf("conflict: cannot perform %s on %s in state %s",
		e.RequestedOp, e.Resource, e.CurrentState)
	if e.Message != "" {
		msg += ": " + e.Message
	}
	if e.Err != nil {
		msg += fmt.Sprintf(": %v", e.Err)
	}
	return msg
}

// Unwrap implements the unwrap interface
func (e *ConflictError) Unwrap() error {
	return e.Err
}

// Is checks if the error is ErrConflict or ErrInvalidStatusTransition
func (e *ConflictError) Is(target error) bool {
	return target == ErrConflict || target == ErrInvalidStatusTransition
}
