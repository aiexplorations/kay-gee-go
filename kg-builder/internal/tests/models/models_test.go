package models_test

import (
	"encoding/json"
	"testing"

	"kg-builder/internal/models"
)

func TestConceptSerialization(t *testing.T) {
	// Create a concept
	concept := models.Concept{
		Name:      "Machine Learning",
		Relation:  "IsA",
		RelatedTo: "Artificial Intelligence",
	}

	// Serialize to JSON
	data, err := json.Marshal(concept)
	if err != nil {
		t.Fatalf("Failed to marshal concept: %v", err)
	}

	// Deserialize from JSON
	var deserializedConcept models.Concept
	err = json.Unmarshal(data, &deserializedConcept)
	if err != nil {
		t.Fatalf("Failed to unmarshal concept: %v", err)
	}

	// Verify the deserialized concept matches the original
	if deserializedConcept.Name != concept.Name {
		t.Errorf("Expected name %s, got %s", concept.Name, deserializedConcept.Name)
	}
	if deserializedConcept.Relation != concept.Relation {
		t.Errorf("Expected relation %s, got %s", concept.Relation, deserializedConcept.Relation)
	}
	if deserializedConcept.RelatedTo != concept.RelatedTo {
		t.Errorf("Expected relatedTo %s, got %s", concept.RelatedTo, deserializedConcept.RelatedTo)
	}
}

func TestConceptJSONTags(t *testing.T) {
	// Create a concept
	concept := models.Concept{
		Name:      "Machine Learning",
		Relation:  "IsA",
		RelatedTo: "Artificial Intelligence",
	}

	// Serialize to JSON
	data, err := json.Marshal(concept)
	if err != nil {
		t.Fatalf("Failed to marshal concept: %v", err)
	}

	// Convert to string for easier inspection
	jsonStr := string(data)

	// Check that the JSON contains the expected field names
	if jsonStr != `{"name":"Machine Learning","relation":"IsA","relatedTo":"Artificial Intelligence"}` {
		t.Errorf("Unexpected JSON format: %s", jsonStr)
	}
} 