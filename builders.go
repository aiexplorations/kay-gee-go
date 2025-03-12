package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Builder represents a builder container
type Builder struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Concept   string    `json:"concept"`
	Status    string    `json:"status"`
	Progress  int       `json:"progress"`
	StartTime time.Time `json:"start_time"`
}

// handleBuilders handles GET requests for builders
func handleBuilders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if a specific builder ID is requested
	builderID := r.URL.Query().Get("id")
	if builderID != "" {
		builder, err := getBuilderByID(builderID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(builder)
		return
	}

	// Otherwise return all active builders
	builders := getActiveBuilders()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(builders)
}

// getBuilderByID returns a specific builder by ID
func getBuilderByID(id string) (*Builder, error) {
	builders := getActiveBuilders()
	
	for _, builder := range builders {
		if builder.ID == id {
			return &builder, nil
		}
	}
	
	return nil, fmt.Errorf("builder with ID %s not found", id)
}

// getActiveBuilders returns a list of active builders
// Currently returns hardcoded data for testing
func getActiveBuilders() []Builder {
	// Hardcoded data for testing
	startTime, _ := time.Parse(time.RFC3339, "2025-03-12T12:22:30.032613131Z")
	
	return []Builder{
		{
			ID:        "1e2b456f2caa",
			Name:      "kaygeego-builder-Machine_Learning-1741782270",
			Concept:   "Machine Learning",
			Status:    "running",
			Progress:  50,
			StartTime: startTime,
		},
	}
}

// extractConceptFromName extracts the concept from the container name
func extractConceptFromName(name string) string {
	// Expected format: kaygeego-builder-Concept_Name-timestamp
	parts := strings.Split(name, "-")
	if len(parts) < 3 {
		return "Unknown"
	}
	
	conceptWithTimestamp := strings.Join(parts[2:], "-")
	conceptParts := strings.Split(conceptWithTimestamp, "-")
	if len(conceptParts) < 2 {
		return strings.ReplaceAll(conceptWithTimestamp, "_", " ")
	}
	
	concept := strings.Join(conceptParts[:len(conceptParts)-1], "-")
	return strings.ReplaceAll(concept, "_", " ")
}

// getBuilderProgress returns the progress of a builder
// Currently returns a hardcoded value for testing
func getBuilderProgress(containerID string) int {
	// Hardcoded progress for testing
	return 50
} 