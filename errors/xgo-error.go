package errors

import (
	"fmt"
	"strings"
)

type XgoError struct {
	part    string
	err     error
	message string
	// stack   []Trace
}

func NewXgoError(part string, err error) *XgoError {
	msg := err.Error()
	parts := []string{}

	if me, ok := err.(*XgoError); ok {
		parts = append(parts, me.part)
		msg = me.message
	}

	parts = append(parts, part)

	return &XgoError{
		part:    strings.Join(parts, "\r\n"),
		err:     err,
		message: msg,
		// stack:   getStack(),
	}
}

func (e *XgoError) Error() string {
	return fmt.Sprintf("%s\r\n%s", e.part, e.message)
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
