package enricher

import (
	"context"
	"testing"
	"time"

	"github.com/kay-gee-go/internal/common/config"
	"github.com/stretchr/testify/assert"
)

func TestNewEnricher(t *testing.T) {
	// Create a test configuration
	cfg := &config.EnricherAppConfig{
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
		Enricher: config.EnricherConfig{
			BatchSize:        5,
			IntervalSeconds:  10,
			MaxRelationships: 10,
			Concurrency:      2,
		},
	}

	// Skip the test if we're not running in an environment with Neo4j
	t.Skip("Skipping test that requires Neo4j")

	// Create a new enricher
	e, err := NewEnricher(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, e)
	defer e.Close()

	// Check that the enricher was initialized correctly
	assert.Equal(t, cfg, e.config)
	assert.NotNil(t, e.neo4jClient)
	assert.NotNil(t, e.llmClient)
	assert.NotNil(t, e.stats)
	assert.Equal(t, 0, e.stats.PairsProcessed)
	assert.Equal(t, 0, e.stats.RelationshipsCreated)
	assert.NotZero(t, e.stats.StartTime)
	assert.Zero(t, e.stats.EndTime)
	assert.Empty(t, e.stats.Duration)
}

func TestEnrichOnce(t *testing.T) {
	// Create a test configuration
	cfg := &config.EnricherAppConfig{
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
		Enricher: config.EnricherConfig{
			BatchSize:        2,
			IntervalSeconds:  10,
			MaxRelationships: 10,
			Concurrency:      1,
		},
	}

	// Skip the test if we're not running in an environment with Neo4j and LLM
	t.Skip("Skipping test that requires Neo4j and LLM")

	// Create a new enricher
	e, err := NewEnricher(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, e)
	defer e.Close()

	// Enrich the knowledge graph once
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stats, err := e.EnrichOnce(ctx, 2)
	assert.NoError(t, err)
	assert.NotNil(t, stats)

	// Check that the stats were updated correctly
	assert.GreaterOrEqual(t, stats.PairsProcessed, 0)
	assert.GreaterOrEqual(t, stats.RelationshipsCreated, 0)
	assert.NotZero(t, stats.StartTime)
	assert.NotZero(t, stats.EndTime)
	assert.NotEmpty(t, stats.Duration)
}

func TestEnrichContinuously(t *testing.T) {
	// Create a test configuration
	cfg := &config.EnricherAppConfig{
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
		Enricher: config.EnricherConfig{
			BatchSize:        2,
			IntervalSeconds:  1,
			MaxRelationships: 1,
			Concurrency:      1,
		},
	}

	// Skip the test if we're not running in an environment with Neo4j and LLM
	t.Skip("Skipping test that requires Neo4j and LLM")

	// Create a new enricher
	e, err := NewEnricher(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, e)
	defer e.Close()

	// Enrich the knowledge graph continuously
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = e.EnrichContinuously(ctx)
	assert.Equal(t, context.DeadlineExceeded, err)

	// Check that the stats were updated correctly
	stats := e.GetStats()
	assert.GreaterOrEqual(t, stats.PairsProcessed, 0)
	assert.GreaterOrEqual(t, stats.RelationshipsCreated, 0)
	assert.NotZero(t, stats.StartTime)
} 