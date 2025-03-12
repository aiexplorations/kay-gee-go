package neo4j

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	apperrors "kg-builder/internal/errors"
	"kg-builder/internal/config"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// Default configuration values
const (
	DefaultMaxRetries    = 5
	DefaultRetryInterval = 5 * time.Second
	DefaultMaxBackoff    = 30 * time.Second
)

// Function types for cleanup operations
type CleanupOrphanRelationshipsFunc func(driver neo4j.Driver) (int, error)
type CleanupOrphanNodesFunc func(driver neo4j.Driver) (int, error)
type GetAllConceptsFunc func(driver neo4j.Driver) ([]string, error)

// Exported variables for functions that can be mocked in tests
var CleanupOrphanRelationships CleanupOrphanRelationshipsFunc = cleanupOrphanRelationships
var CleanupOrphanNodes CleanupOrphanNodesFunc = cleanupOrphanNodes
var GetAllConcepts GetAllConceptsFunc = getAllConcepts

// SetupNeo4jConnection establishes a connection to the Neo4j database with retry logic to handle connection failures.
func SetupNeo4jConnection(cfg *config.Neo4jConfig) (neo4j.Driver, error) {
	log.Printf("Connecting to Neo4j at %s", cfg.URI)
	
	if cfg.URI == "" {
		return nil, apperrors.NewConfigError(apperrors.ErrInvalidInput, "Neo4j URI is not set")
	}
	
	if cfg.User == "" {
		return nil, apperrors.NewConfigError(apperrors.ErrInvalidInput, "Neo4j user is not set")
	}
	
	if cfg.Password == "" {
		return nil, apperrors.NewConfigError(apperrors.ErrInvalidInput, "Neo4j password is not set")
	}
	
	var driver neo4j.Driver
	var lastErr error
	
	for i := 0; i < cfg.MaxRetries; i++ {
		driver, lastErr = neo4j.NewDriver(cfg.URI, neo4j.BasicAuth(cfg.User, cfg.Password, ""))
		if lastErr != nil {
			log.Printf("Failed to create Neo4j driver (attempt %d/%d): %v", i+1, cfg.MaxRetries, lastErr)
			time.Sleep(cfg.RetryInterval)
			continue
		}
		
		lastErr = driver.VerifyConnectivity()
		if lastErr == nil {
			log.Printf("Successfully connected to Neo4j on attempt %d", i+1)
			return driver, nil
		}
		
		log.Printf("Failed to verify connectivity (attempt %d/%d): %v", i+1, cfg.MaxRetries, lastErr)
		time.Sleep(cfg.RetryInterval)
	}

	return nil, apperrors.NewDatabaseError(lastErr, fmt.Sprintf("failed to connect to Neo4j after %d attempts", cfg.MaxRetries))
}

// ConceptExists checks if a concept exists in the Neo4j database
func ConceptExists(driver neo4j.Driver, conceptName string) (bool, error) {
	if driver == nil {
		return false, apperrors.NewDatabaseError(apperrors.ErrInvalidInput, "Neo4j driver is nil")
	}

	if conceptName == "" {
		return false, apperrors.NewDatabaseError(apperrors.ErrInvalidInput, "concept name must not be empty")
	}

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `
			MATCH (c:Concept {name: $name})
			RETURN count(c) AS count
		`
		params := map[string]interface{}{
			"name": conceptName,
		}
		result, err := tx.Run(query, params)
		if err != nil {
			return false, apperrors.NewDatabaseError(err, "failed to execute Cypher query")
		}
		
		record, err := result.Single()
		if err != nil {
			return false, apperrors.NewDatabaseError(err, "failed to get query result")
		}
		
		count, _ := record.Get("count")
		return count.(int64) > 0, nil
	})

	if err != nil {
		return false, err
	}

	return result.(bool), nil
}

// CreateRelationship creates a relationship between two concepts in the Neo4j database using a Cypher query.
func CreateRelationship(driver neo4j.Driver, from, to, relation string) error {
	if driver == nil {
		return apperrors.NewDatabaseError(apperrors.ErrInvalidInput, "Neo4j driver is nil")
	}

	if from == "" || to == "" {
		return apperrors.NewDatabaseError(apperrors.ErrInvalidInput, "from and to must not be empty")
	}
	
	// Enhanced validation for relationship type
	if relation == "" {
		return apperrors.NewDatabaseError(apperrors.ErrInvalidInput, "relation must not be empty")
	}
	
	// Skip generic or vague relation types
	if relation == "related to" || relation == "is related to" || relation == "relates to" {
		return apperrors.NewDatabaseError(apperrors.ErrInvalidInput, "relation type is too generic")
	}

	// Check if both concepts exist or create them
	fromExists, err := ConceptExists(driver, from)
	if err != nil {
		return err
	}

	toExists, err := ConceptExists(driver, to)
	if err != nil {
		return err
	}

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		var query string
		
		if !fromExists && !toExists {
			// Create both nodes and the relationship
			query = `
				CREATE (a:Concept {name: $from})
				CREATE (b:Concept {name: $to})
				CREATE (a)-[r:RELATED_TO {type: $relation}]->(b)
			`
		} else if !fromExists {
			// Create only the 'from' node and the relationship
			query = `
				CREATE (a:Concept {name: $from})
				MATCH (b:Concept {name: $to})
				CREATE (a)-[r:RELATED_TO {type: $relation}]->(b)
			`
		} else if !toExists {
			// Create only the 'to' node and the relationship
			query = `
				MATCH (a:Concept {name: $from})
				CREATE (b:Concept {name: $to})
				CREATE (a)-[r:RELATED_TO {type: $relation}]->(b)
			`
		} else {
			// Both nodes exist, just create the relationship
			query = `
				MATCH (a:Concept {name: $from})
				MATCH (b:Concept {name: $to})
				CREATE (a)-[r:RELATED_TO {type: $relation}]->(b)
			`
		}
		
		params := map[string]interface{}{
			"from":     from,
			"to":       to,
			"relation": relation,
		}
		result, err := tx.Run(query, params)
		if err != nil {
			return nil, apperrors.NewDatabaseError(err, "failed to execute Cypher query")
		}
		
		// Consume the result to ensure the transaction is executed
		_, err = result.Consume()
		if err != nil {
			return nil, apperrors.NewDatabaseError(err, "failed to consume query result")
		}
		
		return nil, nil
	})

	if err != nil {
		return apperrors.NewDatabaseError(err, "failed to create relationship")
	}

	return nil
}

// QueryConcepts retrieves all concepts from the Neo4j database
func QueryConcepts(driver neo4j.Driver) ([]string, error) {
	if driver == nil {
		return nil, apperrors.NewDatabaseError(apperrors.ErrInvalidInput, "Neo4j driver is nil")
	}

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := "MATCH (c:Concept) RETURN c.name AS name"
		result, err := tx.Run(query, nil)
		if err != nil {
			return nil, apperrors.NewDatabaseError(err, "failed to execute Cypher query")
		}

		var concepts []string
		for result.Next() {
			record := result.Record()
			name, _ := record.Get("name")
			concepts = append(concepts, name.(string))
		}

		if err = result.Err(); err != nil {
			return nil, apperrors.NewDatabaseError(err, "error while iterating over results")
		}

		return concepts, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]string), nil
}

// QueryRelationships retrieves all relationships from the Neo4j database
func QueryRelationships(driver neo4j.Driver) ([]map[string]string, error) {
	if driver == nil {
		return nil, apperrors.NewDatabaseError(apperrors.ErrInvalidInput, "Neo4j driver is nil")
	}

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `
			MATCH (a:Concept)-[r:RELATED_TO]->(b:Concept)
			RETURN a.name AS from, r.type AS relation, b.name AS to
		`
		result, err := tx.Run(query, nil)
		if err != nil {
			return nil, apperrors.NewDatabaseError(err, "failed to execute Cypher query")
		}

		var relationships []map[string]string
		for result.Next() {
			record := result.Record()
			from, _ := record.Get("from")
			relation, _ := record.Get("relation")
			to, _ := record.Get("to")

			relationship := map[string]string{
				"from":     from.(string),
				"relation": relation.(string),
				"to":       to.(string),
			}
			relationships = append(relationships, relationship)
		}

		if err = result.Err(); err != nil {
			return nil, apperrors.NewDatabaseError(err, "error while iterating over results")
		}

		return relationships, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]map[string]string), nil
}

// GetLowConnectivityConcepts retrieves concepts with the least number of connections from the Neo4j database
func GetLowConnectivityConcepts(driver neo4j.Driver, limit int) ([]string, error) {
	if driver == nil {
		return nil, apperrors.NewDatabaseError(apperrors.ErrInvalidInput, "Neo4j driver is nil")
	}

	if limit <= 0 {
		limit = 10 // Default to 10 concepts if limit is not specified or invalid
	}

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		// This query counts both incoming and outgoing relationships for each concept
		// and orders them by the total count in ascending order (least connected first)
		query := `
			MATCH (c:Concept)
			OPTIONAL MATCH (c)-[r1:RELATED_TO]->()
			OPTIONAL MATCH ()-[r2:RELATED_TO]->(c)
			WITH c, count(r1) + count(r2) AS connectivity
			ORDER BY connectivity ASC
			LIMIT $limit
			RETURN c.name AS name, connectivity
		`
		params := map[string]interface{}{
			"limit": limit,
		}
		
		result, err := tx.Run(query, params)
		if err != nil {
			return nil, apperrors.NewDatabaseError(err, "failed to execute Cypher query")
		}

		var concepts []string
		for result.Next() {
			record := result.Record()
			name, _ := record.Get("name")
			concepts = append(concepts, name.(string))
		}

		if err = result.Err(); err != nil {
			return nil, apperrors.NewDatabaseError(err, "error while iterating over results")
		}

		return concepts, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]string), nil
}

// GetRandomLowConnectivityConcept retrieves a random concept from the low connectivity concepts
func GetRandomLowConnectivityConcept(driver neo4j.Driver, limit int) (string, error) {
	concepts, err := GetLowConnectivityConcepts(driver, limit)
	if err != nil {
		return "", err
	}

	if len(concepts) == 0 {
		return "", apperrors.NewDatabaseError(apperrors.ErrNotFound, "no concepts found")
	}

	// Select a random concept from the low connectivity concepts
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(concepts))
	return concepts[randomIndex], nil
}

// RelationshipExists checks if a relationship exists between two concepts
func RelationshipExists(driver neo4j.Driver, conceptA, conceptB string) (bool, error) {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	
	query := `
		MATCH (a:Concept {name: $conceptA})-[r]->(b:Concept {name: $conceptB})
		RETURN count(r) > 0 as exists
	`
	
	result, err := session.Run(query, map[string]interface{}{
		"conceptA": conceptA,
		"conceptB": conceptB,
	})
	
	if err != nil {
		return false, err
	}
	
	if result.Next() {
		return result.Record().GetByIndex(0).(bool), nil
	}
	
	return false, nil
}

// getAllConcepts returns all concepts in the graph
func getAllConcepts(driver neo4j.Driver) ([]string, error) {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	
	query := `
		MATCH (c:Concept)
		RETURN c.name as name
	`
	
	result, err := session.Run(query, nil)
	if err != nil {
		return nil, err
	}
	
	var concepts []string
	for result.Next() {
		name, _ := result.Record().Get("name")
		if nameStr, ok := name.(string); ok {
			concepts = append(concepts, nameStr)
		}
	}
	
	return concepts, nil
}

// cleanupOrphanRelationships removes relationships that connect to non-existent nodes
func cleanupOrphanRelationships(driver neo4j.Driver) (int, error) {
	if driver == nil {
		return 0, apperrors.NewDatabaseError(apperrors.ErrInvalidInput, "Neo4j driver is nil")
	}

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		// Find and delete relationships where either the start or end node doesn't exist
		query := `
			MATCH ()-[r:RELATED_TO]->()
			WHERE NOT EXISTS(r.type) OR r.type = ""
			WITH r, count(r) as count
			DELETE r
			RETURN count
		`
		result, err := tx.Run(query, nil)
		if err != nil {
			return 0, apperrors.NewDatabaseError(err, "failed to execute cleanup query")
		}
		
		record, err := result.Single()
		if err != nil {
			if err.Error() == "neo4j: no records returned" {
				return 0, nil
			}
			return 0, apperrors.NewDatabaseError(err, "failed to get query result")
		}
		
		count, _ := record.Get("count")
		// Handle different types of count (int or int64)
		switch c := count.(type) {
		case int64:
			return c, nil
		case int:
			return int64(c), nil
		case float64:
			return int64(c), nil
		default:
			// Return the count as is and let the caller handle the conversion
			return count, nil
		}
	})

	if err != nil {
		return 0, err
	}

	// Handle the result based on its type
	switch v := result.(type) {
	case int64:
		return int(v), nil
	case int:
		return v, nil
	case float64:
		return int(v), nil
	default:
		return 0, apperrors.NewDatabaseError(nil, "unexpected count type")
	}
}

// cleanupOrphanNodes removes nodes that don't have any relationships
func cleanupOrphanNodes(driver neo4j.Driver) (int, error) {
	if driver == nil {
		return 0, apperrors.NewDatabaseError(apperrors.ErrInvalidInput, "Neo4j driver is nil")
	}

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		// Find and delete nodes that don't have any relationships
		query := `
			MATCH (n:Concept)
			WHERE NOT (n)-[]-() 
			WITH n, count(n) as count
			DELETE n
			RETURN count
		`
		result, err := tx.Run(query, nil)
		if err != nil {
			return 0, apperrors.NewDatabaseError(err, "failed to execute cleanup query")
		}
		
		record, err := result.Single()
		if err != nil {
			if err.Error() == "neo4j: no records returned" {
				return 0, nil
			}
			return 0, apperrors.NewDatabaseError(err, "failed to get query result")
		}
		
		count, _ := record.Get("count")
		// Handle different types of count (int or int64)
		switch c := count.(type) {
		case int64:
			return c, nil
		case int:
			return int64(c), nil
		case float64:
			return int64(c), nil
		default:
			// Return the count as is and let the caller handle the conversion
			return count, nil
		}
	})

	if err != nil {
		return 0, err
	}

	// Handle the result based on its type
	switch v := result.(type) {
	case int64:
		return int(v), nil
	case int:
		return v, nil
	case float64:
		return int(v), nil
	default:
		return 0, apperrors.NewDatabaseError(nil, "unexpected count type")
	}
}
