package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfigDefaults(t *testing.T) {
	// Save current environment variables
	oldNeo4jURI := os.Getenv("NEO4J_URI")
	oldNeo4jUser := os.Getenv("NEO4J_USER")
	oldNeo4jPassword := os.Getenv("NEO4J_PASSWORD")
	oldLLMURL := os.Getenv("LLM_URL")
	oldLLMModel := os.Getenv("LLM_MODEL")
	oldConfigFile := os.Getenv("CONFIG_FILE")

	// Clean up environment variables
	os.Unsetenv("NEO4J_URI")
	os.Unsetenv("NEO4J_USER")
	os.Unsetenv("NEO4J_PASSWORD")
	os.Unsetenv("LLM_URL")
	os.Unsetenv("LLM_MODEL")
	os.Unsetenv("CONFIG_FILE")

	// Restore environment variables after test
	defer func() {
		os.Setenv("NEO4J_URI", oldNeo4jURI)
		os.Setenv("NEO4J_USER", oldNeo4jUser)
		os.Setenv("NEO4J_PASSWORD", oldNeo4jPassword)
		os.Setenv("LLM_URL", oldLLMURL)
		os.Setenv("LLM_MODEL", oldLLMModel)
		os.Setenv("CONFIG_FILE", oldConfigFile)
	}()

	// Set a non-existent config file to force using defaults
	os.Setenv("CONFIG_FILE", "non_existent_config.yaml")

	// Load config
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check default values
	if config.Neo4j.URI != "bolt://neo4j:7687" {
		t.Errorf("Expected Neo4j URI to be 'bolt://neo4j:7687', got '%s'", config.Neo4j.URI)
	}
	if config.Neo4j.User != "neo4j" {
		t.Errorf("Expected Neo4j User to be 'neo4j', got '%s'", config.Neo4j.User)
	}
	if config.Neo4j.Password != "password" {
		t.Errorf("Expected Neo4j Password to be 'password', got '%s'", config.Neo4j.Password)
	}
	if config.Neo4j.MaxRetries != 5 {
		t.Errorf("Expected Neo4j MaxRetries to be 5, got %d", config.Neo4j.MaxRetries)
	}
	if config.Neo4j.RetryInterval != 5*time.Second {
		t.Errorf("Expected Neo4j RetryInterval to be 5s, got %s", config.Neo4j.RetryInterval)
	}

	if config.LLM.URL != "http://host.docker.internal:11434/api/generate" {
		t.Errorf("Expected LLM URL to be 'http://host.docker.internal:11434/api/generate', got '%s'", config.LLM.URL)
	}
	if config.LLM.Model != "qwen2.5:3b" {
		t.Errorf("Expected LLM Model to be 'qwen2.5:3b', got '%s'", config.LLM.Model)
	}
	if config.LLM.CacheDir != "./cache/llm" {
		t.Errorf("Expected LLM CacheDir to be './cache/llm', got '%s'", config.LLM.CacheDir)
	}

	if config.Graph.SeedConcept != "Artificial Intelligence" {
		t.Errorf("Expected Graph SeedConcept to be 'Artificial Intelligence', got '%s'", config.Graph.SeedConcept)
	}
	if config.Graph.MaxNodes != 100 {
		t.Errorf("Expected Graph MaxNodes to be 100, got %d", config.Graph.MaxNodes)
	}
	if config.Graph.Timeout != 30*time.Minute {
		t.Errorf("Expected Graph Timeout to be 30m, got %s", config.Graph.Timeout)
	}
	if config.Graph.WorkerCount != 10 {
		t.Errorf("Expected Graph WorkerCount to be 10, got %d", config.Graph.WorkerCount)
	}
	if config.Graph.RandomRelationships != 50 {
		t.Errorf("Expected Graph RandomRelationships to be 50, got %d", config.Graph.RandomRelationships)
	}
	if config.Graph.Concurrency != 5 {
		t.Errorf("Expected Graph Concurrency to be 5, got %d", config.Graph.Concurrency)
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	// Save current environment variables
	oldNeo4jURI := os.Getenv("NEO4J_URI")
	oldNeo4jUser := os.Getenv("NEO4J_USER")
	oldNeo4jPassword := os.Getenv("NEO4J_PASSWORD")
	oldLLMURL := os.Getenv("LLM_URL")
	oldLLMModel := os.Getenv("LLM_MODEL")
	oldConfigFile := os.Getenv("CONFIG_FILE")

	// Set environment variables
	os.Setenv("NEO4J_URI", "bolt://test-neo4j:7687")
	os.Setenv("NEO4J_USER", "test-user")
	os.Setenv("NEO4J_PASSWORD", "test-password")
	os.Setenv("LLM_URL", "http://test-llm:11434/api/generate")
	os.Setenv("LLM_MODEL", "test-model")
	os.Setenv("CONFIG_FILE", "non_existent_config.yaml")

	// Restore environment variables after test
	defer func() {
		os.Setenv("NEO4J_URI", oldNeo4jURI)
		os.Setenv("NEO4J_USER", oldNeo4jUser)
		os.Setenv("NEO4J_PASSWORD", oldNeo4jPassword)
		os.Setenv("LLM_URL", oldLLMURL)
		os.Setenv("LLM_MODEL", oldLLMModel)
		os.Setenv("CONFIG_FILE", oldConfigFile)
	}()

	// Load config
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check values from environment variables
	if config.Neo4j.URI != "bolt://test-neo4j:7687" {
		t.Errorf("Expected Neo4j URI to be 'bolt://test-neo4j:7687', got '%s'", config.Neo4j.URI)
	}
	if config.Neo4j.User != "test-user" {
		t.Errorf("Expected Neo4j User to be 'test-user', got '%s'", config.Neo4j.User)
	}
	if config.Neo4j.Password != "test-password" {
		t.Errorf("Expected Neo4j Password to be 'test-password', got '%s'", config.Neo4j.Password)
	}
	if config.LLM.URL != "http://test-llm:11434/api/generate" {
		t.Errorf("Expected LLM URL to be 'http://test-llm:11434/api/generate', got '%s'", config.LLM.URL)
	}
	if config.LLM.Model != "test-model" {
		t.Errorf("Expected LLM Model to be 'test-model', got '%s'", config.LLM.Model)
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	// Create a temporary config file
	configContent := `
neo4j:
  uri: "bolt://file-neo4j:7687"
  user: "file-user"
  password: "file-password"
  max_retries: 10
  retry_interval_seconds: 10

llm:
  url: "http://file-llm:11434/api/generate"
  model: "file-model"
  cache_dir: "./file-cache/llm"

graph:
  seed_concept: "File Concept"
  max_nodes: 200
  timeout_minutes: 60
  worker_count: 20
  random_relationships: 100
  concurrency: 10
`
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(configContent)); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	// Save current environment variables
	oldNeo4jURI := os.Getenv("NEO4J_URI")
	oldNeo4jUser := os.Getenv("NEO4J_USER")
	oldNeo4jPassword := os.Getenv("NEO4J_PASSWORD")
	oldLLMURL := os.Getenv("LLM_URL")
	oldLLMModel := os.Getenv("LLM_MODEL")
	oldLLMCacheDir := os.Getenv("LLM_CACHE_DIR")
	oldConfigFile := os.Getenv("CONFIG_FILE")

	// Unset environment variables to ensure they don't interfere
	os.Unsetenv("NEO4J_URI")
	os.Unsetenv("NEO4J_USER")
	os.Unsetenv("NEO4J_PASSWORD")
	os.Unsetenv("LLM_URL")
	os.Unsetenv("LLM_MODEL")
	os.Unsetenv("LLM_CACHE_DIR")
	
	// Set config file environment variable
	os.Setenv("CONFIG_FILE", tmpfile.Name())

	// Restore environment variables after test
	defer func() {
		os.Setenv("NEO4J_URI", oldNeo4jURI)
		os.Setenv("NEO4J_USER", oldNeo4jUser)
		os.Setenv("NEO4J_PASSWORD", oldNeo4jPassword)
		os.Setenv("LLM_URL", oldLLMURL)
		os.Setenv("LLM_MODEL", oldLLMModel)
		os.Setenv("LLM_CACHE_DIR", oldLLMCacheDir)
		os.Setenv("CONFIG_FILE", oldConfigFile)
	}()

	// Load config
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check values from file
	if config.Neo4j.URI != "bolt://file-neo4j:7687" {
		t.Errorf("Expected Neo4j URI to be 'bolt://file-neo4j:7687', got '%s'", config.Neo4j.URI)
	}
	if config.Neo4j.User != "file-user" {
		t.Errorf("Expected Neo4j User to be 'file-user', got '%s'", config.Neo4j.User)
	}
	if config.Neo4j.Password != "file-password" {
		t.Errorf("Expected Neo4j Password to be 'file-password', got '%s'", config.Neo4j.Password)
	}
	if config.Neo4j.MaxRetries != 10 {
		t.Errorf("Expected Neo4j MaxRetries to be 10, got %d", config.Neo4j.MaxRetries)
	}
	if config.Neo4j.RetryInterval != 10*time.Second {
		t.Errorf("Expected Neo4j RetryInterval to be 10s, got %s", config.Neo4j.RetryInterval)
	}

	if config.LLM.URL != "http://file-llm:11434/api/generate" {
		t.Errorf("Expected LLM URL to be 'http://file-llm:11434/api/generate', got '%s'", config.LLM.URL)
	}
	if config.LLM.Model != "file-model" {
		t.Errorf("Expected LLM Model to be 'file-model', got '%s'", config.LLM.Model)
	}
	if config.LLM.CacheDir != "./file-cache/llm" {
		t.Errorf("Expected LLM CacheDir to be './file-cache/llm', got '%s'", config.LLM.CacheDir)
	}

	if config.Graph.SeedConcept != "File Concept" {
		t.Errorf("Expected Graph SeedConcept to be 'File Concept', got '%s'", config.Graph.SeedConcept)
	}
	if config.Graph.MaxNodes != 200 {
		t.Errorf("Expected Graph MaxNodes to be 200, got %d", config.Graph.MaxNodes)
	}
	if config.Graph.Timeout != 60*time.Minute {
		t.Errorf("Expected Graph Timeout to be 60m, got %s", config.Graph.Timeout)
	}
	if config.Graph.WorkerCount != 20 {
		t.Errorf("Expected Graph WorkerCount to be 20, got %d", config.Graph.WorkerCount)
	}
	if config.Graph.RandomRelationships != 100 {
		t.Errorf("Expected Graph RandomRelationships to be 100, got %d", config.Graph.RandomRelationships)
	}
	if config.Graph.Concurrency != 10 {
		t.Errorf("Expected Graph Concurrency to be 10, got %d", config.Graph.Concurrency)
	}
}

func TestEnvOverridesFile(t *testing.T) {
	// Create a temporary config file
	configContent := `
neo4j:
  uri: "bolt://file-neo4j:7687"
  user: "file-user"
  password: "file-password"

llm:
  url: "http://file-llm:11434/api/generate"
  model: "file-model"
`
	tmpfile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(configContent)); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	// Save current environment variables
	oldNeo4jURI := os.Getenv("NEO4J_URI")
	oldNeo4jUser := os.Getenv("NEO4J_USER")
	oldNeo4jPassword := os.Getenv("NEO4J_PASSWORD")
	oldLLMURL := os.Getenv("LLM_URL")
	oldLLMModel := os.Getenv("LLM_MODEL")
	oldConfigFile := os.Getenv("CONFIG_FILE")

	// Set environment variables
	os.Setenv("NEO4J_URI", "bolt://env-neo4j:7687")
	os.Setenv("NEO4J_USER", "env-user")
	os.Setenv("NEO4J_PASSWORD", "env-password")
	os.Setenv("LLM_URL", "http://env-llm:11434/api/generate")
	os.Setenv("LLM_MODEL", "env-model")
	os.Setenv("CONFIG_FILE", tmpfile.Name())

	// Restore environment variables after test
	defer func() {
		os.Setenv("NEO4J_URI", oldNeo4jURI)
		os.Setenv("NEO4J_USER", oldNeo4jUser)
		os.Setenv("NEO4J_PASSWORD", oldNeo4jPassword)
		os.Setenv("LLM_URL", oldLLMURL)
		os.Setenv("LLM_MODEL", oldLLMModel)
		os.Setenv("CONFIG_FILE", oldConfigFile)
	}()

	// Load config
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Check that environment variables override file values
	if config.Neo4j.URI != "bolt://env-neo4j:7687" {
		t.Errorf("Expected Neo4j URI to be 'bolt://env-neo4j:7687', got '%s'", config.Neo4j.URI)
	}
	if config.Neo4j.User != "env-user" {
		t.Errorf("Expected Neo4j User to be 'env-user', got '%s'", config.Neo4j.User)
	}
	if config.Neo4j.Password != "env-password" {
		t.Errorf("Expected Neo4j Password to be 'env-password', got '%s'", config.Neo4j.Password)
	}
	if config.LLM.URL != "http://env-llm:11434/api/generate" {
		t.Errorf("Expected LLM URL to be 'http://env-llm:11434/api/generate', got '%s'", config.LLM.URL)
	}
	if config.LLM.Model != "env-model" {
		t.Errorf("Expected LLM Model to be 'env-model', got '%s'", config.LLM.Model)
	}
} 