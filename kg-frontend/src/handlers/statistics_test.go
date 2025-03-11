package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"kg-frontend/src/models"
)

func TestGetStatistics(t *testing.T) {
	// Set up Gin in test mode
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create mock Neo4j driver and session
	mockDriver := new(MockNeo4jDriver)
	mockSession := new(MockNeo4jSession)

	// Set up expectations
	mockDriver.On("NewSession", mock.Anything).Return(mockSession)
	mockSession.On("Close").Return(nil)

	// Mock statistics data
	mockSession.On("ReadTransaction", mock.Anything).Return(int64(42), nil).Once()
	mockSession.On("ReadTransaction", mock.Anything).Return(int64(123), nil).Once()

	// Set up the route
	router.GET("/api/statistics", GetStatistics(mockDriver))

	// Create a test request
	req, _ := http.NewRequest("GET", "/api/statistics", nil)
	resp := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(resp, req)

	// Check the response
	assert.Equal(t, http.StatusOK, resp.Code)

	// Parse the response body
	var stats models.Statistics
	err := json.Unmarshal(resp.Body.Bytes(), &stats)
	assert.NoError(t, err)

	// Check the response data
	assert.Equal(t, 42, stats.ConceptCount)
	assert.Equal(t, 123, stats.RelationshipCount)

	// Verify expectations
	mockDriver.AssertExpectations(t)
	mockSession.AssertExpectations(t)
} 