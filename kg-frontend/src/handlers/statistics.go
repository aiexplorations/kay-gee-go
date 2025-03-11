package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"

	"kg-frontend/src/models"
)

// GetStatistics retrieves statistics about the knowledge graph
func GetStatistics(neo4jDriver neo4j.Driver) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := neo4jDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
		defer session.Close()

		// Query concept count
		conceptCount, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			result, err := tx.Run(
				"MATCH (n:Concept) RETURN count(n) AS count",
				nil,
			)
			if err != nil {
				return nil, err
			}

			if !result.Next() {
				return 0, nil
			}

			record := result.Record()
			count, _ := record.Get("count")
			return count, nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch concept count: " + err.Error()})
			return
		}

		// Query relationship count
		relationshipCount, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			result, err := tx.Run(
				"MATCH ()-[r]->() RETURN count(r) AS count",
				nil,
			)
			if err != nil {
				return nil, err
			}

			if !result.Next() {
				return 0, nil
			}

			record := result.Record()
			count, _ := record.Get("count")
			return count, nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch relationship count: " + err.Error()})
			return
		}

		// Return statistics
		c.JSON(http.StatusOK, models.Statistics{
			ConceptCount:      int(conceptCount.(int64)),
			RelationshipCount: int(relationshipCount.(int64)),
		})
	}
} 