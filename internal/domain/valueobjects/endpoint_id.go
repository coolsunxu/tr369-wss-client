// Package valueobjects 定义领域层的值对象
package valueobjects

import (
	"errors"
	"strings"
)

// ErrEmptyEndpointID 表示端点 ID 为空的错误
var ErrEmptyEndpointID = errors.New("端点 ID 不能为空")

// ErrInvalidEndpointID 表示端点 ID 格式无效的错误
var ErrInvalidEndpointID = errors.New("端点 ID 格式无效")

// EndpointID 表示端点标识符值对象
// 端点 ID 是不可变的，一旦创建就不能修改
type EndpointID struct {
	value string
}

// NewEndpointID 创建新的端点 ID 值对象
// 如果 ID 为空或格式无效，返回错误
func NewEndpointID(id string) (*EndpointID, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrEmptyEndpointID
	}
	return &EndpointID{value: id}, nil
}

// Value 返回端点 ID 的字符串值
func (e *EndpointID) Value() string {
	return e.value
}

// String 实现 Stringer 接口
func (e *EndpointID) String() string {
	return e.value
}

// Equals 比较两个端点 ID 是否相等
func (e *EndpointID) Equals(other *EndpointID) bool {
	if other == nil {
		return false
	}
	return e.value == other.value
}

// MessageType 表示消息类型值对象
type MessageType struct {
	value string
}

// NewMessageType 创建新的消息类型值对象
func NewMessageType(msgType string) *MessageType {
	return &MessageType{value: msgType}
}

// Value 返回消息类型的字符串值
func (m *MessageType) Value() string {
	return m.value
}

// String 实现 Stringer 接口
func (m *MessageType) String() string {
	return m.value
}
