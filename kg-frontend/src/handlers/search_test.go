package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"kg-frontend/src/models"
)

func TestSearchConcepts(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	
	// Create mock Neo4j driver and session
	mockDriver := new(MockNeo4jDriver)
	mockSession := new(MockNeo4jSession)
	
	// Setup mock behavior
	mockDriver.On("NewSession", mock.Anything).Return(mockSession)
	mockSession.On("Close").Return(nil)
	
	// Mock search results
	searchResults := []models.SearchResult{
		{ID: "1", Name: "Artificial Intelligence", Type: "Concept"},
		{ID: "2", Name: "Machine Learning", Type: "Concept"},
	}
	
	// Mock read transaction for search
	mockSession.On("ReadTransaction", mock.Anything, mock.Anything).Return(searchResults, nil).Once()
	
	// Create a test context with query parameter
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	// Add query parameter
	c.Request, _ = http.NewRequest("GET", "/api/search?q=artificial", nil)
	c.Request.URL.RawQuery = "q=artificial"
	
	// Call the handler
	handler := SearchConcepts(mockDriver)
	handler(c)
	
	// Verify
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Parse response
	var response []models.SearchResult
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	// Verify response
	assert.Equal(t, 2, len(response))
	assert.Equal(t, "Artificial Intelligence", response[0].Name)
	assert.Equal(t, "Machine Learning", response[1].Name)
	
	// Verify mocks
	mockDriver.AssertExpectations(t)
	mockSession.AssertExpectations(t)
}

func TestSearchConceptsEmptyQuery(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	
	// Create mock Neo4j driver
	mockDriver := new(MockNeo4jDriver)
	
	// Create a test context with empty query parameter
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	// Add empty query parameter
	c.Request, _ = http.NewRequest("GET", "/api/search?q=", nil)
	c.Request.URL.RawQuery = "q="
	
	// Call the handler
	handler := SearchConcepts(mockDriver)
	handler(c)
	
	// Verify
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	// Parse response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	// Verify response
	assert.Equal(t, "Query parameter 'q' is required", response["error"])
}

func TestSearchConceptsNoQuery(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	
	// Create mock Neo4j driver
	mockDriver := new(MockNeo4jDriver)
	
	// Create a test context with no query parameter
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	// Add no query parameter
	c.Request, _ = http.NewRequest("GET", "/api/search", nil)
	
	// Call the handler
	handler := SearchConcepts(mockDriver)
	handler(c)
	
	// Verify
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	// Parse response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	// Verify response
	assert.Equal(t, "Query parameter 'q' is required", response["error"])
}

func TestSearchConceptsError(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	
	// Create mock Neo4j driver and session
	mockDriver := new(MockNeo4jDriver)
	mockSession := new(MockNeo4jSession)
	
	// Setup mock behavior
	mockDriver.On("NewSession", mock.Anything).Return(mockSession)
	mockSession.On("Close").Return(nil)
	
	// Mock read transaction error
	mockSession.On("ReadTransaction", mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
	
	// Create a test context with query parameter
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	// Add query parameter
	c.Request, _ = http.NewRequest("GET", "/api/search?q=artificial", nil)
	c.Request.URL.RawQuery = "q=artificial"
	
	// Call the handler
	handler := SearchConcepts(mockDriver)
	handler(c)
	
	// Verify
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	
	// Parse response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	// Verify response
	assert.Contains(t, response["error"], "Failed to search concepts")
	
	// Verify mocks
	mockDriver.AssertExpectations(t)
	mockSession.AssertExpectations(t)
} 