package xgo

import (
	"fmt"
	"io"
	"runtime"
	"strings"
	"time"

	"github.com/anoaland/xgo/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type UseLoggerConfig struct {
	Writer io.Writer
	Logger *zerolog.Logger
}

// LoggerFactory creates a new logger instance for each request
type LoggerFactory struct {
	baseWriter io.Writer
	baseLogger *zerolog.Logger
}

func NewLoggerFactory(config ...UseLoggerConfig) *LoggerFactory {
	factory := &LoggerFactory{}

	if len(config) > 0 {
		if config[0].Logger != nil {
			factory.baseLogger = config[0].Logger
		} else if config[0].Writer != nil {
			factory.baseWriter = config[0].Writer
		} else {
			factory.baseWriter = DefaultLogWriter()
		}
	} else {
		factory.baseWriter = DefaultLogWriter()
	}

	return factory
}

func (f *LoggerFactory) CreateRequestLogger(requestID string) *zerolog.Logger {
	var baseLogger zerolog.Logger

	if f.baseLogger != nil {
		// Create a copy of the base logger to avoid shared state
		baseLogger = f.baseLogger.With().Logger()
	} else {
		baseLogger = zerolog.New(f.baseWriter).With().Timestamp().Logger()
	}

	// Create request-specific logger with request_id
	requestLogger := baseLogger.With().Str("request_id", requestID).Logger()
	return &requestLogger
}

// UseLogger is a middleware function that provides error handling and logging for a WebServer.
// It sets up three middleware functions:
//  1. startTimeHandler - Stores the start time of the request and generates a unique request ID if not provided.
//     The request ID is stored in the context for tracking purposes.
//  2. panicRecoverHandler - Recovers from panics and stores the stack trace in the context.
//  3. errorHandler - Logs errors that occur during the request, including the request details, latency, and stack trace.
//     It also logs successful requests with their details.
//
// The error handling can be configured by passing a UseLoggerConfig struct, which allows setting a custom logger or writer.
//
// The request ID is used to uniquely identify each request, which helps in tracking and debugging issues across different parts of the system.
//
// Activity logs are generated for both successful and failed requests. For successful requests, the log includes the request path, method, IP, and latency.
// For failed requests, the log includes additional details such as the error message, HTTP status code, and stack trace.
//
// Parameters:
// - config: Optional configuration for the logger, allowing customization of the logger or writer.
func (server *WebServer) UseLogger(config ...UseLoggerConfig) {
	// Create logger factory to avoid global logger issues
	loggerFactory := NewLoggerFactory(config...)

	startTimeHandler := func(ctx *fiber.Ctx) error {
		requestID := ctx.Get("x-request-id")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx.Locals("xgo_use_logger_requestID", requestID)
		ctx.Locals("xgo_use_logger_startTime", time.Now()) // Store the start time in locals

		// Create a fresh logger instance for this request - no shared state
		requestLogger := loggerFactory.CreateRequestLogger(requestID)
		ctx.Locals("xgo_request_logger", requestLogger)
		return ctx.Next()
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
			c.Locals("xgo_use_logger_stackError", stack)
		},
	})

	errorHandler := func(ctx *fiber.Ctx) error {
		err := ctx.Next()
		start := ctx.Locals("xgo_use_logger_startTime").(time.Time)
		latency := time.Since(start)

		// Get the per-request logger
		requestLogger := ctx.Locals("xgo_request_logger").(*zerolog.Logger)

		if err == nil {
			requestLogger.Info().Ctx(ctx.UserContext()).
				Str("path", ctx.Path()).
				Str("method", ctx.Method()).
				Str("ip", ctx.IP()).
				Str("latency", latency.String()).
				Msg("success")
			return nil
		}

		xgoError := AsXgoError(err)

		evt := requestLogger.Error().
			Str("path", ctx.Path()).
			Str("method", ctx.Method()).
			Str("ip", ctx.IP()).
			Int("status", xgoError.HttpErrorCode).
			Str("latency", latency.String()).
			Str("message", xgoError.Message).
			Str("part", xgoError.Part)

		var stack []string
		if stackError := ctx.Locals("xgo_use_logger_stackError"); stackError != nil {
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

		if xgoError.Stack == "" {
			xgoError.Stack = strings.Join(stack, "\n")
		}

		return xgoError
	}

	server.App.Use(startTimeHandler)
	server.App.Use(errorHandler)
	server.App.Use(panicRecoverHandler)
}

// GetRequestLogger retrieves the request-specific logger from the Fiber context.
// This logger already includes the request_id in its context.
// Returns nil if no logger is found (middleware not properly set up).
func GetRequestLogger(ctx *fiber.Ctx) *zerolog.Logger {
	logger := ctx.Locals("xgo_request_logger")
	if logger == nil {
		return nil
	}
	return logger.(*zerolog.Logger)
}

// GetRequestID retrieves the request ID from the Fiber context.
// Returns empty string if no request ID is found.
func GetRequestID(ctx *fiber.Ctx) string {
	requestID := ctx.Locals("xgo_use_logger_requestID")
	if requestID == nil {
		return ""
	}
	return requestID.(string)
}

// LogWithContext is a helper function that logs with the request context if available.
// If no request logger is found, it falls back to a default logger.
// This is useful for logging outside of HTTP handlers.
func LogWithContext(ctx *fiber.Ctx, level zerolog.Level, msg string) {
	logger := GetRequestLogger(ctx)
	if logger == nil {
		// Fallback to default logger if no request logger available
		defaultLogger := zerolog.New(DefaultLogWriter()).With().Timestamp().Logger()
		logger = &defaultLogger
	}

	var event *zerolog.Event
	switch level {
	case zerolog.InfoLevel:
		event = logger.Info()
	case zerolog.WarnLevel:
		event = logger.Warn()
	case zerolog.ErrorLevel:
		event = logger.Error()
	case zerolog.DebugLevel:
		event = logger.Debug()
	default:
		event = logger.Info()
	}

	event.Msg(msg)
}

func DefaultLogWriter() io.Writer {
	return utils.JsonWriter{}
}
