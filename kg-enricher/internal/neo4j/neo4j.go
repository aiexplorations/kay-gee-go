package neo4j

import (
	"fmt"
	"log"
	"time"

	"kg-enricher/internal/config"
	"kg-enricher/internal/models"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// Default configuration values
const (
	DefaultMaxRetries    = 5
	DefaultRetryInterval = 5 * time.Second
	DefaultMaxBackoff    = 30 * time.Second
)

// SetupNeo4jConnection establishes a connection to the Neo4j database with retry logic to handle connection failures.
func SetupNeo4jConnection(cfg *config.Neo4jConfig) (neo4j.Driver, error) {
	log.Printf("Connecting to Neo4j at %s", cfg.URI)
	
	if cfg.URI == "" {
		return nil, fmt.Errorf("Neo4j URI is not set")
	}
	
	if cfg.User == "" {
		return nil, fmt.Errorf("Neo4j user is not set")
	}
	
	if cfg.Password == "" {
		return nil, fmt.Errorf("Neo4j password is not set")
	}
	
	var driver neo4j.Driver
	var err error
	
	maxRetries := cfg.MaxRetries
	if maxRetries <= 0 {
		maxRetries = DefaultMaxRetries
	}
	
	retryInterval := cfg.RetryInterval
	if retryInterval <= 0 {
		retryInterval = DefaultRetryInterval
	}
	
	// Try to connect with retries
	for attempt := 1; attempt <= maxRetries; attempt++ {
		driver, err = neo4j.NewDriver(cfg.URI, neo4j.BasicAuth(cfg.User, cfg.Password, ""))
		if err == nil {
			// Test the connection
			err = driver.VerifyConnectivity()
			if err == nil {
				log.Printf("Successfully connected to Neo4j at %s", cfg.URI)
				return driver, nil
			}
		}
		
		log.Printf("Failed to create Neo4j driver (attempt %d/%d): %v", attempt, maxRetries, err)
		
		if attempt < maxRetries {
			time.Sleep(retryInterval)
		}
	}
	
	return nil, fmt.Errorf("failed to connect to Neo4j after %d attempts: %v", maxRetries, err)
}

// GetRandomNodes retrieves a batch of random nodes from the Neo4j database
func GetRandomNodes(driver neo4j.Driver, batchSize int) ([]models.Node, error) {
	if driver == nil {
		return nil, fmt.Errorf("Neo4j driver is nil")
	}
	
	if batchSize <= 0 {
		return nil, fmt.Errorf("batch size must be greater than 0")
	}
	
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	
	result, err := session.Run(
		`MATCH (n:Concept) 
		 RETURN id(n) AS id, n.name AS name, labels(n)[0] AS label 
		 ORDER BY rand() 
		 LIMIT $batchSize`,
		map[string]interface{}{
			"batchSize": batchSize,
		},
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to query random nodes: %w", err)
	}
	
	var nodes []models.Node
	for result.Next() {
		record := result.Record()
		
		id, _ := record.Get("id")
		name, _ := record.Get("name")
		label, _ := record.Get("label")
		
		node := models.Node{
			ID:    id.(int64),
			Name:  name.(string),
			Label: label.(string),
		}
		
		nodes = append(nodes, node)
	}
	
	if err = result.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating results: %w", err)
	}
	
	return nodes, nil
}

// CheckExistingRelationship checks if a relationship already exists between two nodes
func CheckExistingRelationship(driver neo4j.Driver, source, target string) (bool, error) {
	if driver == nil {
		return false, fmt.Errorf("Neo4j driver is nil")
	}
	
	if source == "" || target == "" {
		return false, fmt.Errorf("source and target must not be empty")
	}
	
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	
	result, err := session.Run(
		`MATCH (a:Concept {name: $source})-[r]-(b:Concept {name: $target}) 
		 RETURN count(r) > 0 AS exists`,
		map[string]interface{}{
			"source": source,
			"target": target,
		},
	)
	
	if err != nil {
		return false, fmt.Errorf("failed to check existing relationship: %w", err)
	}
	
	if result.Next() {
		record := result.Record()
		exists, _ := record.Get("exists")
		return exists.(bool), nil
	}
	
	if err = result.Err(); err != nil {
		return false, fmt.Errorf("error while checking relationship: %w", err)
	}
	
	return false, nil
}

// CreateRelationship creates a relationship between two concepts in the Neo4j database
func CreateRelationship(driver neo4j.Driver, from, to, relation string) error {
	if driver == nil {
		return fmt.Errorf("Neo4j driver is nil")
	}
	
	if from == "" {
		return fmt.Errorf("from concept cannot be empty")
	}
	
	if to == "" {
		return fmt.Errorf("to concept cannot be empty")
	}
	
	if relation == "" {
		return fmt.Errorf("relation cannot be empty")
	}
	
	// Convert relation to a valid relationship type (lowercase with underscores)
	relationType := formatRelationType(relation)
	
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			`MATCH (a:Concept {name: $from}), (b:Concept {name: $to}) 
			 MERGE (a)-[r:`+relationType+` {description: $relation}]->(b) 
			 RETURN r`,
			map[string]interface{}{
				"from":     from,
				"to":       to,
				"relation": relation,
			},
		)
		
		if err != nil {
			return nil, err
		}
		
		// Check if the relationship was created
		if result.Next() {
			return true, nil
		}
		
		return false, result.Err()
	})
	
	if err != nil {
		return fmt.Errorf("failed to create relationship: %w", err)
	}
	
	log.Printf("Successfully created relationship: %s -[%s]-> %s", from, relation, to)
	return nil
}

// QueryAllConcepts retrieves all concepts from the Neo4j database
func QueryAllConcepts(driver neo4j.Driver) ([]string, error) {
	if driver == nil {
		return nil, fmt.Errorf("Neo4j driver is nil")
	}
	
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	
	result, err := session.Run(
		`MATCH (n:Concept) RETURN n.name AS name ORDER BY n.name`,
		nil,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to query concepts: %w", err)
	}
	
	var concepts []string
	for result.Next() {
		record := result.Record()
		name, _ := record.Get("name")
		concepts = append(concepts, name.(string))
	}
	
	if err = result.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating results: %w", err)
	}
	
	return concepts, nil
}

// QueryRelationships retrieves all relationships from the Neo4j database
func QueryRelationships(driver neo4j.Driver) ([]models.Relationship, error) {
	if driver == nil {
		return nil, fmt.Errorf("Neo4j driver is nil")
	}
	
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	
	result, err := session.Run(
		`MATCH (a:Concept)-[r]->(b:Concept) 
		 RETURN a.name AS source, type(r) AS type, b.name AS target, r.description AS description`,
		nil,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to query relationships: %w", err)
	}
	
	var relationships []models.Relationship
	for result.Next() {
		record := result.Record()
		
		source, _ := record.Get("source")
		relType, _ := record.Get("type")
		target, _ := record.Get("target")
		
		description, _ := record.Get("description")
		var desc string
		if description != nil {
			desc = description.(string)
		} else {
			desc = relType.(string)
		}
		
		relationship := models.Relationship{
			Source:      source.(string),
			Type:        relType.(string),
			Target:      target.(string),
			Description: desc,
		}
		
		relationships = append(relationships, relationship)
	}
	
	if err = result.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating results: %w", err)
	}
	
	return relationships, nil
}

// QueryRandomConceptPairs retrieves random pairs of concepts from the database
// that don't already have a relationship between them.
func QueryRandomConceptPairs(driver neo4j.Driver, batchSize int) ([][2]string, error) {
	if batchSize <= 0 {
		batchSize = 10 // Default batch size
	}

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	// Query to find random pairs of concepts that don't have a relationship
	query := `
		MATCH (a:Concept), (b:Concept)
		WHERE a <> b
		AND NOT (a)-[]->(b)
		AND NOT (b)-[]->(a)
		WITH a, b, rand() AS r
		ORDER BY r
		LIMIT $batchSize
		RETURN a.name AS source, b.name AS target
	`

	result, err := session.Run(query, map[string]interface{}{
		"batchSize": batchSize,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query random concept pairs: %w", err)
	}

	pairs := make([][2]string, 0, batchSize)
	for result.Next() {
		record := result.Record()
		source, _ := record.Get("source")
		target, _ := record.Get("target")
		
		sourceStr, ok1 := source.(string)
		targetStr, ok2 := target.(string)
		
		if ok1 && ok2 {
			pairs = append(pairs, [2]string{sourceStr, targetStr})
		}
	}

	if err = result.Err(); err != nil {
		return nil, fmt.Errorf("error while processing query results: %w", err)
	}

	return pairs, nil
}

// Helper function to format relation type for Neo4j
func formatRelationType(relation string) string {
	// Convert spaces to underscores and make lowercase
	// This is a simplified version, you might want to add more transformations
	result := ""
	for _, c := range relation {
		if c == ' ' {
			result += "_"
		} else {
			result += string(c)
		}
	}
	return result
} 