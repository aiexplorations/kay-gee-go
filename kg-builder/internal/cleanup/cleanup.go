package cleanup

import (
	"log"
	"time"

	"kg-builder/internal/neo4j"

	neo4jdriver "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// CleanupResult represents the result of a cleanup operation
type CleanupResult struct {
	OrphanRelationshipsRemoved int
	OrphanNodesRemoved         int
}

// CleanupOrphans removes orphan relationships and nodes from the database
func CleanupOrphans(driver neo4jdriver.Driver) (*CleanupResult, error) {
	result := &CleanupResult{}

	// Clean up orphan relationships
	log.Println("Cleaning up orphan relationships...")
	startTime := time.Now()
	relCount, err := neo4j.CleanupOrphanRelationships(driver)
	if err != nil {
		return nil, err
	}
	result.OrphanRelationshipsRemoved = relCount
	log.Printf("Removed %d orphan relationships in %v", relCount, time.Since(startTime))

	// Clean up orphan nodes
	log.Println("Cleaning up orphan nodes...")
	startTime = time.Now()
	nodeCount, err := neo4j.CleanupOrphanNodes(driver)
	if err != nil {
		return nil, err
	}
	result.OrphanNodesRemoved = nodeCount
	log.Printf("Removed %d orphan nodes in %v", nodeCount, time.Since(startTime))

	return result, nil
} 