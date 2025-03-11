package llm

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kay-gee-go/internal/common/config"
	"github.com/kay-gee-go/internal/common/errors"
	"github.com/kay-gee-go/internal/common/models"
)

// Client represents an LLM client
type Client struct {
	config config.LLMConfig
}

// NewClient creates a new LLM client
func NewClient(config config.LLMConfig) *Client {
	return &Client{
		config: config,
	}
}

// GetRelatedConcepts retrieves concepts related to the given concept
func (c *Client) GetRelatedConcepts(concept string) ([]models.RelatedConcept, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("related_concepts_%s", concept)
	cachedResponse, err := c.loadFromCache(cacheKey)
	if err == nil && cachedResponse != nil {
		var relatedConcepts []models.RelatedConcept
		if err := json.Unmarshal(cachedResponse, &relatedConcepts); err == nil {
			return relatedConcepts, nil
		}
	}

	// Prepare the prompt
	prompt := fmt.Sprintf(`
You are a knowledge graph builder. Your task is to identify concepts related to "%s".
For each related concept, provide:
1. The name of the concept
2. A brief description (1-2 sentences)
3. A relevance score from 0.0 to 1.0 indicating how closely related it is to "%s"

Return your response as a JSON array with the following structure:
[
  {
    "name": "Related Concept Name",
    "description": "Brief description of the concept",
    "relevance": 0.95
  },
  ...
]

Provide 5-10 related concepts, focusing on the most relevant ones.
`, concept, concept)

	// Call the LLM API
	response, err := c.callLLM(prompt)
	if err != nil {
		return nil, err
	}

	// Extract the JSON part from the response
	jsonStr := extractJSON(response)
	if jsonStr == "" {
		return nil, errors.NewLLMError("failed to extract JSON from LLM response", nil)
	}

	// Parse the response
	var relatedConcepts []models.RelatedConcept
	if err := json.Unmarshal([]byte(jsonStr), &relatedConcepts); err != nil {
		return nil, errors.NewLLMError("failed to parse LLM response", err)
	}

	// Cache the response
	c.saveToCache(cacheKey, []byte(jsonStr))

	return relatedConcepts, nil
}

// GetRelationship determines if a relationship exists between two concepts
func (c *Client) GetRelationship(concept1, concept2 string) (*models.Relationship, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("relationship_%s_%s", concept1, concept2)
	cachedResponse, err := c.loadFromCache(cacheKey)
	if err == nil && cachedResponse != nil {
		var relationship models.Relationship
		if err := json.Unmarshal(cachedResponse, &relationship); err == nil {
			return &relationship, nil
		}
	}

	// Prepare the prompt
	prompt := fmt.Sprintf(`
You are a knowledge graph builder. Your task is to determine if a relationship exists between the concepts "%s" and "%s".

If a relationship exists, provide:
1. The type of relationship (e.g., "IS_A", "PART_OF", "USED_IN", etc.)
2. A brief description of the relationship (1-2 sentences)
3. A strength score from 0.0 to 1.0 indicating how strong the relationship is

Return your response as a JSON object with the following structure:
{
  "source_id": "%s",
  "target_id": "%s",
  "type": "RELATIONSHIP_TYPE",
  "description": "Description of the relationship",
  "strength": 0.85
}

If no meaningful relationship exists, return null.
`, concept1, concept2, concept1, concept2)

	// Call the LLM API
	response, err := c.callLLM(prompt)
	if err != nil {
		return nil, err
	}

	// Extract the JSON part from the response
	jsonStr := extractJSON(response)
	if jsonStr == "" || jsonStr == "null" {
		return nil, nil
	}

	// Parse the response
	var relationship models.Relationship
	if err := json.Unmarshal([]byte(jsonStr), &relationship); err != nil {
		return nil, errors.NewLLMError("failed to parse LLM response", err)
	}

	// Cache the response
	c.saveToCache(cacheKey, []byte(jsonStr))

	return &relationship, nil
}

// Helper functions

// callLLM calls the LLM API with the given prompt
func (c *Client) callLLM(prompt string) (string, error) {
	// Prepare the request
	requestBody, err := json.Marshal(map[string]interface{}{
		"model":  c.config.Model,
		"prompt": prompt,
	})
	if err != nil {
		return "", errors.NewLLMError("failed to marshal request", err)
	}

	// Send the request
	resp, err := http.Post(c.config.URL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", errors.NewLLMError("failed to call LLM API", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.NewLLMError("failed to read LLM response", err)
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		return "", errors.NewLLMError(fmt.Sprintf("LLM API returned status %d: %s", resp.StatusCode, string(body)), nil)
	}

	// Parse the response
	var response struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", errors.NewLLMError("failed to parse LLM API response", err)
	}

	return response.Response, nil
}

// extractJSON extracts the JSON part from the LLM response
func extractJSON(response string) string {
	// Find the first opening bracket
	start := strings.Index(response, "[")
	if start == -1 {
		start = strings.Index(response, "{")
		if start == -1 {
			return ""
		}
	}

	// Find the matching closing bracket
	var end int
	if response[start] == '[' {
		end = findMatchingBracket(response, start, '[', ']')
	} else {
		end = findMatchingBracket(response, start, '{', '}')
	}

	if end == -1 {
		return ""
	}

	return response[start : end+1]
}

// findMatchingBracket finds the matching closing bracket
func findMatchingBracket(s string, start int, openBracket, closeBracket rune) int {
	count := 0
	for i := start; i < len(s); i++ {
		if rune(s[i]) == openBracket {
			count++
		} else if rune(s[i]) == closeBracket {
			count--
			if count == 0 {
				return i
			}
		}
	}
	return -1
}

// loadFromCache loads a response from the cache
func (c *Client) loadFromCache(key string) ([]byte, error) {
	if c.config.CacheDir == "" {
		return nil, fmt.Errorf("cache directory not configured")
	}

	// Create a hash of the key
	hash := md5.Sum([]byte(key))
	filename := filepath.Join(c.config.CacheDir, hex.EncodeToString(hash[:])+".json")

	// Check if the file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("cache file not found")
	}

	// Read the file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// saveToCache saves a response to the cache
func (c *Client) saveToCache(key string, data []byte) error {
	if c.config.CacheDir == "" {
		return fmt.Errorf("cache directory not configured")
	}

	// Create the cache directory if it doesn't exist
	if err := os.MkdirAll(c.config.CacheDir, 0755); err != nil {
		return err
	}

	// Create a hash of the key
	hash := md5.Sum([]byte(key))
	filename := filepath.Join(c.config.CacheDir, hex.EncodeToString(hash[:])+".json")

	// Write the file
	return ioutil.WriteFile(filename, data, 0644)
} 