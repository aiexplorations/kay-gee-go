package neo4j

import (
	"testing"

	"kg-enricher/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func TestSetupNeo4jConnection(t *testing.T) {
	// Skip this test if Neo4j is not available
	t.Skip("Skipping test that requires Neo4j")
	
	// Create a config
	cfg := &config.Neo4jConfig{
		URI:      "bolt://localhost:7687",
		Username: "neo4j",
		Password: "password",
	}
	
	// Setup the Neo4j connection
	driver, err := SetupNeo4jConnection(cfg)
	
	// Assert that there was no error
	assert.NoError(t, err)
	
	// Assert that the driver was created
	assert.NotNil(t, driver)
	
	// Close the driver
	err = driver.Close()
	assert.NoError(t, err)
}

func TestQueryRandomConceptPairs(t *testing.T) {
	// Skip this test if Neo4j is not available
	t.Skip("Skipping test that requires Neo4j")
	
	// Create a config
	cfg := &config.Neo4jConfig{
		URI:      "bolt://localhost:7687",
		Username: "neo4j",
		Password: "password",
	}
	
	// Setup the Neo4j connection
	driver, err := SetupNeo4jConnection(cfg)
	assert.NoError(t, err)
	defer driver.Close()
	
	// Query random concept pairs
	pairs, err := QueryRandomConceptPairs(driver, 5)
	
	// Assert that there was no error
	assert.NoError(t, err)
	
	// Assert that the pairs were returned
	assert.NotNil(t, pairs)
}

func TestCreateRelationship(t *testing.T) {
	// Skip this test if Neo4j is not available
	t.Skip("Skipping test that requires Neo4j")
	
	// Create a config
	cfg := &config.Neo4jConfig{
		URI:      "bolt://localhost:7687",
		Username: "neo4j",
		Password: "password",
	}
	
	// Setup the Neo4j connection
	driver, err := SetupNeo4jConnection(cfg)
	assert.NoError(t, err)
	defer driver.Close()
	
	// Create a relationship
	err = CreateRelationship(driver, "Test Concept 1", "Test Concept 2", "TEST_RELATIONSHIP")
	
	// Assert that there was no error
	assert.NoError(t, err)
}

func TestQueryAllConcepts(t *testing.T) {
	// Skip this test if Neo4j is not available
	t.Skip("Skipping test that requires Neo4j")
	
	// Create a config
	cfg := &config.Neo4jConfig{
		URI:      "bolt://localhost:7687",
		Username: "neo4j",
		Password: "password",
	}
	
	// Setup the Neo4j connection
	driver, err := SetupNeo4jConnection(cfg)
	assert.NoError(t, err)
	defer driver.Close()
	
	// Query all concepts
	concepts, err := QueryAllConcepts(driver)
	
	// Assert that there was no error
	assert.NoError(t, err)
	
	// Assert that the concepts were returned
	assert.NotNil(t, concepts)
}

func TestQueryRelationships(t *testing.T) {
	// Skip this test if Neo4j is not available
	t.Skip("Skipping test that requires Neo4j")
	
	// Create a config
	cfg := &config.Neo4jConfig{
		URI:      "bolt://localhost:7687",
		Username: "neo4j",
		Password: "password",
	}
	
	// Setup the Neo4j connection
	driver, err := SetupNeo4jConnection(cfg)
	assert.NoError(t, err)
	defer driver.Close()
	
	// Query relationships
	relationships, err := QueryRelationships(driver)
	
	// Assert that there was no error
	assert.NoError(t, err)
	
	// Assert that the relationships were returned
	assert.NotNil(t, relationships)
}

func TestGetRandomNodes(t *testing.T) {
	// Skip this test if Neo4j is not available
	t.Skip("Skipping test that requires Neo4j")
	
	// Create a test configuration
	cfg := &config.Neo4jConfig{
		URI:      "bolt://localhost:7687",
		Username: "neo4j",
		Password: "password",
	}
	
	// Setup the Neo4j connection
	driver, err := SetupNeo4jConnection(cfg)
	assert.NoError(t, err)
	defer driver.Close()
	
	// Create some test nodes
	testNodes := []string{
		"TestRandomNode1",
		"TestRandomNode2",
		"TestRandomNode3",
		"TestRandomNode4",
		"TestRandomNode5",
	}
	
	// Create the test nodes in Neo4j
	session := driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close()
	
	for _, nodeName := range testNodes {
		_, err := session.Run("CREATE (n:Concept {name: $name}) RETURN n", map[string]interface{}{
			"name": nodeName,
		})
		assert.NoError(t, err)
	}
	
	// Get random nodes
	count := 3
	nodes, err := GetRandomNodes(driver, count)
	
	// Check that there was no error
	assert.NoError(t, err)
	
	// Check that we got the expected number of nodes
	assert.LessOrEqual(t, len(nodes), count)
	
	// Check that the nodes have the expected fields
	for _, node := range nodes {
		assert.NotEmpty(t, node.Name)
		assert.Equal(t, "Concept", node.Label)
	}
	
	// Clean up the test nodes
	_, err = session.Run("MATCH (n) WHERE n.name IN $names DETACH DELETE n", map[string]interface{}{
		"names": testNodes,
	})
	assert.NoError(t, err)
}

func TestCheckExistingRelationship(t *testing.T) {
	// Skip this test if Neo4j is not available
	t.Skip("Skipping test that requires Neo4j")
	
	// Create a test configuration
	cfg := &config.Neo4jConfig{
		URI:      "bolt://localhost:7687",
		Username: "neo4j",
		Password: "password",
	}
	
	// Setup the Neo4j connection
	driver, err := SetupNeo4jConnection(cfg)
	assert.NoError(t, err)
	defer driver.Close()
	
	// Create some test nodes and relationships
	source := "TestSource"
	target := "TestTarget"
	
	// Create the test nodes and relationship in Neo4j
	session := driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close()
	
	_, err = session.Run("CREATE (s:Concept {name: $source}), (t:Concept {name: $target}) RETURN s, t", map[string]interface{}{
		"source": source,
		"target": target,
	})
	assert.NoError(t, err)
	
	_, err = session.Run("MATCH (s:Concept {name: $source}), (t:Concept {name: $target}) CREATE (s)-[r:TestRelation]->(t) RETURN r", map[string]interface{}{
		"source": source,
		"target": target,
	})
	assert.NoError(t, err)
	
	// Check that the relationship exists
	exists, err := CheckExistingRelationship(driver, source, target)
	
	// Check that there was no error
	assert.NoError(t, err)
	
	// Check that the relationship exists
	assert.True(t, exists)
	
	// Check a non-existent relationship
	exists, err = CheckExistingRelationship(driver, "NonExistentSource", "NonExistentTarget")
	
	// Check that there was no error
	assert.NoError(t, err)
	
	// Check that the relationship does not exist
	assert.False(t, exists)
	
	// Clean up the test nodes and relationship
	_, err = session.Run("MATCH (n) WHERE n.name IN [$source, $target] DETACH DELETE n", map[string]interface{}{
		"source": source,
		"target": target,
	})
	assert.NoError(t, err)
} 