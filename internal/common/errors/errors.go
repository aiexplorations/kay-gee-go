package errors

import (
	"fmt"
)

// ErrorType represents the type of error
type ErrorType string

const (
	// ErrorTypeDatabase represents a database error
	ErrorTypeDatabase ErrorType = "database"
	
	// ErrorTypeLLM represents an LLM service error
	ErrorTypeLLM ErrorType = "llm"
	
	// ErrorTypeConfig represents a configuration error
	ErrorTypeConfig ErrorType = "config"
	
	// ErrorTypeInternal represents an internal error
	ErrorTypeInternal ErrorType = "internal"
	
	// ErrorTypeTimeout represents a timeout error
	ErrorTypeTimeout ErrorType = "timeout"
)

// AppError represents an application error
type AppError struct {
	Type    ErrorType
	Message string
	Cause   error
}

// Error returns the error message
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (cause: %s)", e.Type, e.Message, e.Cause.Error())
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Cause
}

// NewDatabaseError creates a new database error
func NewDatabaseError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeDatabase,
		Message: message,
		Cause:   cause,
	}
}

// NewLLMError creates a new LLM service error
func NewLLMError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeLLM,
		Message: message,
		Cause:   cause,
	}
}

// NewConfigError creates a new configuration error
func NewConfigError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeConfig,
		Message: message,
		Cause:   cause,
	}
}

// NewInternalError creates a new internal error
func NewInternalError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeInternal,
		Message: message,
		Cause:   cause,
	}
}

// NewTimeoutError creates a new timeout error
func NewTimeoutError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeTimeout,
		Message: message,
		Cause:   cause,
	}
}

// IsType checks if the error is of the specified type
func IsType(err error, errorType ErrorType) bool {
	if err == nil {
		return false
	}
	
	if e, ok := err.(*AppError); ok {
		return e.Type == errorType
	}
	
	return false
} 