package neo4j

import (
	"fmt"
	"strings"
	"time"

	"github.com/kay-gee-go/internal/common/config"
	"github.com/kay-gee-go/internal/common/errors"
	"github.com/kay-gee-go/internal/common/models"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// Client represents a Neo4j client
type Client struct {
	driver  neo4j.Driver
	config  config.Neo4jConfig
	session neo4j.Session
}

// NewClient creates a new Neo4j client
func NewClient(config config.Neo4jConfig) (*Client, error) {
	var driver neo4j.Driver
	var err error

	// Try to connect with retries
	for i := 0; i < config.MaxRetries; i++ {
		driver, err = neo4j.NewDriver(
			config.URI,
			neo4j.BasicAuth(config.User, config.Password, ""),
		)
		
		if err == nil {
			// Test the connection
			session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
			_, err = session.Run("RETURN 1", nil)
			session.Close()
			
			if err == nil {
				break
			}
		}
		
		fmt.Printf("Failed to connect to Neo4j (attempt %d/%d): %s\n", i+1, config.MaxRetries, err.Error())
		
		if i < config.MaxRetries-1 {
			time.Sleep(time.Duration(config.RetryIntervalSecs) * time.Second)
		}
	}
	
	if err != nil {
		return nil, errors.NewDatabaseError("failed to connect to Neo4j", err)
	}
	
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	
	return &Client{
		driver:  driver,
		config:  config,
		session: session,
	}, nil
}

// Close closes the Neo4j client
func (c *Client) Close() error {
	if c.session != nil {
		c.session.Close()
	}
	
	if c.driver != nil {
		return c.driver.Close()
	}
	
	return nil
}

// InitializeSchema initializes the Neo4j schema
func (c *Client) InitializeSchema() error {
	// Create constraints
	constraints := []string{
		"CREATE CONSTRAINT concept_name IF NOT EXISTS ON (c:Concept) ASSERT c.name IS UNIQUE",
	}
	
	for _, constraint := range constraints {
		_, err := c.session.Run(constraint, nil)
		if err != nil {
			// Check if the error is because the constraint already exists
			if strings.Contains(err.Error(), "already exists") {
				fmt.Printf("Constraint already exists: %s\n", constraint)
				continue
			}
			
			// Check if the error is because of duplicate nodes
			if strings.Contains(err.Error(), "Constraint") && strings.Contains(err.Error(), "Complexity Theory") {
				fmt.Println("Error creating constraint due to duplicate 'Complexity Theory' nodes. This will be handled separately.")
				continue
			}
			
			return errors.NewDatabaseError("failed to create constraint", err)
		}
	}
	
	return nil
}

// CreateConcept creates a new concept in the database
func (c *Client) CreateConcept(concept models.Concept) (string, error) {
	result, err := c.session.Run(
		`
		MERGE (c:Concept {name: $name})
		ON CREATE SET c.description = $description, c.created_at = datetime()
		RETURN ID(c) as id
		`,
		map[string]interface{}{
			"name":        concept.Name,
			"description": concept.Description,
		},
	)
	
	if err != nil {
		return "", errors.NewDatabaseError("failed to create concept", err)
	}
	
	if result.Next() {
		record := result.Record()
		id, _ := record.Get("id")
		return fmt.Sprintf("%v", id), nil
	}
	
	return "", errors.NewDatabaseError("failed to get concept ID", nil)
}

// CreateRelationship creates a new relationship between two concepts
func (c *Client) CreateRelationship(relationship models.Relationship) (string, error) {
	result, err := c.session.Run(
		`
		MATCH (source:Concept {name: $source_name})
		MATCH (target:Concept {name: $target_name})
		WHERE source <> target
		MERGE (source)-[r:RELATED_TO {type: $type}]->(target)
		ON CREATE SET r.description = $description, r.strength = $strength, r.created_at = datetime()
		RETURN ID(r) as id
		`,
		map[string]interface{}{
			"source_name": relationship.SourceID,
			"target_name": relationship.TargetID,
			"type":        relationship.Type,
			"description": relationship.Description,
			"strength":    relationship.Strength,
		},
	)
	
	if err != nil {
		return "", errors.NewDatabaseError("failed to create relationship", err)
	}
	
	if result.Next() {
		record := result.Record()
		id, _ := record.Get("id")
		return fmt.Sprintf("%v", id), nil
	}
	
	return "", errors.NewDatabaseError("failed to get relationship ID", nil)
}

// GetConceptByName retrieves a concept by name
func (c *Client) GetConceptByName(name string) (*models.Concept, error) {
	result, err := c.session.Run(
		`
		MATCH (c:Concept {name: $name})
		RETURN ID(c) as id, c.name as name, c.description as description, c.created_at as created_at
		`,
		map[string]interface{}{
			"name": name,
		},
	)
	
	if err != nil {
		return nil, errors.NewDatabaseError("failed to get concept", err)
	}
	
	if result.Next() {
		record := result.Record()
		id, _ := record.Get("id")
		name, _ := record.Get("name")
		description, _ := record.Get("description")
		createdAt, _ := record.Get("created_at")
		
		return &models.Concept{
			ID:          fmt.Sprintf("%v", id),
			Name:        name.(string),
			Description: description.(string),
			CreatedAt:   createdAt.(time.Time),
		}, nil
	}
	
	return nil, nil
}

// GetRandomConcepts retrieves a random set of concepts
func (c *Client) GetRandomConcepts(limit int) ([]models.Concept, error) {
	result, err := c.session.Run(
		`
		MATCH (c:Concept)
		RETURN ID(c) as id, c.name as name, c.description as description, c.created_at as created_at
		ORDER BY rand()
		LIMIT $limit
		`,
		map[string]interface{}{
			"limit": limit,
		},
	)
	
	if err != nil {
		return nil, errors.NewDatabaseError("failed to get random concepts", err)
	}
	
	concepts := []models.Concept{}
	
	for result.Next() {
		record := result.Record()
		id, _ := record.Get("id")
		name, _ := record.Get("name")
		description, _ := record.Get("description")
		createdAt, _ := record.Get("created_at")
		
		// Handle nil values safely
		nameStr := ""
		if name != nil {
			nameStr = name.(string)
		}
		
		descStr := ""
		if description != nil {
			descStr = description.(string)
		}
		
		var createdTime time.Time
		if createdAt != nil {
			createdTime = createdAt.(time.Time)
		} else {
			createdTime = time.Now()
		}
		
		concepts = append(concepts, models.Concept{
			ID:          fmt.Sprintf("%v", id),
			Name:        nameStr,
			Description: descStr,
			CreatedAt:   createdTime,
		})
	}
	
	return concepts, nil
}

// GetRandomConceptPairs retrieves random pairs of concepts
func (c *Client) GetRandomConceptPairs(limit int) ([][]models.Concept, error) {
	result, err := c.session.Run(
		`
		MATCH (c1:Concept), (c2:Concept)
		WHERE c1 <> c2
		AND NOT (c1)-[:RELATED_TO]->(c2)
		AND NOT (c2)-[:RELATED_TO]->(c1)
		RETURN 
			ID(c1) as id1, c1.name as name1, c1.description as description1, c1.created_at as created_at1,
			ID(c2) as id2, c2.name as name2, c2.description as description2, c2.created_at as created_at2
		ORDER BY rand()
		LIMIT $limit
		`,
		map[string]interface{}{
			"limit": limit,
		},
	)
	
	if err != nil {
		return nil, errors.NewDatabaseError("failed to get random concept pairs", err)
	}
	
	pairs := [][]models.Concept{}
	
	for result.Next() {
		record := result.Record()
		
		id1, _ := record.Get("id1")
		name1, _ := record.Get("name1")
		description1, _ := record.Get("description1")
		createdAt1, _ := record.Get("created_at1")
		
		id2, _ := record.Get("id2")
		name2, _ := record.Get("name2")
		description2, _ := record.Get("description2")
		createdAt2, _ := record.Get("created_at2")
		
		// Handle nil values safely for concept1
		name1Str := ""
		if name1 != nil {
			name1Str = name1.(string)
		}
		
		desc1Str := ""
		if description1 != nil {
			desc1Str = description1.(string)
		}
		
		var created1Time time.Time
		if createdAt1 != nil {
			created1Time = createdAt1.(time.Time)
		} else {
			created1Time = time.Now()
		}
		
		// Handle nil values safely for concept2
		name2Str := ""
		if name2 != nil {
			name2Str = name2.(string)
		}
		
		desc2Str := ""
		if description2 != nil {
			desc2Str = description2.(string)
		}
		
		var created2Time time.Time
		if createdAt2 != nil {
			created2Time = createdAt2.(time.Time)
		} else {
			created2Time = time.Now()
		}
		
		concept1 := models.Concept{
			ID:          fmt.Sprintf("%v", id1),
			Name:        name1Str,
			Description: desc1Str,
			CreatedAt:   created1Time,
		}
		
		concept2 := models.Concept{
			ID:          fmt.Sprintf("%v", id2),
			Name:        name2Str,
			Description: desc2Str,
			CreatedAt:   created2Time,
		}
		
		pairs = append(pairs, []models.Concept{concept1, concept2})
	}
	
	return pairs, nil
}

// GetGraphStats retrieves statistics about the graph
func (c *Client) GetGraphStats() (*models.GraphStats, error) {
	result, err := c.session.Run(
		`
		MATCH (c:Concept)
		WITH count(c) as node_count
		MATCH ()-[r:RELATED_TO]->()
		RETURN node_count, count(r) as relationship_count
		`,
		nil,
	)
	
	if err != nil {
		return nil, errors.NewDatabaseError("failed to get graph stats", err)
	}
	
	if result.Next() {
		record := result.Record()
		nodeCount, _ := record.Get("node_count")
		relationshipCount, _ := record.Get("relationship_count")
		
		return &models.GraphStats{
			NodeCount:         int(nodeCount.(int64)),
			RelationshipCount: int(relationshipCount.(int64)),
			LastUpdated:       time.Now(),
		}, nil
	}
	
	return &models.GraphStats{
		NodeCount:         0,
		RelationshipCount: 0,
		LastUpdated:       time.Now(),
	}, nil
}

// FixDuplicateNodes fixes duplicate nodes in the database
func (c *Client) FixDuplicateNodes() error {
	// Try to fix the specific duplicate nodes mentioned in the error message using direct Cypher
	fmt.Println("Attempting to fix duplicate 'Complexity Theory' nodes...")
	
	// First, find all nodes with the name 'Complexity Theory'
	result, err := c.session.Run(
		`
		MATCH (c:Concept)
		WHERE c.name = 'Complexity Theory'
		RETURN id(c) as id
		ORDER BY id(c)
		`,
		nil,
	)
	
	if err != nil {
		return errors.NewDatabaseError("failed to find duplicate concepts", err)
	}
	
	// Keep the first node, delete the rest
	var firstNodeId string
	var nodesToDelete []string
	
	for result.Next() {
		record := result.Record()
		id, _ := record.Get("id")
		idStr := fmt.Sprintf("%v", id)
		
		if firstNodeId == "" {
			firstNodeId = idStr
			fmt.Printf("Keeping node with ID: %s\n", idStr)
		} else {
			nodesToDelete = append(nodesToDelete, idStr)
			fmt.Printf("Will delete duplicate node with ID: %s\n", idStr)
		}
	}
	
	// If we found duplicates, delete them
	if len(nodesToDelete) > 0 {
		fmt.Printf("Found %d duplicate nodes to delete\n", len(nodesToDelete))
		
		for _, id := range nodesToDelete {
			// Now delete the duplicate node
			_, err = c.session.Run(
				`
				MATCH (c:Concept)
				WHERE id(c) = $id
				DETACH DELETE c
				`,
				map[string]interface{}{
					"id": id,
				},
			)
			
			if err != nil {
				fmt.Printf("Error deleting duplicate node: %s\n", err.Error())
			} else {
				fmt.Printf("Deleted duplicate node with ID: %s\n", id)
			}
		}
	} else {
		fmt.Println("No duplicate 'Complexity Theory' nodes found")
	}
	
	return nil
} 