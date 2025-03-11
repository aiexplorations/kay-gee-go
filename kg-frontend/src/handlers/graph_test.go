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

// Mock Neo4j driver and session
type MockNeo4jDriver struct {
	mock.Mock
}

func (m *MockNeo4jDriver) NewSession(config neo4j.SessionConfig) neo4j.Session {
	args := m.Called(config)
	return args.Get(0).(neo4j.Session)
}

func (m *MockNeo4jDriver) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockNeo4jDriver) VerifyConnectivity() error {
	args := m.Called()
	return args.Error(0)
}

type MockNeo4jSession struct {
	mock.Mock
}

func (m *MockNeo4jSession) Run(cypher string, params map[string]interface{}) (neo4j.Result, error) {
	args := m.Called(cypher, params)
	return args.Get(0).(neo4j.Result), args.Error(1)
}

func (m *MockNeo4jSession) ReadTransaction(work neo4j.TransactionWork) (interface{}, error) {
	args := m.Called(work)
	return args.Get(0), args.Error(1)
}

func (m *MockNeo4jSession) WriteTransaction(work neo4j.TransactionWork) (interface{}, error) {
	args := m.Called(work)
	return args.Get(0), args.Error(1)
}

func (m *MockNeo4jSession) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestGetGraphData(t *testing.T) {
	// Set up Gin in test mode
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create mock Neo4j driver and session
	mockDriver := new(MockNeo4jDriver)
	mockSession := new(MockNeo4jSession)

	// Set up expectations
	mockDriver.On("NewSession", mock.Anything).Return(mockSession)
	mockSession.On("Close").Return(nil)

	// Mock nodes data
	nodes := []models.Node{
		{ID: "1", Name: "Artificial Intelligence", Size: 5},
		{ID: "2", Name: "Machine Learning", Size: 5},
	}

	// Mock links data
	links := []models.Link{
		{Source: "1", Target: "2", Type: "RELATES_TO"},
	}

	// Set up mock session to return nodes and links
	mockSession.On("ReadTransaction", mock.Anything).Return(nodes, nil).Once()
	mockSession.On("ReadTransaction", mock.Anything).Return(links, nil).Once()

	// Set up the route
	router.GET("/api/graph", GetGraphData(mockDriver))

	// Create a test request
	req, _ := http.NewRequest("GET", "/api/graph", nil)
	resp := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(resp, req)

	// Check the response
	assert.Equal(t, http.StatusOK, resp.Code)

	// Parse the response body
	var graphData models.GraphData
	err := json.Unmarshal(resp.Body.Bytes(), &graphData)
	assert.NoError(t, err)

	// Check the response data
	assert.Equal(t, 2, len(graphData.Nodes))
	assert.Equal(t, 1, len(graphData.Links))
	assert.Equal(t, "Artificial Intelligence", graphData.Nodes[0].Name)
	assert.Equal(t, "Machine Learning", graphData.Nodes[1].Name)
	assert.Equal(t, "1", graphData.Links[0].Source)
	assert.Equal(t, "2", graphData.Links[0].Target)
	assert.Equal(t, "RELATES_TO", graphData.Links[0].Type)

	// Verify expectations
	mockDriver.AssertExpectations(t)
	mockSession.AssertExpectations(t)
} 