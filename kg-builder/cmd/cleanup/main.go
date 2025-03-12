package main

import (
	"flag"
	"log"
	"time"

	"kg-builder/internal/config"
	"kg-builder/internal/neo4j"
)

func main() {
	log.Println("Starting Knowledge Graph Cleanup")

	// Parse command line arguments
	configFile := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfigFromFile(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set up Neo4j connection
	driver, err := neo4j.SetupNeo4jConnection(&cfg.Neo4j)
	if err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err)
	}
	defer driver.Close()

	// Clean up orphan relationships
	log.Println("Cleaning up orphan relationships...")
	startTime := time.Now()
	count, err := neo4j.CleanupOrphanRelationships(driver)
	if err != nil {
		log.Fatalf("Error cleaning up orphan relationships: %v", err)
	}
	log.Printf("Removed %d orphan relationships in %v", count, time.Since(startTime))

	// Clean up orphan nodes
	log.Println("Cleaning up orphan nodes...")
	startTime = time.Now()
	count, err = neo4j.CleanupOrphanNodes(driver)
	if err != nil {
		log.Fatalf("Error cleaning up orphan nodes: %v", err)
	}
	log.Printf("Removed %d orphan nodes in %v", count, time.Since(startTime))

	log.Println("Knowledge Graph Cleanup completed successfully")
} 