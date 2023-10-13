package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/avyukth/search-app/docs"
	appTrace "github.com/avyukth/search-app/foundations/tracing"
	"github.com/avyukth/search-app/pkg/api/router"
	"github.com/avyukth/search-app/pkg/config"
	"github.com/avyukth/search-app/pkg/database/mongo"
	"github.com/avyukth/search-app/pkg/downloader"
	"github.com/avyukth/search-app/pkg/indexer"
	"github.com/avyukth/search-app/pkg/parser"
	"github.com/avyukth/search-app/pkg/queue"
	"github.com/avyukth/search-app/pkg/worker"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
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


	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	otelShutdown, err := appTrace.SetupOTelSDK(ctx, cfg.ServerConfig.ServiceName, cfg.ServerConfig.ServiceVersion)

	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// Setup Database Connection
	db, err := mongo.SetupDatabase(cfg)
	if err != nil {
		log.Fatalf("Error setting up database: %v", err)
	}
	defer db.Client.Disconnect(context.TODO())

	// Initialize components
	httpClient := &http.Client{Timeout: 30 * time.Second,}
	parser := parser.NewParser()

	indexer, err := indexer.NewSearchEngine(cfg.ServerConfig.Storage + cfg.ServerConfig.IndexDirectory)
	if err != nil {
		log.Fatalf("Error loading configurations: %v", err)
	}

	dl := downloader.NewDownloader(httpClient, &cfg.ServerConfig)
	wk := worker.NewWorker(dl, parser, db, indexer)
	q := queue.NewTaskQueue(10, wk)
	swaggerURL := fmt.Sprintf("http://127.0.0.1:%d/docs/swagger.yaml", cfg.ServerConfig.Port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	q.Start(ctx)

	//testing
	// Setup Fiber App
	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
	})

	app.Static("/docs", "./docs")
	app.Use(cors.New())

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: swaggerURL,
	}))

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

