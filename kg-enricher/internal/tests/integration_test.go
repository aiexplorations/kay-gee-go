package tests

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"kg-enricher/internal/config"
	"kg-enricher/internal/enricher"
	"kg-enricher/internal/llm"
	"kg-enricher/internal/models"
	"kg-enricher/internal/neo4j"

	neo4jdriver "github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/stretchr/testify/assert"
)

func TestEnrichmentWorkflow(t *testing.T) {
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
		URI:      "bolt://localhost:7687",
		Username: "neo4j",
		Password: "password",
	}
	
	llmConfig := &config.LLMConfig{
		URL:      "http://localhost:11434/api/generate",
		Model:    "qwen2.5:3b",
		CacheDir: tempDir,
	}
	
	enricherConfig := &config.EnricherConfig{
		BatchSize:        5,
		Interval:         time.Second * 1,
		MaxRelationships: 10,
		Concurrency:      2,
	}
	
	// Initialize the LLM service
	err = llm.Initialize(llmConfig)
	assert.NoError(t, err)
	
	// Setup the Neo4j connection
	driver, err := neo4j.SetupNeo4jConnection(neo4jConfig)
	assert.NoError(t, err)
	defer driver.Close()
	
	// Create some test nodes
	testNodes := []string{
		"TestNode1",
		"TestNode2",
		"TestNode3",
		"TestNode4",
		"TestNode5",
		"TestNode6",
		"TestNode7",
		"TestNode8",
		"TestNode9",
		"TestNode10",
	}
	
	// Create the test nodes in Neo4j
	session := driver.NewSession(neo4jdriver.SessionConfig{
		AccessMode: neo4jdriver.AccessModeWrite,
	})
	defer session.Close()
	
	for _, nodeName := range testNodes {
		_, err := session.Run("CREATE (n:Concept {name: $name}) RETURN n", map[string]interface{}{
			"name": nodeName,
		})
		assert.NoError(t, err)
	}
	
	// Create an enricher
	e := enricher.NewEnricher(driver, enricherConfig)
	
	// Run the enricher once
	err = e.RunOnce(5)
	assert.NoError(t, err)
	
	// Check that the stats were updated
	stats := e.GetStats()
	assert.GreaterOrEqual(t, stats.TotalBatches, 1)
	assert.GreaterOrEqual(t, stats.TotalPairsProcessed, 1)
	
	// Query all relationships (just to verify no errors)
	_, err = neo4j.QueryRelationships(driver)
	assert.NoError(t, err)
	
	// Clean up the test data
	_, err = session.Run("MATCH (n) DETACH DELETE n", nil)
	assert.NoError(t, err)
}

func TestEnrichmentWithExistingGraph(t *testing.T) {
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
		URI:      "bolt://localhost:7687",
		Username: "neo4j",
		Password: "password",
	}
	
	llmConfig := &config.LLMConfig{
		URL:      "http://localhost:11434/api/generate",
		Model:    "qwen2.5:3b",
		CacheDir: tempDir,
	}
	
	enricherConfig := &config.EnricherConfig{
		BatchSize:        5,
		Interval:         time.Second * 1,
		MaxRelationships: 10,
		Concurrency:      2,
	}
	
	// Initialize the LLM service
	err = llm.Initialize(llmConfig)
	assert.NoError(t, err)
	
	// Setup the Neo4j connection
	driver, err := neo4j.SetupNeo4jConnection(neo4jConfig)
	assert.NoError(t, err)
	defer driver.Close()
	
	// Create some test nodes and relationships
	testNodes := []models.Node{
		{Name: "Machine Learning", Label: "Concept"},
		{Name: "Artificial Intelligence", Label: "Concept"},
		{Name: "Neural Networks", Label: "Concept"},
		{Name: "Deep Learning", Label: "Concept"},
		{Name: "Natural Language Processing", Label: "Concept"},
	}
	
	testRelationships := []models.Relationship{
		{Source: "Machine Learning", Target: "Artificial Intelligence", Type: "IsA"},
		{Source: "Neural Networks", Target: "Machine Learning", Type: "UsedIn"},
	}
	
	// Create the test nodes in Neo4j
	session := driver.NewSession(neo4jdriver.SessionConfig{
		AccessMode: neo4jdriver.AccessModeWrite,
	})
	defer session.Close()
	
	for _, node := range testNodes {
		_, err := session.Run("CREATE (n:Concept {name: $name}) RETURN n", map[string]interface{}{
			"name": node.Name,
		})
		assert.NoError(t, err)
	}
	
	// Create the test relationships in Neo4j
	for _, rel := range testRelationships {
		_, err := session.Run("MATCH (s:Concept {name: $source}), (t:Concept {name: $target}) CREATE (s)-[r:"+rel.Type+"]->(t) RETURN r", map[string]interface{}{
			"source": rel.Source,
			"target": rel.Target,
		})
		assert.NoError(t, err)
	}
	
	// Create an enricher
	e := enricher.NewEnricher(driver, enricherConfig)
	
	// Run the enricher once
	err = e.RunOnce(5)
	assert.NoError(t, err)
	
	// Check that the stats were updated
	stats := e.GetStats()
	assert.GreaterOrEqual(t, stats.TotalBatches, 1)
	assert.GreaterOrEqual(t, stats.TotalPairsProcessed, 1)
	
	// Query all relationships (just to verify no errors)
	_, err = neo4j.QueryRelationships(driver)
	assert.NoError(t, err)
	
	// Clean up the test data
	_, err = session.Run("MATCH (n) DETACH DELETE n", nil)
	assert.NoError(t, err)
}

func TestConcurrentEnrichment(t *testing.T) {
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
		URI:      "bolt://localhost:7687",
		Username: "neo4j",
		Password: "password",
	}
	
	llmConfig := &config.LLMConfig{
		URL:      "http://localhost:11434/api/generate",
		Model:    "qwen2.5:3b",
		CacheDir: tempDir,
	}
	
	enricherConfig := &config.EnricherConfig{
		BatchSize:        10,
		Interval:         time.Second * 1,
		MaxRelationships: 20,
		Concurrency:      5, // High concurrency for testing
	}
	
	// Initialize the LLM service
	err = llm.Initialize(llmConfig)
	assert.NoError(t, err)
	
	// Setup the Neo4j connection
	driver, err := neo4j.SetupNeo4jConnection(neo4jConfig)
	assert.NoError(t, err)
	defer driver.Close()
	
	// Create some test nodes
	testNodes := []string{
		"TestNode1", "TestNode2", "TestNode3", "TestNode4", "TestNode5",
		"TestNode6", "TestNode7", "TestNode8", "TestNode9", "TestNode10",
		"TestNode11", "TestNode12", "TestNode13", "TestNode14", "TestNode15",
		"TestNode16", "TestNode17", "TestNode18", "TestNode19", "TestNode20",
	}
	
	// Create the test nodes in Neo4j
	session := driver.NewSession(neo4jdriver.SessionConfig{
		AccessMode: neo4jdriver.AccessModeWrite,
	})
	defer session.Close()
	
	for _, nodeName := range testNodes {
		_, err := session.Run("CREATE (n:Concept {name: $name}) RETURN n", map[string]interface{}{
			"name": nodeName,
		})
		assert.NoError(t, err)
	}
	
	// Create an enricher
	e := enricher.NewEnricher(driver, enricherConfig)
	
	// Start the enricher
	err = e.Start()
	assert.NoError(t, err)
	
	// Wait for the enricher to process some batches
	time.Sleep(time.Second * 3)
	
	// Stop the enricher
	err = e.Stop()
	assert.NoError(t, err)
	
	// Check that the stats were updated
	stats := e.GetStats()
	assert.GreaterOrEqual(t, stats.TotalBatches, 1)
	assert.GreaterOrEqual(t, stats.TotalPairsProcessed, 1)
	
	// Query all relationships (just to verify no errors)
	_, err = neo4j.QueryRelationships(driver)
	assert.NoError(t, err)
	
	// Clean up the test data
	_, err = session.Run("MATCH (n) DETACH DELETE n", nil)
	assert.NoError(t, err)
} 