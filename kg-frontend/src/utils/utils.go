package utils

import (
	"os"
	"path/filepath"
)

// GetScriptPath gets the absolute path to a script
func GetScriptPath(scriptName string) string {
	// Check if the script exists in the current directory
	if _, err := os.Stat(scriptName); err == nil {
		absPath, _ := filepath.Abs(scriptName)
		return absPath
	}

	// Check if the script exists in the /app directory (for Docker)
	appPath := filepath.Join("/app", scriptName)
	if _, err := os.Stat(appPath); err == nil {
		return appPath
	}

	// Return the script name as is (rely on PATH)
	return scriptName
}

// GetEnv gets an environment variable or returns a default value
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
} 