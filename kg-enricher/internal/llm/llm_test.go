package llm

import (
	"os"
	"testing"

	"kg-enricher/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestInitialize(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "llm-test-cache")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Create a config
	testCfg := &config.LLMConfig{
		URL:      "http://localhost:11434/api/generate",
		Model:    "qwen2.5:3b",
		CacheDir: tempDir,
	}
	
	// Initialize the LLM service
	err = Initialize(testCfg)
	
	// Assert that there was no error
	assert.NoError(t, err)
	
	// Assert that the config was set correctly
	assert.Equal(t, testCfg, cfg)
	assert.Equal(t, "http://localhost:11434/api/generate", cfg.URL)
	assert.Equal(t, "qwen2.5:3b", cfg.Model)
}

func TestFindRelationship(t *testing.T) {
	// Skip this test if the LLM service is not available
	t.Skip("Skipping test that requires LLM service")
	
	// Create a config
	cfg := &config.LLMConfig{
		URL:   "http://localhost:11434/api/generate",
		Model: "qwen2.5:3b",
	}
	
	// Initialize the LLM service
	err := Initialize(cfg)
	assert.NoError(t, err)
	
	// Test cases
	testCases := []struct {
		name      string
		concept1  string
		concept2  string
		wantEmpty bool
	}{
		{
			name:      "Related concepts",
			concept1:  "Machine Learning",
			concept2:  "Artificial Intelligence",
			wantEmpty: false,
		},
		{
			name:      "Unrelated concepts",
			concept1:  "Quantum Physics",
			concept2:  "Ancient History",
			wantEmpty: true,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Find the relationship
			relationship, err := FindRelationship(tc.concept1, tc.concept2)
			
			// Assert that there was no error
			assert.NoError(t, err)
			
			// Assert that the relationship is as expected
			if tc.wantEmpty {
				assert.Empty(t, relationship)
			} else {
				assert.NotEmpty(t, relationship)
			}
		})
	}
} 