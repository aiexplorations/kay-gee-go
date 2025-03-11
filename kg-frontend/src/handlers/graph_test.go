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

// MockNeo4jSession is a mock implementation of the Neo4j session
type MockNeo4jSession struct {
	mock.Mock
}

// LastBookmark mocks the LastBookmark method
func (m *MockNeo4jSession) LastBookmark() string {
	args := m.Called()
	return args.String(0)
}

// BeginTransaction mocks the BeginTransaction method
func (m *MockNeo4jSession) BeginTransaction(configurers ...func(*neo4j.TransactionConfig)) (neo4j.Transaction, error) {
	args := m.Called(configurers)
	return args.Get(0).(neo4j.Transaction), args.Error(1)
}

// ReadTransaction mocks the ReadTransaction method
func (m *MockNeo4jSession) ReadTransaction(work neo4j.TransactionWork, configurers ...func(*neo4j.TransactionConfig)) (interface{}, error) {
	args := m.Called(work, configurers)
	return args.Get(0), args.Error(1)
}

// WriteTransaction mocks the WriteTransaction method
func (m *MockNeo4jSession) WriteTransaction(work neo4j.TransactionWork, configurers ...func(*neo4j.TransactionConfig)) (interface{}, error) {
	args := m.Called(work, configurers)
	return args.Get(0), args.Error(1)
}

// Run mocks the Run method
func (m *MockNeo4jSession) Run(cypher string, params map[string]interface{}, configurers ...func(*neo4j.TransactionConfig)) (neo4j.Result, error) {
	args := m.Called(cypher, params, configurers)
	return args.Get(0).(neo4j.Result), args.Error(1)
}

// Close mocks the Close method
func (m *MockNeo4jSession) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockNeo4jTransaction is a mock implementation of the Neo4j transaction
type MockNeo4jTransaction struct {
	mock.Mock
}

// Run mocks the Run method
func (m *MockNeo4jTransaction) Run(cypher string, params map[string]interface{}) (neo4j.Result, error) {
	args := m.Called(cypher, params)
	return args.Get(0).(neo4j.Result), args.Error(1)
}

// Commit mocks the Commit method
func (m *MockNeo4jTransaction) Commit() error {
	args := m.Called()
	return args.Error(0)
}

// Rollback mocks the Rollback method
func (m *MockNeo4jTransaction) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

// Close mocks the Close method
func (m *MockNeo4jTransaction) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockNeo4jResult is a mock implementation of the Neo4j result
type MockNeo4jResult struct {
	mock.Mock
}

// Keys mocks the Keys method
func (m *MockNeo4jResult) Keys() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

// Next mocks the Next method
func (m *MockNeo4jResult) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

// Err mocks the Err method
func (m *MockNeo4jResult) Err() error {
	args := m.Called()
	return args.Error(0)
}

// Record mocks the Record method
func (m *MockNeo4jResult) Record() neo4j.Record {
	args := m.Called()
	return args.Get(0).(neo4j.Record)
}

// Consume mocks the Consume method
func (m *MockNeo4jResult) Consume() (neo4j.ResultSummary, error) {
	args := m.Called()
	return args.Get(0).(neo4j.ResultSummary), args.Error(1)
}

// MockNeo4jRecord is a mock implementation of the Neo4j record
type MockNeo4jRecord struct {
	mock.Mock
}

// Get mocks the Get method
func (m *MockNeo4jRecord) Get(key string) (interface{}, bool) {
	args := m.Called(key)
	return args.Get(0), args.Bool(1)
}

// GetByIndex mocks the GetByIndex method
func (m *MockNeo4jRecord) GetByIndex(index int) interface{} {
	args := m.Called(index)
	return args.Get(0)
}

// Keys mocks the Keys method
func (m *MockNeo4jRecord) Keys() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

// Values mocks the Values method
func (m *MockNeo4jRecord) Values() []interface{} {
	args := m.Called()
	return args.Get(0).([]interface{})
}

func TestGetGraphData(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	
	// Create mock Neo4j driver and session
	mockDriver := new(MockNeo4jDriver)
	mockSession := new(MockNeo4jSession)
	
	// Setup mock behavior
	mockDriver.On("NewSession", mock.Anything).Return(mockSession)
	mockSession.On("Close").Return(nil)
	
	// Mock nodes
	nodes := []models.Node{
		{ID: "1", Name: "Node 1", Size: 5},
		{ID: "2", Name: "Node 2", Size: 7},
	}
	
	// Mock links
	links := []models.Link{
		{Source: "1", Target: "2", Type: "RELATES_TO"},
	}
	
	// Mock read transaction for nodes
	mockSession.On("ReadTransaction", mock.Anything, mock.Anything).Return(nodes, nil).Once()
	
	// Mock read transaction for links
	mockSession.On("ReadTransaction", mock.Anything, mock.Anything).Return(links, nil).Once()
	
	// Create a test context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	// Call the handler
	handler := GetGraphData(mockDriver)
	handler(c)
	
	// Verify
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Parse response
	var response models.GraphData
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	// Verify response
	assert.Equal(t, 2, len(response.Nodes))
	assert.Equal(t, 1, len(response.Links))
	
	// Verify mocks
	mockDriver.AssertExpectations(t)
	mockSession.AssertExpectations(t)
} 