package main

import (
	"flag"
	"kg-builder/internal/config"
	"kg-builder/internal/graph"
	"kg-builder/internal/llm"
	"kg-builder/internal/neo4j"
	"log"
	"os"
	"time"

	neo4jdriver "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

const (
	// Version is the current version of the application
	Version = "0.1.0"
)

func main() {
	log.Printf("Starting Knowledge Graph Builder v%s", Version) // Log the start of the application with version

	// Parse command line arguments
	seedConcept := flag.String("seed", "", "Seed concept for graph building")
	maxNodes := flag.Int("max-nodes", 0, "Maximum number of nodes to build")
	timeoutMinutes := flag.Int("timeout", 0, "Timeout in minutes for graph building")
	randomRelationships := flag.Int("random-relationships", 0, "Number of random relationships to mine")
	concurrency := flag.Int("concurrency", 0, "Number of concurrent workers for mining random relationships")
	statsOnly := flag.Bool("stats-only", false, "Only show statistics without building the graph")
	showVersion := flag.Bool("version", false, "Show version information and exit")
	useLowConnectivity := flag.Bool("use-low-connectivity", false, "Use low connectivity concepts as seeds for subsequent graph building")
	flag.Parse()

	// Load configuration from environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Show version information if requested
	if *showVersion {
		log.Printf("Knowledge Graph Builder v%s", Version)
		log.Printf("Default configuration:")
		log.Printf("  Seed concept: Artificial Intelligence")
		log.Printf("  Max nodes: 100")
		log.Printf("  Timeout: 30 minutes")
		log.Printf("  Random relationships: 50")
		log.Printf("  Concurrency: 5")
		log.Printf("  LLM model: %s", cfg.LLM.Model)
		return
	}

	// Override configuration with command line arguments if provided
	if *seedConcept != "" {
		cfg.Graph.SeedConcept = *seedConcept
	}
	if *maxNodes > 0 {
		cfg.Graph.MaxNodes = *maxNodes
	}
	if *timeoutMinutes > 0 {
		cfg.Graph.Timeout = time.Duration(*timeoutMinutes) * time.Minute
	}
	if *randomRelationships > 0 {
		cfg.Graph.RandomRelationships = *randomRelationships
	}
	if *concurrency > 0 {
		cfg.Graph.Concurrency = *concurrency
	}

	// Log all environment variables
	log.Println("Environment variables:")
	for _, env := range os.Environ() { // Iterate through environment variables
		log.Println(env) // Log each environment variable
	}

	// Log configuration
	log.Printf("Configuration:")
	log.Printf("  Neo4j URI: %s", cfg.Neo4j.URI)
	log.Printf("  Neo4j User: %s", cfg.Neo4j.User)
	log.Printf("  LLM URL: %s", cfg.LLM.URL)
	log.Printf("  LLM Model: %s", cfg.LLM.Model)
	log.Printf("  LLM Cache Directory: %s", cfg.LLM.CacheDir)
	log.Printf("  Seed Concept: %s", cfg.Graph.SeedConcept)
	log.Printf("  Max Nodes: %d", cfg.Graph.MaxNodes)
	log.Printf("  Timeout: %v", cfg.Graph.Timeout)
	log.Printf("  Worker Count: %d", cfg.Graph.WorkerCount)
	log.Printf("  Random Relationships: %d", cfg.Graph.RandomRelationships)
	log.Printf("  Concurrency: %d", cfg.Graph.Concurrency)
	log.Printf("  Use Low Connectivity: %t", *useLowConnectivity)

	// Set up Neo4j connection
	neo4jDriver, err := neo4j.SetupNeo4jConnection(&cfg.Neo4j) // Set up connection to Neo4j database
	if err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err) // Log fatal error if connection fails
	}
	defer neo4jDriver.Close() // Ensure the Neo4j driver is closed when main exits

	// If stats-only flag is set, only show statistics and exit
	if *statsOnly {
		showStatistics(neo4jDriver)
		return
	}

	// Ensure cache directory exists
	if err := os.MkdirAll(cfg.LLM.CacheDir, 0755); err != nil {
		log.Printf("Warning: Failed to create cache directory: %v", err)
	}

	// Initialize LLM service
	if err := llm.Initialize(&cfg.LLM); err != nil {
		log.Fatalf("Failed to initialize LLM service: %v", err)
	}

	// Create graph builder
	graphBuilder := graph.NewGraphBuilder(neo4jDriver, llm.GetRelatedConcepts, llm.MineRelationship, &cfg.Graph) // Create a new graph builder

	// Build the graph
	log.Printf("Starting graph building with seed concept: %s", cfg.Graph.SeedConcept) // Log the start of graph building
	
	var buildErr error
	if *useLowConnectivity {
		// Build the graph using low connectivity concepts as seeds
		log.Printf("Using low connectivity concepts as seeds for subsequent graph building")
		buildErr = graphBuilder.BuildGraphWithLowConnectivitySeeds(cfg.Graph.SeedConcept, cfg.Graph.MaxNodes, cfg.Graph.Timeout)
	} else {
		// Build the graph using the traditional method
		buildErr = graphBuilder.BuildGraph(cfg.Graph.SeedConcept, cfg.Graph.MaxNodes, cfg.Graph.Timeout)
	}
	
	if buildErr != nil {
		log.Printf("Graph building stopped: %v", buildErr) // Log any errors during graph building
	} else {
		log.Printf("Graph building completed successfully")
	}

	// Add a small delay to allow for graph building
	time.Sleep(5 * time.Second) // Sleep for 5 seconds

	// Mine random relationships
	log.Println("Starting random relationship mining") // Log the start of random relationship mining
	err = graphBuilder.MineRandomRelationships(cfg.Graph.RandomRelationships, cfg.Graph.Concurrency) // Mine random relationships
	if err != nil {
		log.Printf("Random relationship mining stopped: %v", err) // Log any errors during random relationship mining
	} else {
		log.Printf("Random relationship mining completed successfully")
	}

	// Show statistics
	showStatistics(neo4jDriver)

	log.Println("Knowledge Graph Builder completed successfully") // Log successful completion of the application
}

// showStatistics shows statistics about the knowledge graph
func showStatistics(driver neo4jdriver.Driver) {
	// Query the graph to show statistics
	concepts, err := neo4j.QueryConcepts(driver)
	if err != nil {
		log.Printf("Failed to query concepts: %v", err)
	} else {
		log.Printf("Total concepts in graph: %d", len(concepts))
	}

	relationships, err := neo4j.QueryRelationships(driver)
	if err != nil {
		log.Printf("Failed to query relationships: %v", err)
	} else {
		log.Printf("Total relationships in graph: %d", len(relationships))
	}
}
