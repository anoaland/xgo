package xgo

import (
	"github.com/gofiber/fiber/v2"
)

// Deprecated: Use Response() function directly instead.
func (server *WebServer) Response(ctx *fiber.Ctx, response any, successCode int, err error) error {
	return Response(ctx, response, successCode, err)
}

// Response sends a JSON response with the given success code if there is no error.
// If an error is provided, it returns the error instead.
//
// Parameters:
//   - ctx: The Fiber context to send the response to.
//   - response: The response data to be sent as JSON.
//   - successCode: The HTTP status code to be used for the response if there is no error.
//   - err: An error that, if not nil, will be returned instead of sending the response.
//
// Returns:
//   - error: The provided error if it is not nil, otherwise nil.
func Response(ctx *fiber.Ctx, response any, successCode int, err error) error {
	if err != nil {
		return err
	}

	switch v := response.(type) {
	case nil:
		return ctx.Status(successCode).JSON(fiber.Map{
			"data": nil,
		})
	case string, int, float64, bool:
		return ctx.Status(successCode).JSON(fiber.Map{
			"data": v,
		})
	}

	return ctx.Status(successCode).JSON(response)
}

type DefaultErrorHandlerConfig struct {
	FatalErrorMessage string
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

	if cfg.FatalErrorMessage == "" {
		cfg.FatalErrorMessage = "Something went wrong"
	}

	return func(ctx *fiber.Ctx, err error) error {
		xgoError := AsXgoError(err)
		return xgoError.FiberJsonResponse(ctx, cfg.FatalErrorMessage)
	}
}
