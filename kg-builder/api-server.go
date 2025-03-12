package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

// StartRequest represents the request to start the builder
type StartRequest struct {
	SeedConcept        string `json:"seedConcept"`
	MaxNodes           int    `json:"maxNodes"`
	Timeout            int    `json:"timeout"`
	RandomRelationships int    `json:"randomRelationships"`
	Concurrency        int    `json:"concurrency"`
}

// Response represents the API response
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

var builderProcess *exec.Cmd

// startHandler handles the request to start the builder
func startHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	// Handle preflight request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	
	// Check if the request method is POST
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Parse the request body
	var req StartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error parsing request body: %v", err)
		sendErrorResponse(w, "Invalid request body")
		return
	}
	
	// Check if the builder is already running
	if builderProcess != nil && builderProcess.Process != nil {
		// Check if the process is still running
		if err := builderProcess.Process.Signal(os.Signal(0)); err == nil {
			sendErrorResponse(w, "Builder is already running")
			return
		}
	}
	
	// Start the builder process
	args := []string{
		"--seed", req.SeedConcept,
		"--max-nodes", strconv.Itoa(req.MaxNodes),
		"--timeout", strconv.Itoa(req.Timeout),
		"--random-relationships", strconv.Itoa(req.RandomRelationships),
		"--concurrency", strconv.Itoa(req.Concurrency),
		"--use-low-connectivity",
	}
	
	log.Printf("Starting builder with args: %v", args)
	builderProcess = exec.Command("/kg-builder", args...)
	builderProcess.Stdout = os.Stdout
	builderProcess.Stderr = os.Stderr
	
	if err := builderProcess.Start(); err != nil {
		log.Printf("Error starting builder: %v", err)
		sendErrorResponse(w, fmt.Sprintf("Failed to start builder: %v", err))
		return
	}
	
	// Send success response
	sendSuccessResponse(w, "Builder started successfully")
	
	// Start a goroutine to wait for the process to complete
	go func() {
		if err := builderProcess.Wait(); err != nil {
			log.Printf("Builder process exited with error: %v", err)
		} else {
			log.Printf("Builder process completed successfully")
		}
	}()
}

// stopHandler handles the request to stop the builder
func stopHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	// Handle preflight request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	
	// Check if the request method is POST
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Check if the builder is running
	if builderProcess == nil || builderProcess.Process == nil {
		sendErrorResponse(w, "Builder is not running")
		return
	}
	
	// Stop the builder process
	if err := builderProcess.Process.Kill(); err != nil {
		log.Printf("Error stopping builder: %v", err)
		sendErrorResponse(w, fmt.Sprintf("Failed to stop builder: %v", err))
		return
	}
	
	// Send success response
	sendSuccessResponse(w, "Builder stopped successfully")
}

// sendSuccessResponse sends a success response
func sendSuccessResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: message,
	})
}

// sendErrorResponse sends an error response
func sendErrorResponse(w http.ResponseWriter, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(Response{
		Status: "error",
		Error:  errorMessage,
	})
}

// startAPIServer starts the API server
func startAPIServer() {
	// Register handlers
	http.HandleFunc("/start", startHandler)
	http.HandleFunc("/stop", stopHandler)
	
	// Start the server
	port := "5000"
	log.Printf("Starting API server on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start API server: %v", err)
	}
}

func main() {
	// Start the API server
	startAPIServer()
} 