package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
)

// Config holds all configuration for the application
type Config struct {
	Neo4j    Neo4jConfig    `yaml:"neo4j"`
	LLM      LLMConfig      `yaml:"llm"`
	Enricher EnricherConfig `yaml:"enricher"`
}

// Neo4jConfig holds Neo4j database configuration
type Neo4jConfig struct {
	URI                string        `yaml:"uri"`
	User               string        `yaml:"user"`
	Username           string        `yaml:"-"` // Alias for User for backward compatibility
	Password           string        `yaml:"password"`
	MaxRetries         int           `yaml:"max_retries"`
	RetryInterval      time.Duration `yaml:"-"`
	RetryIntervalSeconds int         `yaml:"retry_interval_seconds"`
}

// LLMConfig holds LLM service configuration
type LLMConfig struct {
	URL      string `yaml:"url"`
	Model    string `yaml:"model"`
	CacheDir string `yaml:"cache_dir"`
}

// EnricherConfig holds enricher configuration
type EnricherConfig struct {
	BatchSize          int           `yaml:"batch_size"`
	Interval           time.Duration `yaml:"-"`
	IntervalSeconds    int           `yaml:"interval_seconds"`
	MaxRelationships   int           `yaml:"max_relationships"`
	Concurrency        int           `yaml:"concurrency"`
}

// LoadConfig loads configuration from the YAML file and environment variables
func LoadConfig() (*Config, error) {
	// Default configuration
	config := &Config{
		Neo4j: Neo4jConfig{
			URI:                  getEnv("NEO4J_URI", "bolt://neo4j:7687"),
			User:                 getEnv("NEO4J_USER", "neo4j"),
			Password:             getEnv("NEO4J_PASSWORD", "password"),
			MaxRetries:           getEnvAsInt("NEO4J_MAX_RETRIES", 5),
			RetryIntervalSeconds: getEnvAsInt("NEO4J_RETRY_INTERVAL_SECONDS", 5),
		},
		LLM: LLMConfig{
			URL:      getEnv("LLM_URL", "http://host.docker.internal:11434/api/generate"),
			Model:    getEnv("LLM_MODEL", "qwen2.5:3b"),
			CacheDir: getEnv("LLM_CACHE_DIR", "./cache/llm"),
		},
		Enricher: EnricherConfig{
			BatchSize:        getEnvAsInt("ENRICHER_BATCH_SIZE", 10),
			IntervalSeconds:  getEnvAsInt("ENRICHER_INTERVAL_SECONDS", 60),
			MaxRelationships: getEnvAsInt("ENRICHER_MAX_RELATIONSHIPS", 100),
			Concurrency:      getEnvAsInt("ENRICHER_CONCURRENCY", 5),
		},
	}

	// Set Username as an alias for User
	config.Neo4j.Username = config.Neo4j.User

	// Try to load configuration from YAML file
	configFile := getEnv("CONFIG_FILE", "config.yaml")
	if _, err := os.Stat(configFile); err == nil {
		yamlFile, err := ioutil.ReadFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		// Parse YAML file
		var yamlConfig Config
		if err := yaml.Unmarshal(yamlFile, &yamlConfig); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}

		// Override default configuration with values from YAML file
		if yamlConfig.Neo4j.URI != "" {
			config.Neo4j.URI = yamlConfig.Neo4j.URI
		}
		if yamlConfig.Neo4j.User != "" {
			config.Neo4j.User = yamlConfig.Neo4j.User
			config.Neo4j.Username = yamlConfig.Neo4j.User // Update Username alias
		}
		if yamlConfig.Neo4j.Password != "" {
			config.Neo4j.Password = yamlConfig.Neo4j.Password
		}
		if yamlConfig.Neo4j.MaxRetries != 0 {
			config.Neo4j.MaxRetries = yamlConfig.Neo4j.MaxRetries
		}
		if yamlConfig.Neo4j.RetryIntervalSeconds != 0 {
			config.Neo4j.RetryIntervalSeconds = yamlConfig.Neo4j.RetryIntervalSeconds
		}

		if yamlConfig.LLM.URL != "" {
			config.LLM.URL = yamlConfig.LLM.URL
		}
		if yamlConfig.LLM.Model != "" {
			config.LLM.Model = yamlConfig.LLM.Model
		}
		if yamlConfig.LLM.CacheDir != "" {
			config.LLM.CacheDir = yamlConfig.LLM.CacheDir
		}

		if yamlConfig.Enricher.BatchSize != 0 {
			config.Enricher.BatchSize = yamlConfig.Enricher.BatchSize
		}
		if yamlConfig.Enricher.IntervalSeconds != 0 {
			config.Enricher.IntervalSeconds = yamlConfig.Enricher.IntervalSeconds
		}
		if yamlConfig.Enricher.MaxRelationships != 0 {
			config.Enricher.MaxRelationships = yamlConfig.Enricher.MaxRelationships
		}
		if yamlConfig.Enricher.Concurrency != 0 {
			config.Enricher.Concurrency = yamlConfig.Enricher.Concurrency
		}
	}

	// Convert seconds to duration
	config.Neo4j.RetryInterval = time.Duration(config.Neo4j.RetryIntervalSeconds) * time.Second
	config.Enricher.Interval = time.Duration(config.Enricher.IntervalSeconds) * time.Second

	// Ensure cache directory exists
	if config.LLM.CacheDir != "" {
		if err := os.MkdirAll(config.LLM.CacheDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create cache directory: %w", err)
		}
	}

	return config, nil
}

// Helper function to get an environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Helper function to get an environment variable as an integer with a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		fmt.Printf("Warning: Invalid value for %s: %s. Using default: %d\n", key, valueStr, defaultValue)
		return defaultValue
	}
	
	return value
}

// Helper function to get an environment variable as a boolean with a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		fmt.Printf("Warning: Invalid value for %s: %s. Using default: %t\n", key, valueStr, defaultValue)
		return defaultValue
	}
	
	return value
} 