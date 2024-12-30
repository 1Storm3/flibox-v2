package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

func Init(env string) {
	atomic := zap.NewAtomicLevel()

	var enconfig zapcore.EncoderConfig
	var enc zapcore.Encoder

	switch env {
	case "production":
		atomic.SetLevel(zapcore.InfoLevel)
		enconfig = zap.NewProductionEncoderConfig()
		enc = zapcore.NewJSONEncoder(enconfig)

	default:
		atomic.SetLevel(zapcore.DebugLevel)
		enconfig = zap.NewDevelopmentEncoderConfig()
		enc = zapcore.NewConsoleEncoder(enconfig)
	}

	globalLogger = zap.New(zapcore.NewCore(
		enc,
		zapcore.Lock(os.Stdout),
		atomic,
	))
}

func ensureInitialized() {
	if globalLogger == nil {
		Init("development")
	}
}

func Debug(msg string, fields ...zap.Field) {
	ensureInitialized()
	globalLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	ensureInitialized()
	globalLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	ensureInitialized()
	globalLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	ensureInitialized()
	globalLogger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	ensureInitialized()
	globalLogger.Fatal(msg, fields...)
}

func WithOptions(opts ...zap.Option) *zap.Logger {
	ensureInitialized()
	return globalLogger.WithOptions(opts...)
}
