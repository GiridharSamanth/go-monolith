package logger

import (
	"context"
	"os"

	appctx "go-monolith/pkg/context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultLogger Logger

// Config holds logger configuration
type Config struct {
	Environment    string
	EnableHTTPLogs bool
	LogLevel       string
	Format         string // "json" or "console"
}

// Field represents a key-value pair for structured logging
type Field struct {
	Key   string
	Value interface{}
}

// Logger interface defines the logging methods
type Logger interface {
	Debug(ctx context.Context, msg string, fields ...Field)
	Info(ctx context.Context, msg string, fields ...Field)
	Warn(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, fields ...Field)
	Fatal(ctx context.Context, msg string, fields ...Field)
}

// zapLogger implements Logger interface using zap
type zapLogger struct {
	*zap.Logger
}

// Helper functions to create fields
func String(key string, value string) Field {
	return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

func Float64(key string, value float64) Field {
	return Field{Key: key, Value: value}
}

func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// Initialize creates a new logger instance with the given configuration
func Initialize(config Config) (Logger, error) {
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
	zl := zap.New(core)
	if config.Environment != "production" {
		zl = zl.WithOptions(zap.Development())
	}

	logger := &zapLogger{zl}
	defaultLogger = logger
	return logger, nil
}

// Default returns the global logger instance
func Default() Logger {
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
			zl, _ := zap.NewDevelopment()
			defaultLogger = &zapLogger{zl}
		} else {
			defaultLogger = logger
		}
	}
	return defaultLogger
}

// convertFields converts our Field type to zap.Field
func convertFields(fields ...Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		switch v := f.Value.(type) {
		case string:
			zapFields[i] = zap.String(f.Key, v)
		case int:
			zapFields[i] = zap.Int(f.Key, v)
		case int64:
			zapFields[i] = zap.Int64(f.Key, v)
		case float64:
			zapFields[i] = zap.Float64(f.Key, v)
		case bool:
			zapFields[i] = zap.Bool(f.Key, v)
		default:
			zapFields[i] = zap.Any(f.Key, v)
		}
	}
	return zapFields
}

// withContext returns a zap logger with context values
func (l *zapLogger) withContext(ctx context.Context) *zap.Logger {
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
func (l *zapLogger) Debug(ctx context.Context, msg string, fields ...Field) {
	l.withContext(ctx).Debug(msg, convertFields(fields...)...)
}

// Info logs an info message with context
func (l *zapLogger) Info(ctx context.Context, msg string, fields ...Field) {
	l.withContext(ctx).Info(msg, convertFields(fields...)...)
}

// Warn logs a warning message with context
func (l *zapLogger) Warn(ctx context.Context, msg string, fields ...Field) {
	l.withContext(ctx).Warn(msg, convertFields(fields...)...)
}

// Error logs an error message with context
func (l *zapLogger) Error(ctx context.Context, msg string, fields ...Field) {
	l.withContext(ctx).Error(msg, convertFields(fields...)...)
}

// Fatal logs a fatal message with context
func (l *zapLogger) Fatal(ctx context.Context, msg string, fields ...Field) {
	l.withContext(ctx).Fatal(msg, convertFields(fields...)...)
}
