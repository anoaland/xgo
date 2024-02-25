package xgo

import (
	"errors"
	"strings"

	"github.com/Nerzal/gocloak"
	"github.com/gofiber/fiber/v2"
)

type HttpError struct {
	ErrorCode int
	Message   string
}

func NewHttpError(message string, errorCode int) *HttpError {
	return &HttpError{Message: message, ErrorCode: errorCode}
}

func (e *HttpError) Error() string {
	return e.Message
}

func (server *WebServer) Error(ctx *fiber.Ctx, err error) error {
	var httpErr *HttpError
	if errors.As(err, &httpErr) {
		return ctx.Status(httpErr.ErrorCode).SendString(err.Error())
	}

	var apiErr *gocloak.APIError
	if errors.As(err, &apiErr) {

		// parse something like '401 Unauthorized: invalid_grant: Invalid user credentials'
		parts := strings.Split(apiErr.Message, ":")
		message := strings.TrimSpace(parts[len(parts)-1])
		errorCode := apiErr.Code

		if errorCode == 0 {
			errorCode = 500
		}

		return ctx.Status(errorCode).SendString(message)
	}

	return ctx.Status(500).SendString(err.Error())
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
