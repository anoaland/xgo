package errors

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type XgoError struct {
	Part          string
	IsFatal       bool
	Err           error
	Message       string
	File          string
	Line          int
	HttpErrorCode int
	Stack         string
}

func (e *XgoError) AsFiberError(ctx *fiber.Ctx) error {
	if e.IsFatal {
		return ctx.Status(500).JSON(fiber.Map{
			"message":    "Terjadi kesalahan",
			"code":       500,
			"statusCode": "INTERNAL_SERVER_ERROR",
		})
	}

	return ctx.Status(e.HttpErrorCode).JSON(fiber.Map{
		"message":    e.Message,
		"code":       e.HttpErrorCode,
		"statusCode": e.Part,
	})
}

// Deprecated: Use NewError instead
func NewXgoError(part string, err error) *XgoError {
	return NewHttpError(part, err, 500, 0)
}

func NewError(part string, err error) *XgoError {
	return NewHttpError(part, err, 500, 0)
}

func NewHttpError(part string, err error, httpErrorCode int, callerSkip int) *XgoError {
	_, file, line, _ := runtime.Caller(callerSkip + 1)
	msg := err.Error()
	parts := []string{}

	if me, ok := err.(*XgoError); ok {
		parts = append([]string{me.Part}, parts...)
		msg = me.Message
	}

	parts = append(parts, part)

	return &XgoError{
		Part:          strings.Join(parts, " -> "),
		Err:           err,
		Message:       msg,
		File:          file,
		Line:          line,
		HttpErrorCode: httpErrorCode,
		IsFatal:       httpErrorCode >= 500,
		Stack:         string(debug.Stack()),
	}
}

func (e *XgoError) Error() string {
	identity := fmt.Sprintf("[%d]", e.HttpErrorCode)
	if e.Part != "" {
		identity = fmt.Sprintf("[%d | %s]", e.HttpErrorCode, e.Part)
	}

	return fmt.Sprintf("%s %s\r\n\t%s:%d", identity, e.Message, e.File, e.Line)
}

// see: https://mdcfrancis.medium.com/tracing-errors-in-go-using-custom-error-types-9aaf3bba1a64
// func (e *XgoError) Trace() string {
// 	strs := []string{}

// 	for _, st := range e.stack {
// 		strs = append(strs, fmt.Sprintf("%s -- %s:%v", st.File, st.Function, st.Line))
// 	}

// 	return strings.Join(strs, "\r\n")
// }

// type Trace struct {
// 	Function string `json:"function"`
// 	File     string `json:"file"`
// 	Line     int    `json:"line"`
// }

// func getStack() []Trace {
// 	pcs := make([]uintptr, 32)
// 	// Skip 4 stack frames
// 	npcs := runtime.Callers(4, pcs)
// 	traces := make([]Trace, 0, npcs)
// 	callers := pcs[:npcs]
// 	ci := runtime.CallersFrames(callers)
// 	for {
// 		frame, more := ci.Next()
// 		traces = append(traces, Trace{
// 			File:     frame.File,
// 			Line:     frame.Line,
// 			Function: frame.Function,
// 		})
// 		if !more || frame.Function == "main.main" {
// 			break
// 		}
// 	}
// 	return traces
// }
