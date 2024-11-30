package errors

import (
	"fmt"
	"runtime"
	"strings"
)

type XgoError struct {
	Part          string
	Err           error
	Message       string
	File          string
	Line          int
	HttpErrorCode int
	// stack   []Trace
}

// Deprecated: Use NewError instead
func NewXgoError(part string, err error) *XgoError {
	return NewHttpError(part, err, 500)
}

func NewError(part string, err error) *XgoError {
	return NewHttpError(part, err, 500)
}

func NewHttpError(part string, err error, httpErrorCode int) *XgoError {
	_, file, line, _ := runtime.Caller(1)
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
		// stack:   getStack(),
	}
}

func (e *XgoError) Error() string {
	return fmt.Sprintf("[%d | %s] %s | %s:%d", e.HttpErrorCode, e.Part, e.Message, e.File, e.Line)
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
