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

	// logger Middleware
	app.Use(logger.New())

	// Routes
	api := app.Group("/api")
	v1 := api.Group("/v1")
	//go:generate swagger generate spec -o swagger.json
	v1.Get("/search", handler.SearchHandler(db, searchEngine))
	v1.Get("/download", handler.DownloadHandler(db, q))
	v1.Get("/crawl", handler.CrawlerHandler(db, q))
	v1.Get("/live", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
}
