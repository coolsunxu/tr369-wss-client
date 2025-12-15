package model

import (
	"tr369-wss-client/pkg/api"
	tr181Model "tr369-wss-client/tr181/model"
)

// 路径常量定义
const (
	PathDevice       = "Device."
	PathLocalAgent   = "Device.LocalAgent."
	PathSubscription = "Device.LocalAgent.Subscription."
)

// ParamSetting 定义参数设置的通用接口
// 用于统一处理 Set_UpdateParamSetting 和 Add_CreateParamSetting
type ParamSetting interface {
	GetParam() string
	GetValue() string
}

// SubscriptionParams 订阅参数结构体
type SubscriptionParams struct {
	Id            string
	ReferenceList string
	NotifType     string
}

// WSClient defines the interface for TR369 WebSocket client
type WSClient interface {
	// Connect establishes a WebSocket connection to the server
	Connect() error

	// Disconnect closes the WebSocket connection
	Disconnect()

	// StartMessageHandler starts the message handling goroutines
	StartMessageHandler()
}

// DataRepository 定义数据访问接口
// 负责 TR181 数据模型的纯数据 CRUD 操作
type DataRepository interface {
	// GetValue 获取指定路径的值
	GetValue(path string) (interface{}, error)

	// GetParameters 获取底层参数数据（供 UseCase 构建响应使用）
	GetParameters() map[string]interface{}

	// SetValue 设置指定路径的值
	// 返回: changed (是否发生变化), oldValue (旧值)
	SetValue(path string, key string, value string) (changed bool, oldValue string)

	// DeleteNode 删除指定路径的节点
	DeleteNode(path string) (nodePath string, isFound bool)

	// Start 启动数据仓库（初始化和数据同步）
	Start()
}

// ListenerManager 定义监听器管理接口
// 负责事件监听器的管理
type ListenerManager interface {
	// AddListener 添加参数变化监听器
	AddListener(paramName string, listener tr181Model.Listener) error

	// RemoveListener 移除指定参数的监听器
	RemoveListener(paramName string) error

	// ResetListener 重置所有监听器
	ResetListener() error

	// NotifyListeners 通知指定参数的所有监听器
	NotifyListeners(paramName string, value interface{})
}

// ClientUseCase defines the interface for client use case
type ClientUseCase interface {
	// HandleMessage processes incoming USP messages
	HandleMessage(msg *api.Msg)
}
