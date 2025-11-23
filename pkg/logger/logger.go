package logger

import (
	"context"
	"os"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger with additional functionality
type Logger struct {
	*zap.Logger
	sugar *zap.SugaredLogger
}

// Config holds logger configuration
type Config struct {
	Level       string // debug, info, warn, error
	Environment string // development, production
	ServiceName string
}

// New creates a new logger instance
func New(cfg Config) (*Logger, error) {
	var zapConfig zap.Config

	if cfg.Environment == "production" {
		zapConfig = zap.NewProductionConfig()
	} else {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Set log level
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)

	// Add service name to all logs
	zapConfig.InitialFields = map[string]interface{}{
		"service": cfg.ServiceName,
	}

	zapLogger, err := zapConfig.Build(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, err
	}

	return &Logger{
		Logger: zapLogger,
		sugar:  zapLogger.Sugar(),
	}, nil
}

// NewDefault creates a logger with default configuration
func NewDefault(serviceName string) (*Logger, error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}

	return New(Config{
		Level:       level,
		Environment: env,
		ServiceName: serviceName,
	})
}

// WithContext adds trace context to logger
func (l *Logger) WithContext(ctx context.Context) *Logger {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return l
	}

	spanContext := span.SpanContext()
	return &Logger{
		Logger: l.With(
			zap.String("trace_id", spanContext.TraceID().String()),
			zap.String("span_id", spanContext.SpanID().String()),
		),
		sugar: l.sugar.With(
			"trace_id", spanContext.TraceID().String(),
			"span_id", spanContext.SpanID().String(),
		),
	}
}

// WithFields adds structured fields to logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}

	return &Logger{
		Logger: l.With(zapFields...),
		sugar:  l.sugar.With(fields),
	}
}

// Sugar returns sugared logger for printf-style logging
func (l *Logger) Sugar() *zap.SugaredLogger {
	return l.sugar
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

// Global logger instance (optional, for convenience)
var globalLogger *Logger

// InitGlobal initializes the global logger
func InitGlobal(cfg Config) error {
	logger, err := New(cfg)
	if err != nil {
		return err
	}
	globalLogger = logger
	return nil
}

// Global returns the global logger instance
func Global() *Logger {
	if globalLogger == nil {
		// Fallback to default logger
		logger, _ := NewDefault("unknown")
		return logger
	}
	return globalLogger
}
