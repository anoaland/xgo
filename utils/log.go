package utils

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pterm/pterm"
	"github.com/tidwall/pretty"
)

func LogJson(message string, obj any) {
	jsonBytes, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}

	jsonText := string(pretty.Color(jsonBytes, nil))

	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace).WithCallerOffset(1)
	logger.ShowCaller = true
	logger.Debug(pterm.Cyan("JSON ")+message, logger.Args("JSON", jsonText))
}

func Log(message string, args ...any) {
	vargs := make(map[string]any)

	if len(args) == 1 {
		vargs[message] = pterm.Yellow(args[0])
		message = ""
	} else {
		for i, a := range args {
			vargs[strconv.Itoa(i+1)] = pterm.Yellow(a)
		}
	}

	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace).WithCallerOffset(1)
	logger.ShowCaller = true
	logger.Debug(message, logger.ArgsFromMap(vargs))
}

func Info(message string, args ...any) {
	pterm.Info.Printfln(message, args...)
}

func Warn(message string, args ...any) {
	pterm.Warning.Printfln(message, args...)
}

type JsonWriter struct {
	Message string
}

func (w JsonWriter) Write(p []byte) (int, error) {
	var obj any
	if err := json.Unmarshal(p, &obj); err != nil {
		return 0, err
	}

	jsonBytes, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return 0, err
	}
	jsonText := string(pretty.Color(jsonBytes, nil))

	totalLen := len(p) + 1
	if w.Message != "" {
		pterm.Println(w.Message)
		totalLen += len(w.Message) + 1
	}
	pterm.Println(jsonText)
	return totalLen, nil
}
