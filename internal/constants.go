package internal

// Context key constants to avoid typos and ensure consistency across packages
const (
	RequestLoggerKey = "xgo_request_logger"
	RequestIDKey     = "xgo_use_logger_requestID"
	StartTimeKey     = "xgo_use_logger_startTime"
	StackErrorKey    = "xgo_use_logger_stackError"
)

// Define context key type to avoid collisions
type ContextKey string

const FiberContextKey ContextKey = "fiber"
