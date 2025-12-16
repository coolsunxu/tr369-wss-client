// Package services 定义领域层的服务接口
package services

// Logger 定义日志服务接口
// 该接口定义了日志记录的基本操作，具体实现由基础设施层提供
type Logger interface {
	// Debug 输出调试级别日志
	Debug(msg string, args ...interface{})

	// Info 输出信息级别日志
	Info(msg string, args ...interface{})

	// Warn 输出警告级别日志
	Warn(msg string, args ...interface{})

	// Error 输出错误级别日志
	Error(msg string, args ...interface{})

	// Fatal 输出致命错误日志并退出程序
	Fatal(msg string, args ...interface{})
}
