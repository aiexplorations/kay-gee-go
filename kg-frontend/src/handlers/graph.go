package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"

	"kg-frontend/src/models"
)

// GetGraphData retrieves the current graph data from Neo4j
func GetGraphData(neo4jDriver neo4j.Driver) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := neo4jDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
		defer session.Close()

		// Query nodes
		nodes, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			result, err := tx.Run(
				"MATCH (n:Concept) RETURN id(n) AS id, n.name AS name",
				nil,
			)
			if err != nil {
				return nil, err
			}

			var nodes []models.Node
			for result.Next() {
				record := result.Record()
				id, _ := record.Get("id")
				name, _ := record.Get("name")

				nodes = append(nodes, models.Node{
					ID:   strconv.FormatInt(id.(int64), 10),
					Name: name.(string),
					Size: 5, // Default size
				})
			}

			return nodes, nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch nodes: " + err.Error()})
			return
		}

		// Query links
		links, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			result, err := tx.Run(
				"MATCH (a:Concept)-[r]->(b:Concept) RETURN id(a) AS source, id(b) AS target, type(r) AS type",
				nil,
			)
			if err != nil {
				return nil, err
			}

			var links []models.Link
			for result.Next() {
				record := result.Record()
				source, _ := record.Get("source")
				target, _ := record.Get("target")
				relType, _ := record.Get("type")

				links = append(links, models.Link{
					Source: strconv.FormatInt(source.(int64), 10),
					Target: strconv.FormatInt(target.(int64), 10),
					Type:   relType.(string),
				})
			}

			return links, nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch links: " + err.Error()})
			return
		}

		// Return graph data
		c.JSON(http.StatusOK, models.GraphData{
			Nodes: nodes.([]models.Node),
			Links: links.([]models.Link),
		})
	}
} 