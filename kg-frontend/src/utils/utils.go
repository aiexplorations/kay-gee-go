package utils

import (
	"os"
	"path/filepath"
)

// GetScriptPath gets the absolute path to a script
func GetScriptPath(scriptName string) string {
	// First, check if the script exists in the current directory
	if _, err := os.Stat(scriptName); err == nil {
		absPath, _ := filepath.Abs(scriptName)
		return absPath
	}

	// Check if the script exists in the src directory
	srcPath := filepath.Join("src", scriptName)
	if _, err := os.Stat(srcPath); err == nil {
		absPath, _ := filepath.Abs(srcPath)
		return absPath
	}

	// Check if the script exists in the /app/src directory (for Docker)
	appPath := filepath.Join("/app/src", scriptName)
	if _, err := os.Stat(appPath); err == nil {
		return appPath
	}

	// Get the executable directory
	execDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err == nil {
		// Check if the script exists in the executable's directory
		execPath := filepath.Join(execDir, scriptName)
		if _, err := os.Stat(execPath); err == nil {
			return execPath
		}

		// Check if the script exists in the src subdirectory of the executable's directory
		execSrcPath := filepath.Join(execDir, "src", scriptName)
		if _, err := os.Stat(execSrcPath); err == nil {
			return execSrcPath
		}
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