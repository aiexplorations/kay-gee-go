package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

// ConceptRequest represents a request to add a new concept
type ConceptRequest struct {
	Concept             string `json:"concept"`
	MaxNodes            int    `json:"max_nodes,omitempty"`
	Timeout             int    `json:"timeout,omitempty"`
	RandomRelationships int    `json:"random_relationships,omitempty"`
	Concurrency         int    `json:"concurrency,omitempty"`
}

// ConceptResponse represents the response to a concept request
type ConceptResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	ContainerID  string `json:"container_id,omitempty"`
	ContainerName string `json:"container_name,omitempty"`
}

func main() {
	// Set up HTTP server
	http.HandleFunc("/api/concepts", handleConcepts)
	http.HandleFunc("/api/health", handleHealth)
	http.HandleFunc("/api/builders", handleBuilders)
	http.HandleFunc("/api/statistics", handleStatistics)

	// Get port from environment variable or use default
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "5000"
	}

	// Start server
	log.Printf("Starting API server on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// handleHealth handles health check requests
func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"status": true})
}

// handleConcepts handles requests to add new concepts
func handleConcepts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req ConceptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Concept == "" {
		http.Error(w, "Concept is required", http.StatusBadRequest)
		return
	}

	// Set default values if not provided
	if req.MaxNodes <= 0 {
		req.MaxNodes = 50
	}
	if req.Timeout <= 0 {
		req.Timeout = 15
	}
	if req.RandomRelationships <= 0 {
		req.RandomRelationships = 10
	}
	if req.Concurrency <= 0 {
		req.Concurrency = 3
	}

	// Get the project root directory
	projectRoot, err := getProjectRoot()
	if err != nil {
		log.Printf("Error getting project root: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Build the command to launch a new builder container
	scriptPath := filepath.Join(projectRoot, "add-concept.sh")
	cmd := exec.Command(scriptPath, 
		req.Concept,
		"--max-nodes="+strconv.Itoa(req.MaxNodes),
		"--timeout="+strconv.Itoa(req.Timeout),
		"--random-relationships="+strconv.Itoa(req.RandomRelationships),
		"--concurrency="+strconv.Itoa(req.Concurrency),
	)

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error launching builder container: %v\nOutput: %s", err, output)
		http.Error(w, "Failed to launch builder container", http.StatusInternalServerError)
		return
	}

	// Extract container ID from output (this is a simplification, you might need to parse the output more carefully)
	containerID := "unknown"
	containerName := "unknown"
	
	// Prepare response
	resp := ConceptResponse{
		Success:      true,
		Message:      fmt.Sprintf("Successfully launched builder container for concept: %s", req.Concept),
		ContainerID:  containerID,
		ContainerName: containerName,
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// getProjectRoot returns the root directory of the project
func getProjectRoot() (string, error) {
	// In a real implementation, you might want to use a more robust method
	// For now, we'll assume the API server is running in the kg-builder/cmd/api-server directory
	// and the project root is three levels up
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	
	// Navigate up to the project root
	return filepath.Join(dir, "..", "..", ".."), nil
} 