package neo4j

import (
	"kg-enricher/internal/models"

	neo4jdriver "github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// Neo4jService defines the interface for Neo4j operations
type Neo4jService interface {
	GetRandomNodes(batchSize int) ([]models.Node, error)
	CheckExistingRelationship(source, target string) (bool, error)
	CreateRelationship(from, to, relation string) error
	QueryAllConcepts() ([]string, error)
	QueryRelationships() ([]models.Relationship, error)
	QueryRandomConceptPairs(batchSize int) ([][2]string, error)
}

// RealNeo4jService implements the Neo4jService interface using a real Neo4j driver
type RealNeo4jService struct {
	driver neo4jdriver.Driver
}

// NewNeo4jService creates a new Neo4j service with the given driver
func NewNeo4jService(driver neo4jdriver.Driver) Neo4jService {
	return &RealNeo4jService{
		driver: driver,
	}
}

// GetRandomNodes gets random nodes from the database
func (s *RealNeo4jService) GetRandomNodes(batchSize int) ([]models.Node, error) {
	return GetRandomNodes(s.driver, batchSize)
}

// CheckExistingRelationship checks if a relationship already exists between two nodes
func (s *RealNeo4jService) CheckExistingRelationship(source, target string) (bool, error) {
	return CheckExistingRelationship(s.driver, source, target)
}

// CreateRelationship creates a relationship between two nodes
func (s *RealNeo4jService) CreateRelationship(from, to, relation string) error {
	return CreateRelationship(s.driver, from, to, relation)
}

// QueryAllConcepts queries all concepts from the database
func (s *RealNeo4jService) QueryAllConcepts() ([]string, error) {
	return QueryAllConcepts(s.driver)
}

// QueryRelationships queries all relationships from the database
func (s *RealNeo4jService) QueryRelationships() ([]models.Relationship, error) {
	return QueryRelationships(s.driver)
}

// QueryRandomConceptPairs queries random pairs of concepts from the database
func (s *RealNeo4jService) QueryRandomConceptPairs(batchSize int) ([][2]string, error) {
	return QueryRandomConceptPairs(s.driver, batchSize)
} 