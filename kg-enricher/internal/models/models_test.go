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

func TestNodeSerialization(t *testing.T) {
	// Create a node
	node := Node{
		ID:    123,
		Name:  "Machine Learning",
		Label: "Concept",
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(node)
	if err != nil {
		t.Fatalf("Failed to marshal node: %v", err)
	}

	// Deserialize from JSON
	var deserializedNode Node
	err = json.Unmarshal(jsonData, &deserializedNode)
	if err != nil {
		t.Fatalf("Failed to unmarshal node: %v", err)
	}

	// Check that the deserialized node matches the original
	if deserializedNode.ID != node.ID {
		t.Errorf("Expected ID to be %d, got %d", node.ID, deserializedNode.ID)
	}
	if deserializedNode.Name != node.Name {
		t.Errorf("Expected Name to be '%s', got '%s'", node.Name, deserializedNode.Name)
	}
	if deserializedNode.Label != node.Label {
		t.Errorf("Expected Label to be '%s', got '%s'", node.Label, deserializedNode.Label)
	}
}

func TestRelationshipSerialization(t *testing.T) {
	// Create a relationship
	relationship := Relationship{
		Source:      "Machine Learning",
		Target:      "Artificial Intelligence",
		Type:        "IsA",
		Description: "Machine Learning is a subset of Artificial Intelligence",
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(relationship)
	if err != nil {
		t.Fatalf("Failed to marshal relationship: %v", err)
	}

	// Deserialize from JSON
	var deserializedRelationship Relationship
	err = json.Unmarshal(jsonData, &deserializedRelationship)
	if err != nil {
		t.Fatalf("Failed to unmarshal relationship: %v", err)
	}

	// Check that the deserialized relationship matches the original
	if deserializedRelationship.Source != relationship.Source {
		t.Errorf("Expected Source to be '%s', got '%s'", relationship.Source, deserializedRelationship.Source)
	}
	if deserializedRelationship.Target != relationship.Target {
		t.Errorf("Expected Target to be '%s', got '%s'", relationship.Target, deserializedRelationship.Target)
	}
	if deserializedRelationship.Type != relationship.Type {
		t.Errorf("Expected Type to be '%s', got '%s'", relationship.Type, deserializedRelationship.Type)
	}
	if deserializedRelationship.Description != relationship.Description {
		t.Errorf("Expected Description to be '%s', got '%s'", relationship.Description, deserializedRelationship.Description)
	}
}

func TestRelationshipWithoutDescription(t *testing.T) {
	// Create a relationship without a description
	relationship := Relationship{
		Source: "Machine Learning",
		Target: "Artificial Intelligence",
		Type:   "IsA",
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(relationship)
	if err != nil {
		t.Fatalf("Failed to marshal relationship: %v", err)
	}

	// Convert to string
	jsonString := string(jsonData)

	// Check that the JSON does not contain the description field
	if jsonString != `{"source":"Machine Learning","target":"Artificial Intelligence","type":"IsA"}` {
		t.Errorf("Unexpected JSON format: %s", jsonString)
	}

	// Deserialize from JSON
	var deserializedRelationship Relationship
	err = json.Unmarshal(jsonData, &deserializedRelationship)
	if err != nil {
		t.Fatalf("Failed to unmarshal relationship: %v", err)
	}

	// Check that the deserialized relationship matches the original
	if deserializedRelationship.Source != relationship.Source {
		t.Errorf("Expected Source to be '%s', got '%s'", relationship.Source, deserializedRelationship.Source)
	}
	if deserializedRelationship.Target != relationship.Target {
		t.Errorf("Expected Target to be '%s', got '%s'", relationship.Target, deserializedRelationship.Target)
	}
	if deserializedRelationship.Type != relationship.Type {
		t.Errorf("Expected Type to be '%s', got '%s'", relationship.Type, deserializedRelationship.Type)
	}
	if deserializedRelationship.Description != relationship.Description {
		t.Errorf("Expected Description to be '%s', got '%s'", relationship.Description, deserializedRelationship.Description)
	}
} 