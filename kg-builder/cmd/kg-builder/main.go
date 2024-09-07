package main

import (
    "log"
    "time"
    "kg-builder/internal/graph"
    "kg-builder/internal/llm"
    "kg-builder/internal/neo4j"
    "os"
)

func main() {
    log.Println("Starting Knowledge Graph Builder")

    // Log all environment variables
    log.Println("Environment variables:")
    for _, env := range os.Environ() {
        log.Println(env)
    }

    neo4jDriver, err := neo4j.SetupNeo4jConnection()
    if err != nil {
        log.Fatalf("Failed to connect to Neo4j: %v", err)
    }
    defer neo4jDriver.Close()

    graphBuilder := graph.NewGraphBuilder(neo4jDriver, llm.GetRelatedConcepts, llm.MineRelationship)

    seedConcept := "Artificial Intelligence"
    maxNodes := 100
    timeout := 30 * time.Minute

    log.Printf("Starting graph building with seed concept: %s", seedConcept)
    err = graphBuilder.BuildGraph(seedConcept, maxNodes, timeout)
    if err != nil {
        log.Printf("Graph building stopped: %v", err)
    }

    // Add a small delay to allow for graph building
    time.Sleep(5 * time.Second)

    log.Println("Starting random relationship mining")
    graphBuilder.MineRandomRelationships(50, 5) // Mine 50 random relationships with 5 concurrent goroutines

    log.Println("Knowledge Graph Builder completed successfully")
}