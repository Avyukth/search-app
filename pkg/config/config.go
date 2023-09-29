package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for our program
type Config struct {
	MongoUsername                string
	MongoPassword                string
	RedisPassword                string
	RedisAddress                 string
	RedisDB                      int
	MongoURI                     string
	MongoDatabase                string
	MongoMaxPoolSize             uint64
	ServerPort                   int
	RedisTimeout                 time.Duration
	Server                       string
	MongoContainerName           string
	RedisContainerName           string
	MongoDBStorageCollectionName string
	MongoDBIndexCollectionName   string
	MongoDBLinkCollectionName    string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	var cfg Config

	cfg.MongoUsername = getEnv("MONGO_USERNAME", "root")
	cfg.MongoPassword = getEnv("MONGO_PASSWORD", "example")
	cfg.RedisPassword = getEnv("REDIS_PASSWORD", "yourpassword")
	cfg.RedisAddress = getEnv("REDIS_ADDRESS", "localhost:6379")
	cfg.MongoURI = getEnv("MONGODB_URI", "mongodb://%s:%s@localhost:27017")
	cfg.MongoDatabase = getEnv("MONGODB_DATABASE", "searchDB")
	cfg.Server = getEnv("SERVER", "localhost")
	cfg.MongoContainerName = getEnv("MONGO_CONTAINER_NAME", "mongodb-search")
	cfg.RedisContainerName = getEnv("REDIS_CONTAINER_NAME", "redis-search")
	cfg.MongoDBStorageCollectionName = getEnv("MONGODB_STORAGE_COLLECTION_NAME", "patent")
	cfg.MongoDBIndexCollectionName = getEnv("MONGODB_INDEX_COLLECTION_NAME", "indexPatent")
	cfg.MongoDBLinkCollectionName = getEnv("MONGODB_LINK_COLLECTION_NAME", "downloadLink")

	cfg.ServerPort = getEnvAsInt("PORT", 40051)
	cfg.RedisDB = getEnvAsInt("REDIS_DB", 0)
	cfg.MongoMaxPoolSize = getEnvAsUInt("MONGODB_MAX_POOL_SIZE", 10)
	cfg.RedisTimeout = time.Duration(getEnvAsInt("REDIS_TIMEOUT", 10)) * time.Second

	return &cfg, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Printf("Warning: %s environment variable not set, using default: %s", key, defaultValue)
		return defaultValue
	}
	return value
}

// getEnvAsInt gets an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		log.Printf("Warning: %s environment variable not set, using default: %d", key, defaultValue)
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: %s environment variable is not a valid integer, using default: %d", key, defaultValue)
		return defaultValue
	}
	return value
}

// getEnvAsUInt gets an environment variable as an unsigned integer or returns a default value
func getEnvAsUInt(key string, defaultValue uint64) uint64 {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		log.Printf("Warning: %s environment variable not set, using default: %d", key, defaultValue)
		return defaultValue
	}
	value, err := strconv.ParseUint(valueStr, 10, 64)
	if err != nil {
		log.Printf("Warning: %s environment variable is not a valid unsigned integer, using default: %d", key, defaultValue)
		return defaultValue
	}
	return value
}
