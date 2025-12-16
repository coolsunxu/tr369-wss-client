// Package tr181 定义 TR181 数据模型相关的领域实体
package tr181

// 订阅节点路径常量
const (
	DeviceLocalAgentSubscription = "Device.LocalAgent.Subscription."
)

// 订阅类型常量
const (
	ValueChange       = "ValueChange"       // 值变更订阅
	ObjectCreation    = "ObjectCreation"    // 对象创建订阅
	ObjectDeletion    = "ObjectDeletion"    // 对象删除订阅
	OperationComplete = "OperationComplete" // 操作完成订阅
	Event             = "Event"             // 事件订阅
)

// Handler 定义变更回调函数类型
type Handler func(subscriptionId string, change interface{})

// Listener 表示订阅监听器
type Listener struct {
	SubscriptionId string  // 订阅标识符
	Listener       Handler // 回调处理函数
}

// NewListener 创建新的监听器实例
func NewListener(subscriptionId string, handler Handler) *Listener {
	return &Listener{
		SubscriptionId: subscriptionId,
		Listener:       handler,
	}
}

// Subscription 表示订阅信息
type Subscription struct {
	ID            string // 订阅 ID
	ReferenceList string // 引用列表
	NotifType     string // 通知类型
}

// NewSubscription 创建新的订阅实例
func NewSubscription(id, referenceList, notifType string) *Subscription {
	return &Subscription{
		ID:            id,
		ReferenceList: referenceList,
		NotifType:     notifType,
	}
}
