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

// CreateRelationship creates a relationship between two concepts in the Neo4j database using a Cypher query.
func CreateRelationship(driver neo4j.Driver, from, to, relation string) error {
	if driver == nil {
		return apperrors.NewDatabaseError(apperrors.ErrInvalidInput, "Neo4j driver is nil")
	}

	if from == "" || to == "" || relation == "" {
		return apperrors.NewDatabaseError(apperrors.ErrInvalidInput, "from, to, and relation must not be empty")
	}

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `
			MERGE (a:Concept {name: $from})
			MERGE (b:Concept {name: $to})
			MERGE (a)-[r:RELATED_TO {type: $relation}]->(b)
		`
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

// GetAllConcepts returns all concepts in the graph
func GetAllConcepts(driver neo4j.Driver) ([]string, error) {
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
