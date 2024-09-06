package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"kg-builder/internal/models"
)

func GetRelatedConcepts(concept string) ([]models.Concept, error) {
	url := "http://localhost:11434/api/generate"
	prompt := fmt.Sprintf("Given the concept '%s', provide 5 related concepts. For each, specify the relationship type. Format as JSON array with 'name', 'relation', and 'relatedTo' keys.", concept)

	requestBody, err := json.Marshal(map[string]string{
		"model":  "llama3.1:latest",
		"prompt": prompt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Response string `json:"response"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var concepts []models.Concept
	err = json.Unmarshal([]byte(result.Response), &concepts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal concepts: %w", err)
	}

	return concepts, nil
}
