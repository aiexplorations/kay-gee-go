package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"kg-frontend/src/models"
)

// Mock exec.Command
func mockExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// TestHelperProcess isn't a real test. It's used as a helper process for mocking exec.Command
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// Simulate successful command execution
	os.Exit(0)
}

func TestStartBuilder(t *testing.T) {
	// Save and restore the exec.Command function
	origExecCommand := exec.Command
	defer func() { exec.Command = origExecCommand }()
	exec.Command = mockExecCommand

	// Set up Gin in test mode
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Set up the route
	router.POST("/api/builder/start", StartBuilder())

	// Create test request body
	params := models.BuilderParams{
		SeedConcept:        "Artificial Intelligence",
		MaxNodes:           100,
		Timeout:            30,
		RandomRelationships: 50,
		Concurrency:        5,
	}
	jsonData, _ := json.Marshal(params)

	// Create a test request
	req, _ := http.NewRequest("POST", "/api/builder/start", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(resp, req)

	// Check the response
	assert.Equal(t, http.StatusOK, resp.Code)

	// Parse the response body
	var response models.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check the response data
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Builder started successfully", response.Message)
}

func TestStopBuilder(t *testing.T) {
	// Save and restore the exec.Command function
	origExecCommand := exec.Command
	defer func() { exec.Command = origExecCommand }()
	exec.Command = mockExecCommand

	// Set up Gin in test mode
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Set up the route
	router.POST("/api/builder/stop", StopBuilder())

	// Create a test request
	req, _ := http.NewRequest("POST", "/api/builder/stop", nil)
	resp := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(resp, req)

	// Check the response
	assert.Equal(t, http.StatusOK, resp.Code)

	// Parse the response body
	var response models.Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check the response data
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Builder stopped successfully", response.Message)
} 