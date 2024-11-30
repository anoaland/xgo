package xgo

import (
	"errors"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/Nerzal/gocloak"
	"github.com/gofiber/fiber/v2"
)

type HttpError struct {
	ErrorCode     int     `json:"code"`
	Message       string  `json:"message"`
	StatusCode    *string `json:"statusCode"`
	InternalError *error  `json:"internalError"`
}

func NewHttpError(message string, errorCode int) *HttpError {
	return &HttpError{Message: message, ErrorCode: errorCode}
}

func NewHttpBadRequestError(statusCode string, err error) *HttpError {
	return &HttpError{Message: err.Error(), ErrorCode: fiber.StatusBadRequest, InternalError: &err, StatusCode: &statusCode}
}

func NewHttpForbiddenError(statusCode string, err error) *HttpError {
	return &HttpError{Message: err.Error(), ErrorCode: fiber.StatusForbidden, InternalError: &err, StatusCode: &statusCode}
}

func NewHttpNotFoundError(statusCode string, err error) *HttpError {
	return &HttpError{Message: err.Error(), ErrorCode: fiber.StatusNotFound, InternalError: &err, StatusCode: &statusCode}
}

func NewHttpInternalError(statusCode string, err error) *HttpError {
	return &HttpError{Message: "Terjadi kesalahan", ErrorCode: fiber.StatusInternalServerError, InternalError: &err, StatusCode: &statusCode}
}

func (e *HttpError) Error() string {
	return e.Message
}

func (s *WebServer) Error(ctx *fiber.Ctx, err error) error {
	return s.FinalError(ctx, err)
}

func (s *WebServer) FinalError(ctx *fiber.Ctx, err error) error {
	_, file, line, _ := runtime.Caller(1)
	if s.errorHandler != nil {
		fn := *s.errorHandler
		fn(WebTraceError{
			Error: err,
			File:  file,
			Line:  line,
			Stack: string(debug.Stack()),
		})
	}

	var httpErr *HttpError
	if errors.As(err, &httpErr) {

		return ctx.Status(httpErr.ErrorCode).JSON(fiber.Map{
			"message":    httpErr.Message,
			"code":       httpErr.ErrorCode,
			"statusCode": httpErr.StatusCode,
		})
	}

	if fiberError, ok := err.(*fiber.Error); ok {
		if fiberError.Code < 500 {
			return ctx.Status(fiberError.Code).JSON(fiber.Map{
				"message":    fiberError.Message,
				"code":       fiberError.Code,
				"statusCode": "INTERNAL_HTTP_ERROR",
			})
		}
	}

	var goCloakErr *gocloak.APIError
	if errors.As(err, &goCloakErr) {

		// parse something like '401 Unauthorized: invalid_grant: Invalid user credentials'
		parts := strings.Split(goCloakErr.Message, ":")
		message := strings.TrimSpace(parts[len(parts)-1])
		errorCode := goCloakErr.Code

		if errorCode == 0 {
			// TODO: Connect sentry on this line
			errorCode = 500
		}

		return ctx.Status(errorCode).JSON(fiber.Map{

			"message":    message,
			"code":       errorCode,
			"statusCode": "AUTH_ERROR",
		})
	}

	return ctx.Status(500).JSON(fiber.Map{
		"message":    "Terjadi kesalahan",
		"code":       500,
		"statusCode": "INTERNAL_SERVER_ERROR",
	})

}

func (server *WebServer) Response(ctx *fiber.Ctx, response interface{}, successCode int, err error) error {
	if err != nil {
		return server.Error(ctx, err)
	}

	return ctx.Status(successCode).JSON(response)
}

func (server *WebServer) BadRequest(ctx *fiber.Ctx, message string) error {
	return ctx.Status(fiber.StatusBadRequest).SendString(message)
}

func (server *WebServer) InvalidPayload(ctx *fiber.Ctx) error {
	return server.BadRequest(ctx, "Invalid request payload")
}

func (server *WebServer) InvalidParameters(ctx *fiber.Ctx) error {
	return server.BadRequest(ctx, "Invalid request parameters")
}

func (server *WebServer) Status(ctx *fiber.Ctx, successCode int, err error) error {
	if err != nil {
		return server.Error(ctx, err)
	}

	return ctx.SendStatus(successCode)
}
