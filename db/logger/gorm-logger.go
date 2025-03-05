package logger

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm/logger"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// ZerologGormLogger implements gorm/logger.Interface
type ZerologGormLogger struct {
	logger   zerolog.Logger
	config   gormlogger.Config
	LogLevel logger.LogLevel
}

func NewZerologGormLogger(logger *zerolog.Logger, configs ...gormlogger.Config) *ZerologGormLogger {
	var config gormlogger.Config
	if len(configs) > 0 {
		config = configs[0]
	} else {
		config = gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  gormlogger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		}
	}

	return &ZerologGormLogger{
		logger:   *logger,
		config:   config,
		LogLevel: config.LogLevel,
	}
}

// LogMode sets the log level
func (l *ZerologGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info logs general info messages
func (l *ZerologGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.logger.Info().Msgf(msg, data...)
	}
}

// Warn logs warning messages
func (l *ZerologGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.logger.Warn().Msgf(msg, data...)
	}
}

// Error logs error messages
func (l *ZerologGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.logger.Error().Msgf(msg, data...)
	}
}

// Trace logs SQL statements with duration
func (l *ZerologGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	latency := time.Since(begin)
	sql, rows := fc()

	msg := "SQL query"
	event := l.logger.Info()
	if err != nil {
		event = l.logger.Error().Err(err)
	} else if latency > l.config.SlowThreshold {
		msg = "Slow query"
		event = l.logger.Warn()
	}

	// Get Fiber context (see: server.LoggerContext)
	if fiberCtx, ok := ctx.Value("fiber").(*fiber.Ctx); ok {
		requestID := fiberCtx.Locals("xgo_use_logger_requestID")
		event.Any("request_id", requestID)
	}

	arr := zerolog.Arr()
	arr.Str(utils.FileWithLineNum())

	event.Str("sql", sql).
		Int64("rows", rows).
		Str("latency", latency.String()).
		Array("stack", arr).
		Msg(msg)
}
