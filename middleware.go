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
	"github.com/pterm/pterm"
	"github.com/rs/zerolog"
)

type UseErrorHandlerConfig struct {
	Writer io.Writer
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
		logger *zerolog.Logger
	)

	if len(config) == 0 {
		zeroLog := zerolog.New(defaultErrorWriter())
		logger = &zeroLog
	} else {
		if config[0].Logger != nil {
			logger = config[0].Logger
		} else if config[0].Writer != nil {
			zerolog := zerolog.New(config[0].Writer)
			logger = &zerolog
		} else {
			zeroLog := zerolog.New(defaultErrorWriter())
			logger = &zeroLog
		}
	}

	startTimeHandler := func(ctx *fiber.Ctx) error {
		requestID := ctx.Get("x-request-id")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx.Locals("xgo_use_logger_requestID", requestID)
		ctx.Locals("xgo_use_logger_startTime", time.Now()) // Store the start time in locals
		logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("request_id", requestID)
		})
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
		if err == nil {

			logger.Trace().Ctx(ctx.UserContext()).
				Str("path", ctx.Path()).
				Str("method", ctx.Method()).
				Str("ip", ctx.IP()).
				Str("latency", latency.String()).
				Msg("success")
			return nil
		}

		xgoError := AsXgoError(err)

		evt := logger.Error().
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

func defaultErrorWriter() io.Writer {
	return utils.JsonWriter{
		Message: pterm.BgRed.Sprint(" ERROR "),
	}
}
