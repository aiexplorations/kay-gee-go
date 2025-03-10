package models

// Concept represents a concept in the knowledge graph
type Concept struct {
	Name      string `json:"name"`
	Relation  string `json:"relation"`
	RelatedTo string `json:"relatedTo"`
}

// Node represents a node in the Neo4j database
type Node struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Label string `json:"label"`
}

// Relationship represents a relationship between two nodes
type Relationship struct {
	Source      string `json:"source"`
	Target      string `json:"target"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
} 