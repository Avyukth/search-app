package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for our program
// type Config struct {
// 	MongoUsername                string
// 	MongoPassword                string
// 	RedisPassword                string
// 	RedisAddress                 string
// 	RedisDB                      int
// 	MongoURI                     string
// 	MongoDatabase                string
// 	MongoMaxPoolSize             uint64
// 	ServerPort                   int
// 	RedisTimeout                 time.Duration
// 	Server                       string
// 	MongoContainerName           string
// 	RedisContainerName           string
// 	MongoDBStorageCollectionName string
// 	MongoDBIndexCollectionName   string
// 	MongoDBLinkCollectionName    string
// 	IndexDirectory               string
// 	DataStoreDirectory           string
// }

// // LoadConfig loads configuration from environment variables
// func LoadConfig() (*Config, error) {
// 	var cfg Config

// 	cfg.MongoUsername = getEnv("MONGO_USERNAME", "root")
// 	cfg.MongoPassword = getEnv("MONGO_PASSWORD", "example")
// 	cfg.RedisPassword = getEnv("REDIS_PASSWORD", "yourpassword")
// 	cfg.RedisAddress = getEnv("REDIS_ADDRESS", "localhost:6379")
// 	cfg.MongoURI = getEnv("MONGODB_URI", "mongodb://%s:%s@localhost:27017")
// 	cfg.MongoDatabase = getEnv("MONGODB_DATABASE", "searchDB")
// 	cfg.Server = getEnv("SERVER", "localhost")
// 	cfg.MongoContainerName = getEnv("MONGO_CONTAINER_NAME", "mongodb-search")
// 	cfg.RedisContainerName = getEnv("REDIS_CONTAINER_NAME", "redis-search")
// 	cfg.MongoDBStorageCollectionName = getEnv("MONGODB_STORAGE_COLLECTION_NAME", "patent")
// 	cfg.MongoDBIndexCollectionName = getEnv("MONGODB_INDEX_COLLECTION_NAME", "indexPatent")
// 	cfg.MongoDBLinkCollectionName = getEnv("MONGODB_LINK_COLLECTION_NAME", "downloadLink")

// 	cfg.ServerPort = getEnvAsInt("PORT", 40051)
// 	cfg.RedisDB = getEnvAsInt("REDIS_DB", 0)
// 	cfg.MongoMaxPoolSize = getEnvAsUInt("MONGODB_MAX_POOL_SIZE", 10)
// 	cfg.RedisTimeout = time.Duration(getEnvAsInt("REDIS_TIMEOUT", 10)) * time.Second
// 	cfg.IndexDirectory = getEnv("INDEX_DIRECTORY", "index")
// 	cfg.DataStoreDirectory = getEnv("DATA_STORE_DIRECTORY", "search-data")

// 	return &cfg, nil
// }

// MongoDBConfig holds the configuration related to MongoDB.
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
