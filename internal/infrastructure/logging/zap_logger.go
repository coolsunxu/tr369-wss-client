// Package logging 提供基于 zap 的日志实现
package logging

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger 基于 zap 的日志实现
type ZapLogger struct {
	logger *zap.SugaredLogger
}

// NewZapLogger 创建新的 zap 日志实例
func NewZapLogger() *ZapLogger {
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
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	log := zap.New(core, zap.AddCaller())
	return &ZapLogger{logger: log.Sugar()}
}

// Sync 刷盘
func (l *ZapLogger) Sync() {
	_ = l.logger.Sync()
}

// Debug 输出调试日志
func (l *ZapLogger) Debug(msg string, args ...interface{}) {
	l.logger.Debugf(msg, args...)
}

// Info 输出信息日志
func (l *ZapLogger) Info(msg string, args ...interface{}) {
	l.logger.Infof(msg, args...)
}

// Warn 输出警告日志
func (l *ZapLogger) Warn(msg string, args ...interface{}) {
	l.logger.Warnf(msg, args...)
}

// Error 输出错误日志
func (l *ZapLogger) Error(msg string, args ...interface{}) {
	l.logger.Errorf(msg, args...)
}

// Fatal 输出致命错误日志
func (l *ZapLogger) Fatal(msg string, args ...interface{}) {
	l.logger.Fatalf(msg, args...)
}

// GetTraceId 从上下文获取追踪 ID
func GetTraceId(ctx context.Context) string {
	traceId, ok := ctx.Value("trace_id").(string)
	if ok {
		return traceId
	}
	return ""
}
