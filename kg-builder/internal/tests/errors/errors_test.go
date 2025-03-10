package errors_test

import (
	"errors"
	"testing"
	"time"

	apperrors "kg-builder/internal/errors"
)

func TestAppError(t *testing.T) {
	// Test creating a new database error
	dbErr := apperrors.NewDatabaseError(errors.New("connection failed"), "failed to connect to database")
	if dbErr.Type != apperrors.ErrorTypeDatabase {
		t.Errorf("Expected error type %s, got %s", apperrors.ErrorTypeDatabase, dbErr.Type)
	}
	if dbErr.Message != "failed to connect to database" {
		t.Errorf("Expected message 'failed to connect to database', got '%s'", dbErr.Message)
	}
	if !dbErr.Retryable {
		t.Error("Expected database error to be retryable by default")
	}

	// Test creating a new LLM error
	llmErr := apperrors.NewLLMError(errors.New("request failed"), "failed to call LLM service")
	if llmErr.Type != apperrors.ErrorTypeLLM {
		t.Errorf("Expected error type %s, got %s", apperrors.ErrorTypeLLM, llmErr.Type)
	}
	if llmErr.Message != "failed to call LLM service" {
		t.Errorf("Expected message 'failed to call LLM service', got '%s'", llmErr.Message)
	}
	if !llmErr.Retryable {
		t.Error("Expected LLM error to be retryable by default")
	}

	// Test creating a new graph error
	graphErr := apperrors.NewGraphError(errors.New("processing failed"), "failed to process concept")
	if graphErr.Type != apperrors.ErrorTypeGraph {
		t.Errorf("Expected error type %s, got %s", apperrors.ErrorTypeGraph, graphErr.Type)
	}
	if graphErr.Message != "failed to process concept" {
		t.Errorf("Expected message 'failed to process concept', got '%s'", graphErr.Message)
	}
	if graphErr.Retryable {
		t.Error("Expected graph error to be non-retryable by default")
	}

	// Test creating a new config error
	configErr := apperrors.NewConfigError(errors.New("invalid config"), "failed to load configuration")
	if configErr.Type != apperrors.ErrorTypeConfig {
		t.Errorf("Expected error type %s, got %s", apperrors.ErrorTypeConfig, configErr.Type)
	}
	if configErr.Message != "failed to load configuration" {
		t.Errorf("Expected message 'failed to load configuration', got '%s'", configErr.Message)
	}
	if configErr.Retryable {
		t.Error("Expected config error to be non-retryable by default")
	}

	// Test Error() method
	expectedErrStr := "[database] failed to connect to database: connection failed"
	if dbErr.Error() != expectedErrStr {
		t.Errorf("Expected error string '%s', got '%s'", expectedErrStr, dbErr.Error())
	}

	// Test Unwrap() method
	unwrappedErr := dbErr.Unwrap()
	if unwrappedErr.Error() != "connection failed" {
		t.Errorf("Expected unwrapped error 'connection failed', got '%s'", unwrappedErr.Error())
	}

	// Test WithRetryable() method
	nonRetryableDbErr := dbErr.WithRetryable(false)
	if nonRetryableDbErr.Retryable {
		t.Error("Expected error to be non-retryable after WithRetryable(false)")
	}

	// Test IsRetryable() function
	if !apperrors.IsRetryable(dbErr) {
		t.Errorf("Expected IsRetryable to return true for retryable error. Error: %v, Type: %T, Retryable: %v", 
			dbErr, dbErr, dbErr.Retryable)
	}
	if apperrors.IsRetryable(nonRetryableDbErr) {
		t.Errorf("Expected IsRetryable to return false for non-retryable error. Error: %v, Type: %T, Retryable: %v", 
			nonRetryableDbErr, nonRetryableDbErr, nonRetryableDbErr.Retryable)
	}
	if apperrors.IsRetryable(errors.New("standard error")) {
		t.Error("Expected IsRetryable to return false for standard error")
	}
}

func TestRetryWithBackoff(t *testing.T) {
	// Test successful execution without retries
	callCount := 0
	err := apperrors.RetryWithBackoff(3, 1*time.Millisecond, 10*time.Millisecond, func() error {
		callCount++
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if callCount != 1 {
		t.Errorf("Expected function to be called once, got %d", callCount)
	}

	// Test with non-retryable error
	callCount = 0
	nonRetryableErr := apperrors.NewConfigError(errors.New("config error"), "config error").WithRetryable(false)
	err = apperrors.RetryWithBackoff(3, 1*time.Millisecond, 10*time.Millisecond, func() error {
		callCount++
		return nonRetryableErr
	})
	if err != nonRetryableErr {
		t.Errorf("Expected error %v, got %v", nonRetryableErr, err)
	}
	if callCount != 1 {
		t.Errorf("Expected function to be called once for non-retryable error, got %d", callCount)
	}

	// Test with retryable error that eventually succeeds
	callCount = 0
	err = apperrors.RetryWithBackoff(3, 1*time.Millisecond, 10*time.Millisecond, func() error {
		callCount++
		if callCount < 3 {
			return apperrors.NewDatabaseError(errors.New("temp error"), "temporary error")
		}
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error after retries, got %v", err)
	}
	if callCount != 3 {
		t.Errorf("Expected function to be called 3 times, got %d", callCount)
	}

	// Test with retryable error that never succeeds
	callCount = 0
	err = apperrors.RetryWithBackoff(3, 1*time.Millisecond, 10*time.Millisecond, func() error {
		callCount++
		return apperrors.NewDatabaseError(errors.New("persistent error"), "persistent error")
	})
	if err == nil {
		t.Error("Expected error after max retries, got nil")
	}
	if callCount != 3 {
		t.Errorf("Expected function to be called 3 times for max retries, got %d", callCount)
	}
} 