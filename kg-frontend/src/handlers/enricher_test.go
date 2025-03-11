package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"kg-frontend/src/models"
)

func TestStartEnricher(t *testing.T) {
	// Save and restore the exec.Command function
	origExecCommand := exec.Command
	defer func() { exec.Command = origExecCommand }()
	exec.Command = mockExecCommand

	// Set up Gin in test mode
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Set up the route
	router.POST("/api/enricher/start", StartEnricher())

	// Create test request body
	params := models.EnricherParams{
		BatchSize:        10,
		Interval:         60,
		MaxRelationships: 100,
		Concurrency:      5,
	}
	jsonData, _ := json.Marshal(params)

	// Create a test request
	req, _ := http.NewRequest("POST", "/api/enricher/start", bytes.NewBuffer(jsonData))
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
	assert.Equal(t, "Enricher started successfully", response.Message)
}

func TestStopEnricher(t *testing.T) {
	// Save and restore the exec.Command function
	origExecCommand := exec.Command
	defer func() { exec.Command = origExecCommand }()
	exec.Command = mockExecCommand

	// Set up Gin in test mode
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Set up the route
	router.POST("/api/enricher/stop", StopEnricher())

	// Create a test request
	req, _ := http.NewRequest("POST", "/api/enricher/stop", nil)
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
	assert.Equal(t, "Enricher stopped successfully", response.Message)
} 