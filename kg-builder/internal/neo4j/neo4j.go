package neo4j

import (
	"fmt"
	"os"
	"time"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func SetupNeo4jConnection() (neo4j.Driver, error) {
	uri := "bolt://localhost:7687"
	username := "neo4j"
	password := "password"

	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to create driver: %w", err)
	}

	err = driver.VerifyConnectivity()
	if err != nil {
		return nil, fmt.Errorf("failed to verify connectivity: %w", err)
	}

	return driver, nil
}

func CreateRelationship(driver neo4j.Driver, from, to, relation string) error {
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := `
			MERGE (a:Concept {name: })
			MERGE (b:Concept {name: })
			MERGE (a)-[:]->(b)
		`
		params := map[string]interface{}{
			"from":     from,
			"to":       to,
			"relation": relation,
		}
		_, err := tx.Run(query, params)
		return nil, err
	})

	return err
}

func connectToNeo4jWithRetry(maxRetries int, retryInterval time.Duration) (neo4j.Driver, error) {
	var (
		driver neo4j.Driver
		err    error
	)

	neo4jURI := os.Getenv("NEO4J_URI")
	neo4jUser := os.Getenv("NEO4J_USER")
	neo4jPassword := os.Getenv("NEO4J_PASSWORD")

	for i := 0; i < maxRetries; i++ {
		driver, err = neo4j.NewDriver(neo4jURI, neo4j.BasicAuth(neo4jUser, neo4jPassword, ""))
		if err == nil {
			err = driver.VerifyConnectivity()
			if err == nil {
				fmt.Printf("Successfully connected to Neo4j on attempt %d\n", i+1)
				return driver, nil
			}
		}
		fmt.Printf("Failed to connect to Neo4j (attempt %d/%d): %v\n", i+1, maxRetries, err)
		time.Sleep(retryInterval)
	}

	return nil, fmt.Errorf("failed to connect to Neo4j after %d attempts", maxRetries)
}
