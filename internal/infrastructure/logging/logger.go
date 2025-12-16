// Package logging 提供日志记录功能
package logging

import (
	"fmt"
	"log"
	"os"

	"tr369-wss-client/internal/domain/services"
)

// Logger 是 services.Logger 的别名，保持向后兼容
type Logger = services.Logger

// DefaultLogger 默认日志实现
type DefaultLogger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
}

// NewLogger 创建新的日志实例
func NewLogger() Logger {
	return &DefaultLogger{
		debugLogger: log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile),
		infoLogger:  log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger:  log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
		fatalLogger: log.New(os.Stderr, "[FATAL] ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Debug 输出调试日志
func (l *DefaultLogger) Debug(msg string, args ...interface{}) {
	l.debugLogger.Output(2, fmt.Sprintf(msg, args...))
}

// Info 输出信息日志
func (l *DefaultLogger) Info(msg string, args ...interface{}) {
	l.infoLogger.Output(2, fmt.Sprintf(msg, args...))
}

// Warn 输出警告日志
func (l *DefaultLogger) Warn(msg string, args ...interface{}) {
	l.warnLogger.Output(2, fmt.Sprintf(msg, args...))
}

// Error 输出错误日志
func (l *DefaultLogger) Error(msg string, args ...interface{}) {
	l.errorLogger.Output(2, fmt.Sprintf(msg, args...))
}

// Fatal 输出致命错误日志并退出
func (l *DefaultLogger) Fatal(msg string, args ...interface{}) {
	l.fatalLogger.Output(2, fmt.Sprintf(msg, args...))
	os.Exit(1)
}
