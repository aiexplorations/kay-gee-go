package builder

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

// Builder represents a knowledge graph builder
type Builder struct {
	config     *config.BuilderConfig
	neo4jClient *neo4j.Client
	llmClient  *llm.Client
	stats      *models.BuilderStats
	mutex      sync.Mutex
	concepts   map[string]bool
}

// NewBuilder creates a new knowledge graph builder
func NewBuilder(config *config.BuilderConfig) (*Builder, error) {
	// Create Neo4j client
	neo4jClient, err := neo4j.NewClient(config.Neo4j)
	if err != nil {
		return nil, err
	}

	// Fix duplicate nodes first
	if err := neo4jClient.FixDuplicateNodes(); err != nil {
		fmt.Printf("Warning: Error fixing duplicate nodes: %s\n", err.Error())
		// Continue anyway, as this is just a best-effort attempt
	}

	// Initialize Neo4j schema
	if err := neo4jClient.InitializeSchema(); err != nil {
		return nil, err
	}

	// Create LLM client
	llmClient := llm.NewClient(config.LLM)

	return &Builder{
		config:     config,
		neo4jClient: neo4jClient,
		llmClient:  llmClient,
		stats: &models.BuilderStats{
			SeedConcept:  config.Graph.SeedConcept,
			NodesCreated: 0,
			RelationshipsCreated: 0,
			StartTime:    time.Now(),
		},
		concepts: make(map[string]bool),
	}, nil
}

// Close closes the builder
func (b *Builder) Close() error {
	return b.neo4jClient.Close()
}

// Build builds the knowledge graph
func (b *Builder) Build(ctx context.Context) (*models.BuilderStats, error) {
	fmt.Println("Starting knowledge graph building process...")
	fmt.Printf("Seed concept: %s\n", b.config.Graph.SeedConcept)
	fmt.Printf("Max nodes: %d\n", b.config.Graph.MaxNodes)
	fmt.Printf("Timeout: %d minutes\n", b.config.Graph.TimeoutMinutes)

	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(b.config.Graph.TimeoutMinutes)*time.Minute)
	defer cancel()

	// Create the seed concept
	seedConcept := models.Concept{
		Name:        b.config.Graph.SeedConcept,
		Description: fmt.Sprintf("Seed concept for the knowledge graph: %s", b.config.Graph.SeedConcept),
	}

	seedID, err := b.neo4jClient.CreateConcept(seedConcept)
	if err != nil {
		return nil, err
	}

	b.mutex.Lock()
	b.concepts[seedConcept.Name] = true
	b.stats.NodesCreated++
	b.mutex.Unlock()

	fmt.Printf("Created seed concept: %s (ID: %s)\n", seedConcept.Name, seedID)

	// Create a queue for BFS
	queue := []string{seedConcept.Name}
	
	// Create a worker pool
	workerCount := b.config.Graph.WorkerCount
	if workerCount <= 0 {
		workerCount = 5
	}
	
	// Create channels for the worker pool
	jobs := make(chan string, workerCount)
	results := make(chan error, workerCount)
	
	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for concept := range jobs {
				err := b.processRelatedConcepts(timeoutCtx, concept)
				results <- err
			}
		}()
	}
	
	// Process the queue
	go func() {
		for len(queue) > 0 && b.stats.NodesCreated < b.config.Graph.MaxNodes {
			// Check if the context is done
			select {
			case <-timeoutCtx.Done():
				close(jobs)
				return
			default:
				// Get the next concept from the queue
				concept := queue[0]
				queue = queue[1:]
				
				// Send the concept to a worker
				jobs <- concept
				
				// Wait for the result
				err := <-results
				if err != nil {
					fmt.Printf("Error processing related concepts for %s: %s\n", concept, err.Error())
				}
				
				// Check if we've reached the maximum number of nodes
				if b.stats.NodesCreated >= b.config.Graph.MaxNodes {
					close(jobs)
					break
				}
			}
		}
		close(jobs)
	}()
	
	// Wait for all workers to finish
	wg.Wait()
	
	// Mine random relationships
	if b.config.Graph.RandomRelationships > 0 {
		fmt.Printf("Mining %d random relationships...\n", b.config.Graph.RandomRelationships)
		b.mineRandomRelationships(timeoutCtx, b.config.Graph.RandomRelationships)
	}
	
	// Update stats
	b.mutex.Lock()
	b.stats.EndTime = time.Now()
	b.stats.Duration = b.stats.EndTime.Sub(b.stats.StartTime).String()
	b.mutex.Unlock()
	
	fmt.Println("Knowledge graph building process completed.")
	fmt.Printf("Nodes created: %d\n", b.stats.NodesCreated)
	fmt.Printf("Relationships created: %d\n", b.stats.RelationshipsCreated)
	fmt.Printf("Duration: %s\n", b.stats.Duration)
	
	return b.stats, nil
}

// processRelatedConcepts processes the concepts related to the given concept
func (b *Builder) processRelatedConcepts(ctx context.Context, concept string) error {
	// Check if the context is done
	select {
	case <-ctx.Done():
		return errors.NewTimeoutError("context deadline exceeded", ctx.Err())
	default:
		// Continue processing
	}
	
	// Get related concepts from the LLM
	relatedConcepts, err := b.llmClient.GetRelatedConcepts(concept)
	if err != nil {
		return err
	}
	
	fmt.Printf("Found %d concepts related to %s\n", len(relatedConcepts), concept)
	
	// Process each related concept
	for _, relatedConcept := range relatedConcepts {
		// Check if the context is done
		select {
		case <-ctx.Done():
			return errors.NewTimeoutError("context deadline exceeded", ctx.Err())
		default:
			// Continue processing
		}
		
		// Check if we've reached the maximum number of nodes
		b.mutex.Lock()
		if b.stats.NodesCreated >= b.config.Graph.MaxNodes {
			b.mutex.Unlock()
			return nil
		}
		
		// Check if the concept already exists
		if b.concepts[relatedConcept.Name] {
			b.mutex.Unlock()
			continue
		}
		b.mutex.Unlock()
		
		// Create the concept
		newConcept := models.Concept{
			Name:        relatedConcept.Name,
			Description: relatedConcept.Description,
		}
		
		_, err := b.neo4jClient.CreateConcept(newConcept)
		if err != nil {
			fmt.Printf("Error creating concept %s: %s\n", newConcept.Name, err.Error())
			continue
		}
		
		b.mutex.Lock()
		b.concepts[newConcept.Name] = true
		b.stats.NodesCreated++
		b.mutex.Unlock()
		
		fmt.Printf("Created concept: %s\n", newConcept.Name)
		
		// Create a relationship between the concept and the related concept
		relationship := models.Relationship{
			SourceID:    concept,
			TargetID:    newConcept.Name,
			Type:        "RELATED_TO",
			Description: fmt.Sprintf("%s is related to %s", concept, newConcept.Name),
			Strength:    relatedConcept.Relevance,
		}
		
		_, err = b.neo4jClient.CreateRelationship(relationship)
		if err != nil {
			fmt.Printf("Error creating relationship between %s and %s: %s\n", concept, newConcept.Name, err.Error())
			continue
		}
		
		b.mutex.Lock()
		b.stats.RelationshipsCreated++
		b.mutex.Unlock()
		
		fmt.Printf("Created relationship: %s -> %s\n", concept, newConcept.Name)
	}
	
	return nil
}

// mineRandomRelationships mines random relationships between existing concepts
func (b *Builder) mineRandomRelationships(ctx context.Context, count int) {
	// Get all concepts
	concepts, err := b.neo4jClient.GetRandomConcepts(b.config.Graph.MaxNodes)
	if err != nil {
		fmt.Printf("Error getting concepts: %s\n", err.Error())
		return
	}
	
	if len(concepts) < 2 {
		fmt.Println("Not enough concepts to mine relationships")
		return
	}
	
	// Create a worker pool
	concurrency := b.config.Graph.Concurrency
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
				relationship, err := b.llmClient.GetRelationship(pair.source.Name, pair.target.Name)
				if err != nil {
					results <- err
					continue
				}
				
				if relationship == nil {
					results <- nil
					continue
				}
				
				// Create the relationship
				_, err = b.neo4jClient.CreateRelationship(*relationship)
				if err != nil {
					results <- err
					continue
				}
				
				b.mutex.Lock()
				b.stats.RelationshipsCreated++
				b.mutex.Unlock()
				
				fmt.Printf("Created relationship: %s -> %s (%s)\n", pair.source.Name, pair.target.Name, relationship.Type)
				
				results <- nil
			}
		}()
	}
	
	// Generate random pairs of concepts
	processed := 0
	for i := 0; i < len(concepts) && processed < count; i++ {
		for j := i + 1; j < len(concepts) && processed < count; j++ {
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
				source: concepts[i],
				target: concepts[j],
			}
			
			// Wait for the result
			err := <-results
			if err != nil {
				fmt.Printf("Error mining relationship between %s and %s: %s\n", concepts[i].Name, concepts[j].Name, err.Error())
			}
			
			processed++
		}
	}
	
	close(jobs)
	wg.Wait()
} 