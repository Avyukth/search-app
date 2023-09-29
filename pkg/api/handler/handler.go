package handler

import (
	"log"
	"net/http"

	"github.com/avyukth/search-app/pkg/database/mongo"
	"github.com/gofiber/fiber/v2"
)

// DownloadHandler handles download requests for tar files
func DownloadHandler(db *mongo.Database) fiber.Handler {
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

		// Send link to processing queue and return response
		// TODO: Send link to processing queue
		return c.SendString("Link is sent for processing")
	}
}

func SearchHandler(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract search parameters from the request
		query := c.Query("query")
		log.Println("******************Query", query)
		// Perform search operation using the database instance
		// results, err := db.Search(query)
		// if err != nil {
		// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
		// }
		return c.JSON("Success")
	}
}
