package graph

import (
	"fmt"
	"log"
	"sync"
	"time"

	"kg-builder/internal/models"
	kgneo4j "kg-builder/internal/neo4j"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type GraphBuilder struct {
	driver             neo4j.Driver
	getRelatedConcepts func(string) ([]models.Concept, error)
	processedConcepts  map[string]bool
	mutex              sync.Mutex
}

func NewGraphBuilder(driver neo4j.Driver, getRelatedConcepts func(string) ([]models.Concept, error)) *GraphBuilder {
	return &GraphBuilder{
		driver:             driver,
		getRelatedConcepts: getRelatedConcepts,
		processedConcepts:  make(map[string]bool),
	}
}

func (gb *GraphBuilder) BuildGraph(seedConcept string, maxNodes int, timeout time.Duration) error {
	queue := []string{seedConcept}
	startTime := time.Now()

	for len(gb.processedConcepts) < maxNodes && len(queue) > 0 {
		concept := queue[0]
		queue = queue[1:]

		gb.mutex.Lock()
		if gb.processedConcepts[concept] {
			gb.mutex.Unlock()
			continue
		}
		gb.processedConcepts[concept] = true
		gb.mutex.Unlock()

		relatedConcepts, err := gb.getRelatedConcepts(concept)
		if err != nil {
			log.Printf("Error getting related concepts for %s: %v", concept, err)
			continue
		}

		for _, rc := range relatedConcepts {
			err := kgneo4j.CreateRelationship(gb.driver, concept, rc.Name, rc.Relation)
			if err != nil {
				log.Printf("Error creating relationship: %v", err)
				continue
			}

			gb.mutex.Lock()
			if !gb.processedConcepts[rc.Name] {
				queue = append(queue, rc.Name)
			}
			gb.mutex.Unlock()
		}

		if time.Since(startTime) > timeout {
			return fmt.Errorf("timeout reached after processing %d concepts", len(gb.processedConcepts))
		}
	}

	return nil
}
