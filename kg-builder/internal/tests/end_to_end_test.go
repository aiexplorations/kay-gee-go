package tests

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"kg-builder/internal/config"
	"kg-builder/internal/graph"
	"kg-builder/internal/llm"
	"kg-builder/internal/neo4j"

	neo4jdriver "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func TestEndToEndWorkflow(t *testing.T) {
	// Skip this test if Neo4j is not available
	t.Skip("Skipping test that requires Neo4j")
	
	// Create a temporary directory for the cache
	tempDir, err := ioutil.TempDir("", "llm-cache-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create a test configuration
	neo4jConfig := &config.Neo4jConfig{
		URI:           "bolt://localhost:7687",
		User:          "neo4j",
		Password:      "password",
		MaxRetries:    3,
		RetryInterval: 2 * time.Second,
	}
	
	llmConfig := &config.LLMConfig{
		URL:      "http://localhost:11434/api/generate",
		Model:    "qwen2.5:3b",
		CacheDir: tempDir,
	}
	
	graphConfig := &config.GraphConfig{
		SeedConcept:         "Artificial Intelligence",
		MaxNodes:            5, // Small number for testing
		Timeout:             5 * time.Minute,
		WorkerCount:         2,
		RandomRelationships: 2,
		Concurrency:         2,
	}
	
	// Initialize the LLM service
	err = llm.Initialize(llmConfig)
	if err != nil {
		t.Fatalf("Failed to initialize LLM service: %v", err)
	}
	
	// Setup the Neo4j connection
	driver, err := neo4j.SetupNeo4jConnection(neo4jConfig)
	if err != nil {
		t.Fatalf("Failed to setup Neo4j connection: %v", err)
	}
	defer driver.Close()
	
	// Clean up any existing data
	session := driver.NewSession(neo4jdriver.SessionConfig{
		AccessMode: neo4jdriver.AccessModeWrite,
	})
	_, err = session.Run("MATCH (n) DETACH DELETE n", nil)
	if err != nil {
		t.Fatalf("Failed to clean up existing data: %v", err)
	}
	session.Close()
	
	// Step 1: Build the graph
	t.Log("Step 1: Building the graph")
	gb := graph.NewGraphBuilder(
		driver,
		llm.GetRelatedConcepts,
		llm.MineRelationship,
		graphConfig,
	)
	
	err = gb.BuildGraph(
		graphConfig.SeedConcept,
		graphConfig.MaxNodes,
		graphConfig.Timeout,
	)
	
	if err != nil {
		t.Fatalf("Failed to build graph: %v", err)
	}
	
	// Check that we processed some nodes
	nodeCount := gb.GetNodeCount()
	if nodeCount == 0 {
		t.Errorf("Expected to process some nodes, got 0")
	}
	
	t.Logf("Built graph with %d nodes", nodeCount)
	
	// Step 2: Mine random relationships
	t.Log("Step 2: Mining random relationships")
	err = gb.MineRandomRelationships(
		graphConfig.RandomRelationships,
		graphConfig.Concurrency,
	)
	
	if err != nil {
		t.Fatalf("Failed to mine random relationships: %v", err)
	}
	
	// Step 3: Query the graph
	t.Log("Step 3: Querying the graph")
	concepts, err := neo4j.QueryConcepts(driver)
	if err != nil {
		t.Fatalf("Failed to query concepts: %v", err)
	}
	
	if len(concepts) == 0 {
		t.Errorf("Expected to get some concepts, got 0")
	}
	
	t.Logf("Found %d concepts in the graph", len(concepts))
	
	relationships, err := neo4j.QueryRelationships(driver)
	if err != nil {
		t.Fatalf("Failed to query relationships: %v", err)
	}
	
	if len(relationships) == 0 {
		t.Errorf("Expected to get some relationships, got 0")
	}
	
	t.Logf("Found %d relationships in the graph", len(relationships))
	
	// Step 4: Enrich the graph using the enricher
	t.Log("Step 4: Enriching the graph")
	
	// Import the enricher package dynamically to avoid circular dependencies
	// In a real-world scenario, you would have a separate binary for the enricher
	// and would call it via a command-line interface or API
	
	// For this test, we'll simulate the enricher by creating more relationships
	// between random pairs of concepts
	
	// Get all concepts
	session = driver.NewSession(neo4jdriver.SessionConfig{
		AccessMode: neo4jdriver.AccessModeRead,
	})
	result, err := session.Run("MATCH (n:Concept) RETURN n.name AS name", nil)
	if err != nil {
		t.Fatalf("Failed to query concepts: %v", err)
	}
	
	var conceptNames []string
	for result.Next() {
		record := result.Record()
		name, _ := record.Get("name")
		conceptNames = append(conceptNames, name.(string))
	}
	session.Close()
	
	// Create relationships between random pairs of concepts
	if len(conceptNames) >= 2 {
		session = driver.NewSession(neo4jdriver.SessionConfig{
			AccessMode: neo4jdriver.AccessModeWrite,
		})
		for i := 0; i < 3; i++ {
			source := conceptNames[i%len(conceptNames)]
			target := conceptNames[(i+1)%len(conceptNames)]
			
			if source == target {
				continue
			}
			
			// Check if a relationship already exists
			checkResult, err := session.Run("MATCH (s:Concept {name: $source})-[r]->(t:Concept {name: $target}) RETURN count(r) AS count", map[string]interface{}{
				"source": source,
				"target": target,
			})
			if err != nil {
				t.Fatalf("Failed to check existing relationship: %v", err)
			}
			
			var count int64
			if checkResult.Next() {
				record := checkResult.Record()
				countValue, _ := record.Get("count")
				count = countValue.(int64)
			}
			
			if count == 0 {
				// Create a new relationship
				_, err = session.Run("MATCH (s:Concept {name: $source}), (t:Concept {name: $target}) CREATE (s)-[r:EnrichedRelation]->(t) RETURN r", map[string]interface{}{
					"source": source,
					"target": target,
				})
				if err != nil {
					t.Fatalf("Failed to create relationship: %v", err)
				}
				
				t.Logf("Created enriched relationship: %s -> %s", source, target)
			}
		}
		session.Close()
	}
	
	// Step 5: Query the enriched graph
	t.Log("Step 5: Querying the enriched graph")
	relationships, err = neo4j.QueryRelationships(driver)
	if err != nil {
		t.Fatalf("Failed to query relationships: %v", err)
	}
	
	t.Logf("Found %d relationships in the enriched graph", len(relationships))
	
	// Clean up the test data
	session = driver.NewSession(neo4jdriver.SessionConfig{
		AccessMode: neo4jdriver.AccessModeWrite,
	})
	_, err = session.Run("MATCH (n) DETACH DELETE n", nil)
	if err != nil {
		t.Fatalf("Failed to clean up test data: %v", err)
	}
	session.Close()
} 