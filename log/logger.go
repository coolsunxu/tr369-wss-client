package logger

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger // 包级全局

// InitLogger 初始化
func InitLogger() {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,

		// 1. 时间格式化成 “2006-01-02 15:04:05”
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg), // 也可以 NewConsoleEncoder
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	log := zap.New(core, zap.AddCaller())
	logger = log.Sugar()
}

// Sync 刷盘，main() 退出前调用
func Sync() {
	_ = logger.Sync()
}

// Debugf printf 风格
func Debugf(format string, args ...interface{}) { logger.Debugf(format, args...) }
func Infof(format string, args ...interface{})  { logger.Infof(format, args...) }
func Warnf(format string, args ...interface{})  { logger.Warnf(format, args...) }
func Errorf(format string, args ...interface{}) { logger.Errorf(format, args...) }
func Fatalf(format string, args ...interface{}) { logger.Fatalf(format, args...) }

// Debug println 风格
func Debug(args ...interface{}) { logger.Debug(args...) }
func Info(args ...interface{})  { logger.Info(args...) }
func Warn(args ...interface{})  { logger.Warn(args...) }
func Error(args ...interface{}) { logger.Error(args...) }
func Fatal(args ...interface{}) { logger.Fatal(args...) }

// Debugw 结构化风格
func Debugw(msg string, keysAndValues ...interface{}) { logger.Debugw(msg, keysAndValues...) }
func Infow(msg string, keysAndValues ...interface{})  { logger.Infow(msg, keysAndValues...) }
func Warnw(msg string, keysAndValues ...interface{})  { logger.Warnw(msg, keysAndValues...) }
func Errorw(msg string, keysAndValues ...interface{}) { logger.Errorw(msg, keysAndValues...) }
func Fatalw(msg string, keysAndValues ...interface{}) { logger.Fatalw(msg, keysAndValues...) }

func GetTraceId(ctx context.Context) string {
	traceId, ok := ctx.Value("trace_id").(string)
	if ok {
		return traceId
	} else {
		return ""
	}
}
