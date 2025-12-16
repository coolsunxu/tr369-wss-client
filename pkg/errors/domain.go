// Package errors 定义项目中使用的错误类型
package errors

import "fmt"

// DomainError 表示领域层错误
type DomainError struct {
	Code    string // 错误代码
	Message string // 错误消息
	Cause   error  // 原因错误
}

// Error 实现 error 接口
func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 返回原因错误
func (e *DomainError) Unwrap() error {
	return e.Cause
}

// NewDomainError 创建新的领域错误
func NewDomainError(code, message string, cause error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// 预定义的领域错误代码
const (
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeInvalidInput   = "INVALID_INPUT"
	ErrCodeInternalError  = "INTERNAL_ERROR"
	ErrCodeConnectionFail = "CONNECTION_FAIL"
	ErrCodeTimeout        = "TIMEOUT"
)

// ErrNotFound 创建未找到错误
func ErrNotFound(message string) *DomainError {
	return NewDomainError(ErrCodeNotFound, message, nil)
}

// ErrInvalidInput 创建无效输入错误
func ErrInvalidInput(message string) *DomainError {
	return NewDomainError(ErrCodeInvalidInput, message, nil)
}

// ErrInternal 创建内部错误
func ErrInternal(message string, cause error) *DomainError {
	return NewDomainError(ErrCodeInternalError, message, cause)
}

// ErrConnectionFail 创建连接失败错误
func ErrConnectionFail(message string, cause error) *DomainError {
	return NewDomainError(ErrCodeConnectionFail, message, cause)
}

// ErrTimeout 创建超时错误
func ErrTimeout(message string) *DomainError {
	return NewDomainError(ErrCodeTimeout, message, nil)
}
