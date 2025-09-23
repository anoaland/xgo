# XGO

⚠️ **WARNING: This library is in an early stage of development.**  
It is **not stable**, lacks unit tests, and has minimal documentation.  
Use at your own risk and expect breaking changes in future updates.

![Status: Experimental](https://img.shields.io/badge/status-experimental-orange)

## Logger

XGO provides an integrated logging system built on [zerolog](https://github.com/rs/zerolog) with automatic request tracing and GORM integration.

### Features

- **Request Tracing**: Automatic `request_id` generation and propagation across all logs
- **GORM Integration**: Database queries automatically include the same `request_id`
- **Per-Request Isolation**: Each request gets its own logger instance (no shared state)
- **Zero Configuration**: Works out of the box with sensible defaults
- **Memory Efficient**: Proper cleanup of per-request resources

### Quick Start

```go
package main

import (
    "github.com/anoaland/xgo"
    "github.com/gofiber/fiber/v2"
)

func main() {
    server := xgo.New()

    // Enable logging with default configuration
    server.UseLogger()

    server.App.Get("/api/users", func(ctx *fiber.Ctx) error {
        // Get request-specific logger
        logger := xgo.GetRequestLogger(ctx)
        logger.Info().Msg("Processing user request")

        return ctx.JSON(fiber.Map{"users": []string{}})
    })

    server.Run(3000, func() error { return nil })
}
```

### Configuration

```go
server.UseLogger(xgo.UseLoggerConfig{
    Writer: os.Stdout,              // Custom writer
    Logger: customZerologLogger,    // Custom zerolog instance
})
```

### GORM Integration

XGO automatically integrates with GORM to provide consistent request tracing across application and database logs:

```go
type Service struct {
    db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
    return &Service{db: db}
}

// In your middleware/route setup
router.Use(func(ctx *fiber.Ctx) error {
    service := NewService(
        db.WithContext(server.LoggerContext(ctx)), // Enables request tracing
    )
    // Your service now has request-traced database operations
    return ctx.Next()
})
```

### Logging Helpers

XGO provides convenient helpers for logging within handlers:

```go
func handler(ctx *fiber.Ctx) error {
    // Get request logger (includes request_id automatically)
    logger := xgo.GetRequestLogger(ctx)
    logger.Info().Str("action", "fetch_data").Msg("User action")

    // Get request ID if needed
    requestID := xgo.GetRequestID(ctx)

    // Helper for consistent logging
    xgo.LogWithContext(ctx, zerolog.InfoLevel, "Operation completed")

    return ctx.JSON(fiber.Map{"request_id": requestID})
}
```

### Log Format

All logs include structured data with consistent request tracing:

```json
{
  "level": "info",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "path": "/api/users",
  "method": "GET",
  "ip": "127.0.0.1",
  "latency": "45.123ms",
  "message": "success"
}
```

Database operations automatically include the same `request_id`:

```json
{
  "level": "info",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "sql": "SELECT * FROM users WHERE active = ?",
  "rows": 5,
  "latency": "2.45ms",
  "message": "SQL query"
}
```

### Architecture

XGO uses a **Logger Factory Pattern** to ensure clean separation between requests:

- **No Global State**: Each request gets a fresh logger instance
- **Request Isolation**: No shared mutable state between concurrent requests
- **Memory Safe**: Automatic cleanup and garbage collection of per-request loggers
- **Thread Safe**: Complete isolation prevents race conditions
