package neo4j

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// SetupNeo4jConnection establishes a connection to the Neo4j database with retry logic to handle connection failures.
func SetupNeo4jConnection() (neo4j.Driver, error) {
	return connectToNeo4jWithRetry(5, 5*time.Second)
}

// CreateRelationship creates a relationship between two concepts in the Neo4j database using a Cypher query.
func CreateRelationship(driver neo4j.Driver, from, to, relation string) error {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	// Write a transaction to create the relationship
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
		_, err := tx.Run(query, params)
		return nil, err // Return the error from the transaction
	})
	// Return the error from the transaction
	return err
}

// connectToNeo4jWithRetry attempts to connect to the Neo4j database multiple times with retry logic.

func connectToNeo4jWithRetry(maxRetries int, retryInterval time.Duration) (neo4j.Driver, error) {
	neo4jURI := os.Getenv("NEO4J_URI")
	if neo4jURI == "" {
		return nil, fmt.Errorf("NEO4J_URI environment variable is not set")
	}

	// Parse the URI to ensure it's valid
	_, err := url.Parse(neo4jURI)
	if err != nil {
		return nil, fmt.Errorf("invalid NEO4J_URI: %v", err)
	}

	// Get the Neo4j user and password from the environment variables
	neo4jUser := os.Getenv("NEO4J_USER")
	if neo4jUser == "" {
		return nil, fmt.Errorf("NEO4J_USER environment variable is not set")
	}

	neo4jPassword := os.Getenv("NEO4J_PASSWORD")
	if neo4jPassword == "" {
		return nil, fmt.Errorf("NEO4J_PASSWORD environment variable is not set")
	}

	log.Printf("Attempting to connect to Neo4j at %s", neo4jURI)

	// Attempt to create a driver with retry logic
	var driver neo4j.Driver
	for i := 0; i < maxRetries; i++ {
		driver, err = neo4j.NewDriver(neo4jURI, neo4j.BasicAuth(neo4jUser, neo4jPassword, ""))
		if err == nil {
			log.Printf("Driver created successfully, verifying connectivity...")
			err = driver.VerifyConnectivity()
			if err == nil {
				log.Printf("Successfully connected to Neo4j on attempt %d", i+1)
				return driver, nil
			}
		}
		// Log the failure and wait before retrying
		log.Printf("Failed to connect to Neo4j (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(retryInterval)
	}
	// If all attempts fail, return an error
	return nil, fmt.Errorf("failed to connect to Neo4j after %d attempts", maxRetries)
}
