package main

import (
    "log"
    "time"
    "kg-builder/internal/graph"
    "kg-builder/internal/llm"
    "kg-builder/internal/neo4j"
)

func main() {
    log.Println("Starting Knowledge Graph Builder")

    neo4jDriver, err := neo4j.SetupNeo4jConnection()
    if err != nil {
        log.Fatalf("Failed to connect to Neo4j: %v", err)
    }
    defer neo4jDriver.Close()

    graphBuilder := graph.NewGraphBuilder(neo4jDriver, llm.GetRelatedConcepts)

    seedConcept := "Artificial Intelligence"
    maxNodes := 100
    timeout := 20 * time.Minute

    err = graphBuilder.BuildGraph(seedConcept, maxNodes, timeout)
    if err != nil {
        log.Fatalf("Failed to build graph: %v", err)
    }

    log.Println("Knowledge Graph Builder completed")
}