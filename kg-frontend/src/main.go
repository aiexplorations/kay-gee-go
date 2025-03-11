package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"

	"kg-frontend/src/handlers"
	"kg-frontend/src/utils"
)

// Global variables
var (
	neo4jDriver neo4j.Driver
	logger      *log.Logger
)

func main() {
	// Initialize logger
	logger = log.New(os.Stdout, "[KG-FRONTEND] ", log.LstdFlags)
	logger.Println("Starting Knowledge Graph Visualizer API")

	// Connect to Neo4j
	neo4jURI := utils.GetEnv("NEO4J_URI", "bolt://neo4j:7687")
	neo4jUser := utils.GetEnv("NEO4J_USER", "neo4j")
	neo4jPassword := utils.GetEnv("NEO4J_PASSWORD", "password")

	var err error
	neo4jDriver, err = neo4j.NewDriver(neo4jURI, neo4j.BasicAuth(neo4jUser, neo4jPassword, ""))
	if err != nil {
		logger.Fatalf("Failed to connect to Neo4j: %v", err)
	}
	defer neo4jDriver.Close()

	// Test Neo4j connection
	err = neo4jDriver.VerifyConnectivity()
	if err != nil {
		logger.Fatalf("Failed to verify Neo4j connectivity: %v", err)
	}
	logger.Println("Connected to Neo4j successfully")

	// Set up Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	// API routes
	api := router.Group("/api")
	{
		// Graph data
		api.GET("/graph", handlers.GetGraphData(neo4jDriver))

		// Builder controls
		api.POST("/builder/start", handlers.StartBuilder())
		api.POST("/builder/stop", handlers.StopBuilder())

		// Enricher controls
		api.POST("/enricher/start", handlers.StartEnricher())
		api.POST("/enricher/stop", handlers.StopEnricher())

		// Concept search
		api.GET("/concepts/search", handlers.SearchConcepts(neo4jDriver))

		// Relationship creation
		api.POST("/relationships", handlers.CreateRelationship(neo4jDriver))

		// Statistics
		api.GET("/statistics", handlers.GetStatistics(neo4jDriver))
	}

	// Serve static files
	router.Static("/", "./public")
	
	// Handle all other routes by serving index.html
	router.NoRoute(func(c *gin.Context) {
		c.File("./public/index.html")
	})

	// Start server
	port := utils.GetEnv("PORT", "8080")
	logger.Printf("Server listening on port %s", port)
	if err := router.Run(":" + port); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
} 