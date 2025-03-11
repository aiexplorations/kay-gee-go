package frontend

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kay-gee-go/internal/common/config"
	"github.com/stretchr/testify/assert"
)

func TestNewFrontend(t *testing.T) {
	// Create a test configuration
	cfg := &config.Neo4jConfig{
		URI:               "bolt://localhost:7687",
		User:              "neo4j",
		Password:          "password",
		MaxRetries:        3,
		RetryIntervalSecs: 1,
	}

	// Skip the test if we're not running in an environment with Neo4j
	t.Skip("Skipping test that requires Neo4j")

	// Create a new frontend
	f, err := NewFrontend(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, f)
	defer f.Close()

	// Check that the frontend was initialized correctly
	assert.Equal(t, cfg, f.config)
	assert.NotNil(t, f.neo4jClient)
	assert.NotNil(t, f.templates)
}

func TestCreateFiles(t *testing.T) {
	// Create a test configuration
	cfg := &config.Neo4jConfig{
		URI:               "bolt://localhost:7687",
		User:              "neo4j",
		Password:          "password",
		MaxRetries:        3,
		RetryIntervalSecs: 1,
	}

	// Skip the test if we're not running in an environment with Neo4j
	t.Skip("Skipping test that requires Neo4j")

	// Create a new frontend
	f, err := NewFrontend(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, f)
	defer f.Close()

	// Create a temporary directory for testing
	tempDir := t.TempDir()
	os.Chdir(tempDir)

	// Create the files
	err = f.createIndexHTML()
	assert.NoError(t, err)
	assert.FileExists(t, "public/index.html")

	err = f.createStyleCSS()
	assert.NoError(t, err)
	assert.FileExists(t, "public/css/style.css")

	err = f.createScriptJS()
	assert.NoError(t, err)
	assert.FileExists(t, "public/js/script.js")
}

func TestHandleIndex(t *testing.T) {
	// Create a test configuration
	cfg := &config.Neo4jConfig{
		URI:               "bolt://localhost:7687",
		User:              "neo4j",
		Password:          "password",
		MaxRetries:        3,
		RetryIntervalSecs: 1,
	}

	// Skip the test if we're not running in an environment with Neo4j
	t.Skip("Skipping test that requires Neo4j")

	// Create a new frontend
	f, err := NewFrontend(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, f)
	defer f.Close()

	// Create a request to the index page
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(f.handleIndex)
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestHandleStats(t *testing.T) {
	// Create a test configuration
	cfg := &config.Neo4jConfig{
		URI:               "bolt://localhost:7687",
		User:              "neo4j",
		Password:          "password",
		MaxRetries:        3,
		RetryIntervalSecs: 1,
	}

	// Skip the test if we're not running in an environment with Neo4j
	t.Skip("Skipping test that requires Neo4j")

	// Create a new frontend
	f, err := NewFrontend(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, f)
	defer f.Close()

	// Create a request to the stats API
	req, err := http.NewRequest("GET", "/api/stats", nil)
	assert.NoError(t, err)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(f.handleStats)
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the content type
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestHandleGraph(t *testing.T) {
	// Create a test configuration
	cfg := &config.Neo4jConfig{
		URI:               "bolt://localhost:7687",
		User:              "neo4j",
		Password:          "password",
		MaxRetries:        3,
		RetryIntervalSecs: 1,
	}

	// Skip the test if we're not running in an environment with Neo4j
	t.Skip("Skipping test that requires Neo4j")

	// Create a new frontend
	f, err := NewFrontend(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, f)
	defer f.Close()

	// Create a request to the graph API
	req, err := http.NewRequest("GET", "/api/graph", nil)
	assert.NoError(t, err)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(f.handleGraph)
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the content type
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
} 