package builder

import (
	"context"
	"testing"
	"time"

	"github.com/kay-gee-go/internal/common/config"
	"github.com/stretchr/testify/assert"
)

func TestNewBuilder(t *testing.T) {
	// Create a test configuration
	cfg := &config.BuilderConfig{
		Neo4j: config.Neo4jConfig{
			URI:               "bolt://localhost:7687",
			User:              "neo4j",
			Password:          "password",
			MaxRetries:        3,
			RetryIntervalSecs: 1,
		},
		LLM: config.LLMConfig{
			URL:      "http://localhost:11434/api/generate",
			Model:    "qwen2.5:3b",
			CacheDir: "./test-cache",
		},
		Graph: config.GraphConfig{
			SeedConcept:        "Test",
			MaxNodes:           10,
			TimeoutMinutes:     1,
			WorkerCount:        2,
			RandomRelationships: 5,
			Concurrency:        2,
		},
	}

	// Skip the test if we're not running in an environment with Neo4j
	t.Skip("Skipping test that requires Neo4j")

	// Create a new builder
	b, err := NewBuilder(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, b)
	defer b.Close()

	// Check that the builder was initialized correctly
	assert.Equal(t, cfg, b.config)
	assert.NotNil(t, b.neo4jClient)
	assert.NotNil(t, b.llmClient)
	assert.NotNil(t, b.stats)
	assert.Equal(t, cfg.Graph.SeedConcept, b.stats.SeedConcept)
	assert.Equal(t, 0, b.stats.NodesCreated)
	assert.Equal(t, 0, b.stats.RelationshipsCreated)
	assert.NotZero(t, b.stats.StartTime)
	assert.Zero(t, b.stats.EndTime)
	assert.Empty(t, b.stats.Duration)
	assert.NotNil(t, b.concepts)
	assert.Empty(t, b.concepts)
}

func TestBuild(t *testing.T) {
	// Create a test configuration
	cfg := &config.BuilderConfig{
		Neo4j: config.Neo4jConfig{
			URI:               "bolt://localhost:7687",
			User:              "neo4j",
			Password:          "password",
			MaxRetries:        3,
			RetryIntervalSecs: 1,
		},
		LLM: config.LLMConfig{
			URL:      "http://localhost:11434/api/generate",
			Model:    "qwen2.5:3b",
			CacheDir: "./test-cache",
		},
		Graph: config.GraphConfig{
			SeedConcept:        "Test",
			MaxNodes:           2,
			TimeoutMinutes:     1,
			WorkerCount:        1,
			RandomRelationships: 0,
			Concurrency:        1,
		},
	}

	// Skip the test if we're not running in an environment with Neo4j and LLM
	t.Skip("Skipping test that requires Neo4j and LLM")

	// Create a new builder
	b, err := NewBuilder(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, b)
	defer b.Close()

	// Build the knowledge graph
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stats, err := b.Build(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, stats)

	// Check that the stats were updated correctly
	assert.Equal(t, cfg.Graph.SeedConcept, stats.SeedConcept)
	assert.GreaterOrEqual(t, stats.NodesCreated, 1)
	assert.GreaterOrEqual(t, stats.RelationshipsCreated, 0)
	assert.NotZero(t, stats.StartTime)
	assert.NotZero(t, stats.EndTime)
	assert.NotEmpty(t, stats.Duration)
} 