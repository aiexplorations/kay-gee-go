package enricher

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kay-gee-go/internal/common/config"
	"github.com/kay-gee-go/internal/common/errors"
	"github.com/kay-gee-go/internal/common/llm"
	"github.com/kay-gee-go/internal/common/models"
	"github.com/kay-gee-go/internal/common/neo4j"
)

// Enricher represents a knowledge graph enricher
type Enricher struct {
	config     *config.EnricherAppConfig
	neo4jClient *neo4j.Client
	llmClient  *llm.Client
	stats      *models.EnricherStats
	mutex      sync.Mutex
}

// NewEnricher creates a new knowledge graph enricher
func NewEnricher(config *config.EnricherAppConfig) (*Enricher, error) {
	// Create Neo4j client
	neo4jClient, err := neo4j.NewClient(config.Neo4j)
	if err != nil {
		return nil, err
	}

	// Create LLM client
	llmClient := llm.NewClient(config.LLM)

	return &Enricher{
		config:     config,
		neo4jClient: neo4jClient,
		llmClient:  llmClient,
		stats: &models.EnricherStats{
			PairsProcessed:       0,
			RelationshipsCreated: 0,
			StartTime:            time.Now(),
		},
	}, nil
}

// Close closes the enricher
func (e *Enricher) Close() error {
	return e.neo4jClient.Close()
}

// EnrichOnce enriches the knowledge graph once
func (e *Enricher) EnrichOnce(ctx context.Context, count int) (*models.EnricherStats, error) {
	fmt.Println("Starting knowledge graph enrichment process...")
	fmt.Printf("Batch size: %d\n", count)

	// Reset stats
	e.mutex.Lock()
	e.stats = &models.EnricherStats{
		PairsProcessed:       0,
		RelationshipsCreated: 0,
		StartTime:            time.Now(),
	}
	e.mutex.Unlock()

	// Get random pairs of concepts
	pairs, err := e.neo4jClient.GetRandomConceptPairs(count)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Found %d random pairs of concepts\n", len(pairs))

	// Process each pair
	e.processPairs(ctx, pairs)

	// Update stats
	e.mutex.Lock()
	e.stats.EndTime = time.Now()
	e.stats.Duration = e.stats.EndTime.Sub(e.stats.StartTime).String()
	e.mutex.Unlock()

	fmt.Println("Knowledge graph enrichment process completed.")
	fmt.Printf("Pairs processed: %d\n", e.stats.PairsProcessed)
	fmt.Printf("Relationships created: %d\n", e.stats.RelationshipsCreated)
	fmt.Printf("Duration: %s\n", e.stats.Duration)

	return e.stats, nil
}

// EnrichContinuously enriches the knowledge graph continuously
func (e *Enricher) EnrichContinuously(ctx context.Context) error {
	fmt.Println("Starting continuous knowledge graph enrichment process...")
	fmt.Printf("Batch size: %d\n", e.config.Enricher.BatchSize)
	fmt.Printf("Interval: %d seconds\n", e.config.Enricher.IntervalSeconds)
	fmt.Printf("Max relationships: %d\n", e.config.Enricher.MaxRelationships)

	// Reset stats
	e.mutex.Lock()
	e.stats = &models.EnricherStats{
		PairsProcessed:       0,
		RelationshipsCreated: 0,
		StartTime:            time.Now(),
	}
	e.mutex.Unlock()

	ticker := time.NewTicker(time.Duration(e.config.Enricher.IntervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// Check if we've reached the maximum number of relationships
			e.mutex.Lock()
			if e.config.Enricher.MaxRelationships > 0 && e.stats.RelationshipsCreated >= e.config.Enricher.MaxRelationships {
				e.mutex.Unlock()
				fmt.Printf("Reached maximum number of relationships (%d)\n", e.config.Enricher.MaxRelationships)
				return nil
			}
			e.mutex.Unlock()

			// Get random pairs of concepts
			pairs, err := e.neo4jClient.GetRandomConceptPairs(e.config.Enricher.BatchSize)
			if err != nil {
				fmt.Printf("Error getting random pairs: %s\n", err.Error())
				continue
			}

			fmt.Printf("Found %d random pairs of concepts\n", len(pairs))

			// Process each pair
			e.processPairs(ctx, pairs)

			// Print stats
			e.mutex.Lock()
			fmt.Printf("Pairs processed: %d\n", e.stats.PairsProcessed)
			fmt.Printf("Relationships created: %d\n", e.stats.RelationshipsCreated)
			e.mutex.Unlock()
		}
	}
}

// processPairs processes a batch of concept pairs
func (e *Enricher) processPairs(ctx context.Context, pairs [][]models.Concept) {
	if len(pairs) == 0 {
		return
	}

	// Create a worker pool
	concurrency := e.config.Enricher.Concurrency
	if concurrency <= 0 {
		concurrency = 5
	}

	// Create channels for the worker pool
	type conceptPair struct {
		source models.Concept
		target models.Concept
	}

	jobs := make(chan conceptPair, concurrency)
	results := make(chan error, concurrency)

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for pair := range jobs {
				// Check if the context is done
				select {
				case <-ctx.Done():
					results <- errors.NewTimeoutError("context deadline exceeded", ctx.Err())
					continue
				default:
					// Continue processing
				}

				// Get the relationship between the two concepts
				relationship, err := e.llmClient.GetRelationship(pair.source.Name, pair.target.Name)
				if err != nil {
					results <- err
					continue
				}

				if relationship == nil {
					results <- nil
					continue
				}

				// Create the relationship
				_, err = e.neo4jClient.CreateRelationship(*relationship)
				if err != nil {
					results <- err
					continue
				}

				e.mutex.Lock()
				e.stats.RelationshipsCreated++
				e.mutex.Unlock()

				fmt.Printf("Created relationship: %s -> %s (%s)\n", pair.source.Name, pair.target.Name, relationship.Type)

				results <- nil
			}
		}()
	}

	// Process each pair
	for _, pair := range pairs {
		// Check if the context is done
		select {
		case <-ctx.Done():
			close(jobs)
			wg.Wait()
			return
		default:
			// Continue processing
		}

		// Send the pair to a worker
		jobs <- conceptPair{
			source: pair[0],
			target: pair[1],
		}

		// Wait for the result
		err := <-results
		if err != nil {
			fmt.Printf("Error processing pair %s and %s: %s\n", pair[0].Name, pair[1].Name, err.Error())
		}

		e.mutex.Lock()
		e.stats.PairsProcessed++
		e.mutex.Unlock()
	}

	close(jobs)
	wg.Wait()
}

// GetStats returns the current enricher stats
func (e *Enricher) GetStats() *models.EnricherStats {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return e.stats
} 