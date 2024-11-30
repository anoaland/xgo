package xgo

import (
	"errors"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/Nerzal/gocloak"
	xgoErrors "github.com/anoaland/xgo/errors"
	"github.com/gofiber/fiber/v2"
)

func NewHttpError(message string, errorCode int) *xgoErrors.XgoError {
	return xgoErrors.NewHttpError("", errors.New(message), errorCode, 1)
}

func NewHttpBadRequestError(statusCode string, err error) *xgoErrors.XgoError {
	return xgoErrors.NewHttpError("", err, fiber.StatusBadRequest, 1)
}

func NewHttpForbiddenError(statusCode string, err error) *xgoErrors.XgoError {
	return xgoErrors.NewHttpError("", err, fiber.StatusForbidden, 1)
}

func NewHttpNotFoundError(statusCode string, err error) *xgoErrors.XgoError {
	return xgoErrors.NewHttpError("", err, fiber.StatusNotFound, 1)
}

func NewHttpInternalError(statusCode string, err error) *xgoErrors.XgoError {
	return xgoErrors.NewHttpError("", err, fiber.StatusInternalServerError, 1)
}

func (s *WebServer) Error(ctx *fiber.Ctx, err error) error {
	return s.FinalError(ctx, err)
}

func (s *WebServer) FinalError(ctx *fiber.Ctx, err error) error {

	var xgoErr *xgoErrors.XgoError
	if errors.As(err, &xgoErr) {
		s.traceXgoError(xgoErr)
		return xgoErr.AsFiberError(ctx)
	}

	if fiberError, ok := err.(*fiber.Error); ok {
		if fiberError.Code < 500 {
			s.traceError("FIBER", fiberError.Code, err)
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

		s.traceError("KEYCLOAK", errorCode, err)
		return ctx.Status(errorCode).JSON(fiber.Map{
			"message":    message,
			"code":       errorCode,
			"statusCode": "AUTH_ERROR",
		})
	}

	s.traceError("FATAL", 500, err)
	return ctx.Status(500).JSON(fiber.Map{
		"message":    "Terjadi kesalahan",
		"code":       500,
		"statusCode": "INTERNAL_SERVER_ERROR",
	})

}

func (s *WebServer) traceError(part string, httpErrorCode int, err error) {
	if s.errorHandler == nil {
		return
	}
	fn := *s.errorHandler

	_, file, line, _ := runtime.Caller(1)
	fn(xgoErrors.XgoError{
		IsFatal:       httpErrorCode >= 500,
		Part:          part,
		Err:           err,
		Message:       err.Error(),
		File:          file,
		Line:          line,
		HttpErrorCode: httpErrorCode,
		Stack:         string(debug.Stack()),
	})
}

func (s *WebServer) traceXgoError(err *xgoErrors.XgoError) {
	if s.errorHandler == nil {
		return
	}
	fn := *s.errorHandler
	fn(*err)
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
