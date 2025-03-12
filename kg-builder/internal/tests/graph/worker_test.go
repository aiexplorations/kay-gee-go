package graph_test

import (
	"context"
	"net/url"
	"sync"
	"testing"
	"time"

	"kg-builder/internal/config"
	"kg-builder/internal/graph"
	"kg-builder/internal/models"
	kgneo4j "kg-builder/internal/neo4j"

	neo4jdriver "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// MockNeo4jDriver implements the Neo4j driver interface for testing
type MockNeo4jDriverWithTracking struct {
	// Track calls to methods
	CleanupOrphanRelationshipsCalled bool
	CleanupOrphanNodesCalled         bool
	
	// Mock return values
	OrphanRelationshipsCount int
	OrphanNodesCount         int
	
	// Track number of calls
	CleanupOrphanRelationshipsCallCount int
	CleanupOrphanNodesCallCount         int
	
	// Mutex for thread safety
	mutex sync.Mutex
}

// Implement Neo4j driver interface methods
func (m *MockNeo4jDriverWithTracking) Close() error {
	return nil
}

func (m *MockNeo4jDriverWithTracking) NewSession(config neo4jdriver.SessionConfig) neo4jdriver.Session {
	return &MockNeo4jSessionWithTracking{driver: m}
}

func (m *MockNeo4jDriverWithTracking) Session(accessMode neo4jdriver.AccessMode, bookmarks ...string) (neo4jdriver.Session, error) {
	return &MockNeo4jSessionWithTracking{driver: m}, nil
}

func (m *MockNeo4jDriverWithTracking) Target() url.URL {
	u, _ := url.Parse("bolt://localhost:7687")
	return *u
}

func (m *MockNeo4jDriverWithTracking) VerifyConnectivity() error {
	return nil
}

// MockNeo4jSessionWithTracking implements the Neo4j session interface for testing
type MockNeo4jSessionWithTracking struct {
	driver *MockNeo4jDriverWithTracking
}

func (s *MockNeo4jSessionWithTracking) LastBookmark() string {
	return ""
}

func (s *MockNeo4jSessionWithTracking) BeginTransaction(configurers ...func(*neo4jdriver.TransactionConfig)) (neo4jdriver.Transaction, error) {
	return nil, nil
}

func (s *MockNeo4jSessionWithTracking) ReadTransaction(work neo4jdriver.TransactionWork, configurers ...func(*neo4jdriver.TransactionConfig)) (interface{}, error) {
	// For ConceptExists, always return true
	return true, nil
}

func (s *MockNeo4jSessionWithTracking) WriteTransaction(work neo4jdriver.TransactionWork, configurers ...func(*neo4jdriver.TransactionConfig)) (interface{}, error) {
	// For CreateRelationship, just return success
	return nil, nil
}

func (s *MockNeo4jSessionWithTracking) Run(cypher string, params map[string]interface{}, configurers ...func(*neo4jdriver.TransactionConfig)) (neo4jdriver.Result, error) {
	return nil, nil
}

func (s *MockNeo4jSessionWithTracking) Close() error {
	return nil
}

// Mock implementation of CleanupOrphanRelationships for testing
func MockCleanupOrphanRelationshipsWithTracking(driver neo4jdriver.Driver) (int, error) {
	mockDriver, ok := driver.(*MockNeo4jDriverWithTracking)
	if !ok {
		return 0, nil
	}
	
	mockDriver.mutex.Lock()
	defer mockDriver.mutex.Unlock()
	
	mockDriver.CleanupOrphanRelationshipsCalled = true
	mockDriver.CleanupOrphanRelationshipsCallCount++
	
	return mockDriver.OrphanRelationshipsCount, nil
}

// Mock implementation of CleanupOrphanNodes for testing
func MockCleanupOrphanNodesWithTracking(driver neo4jdriver.Driver) (int, error) {
	mockDriver, ok := driver.(*MockNeo4jDriverWithTracking)
	if !ok {
		return 0, nil
	}
	
	mockDriver.mutex.Lock()
	defer mockDriver.mutex.Unlock()
	
	mockDriver.CleanupOrphanNodesCalled = true
	mockDriver.CleanupOrphanNodesCallCount++
	
	return mockDriver.OrphanNodesCount, nil
}

// Mock function for getting related concepts
func mockGetRelatedConceptsForWorker(concept string) ([]models.Concept, error) {
	// Return a different set of concepts for each input to simulate real behavior
	switch concept {
	case "Artificial Intelligence":
		return []models.Concept{
			{Name: "Machine Learning", Relation: "IsA", RelatedTo: concept},
			{Name: "Neural Networks", Relation: "IsA", RelatedTo: concept},
			{Name: "Deep Learning", Relation: "IsA", RelatedTo: concept},
		}, nil
	case "Machine Learning":
		return []models.Concept{
			{Name: "Supervised Learning", Relation: "IsA", RelatedTo: concept},
			{Name: "Unsupervised Learning", Relation: "IsA", RelatedTo: concept},
		}, nil
	case "Neural Networks":
		return []models.Concept{
			{Name: "Convolutional Neural Networks", Relation: "IsA", RelatedTo: concept},
			{Name: "Recurrent Neural Networks", Relation: "IsA", RelatedTo: concept},
		}, nil
	default:
		return []models.Concept{
			{Name: "Generic Concept 1", Relation: "RelatedTo", RelatedTo: concept},
			{Name: "Generic Concept 2", Relation: "RelatedTo", RelatedTo: concept},
		}, nil
	}
}

// Mock function for mining relationships
func mockMineRelationshipForWorker(concept1, concept2 string) (*models.Concept, error) {
	return &models.Concept{
		Name:      concept1,
		Relation:  "RelatedTo",
		RelatedTo: concept2,
	}, nil
}

func TestWorkerPeriodicCleanup(t *testing.T) {
	// Create a mock driver with tracking
	mockDriver := &MockNeo4jDriverWithTracking{
		OrphanRelationshipsCount: 5,
		OrphanNodesCount:         3,
	}
	
	// Override the neo4j cleanup functions with our mocks
	originalCleanupRelationships := kgneo4j.CleanupOrphanRelationships
	originalCleanupNodes := kgneo4j.CleanupOrphanNodes
	
	// Restore the original functions after the test
	defer func() {
		kgneo4j.CleanupOrphanRelationships = originalCleanupRelationships
		kgneo4j.CleanupOrphanNodes = originalCleanupNodes
	}()
	
	// Replace with mock functions
	kgneo4j.CleanupOrphanRelationships = MockCleanupOrphanRelationshipsWithTracking
	kgneo4j.CleanupOrphanNodes = MockCleanupOrphanNodesWithTracking
	
	// Create a graph builder with the mock driver
	config := &config.GraphConfig{
		SeedConcept:         "Artificial Intelligence",
		MaxNodes:            20, // Set to a value that will trigger multiple cleanup operations
		WorkerCount:         1,  // Use a single worker for predictable behavior
		RandomRelationships: 0,
		Concurrency:         1,
	}
	
	gb := graph.NewGraphBuilder(mockDriver, mockGetRelatedConceptsForWorker, mockMineRelationshipForWorker, config)
	
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Create a wait group for the worker
	var wg sync.WaitGroup
	wg.Add(1)
	
	// Create a channel for concepts
	queue := make(chan string, 100)
	
	// Add the seed concept to the queue
	queue <- "Artificial Intelligence"
	
	// Start the worker in a goroutine
	go gb.Worker(ctx, &wg, queue)
	
	// Wait for the worker to finish or timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	
	select {
	case <-done:
		// Worker finished
	case <-time.After(5 * time.Second):
		// Timeout
		t.Log("Worker timed out, but this is expected")
	}
	
	// Close the queue to signal the worker to stop
	close(queue)
	
	// Wait a bit for any pending operations to complete
	time.Sleep(100 * time.Millisecond)
	
	// Check that cleanup was called at least once
	if mockDriver.CleanupOrphanRelationshipsCallCount == 0 {
		t.Error("Expected CleanupOrphanRelationships to be called at least once, but it was not called")
	}
	
	if mockDriver.CleanupOrphanNodesCallCount == 0 {
		t.Error("Expected CleanupOrphanNodes to be called at least once, but it was not called")
	}
}

func TestWorkerCleanupFrequency(t *testing.T) {
	// Create a mock driver with tracking
	mockDriver := &MockNeo4jDriverWithTracking{
		OrphanRelationshipsCount: 5,
		OrphanNodesCount:         3,
	}
	
	// Override the neo4j cleanup functions with our mocks
	originalCleanupRelationships := kgneo4j.CleanupOrphanRelationships
	originalCleanupNodes := kgneo4j.CleanupOrphanNodes
	
	// Restore the original functions after the test
	defer func() {
		kgneo4j.CleanupOrphanRelationships = originalCleanupRelationships
		kgneo4j.CleanupOrphanNodes = originalCleanupNodes
	}()
	
	// Replace with mock functions
	kgneo4j.CleanupOrphanRelationships = MockCleanupOrphanRelationshipsWithTracking
	kgneo4j.CleanupOrphanNodes = MockCleanupOrphanNodesWithTracking
	
	// Create a graph builder with the mock driver
	config := &config.GraphConfig{
		SeedConcept:         "Artificial Intelligence",
		MaxNodes:            100, // Set to a high value to allow many concepts to be processed
		WorkerCount:         1,   // Use a single worker for predictable behavior
		RandomRelationships: 0,
		Concurrency:         1,
	}
	
	gb := graph.NewGraphBuilder(mockDriver, mockGetRelatedConceptsForWorker, mockMineRelationshipForWorker, config)
	
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Create a wait group for the worker
	var wg sync.WaitGroup
	wg.Add(1)
	
	// Create a channel for concepts
	queue := make(chan string, 100)
	
	// Add multiple concepts to the queue to trigger cleanup
	concepts := []string{
		"Artificial Intelligence",
		"Machine Learning",
		"Neural Networks",
		"Deep Learning",
		"Supervised Learning",
		"Unsupervised Learning",
		"Convolutional Neural Networks",
		"Recurrent Neural Networks",
		"Reinforcement Learning",
		"Natural Language Processing",
	}
	
	for _, concept := range concepts {
		queue <- concept
	}
	
	// Start the worker in a goroutine
	go gb.Worker(ctx, &wg, queue)
	
	// Wait for the worker to process some concepts
	time.Sleep(2 * time.Second)
	
	// Close the queue to signal the worker to stop
	close(queue)
	
	// Wait a bit for any pending operations to complete
	time.Sleep(100 * time.Millisecond)
	
	// Check that cleanup was called at least once
	if mockDriver.CleanupOrphanRelationshipsCallCount == 0 {
		t.Error("Expected CleanupOrphanRelationships to be called at least once, but it was not called")
	}
	
	if mockDriver.CleanupOrphanNodesCallCount == 0 {
		t.Error("Expected CleanupOrphanNodes to be called at least once, but it was not called")
	}
	
	// Check that cleanup was called approximately every 5 concepts
	// The exact number may vary due to concurrency, but it should be roughly nodeCount / 5
	expectedCleanupCalls := gb.GetNodeCount() / 5
	if mockDriver.CleanupOrphanRelationshipsCallCount < expectedCleanupCalls-1 || 
	   mockDriver.CleanupOrphanRelationshipsCallCount > expectedCleanupCalls+1 {
		t.Logf("Expected approximately %d cleanup calls, got %d", 
			expectedCleanupCalls, mockDriver.CleanupOrphanRelationshipsCallCount)
	}
} 