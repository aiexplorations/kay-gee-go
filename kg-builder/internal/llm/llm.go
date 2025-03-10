package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"kg-builder/internal/config"
	apperrors "kg-builder/internal/errors"
	"kg-builder/internal/models"
)

const (
	DefaultMaxRetries    = 3
	DefaultRetryInterval = 2 * time.Second
	DefaultMaxBackoff    = 15 * time.Second
)

var (
	// Cache for LLM responses
	conceptCache     = make(map[string][]models.Concept)
	relationshipCache = make(map[string]*models.Concept)
	cacheMutex       sync.RWMutex
	
	// Configuration
	cfg *config.LLMConfig
)

// Initialize the LLM service with configuration
func Initialize(config *config.LLMConfig) error {
	cfg = config
	
	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(cfg.CacheDir, 0755); err != nil {
		log.Printf("Failed to create cache directory: %v", err)
		return err
	}
	
	// Load cached data if available
	loadConceptCache()
	loadRelationshipCache()
	
	return nil
}

// GetRelatedConcepts sends a request to the LLM service to get related concepts for a given concept.
func GetRelatedConcepts(concept string) ([]models.Concept, error) {
	if cfg == nil {
		return nil, apperrors.NewLLMError(apperrors.ErrInvalidInput, "LLM service not initialized")
	}

	// Check cache first
	cacheMutex.RLock()
	if concepts, ok := conceptCache[concept]; ok {
		cacheMutex.RUnlock()
		log.Printf("Cache hit for concept: %s", concept)
		return concepts, nil
	}
	cacheMutex.RUnlock()

	prompt := fmt.Sprintf(`You are an expert ontologist with an understanding of concepts and the relationships between them. You respond only in JSON. 
	Given the concept '%s', provide 5 related concepts. 
	For each, specify the relationship type. 
	Return ONLY a JSON array with 'name', 'relation', and 'relatedTo' keys. 
	Do not include any explanations, markdown formatting, or additional text. 
	The response should be valid JSON that can be directly parsed. Example format:
    [
        {
            "name": "Related Concept 1",
            "relation": "RelationType",
            "relatedTo": "%s"
        },
        ...
    ]
	Do not return any explanations, markdown formatting, or additional text.
	`, concept, concept)

	var concepts []models.Concept
	var err error

	err = apperrors.RetryWithBackoff(DefaultMaxRetries, DefaultRetryInterval, DefaultMaxBackoff, func() error {
		// Marshal the request body
		requestBody, err := json.Marshal(map[string]string{
			"model":  cfg.Model,
			"prompt": prompt,
		})

		// Check if the request body was marshalled successfully
		if err != nil {
			return apperrors.NewLLMError(err, "failed to marshal request")
		}

		// Send the request to the LLM service
		resp, err := http.Post(cfg.URL, "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			return apperrors.NewLLMError(err, "failed to make request").WithRetryable(true)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return apperrors.NewLLMError(
				fmt.Errorf("unexpected status code: %d", resp.StatusCode),
				"LLM service returned non-OK status",
			).WithRetryable(resp.StatusCode >= 500) // Only retry on server errors
		}

		// Read the response from the LLM service
		var fullResponse strings.Builder
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			var streamResponse struct {
				Response string `json:"response"`
			}
			if err := json.Unmarshal([]byte(line), &streamResponse); err == nil {
				fullResponse.WriteString(streamResponse.Response)
			}
		}

		// Check if there was an error reading the response
		if err := scanner.Err(); err != nil {
			return apperrors.NewLLMError(err, "error reading response").WithRetryable(true)
		}

		// Unmarshal the response into a slice of Concept structs
		if err := json.Unmarshal([]byte(fullResponse.String()), &concepts); err != nil {
			log.Printf("Raw LLM response: %s", fullResponse.String())
			return apperrors.NewLLMError(err, "failed to unmarshal concepts").WithRetryable(false)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Cache the result
	cacheMutex.Lock()
	conceptCache[concept] = concepts
	cacheMutex.Unlock()
	
	// Save to disk
	saveConceptCache(concept, concepts)

	return concepts, nil
}

// MineRelationship sends a request to the LLM service to determine if there is a relationship between two concepts.
func MineRelationship(concept1, concept2 string) (*models.Concept, error) {
	if cfg == nil {
		return nil, apperrors.NewLLMError(apperrors.ErrInvalidInput, "LLM service not initialized")
	}

	// Create a cache key for the relationship
	cacheKey := fmt.Sprintf("%s|%s", concept1, concept2)
	
	// Check cache first
	cacheMutex.RLock()
	if concept, ok := relationshipCache[cacheKey]; ok {
		cacheMutex.RUnlock()
		log.Printf("Cache hit for relationship: %s -> %s", concept1, concept2)
		return concept, nil
	}
	cacheMutex.RUnlock()

	prompt := fmt.Sprintf(`You are an expert ontologist and respond only in JSON. 
	Determine if there's a relationship between the concepts '%s' and '%s'. If there is, provide the relationship type. 
	If not, respond with "No relationship". 
	Return the response as a JSON object with 'name', 'relation', and 'relatedTo' keys. The response should be valid JSON that can be directly parsed. 
	Example format:
    {
        "name": "%s",
        "relation": "RelationType",
        "relatedTo": "%s"
    }
    Or if there's no relationship:
    {
        "name": "",
        "relation": "",
        "relatedTo": ""
    }
	Do not return any explanations, markdown formatting, or additional text.`, concept1, concept2, concept2, concept1)

	var concept *models.Concept
	var err error

	err = apperrors.RetryWithBackoff(DefaultMaxRetries, DefaultRetryInterval, DefaultMaxBackoff, func() error {
		requestBody, err := json.Marshal(map[string]string{
			"model":  cfg.Model,
			"prompt": prompt,
		})
		if err != nil {
			return apperrors.NewLLMError(err, "failed to marshal request")
		}

		// Send the request to the LLM service
		resp, err := http.Post(cfg.URL, "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			return apperrors.NewLLMError(err, "failed to make request").WithRetryable(true)
		}
		defer resp.Body.Close()

		// Check if the response status code is OK
		if resp.StatusCode != http.StatusOK {
			return apperrors.NewLLMError(
				fmt.Errorf("unexpected status code: %d", resp.StatusCode),
				"LLM service returned non-OK status",
			).WithRetryable(resp.StatusCode >= 500) // Only retry on server errors
		}

		// Read the response from the LLM service
		var fullResponse strings.Builder
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			var streamResponse struct {
				Response string `json:"response"`
			}
			if err := json.Unmarshal([]byte(line), &streamResponse); err == nil {
				fullResponse.WriteString(streamResponse.Response)
			}
		}

		// Check if there was an error reading the response
		if err := scanner.Err(); err != nil {
			return apperrors.NewLLMError(err, "error reading response").WithRetryable(true)
		}

		// Unmarshal the response into a Concept struct
		var conceptData models.Concept
		if err := json.Unmarshal([]byte(fullResponse.String()), &conceptData); err != nil {
			log.Printf("Raw LLM response: %s", fullResponse.String())
			return apperrors.NewLLMError(err, "failed to unmarshal concept").WithRetryable(false)
		}

		// Check if the relationship is empty
		if conceptData.Relation == "" {
			concept = nil
		} else {
			concept = &conceptData
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Cache the result
	cacheMutex.Lock()
	relationshipCache[cacheKey] = concept
	cacheMutex.Unlock()
	
	// Save to disk
	saveRelationshipCache(cacheKey, concept)

	return concept, nil
}

// Helper functions for configuration

// GetLLMURL returns the URL for the LLM service (exported for testing)
func GetLLMURL() string {
	if cfg != nil {
		return cfg.URL
	}
	return ""
}

// GetLLMModel returns the model for the LLM service (exported for testing)
func GetLLMModel() string {
	if cfg != nil {
		return cfg.Model
	}
	return ""
}

// Cache management functions

func saveConceptCache(concept string, concepts []models.Concept) {
	if cfg == nil {
		log.Printf("Cannot save concept cache: LLM service not initialized")
		return
	}
	
	filename := filepath.Join(cfg.CacheDir, fmt.Sprintf("concept_%s.json", sanitizeFilename(concept)))
	data, err := json.Marshal(concepts)
	if err != nil {
		log.Printf("Failed to marshal concept cache: %v", err)
		return
	}
	
	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		log.Printf("Failed to write concept cache to file: %v", err)
	}
}

func saveRelationshipCache(key string, concept *models.Concept) {
	if cfg == nil {
		log.Printf("Cannot save relationship cache: LLM service not initialized")
		return
	}
	
	filename := filepath.Join(cfg.CacheDir, fmt.Sprintf("rel_%s.json", sanitizeFilename(key)))
	data, err := json.Marshal(concept)
	if err != nil {
		log.Printf("Failed to marshal relationship cache: %v", err)
		return
	}
	
	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		log.Printf("Failed to write relationship cache to file: %v", err)
	}
}

func loadConceptCache() {
	if cfg == nil {
		log.Printf("Cannot load concept cache: LLM service not initialized")
		return
	}
	
	files, err := filepath.Glob(filepath.Join(cfg.CacheDir, "concept_*.json"))
	if err != nil {
		log.Printf("Failed to list concept cache files: %v", err)
		return
	}
	
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Printf("Failed to read concept cache file %s: %v", file, err)
			continue
		}
		
		var concepts []models.Concept
		if err := json.Unmarshal(data, &concepts); err != nil {
			log.Printf("Failed to unmarshal concept cache from %s: %v", file, err)
			continue
		}
		
		// Extract concept name from filename
		baseName := filepath.Base(file)
		conceptName := strings.TrimPrefix(baseName, "concept_")
		conceptName = strings.TrimSuffix(conceptName, ".json")
		conceptName = unsanitizeFilename(conceptName)
		
		cacheMutex.Lock()
		conceptCache[conceptName] = concepts
		cacheMutex.Unlock()
		
		log.Printf("Loaded %d concepts for '%s' from cache", len(concepts), conceptName)
	}
}

func loadRelationshipCache() {
	if cfg == nil {
		log.Printf("Cannot load relationship cache: LLM service not initialized")
		return
	}
	
	files, err := filepath.Glob(filepath.Join(cfg.CacheDir, "rel_*.json"))
	if err != nil {
		log.Printf("Failed to list relationship cache files: %v", err)
		return
	}
	
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Printf("Failed to read relationship cache file %s: %v", file, err)
			continue
		}
		
		var concept models.Concept
		if err := json.Unmarshal(data, &concept); err != nil {
			log.Printf("Failed to unmarshal relationship cache from %s: %v", file, err)
			continue
		}
		
		// Extract relationship key from filename
		baseName := filepath.Base(file)
		relKey := strings.TrimPrefix(baseName, "rel_")
		relKey = strings.TrimSuffix(relKey, ".json")
		relKey = unsanitizeFilename(relKey)
		
		cacheMutex.Lock()
		relationshipCache[relKey] = &concept
		cacheMutex.Unlock()
		
		log.Printf("Loaded relationship '%s' from cache", relKey)
	}
}

// sanitizeFilename replaces invalid characters in a filename with underscores
func sanitizeFilename(s string) string {
	// Replace invalid characters with underscores
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", "\n", "\r", "\t"}
	result := s
	
	for _, char := range invalid {
		result = strings.ReplaceAll(result, char, "_")
	}
	
	// Replace spaces with underscores
	result = strings.ReplaceAll(result, " ", "_")
	
	// Replace control characters and extended ASCII
	var sb strings.Builder
	for _, r := range result {
		if r < 32 || r > 126 {
			sb.WriteRune('_')
		} else {
			sb.WriteRune(r)
		}
	}
	
	return sb.String()
}

// SanitizeFilename is an exported version of sanitizeFilename for testing
func SanitizeFilename(s string) string {
	return sanitizeFilename(s)
}

// unsanitizeFilename reverses the sanitization process
func unsanitizeFilename(s string) string {
	// For now, we just return the string as is
	// In a real implementation, we might want to reverse some of the sanitization
	return s
}
