package graph_test

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
	"unsafe"

	"kg-builder/internal/config"
	"kg-builder/internal/graph"
	"kg-builder/internal/models"

	neo4jdriver "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// Mock Neo4j driver for testing
type MockNeo4jDriver struct{}

func (m *MockNeo4jDriver) Close() error {
	return nil
}

func (m *MockNeo4jDriver) NewSession(config neo4jdriver.SessionConfig) neo4jdriver.Session {
	return nil
}

func (m *MockNeo4jDriver) Session(accessMode neo4jdriver.AccessMode, bookmarks ...string) (neo4jdriver.Session, error) {
	return nil, nil
}

func (m *MockNeo4jDriver) VerifyConnectivity() error {
	return nil
}

func (m *MockNeo4jDriver) Target() url.URL {
	u, _ := url.Parse("bolt://localhost:7687")
	return *u
}

// Mock function for getting related concepts
func mockGetRelatedConcepts(concept string) ([]models.Concept, error) {
	if concept == "error" {
		return nil, errors.New("mock error")
	}
	
	if concept == "empty" {
		return []models.Concept{}, nil
	}
	
	return []models.Concept{
		{Name: "Related1", Relation: "IsA", RelatedTo: concept},
		{Name: "Related2", Relation: "HasPart", RelatedTo: concept},
	}, nil
}

// Mock function for mining relationships
func mockMineRelationship(concept1, concept2 string) (*models.Concept, error) {
	if concept1 == "error" || concept2 == "error" {
		return nil, errors.New("mock error")
	}
	
	if concept1 == "empty" || concept2 == "empty" {
		return nil, nil
	}
	
	return &models.Concept{
		Name: concept2, 
		Relation: "IsRelatedTo", 
		RelatedTo: concept1,
	}, nil
}

func TestGraphBuilderCreation(t *testing.T) {
	// Test with valid parameters
	driver := &MockNeo4jDriver{}
	config := &config.GraphConfig{
		SeedConcept:         "Artificial Intelligence",
		MaxNodes:            100,
		Timeout:             30 * time.Minute,
		WorkerCount:         10,
		RandomRelationships: 50,
		Concurrency:         5,
	}
	gb := graph.NewGraphBuilder(driver, mockGetRelatedConcepts, mockMineRelationship, config)
	
	if gb == nil {
		t.Error("Expected GraphBuilder to be created, got nil")
	}
}

func TestBuildGraphValidation(t *testing.T) {
	// Test with empty seed concept
	driver := &MockNeo4jDriver{}
	config := &config.GraphConfig{
		SeedConcept:         "Artificial Intelligence",
		MaxNodes:            100,
		Timeout:             30 * time.Minute,
		WorkerCount:         10,
		RandomRelationships: 50,
		Concurrency:         5,
	}
	gb := graph.NewGraphBuilder(driver, mockGetRelatedConcepts, mockMineRelationship, config)
	
	err := gb.BuildGraph("", 10, 1*time.Minute)
	if err == nil {
		t.Error("Expected error with empty seed concept, got nil")
	}
	
	// Test with zero max nodes
	err = gb.BuildGraph("test", 0, 1*time.Minute)
	if err == nil {
		t.Error("Expected error with zero max nodes, got nil")
	}
	
	// Test with zero timeout
	err = gb.BuildGraph("test", 10, 0)
	if err == nil {
		t.Error("Expected error with zero timeout, got nil")
	}
}

func TestMineRandomRelationshipsValidation(t *testing.T) {
	// Test with zero count
	driver := &MockNeo4jDriver{}
	config := &config.GraphConfig{
		SeedConcept:         "Artificial Intelligence",
		MaxNodes:            100,
		Timeout:             30 * time.Minute,
		WorkerCount:         10,
		RandomRelationships: 50,
		Concurrency:         5,
	}
	gb := graph.NewGraphBuilder(driver, mockGetRelatedConcepts, mockMineRelationship, config)
	
	err := gb.MineRandomRelationships(0, 1)
	if err == nil {
		t.Error("Expected error with zero count, got nil")
	}
	
	// Test with zero concurrency
	err = gb.MineRandomRelationships(10, 0)
	if err == nil {
		t.Error("Expected error with zero concurrency, got nil")
	}
	
	// Test with not enough processed concepts
	err = gb.MineRandomRelationships(10, 1)
	if err == nil {
		t.Error("Expected error with not enough processed concepts, got nil")
	}
}

func TestGetProcessedConcepts(t *testing.T) {
	// Test with empty processed concepts
	driver := &MockNeo4jDriver{}
	config := &config.GraphConfig{
		SeedConcept:         "Artificial Intelligence",
		MaxNodes:            100,
		Timeout:             30 * time.Minute,
		WorkerCount:         10,
		RandomRelationships: 50,
		Concurrency:         5,
	}
	gb := graph.NewGraphBuilder(driver, mockGetRelatedConcepts, mockMineRelationship, config)
	
	concepts := gb.GetProcessedConcepts()
	if len(concepts) != 0 {
		t.Errorf("Expected 0 processed concepts, got %d", len(concepts))
	}
}

func TestGetNodeCount(t *testing.T) {
	// Test with zero node count
	driver := &MockNeo4jDriver{}
	config := &config.GraphConfig{
		SeedConcept:         "Artificial Intelligence",
		MaxNodes:            100,
		Timeout:             30 * time.Minute,
		WorkerCount:         10,
		RandomRelationships: 50,
		Concurrency:         5,
	}
	gb := graph.NewGraphBuilder(driver, mockGetRelatedConcepts, mockMineRelationship, config)
	
	count := gb.GetNodeCount()
	if count != 0 {
		t.Errorf("Expected node count 0, got %d", count)
	}
}

func TestBuildGraphConcurrency(t *testing.T) {
	// Create a mock driver
	driver := &MockNeo4jDriver{}
	
	// Create a config with a high worker count
	config := &config.GraphConfig{
		SeedConcept:         "Artificial Intelligence",
		MaxNodes:            100,
		Timeout:             30 * time.Minute,
		WorkerCount:         10, // High worker count to test concurrency
		RandomRelationships: 50,
		Concurrency:         5,
	}
	
	// Create a channel to track concurrent executions
	concurrentExecutions := make(chan struct{}, config.WorkerCount)
	executionCount := 0
	var executionMutex sync.Mutex
	
	// Mock function that simulates concurrent executions
	mockGetRelatedConcepts := func(concept string) ([]models.Concept, error) {
		// Signal that we're executing
		concurrentExecutions <- struct{}{}
		
		// Track the maximum concurrent executions
		executionMutex.Lock()
		executionCount++
		executionMutex.Unlock()
		
		// Simulate work
		time.Sleep(100 * time.Millisecond)
		
		// Signal that we're done
		<-concurrentExecutions
		
		// Return some concepts
		return []models.Concept{
			{Name: fmt.Sprintf("Related1-%s", concept), Relation: "IsA", RelatedTo: concept},
			{Name: fmt.Sprintf("Related2-%s", concept), Relation: "HasPart", RelatedTo: concept},
		}, nil
	}
	
	// Create a graph builder with our mock functions
	gb := graph.NewGraphBuilder(driver, mockGetRelatedConcepts, mockMineRelationship, config)
	
	// Set a low max nodes to make the test run faster
	maxNodes := 20
	
	// Build the graph with a short timeout
	err := gb.BuildGraph("test", maxNodes, 5*time.Second)
	
	// Check that there was no error
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Check that we processed the expected number of nodes
	if gb.GetNodeCount() != maxNodes {
		t.Errorf("Expected %d nodes, got %d", maxNodes, gb.GetNodeCount())
	}
	
	// Check that we had concurrent executions
	if len(concurrentExecutions) == 0 && executionCount < 2 {
		t.Errorf("Expected concurrent executions, got sequential execution")
	}
}

func TestBuildGraphMaxNodes(t *testing.T) {
	// Create a mock driver
	driver := &MockNeo4jDriver{}
	
	// Create a config
	config := &config.GraphConfig{
		SeedConcept:         "Artificial Intelligence",
		MaxNodes:            100,
		Timeout:             30 * time.Minute,
		WorkerCount:         5,
		RandomRelationships: 50,
		Concurrency:         5,
	}
	
	// Counter for the number of times getRelatedConcepts is called
	callCount := 0
	
	// Mock function that returns a large number of concepts
	mockGetRelatedConcepts := func(concept string) ([]models.Concept, error) {
		callCount++
		
		// Return a lot of concepts to quickly reach max nodes
		concepts := make([]models.Concept, 10)
		for i := 0; i < 10; i++ {
			concepts[i] = models.Concept{
				Name:      fmt.Sprintf("Related%d-%s", i, concept),
				Relation:  "IsA",
				RelatedTo: concept,
			}
		}
		return concepts, nil
	}
	
	// Create a graph builder with our mock functions
	gb := graph.NewGraphBuilder(driver, mockGetRelatedConcepts, mockMineRelationship, config)
	
	// Set a low max nodes
	maxNodes := 15
	
	// Build the graph
	err := gb.BuildGraph("test", maxNodes, 30*time.Second)
	
	// Check that there was no error
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Check that we processed exactly maxNodes nodes
	if gb.GetNodeCount() != maxNodes {
		t.Errorf("Expected %d nodes, got %d", maxNodes, gb.GetNodeCount())
	}
	
	// Check that getRelatedConcepts was called at least once
	if callCount == 0 {
		t.Errorf("Expected getRelatedConcepts to be called at least once")
	}
}

func TestBuildGraphTimeout(t *testing.T) {
	// Create a mock driver
	driver := &MockNeo4jDriver{}
	
	// Create a config
	config := &config.GraphConfig{
		SeedConcept:         "Artificial Intelligence",
		MaxNodes:            100,
		Timeout:             30 * time.Minute,
		WorkerCount:         5,
		RandomRelationships: 50,
		Concurrency:         5,
	}
	
	// Mock function that takes a long time to execute
	mockGetRelatedConcepts := func(concept string) ([]models.Concept, error) {
		// Simulate a long-running operation
		time.Sleep(500 * time.Millisecond)
		
		return []models.Concept{
			{Name: fmt.Sprintf("Related1-%s", concept), Relation: "IsA", RelatedTo: concept},
			{Name: fmt.Sprintf("Related2-%s", concept), Relation: "HasPart", RelatedTo: concept},
		}, nil
	}
	
	// Create a graph builder with our mock functions
	gb := graph.NewGraphBuilder(driver, mockGetRelatedConcepts, mockMineRelationship, config)
	
	// Set a very short timeout
	timeout := 100 * time.Millisecond
	
	// Build the graph with a very short timeout
	err := gb.BuildGraph("test", 100, timeout)
	
	// Check that we got a timeout error
	if err == nil {
		t.Errorf("Expected timeout error, got nil")
	} else if !strings.Contains(err.Error(), "timed out") {
		t.Errorf("Expected timeout error, got %v", err)
	}
}

func TestGetRandomPair(t *testing.T) {
	// Create a mock driver
	driver := &MockNeo4jDriver{}
	
	// Create a config
	config := &config.GraphConfig{
		SeedConcept:         "Artificial Intelligence",
		MaxNodes:            100,
		Timeout:             30 * time.Minute,
		WorkerCount:         5,
		RandomRelationships: 50,
		Concurrency:         5,
	}
	
	// Create a graph builder
	gb := graph.NewGraphBuilder(driver, mockGetRelatedConcepts, mockMineRelationship, config)
	
	// Add some processed concepts
	processedConcepts := gb.GetProcessedConcepts()
	if len(processedConcepts) != 0 {
		t.Errorf("Expected 0 processed concepts, got %d", len(processedConcepts))
	}
	
	// Test with no processed concepts
	pair := gb.GetRandomPair()
	if pair[0] != "" || pair[1] != "" {
		t.Errorf("Expected empty pair with no processed concepts, got %v", pair)
	}
	
	// Add some processed concepts using reflection
	// Note: This is a bit of a hack, but it's the only way to test this function
	// without exposing the processedConcepts field
	processedConceptsField := reflect.ValueOf(gb).Elem().FieldByName("processedConcepts")
	if !processedConceptsField.IsValid() {
		t.Fatalf("Could not find processedConcepts field")
	}
	
	// Make the field accessible
	processedConceptsField = reflect.NewAt(processedConceptsField.Type(), unsafe.Pointer(processedConceptsField.UnsafeAddr())).Elem()
	
	// Set some processed concepts
	processedConceptsMap := make(map[string]bool)
	processedConceptsMap["Concept1"] = true
	processedConceptsMap["Concept2"] = true
	processedConceptsMap["Concept3"] = true
	processedConceptsField.Set(reflect.ValueOf(processedConceptsMap))
	
	// Test with processed concepts
	pair = gb.GetRandomPair()
	if pair[0] == "" || pair[1] == "" {
		t.Errorf("Expected non-empty pair with processed concepts, got %v", pair)
	}
	if pair[0] == pair[1] {
		t.Errorf("Expected different concepts in pair, got %v", pair)
	}
	if !processedConceptsMap[pair[0]] || !processedConceptsMap[pair[1]] {
		t.Errorf("Expected pair to contain processed concepts, got %v", pair)
	}
} 