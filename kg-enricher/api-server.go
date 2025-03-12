package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

// StartRequest represents the request to start the enricher
type StartRequest struct {
	BatchSize       int `json:"batchSize"`
	Interval        int `json:"interval"`
	MaxRelationships int `json:"maxRelationships"`
	Concurrency     int `json:"concurrency"`
}

// Response represents the API response
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

var enricherProcess *exec.Cmd

// startHandler handles the request to start the enricher
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
	
	// Check if the enricher is already running
	if enricherProcess != nil && enricherProcess.Process != nil {
		// Check if the process is still running
		if err := enricherProcess.Process.Signal(syscall.Signal(0)); err == nil {
			sendErrorResponse(w, "Enricher is already running")
			return
		}
	}
	
	// Start the enricher process
	args := []string{
		"--batch-size", strconv.Itoa(req.BatchSize),
		"--interval", strconv.Itoa(req.Interval),
		"--max-relationships", strconv.Itoa(req.MaxRelationships),
		"--concurrency", strconv.Itoa(req.Concurrency),
	}
	
	log.Printf("Starting enricher with args: %v", args)
	enricherProcess = exec.Command("/app/enricher", args...)
	enricherProcess.Stdout = os.Stdout
	enricherProcess.Stderr = os.Stderr
	
	if err := enricherProcess.Start(); err != nil {
		log.Printf("Error starting enricher: %v", err)
		sendErrorResponse(w, fmt.Sprintf("Failed to start enricher: %v", err))
		return
	}
	
	// Send success response
	sendSuccessResponse(w, "Enricher started successfully")
	
	// Start a goroutine to wait for the process to complete
	go func() {
		if err := enricherProcess.Wait(); err != nil {
			log.Printf("Enricher process exited with error: %v", err)
		} else {
			log.Printf("Enricher process completed successfully")
		}
	}()
}

// stopHandler handles the request to stop the enricher
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
	
	// Check if the enricher is running
	if enricherProcess == nil || enricherProcess.Process == nil {
		sendErrorResponse(w, "Enricher is not running")
		return
	}
	
	// Stop the enricher process
	if err := enricherProcess.Process.Kill(); err != nil {
		log.Printf("Error stopping enricher: %v", err)
		sendErrorResponse(w, fmt.Sprintf("Failed to stop enricher: %v", err))
		return
	}
	
	// Send success response
	sendSuccessResponse(w, "Enricher stopped successfully")
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
	port := "5001"
	log.Printf("Starting API server on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start API server: %v", err)
	}
}

func main() {
	// Start the API server
	startAPIServer()
} 