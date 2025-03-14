package logger

import (
	"context"
	"os"

	appctx "go-monolith/pkg/context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultLogger *Logger

// Config holds logger configuration
type Config struct {
	Environment    string
	EnableHTTPLogs bool
	LogLevel       string
	Format         string // "json" or "console"
}

// Logger wraps zap logger with additional context methods
type Logger struct {
	*zap.Logger
}

// Initialize creates a new logger instance with the given configuration
func Initialize(config Config) (*Logger, error) {
	// Set log level
	level := zap.InfoLevel
	if config.LogLevel != "" {
		if err := level.UnmarshalText([]byte(config.LogLevel)); err != nil {
			return nil, err
		}
	}

	// Configure encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Choose encoder based on format
	var encoder zapcore.Encoder
	if config.Format == "console" && config.Environment != "production" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// Create core
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		level,
	)

	// Create logger
	zapLogger := zap.New(core)
	if config.Environment != "production" {
		zapLogger = zapLogger.WithOptions(zap.Development())
	}

	logger := &Logger{zapLogger}
	defaultLogger = logger
	return logger, nil
}

// Default returns the global logger instance
func Default() *Logger {
	if defaultLogger == nil {
		// Initialize with default development configuration
		cfg := Config{
			Environment:    "development",
			EnableHTTPLogs: true,
			LogLevel:       "debug",
			Format:         "console",
		}
		logger, err := Initialize(cfg)
		if err != nil {
			// Fallback to basic development logger
			zapLogger, _ := zap.NewDevelopment()
			defaultLogger = &Logger{zapLogger}
		} else {
			defaultLogger = logger
		}
	}
	return defaultLogger
}

// WithContext returns a logger with context values
func (l *Logger) WithContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return l.Logger
	}

	// Get app context
	appCtx := appctx.FromContext(ctx)

	// Add context fields
	fields := []zap.Field{
		zap.String("trace_id", appCtx.TraceID()),
	}

	// Add user ID if available
	if userID := appCtx.UserID(); userID != "" {
		fields = append(fields, zap.String("user_id", userID))
	}

	// Add client info
	if clientIP := appCtx.ClientIP(); clientIP != "" {
		fields = append(fields, zap.String("client_ip", clientIP))
	}

	// Add API version if available
	if version := appCtx.APIVersion(); version != "" {
		fields = append(fields, zap.String("api_version", version))
	}

	return l.Logger.With(fields...)
}

// Debug logs a debug message with context
func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Debug(msg, fields...)
}

// Info logs an info message with context
func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Info(msg, fields...)
}

// Warn logs a warning message with context
func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Warn(msg, fields...)
}

// Error logs an error message with context
func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Error(msg, fields...)
}

// Fatal logs a fatal message with context
func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.WithContext(ctx).Fatal(msg, fields...)
}
