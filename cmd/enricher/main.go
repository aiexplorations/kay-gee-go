package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kay-gee-go/internal/common/config"
	"github.com/kay-gee-go/internal/enricher"
)

func main() {
	// Parse command-line flags
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	runOnce := flag.Bool("run-once", false, "Run once and exit")
	count := flag.Int("count", 10, "Number of relationships to mine when running once")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadEnricherConfig(*configPath)
	if err != nil {
		fmt.Printf("Error loading configuration: %s\n", err.Error())
		os.Exit(1)
	}

	// Create enricher
	e, err := enricher.NewEnricher(cfg)
	if err != nil {
		fmt.Printf("Error creating enricher: %s\n", err.Error())
		os.Exit(1)
	}
	defer e.Close()

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

	// Run the enricher
	if *runOnce {
		// Run once and exit
		stats, err := e.EnrichOnce(ctx, *count)
		if err != nil {
			fmt.Printf("Error enriching knowledge graph: %s\n", err.Error())
			os.Exit(1)
		}

		// Print stats
		fmt.Println("Knowledge graph enrichment completed successfully.")
		fmt.Printf("Pairs processed: %d\n", stats.PairsProcessed)
		fmt.Printf("Relationships created: %d\n", stats.RelationshipsCreated)
		fmt.Printf("Duration: %s\n", stats.Duration)
	} else {
		// Run continuously
		err := e.EnrichContinuously(ctx)
		if err != nil && err != context.Canceled {
			fmt.Printf("Error enriching knowledge graph: %s\n", err.Error())
			os.Exit(1)
		}

		// Print stats
		stats := e.GetStats()
		fmt.Println("Knowledge graph enrichment stopped.")
		fmt.Printf("Pairs processed: %d\n", stats.PairsProcessed)
		fmt.Printf("Relationships created: %d\n", stats.RelationshipsCreated)
		if stats.EndTime.IsZero() {
			stats.EndTime = time.Now()
			stats.Duration = stats.EndTime.Sub(stats.StartTime).String()
		}
		fmt.Printf("Duration: %s\n", stats.Duration)
	}
} 