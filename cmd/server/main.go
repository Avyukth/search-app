package main

import (
	"context"
	// "sync"

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
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	// "github.com/ardanlabs/conf"
)


var build =  "development"


// func run() {
// 	cfg := struct {
// 			conf.Version
// 			Web struct {
// 				APIHost         string        `conf:"default:0.0.0.0:3000"`
// 				DebugHost       string        `conf:"default:0.0.0.0:4000"`
// 				ReadTimeout     time.Duration `conf:"default:5s"`
// 				WriteTimeout    time.Duration `conf:"default:10s"`
// 				IdleTimeout     time.Duration `conf:"default:120s"`
// 				ShutdownTimeout time.Duration `conf:"default:20s"`
// 			}
// 		}{
// 			Version: conf.Version{
// 				SVN: build,
// 				Desc:  "This is sales-api part of developing full production ready backend infrastructure Author: @Avyukth",
// 			},
// 		}
// }

func main() {
	// var wg sync.WaitGroup
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	cfg := loadConfig()
	ctx, stop := setupContext()
	defer stop()
	otelShutdown := setupTracing(ctx, cfg)
	defer shutdownTracing(ctx, otelShutdown)

	db := setupDatabase(cfg)
	defer db.Client.Disconnect(ctx)

	httpClient, parser, indexer := initializeComponents(cfg)
	q := setupWorkerComponents(ctx, httpClient, parser, db, indexer, cfg)
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

func setupContext() (context.Context, context.CancelFunc) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	return ctx, stop
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

func setupWorkerComponents(ctx context.Context, httpClient *http.Client, parser *parser.Parser, db *mongo.Database, indexer *indexer.SearchEngine, cfg *config.Config) (*queue.TaskQueue) {
	dl := downloader.NewDownloader(httpClient, &cfg.ServerConfig)
	wk := worker.NewWorker(dl, parser, db, indexer)
	q := queue.NewTaskQueue(10, wk)
	q.Start(ctx)
	return  q
}

func setupFiberApp(cfg *config.Config) *fiber.App {
	engine := html.New("./front-ui/templates", ".html")
	app := fiber.New(fiber.Config{
		Views: 				engine,
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
