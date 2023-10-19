package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type MongoDBConfig struct {
	Host                  string
	Port                  int
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
	ServiceHost        string
	ServicePort        int
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

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Define the default values

	viper.SetDefault("MONGODB_USERNAME", "root")
	viper.SetDefault("MONGODB_PASSWORD", "example")
	viper.SetDefault("MONGODB_DATABASE", "searchDB")
	viper.SetDefault("MONGODB_MAX_POOL_SIZE", 10)
	viper.SetDefault("MONGODB_STORAGE_COLLECTION_NAME", "patent")
	viper.SetDefault("MONGODB_INDEX_COLLECTION_NAME", "indexPatent")
	viper.SetDefault("MONGODB_LINK_COLLECTION_NAME", "downloadLink")
	viper.SetDefault("REDIS_PASSWORD", "yourpassword")
	viper.SetDefault("REDIS_ADDRESS", "localhost:6379")
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("SERVER_PORT", 40051)
	viper.SetDefault("SERVER_HOST", "localhost")
	viper.SetDefault("REDIS_TIMEOUT", 10)
	viper.SetDefault("SERVER", "localhost")
	viper.SetDefault("MONGO_CONTAINER_NAME", "mongodb-search")
	viper.SetDefault("REDIS_CONTAINER_NAME", "redis-search")
	viper.SetDefault("INDEX_DIRECTORY", "/index")
	viper.SetDefault("DATA_STORE_DIRECTORY", "./search-data")
	viper.SetDefault("STORAGE_DIRECTORY", "./storage")
	viper.SetDefault("VERSION", "1.0")

	// Read the configuration from environment variables
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
		return nil, err
	}

	return &Config{
		MongoDBConfig: MongoDBConfig{
			Host:                  viper.GetString("MONGODB_HOST"),
			Port:                  viper.GetInt("MONGODB_PORT"),
			Username:              viper.GetString("MONGODB_USERNAME"),
			Password:              viper.GetString("MONGODB_PASSWORD"),
			Database:              viper.GetString("MONGODB_DATABASE"),
			MaxPoolSize:           viper.GetUint64("MONGODB_MAX_POOL_SIZE"),
			URI:                   viper.GetString("MONGODB_URI"),
			ContainerName:         viper.GetString("MONGO_CONTAINER_NAME"),
			StorageCollectionName: viper.GetString("MONGODB_STORAGE_COLLECTION_NAME"),
			IndexCollectionName:   viper.GetString("MONGODB_INDEX_COLLECTION_NAME"),
			LinkCollectionName:    viper.GetString("MONGODB_LINK_COLLECTION_NAME"),
		},
		RedisConfig: RedisConfig{
			Password:      viper.GetString("REDIS_PASSWORD"),
			Address:       viper.GetString("REDIS_ADDRESS"),
			DB:            viper.GetInt("REDIS_DB"),
			Timeout:       time.Duration(viper.GetInt("REDIS_TIMEOUT")) * time.Second,
			ContainerName: viper.GetString("REDIS_CONTAINER_NAME"),
		},
		ServerConfig: ServerConfig{
			ServiceHost:        viper.GetString("SERVER_HOST"),
			ServicePort:        viper.GetInt("SERVER_PORT"),
			Server:             viper.GetString("SERVER"),
			IndexDirectory:     viper.GetString("INDEX_DIRECTORY"),
			DataStoreDirectory: viper.GetString("DATA_STORE_DIRECTORY"),
			Storage:            viper.GetString("STORAGE_DIRECTORY"),
			ServiceName:        viper.GetString("SERVICE_NAME"),
			ServiceVersion:     viper.GetString("VERSION"),
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
