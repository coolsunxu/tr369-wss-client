// Package services 定义领域层的服务接口
package services

import "context"

// WebSocketClient 定义 WebSocket 客户端接口
// 该接口定义了 WebSocket 连接的基本操作
type WebSocketClient interface {
	// Connect 建立 WebSocket 连接
	Connect() error

	// Disconnect 关闭 WebSocket 连接
	Disconnect()

	// IsConnected 返回连接状态
	IsConnected() bool

	// Send 发送消息
	Send(data []byte) error

	// Read 读取消息
	Read() ([]byte, error)

	// Context 返回客户端上下文
	Context() context.Context

	// Cancel 取消客户端上下文
	Cancel()
}

// MessageTransmitter 定义消息传输接口
type MessageTransmitter interface {
	// TransmitMessage 传输消息
	TransmitMessage(payload []byte) error
}

// NotificationService 定义通知服务接口
type NotificationService interface {
	// SendValueChangeNotify 发送值变更通知
	SendValueChangeNotify(subscriptionId, path, value string) error

	// SendObjectCreationNotify 发送对象创建通知
	SendObjectCreationNotify(subscriptionId, path string, uniqueKeys map[string]string) error

	// SendObjectDeletionNotify 发送对象删除通知
	SendObjectDeletionNotify(subscriptionId, path string) error

	// SendOperationCompleteNotify 发送操作完成通知
	SendOperationCompleteNotify(subscriptionId, objPath, commandName, commandKey string, outputArgs map[string]string) error

	// SendEventNotify 发送事件通知
	SendEventNotify(subscriptionId, eventPath string, params map[string]string) error
}
