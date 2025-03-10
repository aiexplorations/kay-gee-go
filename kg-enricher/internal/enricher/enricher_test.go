package enricher

import (
	"fmt"
	"testing"
	"time"

	"kg-enricher/internal/config"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDriver is a mock implementation of neo4j.Driver
type MockDriver struct {
	mock.Mock
}

func (m *MockDriver) Session(accessMode neo4j.AccessMode, bookmarks ...string) neo4j.Session {
	args := m.Called(accessMode, bookmarks)
	return args.Get(0).(neo4j.Session)
}

func (m *MockDriver) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockSession is a mock implementation of neo4j.Session
type MockSession struct {
	mock.Mock
}

func (m *MockSession) Run(cypher string, params map[string]interface{}) (neo4j.Result, error) {
	args := m.Called(cypher, params)
	return args.Get(0).(neo4j.Result), args.Error(1)
}

func (m *MockSession) ReadTransaction(work neo4j.TransactionWork) (interface{}, error) {
	args := m.Called(work)
	return args.Get(0), args.Error(1)
}

func (m *MockSession) WriteTransaction(work neo4j.TransactionWork) (interface{}, error) {
	args := m.Called(work)
	return args.Get(0), args.Error(1)
}

func (m *MockSession) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSession) LastBookmark() string {
	args := m.Called()
	return args.String(0)
}

// MockResult is a mock implementation of neo4j.Result
type MockResult struct {
	mock.Mock
}

func (m *MockResult) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockResult) Record() *neo4j.Record {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*neo4j.Record)
}

func (m *MockResult) Err() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockResult) Consume() (neo4j.ResultSummary, error) {
	args := m.Called()
	return args.Get(0).(neo4j.ResultSummary), args.Error(1)
}

func TestNewEnricher(t *testing.T) {
	// Create a mock driver
	mockDriver := new(MockDriver)
	
	// Create a config
	cfg := &config.EnricherConfig{
		BatchSize:       10,
		Interval:        time.Second * 60,
		MaxRelationships: 100,
		Concurrency:     5,
	}
	
	// Create an enricher
	e := NewEnricher(mockDriver, cfg)
	
	// Assert that the enricher was created correctly
	assert.NotNil(t, e)
	assert.Equal(t, mockDriver, e.driver)
	assert.Equal(t, cfg, e.config)
	assert.Equal(t, 10, e.config.BatchSize)
	assert.Equal(t, time.Second*60, e.config.Interval)
	assert.Equal(t, 100, e.config.MaxRelationships)
	assert.Equal(t, 5, e.config.Concurrency)
}

func TestGetStats(t *testing.T) {
	// Create a mock driver
	mockDriver := new(MockDriver)
	
	// Create a config
	cfg := &config.EnricherConfig{
		BatchSize:       10,
		Interval:        time.Second * 60,
		MaxRelationships: 100,
		Concurrency:     5,
	}
	
	// Create an enricher
	e := NewEnricher(mockDriver, cfg)
	
	// Set some stats
	e.stats.TotalBatches = 5
	e.stats.TotalPairsProcessed = 50
	e.stats.TotalRelationsFound = 20
	e.stats.TotalRelationsCreated = 15
	e.stats.StartTime = time.Now().Add(-time.Hour)
	e.stats.LastBatchTime = time.Now().Add(-time.Minute)
	
	// Get the stats
	stats := e.GetStats()
	
	// Assert that the stats are correct
	assert.Equal(t, e.stats, stats)
	assert.Equal(t, 5, stats.TotalBatches)
	assert.Equal(t, 50, stats.TotalPairsProcessed)
	assert.Equal(t, 20, stats.TotalRelationsFound)
	assert.Equal(t, 15, stats.TotalRelationsCreated)
	assert.WithinDuration(t, e.stats.StartTime, stats.StartTime, time.Second)
	assert.WithinDuration(t, e.stats.LastBatchTime, stats.LastBatchTime, time.Second)
}

func TestRunOnce(t *testing.T) {
	// Create a mock driver
	mockDriver := new(MockDriver)
	mockSession := new(MockSession)
	mockResult := new(MockResult)
	
	// Set up mock behavior
	mockDriver.On("Session", mock.Anything, mock.Anything).Return(mockSession)
	mockSession.On("ReadTransaction", mock.AnythingOfType("neo4j.TransactionWork")).Return([][]string{
		{"Concept1", "Concept2"},
		{"Concept3", "Concept4"},
		{"Concept5", "Concept6"},
	}, nil)
	mockSession.On("WriteTransaction", mock.AnythingOfType("neo4j.TransactionWork")).Return(nil, nil)
	mockSession.On("Close").Return(nil)
	mockResult.On("Next").Return(false)
	mockResult.On("Err").Return(nil)
	
	// Create a config
	cfg := &config.EnricherConfig{
		BatchSize:        10,
		Interval:         time.Second * 60,
		MaxRelationships: 100,
		Concurrency:      1, // Use 1 for predictable testing
	}
	
	// Create an enricher with a mock LLM function
	e := NewEnricher(mockDriver, cfg)
	
	// Replace the findRelationship function with a mock
	originalFindRelationship := findRelationship
	defer func() { findRelationship = originalFindRelationship }()
	
	// Mock the findRelationship function to always return a relationship
	findRelationship = func(concept1, concept2 string) (string, error) {
		return "IsRelatedTo", nil
	}
	
	// Run once with 3 pairs
	err := e.RunOnce(3)
	
	// Assert that there was no error
	assert.NoError(t, err)
	
	// Assert that the stats were updated correctly
	assert.Equal(t, 1, e.stats.TotalBatches)
	assert.Equal(t, 3, e.stats.TotalPairsProcessed)
	assert.Equal(t, 3, e.stats.TotalRelationsFound)
	assert.Equal(t, 3, e.stats.TotalRelationsCreated)
	
	// Assert that the mock methods were called the expected number of times
	mockDriver.AssertNumberOfCalls(t, "Session", 2) // One for reading pairs, one for writing relationships
	mockSession.AssertNumberOfCalls(t, "ReadTransaction", 1)
	mockSession.AssertNumberOfCalls(t, "WriteTransaction", 3) // One for each relationship
	mockSession.AssertNumberOfCalls(t, "Close", 2)
}

func TestRunOnceWithNoRelationships(t *testing.T) {
	// Create a mock driver
	mockDriver := new(MockDriver)
	mockSession := new(MockSession)
	mockResult := new(MockResult)
	
	// Set up mock behavior
	mockDriver.On("Session", mock.Anything, mock.Anything).Return(mockSession)
	mockSession.On("ReadTransaction", mock.AnythingOfType("neo4j.TransactionWork")).Return([][]string{
		{"Concept1", "Concept2"},
		{"Concept3", "Concept4"},
		{"Concept5", "Concept6"},
	}, nil)
	mockSession.On("WriteTransaction", mock.AnythingOfType("neo4j.TransactionWork")).Return(nil, nil)
	mockSession.On("Close").Return(nil)
	mockResult.On("Next").Return(false)
	mockResult.On("Err").Return(nil)
	
	// Create a config
	cfg := &config.EnricherConfig{
		BatchSize:        10,
		Interval:         time.Second * 60,
		MaxRelationships: 100,
		Concurrency:      1, // Use 1 for predictable testing
	}
	
	// Create an enricher with a mock LLM function
	e := NewEnricher(mockDriver, cfg)
	
	// Replace the findRelationship function with a mock
	originalFindRelationship := findRelationship
	defer func() { findRelationship = originalFindRelationship }()
	
	// Mock the findRelationship function to always return no relationship
	findRelationship = func(concept1, concept2 string) (string, error) {
		return "", nil
	}
	
	// Run once with 3 pairs
	err := e.RunOnce(3)
	
	// Assert that there was no error
	assert.NoError(t, err)
	
	// Assert that the stats were updated correctly
	assert.Equal(t, 1, e.stats.TotalBatches)
	assert.Equal(t, 3, e.stats.TotalPairsProcessed)
	assert.Equal(t, 0, e.stats.TotalRelationsFound)
	assert.Equal(t, 0, e.stats.TotalRelationsCreated)
	
	// Assert that the mock methods were called the expected number of times
	mockDriver.AssertNumberOfCalls(t, "Session", 1) // Only for reading pairs
	mockSession.AssertNumberOfCalls(t, "ReadTransaction", 1)
	mockSession.AssertNumberOfCalls(t, "WriteTransaction", 0) // No relationships to create
	mockSession.AssertNumberOfCalls(t, "Close", 1)
}

func TestRunOnceWithError(t *testing.T) {
	// Create a mock driver
	mockDriver := new(MockDriver)
	mockSession := new(MockSession)
	
	// Set up mock behavior to return an error
	mockDriver.On("Session", mock.Anything, mock.Anything).Return(mockSession)
	mockSession.On("ReadTransaction", mock.AnythingOfType("neo4j.TransactionWork")).Return(nil, fmt.Errorf("mock error"))
	mockSession.On("Close").Return(nil)
	
	// Create a config
	cfg := &config.EnricherConfig{
		BatchSize:        10,
		Interval:         time.Second * 60,
		MaxRelationships: 100,
		Concurrency:      1,
	}
	
	// Create an enricher
	e := NewEnricher(mockDriver, cfg)
	
	// Run once
	err := e.RunOnce(3)
	
	// Assert that there was an error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mock error")
	
	// Assert that the stats were not updated
	assert.Equal(t, 0, e.stats.TotalBatches)
	assert.Equal(t, 0, e.stats.TotalPairsProcessed)
	assert.Equal(t, 0, e.stats.TotalRelationsFound)
	assert.Equal(t, 0, e.stats.TotalRelationsCreated)
}

func TestStartAndStop(t *testing.T) {
	// Create a mock driver
	mockDriver := new(MockDriver)
	mockSession := new(MockSession)
	
	// Set up mock behavior
	mockDriver.On("Session", mock.Anything, mock.Anything).Return(mockSession)
	mockSession.On("ReadTransaction", mock.AnythingOfType("neo4j.TransactionWork")).Return([][]string{
		{"Concept1", "Concept2"},
		{"Concept3", "Concept4"},
	}, nil)
	mockSession.On("WriteTransaction", mock.AnythingOfType("neo4j.TransactionWork")).Return(nil, nil)
	mockSession.On("Close").Return(nil)
	
	// Create a config with a short interval
	cfg := &config.EnricherConfig{
		BatchSize:        10,
		Interval:         time.Millisecond * 100, // Short interval for testing
		MaxRelationships: 100,
		Concurrency:      1,
	}
	
	// Create an enricher
	e := NewEnricher(mockDriver, cfg)
	
	// Replace the findRelationship function with a mock
	originalFindRelationship := findRelationship
	defer func() { findRelationship = originalFindRelationship }()
	
	// Mock the findRelationship function to always return a relationship
	findRelationship = func(concept1, concept2 string) (string, error) {
		return "IsRelatedTo", nil
	}
	
	// Check that the enricher is not running
	assert.False(t, e.IsRunning())
	
	// Start the enricher
	err := e.Start()
	
	// Check that there was no error
	assert.NoError(t, err)
	
	// Check that the enricher is running
	assert.True(t, e.IsRunning())
	
	// Try to start the enricher again
	err = e.Start()
	
	// Check that there was an error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already running")
	
	// Wait for the enricher to process at least one batch
	time.Sleep(time.Millisecond * 200)
	
	// Stop the enricher
	err = e.Stop()
	
	// Check that there was no error
	assert.NoError(t, err)
	
	// Check that the enricher is not running
	assert.False(t, e.IsRunning())
	
	// Try to stop the enricher again
	err = e.Stop()
	
	// Check that there was an error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not running")
	
	// Check that the stats were updated
	stats := e.GetStats()
	assert.GreaterOrEqual(t, stats.TotalBatches, 1)
	assert.GreaterOrEqual(t, stats.TotalPairsProcessed, 2)
	assert.GreaterOrEqual(t, stats.TotalRelationsFound, 2)
	assert.GreaterOrEqual(t, stats.TotalRelationsCreated, 2)
}

func TestIsRunning(t *testing.T) {
	// Create a mock driver
	mockDriver := new(MockDriver)
	
	// Create a config
	cfg := &config.EnricherConfig{
		BatchSize:        10,
		Interval:         time.Second * 60,
		MaxRelationships: 100,
		Concurrency:      1,
	}
	
	// Create an enricher
	e := NewEnricher(mockDriver, cfg)
	
	// Check that the enricher is not running
	assert.False(t, e.IsRunning())
	
	// Set the running flag using reflection
	e.mutex.Lock()
	e.running = true
	e.mutex.Unlock()
	
	// Check that the enricher is running
	assert.True(t, e.IsRunning())
	
	// Set the running flag back to false
	e.mutex.Lock()
	e.running = false
	e.mutex.Unlock()
	
	// Check that the enricher is not running
	assert.False(t, e.IsRunning())
}

func TestProcessBatch(t *testing.T) {
	// Create a mock driver
	mockDriver := new(MockDriver)
	mockSession := new(MockSession)
	
	// Set up mock behavior
	mockDriver.On("Session", mock.Anything, mock.Anything).Return(mockSession)
	mockSession.On("ReadTransaction", mock.AnythingOfType("neo4j.TransactionWork")).Return([][]string{
		{"Concept1", "Concept2"},
		{"Concept3", "Concept4"},
		{"Concept5", "Concept6"},
	}, nil)
	mockSession.On("WriteTransaction", mock.AnythingOfType("neo4j.TransactionWork")).Return(nil, nil)
	mockSession.On("Close").Return(nil)
	
	// Create a config
	cfg := &config.EnricherConfig{
		BatchSize:        10,
		Interval:         time.Second * 60,
		MaxRelationships: 100,
		Concurrency:      1,
	}
	
	// Create an enricher
	e := NewEnricher(mockDriver, cfg)
	
	// Replace the findRelationship function with a mock
	originalFindRelationship := findRelationship
	defer func() { findRelationship = originalFindRelationship }()
	
	// Mock the findRelationship function to return different relationships based on the concepts
	findRelationship = func(concept1, concept2 string) (string, error) {
		if concept1 == "Concept1" && concept2 == "Concept2" {
			return "IsA", nil
		} else if concept1 == "Concept3" && concept2 == "Concept4" {
			return "HasPart", nil
		} else if concept1 == "Concept5" && concept2 == "Concept6" {
			return "", nil // No relationship
		}
		return "", nil
	}
	
	// Call processBatch directly
	e.processBatch()
	
	// Check that the stats were updated correctly
	assert.Equal(t, 1, e.stats.TotalBatches)
	assert.Equal(t, 3, e.stats.TotalPairsProcessed)
	assert.Equal(t, 2, e.stats.TotalRelationsFound)
	assert.Equal(t, 2, e.stats.TotalRelationsCreated)
	
	// Check that the mock methods were called the expected number of times
	mockDriver.AssertNumberOfCalls(t, "Session", 2) // One for reading pairs, one for writing relationships
	mockSession.AssertNumberOfCalls(t, "ReadTransaction", 1)
	mockSession.AssertNumberOfCalls(t, "WriteTransaction", 2) // One for each relationship
	mockSession.AssertNumberOfCalls(t, "Close", 2)
} 