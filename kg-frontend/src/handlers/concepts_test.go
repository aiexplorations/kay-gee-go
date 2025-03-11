package handlers

import (
	"bytes"
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

// Mock Neo4j result
type MockNeo4jResult struct {
	mock.Mock
}

func (m *MockNeo4jResult) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockNeo4jResult) Record() neo4j.Record {
	args := m.Called()
	return args.Get(0).(neo4j.Record)
}

func (m *MockNeo4jResult) Err() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockNeo4jResult) Consume() (neo4j.ResultSummary, error) {
	args := m.Called()
	return args.Get(0).(neo4j.ResultSummary), args.Error(1)
}

// Mock Neo4j record
type MockNeo4jRecord struct {
	mock.Mock
}

func (m *MockNeo4jRecord) Keys() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockNeo4jRecord) Get(key string) (interface{}, bool) {
	args := m.Called(key)
	return args.Get(0), args.Bool(1)
}

func (m *MockNeo4jRecord) Values() []interface{} {
	args := m.Called()
	return args.Get(0).([]interface{})
}

func TestSearchConcepts(t *testing.T) {
	// Set up Gin in test mode
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create mock Neo4j driver and session
	mockDriver := new(MockNeo4jDriver)
	mockSession := new(MockNeo4jSession)
	mockResult := new(MockNeo4jResult)
	mockRecord := new(MockNeo4jRecord)

	// Set up expectations
	mockDriver.On("NewSession", mock.Anything).Return(mockSession)
	mockSession.On("Close").Return(nil)

	// Mock search results
	mockSession.On("ReadTransaction", mock.Anything).Return([]models.Node{
		{ID: "1", Name: "Artificial Intelligence"},
		{ID: "2", Name: "Machine Learning"},
	}, nil)

	// Set up the route
	router.GET("/api/concepts/search", SearchConcepts(mockDriver))

	// Create a test request
	req, _ := http.NewRequest("GET", "/api/concepts/search?q=Intelligence", nil)
	resp := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(resp, req)

	// Check the response
	assert.Equal(t, http.StatusOK, resp.Code)

	// Parse the response body
	var concepts []models.Node
	err := json.Unmarshal(resp.Body.Bytes(), &concepts)
	assert.NoError(t, err)

	// Check the response data
	assert.Equal(t, 2, len(concepts))
	assert.Equal(t, "Artificial Intelligence", concepts[0].Name)
	assert.Equal(t, "Machine Learning", concepts[1].Name)

	// Verify expectations
	mockDriver.AssertExpectations(t)
	mockSession.AssertExpectations(t)
}

func TestCreateRelationship(t *testing.T) {
	// Set up Gin in test mode
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create mock Neo4j driver and session
	mockDriver := new(MockNeo4jDriver)
	mockSession := new(MockNeo4jSession)

	// Set up expectations
	mockDriver.On("NewSession", mock.Anything).Return(mockSession)
	mockSession.On("Close").Return(nil)

	// Mock relationship creation
	mockSession.On("WriteTransaction", mock.Anything).Return(models.Link{
		Source: "1",
		Target: "2",
		Type:   "RELATES_TO",
	}, nil)

	// Set up the route
	router.POST("/api/relationships", CreateRelationship(mockDriver))

	// Create test request body
	relationship := models.RelationshipCreate{
		Source: "1",
		Target: "2",
		Type:   "RELATES_TO",
	}
	jsonData, _ := json.Marshal(relationship)

	// Create a test request
	req, _ := http.NewRequest("POST", "/api/relationships", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(resp, req)

	// Check the response
	assert.Equal(t, http.StatusOK, resp.Code)

	// Parse the response body
	var link models.Link
	err := json.Unmarshal(resp.Body.Bytes(), &link)
	assert.NoError(t, err)

	// Check the response data
	assert.Equal(t, "1", link.Source)
	assert.Equal(t, "2", link.Target)
	assert.Equal(t, "RELATES_TO", link.Type)

	// Verify expectations
	mockDriver.AssertExpectations(t)
	mockSession.AssertExpectations(t)
} 