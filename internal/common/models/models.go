package models

import (
	"time"
)

// Concept represents a node in the knowledge graph
type Concept struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// Relationship represents an edge between two concepts in the knowledge graph
type Relationship struct {
	ID          string    `json:"id"`
	SourceID    string    `json:"source_id"`
	TargetID    string    `json:"target_id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Strength    float64   `json:"strength"`
	CreatedAt   time.Time `json:"created_at"`
}

// LLMResponse represents a response from the LLM service
type LLMResponse struct {
	RelatedConcepts []RelatedConcept `json:"related_concepts"`
	Relationships   []Relationship   `json:"relationships"`
}

// RelatedConcept represents a concept related to another concept
type RelatedConcept struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Relevance   float64 `json:"relevance"`
}

// GraphStats represents statistics about the knowledge graph
type GraphStats struct {
	NodeCount         int       `json:"node_count"`
	RelationshipCount int       `json:"relationship_count"`
	LastUpdated       time.Time `json:"last_updated"`
}

// BuilderStats represents statistics about the graph building process
type BuilderStats struct {
	SeedConcept       string    `json:"seed_concept"`
	NodesCreated      int       `json:"nodes_created"`
	RelationshipsCreated int    `json:"relationships_created"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	Duration          string    `json:"duration"`
}

// EnricherStats represents statistics about the graph enrichment process
type EnricherStats struct {
	PairsProcessed       int       `json:"pairs_processed"`
	RelationshipsCreated int       `json:"relationships_created"`
	StartTime            time.Time `json:"start_time"`
	EndTime              time.Time `json:"end_time"`
	Duration             string    `json:"duration"`
} 