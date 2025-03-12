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
	LastUpdated        time.Time        `json:"last_updated,omitempty"`
	MostConnected      []ConnectedNode  `json:"most_connected,omitempty"`
}

// ConnectedNode represents a concept node with its connection count
type ConnectedNode struct {
	Name            string `json:"name"`
	ConnectionCount int    `json:"connection_count"`
}

// handleStatistics handles requests to get statistics
func handleStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if detailed statistics are requested
	detailed := r.URL.Query().Get("detailed") == "true"
	log.Printf("Detailed statistics requested: %v", detailed)

	stats, err := getStatistics(detailed)
	if err != nil {
		log.Printf("Error getting statistics: %v", err)
		http.Error(w, "Failed to get statistics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// getStatistics returns the statistics of the knowledge graph
func getStatistics(detailed bool) (*Statistics, error) {
	// Since we can't execute docker commands from inside the container,
	// we'll return hardcoded statistics for now
	// This can be enhanced later with a more sophisticated approach
	
	// In a real implementation, we would:
	// 1. Use a Neo4j driver to connect directly to the database
	// 2. Execute Cypher queries to get the actual counts
	// 3. Cache the results to avoid frequent database queries
	
	// For now, return placeholder values
	stats := &Statistics{
		TotalConcepts:      100,
		TotalRelationships: 150,
		LastUpdated:        time.Now(),
	}
	
	// Add most connected concepts if detailed statistics are requested
	if detailed {
		stats.MostConnected = getMostConnectedConcepts(5)
	}
	
	return stats, nil
}

// getMostConnectedConcepts returns the most connected concepts
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