package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kg-enricher/internal/config"
	"kg-enricher/internal/enricher"
	"kg-enricher/internal/llm"
	"kg-enricher/internal/neo4j"

	neo4jdriver "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// Version information
const (
	Version = "0.1.0"
)

func main() {
	// Parse command-line flags
	runOnce := flag.Bool("run-once", false, "Run once and exit")
	count := flag.Int("count", 10, "Number of relationships to mine when running once")
	showStats := flag.Bool("stats", false, "Show statistics and exit")
	showVersion := flag.Bool("version", false, "Show version information and exit")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Show version information if requested
	if *showVersion {
		fmt.Printf("Knowledge Graph Enricher v%s\n", Version)
		fmt.Printf("Default configuration:\n")
		fmt.Printf("  Batch size: %d\n", cfg.Enricher.BatchSize)
		fmt.Printf("  Interval: %s\n", cfg.Enricher.Interval)
		fmt.Printf("  Max relationships: %d\n", cfg.Enricher.MaxRelationships)
		fmt.Printf("  Concurrency: %d\n", cfg.Enricher.Concurrency)
		fmt.Printf("  LLM model: %s\n", cfg.LLM.Model)
		return
	}

	// Connect to Neo4j
	driver, err := neo4j.SetupNeo4jConnection(&cfg.Neo4j)
	if err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err)
	}
	defer driver.Close()

	// Initialize LLM service
	if err := llm.Initialize(&cfg.LLM); err != nil {
		log.Fatalf("Failed to initialize LLM service: %v", err)
	}

	// Create enricher
	e := enricher.NewEnricher(driver, &cfg.Enricher)

	// Show statistics if requested
	if *showStats {
		showStatistics(driver, e)
		return
	}

	// Run once if requested
	if *runOnce {
		log.Printf("Running once with count %d", *count)
		if err := e.RunOnce(*count); err != nil {
			log.Fatalf("Failed to run enricher: %v", err)
		}
		showStatistics(driver, e)
		return
	}

	// Start the enricher
	if err := e.Start(); err != nil {
		log.Fatalf("Failed to start enricher: %v", err)
	}

	// Handle signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signal
	<-sigChan
	log.Printf("Received signal, shutting down...")

	// Stop the enricher
	if err := e.Stop(); err != nil {
		log.Printf("Failed to stop enricher: %v", err)
	}

	// Show final statistics
	showStatistics(driver, e)
}

// showStatistics shows statistics about the knowledge graph and enricher
func showStatistics(driver neo4jdriver.Driver, e *enricher.Enricher) {
	// Get concepts
	concepts, err := neo4j.QueryAllConcepts(driver)
	if err != nil {
		log.Printf("Failed to query concepts: %v", err)
	} else {
		log.Printf("Total concepts in graph: %d", len(concepts))
	}

	// Get relationships
	relationships, err := neo4j.QueryRelationships(driver)
	if err != nil {
		log.Printf("Failed to query relationships: %v", err)
	} else {
		log.Printf("Total relationships in graph: %d", len(relationships))
	}

	// Show enricher statistics
	stats := e.GetStats()
	log.Printf("Enricher statistics:")
	log.Printf("  Total batches: %d", stats.TotalBatches)
	log.Printf("  Total pairs processed: %d", stats.TotalPairsProcessed)
	log.Printf("  Total relations found: %d", stats.TotalRelationsFound)
	log.Printf("  Total relations created: %d", stats.TotalRelationsCreated)
	
	if stats.LastBatchTime.IsZero() {
		log.Printf("  Last batch: never")
	} else {
		log.Printf("  Last batch: %s", stats.LastBatchTime.Format(time.RFC3339))
	}
	
	if stats.StartTime.IsZero() {
		log.Printf("  Running since: never")
	} else {
		log.Printf("  Running since: %s", stats.StartTime.Format(time.RFC3339))
		log.Printf("  Running for: %s", time.Since(stats.StartTime).Round(time.Second))
	}
} 