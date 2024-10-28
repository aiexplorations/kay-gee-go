package main

import (
	"kg-builder/internal/graph"
	"kg-builder/internal/llm"
	"kg-builder/internal/neo4j"
	"log"
	"os"
	"time"
)

func main() {
	log.Println("Starting Knowledge Graph Builder") // Log the start of the application

	// Log all environment variables
	log.Println("Environment variables:")
	for _, env := range os.Environ() { // Iterate through environment variables
		log.Println(env) // Log each environment variable
	}

	neo4jDriver, err := neo4j.SetupNeo4jConnection() // Set up connection to Neo4j database
	if err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err) // Log fatal error if connection fails
	}
	defer neo4jDriver.Close() // Ensure the Neo4j driver is closed when main exits

	graphBuilder := graph.NewGraphBuilder(neo4jDriver, llm.GetRelatedConcepts, llm.MineRelationship) // Create a new graph builder

	seedConcept := "Artificial Intelligence" // Define the seed concept for graph building
	maxNodes := 100                          // Set the maximum number of nodes to build
	timeout := 30 * time.Minute              // Set the timeout for graph building

	log.Printf("Starting graph building with seed concept: %s", seedConcept) // Log the start of graph building
	err = graphBuilder.BuildGraph(seedConcept, maxNodes, timeout)            // Build the graph
	if err != nil {
		log.Printf("Graph building stopped: %v", err) // Log any errors during graph building
	}

	// Add a small delay to allow for graph building
	time.Sleep(5 * time.Second) // Sleep for 5 seconds

	log.Println("Starting random relationship mining") // Log the start of random relationship mining
	graphBuilder.MineRandomRelationships(50, 5)        // Mine 50 random relationships with 5 concurrent goroutines

	log.Println("Knowledge Graph Builder completed successfully") // Log successful completion of the application
}
