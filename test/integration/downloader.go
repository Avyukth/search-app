package main

// import (
// 	"context"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/avyukth/search-app/pkg/config"
// 	"github.com/avyukth/search-app/pkg/downloader"
// 	"github.com/joho/godotenv"
// )

// func main() {
// // 		// Load .env file from current directory
// 	if err := godotenv.Load(); err != nil {
// 		log.Println("Error loading .env file")
// 	}

// 	// Load Configuration
// 	cfg, err := config.LoadConfig()
// 	if err != nil {
// 		log.Fatalf("Error loading configurations: %v", err)
// 	}

// 	client := &http.Client{
// 		Timeout: 30 * time.Second, // Example timeout
// 	}

// 	// Initialize Downloader
// 	d := downloader.NewDownloader(client, &cfg.ServerConfig)

// 	// Set up context with timeout
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
// 	defer cancel()

// 	// Test Download
// 	link := "https://bitly.ws/W7f4"
// 	filePath, err := d.Download(ctx, link)
// 	if err != nil {
// 		log.Fatalf("Failed to download: %v", err)
// 	}
// 	log.Printf("Successfully downloaded file to: %s", filePath)

// 	// Test ExtractTarGz
// 	extractedPath, err := d.ExtractTarGz(filePath)
// 	if err != nil {
// 		log.Fatalf("Failed to extract: %v", err)
// 	}
// 	log.Printf("Successfully extracted files to: %s", extractedPath)
// }
