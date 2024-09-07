package llm

import (
    "bufio"
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "strings"

    "kg-builder/internal/models"
)

func GetRelatedConcepts(concept string) ([]models.Concept, error) {
    url := "http://host.docker.internal:11434/api/generate"
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

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

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

    if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("error reading response: %w", err)
    }

    var concepts []models.Concept
    if err := json.Unmarshal([]byte(fullResponse.String()), &concepts); err != nil {
        fmt.Printf("Raw LLM response: %s\n", fullResponse.String())
        return nil, fmt.Errorf("failed to unmarshal concepts: %w", err)
    }

    return concepts, nil
}

func MineRelationship(concept1, concept2 string) (*models.Concept, error) {
    url := "http://host.docker.internal:11434/api/generate"
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

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

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

    if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("error reading response: %w", err)
    }

    var concept models.Concept
    if err := json.Unmarshal([]byte(fullResponse.String()), &concept); err != nil {
        fmt.Printf("Raw LLM response: %s\n", fullResponse.String())
        return nil, fmt.Errorf("failed to unmarshal concept: %w", err)
    }

    if concept.Relation == "" {
        return nil, nil // No relationship found
    }

    return &concept, nil
}