package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type MongoDBConfig struct {
	Username              string
	Password              string
	URI                   string
	Database              string
	MaxPoolSize           uint64
	ContainerName         string
	StorageCollectionName string
	IndexCollectionName   string
	LinkCollectionName    string
}

// RedisConfig holds the configuration related to Redis.
type RedisConfig struct {
	Password      string
	Address       string
	DB            int
	Timeout       time.Duration
	ContainerName string
}

// ServerConfig holds the configuration related to the Server.
type ServerConfig struct {
	Port               int
	Server             string
	IndexDirectory     string
	DataStoreDirectory string
	Storage            string
	ServiceName        string
	ServiceVersion     string
}

// Config holds all configuration for our program.
type Config struct {
	MongoDBConfig
	RedisConfig
	ServerConfig
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() (*Config, error) {
	return &Config{
		MongoDBConfig: MongoDBConfig{
			Username:              getEnv("MONGO_USERNAME", "root"),
			Password:              getEnv("MONGO_PASSWORD", "example"),
			URI:                   getEnv("MONGODB_URI", "mongodb://%s:%s@localhost:27017"),
			Database:              getEnv("MONGODB_DATABASE", "searchDB"),
			MaxPoolSize:           getEnvAsUInt("MONGODB_MAX_POOL_SIZE", 10),
			ContainerName:         getEnv("MONGO_CONTAINER_NAME", "mongodb-search"),
			StorageCollectionName: getEnv("MONGODB_STORAGE_COLLECTION_NAME", "patent"),
			IndexCollectionName:   getEnv("MONGODB_INDEX_COLLECTION_NAME", "indexPatent"),
			LinkCollectionName:    getEnv("MONGODB_LINK_COLLECTION_NAME", "downloadLink"),
		},
		RedisConfig: RedisConfig{
			Password:      getEnv("REDIS_PASSWORD", "yourpassword"),
			Address:       getEnv("REDIS_ADDRESS", "localhost:6379"),
			DB:            getEnvAsInt("REDIS_DB", 0),
			Timeout:       time.Duration(getEnvAsInt("REDIS_TIMEOUT", 10)) * time.Second,
			ContainerName: getEnv("REDIS_CONTAINER_NAME", "redis-search"),
		},
		ServerConfig: ServerConfig{
			Port:               getEnvAsInt("PORT", 40051),
			Server:             getEnv("SERVER", "localhost"),
			IndexDirectory:     getEnv("INDEX_DIRECTORY", "index"),
			DataStoreDirectory: getEnv("DATA_STORE_DIRECTORY", "search-data"),
			Storage:            getEnv("STORAGE_DIRECTORY", "storage"),
			ServiceName:        getEnv("SERVICE_NAME", "search-app"),
			ServiceVersion:     getEnv("SERVICE_VERSION", "1.0.0"),
		},
	}, nil
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
