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
	Neo4j Neo4jConfig `yaml:"neo4j"`
	LLM   LLMConfig   `yaml:"llm"`
	Graph GraphConfig `yaml:"graph"`
}

// Neo4jConfig holds Neo4j database configuration
type Neo4jConfig struct {
	URI           string        `yaml:"uri"`
	User          string        `yaml:"user"`
	Password      string        `yaml:"password"`
	MaxRetries    int           `yaml:"max_retries"`
	RetryInterval time.Duration `yaml:"-"`
	RetryIntervalSeconds int    `yaml:"retry_interval_seconds"`
}

// LLMConfig holds LLM service configuration
type LLMConfig struct {
	URL      string `yaml:"url"`
	Model    string `yaml:"model"`
	CacheDir string `yaml:"cache_dir"`
}

// GraphConfig holds graph building configuration
type GraphConfig struct {
	SeedConcept         string        `yaml:"seed_concept"`
	MaxNodes            int           `yaml:"max_nodes"`
	Timeout             time.Duration `yaml:"-"`
	TimeoutMinutes      int           `yaml:"timeout_minutes"`
	WorkerCount         int           `yaml:"worker_count"`
	RandomRelationships int           `yaml:"random_relationships"`
	Concurrency         int           `yaml:"concurrency"`
}

// LoadConfig loads configuration from the YAML file and environment variables
func LoadConfig() (*Config, error) {
	// Default configuration
	config := &Config{
		Neo4j: Neo4jConfig{
			URI:                "bolt://neo4j:7687",
			User:               "neo4j",
			Password:           "password",
			MaxRetries:         5,
			RetryIntervalSeconds: 5,
		},
		LLM: LLMConfig{
			URL:      "http://host.docker.internal:11434/api/generate",
			Model:    "qwen2.5:3b",
			CacheDir: "./cache/llm",
		},
		Graph: GraphConfig{
			SeedConcept:         "Artificial Intelligence",
			MaxNodes:            100,
			TimeoutMinutes:      30,
			WorkerCount:         10,
			RandomRelationships: 50,
			Concurrency:         5,
		},
	}

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

		if yamlConfig.Graph.SeedConcept != "" {
			config.Graph.SeedConcept = yamlConfig.Graph.SeedConcept
		}
		if yamlConfig.Graph.MaxNodes != 0 {
			config.Graph.MaxNodes = yamlConfig.Graph.MaxNodes
		}
		if yamlConfig.Graph.TimeoutMinutes != 0 {
			config.Graph.TimeoutMinutes = yamlConfig.Graph.TimeoutMinutes
		}
		if yamlConfig.Graph.WorkerCount != 0 {
			config.Graph.WorkerCount = yamlConfig.Graph.WorkerCount
		}
		if yamlConfig.Graph.RandomRelationships != 0 {
			config.Graph.RandomRelationships = yamlConfig.Graph.RandomRelationships
		}
		if yamlConfig.Graph.Concurrency != 0 {
			config.Graph.Concurrency = yamlConfig.Graph.Concurrency
		}
	}

	// Override with environment variables
	if envValue := os.Getenv("NEO4J_URI"); envValue != "" {
		config.Neo4j.URI = envValue
	}
	if envValue := os.Getenv("NEO4J_USER"); envValue != "" {
		config.Neo4j.User = envValue
	}
	if envValue := os.Getenv("NEO4J_PASSWORD"); envValue != "" {
		config.Neo4j.Password = envValue
	}
	if envValue := os.Getenv("NEO4J_MAX_RETRIES"); envValue != "" {
		if value, err := strconv.Atoi(envValue); err == nil {
			config.Neo4j.MaxRetries = value
		}
	}
	if envValue := os.Getenv("NEO4J_RETRY_INTERVAL_SECONDS"); envValue != "" {
		if value, err := strconv.Atoi(envValue); err == nil {
			config.Neo4j.RetryIntervalSeconds = value
		}
	}

	if envValue := os.Getenv("LLM_URL"); envValue != "" {
		config.LLM.URL = envValue
	}
	if envValue := os.Getenv("LLM_MODEL"); envValue != "" {
		config.LLM.Model = envValue
	}
	if envValue := os.Getenv("LLM_CACHE_DIR"); envValue != "" {
		config.LLM.CacheDir = envValue
	}

	if envValue := os.Getenv("SEED_CONCEPT"); envValue != "" {
		config.Graph.SeedConcept = envValue
	}
	if envValue := os.Getenv("MAX_NODES"); envValue != "" {
		if value, err := strconv.Atoi(envValue); err == nil {
			config.Graph.MaxNodes = value
		}
	}
	if envValue := os.Getenv("TIMEOUT_MINUTES"); envValue != "" {
		if value, err := strconv.Atoi(envValue); err == nil {
			config.Graph.TimeoutMinutes = value
		}
	}
	if envValue := os.Getenv("WORKER_COUNT"); envValue != "" {
		if value, err := strconv.Atoi(envValue); err == nil {
			config.Graph.WorkerCount = value
		}
	}
	if envValue := os.Getenv("RANDOM_RELATIONSHIPS"); envValue != "" {
		if value, err := strconv.Atoi(envValue); err == nil {
			config.Graph.RandomRelationships = value
		}
	}
	if envValue := os.Getenv("CONCURRENCY"); envValue != "" {
		if value, err := strconv.Atoi(envValue); err == nil {
			config.Graph.Concurrency = value
		}
	}

	// Convert seconds to duration
	config.Neo4j.RetryInterval = time.Duration(config.Neo4j.RetryIntervalSeconds) * time.Second
	config.Graph.Timeout = time.Duration(config.Graph.TimeoutMinutes) * time.Minute

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