package graph

import (
	"context"
	"fmt"
	"log"
	"math/rand" // Keep this import as we'll use it in getRandomPair
	"strings"
	"sync"
	"time"

	"kg-builder/internal/config"
	apperrors "kg-builder/internal/errors"
	"kg-builder/internal/models"
	kgneo4j "kg-builder/internal/neo4j"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// GraphBuilder struct
type GraphBuilder struct {
	driver             neo4j.Driver
	getRelatedConcepts func(string) ([]models.Concept, error)
	mineRelationship   func(string, string) (*models.Concept, error)
	processedConcepts  map[string]bool
	nodeCount          int
	mutex              sync.Mutex
	config             *config.GraphConfig
}

// NewGraphBuilder creates a new GraphBuilder instance
func NewGraphBuilder(driver neo4j.Driver, getRelatedConcepts func(string) ([]models.Concept, error), mineRelationship func(string, string) (*models.Concept, error), config *config.GraphConfig) *GraphBuilder {
	if driver == nil {
		log.Fatal("Neo4j driver cannot be nil")
	}
	
	if getRelatedConcepts == nil {
		log.Fatal("getRelatedConcepts function cannot be nil")
	}
	
	if mineRelationship == nil {
		log.Fatal("mineRelationship function cannot be nil")
	}
	
	return &GraphBuilder{
		driver:             driver,
		getRelatedConcepts: getRelatedConcepts,
		mineRelationship:   mineRelationship,
		processedConcepts:  make(map[string]bool),
		nodeCount:          0,
		config:             config,
	}
}

// BuildGraph builds the knowledge graph
func (gb *GraphBuilder) BuildGraph(seedConcept string, maxNodes int, timeout time.Duration) error {
	if seedConcept == "" {
		return apperrors.NewGraphError(apperrors.ErrInvalidInput, "seed concept cannot be empty")
	}
	
	if maxNodes <= 0 {
		return apperrors.NewGraphError(apperrors.ErrInvalidInput, "maxNodes must be greater than 0")
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	queue := make(chan string, maxNodes) // Create a channel to hold concepts
	queue <- seedConcept                 // Add the seed concept to the queue

	var wg sync.WaitGroup
	workerCount := gb.config.WorkerCount // Use the configured worker count

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go gb.worker(ctx, &wg, queue)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			log.Printf("Timeout reached after processing %d concepts", gb.nodeCount)
			return apperrors.NewGraphError(apperrors.ErrTimeout, "graph building timed out")
		}
		return apperrors.NewGraphError(ctx.Err(), "context cancelled")
	case <-done:
		log.Printf("Graph building completed, processed %d concepts", gb.nodeCount)
		return nil
	}
}

// BuildGraphWithLowConnectivitySeeds builds the knowledge graph using low connectivity concepts as seeds
func (gb *GraphBuilder) BuildGraphWithLowConnectivitySeeds(initialSeedConcept string, targetNodeCount int, timeout time.Duration) error {
	if initialSeedConcept == "" {
		return apperrors.NewGraphError(apperrors.ErrInvalidInput, "initial seed concept cannot be empty")
	}
	
	if targetNodeCount <= 0 {
		return apperrors.NewGraphError(apperrors.ErrInvalidInput, "targetNodeCount must be greater than 0")
	}
	
	// Start with the initial seed concept
	err := gb.BuildGraph(initialSeedConcept, gb.config.MaxNodes, timeout)
	if err != nil {
		return err
	}
	
	// Continue building the graph with low connectivity concepts until we reach the target node count
	for gb.nodeCount < targetNodeCount {
		// Get a random low connectivity concept
		lowConnectivityConcept, err := kgneo4j.GetRandomLowConnectivityConcept(gb.driver, 10)
		if err != nil {
			log.Printf("Error getting low connectivity concept: %v", err)
			return err
		}
		
		log.Printf("Using low connectivity concept as seed: %s", lowConnectivityConcept)
		
		// Build the graph with the low connectivity concept as seed
		err = gb.BuildGraph(lowConnectivityConcept, targetNodeCount, timeout)
		if err != nil {
			log.Printf("Error building graph with low connectivity concept: %v", err)
			// Continue with the next low connectivity concept even if there's an error
			continue
		}
	}
	
	return nil
}

func (gb *GraphBuilder) worker(ctx context.Context, wg *sync.WaitGroup, queue chan string) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case concept, ok := <-queue:
			if !ok {
				return
			}

			gb.mutex.Lock()
			if gb.processedConcepts[concept] || gb.nodeCount >= gb.config.MaxNodes {
				gb.mutex.Unlock()
				continue
			}
			gb.processedConcepts[concept] = true
			gb.nodeCount++
			currentNodeCount := gb.nodeCount
			gb.mutex.Unlock()

			log.Printf("Processing concept: %s (Node count: %d)", concept, currentNodeCount)

			err := gb.processRelatedConcepts(ctx, concept)
			if err != nil {
				log.Printf("Error processing related concepts for %s: %v", concept, err)
				continue
			}
		}
	}
}

// isValidConcept checks if a concept appears to be valid
func isValidConcept(concept models.Concept) bool {
	// Skip concepts with empty names
	if concept.Name == "" {
		return false
	}
	
	// Skip concepts with very short names (likely abbreviations without context)
	if len(concept.Name) < 3 {
		return false
	}
	
	// Skip concepts with unusual characters that might indicate made-up terms
	if strings.ContainsAny(concept.Name, "!@#$%^&*()_+={}[]|\\:;\"'<>,?/~`") {
		return false
	}
	
	// Skip concepts with relation types that seem generic or vague
	if concept.Relation == "" || 
	   concept.Relation == "related to" || 
	   concept.Relation == "is related to" ||
	   concept.Relation == "relates to" {
		return false
	}
	
	// Skip concepts with unusual capitalization patterns (often indicates made-up terms)
	wordCount := 0
	capitalCount := 0
	for _, word := range strings.Fields(concept.Name) {
		wordCount++
		if len(word) > 0 && word[0] >= 'A' && word[0] <= 'Z' {
			capitalCount++
		}
	}
	
	// If all words are capitalized in a multi-word concept, it might be made up
	if wordCount > 2 && capitalCount == wordCount {
		return false
	}
	
	return true
}

// processRelatedConcepts processes the related concepts for a given concept
func (gb *GraphBuilder) processRelatedConcepts(ctx context.Context, concept string) error {
	// Get related concepts from the LLM
	relatedConcepts, err := gb.getRelatedConcepts(concept)
	if err != nil {
		return err
	}
	
	// Filter out potentially invalid concepts
	var validConcepts []models.Concept
	for _, relatedConcept := range relatedConcepts {
		if isValidConcept(relatedConcept) {
			validConcepts = append(validConcepts, relatedConcept)
		} else {
			log.Printf("Filtered out potentially invalid concept: %s", relatedConcept.Name)
		}
	}
	
	// Process each valid related concept
	for _, relatedConcept := range validConcepts {
		gb.mutex.Lock()
		if gb.nodeCount >= gb.config.MaxNodes {
			gb.mutex.Unlock()
			return nil
		}
		gb.mutex.Unlock()

		log.Printf("Creating relationship: %s -[%s]-> %s", concept, relatedConcept.Relation, relatedConcept.Name)
		err := kgneo4j.CreateRelationship(gb.driver, concept, relatedConcept.Name, relatedConcept.Relation)
		if err != nil {
			log.Printf("Error creating relationship: %v", err)
			continue
		}
		log.Printf("Successfully created relationship: %s -[%s]-> %s", concept, relatedConcept.Relation, relatedConcept.Name)

		gb.mutex.Lock()
		if !gb.processedConcepts[relatedConcept.Name] && gb.nodeCount < gb.config.MaxNodes {
			select {
			case queue <- relatedConcept.Name:
			default:
				// Queue is full, skip this concept
			}
		}
		gb.mutex.Unlock()
	}

	return nil
}

// isValidRelationship checks if a relationship appears to be valid
func isValidRelationship(relationship *models.Concept) bool {
	// Skip nil relationships
	if relationship == nil {
		return false
	}
	
	// Skip relationships with empty fields
	if relationship.Name == "" || relationship.Relation == "" || relationship.RelatedTo == "" {
		return false
	}
	
	// Skip relationships with generic or vague relation types
	if relationship.Relation == "related to" || 
	   relationship.Relation == "is related to" ||
	   relationship.Relation == "relates to" {
		return false
	}
	
	return true
}

// MineRandomRelationships mines random relationships between concepts
func (gb *GraphBuilder) MineRandomRelationships(count int, concurrency int) error {
	// Get all concepts
	concepts, err := kgneo4j.GetAllConcepts(gb.driver)
	if err != nil {
		return err
	}
	
	if len(concepts) < 2 {
		return fmt.Errorf("not enough concepts to mine relationships")
	}
	
	log.Printf("Mining %d random relationships between %d concepts", count, len(concepts))
	
	// Create a channel for pairs of concepts
	pairs := make(chan [2]string, count)
	
	// Create a channel for results
	results := make(chan error, count)
	
	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for pair := range pairs {
				conceptA := pair[0]
				conceptB := pair[1]
				
				// Skip if the relationship already exists
				exists, err := kgneo4j.RelationshipExists(gb.driver, conceptA, conceptB)
				if err != nil {
					results <- err
					continue
				}
				
				if exists {
					log.Printf("Relationship already exists between %s and %s, skipping", conceptA, conceptB)
					results <- nil
					continue
				}
				
				// Mine the relationship
				relationship, err := gb.mineRelationship(conceptA, conceptB)
				if err != nil {
					results <- err
					continue
				}
				
				// Validate the relationship
				if !isValidRelationship(relationship) {
					log.Printf("Filtered out potentially invalid relationship between %s and %s", conceptA, conceptB)
					results <- nil
					continue
				}
				
				// Create the relationship
				err = kgneo4j.CreateRelationship(gb.driver, relationship.Name, relationship.RelatedTo, relationship.Relation)
				if err != nil {
					results <- err
					continue
				}
				
				log.Printf("Created relationship: %s -[%s]-> %s", relationship.Name, relationship.Relation, relationship.RelatedTo)
				results <- nil
			}
		}()
	}
	
	// Generate random pairs of concepts
	go func() {
		for i := 0; i < count; i++ {
			pair := getRandomPair(concepts)
			pairs <- pair
		}
		close(pairs)
	}()
	
	// Collect results
	var errors []error
	for i := 0; i < count; i++ {
		err := <-results
		if err != nil {
			errors = append(errors, err)
		}
	}
	
	// Wait for all workers to finish
	wg.Wait()
	
	if len(errors) > 0 {
		return fmt.Errorf("encountered %d errors while mining relationships", len(errors))
	}
	
	return nil
}

// getRandomPair returns a random pair of concepts
func getRandomPair(concepts []string) [2]string {
	if len(concepts) < 2 {
		return [2]string{"", ""}
	}
	
	// Get two random indices
	idx1 := rand.Intn(len(concepts))
	idx2 := rand.Intn(len(concepts))
	
	// Make sure they're different
	for idx1 == idx2 {
		idx2 = rand.Intn(len(concepts))
	}
	
	return [2]string{concepts[idx1], concepts[idx2]}
}

// GetRandomPair returns a random pair of concepts from the processed concepts
func (gb *GraphBuilder) GetRandomPair() [2]string {
	gb.mutex.Lock()
	defer gb.mutex.Unlock()

	concepts := make([]string, 0, len(gb.processedConcepts))
	for concept := range gb.processedConcepts {
		concepts = append(concepts, concept)
	}

	if len(concepts) < 2 {
		return [2]string{"", ""}
	}

	i := rand.Intn(len(concepts))
	j := rand.Intn(len(concepts) - 1)
	if j >= i {
		j++
	}

	return [2]string{concepts[i], concepts[j]}
}

// GetProcessedConcepts returns a copy of the processed concepts map
func (gb *GraphBuilder) GetProcessedConcepts() map[string]bool {
	gb.mutex.Lock()
	defer gb.mutex.Unlock()
	
	result := make(map[string]bool, len(gb.processedConcepts))
	for k, v := range gb.processedConcepts {
		result[k] = v
	}
	
	return result
}

// GetNodeCount returns the current node count
func (gb *GraphBuilder) GetNodeCount() int {
	gb.mutex.Lock()
	defer gb.mutex.Unlock()
	
	return gb.nodeCount
}
