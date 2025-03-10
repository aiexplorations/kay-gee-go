package enricher

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"kg-enricher/internal/config"
	"kg-enricher/internal/llm"
	"kg-enricher/internal/neo4j"

	neo4jdriver "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// Enricher is responsible for enriching the knowledge graph with new relationships
type Enricher struct {
	driver      neo4jdriver.Driver
	config      *config.EnricherConfig
	stopChan    chan struct{}
	wg          sync.WaitGroup
	running     bool
	mutex       sync.Mutex
	stats       EnricherStats
}

// EnricherStats holds statistics about the enricher
type EnricherStats struct {
	TotalBatches          int
	TotalPairsProcessed   int
	TotalRelationsFound   int
	TotalRelationsCreated int
	StartTime             time.Time
	LastBatchTime         time.Time
}

// NewEnricher creates a new Enricher instance
func NewEnricher(driver neo4jdriver.Driver, config *config.EnricherConfig) *Enricher {
	return &Enricher{
		driver:   driver,
		config:   config,
		stopChan: make(chan struct{}),
		stats: EnricherStats{
			StartTime: time.Now(),
		},
	}
}

// Start starts the enricher service
func (e *Enricher) Start() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.running {
		return fmt.Errorf("enricher is already running")
	}

	e.running = true
	e.stopChan = make(chan struct{})
	e.wg.Add(1)

	go e.run()

	log.Printf("Enricher started with batch size %d, interval %s, max relationships %d, concurrency %d",
		e.config.BatchSize, e.config.Interval, e.config.MaxRelationships, e.config.Concurrency)

	return nil
}

// Stop stops the enricher service
func (e *Enricher) Stop() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if !e.running {
		return fmt.Errorf("enricher is not running")
	}

	close(e.stopChan)
	e.wg.Wait()
	e.running = false

	log.Printf("Enricher stopped")
	return nil
}

// IsRunning returns true if the enricher is running
func (e *Enricher) IsRunning() bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return e.running
}

// GetStats returns the current statistics
func (e *Enricher) GetStats() EnricherStats {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return e.stats
}

// run is the main loop of the enricher
func (e *Enricher) run() {
	defer e.wg.Done()

	ticker := time.NewTicker(e.config.Interval)
	defer ticker.Stop()

	// Process one batch immediately
	e.processBatch()

	for {
		select {
		case <-ticker.C:
			e.processBatch()
		case <-e.stopChan:
			return
		}
	}
}

// processBatch processes a batch of random node pairs
func (e *Enricher) processBatch() {
	log.Printf("Processing batch of random node pairs...")
	e.mutex.Lock()
	e.stats.TotalBatches++
	e.stats.LastBatchTime = time.Now()
	e.mutex.Unlock()

	// Get random nodes
	nodes, err := neo4j.GetRandomNodes(e.driver, e.config.BatchSize*2)
	if err != nil {
		log.Printf("Failed to get random nodes: %v", err)
		return
	}

	if len(nodes) < 2 {
		log.Printf("Not enough nodes to process (got %d, need at least 2)", len(nodes))
		return
	}

	// Shuffle the nodes
	rand.Shuffle(len(nodes), func(i, j int) {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})

	// Create pairs of nodes
	var pairs [][2]string
	for i := 0; i < len(nodes)-1 && len(pairs) < e.config.BatchSize; i += 2 {
		// Skip if the nodes are the same
		if nodes[i].Name == nodes[i+1].Name {
			continue
		}

		// Check if a relationship already exists
		exists, err := neo4j.CheckExistingRelationship(e.driver, nodes[i].Name, nodes[i+1].Name)
		if err != nil {
			log.Printf("Failed to check existing relationship: %v", err)
			continue
		}

		if !exists {
			pairs = append(pairs, [2]string{nodes[i].Name, nodes[i+1].Name})
		}
	}

	log.Printf("Processing %d pairs of nodes", len(pairs))

	// Process pairs concurrently
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, e.config.Concurrency)
	var relationsFound int
	var relationsCreated int
	var mutex sync.Mutex

	for _, pair := range pairs {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(source, target string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			// Check if we should stop
			select {
			case <-e.stopChan:
				return
			default:
				// Continue processing
			}

			// Mine relationship
			log.Printf("Mining relationship between %s and %s", source, target)
			relationship, err := llm.MineRelationship(source, target)
			if err != nil {
				log.Printf("Failed to mine relationship: %v", err)
				return
			}

			mutex.Lock()
			e.stats.TotalPairsProcessed++
			mutex.Unlock()

			// If no relationship was found, skip
			if relationship == nil {
				log.Printf("No relationship found between %s and %s", source, target)
				return
			}

			mutex.Lock()
			relationsFound++
			e.stats.TotalRelationsFound++
			mutex.Unlock()

			// Create relationship in Neo4j
			err = neo4j.CreateRelationship(e.driver, source, target, relationship.Relation)
			if err != nil {
				log.Printf("Failed to create relationship: %v", err)
				return
			}

			mutex.Lock()
			relationsCreated++
			e.stats.TotalRelationsCreated++
			mutex.Unlock()

			log.Printf("Created relationship: %s -[%s]-> %s", source, relationship.Relation, target)
		}(pair[0], pair[1])
	}

	wg.Wait()

	log.Printf("Batch completed: %d pairs processed, %d relations found, %d relations created",
		len(pairs), relationsFound, relationsCreated)
}

// RunOnce runs the enricher once and then stops
func (e *Enricher) RunOnce(count int) error {
	if count <= 0 {
		return fmt.Errorf("count must be greater than 0")
	}

	e.mutex.Lock()
	if e.running {
		e.mutex.Unlock()
		return fmt.Errorf("enricher is already running")
	}
	e.running = true
	e.mutex.Unlock()

	defer func() {
		e.mutex.Lock()
		e.running = false
		e.mutex.Unlock()
	}()

	// Get random nodes
	nodes, err := neo4j.GetRandomNodes(e.driver, count*2)
	if err != nil {
		return fmt.Errorf("failed to get random nodes: %w", err)
	}

	if len(nodes) < 2 {
		return fmt.Errorf("not enough nodes to process (got %d, need at least 2)", len(nodes))
	}

	// Shuffle the nodes
	rand.Shuffle(len(nodes), func(i, j int) {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})

	// Create pairs of nodes
	var pairs [][2]string
	for i := 0; i < len(nodes)-1 && len(pairs) < count; i += 2 {
		// Skip if the nodes are the same
		if nodes[i].Name == nodes[i+1].Name {
			continue
		}

		// Check if a relationship already exists
		exists, err := neo4j.CheckExistingRelationship(e.driver, nodes[i].Name, nodes[i+1].Name)
		if err != nil {
			log.Printf("Failed to check existing relationship: %v", err)
			continue
		}

		if !exists {
			pairs = append(pairs, [2]string{nodes[i].Name, nodes[i+1].Name})
		}
	}

	log.Printf("Processing %d pairs of nodes", len(pairs))

	// Process pairs concurrently
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, e.config.Concurrency)
	var relationsFound int
	var relationsCreated int
	var mutex sync.Mutex

	for _, pair := range pairs {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(source, target string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			// Mine relationship
			log.Printf("Mining relationship between %s and %s", source, target)
			relationship, err := llm.MineRelationship(source, target)
			if err != nil {
				log.Printf("Failed to mine relationship: %v", err)
				return
			}

			mutex.Lock()
			e.stats.TotalPairsProcessed++
			mutex.Unlock()

			// If no relationship was found, skip
			if relationship == nil {
				log.Printf("No relationship found between %s and %s", source, target)
				return
			}

			mutex.Lock()
			relationsFound++
			e.stats.TotalRelationsFound++
			mutex.Unlock()

			// Create relationship in Neo4j
			err = neo4j.CreateRelationship(e.driver, source, target, relationship.Relation)
			if err != nil {
				log.Printf("Failed to create relationship: %v", err)
				return
			}

			mutex.Lock()
			relationsCreated++
			e.stats.TotalRelationsCreated++
			mutex.Unlock()

			log.Printf("Created relationship: %s -[%s]-> %s", source, relationship.Relation, target)
		}(pair[0], pair[1])
	}

	wg.Wait()

	log.Printf("Run completed: %d pairs processed, %d relations found, %d relations created",
		len(pairs), relationsFound, relationsCreated)

	return nil
} 