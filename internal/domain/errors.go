// Package domain 定义领域层的错误类型
package domain

import "errors"

// 领域层预定义错误
var (
	// ErrPathNotFound 路径未找到错误
	ErrPathNotFound = errors.New("路径未找到")

	// ErrInvalidPath 无效路径错误
	ErrInvalidPath = errors.New("无效的路径格式")

	// ErrParameterNotWritable 参数不可写错误
	ErrParameterNotWritable = errors.New("参数不可写")

	// ErrSubscriptionNotFound 订阅未找到错误
	ErrSubscriptionNotFound = errors.New("订阅未找到")

	// ErrInvalidSubscriptionType 无效订阅类型错误
	ErrInvalidSubscriptionType = errors.New("无效的订阅类型")

	// ErrConnectionNotEstablished 连接未建立错误
	ErrConnectionNotEstablished = errors.New("连接未建立")

	// ErrMessageEncodeFailed 消息编码失败错误
	ErrMessageEncodeFailed = errors.New("消息编码失败")

	// ErrMessageDecodeFailed 消息解码失败错误
	ErrMessageDecodeFailed = errors.New("消息解码失败")

	// ErrOperationTimeout 操作超时错误
	ErrOperationTimeout = errors.New("操作超时")

	// ErrInvalidConfiguration 无效配置错误
	ErrInvalidConfiguration = errors.New("无效的配置")
)

// DomainError 领域错误包装器
type DomainError struct {
	Op      string // 操作名称
	Path    string // 相关路径
	Err     error  // 原始错误
	Message string // 错误消息
}

// Error 实现 error 接口
func (e *DomainError) Error() string {
	if e.Path != "" {
		return e.Op + " " + e.Path + ": " + e.Message
	}
	return e.Op + ": " + e.Message
}

// Unwrap 返回原始错误
func (e *DomainError) Unwrap() error {
	return e.Err
}

// NewDomainError 创建新的领域错误
func NewDomainError(op, path, message string, err error) *DomainError {
	return &DomainError{
		Op:      op,
		Path:    path,
		Err:     err,
		Message: message,
	}
}
