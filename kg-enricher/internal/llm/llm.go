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

	"kg-enricher/internal/config"
	"kg-enricher/internal/models"
)

const (
	DefaultMaxRetries    = 3
	DefaultRetryInterval = 2 * time.Second
	DefaultMaxBackoff    = 15 * time.Second
)

var (
	// Cache for LLM responses
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
	loadRelationshipCache()
	
	return nil
}

// MineRelationship sends a request to the LLM service to determine if there is a relationship between two concepts.
func MineRelationship(concept1, concept2 string) (*models.Concept, error) {
	if cfg == nil {
		return nil, fmt.Errorf("LLM service not initialized")
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
	
Determine if there's a REAL, FACTUAL relationship between the concepts '%s' and '%s'.

IMPORTANT RULES:
1. ONLY identify relationships that actually exist in the real world
2. Use standard, established relationship types (e.g., IS_A, PART_OF, USED_IN, DEVELOPED_BY)
3. If no meaningful relationship exists, return an empty response
4. Never invent or make up relationships
5. Be specific and precise in describing the relationship

Return the response as a JSON object with 'name', 'relation', and 'relatedTo' keys. The response should be valid JSON that can be directly parsed. 

Example format for a real relationship:
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

	for attempt := 1; attempt <= DefaultMaxRetries; attempt++ {
		var requestBody []byte
		requestBody, err = json.Marshal(map[string]string{
			"model":  cfg.Model,
			"prompt": prompt,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}

		// Send the request to the LLM service
		resp, err := http.Post(cfg.URL, "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			log.Printf("Attempt %d/%d: Failed to make request: %v", attempt, DefaultMaxRetries, err)
			if attempt < DefaultMaxRetries {
				time.Sleep(DefaultRetryInterval)
				continue
			}
			return nil, fmt.Errorf("failed to make request: %w", err)
		}
		defer resp.Body.Close()

		// Check if the response status code is OK
		if resp.StatusCode != http.StatusOK {
			log.Printf("Attempt %d/%d: LLM service returned non-OK status: %d", attempt, DefaultMaxRetries, resp.StatusCode)
			if attempt < DefaultMaxRetries && resp.StatusCode >= 500 {
				time.Sleep(DefaultRetryInterval)
				continue
			}
			return nil, fmt.Errorf("LLM service returned non-OK status: %d", resp.StatusCode)
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
			log.Printf("Attempt %d/%d: Error reading response: %v", attempt, DefaultMaxRetries, err)
			if attempt < DefaultMaxRetries {
				time.Sleep(DefaultRetryInterval)
				continue
			}
			return nil, fmt.Errorf("error reading response: %w", err)
		}

		// Unmarshal the response into a Concept struct
		var conceptData models.Concept
		if err := json.Unmarshal([]byte(fullResponse.String()), &conceptData); err != nil {
			log.Printf("Raw LLM response: %s", fullResponse.String())
			log.Printf("Attempt %d/%d: Failed to unmarshal concept: %v", attempt, DefaultMaxRetries, err)
			if attempt < DefaultMaxRetries {
				time.Sleep(DefaultRetryInterval)
				continue
			}
			return nil, fmt.Errorf("failed to unmarshal concept: %w", err)
		}

		// Check if the relationship is empty
		if conceptData.Relation == "" || conceptData.Relation == "No relationship" {
			concept = nil
		} else {
			concept = &conceptData
		}

		break
	}

	// Cache the result
	cacheMutex.Lock()
	relationshipCache[cacheKey] = concept
	cacheMutex.Unlock()
	
	// Save to disk
	saveRelationshipCache(cacheKey, concept)

	return concept, nil
}

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

// loadRelationshipCache loads cached relationship data from disk
func loadRelationshipCache() {
	if cfg == nil || cfg.CacheDir == "" {
		return
	}
	
	files, err := filepath.Glob(filepath.Join(cfg.CacheDir, "rel_*.json"))
	if err != nil {
		log.Printf("Failed to list cache files: %v", err)
		return
	}
	
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Printf("Failed to read cache file %s: %v", file, err)
			continue
		}
		
		var concept models.Concept
		if err := json.Unmarshal(data, &concept); err != nil {
			log.Printf("Failed to unmarshal cache file %s: %v", file, err)
			continue
		}
		
		// Extract the cache key from the filename
		filename := filepath.Base(file)
		cacheKey := strings.TrimPrefix(filename, "rel_")
		cacheKey = strings.TrimSuffix(cacheKey, ".json")
		cacheKey = strings.Replace(cacheKey, "_", "|", 1)
		
		cacheMutex.Lock()
		relationshipCache[cacheKey] = &concept
		cacheMutex.Unlock()
	}
	
	log.Printf("Loaded %d relationships from cache", len(relationshipCache))
}

// saveRelationshipCache saves a relationship to disk
func saveRelationshipCache(cacheKey string, concept *models.Concept) {
	if cfg == nil || cfg.CacheDir == "" {
		return
	}
	
	// Create a safe filename from the cache key
	filename := strings.Replace(cacheKey, "|", "_", -1)
	filename = fmt.Sprintf("rel_%s.json", filename)
	filepath := filepath.Join(cfg.CacheDir, filename)
	
	// If the concept is nil, we don't need to save it
	if concept == nil {
		// Create an empty concept to indicate no relationship
		concept = &models.Concept{
			Name:      "",
			Relation:  "",
			RelatedTo: "",
		}
	}
	
	data, err := json.Marshal(concept)
	if err != nil {
		log.Printf("Failed to marshal concept for cache: %v", err)
		return
	}
	
	if err := ioutil.WriteFile(filepath, data, 0644); err != nil {
		log.Printf("Failed to write cache file %s: %v", filepath, err)
	}
}

// FindRelationship finds a relationship between two concepts using the LLM service.
// It returns the relationship name or an empty string if no relationship is found.
func FindRelationship(concept1, concept2 string) (string, error) {
	concept, err := MineRelationship(concept1, concept2)
	if err != nil {
		return "", fmt.Errorf("failed to mine relationship: %w", err)
	}
	
	// If no relationship was found, return an empty string
	if concept == nil || concept.Relation == "" {
		return "", nil
	}
	
	return concept.Relation, nil
} 