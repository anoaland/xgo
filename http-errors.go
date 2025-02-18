package xgo

import (
	"errors"
	"strings"

	"github.com/Nerzal/gocloak"
	xgoErrors "github.com/anoaland/xgo/errors"
	"github.com/gofiber/fiber/v2"
)

func NewHttpBadRequestError(part string, err error) *xgoErrors.XgoError {
	return xgoErrors.NewHttpError(part, err, fiber.StatusBadRequest, 1)
}

func NewHttpForbiddenError(part string, err error) *xgoErrors.XgoError {
	return xgoErrors.NewHttpError(part, err, fiber.StatusForbidden, 1)
}

func NewHttpUnauthorizedError(part string, err error) *xgoErrors.XgoError {
	return xgoErrors.NewHttpError(part, err, fiber.StatusUnauthorized, 1)
}

func NewHttpNotFoundError(part string, err error) *xgoErrors.XgoError {
	return xgoErrors.NewHttpError(part, err, fiber.StatusNotFound, 1)
}

func NewHttpInternalError(part string, err error) *xgoErrors.XgoError {
	return xgoErrors.NewHttpError(part, err, fiber.StatusInternalServerError, 1)
}

func NewHttpBadGatewayError(part string, err error) *xgoErrors.XgoError {
	return xgoErrors.NewHttpError(part, err, fiber.StatusBadGateway, 1)
}

func NewHttpCustomError(part string, httpStatusCode int, err error) *xgoErrors.XgoError {
	return xgoErrors.NewHttpError(part, err, httpStatusCode, 1)
}

// AsXgoError converts a given error into an XgoError. It attempts to match the error
// to known error types and returns a corresponding XgoError. If the error is of type
// *fiber.Error, it creates a new XgoError with the "FIBER" category. If the error is
// of type *gocloak.APIError, it parses the error message and creates a new XgoError
// with the "AUTH_ERROR" category and appropriate HTTP status code. If the error does
// not match any known types, it returns a general XgoError.
//
// Parameters:
//   - err: The error to be converted.
//
// Returns:
//   - *xgoErrors.XgoError: The converted XgoError.
func AsXgoError(err error) *xgoErrors.XgoError {
	var xgoErr *xgoErrors.XgoError
	if errors.As(err, &xgoErr) {
		return xgoErr
	}

	if fiberError, ok := err.(*fiber.Error); ok {
		return &xgoErrors.XgoError{
			Message:       fiberError.Message,
			IsFatal:       fiberError.Code == fiber.StatusInternalServerError,
			HttpErrorCode: fiberError.Code,
			Part:          "FIBER",
		}
	}

	var goCloakErr *gocloak.APIError
	if errors.As(err, &goCloakErr) {
		// parse something like '401 Unauthorized: invalid_grant: Invalid user credentials'
		parts := strings.Split(goCloakErr.Message, ":")
		message := strings.TrimSpace(parts[len(parts)-1])
		errorCode := goCloakErr.Code

		if errorCode == 0 {
			errorCode = 500
		}

		return &xgoErrors.XgoError{
			Message:       message,
			IsFatal:       true,
			HttpErrorCode: errorCode,
			Part:          "AUTH_ERROR",
		}
	}

	return &xgoErrors.XgoError{
		Message:       err.Error(),
		IsFatal:       true,
		HttpErrorCode: 500,
		Part:          "UNSPECIFIED",
	}
}
