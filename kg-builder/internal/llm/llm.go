package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
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

	// Create a retry function
	var concepts []models.Concept
	var err error
	err = apperrors.RetryWithBackoff(DefaultMaxRetries, DefaultRetryInterval, DefaultMaxBackoff, func() error {
		// Create the request body
		prompt := fmt.Sprintf(`Generate 5 concepts related to "%s". 
For each concept, provide a name and a specific relation to "%s".
Return ONLY a JSON array of objects with the following structure:
[
  {
    "name": "Related Concept Name",
    "relation": "specific relation to %s",
    "relatedTo": "%s"
  }
]
Do not include any explanations, markdown formatting, or additional text. Return only the JSON array.`, 
			concept, concept, concept, concept)
		
		requestBody, err := json.Marshal(map[string]interface{}{
			"model":  cfg.Model,
			"prompt": prompt,
			"stream": false,
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
		var responseData struct {
			Response string `json:"response"`
		}
		
		if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
			return apperrors.NewLLMError(err, "error decoding response").WithRetryable(true)
		}
		
		// Log the raw response for debugging
		log.Printf("Raw LLM response: %s", responseData.Response)
		
		// Unmarshal the response into a slice of Concept structs
		if err := json.Unmarshal([]byte(responseData.Response), &concepts); err != nil {
			// If direct unmarshaling fails, try to extract JSON from the response
			jsonStr := extractJSONFromResponse(responseData.Response)
			log.Printf("Extracted JSON: %s", jsonStr)
			
			if err := json.Unmarshal([]byte(jsonStr), &concepts); err != nil {
				return apperrors.NewLLMError(err, "failed to unmarshal concepts").WithRetryable(false)
			}
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

// MineRelationship sends a request to the LLM service to mine a relationship between two concepts.
func MineRelationship(conceptA, conceptB string) (*models.Concept, error) {
	if cfg == nil {
		return nil, apperrors.NewLLMError(apperrors.ErrInvalidInput, "LLM service not initialized")
	}

	// Create a cache key for the relationship
	cacheKey := conceptA + "_" + conceptB

	// Check cache first
	cacheMutex.RLock()
	if relationship, ok := relationshipCache[cacheKey]; ok {
		cacheMutex.RUnlock()
		log.Printf("Cache hit for relationship: %s", cacheKey)
		return relationship, nil
	}
	cacheMutex.RUnlock()

	// Create a retry function
	var relationship *models.Concept
	var err error
	err = apperrors.RetryWithBackoff(DefaultMaxRetries, DefaultRetryInterval, DefaultMaxBackoff, func() error {
		// Create the request body
		prompt := fmt.Sprintf(`Determine a specific relationship between the concepts "%s" and "%s".
If there is a meaningful relationship, return a JSON object with the following structure:
{
  "name": "%s",
  "relation": "specific relationship type",
  "relatedTo": "%s"
}
If there is no meaningful relationship, return null.
Do not include any explanations, markdown formatting, or additional text. Return only the JSON object or null.`,
			conceptA, conceptB, conceptA, conceptB)
		
		requestBody, err := json.Marshal(map[string]interface{}{
			"model":  cfg.Model,
			"prompt": prompt,
			"stream": false,
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
		var responseData struct {
			Response string `json:"response"`
		}
		
		if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
			return apperrors.NewLLMError(err, "error decoding response").WithRetryable(true)
		}
		
		// Log the raw response for debugging
		log.Printf("Raw LLM response: %s", responseData.Response)
		
		// Unmarshal the response into a Concept struct
		var conceptData models.Concept
		if err := json.Unmarshal([]byte(responseData.Response), &conceptData); err != nil {
			// If direct unmarshaling fails, try to extract JSON from the response
			jsonStr := extractJSONFromResponse(responseData.Response)
			log.Printf("Extracted JSON: %s", jsonStr)
			
			if err := json.Unmarshal([]byte(jsonStr), &conceptData); err != nil {
				return apperrors.NewLLMError(err, "failed to unmarshal concept").WithRetryable(false)
			}
		}

		// Check if the relationship is empty
		if conceptData.Relation == "" {
			relationship = nil
		} else {
			relationship = &conceptData
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Cache the result
	cacheMutex.Lock()
	relationshipCache[cacheKey] = relationship
	cacheMutex.Unlock()
	
	// Save to disk
	saveRelationshipCache(cacheKey, relationship)

	return relationship, nil
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

// extractJSONFromResponse extracts JSON from a response that might contain markdown formatting
func extractJSONFromResponse(response string) string {
	// Look for JSON array between markdown code blocks
	arrayRegex := regexp.MustCompile("(?s)```(?:json)?\\s*\\n?(\\[.*?\\])\\s*```")
	matches := arrayRegex.FindStringSubmatch(response)
	if len(matches) > 1 {
		return matches[1]
	}
	
	// Look for JSON object between markdown code blocks
	objectRegex := regexp.MustCompile("(?s)```(?:json)?\\s*\\n?(\\{.*?\\})\\s*```")
	matches = objectRegex.FindStringSubmatch(response)
	if len(matches) > 1 {
		return matches[1]
	}
	
	// If no markdown blocks found, try to find JSON array directly
	arrayRegex = regexp.MustCompile("(?s)\\[.*\\]")
	matches = arrayRegex.FindStringSubmatch(response)
	if len(matches) > 0 {
		return matches[0]
	}
	
	// If no array found, try to find JSON object directly
	objectRegex = regexp.MustCompile("(?s)\\{.*\\}")
	matches = objectRegex.FindStringSubmatch(response)
	if len(matches) > 0 {
		return matches[0]
	}
	
	// Return the original response if no JSON found
	return response
}
