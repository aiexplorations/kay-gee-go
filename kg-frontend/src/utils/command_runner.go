package utils

import (
	"os/exec"
)

// CommandRunnerInterface defines the interface for running commands
type CommandRunnerInterface interface {
	RunCommand(name string, args ...string) ([]byte, error)
}

// CommandRunner implements the CommandRunnerInterface
type CommandRunner struct{}

// Variable to allow mocking in tests
var execCommand = exec.Command

// RunCommand executes a command with the given name and arguments
// and returns the combined output and any error
func (c *CommandRunner) RunCommand(name string, args ...string) ([]byte, error) {
	cmd := execCommand(name, args...)
	return cmd.CombinedOutput()
} 