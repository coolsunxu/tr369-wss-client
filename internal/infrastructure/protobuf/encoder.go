// Package protobuf 提供 Protocol Buffers 编解码功能
package protobuf

import (
	"tr369-wss-client/pkg/api"

	"google.golang.org/protobuf/proto"
)

// Encoder 提供 USP 消息编码功能
type Encoder struct{}

// NewEncoder 创建新的编码器
func NewEncoder() *Encoder {
	return &Encoder{}
}

// EncodeMessage 编码 USP 消息
func (e *Encoder) EncodeMessage(msg *api.Msg) ([]byte, error) {
	return proto.Marshal(msg)
}

// EncodeRecord 编码 USP 记录
func (e *Encoder) EncodeRecord(rec *api.Record) ([]byte, error) {
	return proto.Marshal(rec)
}

// CreateRecord 创建 USP 记录（无会话上下文）
func (e *Encoder) CreateRecord(version, toId, fromId string, msg *api.Msg) (*api.Record, error) {
	msgBytes, err := e.EncodeMessage(msg)
	if err != nil {
		return nil, err
	}

	return &api.Record{
		Version: version,
		ToId:    toId,
		FromId:  fromId,
		RecordType: &api.Record_NoSessionContext{
			NoSessionContext: &api.NoSessionContextRecord{
				Payload: msgBytes,
			},
		},
	}, nil
}
