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
	neo4jDriver, err = createNeo4jDriver(neo4jURI, neo4jUser, neo4jPassword)
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

	// Create command runner
	commandRunner := &utils.CommandRunner{}

	// Set up Gin router
	router := setupRouter(neo4jDriver, commandRunner)

	// Start server
	port := utils.GetEnv("PORT", "8081")
	logger.Printf("Server listening on port %s", port)
	if err := router.Run(":" + port); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}

// createNeo4jDriver creates a new Neo4j driver
func createNeo4jDriver(uri, username, password string) (neo4j.Driver, error) {
	return neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
}

// setupRouter sets up the Gin router with all routes
func setupRouter(driver neo4j.Driver, runner utils.CommandRunnerInterface) *gin.Engine {
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
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Graph data
		api.GET("/graph", handlers.GetGraphData(driver))

		// Builder controls
		api.POST("/builder/start", handlers.StartBuilder(runner))
		api.POST("/builder/stop", handlers.StopBuilder(runner))

		// Enricher controls
		api.POST("/enricher/start", handlers.StartEnricher(runner))
		api.POST("/enricher/stop", handlers.StopEnricher(runner))

		// Concept search
		api.GET("/concepts/search", handlers.SearchConcepts(driver))

		// Relationship creation
		api.POST("/relationships", handlers.CreateRelationship(driver))

		// Statistics
		api.GET("/statistics", handlers.GetStatistics(driver))
	}

	// Serve static files from the public directory
	router.Static("/css", "/app/public/css")
	router.Static("/js", "/app/public/js")
	router.Static("/static", "/app/public")
	
	// Handle all other routes by serving index.html
	router.NoRoute(func(c *gin.Context) {
		c.File("/app/public/index.html")
	})

	return router
} 