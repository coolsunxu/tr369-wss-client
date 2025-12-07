package model

// 订阅节点路径
const (
	DeviceLocalAgentSubscription = "Device.LocalAgent.Subscription."
)

const (
	ValueChange       = "ValueChange"
	ObjectCreation    = "ObjectCreation"
	ObjectDeletion    = "ObjectDeletion"
	OperationComplete = "OperationComplete"
	Event             = "Event"
)

// Listener defines a callback for changes
type Handler func(subscriptionId string, change interface{})

type Listener struct {
	SubscriptionId string
	Listener       Handler
}

type TR181DataModel struct {
	Parameters map[string]interface{}
	Listeners  map[string][]Listener
}
