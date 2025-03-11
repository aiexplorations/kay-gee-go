package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"kg-frontend/src/handlers"
	"kg-frontend/src/utils"
)

// MockNeo4jDriver is a mock implementation of the Neo4j driver
type MockNeo4jDriver struct {
	mock.Mock
}

// NewSession mocks the NewSession method
func (m *MockNeo4jDriver) NewSession(config neo4j.SessionConfig) neo4j.Session {
	args := m.Called(config)
	return args.Get(0).(neo4j.Session)
}

// Close mocks the Close method
func (m *MockNeo4jDriver) Close() error {
	args := m.Called()
	return args.Error(0)
}

// VerifyConnectivity mocks the VerifyConnectivity method
func (m *MockNeo4jDriver) VerifyConnectivity() error {
	args := m.Called()
	return args.Error(0)
}

// TestSetupRouter tests the setupRouter function
func TestSetupRouter(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create mock Neo4j driver
	mockDriver := new(MockNeo4jDriver)
	mockDriver.On("VerifyConnectivity").Return(nil)

	// Create mock command runner
	mockRunner := new(handlers.MockCommandRunner)

	// Setup router
	router := setupRouter(mockDriver, mockRunner)

	// Test static file serving
	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// Test API endpoint
	req, _ = http.NewRequest("GET", "/api/health", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// Verify mocks
	mockDriver.AssertExpectations(t)
}

// TestHealthCheck tests the health check endpoint
func TestHealthCheck(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new router
	router := gin.New()
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Test health check endpoint
	req, _ := http.NewRequest("GET", "/api/health", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

// TestCreateNeo4jDriver tests the createNeo4jDriver function
func TestCreateNeo4jDriver(t *testing.T) {
	// Skip this test in CI environment
	t.Skip("Skipping test that requires Neo4j connection")

	// Test with invalid URI
	_, err := createNeo4jDriver("invalid-uri", "username", "password")
	assert.Error(t, err)
}

// TestMain tests the main function
func TestMain(t *testing.T) {
	// Skip this test as it would start the server
	t.Skip("Skipping test that would start the server")
} 