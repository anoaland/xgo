package xgo

import (
	"github.com/gofiber/fiber/v2"
)

// func (s *WebServer) traceError(part string, httpErrorCode int, err error) {
// 	if s.errorHandler == nil {
// 		return
// 	}
// 	fn := *s.errorHandler

// 	_, file, line, _ := runtime.Caller(4)
// 	fn(xgoErrors.XgoError{
// 		IsFatal:       httpErrorCode >= 500,
// 		Part:          part,
// 		Err:           err,
// 		Message:       err.Error(),
// 		File:          file,
// 		Line:          line,
// 		HttpErrorCode: httpErrorCode,
// 		Stack:         string(debug.Stack()),
// 	})
// }

// func (s *WebServer) traceXgoError(err *xgoErrors.XgoError) {
// 	if s.errorHandler == nil {
// 		return
// 	}
// 	fn := *s.errorHandler
// 	fn(*err)
// }

func (server *WebServer) Response(ctx *fiber.Ctx, response interface{}, successCode int, err error) error {
	if err != nil {
		return err
	}

	return ctx.Status(successCode).JSON(response)
}

type DefaultErrorHandlerConfig struct {
	fatalErrorMessage string
}

// DefaultErrorHandler returns a fiber.ErrorHandler that handles errors by converting them
// to an XgoError and sending a JSON response with a specified fatal error message.
// The function accepts an optional configuration parameter of type DefaultErrorHandlerConfig.
// If no configuration is provided, a default configuration is used.
//
// Parameters:
//   - config: Optional. A variadic parameter of DefaultErrorHandlerConfig. If provided, the first
//     element is used as the configuration.
//
// Returns:
//   - fiber.ErrorHandler: A function that handles errors by converting them to an XgoError and
//     sending a JSON response with the specified fatal error message.
//
// Example usage:
//
//	app := fiber.New()
//	app.Use(DefaultErrorHandler(DefaultErrorHandlerConfig{
//	    fatalErrorMessage: "Custom error message",
//	}))
func DefaultErrorHandler(config ...DefaultErrorHandlerConfig) fiber.ErrorHandler {
	var cfg DefaultErrorHandlerConfig
	if len(config) > 0 {
		cfg = config[0]
	} else {
		cfg = DefaultErrorHandlerConfig{}
	}

	if cfg.fatalErrorMessage == "" {
		cfg.fatalErrorMessage = "Something went wrong"
	}

	return func(ctx *fiber.Ctx, err error) error {
		xgoError := AsXgoError(err)
		return xgoError.FiberJsonResponse(ctx, cfg.fatalErrorMessage)
	}
}
