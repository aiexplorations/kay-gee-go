package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Statistics represents the statistics of the knowledge graph
type Statistics struct {
	TotalConcepts      int `json:"total_concepts"`
	TotalRelationships int `json:"total_relationships"`
}

// handleStatistics handles requests to get statistics
func handleStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := getStatistics()
	if err != nil {
		log.Printf("Error getting statistics: %v", err)
		http.Error(w, "Failed to get statistics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// getStatistics returns the statistics of the knowledge graph
func getStatistics() (*Statistics, error) {
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
	}
	
	return stats, nil
} 