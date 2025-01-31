package xgo

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/pterm/pterm"
	"github.com/rs/zerolog"
)

type UseErrorHandlerConfig struct {
	Writer *io.Writer
	Logger *zerolog.Logger
}

// UseErrorHandler is a middleware function that provides error handling for a WebServer.
// It sets up three middleware functions:
// 1. startTimeHandler - Stores the start time of the request in the context.
// 2. panicRecoverHandler - Recovers from panics and stores the stack trace in the context.
// 3. errorHandler - Logs errors that occur during the request, including the request details, latency, and stack trace.
// The error handling can be configured by passing a UseErrorHandlerConfig struct, which allows setting a custom logger.
func (server *WebServer) UseErrorHandler(config ...UseErrorHandlerConfig) {
	var (
		logger zerolog.Logger
	)

	if len(config) == 0 {
		logger = zerolog.New(os.Stderr)
	} else {
		if config[0].Logger != nil {
			logger = *config[0].Logger
		} else if config[0].Writer != nil {
			logger = zerolog.New(*config[0].Writer)
		} else {
			logger = zerolog.New(pterm.Error.Writer)
		}
	}

	startTimeHandler := func(c *fiber.Ctx) error {
		c.Locals("xgo_use_error_handler_startTime", time.Now()) // Store the start time in locals
		return c.Next()
	}

	panicRecoverHandler := recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e any) {
			var stack []string
			pcs := make([]uintptr, 32)
			n := runtime.Callers(1, pcs)
			frames := runtime.CallersFrames(pcs[:n])
			for {
				frame, more := frames.Next()
				if !more {
					break
				}
				stack = append(stack, fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function))
			}
			c.Locals("xgo_use_error_handler_stackError", stack)
		},
	})

	errorHandler := func(ctx *fiber.Ctx) error {
		err := ctx.Next()
		if err == nil {
			return nil
		}

		xgoError := AsXgoError(err)
		start := ctx.Locals("xgo_use_error_handler_startTime").(time.Time)
		latency := time.Since(start)

		evt := logger.Error().
			Str("path", ctx.Path()).
			Str("method", ctx.Method()).
			Str("ip", ctx.IP()).
			Int("status", xgoError.HttpErrorCode).
			Str("latency", latency.String()).
			Str("message", xgoError.Message).
			Str("part", xgoError.Part)

		var stack []string
		if stackError := ctx.Locals("xgo_use_error_handler_stackError"); stackError != nil {
			stack = stackError.([]string)
		}

		if xgoError.File != "" {
			stack = append(stack, fmt.Sprintf("%s:%d", xgoError.File, xgoError.Line))
		}

		arr := zerolog.Arr()
		for _, st := range stack {
			arr.Str(st)
		}
		evt.Array("stack", arr)

		evt.Send()

		return nil
	}

	server.App.Use(startTimeHandler)
	server.App.Use(panicRecoverHandler)
	server.App.Use(errorHandler)
}
