package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/avyukth/search-app/pkg/api/router"
	"github.com/avyukth/search-app/pkg/config"
	"github.com/avyukth/search-app/pkg/database/mongo"
	"github.com/avyukth/search-app/pkg/downloader"
	"github.com/avyukth/search-app/pkg/indexer"
	"github.com/avyukth/search-app/pkg/parser"
	"github.com/avyukth/search-app/pkg/queue"
	"github.com/avyukth/search-app/pkg/worker"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file from current directory
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	// Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configurations: %v", err)
	}

	// Setup Database Connection
	db, err := mongo.SetupDatabase(cfg)
	if err != nil {
		log.Fatalf("Error setting up database: %v", err)
	}
	defer db.Client.Disconnect(context.TODO())

	// Initialize components
	httpClient := &http.Client{Timeout: 30 * time.Second,}
	parser := parser.NewParser()

	indexer, err := indexer.NewSearchEngine(cfg.IndexDirectory)
	if err != nil {
		log.Fatalf("Error loading configurations: %v", err)
	}

	dl := downloader.NewDownloader(httpClient, &cfg.ServerConfig) // Passing actual dependency
	wk := worker.NewWorker(dl, parser, db, indexer)               // Passing actual dependencies
	q := queue.NewTaskQueue(10, wk)                               // Passing actual dependencies


	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	q.Start(ctx)

	
	// Setup Fiber App
	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
	})

	// Setup Router
	router.SetupRoutes(app, db, indexer, q)

	// Start Server
	go func() {
		log.Printf("Starting Server on port %d", cfg.ServerConfig.Port)
		if err := app.Listen(fmt.Sprintf(":%d", cfg.ServerConfig.Port)); err != nil {
			log.Fatalf("Error starting the app: %v", err)
		}
	}()

	// Graceful Shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	_ = <-c
	log.Println("Graceful Shutdown...")
	_ = app.Shutdown()
}

// package main

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

// 		// Load .env file from current directory
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
// 	d := downloader.NewDownloader(client, &cfg.ServerConfig)

// 	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
// 	defer cancel()

// 	_, err = d.Download(ctx, "https://bitly.ws/W7f4")
// 	if err != nil {
// 		log.Fatalf("Failed to download: %v", err)
// 	}
// }
