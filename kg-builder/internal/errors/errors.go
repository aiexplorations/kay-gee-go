package errors

import (
	"errors"
	"fmt"
	"time"
)

// Standard error types
var (
	ErrNotFound      = errors.New("not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrTimeout       = errors.New("operation timed out")
	ErrUnavailable   = errors.New("service unavailable")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrInternal      = errors.New("internal error")
)

// ErrorType represents the type of an error
type ErrorType string

const (
	// ErrorTypeDatabase represents database errors
	ErrorTypeDatabase ErrorType = "database"
	// ErrorTypeLLM represents LLM service errors
	ErrorTypeLLM ErrorType = "llm"
	// ErrorTypeGraph represents graph building errors
	ErrorTypeGraph ErrorType = "graph"
	// ErrorTypeConfig represents configuration errors
	ErrorTypeConfig ErrorType = "config"
)

// AppError is a custom error type that includes additional context
type AppError struct {
	Type      ErrorType
	Err       error
	Message   string
	Timestamp time.Time
	Retryable bool
}

// Error implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Err)
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// Is reports whether the target error is of the same type as the receiver
func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok {
		return errors.Is(e.Err, target)
	}
	return e.Type == t.Type
}

// NewDatabaseError creates a new database error
func NewDatabaseError(err error, message string) *AppError {
	return &AppError{
		Type:      ErrorTypeDatabase,
		Err:       err,
		Message:   message,
		Timestamp: time.Now(),
		Retryable: true, // Database errors are often retryable
	}
}

// NewLLMError creates a new LLM service error
func NewLLMError(err error, message string) *AppError {
	return &AppError{
		Type:      ErrorTypeLLM,
		Err:       err,
		Message:   message,
		Timestamp: time.Now(),
		Retryable: true, // LLM errors are often retryable
	}
}

// NewGraphError creates a new graph building error
func NewGraphError(err error, message string) *AppError {
	return &AppError{
		Type:      ErrorTypeGraph,
		Err:       err,
		Message:   message,
		Timestamp: time.Now(),
		Retryable: false, // Graph errors are often not retryable
	}
}

// NewConfigError creates a new configuration error
func NewConfigError(err error, message string) *AppError {
	return &AppError{
		Type:      ErrorTypeConfig,
		Err:       err,
		Message:   message,
		Timestamp: time.Now(),
		Retryable: false, // Config errors are not retryable
	}
}

// WithRetryable sets the retryable flag on an AppError
func (e *AppError) WithRetryable(retryable bool) *AppError {
	// Create a new AppError with the same values but different retryable flag
	return &AppError{
		Type:      e.Type,
		Err:       e.Err,
		Message:   e.Message,
		Timestamp: e.Timestamp,
		Retryable: retryable,
	}
}

// IsRetryable returns true if the error is retryable
func IsRetryable(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Retryable
	}
	return false
}

// RetryWithBackoff executes the given function with exponential backoff
func RetryWithBackoff(maxRetries int, initialBackoff time.Duration, maxBackoff time.Duration, fn func() error) error {
	var err error
	backoff := initialBackoff

	for i := 0; i < maxRetries; i++ {
		err = fn()
		if err == nil {
			return nil
		}

		if !IsRetryable(err) {
			return err
		}

		if i < maxRetries-1 {
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * 1.5)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}
	}

	return fmt.Errorf("failed after %d retries: %w", maxRetries, err)
} 