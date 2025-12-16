// Package usp 测试 USP 消息实体
package usp

import (
	"testing"
)

// **Feature: code-structure-optimization, Property 12: 接口模拟测试支持**
// **验证需求: 需求 3.4**
// 对于任何定义的接口，都应该能够创建模拟对象进行单元测试

func TestMessageCreation(t *testing.T) {
	msg := NewMessage("test-id", MessageTypeGet, nil)

	if msg.ID != "test-id" {
		t.Errorf("期望消息 ID 为 'test-id', 实际为 '%s'", msg.ID)
	}

	if msg.Type != MessageTypeGet {
		t.Errorf("期望消息类型为 GET, 实际为 '%s'", msg.Type.String())
	}
}

func TestMessageTypeString(t *testing.T) {
	testCases := []struct {
		msgType  MessageType
		expected string
	}{
		{MessageTypeGet, "GET"},
		{MessageTypeGetResp, "GET_RESP"},
		{MessageTypeSet, "SET"},
		{MessageTypeSetResp, "SET_RESP"},
		{MessageTypeAdd, "ADD"},
		{MessageTypeAddResp, "ADD_RESP"},
		{MessageTypeDelete, "DELETE"},
		{MessageTypeDeleteResp, "DELETE_RESP"},
		{MessageTypeOperate, "OPERATE"},
		{MessageTypeOperateResp, "OPERATE_RESP"},
		{MessageTypeNotify, "NOTIFY"},
		{MessageTypeNotifyResp, "NOTIFY_RESP"},
		{MessageTypeUnknown, "UNKNOWN"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			if tc.msgType.String() != tc.expected {
				t.Errorf("期望 '%s', 实际为 '%s'", tc.expected, tc.msgType.String())
			}
		})
	}
}

func TestHeaderCreation(t *testing.T) {
	header := NewHeader("header-id", MessageTypeSet)

	if header.MessageID != "header-id" {
		t.Errorf("期望消息 ID 为 'header-id', 实际为 '%s'", header.MessageID)
	}

	if header.MessageType != MessageTypeSet {
		t.Errorf("期望消息类型为 SET, 实际为 '%s'", header.MessageType.String())
	}
}

func TestRecordCreation(t *testing.T) {
	payload := []byte("test payload")
	record := NewRecord("1.0", "to-id", "from-id", payload)

	if record.Version != "1.0" {
		t.Errorf("期望版本为 '1.0', 实际为 '%s'", record.Version)
	}

	if record.ToID != "to-id" {
		t.Errorf("期望目标 ID 为 'to-id', 实际为 '%s'", record.ToID)
	}

	if record.FromID != "from-id" {
		t.Errorf("期望来源 ID 为 'from-id', 实际为 '%s'", record.FromID)
	}

	if string(record.Payload) != "test payload" {
		t.Errorf("期望载荷为 'test payload', 实际为 '%s'", string(record.Payload))
	}
}

func TestRecordIsValid(t *testing.T) {
	testCases := []struct {
		name     string
		record   *Record
		expected bool
	}{
		{
			name:     "有效记录",
			record:   NewRecord("1.0", "to", "from", []byte("data")),
			expected: true,
		},
		{
			name:     "空版本",
			record:   NewRecord("", "to", "from", []byte("data")),
			expected: false,
		},
		{
			name:     "空目标 ID",
			record:   NewRecord("1.0", "", "from", []byte("data")),
			expected: false,
		},
		{
			name:     "空来源 ID",
			record:   NewRecord("1.0", "to", "", []byte("data")),
			expected: false,
		},
		{
			name:     "空载荷",
			record:   NewRecord("1.0", "to", "from", []byte{}),
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.record.IsValid() != tc.expected {
				t.Errorf("期望 IsValid() 返回 %v, 实际返回 %v", tc.expected, tc.record.IsValid())
			}
		})
	}
}
