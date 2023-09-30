package router

import (
	"github.com/avyukth/search-app/pkg/api/handler"
	"github.com/avyukth/search-app/pkg/database/mongo"
	"github.com/avyukth/search-app/pkg/indexer"
	"github.com/avyukth/search-app/pkg/queue"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes sets up all the routes for your application
func SetupRoutes(app *fiber.App, db *mongo.Database, searchEngine *indexer.SearchEngine, q *queue.TaskQueue) {

	// Middleware
	app.Use(logger.New()) // add logging to each request

	// Routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Define your routes here
	v1.Get("/search", handler.SearchHandler(db, searchEngine))
	v1.Get("/download", handler.DownloadHandler(db, q))
}
