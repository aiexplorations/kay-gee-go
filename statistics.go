package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Statistics represents the statistics of the knowledge graph
type Statistics struct {
	TotalConcepts      int              `json:"total_concepts"`
	TotalRelationships int              `json:"total_relationships"`
	LastUpdated        time.Time        `json:"last_updated"`
	MostConnected      []ConnectedNode  `json:"most_connected,omitempty"`
}

// ConnectedNode represents a concept node with its connection count
type ConnectedNode struct {
	Name         string `json:"name"`
	ConnectionCount int    `json:"connection_count"`
}

// handleStatistics handles GET requests for statistics
func handleStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if detailed statistics are requested
	detailed := r.URL.Query().Get("detailed") == "true"
	log.Printf("Detailed statistics requested: %v", detailed)
	
	stats := getStatistics(detailed)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// getStatistics returns the statistics of the knowledge graph
// Currently returns hardcoded data for testing
func getStatistics(detailed bool) Statistics {
	stats := Statistics{
		TotalConcepts:      100,
		TotalRelationships: 150,
		LastUpdated:        time.Now(),
	}
	
	if detailed {
		stats.MostConnected = getMostConnectedConcepts(5)
	}
	
	return stats
}

// getMostConnectedConcepts returns the most connected concepts
// Currently returns hardcoded data for testing
func getMostConnectedConcepts(limit int) []ConnectedNode {
	// Hardcoded data for testing
	return []ConnectedNode{
		{Name: "Artificial Intelligence", ConnectionCount: 25},
		{Name: "Machine Learning", ConnectionCount: 18},
		{Name: "Neural Networks", ConnectionCount: 15},
		{Name: "Deep Learning", ConnectionCount: 12},
		{Name: "Natural Language Processing", ConnectionCount: 10},
	}
}

// parseCount parses the count from the output of the Cypher shell
func parseCount(output string) int {
	// Hardcoded count for testing
	return 100
} 