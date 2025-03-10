package models

import (
	"encoding/json"
	"testing"
)

func TestConceptSerialization(t *testing.T) {
	// Create a concept
	concept := Concept{
		Name:      "Machine Learning",
		Relation:  "IsA",
		RelatedTo: "Artificial Intelligence",
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(concept)
	if err != nil {
		t.Fatalf("Failed to marshal concept: %v", err)
	}

	// Deserialize from JSON
	var deserializedConcept Concept
	err = json.Unmarshal(jsonData, &deserializedConcept)
	if err != nil {
		t.Fatalf("Failed to unmarshal concept: %v", err)
	}

	// Check that the deserialized concept matches the original
	if deserializedConcept.Name != concept.Name {
		t.Errorf("Expected Name to be '%s', got '%s'", concept.Name, deserializedConcept.Name)
	}
	if deserializedConcept.Relation != concept.Relation {
		t.Errorf("Expected Relation to be '%s', got '%s'", concept.Relation, deserializedConcept.Relation)
	}
	if deserializedConcept.RelatedTo != concept.RelatedTo {
		t.Errorf("Expected RelatedTo to be '%s', got '%s'", concept.RelatedTo, deserializedConcept.RelatedTo)
	}
}

func TestConceptJSONTags(t *testing.T) {
	// Create a concept
	concept := Concept{
		Name:      "Machine Learning",
		Relation:  "IsA",
		RelatedTo: "Artificial Intelligence",
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(concept)
	if err != nil {
		t.Fatalf("Failed to marshal concept: %v", err)
	}

	// Convert to string
	jsonString := string(jsonData)

	// Check that the JSON contains the expected field names
	if jsonString != `{"name":"Machine Learning","relation":"IsA","relatedTo":"Artificial Intelligence"}` {
		t.Errorf("Unexpected JSON format: %s", jsonString)
	}
}

func TestConceptEmptyFields(t *testing.T) {
	// Create a concept with empty fields
	concept := Concept{
		Name:      "",
		Relation:  "",
		RelatedTo: "",
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(concept)
	if err != nil {
		t.Fatalf("Failed to marshal concept: %v", err)
	}

	// Deserialize from JSON
	var deserializedConcept Concept
	err = json.Unmarshal(jsonData, &deserializedConcept)
	if err != nil {
		t.Fatalf("Failed to unmarshal concept: %v", err)
	}

	// Check that the deserialized concept matches the original
	if deserializedConcept.Name != concept.Name {
		t.Errorf("Expected Name to be '%s', got '%s'", concept.Name, deserializedConcept.Name)
	}
	if deserializedConcept.Relation != concept.Relation {
		t.Errorf("Expected Relation to be '%s', got '%s'", concept.Relation, deserializedConcept.Relation)
	}
	if deserializedConcept.RelatedTo != concept.RelatedTo {
		t.Errorf("Expected RelatedTo to be '%s', got '%s'", concept.RelatedTo, deserializedConcept.RelatedTo)
	}
} 