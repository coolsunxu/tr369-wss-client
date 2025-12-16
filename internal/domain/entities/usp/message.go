// Package usp 定义 USP 消息相关的领域实体
package usp

// MessageType 消息类型
type MessageType int

const (
	MessageTypeUnknown MessageType = iota
	MessageTypeGet
	MessageTypeGetResp
	MessageTypeSet
	MessageTypeSetResp
	MessageTypeAdd
	MessageTypeAddResp
	MessageTypeDelete
	MessageTypeDeleteResp
	MessageTypeOperate
	MessageTypeOperateResp
	MessageTypeNotify
	MessageTypeNotifyResp
)

// String 返回消息类型的字符串表示
func (mt MessageType) String() string {
	switch mt {
	case MessageTypeGet:
		return "GET"
	case MessageTypeGetResp:
		return "GET_RESP"
	case MessageTypeSet:
		return "SET"
	case MessageTypeSetResp:
		return "SET_RESP"
	case MessageTypeAdd:
		return "ADD"
	case MessageTypeAddResp:
		return "ADD_RESP"
	case MessageTypeDelete:
		return "DELETE"
	case MessageTypeDeleteResp:
		return "DELETE_RESP"
	case MessageTypeOperate:
		return "OPERATE"
	case MessageTypeOperateResp:
		return "OPERATE_RESP"
	case MessageTypeNotify:
		return "NOTIFY"
	case MessageTypeNotifyResp:
		return "NOTIFY_RESP"
	default:
		return "UNKNOWN"
	}
}

// Message 表示 USP 消息
type Message struct {
	ID      string      // 消息 ID
	Type    MessageType // 消息类型
	Payload interface{} // 消息载荷
}

// NewMessage 创建新的消息实例
func NewMessage(id string, msgType MessageType, payload interface{}) *Message {
	return &Message{
		ID:      id,
		Type:    msgType,
		Payload: payload,
	}
}

// Header 表示消息头
type Header struct {
	MessageID   string      // 消息 ID
	MessageType MessageType // 消息类型
}

// NewHeader 创建新的消息头
func NewHeader(id string, msgType MessageType) *Header {
	return &Header{
		MessageID:   id,
		MessageType: msgType,
	}
}
