package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Test loading config from environment variables
	os.Setenv("NEO4J_URI", "bolt://test:7687")
	os.Setenv("NEO4J_USER", "testuser")
	os.Setenv("NEO4J_PASSWORD", "testpass")
	os.Setenv("LLM_URL", "http://test:11434/api/generate")
	os.Setenv("LLM_MODEL", "test-model")
	os.Setenv("ENRICHER_BATCH_SIZE", "15")
	os.Setenv("ENRICHER_INTERVAL_SECONDS", "30")
	os.Setenv("ENRICHER_MAX_RELATIONSHIPS", "200")
	os.Setenv("ENRICHER_CONCURRENCY", "3")
	
	// Load the config
	cfg, err := LoadConfig()
	
	// Assert that there was no error
	assert.NoError(t, err)
	
	// Assert that the config was loaded correctly
	assert.Equal(t, "bolt://test:7687", cfg.Neo4j.URI)
	assert.Equal(t, "testuser", cfg.Neo4j.Username)
	assert.Equal(t, "testpass", cfg.Neo4j.Password)
	assert.Equal(t, "http://test:11434/api/generate", cfg.LLM.URL)
	assert.Equal(t, "test-model", cfg.LLM.Model)
	assert.Equal(t, 15, cfg.Enricher.BatchSize)
	assert.Equal(t, time.Second*30, cfg.Enricher.Interval)
	assert.Equal(t, 200, cfg.Enricher.MaxRelationships)
	assert.Equal(t, 3, cfg.Enricher.Concurrency)
	
	// Reset environment variables
	os.Unsetenv("NEO4J_URI")
	os.Unsetenv("NEO4J_USER")
	os.Unsetenv("NEO4J_PASSWORD")
	os.Unsetenv("LLM_URL")
	os.Unsetenv("LLM_MODEL")
	os.Unsetenv("ENRICHER_BATCH_SIZE")
	os.Unsetenv("ENRICHER_INTERVAL_SECONDS")
	os.Unsetenv("ENRICHER_MAX_RELATIONSHIPS")
	os.Unsetenv("ENRICHER_CONCURRENCY")
}

func TestLoadConfigDefaults(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("NEO4J_URI")
	os.Unsetenv("NEO4J_USER")
	os.Unsetenv("NEO4J_PASSWORD")
	os.Unsetenv("LLM_URL")
	os.Unsetenv("LLM_MODEL")
	os.Unsetenv("ENRICHER_BATCH_SIZE")
	os.Unsetenv("ENRICHER_INTERVAL_SECONDS")
	os.Unsetenv("ENRICHER_MAX_RELATIONSHIPS")
	os.Unsetenv("ENRICHER_CONCURRENCY")
	
	// Load the config
	cfg, err := LoadConfig()
	
	// Assert that there was no error
	assert.NoError(t, err)
	
	// Assert that the config was loaded with default values
	assert.Equal(t, "bolt://neo4j:7687", cfg.Neo4j.URI)
	assert.Equal(t, "neo4j", cfg.Neo4j.Username)
	assert.Equal(t, "password", cfg.Neo4j.Password)
	assert.Equal(t, "http://host.docker.internal:11434/api/generate", cfg.LLM.URL)
	assert.Equal(t, "qwen2.5:3b", cfg.LLM.Model)
	assert.Equal(t, 10, cfg.Enricher.BatchSize)
	assert.Equal(t, time.Second*60, cfg.Enricher.Interval)
	assert.Equal(t, 100, cfg.Enricher.MaxRelationships)
	assert.Equal(t, 5, cfg.Enricher.Concurrency)
} 