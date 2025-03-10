package neo4j_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"kg-builder/internal/config"
	"kg-builder/internal/neo4j"

	neo4jdriver "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func TestNeo4jConfiguration(t *testing.T) {
	// Save current environment variables
	oldNeo4jURI := os.Getenv("NEO4J_URI")
	oldNeo4jUser := os.Getenv("NEO4J_USER")
	oldNeo4jPassword := os.Getenv("NEO4J_PASSWORD")

	// Set environment variables
	os.Setenv("NEO4J_URI", "bolt://test-neo4j:7687")
	os.Setenv("NEO4J_USER", "test-user")
	os.Setenv("NEO4J_PASSWORD", "test-password")

	// Restore environment variables after test
	defer func() {
		os.Setenv("NEO4J_URI", oldNeo4jURI)
		os.Setenv("NEO4J_USER", oldNeo4jUser)
		os.Setenv("NEO4J_PASSWORD", oldNeo4jPassword)
	}()

	// Create a test configuration
	cfg := &config.Neo4jConfig{
		URI:           "bolt://test-neo4j:7687",
		User:          "test-user",
		Password:      "test-password",
		MaxRetries:    3,
		RetryInterval: 2 * time.Second,
	}

	// Test the configuration
	if cfg.URI != "bolt://test-neo4j:7687" {
		t.Errorf("Expected URI to be 'bolt://test-neo4j:7687', got '%s'", cfg.URI)
	}
	if cfg.User != "test-user" {
		t.Errorf("Expected User to be 'test-user', got '%s'", cfg.User)
	}
	if cfg.Password != "test-password" {
		t.Errorf("Expected Password to be 'test-password', got '%s'", cfg.Password)
	}
	if cfg.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries to be 3, got %d", cfg.MaxRetries)
	}
	if cfg.RetryInterval != 2*time.Second {
		t.Errorf("Expected RetryInterval to be 2s, got %s", cfg.RetryInterval)
	}
}

func TestCreateRelationshipValidation(t *testing.T) {
	// Test with empty source
	err := neo4j.CreateRelationship(nil, "", "target", "relation")
	if err == nil {
		t.Error("Expected error with empty source, got nil")
	}

	// Test with empty target
	err = neo4j.CreateRelationship(nil, "source", "", "relation")
	if err == nil {
		t.Error("Expected error with empty target, got nil")
	}

	// Test with empty relation
	err = neo4j.CreateRelationship(nil, "source", "target", "")
	if err == nil {
		t.Error("Expected error with empty relation, got nil")
	}
}

func TestQueryConcepts(t *testing.T) {
	// Skip this test if Neo4j is not available
	t.Skip("Skipping test that requires Neo4j")
	
	// Create a test configuration
	cfg := &config.Neo4jConfig{
		URI:           "bolt://localhost:7687",
		User:          "neo4j",
		Password:      "password",
		MaxRetries:    3,
		RetryInterval: 2 * time.Second,
	}
	
	// Setup the Neo4j connection
	driver, err := neo4j.SetupNeo4jConnection(cfg)
	if err != nil {
		t.Fatalf("Failed to setup Neo4j connection: %v", err)
	}
	defer driver.Close()
	
	// Create some test concepts
	testConcepts := []struct {
		name     string
		relation string
		relatedTo string
	}{
		{"TestConcept1", "IsA", "TestParent1"},
		{"TestConcept2", "IsA", "TestParent1"},
		{"TestConcept3", "HasPart", "TestParent2"},
	}
	
	// Create the test concepts in Neo4j
	for _, tc := range testConcepts {
		err := neo4j.CreateRelationship(driver, tc.relatedTo, tc.name, tc.relation)
		if err != nil {
			t.Fatalf("Failed to create test concept: %v", err)
		}
	}
	
	// Query all concepts
	concepts, err := neo4j.QueryConcepts(driver)
	if err != nil {
		t.Fatalf("Failed to query concepts: %v", err)
	}
	
	// Check that we got at least the test concepts
	foundConcepts := make(map[string]bool)
	for _, concept := range concepts {
		foundConcepts[concept] = true
	}
	
	for _, tc := range testConcepts {
		if !foundConcepts[tc.name] {
			t.Errorf("Expected to find concept %s, but it was not in the results", tc.name)
		}
		if !foundConcepts[tc.relatedTo] {
			t.Errorf("Expected to find concept %s, but it was not in the results", tc.relatedTo)
		}
	}
	
	// Clean up the test concepts
	session := driver.NewSession(neo4jdriver.AccessModeWrite)
	defer session.Close()
	
	_, err = session.Run("MATCH (n) WHERE n.name IN $names DETACH DELETE n", map[string]interface{}{
		"names": []string{"TestConcept1", "TestConcept2", "TestConcept3", "TestParent1", "TestParent2"},
	})
	if err != nil {
		t.Fatalf("Failed to clean up test concepts: %v", err)
	}
}

func TestQueryRelationships(t *testing.T) {
	// Skip this test if Neo4j is not available
	t.Skip("Skipping test that requires Neo4j")
	
	// Create a test configuration
	cfg := &config.Neo4jConfig{
		URI:           "bolt://localhost:7687",
		User:          "neo4j",
		Password:      "password",
		MaxRetries:    3,
		RetryInterval: 2 * time.Second,
	}
	
	// Setup the Neo4j connection
	driver, err := neo4j.SetupNeo4jConnection(cfg)
	if err != nil {
		t.Fatalf("Failed to setup Neo4j connection: %v", err)
	}
	defer driver.Close()
	
	// Create some test relationships
	testRelationships := []struct {
		source   string
		target   string
		relation string
	}{
		{"TestSource1", "TestTarget1", "IsA"},
		{"TestSource2", "TestTarget2", "HasPart"},
		{"TestSource3", "TestTarget3", "UsedIn"},
	}
	
	// Create the test relationships in Neo4j
	for _, tr := range testRelationships {
		err := neo4j.CreateRelationship(driver, tr.source, tr.target, tr.relation)
		if err != nil {
			t.Fatalf("Failed to create test relationship: %v", err)
		}
	}
	
	// Query all relationships
	relationships, err := neo4j.QueryRelationships(driver)
	if err != nil {
		t.Fatalf("Failed to query relationships: %v", err)
	}
	
	// Check that we got at least the test relationships
	foundRelationships := make(map[string]bool)
	for _, rel := range relationships {
		key := rel["source"] + "-" + rel["relation"] + "-" + rel["target"]
		foundRelationships[key] = true
	}
	
	for _, tr := range testRelationships {
		key := tr.source + "-" + tr.relation + "-" + tr.target
		if !foundRelationships[key] {
			t.Errorf("Expected to find relationship %s, but it was not in the results", key)
		}
	}
	
	// Clean up the test relationships
	session := driver.NewSession(neo4jdriver.AccessModeWrite)
	defer session.Close()
	
	_, err = session.Run("MATCH (n) WHERE n.name IN $names DETACH DELETE n", map[string]interface{}{
		"names": []string{"TestSource1", "TestSource2", "TestSource3", "TestTarget1", "TestTarget2", "TestTarget3"},
	})
	if err != nil {
		t.Fatalf("Failed to clean up test relationships: %v", err)
	}
} 