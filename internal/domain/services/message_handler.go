// Package services 定义领域层的服务接口
package services

import (
	"tr369-wss-client/pkg/api"
)

// MessageHandler 定义消息处理服务接口
// 该接口定义了处理各种 USP 消息的方法
type MessageHandler interface {
	// HandleMessage 处理传入的 USP 消息
	HandleMessage(msg *api.Msg)

	// HandleGetRequest 处理 GET 请求
	HandleGetRequest(msg *api.Msg)

	// HandleSetRequest 处理 SET 请求
	HandleSetRequest(msg *api.Msg)

	// HandleAddRequest 处理 ADD 请求
	HandleAddRequest(msg *api.Msg)

	// HandleDeleteRequest 处理 DELETE 请求
	HandleDeleteRequest(msg *api.Msg)

	// HandleOperateRequest 处理 OPERATE 请求
	HandleOperateRequest(msg *api.Msg)

	// HandleNotifyResp 处理 NOTIFY 响应
	HandleNotifyResp(msg *api.Msg)
}

// SubscriptionManager 定义订阅管理服务接口
type SubscriptionManager interface {
	// HandleSubscription 处理订阅
	HandleSubscription(path, subscriptionId, subscriptionType string) error

	// HandleAddSubscription 处理添加订阅
	HandleAddSubscription(requestPath string, paramSettings map[string]string) error

	// HandleDeleteSubscription 处理删除订阅
	HandleDeleteSubscription(requestPath string) error
}
