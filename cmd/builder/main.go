package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kay-gee-go/internal/builder"
	"github.com/kay-gee-go/internal/common/config"
)

func main() {
	// Parse command-line flags
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadBuilderConfig(*configPath)
	if err != nil {
		fmt.Printf("Error loading configuration: %s\n", err.Error())
		os.Exit(1)
	}

	// Create builder
	b, err := builder.NewBuilder(cfg)
	if err != nil {
		fmt.Printf("Error creating builder: %s\n", err.Error())
		os.Exit(1)
	}
	defer b.Close()

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("Received signal, shutting down...")
		cancel()
	}()

	// Build the knowledge graph
	stats, err := b.Build(ctx)
	if err != nil {
		fmt.Printf("Error building knowledge graph: %s\n", err.Error())
		os.Exit(1)
	}

	// Print stats
	fmt.Println("Knowledge graph building completed successfully.")
	fmt.Printf("Nodes created: %d\n", stats.NodesCreated)
	fmt.Printf("Relationships created: %d\n", stats.RelationshipsCreated)
	fmt.Printf("Duration: %s\n", stats.Duration)
} 