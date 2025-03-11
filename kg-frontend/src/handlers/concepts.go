package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"

	"kg-frontend/src/models"
)

// SearchConcepts searches for concepts in the graph that match the query
func SearchConcepts(neo4jDriver neo4j.Driver) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get query parameter
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Query parameter 'q' is required"})
			return
		}

		session := neo4jDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
		defer session.Close()

		// Query concepts that match the search term
		result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			result, err := tx.Run(
				"MATCH (n:Concept) WHERE n.name CONTAINS $query RETURN id(n) AS id, n.name AS name LIMIT 10",
				map[string]interface{}{"query": query},
			)
			if err != nil {
				return nil, err
			}

			var concepts []models.Node
			for result.Next() {
				record := result.Record()
				id, _ := record.Get("id")
				name, _ := record.Get("name")

				concepts = append(concepts, models.Node{
					ID:   strconv.FormatInt(id.(int64), 10),
					Name: name.(string),
				})
			}

			return concepts, nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to search concepts: " + err.Error()})
			return
		}

		// Return concepts
		c.JSON(http.StatusOK, result)
	}
}

// CreateRelationship creates a relationship between two concepts
func CreateRelationship(neo4jDriver neo4j.Driver) gin.HandlerFunc {
	return func(c *gin.Context) {
		var relationship models.RelationshipCreate
		if err := c.ShouldBindJSON(&relationship); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid parameters: " + err.Error()})
			return
		}

		session := neo4jDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
		defer session.Close()

		// Create relationship between concepts
		result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			// Convert string IDs to int64
			sourceID, err := strconv.ParseInt(relationship.Source, 10, 64)
			if err != nil {
				return nil, err
			}

			targetID, err := strconv.ParseInt(relationship.Target, 10, 64)
			if err != nil {
				return nil, err
			}

			// Create the relationship
			query := "MATCH (a:Concept), (b:Concept) WHERE id(a) = $source AND id(b) = $target CREATE (a)-[r:`" + relationship.Type + "`]->(b) RETURN id(a) AS source, id(b) AS target, type(r) AS type"
			result, err := tx.Run(
				query,
				map[string]interface{}{
					"source": sourceID,
					"target": targetID,
				},
			)
			if err != nil {
				return nil, err
			}

			if !result.Next() {
				return nil, nil
			}

			record := result.Record()
			source, _ := record.Get("source")
			target, _ := record.Get("target")
			relType, _ := record.Get("type")

			return models.Link{
				Source: strconv.FormatInt(source.(int64), 10),
				Target: strconv.FormatInt(target.(int64), 10),
				Type:   relType.(string),
			}, nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create relationship: " + err.Error()})
			return
		}

		if result == nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Concepts not found"})
			return
		}

		// Return the created relationship
		c.JSON(http.StatusOK, result)
	}
} 