package main

import (
	"log"
	"os"

	"github.com/anoaland/xgo"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func main() {
	// Create a new XGO server
	server := xgo.New()

	// Setup logger with the improved factory pattern
	server.UseLogger(xgo.UseLoggerConfig{
		Writer: os.Stdout, // You can also use a custom writer or logger
		// Logger: customZerologLogger, // Or provide your own zerolog.Logger
	})

	// Example handler showing how to use the request logger
	server.App.Get("/api/example", func(ctx *fiber.Ctx) error {
		// Get the request-specific logger (already includes request_id)
		logger := xgo.GetRequestLogger(ctx)

		if logger != nil {
			logger.Info().Msg("Processing example request")

			// You can add more context as needed
			logger.Info().
				Str("user_id", "123").
				Str("action", "fetch_data").
				Msg("User action logged")
		}

		// Get request ID if needed
		requestID := xgo.GetRequestID(ctx)

		return ctx.JSON(fiber.Map{
			"message":    "Success",
			"request_id": requestID,
		})
	})

	// Example of database operation (GORM will automatically use the same request logger)
	server.App.Post("/api/users", func(ctx *fiber.Ctx) error {
		logger := xgo.GetRequestLogger(ctx)

		if logger != nil {
			logger.Info().Msg("Creating new user")
		}

		// Your GORM operations here will automatically log with the same request_id
		// db.Create(&user) // This will log SQL with the correct request_id

		// Use the helper function for consistent logging
		xgo.LogWithContext(ctx, zerolog.InfoLevel, "User created successfully")

		return ctx.JSON(fiber.Map{"status": "created"})
	})

	// Start server
	log.Println("Starting server on :3000")
	server.Run(3000, func() error {
		return nil
	})
}

/*
Example of how logs will now look:

Before (with duplicated request_id):
{
  "level": "info",
  "request_id": "uuid1",
  "request_id": "uuid2",
  "request_id": "uuid3",
  "message": "success"
}

After (clean single request_id):
{
  "level": "info",
  "request_id": "uuid1",
  "path": "/api/example",
  "method": "GET",
  "message": "success"
}
*/
