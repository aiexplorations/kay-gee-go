package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeModel(t *testing.T) {
	// Create a node
	node := Node{
		ID:    "1",
		Name:  "Test Node",
		Size:  5,
		Color: "#ff0000",
		Type:  "Concept",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(node)
	assert.NoError(t, err)

	// Test JSON unmarshaling
	var unmarshaledNode Node
	err = json.Unmarshal(jsonData, &unmarshaledNode)
	assert.NoError(t, err)

	// Verify the unmarshaled node
	assert.Equal(t, node.ID, unmarshaledNode.ID)
	assert.Equal(t, node.Name, unmarshaledNode.Name)
	assert.Equal(t, node.Size, unmarshaledNode.Size)
	assert.Equal(t, node.Color, unmarshaledNode.Color)
	assert.Equal(t, node.Type, unmarshaledNode.Type)
}

func TestLinkModel(t *testing.T) {
	// Create a link
	link := Link{
		Source: "1",
		Target: "2",
		Type:   "RELATES_TO",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(link)
	assert.NoError(t, err)

	// Test JSON unmarshaling
	var unmarshaledLink Link
	err = json.Unmarshal(jsonData, &unmarshaledLink)
	assert.NoError(t, err)

	// Verify the unmarshaled link
	assert.Equal(t, link.Source, unmarshaledLink.Source)
	assert.Equal(t, link.Target, unmarshaledLink.Target)
	assert.Equal(t, link.Type, unmarshaledLink.Type)
}

func TestGraphDataModel(t *testing.T) {
	// Create nodes
	nodes := []Node{
		{ID: "1", Name: "Node 1", Size: 5},
		{ID: "2", Name: "Node 2", Size: 7},
	}

	// Create links
	links := []Link{
		{Source: "1", Target: "2", Type: "RELATES_TO"},
	}

	// Create graph data
	graphData := GraphData{
		Nodes: nodes,
		Links: links,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(graphData)
	assert.NoError(t, err)

	// Test JSON unmarshaling
	var unmarshaledGraphData GraphData
	err = json.Unmarshal(jsonData, &unmarshaledGraphData)
	assert.NoError(t, err)

	// Verify the unmarshaled graph data
	assert.Equal(t, len(graphData.Nodes), len(unmarshaledGraphData.Nodes))
	assert.Equal(t, len(graphData.Links), len(unmarshaledGraphData.Links))
	assert.Equal(t, graphData.Nodes[0].ID, unmarshaledGraphData.Nodes[0].ID)
	assert.Equal(t, graphData.Nodes[1].Name, unmarshaledGraphData.Nodes[1].Name)
	assert.Equal(t, graphData.Links[0].Source, unmarshaledGraphData.Links[0].Source)
	assert.Equal(t, graphData.Links[0].Target, unmarshaledGraphData.Links[0].Target)
}

func TestSearchResultModel(t *testing.T) {
	// Create a search result
	searchResult := SearchResult{
		ID:   "1",
		Name: "Test Result",
		Type: "Concept",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(searchResult)
	assert.NoError(t, err)

	// Test JSON unmarshaling
	var unmarshaledSearchResult SearchResult
	err = json.Unmarshal(jsonData, &unmarshaledSearchResult)
	assert.NoError(t, err)

	// Verify the unmarshaled search result
	assert.Equal(t, searchResult.ID, unmarshaledSearchResult.ID)
	assert.Equal(t, searchResult.Name, unmarshaledSearchResult.Name)
	assert.Equal(t, searchResult.Type, unmarshaledSearchResult.Type)
}

func TestRelationshipModel(t *testing.T) {
	// Create a relationship
	relationship := Relationship{
		SourceID: "1",
		TargetID: "2",
		Type:     "RELATES_TO",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(relationship)
	assert.NoError(t, err)

	// Test JSON unmarshaling
	var unmarshaledRelationship Relationship
	err = json.Unmarshal(jsonData, &unmarshaledRelationship)
	assert.NoError(t, err)

	// Verify the unmarshaled relationship
	assert.Equal(t, relationship.SourceID, unmarshaledRelationship.SourceID)
	assert.Equal(t, relationship.TargetID, unmarshaledRelationship.TargetID)
	assert.Equal(t, relationship.Type, unmarshaledRelationship.Type)
}

func TestStatisticsModel(t *testing.T) {
	// Create type counts
	nodeTypes := []TypeCount{
		{Type: "Concept", Count: 75},
		{Type: "Entity", Count: 25},
	}

	relationshipTypes := []TypeCount{
		{Type: "RELATES_TO", Count: 150},
		{Type: "PART_OF", Count: 100},
	}

	// Create statistics
	statistics := Statistics{
		NodeCount:         100,
		RelationshipCount: 250,
		NodeTypes:         nodeTypes,
		RelationshipTypes: relationshipTypes,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(statistics)
	assert.NoError(t, err)

	// Test JSON unmarshaling
	var unmarshaledStatistics Statistics
	err = json.Unmarshal(jsonData, &unmarshaledStatistics)
	assert.NoError(t, err)

	// Verify the unmarshaled statistics
	assert.Equal(t, statistics.NodeCount, unmarshaledStatistics.NodeCount)
	assert.Equal(t, statistics.RelationshipCount, unmarshaledStatistics.RelationshipCount)
	assert.Equal(t, len(statistics.NodeTypes), len(unmarshaledStatistics.NodeTypes))
	assert.Equal(t, len(statistics.RelationshipTypes), len(unmarshaledStatistics.RelationshipTypes))
	assert.Equal(t, statistics.NodeTypes[0].Type, unmarshaledStatistics.NodeTypes[0].Type)
	assert.Equal(t, statistics.NodeTypes[0].Count, unmarshaledStatistics.NodeTypes[0].Count)
	assert.Equal(t, statistics.RelationshipTypes[0].Type, unmarshaledStatistics.RelationshipTypes[0].Type)
	assert.Equal(t, statistics.RelationshipTypes[0].Count, unmarshaledStatistics.RelationshipTypes[0].Count)
}

func TestTypeCountModel(t *testing.T) {
	// Create a type count
	typeCount := TypeCount{
		Type:  "Concept",
		Count: 75,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(typeCount)
	assert.NoError(t, err)

	// Test JSON unmarshaling
	var unmarshaledTypeCount TypeCount
	err = json.Unmarshal(jsonData, &unmarshaledTypeCount)
	assert.NoError(t, err)

	// Verify the unmarshaled type count
	assert.Equal(t, typeCount.Type, unmarshaledTypeCount.Type)
	assert.Equal(t, typeCount.Count, unmarshaledTypeCount.Count)
} 