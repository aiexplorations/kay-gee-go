package main

import (
	"encoding/json"
	"fmt"
	"log"
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

// handleBuilders handles requests to list active builders
func handleBuilders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if a specific builder ID is requested
	builderID := r.URL.Query().Get("id")
	if builderID != "" {
		log.Printf("Builder ID requested: %s", builderID)
		builder, err := getBuilderByID(builderID)
		if err != nil {
			log.Printf("Error getting builder by ID: %v", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(builder)
		return
	}

	// Otherwise return all active builders
	builders, err := getActiveBuilders()
	if err != nil {
		log.Printf("Error getting active builders: %v", err)
		http.Error(w, "Failed to get active builders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(builders)
}

// getBuilderByID returns a specific builder by ID
func getBuilderByID(id string) (*Builder, error) {
	builders, err := getActiveBuilders()
	if err != nil {
		return nil, err
	}
	
	for _, builder := range builders {
		if builder.ID == id {
			return &builder, nil
		}
	}
	
	return nil, fmt.Errorf("builder with ID %s not found", id)
}

// getActiveBuilders returns a list of active builder containers
func getActiveBuilders() ([]Builder, error) {
	// Since we can't execute docker commands from inside the container,
	// we'll return a hardcoded list of active builders for now
	// This can be enhanced later with a more sophisticated approach
	
	// Check if we have any active builders from our own tracking
	builders := []Builder{
		{
			ID:        "1e2b456f2caa",
			Name:      "kaygeego-builder-Machine_Learning-1741782270",
			Concept:   "Machine Learning",
			Status:    "running",
			Progress:  50,
			StartTime: time.Now().Add(-5 * time.Minute),
		},
	}
	
	// In a real implementation, we would need to:
	// 1. Use a shared database or API to track builder containers
	// 2. Have the builder containers register themselves when they start
	// 3. Update their progress periodically
	// 4. Remove them when they complete or fail
	
	return builders, nil
}

// extractConceptFromName extracts the concept from the container name
func extractConceptFromName(name string) string {
	// Expected format: kaygeego-builder-Concept_Name-Timestamp
	parts := strings.Split(name, "-")
	if len(parts) < 3 {
		return ""
	}
	
	// Get the concept part (may contain underscores)
	conceptWithTimestamp := strings.Join(parts[2:], "-")
	
	// Split by the last hyphen to separate concept from timestamp
	lastHyphen := strings.LastIndex(conceptWithTimestamp, "-")
	if lastHyphen == -1 {
		return conceptWithTimestamp
	}
	
	concept := conceptWithTimestamp[:lastHyphen]
	// Replace underscores with spaces
	return strings.ReplaceAll(concept, "_", " ")
}

// getBuilderProgress gets the progress of a builder from its logs
func getBuilderProgress(containerID string) int {
	// In a real implementation, we would parse the logs to get the progress
	// For now, return a placeholder value
	return 50
} 