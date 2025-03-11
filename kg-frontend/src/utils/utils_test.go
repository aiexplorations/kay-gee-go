package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetEnv(t *testing.T) {
	// Test with existing environment variable
	os.Setenv("TEST_ENV_VAR", "test_value")
	defer os.Unsetenv("TEST_ENV_VAR")

	value := GetEnv("TEST_ENV_VAR", "default_value")
	if value != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", value)
	}

	// Test with non-existing environment variable
	value = GetEnv("NON_EXISTING_VAR", "default_value")
	if value != "default_value" {
		t.Errorf("Expected 'default_value', got '%s'", value)
	}
}

func TestGetScriptPath(t *testing.T) {
	// Create a temporary file for testing
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test_script.sh")
	
	// Create the file
	f, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	f.Close()

	// Test with existing file
	scriptName := filepath.Base(tempFile)
	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	path := GetScriptPath(scriptName)
	absPath, _ := filepath.Abs(scriptName)
	if path != absPath {
		t.Errorf("Expected '%s', got '%s'", absPath, path)
	}

	// Test with non-existing file
	path = GetScriptPath("non_existing_script.sh")
	if path != "non_existing_script.sh" {
		t.Errorf("Expected 'non_existing_script.sh', got '%s'", path)
	}
} 