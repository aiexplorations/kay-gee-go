package graph

import (
	"context"
	"log"
	"math/rand" // Keep this import as we'll use it in getRandomPair
	"sync"
	"time"

	"kg-builder/internal/models"
	kgneo4j "kg-builder/internal/neo4j"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

const maxNodes = 100

// GraphBuilder struct
type GraphBuilder struct {
	driver             neo4j.Driver
	getRelatedConcepts func(string) ([]models.Concept, error)
	mineRelationship   func(string, string) (*models.Concept, error)
	processedConcepts  map[string]bool
	nodeCount          int
	mutex              sync.Mutex
}

// NewGraphBuilder creates a new GraphBuilder instance
func NewGraphBuilder(driver neo4j.Driver, getRelatedConcepts func(string) ([]models.Concept, error), mineRelationship func(string, string) (*models.Concept, error)) *GraphBuilder {
	return &GraphBuilder{
		driver:             driver,
		getRelatedConcepts: getRelatedConcepts,
		mineRelationship:   mineRelationship,
		processedConcepts:  make(map[string]bool),
		nodeCount:          0,
	}
}

// BuildGraph builds the knowledge graph
func (gb *GraphBuilder) BuildGraph(seedConcept string, maxNodes int, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	queue := make(chan string, maxNodes) // Create a channel to hold concepts
	queue <- seedConcept                 // Add the seed concept to the queue

	var wg sync.WaitGroup
	workerCount := 10 // Adjust this number based on your needs and system capabilities

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
		}
	case <-done:
		log.Printf("Graph building completed, processed %d concepts", gb.nodeCount)
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
			if gb.processedConcepts[concept] || gb.nodeCount >= maxNodes {
				gb.mutex.Unlock()
				continue
			}
			gb.processedConcepts[concept] = true
			gb.nodeCount++
			currentNodeCount := gb.nodeCount
			gb.mutex.Unlock()

			log.Printf("Processing concept: %s (Node count: %d)", concept, currentNodeCount)

			relatedConcepts, err := gb.getRelatedConcepts(concept)
			if err != nil {
				log.Printf("Error getting related concepts for %s: %v", concept, err)
				continue
			}

			log.Printf("Found %d related concepts for %s", len(relatedConcepts), concept)
			for _, rc := range relatedConcepts {
				gb.mutex.Lock()
				if gb.nodeCount >= maxNodes {
					gb.mutex.Unlock()
					return
				}
				gb.mutex.Unlock()

				log.Printf("Creating relationship: %s -[%s]-> %s", concept, rc.Relation, rc.Name)
				err := kgneo4j.CreateRelationship(gb.driver, concept, rc.Name, rc.Relation)
				if err != nil {
					log.Printf("Error creating relationship: %v", err)
					continue
				}
				log.Printf("Successfully created relationship: %s -[%s]-> %s", concept, rc.Relation, rc.Name)

				gb.mutex.Lock()
				if !gb.processedConcepts[rc.Name] && gb.nodeCount < maxNodes {
					select {
					case queue <- rc.Name:
					default:
						// Queue is full, skip this concept
					}
				}
				gb.mutex.Unlock()
			}
		}
	}
}

func (gb *GraphBuilder) MineRandomRelationships(count int, concurrency int) {
	semaphore := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for i := 0; i < count; i++ {
		wg.Add(1)
		semaphore <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-semaphore }()

			concepts := gb.getRandomPair()
			if concepts[0] == concepts[1] {
				return
			}

			log.Printf("Mining relationship between %s and %s", concepts[0], concepts[1])
			concept, err := gb.mineRelationship(concepts[0], concepts[1])
			if err != nil {
				log.Printf("Error mining relationship: %v", err)
				return
			}

			if concept == nil {
				log.Printf("No relationship found between %s and %s", concepts[0], concepts[1])
				return
			}

			log.Printf("Creating relationship: %s -[%s]-> %s", concepts[0], concept.Relation, concepts[1])
			err = kgneo4j.CreateRelationship(gb.driver, concepts[0], concepts[1], concept.Relation)
			if err != nil {
				log.Printf("Error creating relationship: %v", err)
				return
			}
			log.Printf("Successfully created relationship: %s -[%s]-> %s", concepts[0], concept.Relation, concepts[1])
		}()
	}

	wg.Wait()
}

func (gb *GraphBuilder) getRandomPair() [2]string {
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
