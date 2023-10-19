package main

import (
	"context"
	// "errors"
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
	// "github.com/joho/godotenv"
)

type AppConfig struct {
	ListenAddr   string
	ServiceName string
	ServiceVersion string
	FiberConfig  fiber.Config
}


// func main() {
// 	// Load Configuration
// 	cfg, err := config.LoadConfig()
// 	if err != nil {
// 		log.Fatalf("Error loading configurations: %v", err)
// 	}


// 	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
// 	defer stop()

// 	otelShutdown, err := appTrace.SetupOTelSDK(ctx, cfg.ServerConfig.ServiceName, cfg.ServerConfig.ServiceVersion)

// 	defer func() {
// 		err = errors.Join(err, otelShutdown(context.Background()))
// 	}()

// 	// Setup Database Connection
// 	db, err := mongo.SetupDatabase(cfg)
// 	if err != nil {
// 		log.Fatalf("Error setting up database: %v", err)
// 	}
// 	defer db.Client.Disconnect(context.TODO())

// 	// Initialize components
// 	httpClient := &http.Client{Timeout: 30 * time.Second,}
// 	parser := parser.NewParser()

// 	indexer, err := indexer.NewSearchEngine(cfg.ServerConfig.Storage + cfg.ServerConfig.IndexDirectory)
// 	if err != nil {
// 		log.Fatalf("Error loading configurations: %v", err)
// 	}

// 	dl := downloader.NewDownloader(httpClient, &cfg.ServerConfig)
// 	wk := worker.NewWorker(dl, parser, db, indexer)
// 	q := queue.NewTaskQueue(10, wk)
// 	swaggerURL := fmt.Sprintf("/docs/swagger.yaml")

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
// 	q.Start(ctx)

// 	//testing
// 	// Setup Fiber App
// 	app := fiber.New(fiber.Config{
// 		Prefork:       false,
// 		CaseSensitive: true,
// 		StrictRouting: true,
// 	})

// 	app.Static("/docs", "./docs")
// 	app.Use(cors.New())

// 	app.Get("/swagger/*", swagger.New(swagger.Config{
// 		URL: swaggerURL,
// 	}))

// 	// Setup Router
// 	router.SetupRoutes(app, db, indexer, q)


// 	appConfig :=  AppConfig{
// 		ListenAddr: fmt.Sprintf(":%d", cfg.ServerConfig.ServicePort),
// 		ServiceName: cfg.ServerConfig.ServiceName,
// 		ServiceVersion: cfg.ServerConfig.ServiceVersion,
// 		FiberConfig: fiber.Config{
// 			Prefork:               false,
// 			ServerHeader:          "Fiber",
// 			ReadTimeout:           time.Second,
// 			WriteTimeout:          10 * time.Second,
// 			DisableStartupMessage: false,
// 		},
// 	}


// 	go startApp(app, appConfig.ListenAddr)

// 	// Wait for interrupt signal to gracefully shut down the server.
// 	waitForShutdownSignal(app)

// }
func main() {
	cfg := loadConfig()
	ctx := setupContext()

	otelShutdown := setupTracing(ctx, cfg)
	defer shutdownTracing(ctx, otelShutdown)

	db := setupDatabase(cfg)
	defer db.Client.Disconnect(ctx)

	httpClient, parser, indexer := initializeComponents(cfg)
	q := setupWorkerComponents(httpClient, parser, db, indexer, cfg)
	defer q.Stop()

	app := setupFiberApp(cfg)
	router.SetupRoutes(app, db, indexer, q)


	go startApp(app, cfg.ServerConfig)
	waitForShutdownSignal(app)
}

func loadConfig() *config.Config {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configurations: %v", err)
	}
	return cfg
}

func setupContext() context.Context {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	return ctx
}

func setupTracing(ctx context.Context, cfg *config.Config) func(context.Context) error {
	otelShutdown, err := appTrace.SetupOTelSDK(ctx, cfg.ServerConfig.ServiceName, cfg.ServerConfig.ServiceVersion)
	if err != nil {
		log.Fatalf("Error setting up tracing: %v", err)
	}
	return otelShutdown
}

func shutdownTracing(ctx context.Context, otelShutdown func(context.Context) error) {
	if err := otelShutdown(ctx); err != nil {
		log.Fatalf("Error shutting down tracing: %v", err)
	}
}

func setupDatabase(cfg *config.Config) *mongo.Database {
	db, err := mongo.SetupDatabase(cfg)
	if err != nil {
		log.Fatalf("Error setting up database: %v", err)
	}
	return db
}

func initializeComponents(cfg *config.Config) (*http.Client, *parser.Parser, *indexer.SearchEngine) {
	httpClient := &http.Client{Timeout: 30 * time.Second}
	parser := parser.NewParser()
	indexer, err := indexer.NewSearchEngine(cfg.ServerConfig.Storage + cfg.ServerConfig.IndexDirectory)
	if err != nil {
		log.Fatalf("Error initializing indexer: %v", err)
	}
	return httpClient, parser, indexer
}

func setupWorkerComponents(httpClient *http.Client, parser *parser.Parser, db *mongo.Database, indexer *indexer.SearchEngine, cfg *config.Config) (*queue.TaskQueue) {
	dl := downloader.NewDownloader(httpClient, &cfg.ServerConfig)
	wk := worker.NewWorker(dl, parser, db, indexer)
	q := queue.NewTaskQueue(10, wk)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	q.Start(ctx)
	return  q
}

func setupFiberApp(cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:               false,
		CaseSensitive:         true,
		StrictRouting:         true,
		ServerHeader:          "Fiber",
		ReadTimeout:           time.Second,
		WriteTimeout:          10 * time.Second,
		DisableStartupMessage: false,
	})
	app.Static("/docs", "./docs")
	app.Use(cors.New())
	swaggerURL := "/docs/swagger.yaml"
	app.Get("/swagger/*", swagger.New(swagger.Config{URL: swaggerURL}))
	return app
}


func waitForShutdownSignal(app *fiber.App) {
	// Create a channel to listen for interrupt or terminate signals.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal.
	<-c

	// Attempt to gracefully shut down the server.
	log.Println("Gracefully shutting down...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}
	log.Println("Server shutdown complete")
}


func startApp(app *fiber.App, cfg config.ServerConfig) {
	addr:= fmt.Sprintf(":%d", cfg.ServicePort)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
