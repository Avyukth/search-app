package handler

import (
	"net/http"

	"github.com/avyukth/search-app/pkg/database/mongo"
	"github.com/avyukth/search-app/pkg/indexer"
	"github.com/avyukth/search-app/pkg/queue"
	"github.com/gofiber/fiber/v2"
)

// DownloadHandler handles download requests for tar files
func DownloadHandler(db *mongo.Database, q *queue.TaskQueue) fiber.Handler {
	return func(c *fiber.Ctx) error {
		link := c.Query("link")
		if link == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Link is required")
		}

		// Check if the link is live
		resp, err := http.Head(link)
		if err != nil || resp.StatusCode != http.StatusOK {
			return c.Status(fiber.StatusNotFound).SendString("Link is not live")
		}

		// Check and set link status in MongoDB
		isSet, err := db.CheckAndSetLinkStatus(link)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error processing link")
		}
		if !isSet {
			return c.Status(fiber.StatusConflict).SendString("Link is already processed or completed")
		}

		task := queue.Task{
			FilePath: link,
			Type:     queue.DownloadAndProcess,
		}
		// Send link to processing queue and return response
		q.Enqueue(task)
		return c.SendString("Link is sent for processing")
	}
}

func SearchHandler(db *mongo.Database, searchEngine *indexer.SearchEngine) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract search parameters from the request
		query := c.Query("query")

		// Perform search operation using the search engine instance
		results, err := searchEngine.SearchAndRetrievePatents(query)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(results)
	}
}

func CrawlerHandler(db *mongo.Database, q *queue.TaskQueue) fiber.Handler {
	return func(c *fiber.Ctx) error {
		dirPath := c.Query("path")
		if dirPath == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Path is required")
		}

		task := queue.Task{
			FilePath: dirPath,
			Type:     queue.WalkAndProcess,
		}
		q.Enqueue(task)
		return c.SendString("Directory is sent for walking and processing")
	}
}
