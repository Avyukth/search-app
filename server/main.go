package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/avyukth/search-app/pkg/api/router"
	"github.com/avyukth/search-app/pkg/config"
	"github.com/avyukth/search-app/pkg/database/mongo"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
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

	// Setup Fiber App
	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
	})

	// Middleware
	app.Use(logger.New())

	// Setup Router
	router.SetupRouter(app, db)

	// Start Server
	go func() {
		log.Printf("Starting Server on port %d", cfg.ServerPort)
		if err := app.Listen(fmt.Sprintf(":%d", cfg.ServerPort)); err != nil {
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
