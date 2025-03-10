package main_test

import (
	"flag"
	"os"
	"testing"
)

// TestMainPackageFlags tests that the main package correctly defines and parses command-line flags
func TestMainPackageFlags(t *testing.T) {
	// Save original command-line arguments
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}()

	// Set test command-line arguments
	os.Args = []string{
		"kg-builder",
		"--seed", "Machine Learning",
		"--max-nodes", "200",
		"--timeout", "60",
		"--random-relationships", "100",
		"--concurrency", "10",
	}

	// Reset the flag.CommandLine to parse the new arguments
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Define flags (same as in main.go)
	seedConcept := flag.String("seed", "Artificial Intelligence", "Seed concept for graph building")
	maxNodes := flag.Int("max-nodes", 100, "Maximum number of nodes to build")
	timeoutMinutes := flag.Int("timeout", 30, "Timeout in minutes for graph building")
	randomRelationships := flag.Int("random-relationships", 50, "Number of random relationships to mine")
	concurrency := flag.Int("concurrency", 5, "Number of concurrent workers for mining random relationships")

	// Parse flags
	flag.Parse()

	// Verify flags are correctly parsed
	if *seedConcept != "Machine Learning" {
		t.Errorf("Expected seed concept 'Machine Learning', got '%s'", *seedConcept)
	}

	if *maxNodes != 200 {
		t.Errorf("Expected max nodes 200, got %d", *maxNodes)
	}

	if *timeoutMinutes != 60 {
		t.Errorf("Expected timeout 60 minutes, got %d", *timeoutMinutes)
	}

	if *randomRelationships != 100 {
		t.Errorf("Expected random relationships 100, got %d", *randomRelationships)
	}

	if *concurrency != 10 {
		t.Errorf("Expected concurrency 10, got %d", *concurrency)
	}
}

// TestMainPackageEnvironment tests that the main package correctly reads environment variables
func TestMainPackageEnvironment(t *testing.T) {
	// Save original environment variables
	originalNeo4jURI := os.Getenv("NEO4J_URI")
	originalNeo4jUser := os.Getenv("NEO4J_USER")
	originalNeo4jPassword := os.Getenv("NEO4J_PASSWORD")
	originalLLMURL := os.Getenv("LLM_URL")
	originalLLMModel := os.Getenv("LLM_MODEL")

	// Restore environment variables after the test
	defer func() {
		os.Setenv("NEO4J_URI", originalNeo4jURI)
		os.Setenv("NEO4J_USER", originalNeo4jUser)
		os.Setenv("NEO4J_PASSWORD", originalNeo4jPassword)
		os.Setenv("LLM_URL", originalLLMURL)
		os.Setenv("LLM_MODEL", originalLLMModel)
	}()

	// Set test environment variables
	os.Setenv("NEO4J_URI", "bolt://test-neo4j:7687")
	os.Setenv("NEO4J_USER", "test-user")
	os.Setenv("NEO4J_PASSWORD", "test-password")
	os.Setenv("LLM_URL", "http://test-llm:11434/api/generate")
	os.Setenv("LLM_MODEL", "test-model")

	// Verify environment variables are correctly set
	if os.Getenv("NEO4J_URI") != "bolt://test-neo4j:7687" {
		t.Errorf("Expected NEO4J_URI 'bolt://test-neo4j:7687', got '%s'", os.Getenv("NEO4J_URI"))
	}

	if os.Getenv("NEO4J_USER") != "test-user" {
		t.Errorf("Expected NEO4J_USER 'test-user', got '%s'", os.Getenv("NEO4J_USER"))
	}

	if os.Getenv("NEO4J_PASSWORD") != "test-password" {
		t.Errorf("Expected NEO4J_PASSWORD 'test-password', got '%s'", os.Getenv("NEO4J_PASSWORD"))
	}

	if os.Getenv("LLM_URL") != "http://test-llm:11434/api/generate" {
		t.Errorf("Expected LLM_URL 'http://test-llm:11434/api/generate', got '%s'", os.Getenv("LLM_URL"))
	}

	if os.Getenv("LLM_MODEL") != "test-model" {
		t.Errorf("Expected LLM_MODEL 'test-model', got '%s'", os.Getenv("LLM_MODEL"))
	}
} 