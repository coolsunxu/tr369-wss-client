// Package protobuf 提供 Protocol Buffers 编解码功能
package protobuf

import (
	"fmt"

	"tr369-wss-client/pkg/api"

	"google.golang.org/protobuf/proto"
)

// Decoder 提供 USP 消息解码功能
type Decoder struct{}

// NewDecoder 创建新的解码器
func NewDecoder() *Decoder {
	return &Decoder{}
}

// DecodeMessage 解码 USP 消息
func (d *Decoder) DecodeMessage(data []byte) (*api.Msg, error) {
	msg := new(api.Msg)

	opts := proto.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	if err := opts.Unmarshal(data, msg); err != nil {
		return nil, fmt.Errorf("解码 USP 消息失败: %w", err)
	}

	return msg, nil
}

// DecodeRecord 解码 USP 记录
func (d *Decoder) DecodeRecord(data []byte) (*api.Record, error) {
	record := new(api.Record)

	if err := proto.Unmarshal(data, record); err != nil {
		return nil, fmt.Errorf("解码 USP 记录失败: %w", err)
	}

	return record, nil
}

// ExtractPayload 从记录中提取消息载荷
func (d *Decoder) ExtractPayload(record *api.Record) ([]byte, error) {
	noSessionContext := record.GetNoSessionContext()
	if noSessionContext == nil {
		return nil, fmt.Errorf("记录不是 NoSessionContextRecord 类型")
	}

	return noSessionContext.GetPayload(), nil
}
