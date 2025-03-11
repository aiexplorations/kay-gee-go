package utils

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCommandRunner tests the CommandRunner struct
func TestCommandRunner(t *testing.T) {
	// Create a new CommandRunner
	runner := &CommandRunner{}

	// Skip this test if we're in CI or want to avoid actual command execution
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping test that executes real commands in CI environment")
	}

	// Test running a simple command
	output, err := runner.RunCommand("echo", "test")
	assert.NoError(t, err)
	assert.Equal(t, "test\n", string(output))

	// Test running a command that doesn't exist
	_, err = runner.RunCommand("nonexistentcommand")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "executable file not found")

	// Test running a command that fails
	_, err = runner.RunCommand("ls", "/nonexistentdirectory")
	assert.Error(t, err)
}

// TestCommandRunnerWithMock tests the CommandRunner using a mock exec.Command
func TestCommandRunnerWithMock(t *testing.T) {
	// Save the original exec.Command function
	originalExecCommand := execCommand
	defer func() { execCommand = originalExecCommand }()

	// Mock the exec.Command function
	execCommand = func(name string, args ...string) *exec.Cmd {
		// Create a mock command that always succeeds
		if name == "echo" && len(args) > 0 && args[0] == "success" {
			cmd := exec.Command("echo", "mocked success")
			return cmd
		}
		
		// Create a mock command that always fails
		if name == "echo" && len(args) > 0 && args[0] == "failure" {
			cmd := exec.Command("nonexistentcommand")
			return cmd
		}
		
		// Default to the real command
		return originalExecCommand(name, args...)
	}

	// Create a new CommandRunner
	runner := &CommandRunner{}

	// Test running a mocked successful command
	output, err := runner.RunCommand("echo", "success")
	assert.NoError(t, err)
	assert.Equal(t, "mocked success\n", string(output))

	// Test running a mocked failing command
	_, err = runner.RunCommand("echo", "failure")
	assert.Error(t, err)
}

// TestCommandRunnerWithEnv tests the CommandRunner with environment variables
func TestCommandRunnerWithEnv(t *testing.T) {
	// Skip this test if we're in CI or want to avoid actual command execution
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping test that executes real commands in CI environment")
	}

	// Create a new CommandRunner
	runner := &CommandRunner{}

	// Set an environment variable
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	// Test running a command that uses the environment variable
	output, err := runner.RunCommand("sh", "-c", "echo $TEST_VAR")
	assert.NoError(t, err)
	assert.Equal(t, "test_value\n", string(output))
} 