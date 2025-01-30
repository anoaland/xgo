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
