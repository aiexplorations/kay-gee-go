package neo4j

import (
	"fmt"
	"kg-enricher/internal/models"
)

// MockNeo4jService is a mock implementation of the Neo4jService interface for testing
type MockNeo4jService struct {
	// Mock data
	Nodes                []models.Node
	Relationships        []models.Relationship
	ConceptPairs         [][2]string
	ExistingRelationships map[string]bool
	
	// Mock behavior flags
	ShouldFailGetRandomNodes         bool
	ShouldFailCheckExistingRelationship bool
	ShouldFailCreateRelationship     bool
	ShouldFailQueryAllConcepts       bool
	ShouldFailQueryRelationships     bool
	ShouldFailQueryRandomConceptPairs bool
}

// NewMockNeo4jService creates a new mock Neo4j service
func NewMockNeo4jService() *MockNeo4jService {
	return &MockNeo4jService{
		Nodes:                []models.Node{},
		Relationships:        []models.Relationship{},
		ConceptPairs:         [][2]string{},
		ExistingRelationships: make(map[string]bool),
	}
}

// GetRandomNodes returns mock random nodes
func (s *MockNeo4jService) GetRandomNodes(batchSize int) ([]models.Node, error) {
	if s.ShouldFailGetRandomNodes {
		return nil, ErrMockFailure
	}
	
	// If no nodes are set, create some default ones
	if len(s.Nodes) == 0 {
		s.Nodes = []models.Node{
			{Name: "Concept1"},
			{Name: "Concept2"},
			{Name: "Concept3"},
			{Name: "Concept4"},
			{Name: "Concept5"},
			{Name: "Concept6"},
		}
	}
	
	// Return at most batchSize nodes
	if batchSize < len(s.Nodes) {
		return s.Nodes[:batchSize], nil
	}
	
	return s.Nodes, nil
}

// CheckExistingRelationship checks if a relationship exists in the mock data
func (s *MockNeo4jService) CheckExistingRelationship(source, target string) (bool, error) {
	if s.ShouldFailCheckExistingRelationship {
		return false, ErrMockFailure
	}
	
	key := source + "->" + target
	exists, ok := s.ExistingRelationships[key]
	if ok {
		return exists, nil
	}
	
	return false, nil
}

// CreateRelationship creates a mock relationship
func (s *MockNeo4jService) CreateRelationship(from, to, relation string) error {
	if s.ShouldFailCreateRelationship {
		return ErrMockFailure
	}
	
	s.Relationships = append(s.Relationships, models.Relationship{
		Source: from,
		Target: to,
		Type:   relation,
	})
	
	key := from + "->" + to
	s.ExistingRelationships[key] = true
	
	return nil
}

// QueryAllConcepts returns mock concepts
func (s *MockNeo4jService) QueryAllConcepts() ([]string, error) {
	if s.ShouldFailQueryAllConcepts {
		return nil, ErrMockFailure
	}
	
	concepts := make([]string, 0, len(s.Nodes))
	for _, node := range s.Nodes {
		concepts = append(concepts, node.Name)
	}
	
	return concepts, nil
}

// QueryRelationships returns mock relationships
func (s *MockNeo4jService) QueryRelationships() ([]models.Relationship, error) {
	if s.ShouldFailQueryRelationships {
		return nil, ErrMockFailure
	}
	
	return s.Relationships, nil
}

// QueryRandomConceptPairs returns mock concept pairs
func (s *MockNeo4jService) QueryRandomConceptPairs(batchSize int) ([][2]string, error) {
	if s.ShouldFailQueryRandomConceptPairs {
		return nil, ErrMockFailure
	}
	
	// If no pairs are set, create some default ones
	if len(s.ConceptPairs) == 0 {
		s.ConceptPairs = [][2]string{
			{"Concept1", "Concept2"},
			{"Concept3", "Concept4"},
			{"Concept5", "Concept6"},
		}
	}
	
	// Return at most batchSize pairs
	if batchSize < len(s.ConceptPairs) {
		return s.ConceptPairs[:batchSize], nil
	}
	
	return s.ConceptPairs, nil
}

// ErrMockFailure is a mock error
var ErrMockFailure = fmt.Errorf("mock failure") 