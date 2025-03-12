package cleanup_test

import (
	"testing"
	"time"
	"errors"
	"net/url"

	"kg-builder/internal/cleanup"
	"kg-builder/internal/config"
	"kg-builder/internal/graph"
	"kg-builder/internal/models"
	kgneo4j "kg-builder/internal/neo4j"

	neo4jdriver "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// Define a custom error for session expired
var errSessionExpired = errors.New("session expired")

// MockNeo4jDriver implements the Neo4j driver interface for testing
type MockNeo4jDriver struct {
	// Track calls to methods
	CleanupOrphanRelationshipsCalled bool
	CleanupOrphanNodesCalled         bool
	
	// Mock return values
	OrphanRelationshipsCount int
	OrphanNodesCount         int
	
	// Mock error responses
	ShouldReturnError bool
	
	// Track number of calls
	CleanupOrphanRelationshipsCallCount int
	CleanupOrphanNodesCallCount         int
}

// Implement Neo4j driver interface methods
func (m *MockNeo4jDriver) Close() error {
	return nil
}

func (m *MockNeo4jDriver) NewSession(config neo4jdriver.SessionConfig) neo4jdriver.Session {
	return &MockNeo4jSession{driver: m}
}

func (m *MockNeo4jDriver) Session(accessMode neo4jdriver.AccessMode, bookmarks ...string) (neo4jdriver.Session, error) {
	return &MockNeo4jSession{driver: m}, nil
}

func (m *MockNeo4jDriver) Target() url.URL {
	u, _ := url.Parse("bolt://localhost:7687")
	return *u
}

func (m *MockNeo4jDriver) VerifyConnectivity() error {
	return nil
}

// MockNeo4jSession implements the Neo4j session interface for testing
type MockNeo4jSession struct {
	driver *MockNeo4jDriver
}

func (s *MockNeo4jSession) LastBookmark() string {
	return ""
}

func (s *MockNeo4jSession) BeginTransaction(configurers ...func(*neo4jdriver.TransactionConfig)) (neo4jdriver.Transaction, error) {
	return nil, nil
}

func (s *MockNeo4jSession) ReadTransaction(work neo4jdriver.TransactionWork, configurers ...func(*neo4jdriver.TransactionConfig)) (interface{}, error) {
	return nil, nil
}

func (s *MockNeo4jSession) WriteTransaction(work neo4jdriver.TransactionWork, configurers ...func(*neo4jdriver.TransactionConfig)) (interface{}, error) {
	// Check which function is being called based on the query in the transaction work
	if s.driver.ShouldReturnError {
		return 0, errSessionExpired
	}
	
	// Mock the cleanup functions
	if s.driver.CleanupOrphanRelationshipsCalled {
		return int64(s.driver.OrphanRelationshipsCount), nil
	}
	
	if s.driver.CleanupOrphanNodesCalled {
		return int64(s.driver.OrphanNodesCount), nil
	}
	
	return 0, nil
}

func (s *MockNeo4jSession) Run(cypher string, params map[string]interface{}, configurers ...func(*neo4jdriver.TransactionConfig)) (neo4jdriver.Result, error) {
	return nil, nil
}

func (s *MockNeo4jSession) Close() error {
	return nil
}

// Mock implementation of CleanupOrphanRelationships for testing
func MockCleanupOrphanRelationships(driver neo4jdriver.Driver) (int, error) {
	mockDriver, ok := driver.(*MockNeo4jDriver)
	if !ok {
		return 0, errors.New("invalid driver type")
	}
	
	mockDriver.CleanupOrphanRelationshipsCalled = true
	mockDriver.CleanupOrphanRelationshipsCallCount++
	
	if mockDriver.ShouldReturnError {
		return 0, errSessionExpired
	}
	
	return mockDriver.OrphanRelationshipsCount, nil
}

// Mock implementation of CleanupOrphanNodes for testing
func MockCleanupOrphanNodes(driver neo4jdriver.Driver) (int, error) {
	mockDriver, ok := driver.(*MockNeo4jDriver)
	if !ok {
		return 0, errors.New("invalid driver type")
	}
	
	mockDriver.CleanupOrphanNodesCalled = true
	mockDriver.CleanupOrphanNodesCallCount++
	
	if mockDriver.ShouldReturnError {
		return 0, errSessionExpired
	}
	
	return mockDriver.OrphanNodesCount, nil
}

// Mock function for getting related concepts
func mockGetRelatedConcepts(concept string) ([]models.Concept, error) {
	return []models.Concept{
		{
			Name:      "Related1",
			Relation:  "IsA",
			RelatedTo: concept,
		},
		{
			Name:      "Related2",
			Relation:  "HasA",
			RelatedTo: concept,
		},
	}, nil
}

// Mock function for mining relationships
func mockMineRelationship(concept1, concept2 string) (*models.Concept, error) {
	return &models.Concept{
		Name:      concept1,
		Relation:  "Contains",
		RelatedTo: concept2,
	}, nil
}

func TestCleanupOrphans(t *testing.T) {
	// Create a mock driver
	mockDriver := &MockNeo4jDriver{
		OrphanRelationshipsCount: 5,
		OrphanNodesCount:         3,
	}
	
	// Save original functions
	originalCleanupRelationships := kgneo4j.CleanupOrphanRelationships
	originalCleanupNodes := kgneo4j.CleanupOrphanNodes
	
	// Replace with mock functions
	kgneo4j.CleanupOrphanRelationships = MockCleanupOrphanRelationships
	kgneo4j.CleanupOrphanNodes = MockCleanupOrphanNodes
	
	// Restore original functions after test
	defer func() {
		kgneo4j.CleanupOrphanRelationships = originalCleanupRelationships
		kgneo4j.CleanupOrphanNodes = originalCleanupNodes
	}()
	
	// Test successful cleanup
	result, err := cleanup.CleanupOrphans(mockDriver)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if result.OrphanRelationshipsRemoved != 5 {
		t.Errorf("Expected 5 orphan relationships removed, got %d", result.OrphanRelationshipsRemoved)
	}
	
	if result.OrphanNodesRemoved != 3 {
		t.Errorf("Expected 3 orphan nodes removed, got %d", result.OrphanNodesRemoved)
	}
	
	// Test error handling
	mockDriver.ShouldReturnError = true
	_, err = cleanup.CleanupOrphans(mockDriver)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestScheduledCleanup(t *testing.T) {
	// Create a mock driver
	mockDriver := &MockNeo4jDriver{
		OrphanRelationshipsCount: 5,
		OrphanNodesCount:         3,
	}
	
	// Save original functions
	originalCleanupRelationships := kgneo4j.CleanupOrphanRelationships
	originalCleanupNodes := kgneo4j.CleanupOrphanNodes
	
	// Replace with mock functions
	kgneo4j.CleanupOrphanRelationships = MockCleanupOrphanRelationships
	kgneo4j.CleanupOrphanNodes = MockCleanupOrphanNodes
	
	// Restore original functions after test
	defer func() {
		kgneo4j.CleanupOrphanRelationships = originalCleanupRelationships
		kgneo4j.CleanupOrphanNodes = originalCleanupNodes
	}()
	
	// Simulate mining random relationships
	for i := 0; i < 20; i++ {
		// Call cleanup every 5 iterations to simulate periodic cleanup
		if i > 0 && i%5 == 0 {
			cleanup.CleanupOrphans(mockDriver)
		}
	}
	
	// Verify that cleanup was called multiple times
	if mockDriver.CleanupOrphanRelationshipsCallCount < 3 {
		t.Errorf("Expected CleanupOrphanRelationships to be called at least 3 times, got %d", mockDriver.CleanupOrphanRelationshipsCallCount)
	}
	
	if mockDriver.CleanupOrphanNodesCallCount < 3 {
		t.Errorf("Expected CleanupOrphanNodes to be called at least 3 times, got %d", mockDriver.CleanupOrphanNodesCallCount)
	}
}

func TestCleanupWithNoOrphans(t *testing.T) {
	// Create a mock driver with no orphaned relationships or nodes
	mockDriver := &MockNeo4jDriver{
		OrphanRelationshipsCount: 0,
		OrphanNodesCount:         0,
	}
	
	// Save original functions
	originalCleanupRelationships := kgneo4j.CleanupOrphanRelationships
	originalCleanupNodes := kgneo4j.CleanupOrphanNodes
	
	// Replace with mock functions
	kgneo4j.CleanupOrphanRelationships = MockCleanupOrphanRelationships
	kgneo4j.CleanupOrphanNodes = MockCleanupOrphanNodes
	
	// Restore original functions after test
	defer func() {
		kgneo4j.CleanupOrphanRelationships = originalCleanupRelationships
		kgneo4j.CleanupOrphanNodes = originalCleanupNodes
	}()
	
	// Test cleanup with no orphans
	result, err := cleanup.CleanupOrphans(mockDriver)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if result.OrphanRelationshipsRemoved != 0 {
		t.Errorf("Expected 0 orphan relationships removed, got %d", result.OrphanRelationshipsRemoved)
	}
	
	if result.OrphanNodesRemoved != 0 {
		t.Errorf("Expected 0 orphan nodes removed, got %d", result.OrphanNodesRemoved)
	}
}

func TestCleanupWithLargeNumberOfOrphans(t *testing.T) {
	// Create a mock driver with a large number of orphaned relationships and nodes
	mockDriver := &MockNeo4jDriver{
		OrphanRelationshipsCount: 10000,
		OrphanNodesCount:         5000,
	}
	
	// Save original functions
	originalCleanupRelationships := kgneo4j.CleanupOrphanRelationships
	originalCleanupNodes := kgneo4j.CleanupOrphanNodes
	
	// Replace with mock functions
	kgneo4j.CleanupOrphanRelationships = MockCleanupOrphanRelationships
	kgneo4j.CleanupOrphanNodes = MockCleanupOrphanNodes
	
	// Restore original functions after test
	defer func() {
		kgneo4j.CleanupOrphanRelationships = originalCleanupRelationships
		kgneo4j.CleanupOrphanNodes = originalCleanupNodes
	}()
	
	// Test cleanup with a large number of orphans
	result, err := cleanup.CleanupOrphans(mockDriver)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if result.OrphanRelationshipsRemoved != 10000 {
		t.Errorf("Expected 10000 orphan relationships removed, got %d", result.OrphanRelationshipsRemoved)
	}
	
	if result.OrphanNodesRemoved != 5000 {
		t.Errorf("Expected 5000 orphan nodes removed, got %d", result.OrphanNodesRemoved)
	}
}

func TestCleanupWithErrorOnRelationships(t *testing.T) {
	// Create a mock driver that returns an error
	mockDriver := &MockNeo4jDriver{
		ShouldReturnError: true,
	}
	
	// Save original functions
	originalCleanupRelationships := kgneo4j.CleanupOrphanRelationships
	originalCleanupNodes := kgneo4j.CleanupOrphanNodes
	
	// Replace with mock functions
	kgneo4j.CleanupOrphanRelationships = MockCleanupOrphanRelationships
	kgneo4j.CleanupOrphanNodes = MockCleanupOrphanNodes
	
	// Restore original functions after test
	defer func() {
		kgneo4j.CleanupOrphanRelationships = originalCleanupRelationships
		kgneo4j.CleanupOrphanNodes = originalCleanupNodes
	}()
	
	// Test cleanup with an error
	_, err := cleanup.CleanupOrphans(mockDriver)
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
}

func TestBuildGraphWithCleanup(t *testing.T) {
	// Create a mock driver
	mockDriver := &MockNeo4jDriver{
		OrphanRelationshipsCount: 5,
		OrphanNodesCount:         3,
	}
	
	// Save original functions
	originalCleanupRelationships := kgneo4j.CleanupOrphanRelationships
	originalCleanupNodes := kgneo4j.CleanupOrphanNodes
	
	// Replace with mock functions
	kgneo4j.CleanupOrphanRelationships = MockCleanupOrphanRelationships
	kgneo4j.CleanupOrphanNodes = MockCleanupOrphanNodes
	
	// Restore original functions after test
	defer func() {
		kgneo4j.CleanupOrphanRelationships = originalCleanupRelationships
		kgneo4j.CleanupOrphanNodes = originalCleanupNodes
	}()
	
	// Create a configuration
	config := &config.GraphConfig{
		SeedConcept:         "Artificial Intelligence",
		MaxNodes:            100,
		WorkerCount:         5,
		RandomRelationships: 50,
		Concurrency:         5,
	}
	
	// Create a graph builder with our mock functions
	gb := graph.NewGraphBuilder(mockDriver, mockGetRelatedConcepts, mockMineRelationship, config)
	
	// Build the graph with cleanup
	err := gb.BuildGraph("Artificial Intelligence", 10, 5*time.Second)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Verify that cleanup was called at least once
	if mockDriver.CleanupOrphanRelationshipsCallCount < 1 {
		t.Errorf("Expected CleanupOrphanRelationships to be called at least once, got %d", mockDriver.CleanupOrphanRelationshipsCallCount)
	}
	
	if mockDriver.CleanupOrphanNodesCallCount < 1 {
		t.Errorf("Expected CleanupOrphanNodes to be called at least once, got %d", mockDriver.CleanupOrphanNodesCallCount)
	}
} 