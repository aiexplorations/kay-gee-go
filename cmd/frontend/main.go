package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kay-gee-go/internal/common/config"
	"github.com/kay-gee-go/internal/frontend"
)

func main() {
	// Parse command-line flags
	port := flag.Int("port", 8080, "Port to listen on")
	neo4jURI := flag.String("neo4j-uri", "bolt://neo4j:7687", "Neo4j URI")
	neo4jUser := flag.String("neo4j-user", "neo4j", "Neo4j username")
	neo4jPassword := flag.String("neo4j-password", "password", "Neo4j password")
	flag.Parse()

	// Create Neo4j config
	neo4jConfig := &config.Neo4jConfig{
		URI:               *neo4jURI,
		User:              *neo4jUser,
		Password:          *neo4jPassword,
		MaxRetries:        5,
		RetryIntervalSecs: 5,
	}

	// Create frontend
	f, err := frontend.NewFrontend(neo4jConfig)
	if err != nil {
		fmt.Printf("Error creating frontend: %s\n", err.Error())
		os.Exit(1)
	}
	defer f.Close()

	// Handle signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("Received signal, shutting down...")
		os.Exit(0)
	}()

	// Start the frontend
	fmt.Printf("Starting frontend on port %d...\n", *port)
	if err := f.Start(*port); err != nil {
		fmt.Printf("Error starting frontend: %s\n", err.Error())
		os.Exit(1)
	}
} 