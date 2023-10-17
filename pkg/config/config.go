package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type MongoDBConfig struct {
	// Connection details
	Host          string
	Port          int
	Username      string
	Password      string
	URI           string
	Database      string
	MaxPoolSize   uint64
	ContainerName string

	// Collection names
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
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	viper.SetDefault("MONGO_HOST", "localhost")
	viper.SetDefault("MONGO_PORT", 27017)
	viper.SetDefault("MONGO_USERNAME", "root")
	viper.SetDefault("MONGO_PASSWORD", "example")
	viper.SetDefault("MONGO_DATABASE", "search")
	viper.SetDefault("MONGO_MAX_POOL_SIZE", 100)
	viper.SetDefault("MONGO_CONTAINER_NAME", "")
	viper.SetDefault("MONGO_URI", "")
	viper.SetDefault("STORAGE_COLLECTION_NAME", "storage")
	viper.SetDefault("INDEX_COLLECTION_NAME", "index")
	viper.SetDefault("LINK_COLLECTION_NAME", "link")

	// Set defaults for RedisConfig
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_ADDRESS", "localhost:6379")
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("REDIS_TIMEOUT", 10) // Assuming this is in seconds

	// Set defaults for ServerConfig
	viper.SetDefault("SERVER_PORT", 40051)
	viper.SetDefault("SERVER_HOST", "localhost")
	viper.SetDefault("INDEX_DIRECTORY", "index")
	viper.SetDefault("SERVER_DATA_STORE_DIRECTORY", "data")
	viper.SetDefault("STORAGE_DIRECTORY", "local")
	viper.SetDefault("SERVICE_NAME", "search")
	viper.SetDefault("SERVICE_VERSION", "1.0.0")

	return &Config{
		MongoDBConfig: MongoDBConfig{
			Host:                  viper.GetString("MONGO_HOST"),
			Port:                  viper.GetInt("MONGO_PORT"),
			Username:              viper.GetString("MONGODB_USERNAME"),
			Password:              viper.GetString("MONGODB_PASSWORD"),
			Database:              viper.GetString("MONGO_DATABASE"),
			MaxPoolSize:           viper.GetUint64("MONGO_MAX_POOL_SIZE"),
			URI:                   getMongoUri(),
			ContainerName:         viper.GetString("MONGO_CONTAINER_NAME"),
			StorageCollectionName: viper.GetString("STORAGE_COLLECTION_NAME"),
			IndexCollectionName:   viper.GetString("INDEX_COLLECTION_NAME"),
			LinkCollectionName:    viper.GetString("LINK_COLLECTION_NAME"),
		},
		RedisConfig: RedisConfig{
			Password:      viper.GetString("REDIS_PASSWORD"),
			Address:       viper.GetString("REDIS_ADDRESS"),
			DB:            viper.GetInt("REDIS_DB"),
			Timeout:       time.Duration(viper.GetInt("REDIS_TIMEOUT")) * time.Second,
			ContainerName: viper.GetString("REDIS_CONTAINER_NAME"),
		},
		ServerConfig: ServerConfig{
			Port:               viper.GetInt("SERVER_PORT"),
			Server:             viper.GetString("SERVER_HOST"),
			IndexDirectory:     viper.GetString("INDEX_DIRECTORY"),
			DataStoreDirectory: viper.GetString("DATA_STORE_DIRECTORY"),
			Storage:            viper.GetString("STORAGE_DIRECTORY"),
			ServiceName:        viper.GetString("SERVICE_NAME"),
			ServiceVersion:     viper.GetString("SERVICE_VERSION"),
		},
	}, nil
}

func getMongoUri() string {
	if uri := viper.GetString("MONGODB_URI"); uri != "" {
		return uri
	}

	host := viper.GetString("MONGO_HOST")
	port := viper.GetInt("MONGO_PORT")
	username := viper.GetString("MONGO_USERNAME")
	password := viper.GetString("MONGO_PASSWORD")
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/?directConnection=true", username, password, host, port)
	return uri
}
