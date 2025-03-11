package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Neo4jConfig represents the Neo4j database configuration
type Neo4jConfig struct {
	URI                string `mapstructure:"uri"`
	User               string `mapstructure:"user"`
	Password           string `mapstructure:"password"`
	MaxRetries         int    `mapstructure:"max_retries"`
	RetryIntervalSecs  int    `mapstructure:"retry_interval_seconds"`
}

// LLMConfig represents the LLM service configuration
type LLMConfig struct {
	URL      string `mapstructure:"url"`
	Model    string `mapstructure:"model"`
	CacheDir string `mapstructure:"cache_dir"`
}

// GraphConfig represents the graph building configuration
type GraphConfig struct {
	SeedConcept        string `mapstructure:"seed_concept"`
	MaxNodes           int    `mapstructure:"max_nodes"`
	TimeoutMinutes     int    `mapstructure:"timeout_minutes"`
	WorkerCount        int    `mapstructure:"worker_count"`
	RandomRelationships int   `mapstructure:"random_relationships"`
	Concurrency        int    `mapstructure:"concurrency"`
}

// EnricherConfig represents the graph enricher configuration
type EnricherConfig struct {
	BatchSize        int `mapstructure:"batch_size"`
	IntervalSeconds  int `mapstructure:"interval_seconds"`
	MaxRelationships int `mapstructure:"max_relationships"`
	Concurrency      int `mapstructure:"concurrency"`
}

// BuilderConfig represents the complete builder configuration
type BuilderConfig struct {
	Neo4j Neo4jConfig `mapstructure:"neo4j"`
	LLM   LLMConfig   `mapstructure:"llm"`
	Graph GraphConfig `mapstructure:"graph"`
}

// EnricherAppConfig represents the complete enricher configuration
type EnricherAppConfig struct {
	Neo4j    Neo4jConfig    `mapstructure:"neo4j"`
	LLM      LLMConfig      `mapstructure:"llm"`
	Enricher EnricherConfig `mapstructure:"enricher"`
}

// LoadBuilderConfig loads the builder configuration from a YAML file and environment variables
func LoadBuilderConfig(configPath string) (*BuilderConfig, error) {
	v := viper.New()
	
	// Set default values
	setBuilderDefaults(v)
	
	// Read from config file
	if configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}
	
	// Override with environment variables
	v.AutomaticEnv()
	v.SetEnvPrefix("")
	
	// Map environment variables to config fields
	mapBuilderEnvVars(v)
	
	// Unmarshal config
	var config BuilderConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	// Create cache directory if it doesn't exist
	if config.LLM.CacheDir != "" {
		if err := os.MkdirAll(config.LLM.CacheDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create cache directory: %w", err)
		}
	}
	
	return &config, nil
}

// LoadEnricherConfig loads the enricher configuration from a YAML file and environment variables
func LoadEnricherConfig(configPath string) (*EnricherAppConfig, error) {
	v := viper.New()
	
	// Set default values
	setEnricherDefaults(v)
	
	// Read from config file
	if configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}
	
	// Override with environment variables
	v.AutomaticEnv()
	v.SetEnvPrefix("")
	
	// Map environment variables to config fields
	mapEnricherEnvVars(v)
	
	// Unmarshal config
	var config EnricherAppConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	// Create cache directory if it doesn't exist
	if config.LLM.CacheDir != "" {
		if err := os.MkdirAll(config.LLM.CacheDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create cache directory: %w", err)
		}
	}
	
	return &config, nil
}

// Helper functions

func setBuilderDefaults(v *viper.Viper) {
	// Neo4j defaults
	v.SetDefault("neo4j.uri", "bolt://neo4j:7687")
	v.SetDefault("neo4j.user", "neo4j")
	v.SetDefault("neo4j.password", "password")
	v.SetDefault("neo4j.max_retries", 5)
	v.SetDefault("neo4j.retry_interval_seconds", 5)
	
	// LLM defaults
	v.SetDefault("llm.url", "http://host.docker.internal:11434/api/generate")
	v.SetDefault("llm.model", "qwen2.5:3b")
	v.SetDefault("llm.cache_dir", "./cache")
	
	// Graph defaults
	v.SetDefault("graph.seed_concept", "Artificial Intelligence")
	v.SetDefault("graph.max_nodes", 100)
	v.SetDefault("graph.timeout_minutes", 30)
	v.SetDefault("graph.worker_count", 10)
	v.SetDefault("graph.random_relationships", 50)
	v.SetDefault("graph.concurrency", 5)
}

func setEnricherDefaults(v *viper.Viper) {
	// Neo4j defaults
	v.SetDefault("neo4j.uri", "bolt://neo4j:7687")
	v.SetDefault("neo4j.user", "neo4j")
	v.SetDefault("neo4j.password", "password")
	v.SetDefault("neo4j.max_retries", 5)
	v.SetDefault("neo4j.retry_interval_seconds", 5)
	
	// LLM defaults
	v.SetDefault("llm.url", "http://host.docker.internal:11434/api/generate")
	v.SetDefault("llm.model", "qwen2.5:3b")
	v.SetDefault("llm.cache_dir", "./cache")
	
	// Enricher defaults
	v.SetDefault("enricher.batch_size", 10)
	v.SetDefault("enricher.interval_seconds", 60)
	v.SetDefault("enricher.max_relationships", 100)
	v.SetDefault("enricher.concurrency", 5)
}

func mapBuilderEnvVars(v *viper.Viper) {
	// Neo4j environment variables
	v.BindEnv("neo4j.uri", "NEO4J_URI")
	v.BindEnv("neo4j.user", "NEO4J_USER")
	v.BindEnv("neo4j.password", "NEO4J_PASSWORD")
	
	// LLM environment variables
	v.BindEnv("llm.url", "LLM_URL")
	v.BindEnv("llm.model", "LLM_MODEL")
}

func mapEnricherEnvVars(v *viper.Viper) {
	// Neo4j environment variables
	v.BindEnv("neo4j.uri", "NEO4J_URI")
	v.BindEnv("neo4j.user", "NEO4J_USER")
	v.BindEnv("neo4j.password", "NEO4J_PASSWORD")
	
	// LLM environment variables
	v.BindEnv("llm.url", "LLM_URL")
	v.BindEnv("llm.model", "LLM_MODEL")
	
	// Enricher environment variables
	v.BindEnv("enricher.batch_size", "ENRICHER_BATCH_SIZE")
	v.BindEnv("enricher.interval_seconds", "ENRICHER_INTERVAL_SECONDS")
	v.BindEnv("enricher.max_relationships", "ENRICHER_MAX_RELATIONSHIPS")
	v.BindEnv("enricher.concurrency", "ENRICHER_CONCURRENCY")
} 