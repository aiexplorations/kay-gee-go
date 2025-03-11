package models

// BuilderParams represents the parameters for starting the knowledge graph builder
type BuilderParams struct {
	SeedConcept        string `json:"seedConcept" binding:"required"`
	MaxNodes           int    `json:"maxNodes" binding:"required,min=1"`
	Timeout            int    `json:"timeout" binding:"required,min=1"`
	RandomRelationships int    `json:"randomRelationships" binding:"required,min=0"`
	Concurrency        int    `json:"concurrency" binding:"required,min=1"`
}

// EnricherParams represents the parameters for starting the knowledge graph enricher
type EnricherParams struct {
	BatchSize        int `json:"batchSize" binding:"required,min=1"`
	Interval         int `json:"interval" binding:"required,min=10"`
	MaxRelationships int `json:"maxRelationships" binding:"required,min=1"`
	Concurrency      int `json:"concurrency" binding:"required,min=1"`
}

// RelationshipCreate represents the parameters for creating a relationship between concepts
type RelationshipCreate struct {
	Source string `json:"source" binding:"required"`
	Target string `json:"target" binding:"required"`
	Type   string `json:"type" binding:"required"`
}

// Node represents a concept node in the knowledge graph
type Node struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type,omitempty"`
	Size int    `json:"size,omitempty"`
}

// Link represents a relationship between concepts in the knowledge graph
type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"`
}

// GraphData represents the complete graph data with nodes and links
type GraphData struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

// Statistics represents the statistics about the knowledge graph
type Statistics struct {
	ConceptCount      int `json:"conceptCount"`
	RelationshipCount int `json:"relationshipCount"`
}

// Response represents a generic API response
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Output  string `json:"output,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
} 