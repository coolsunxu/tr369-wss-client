// Package usp 定义 USP 消息相关的领域实体
package usp

// Record 表示 USP 记录
type Record struct {
	Version string // 版本
	ToID    string // 目标 ID
	FromID  string // 来源 ID
	Payload []byte // 载荷
}

// NewRecord 创建新的记录实例
func NewRecord(version, toID, fromID string, payload []byte) *Record {
	return &Record{
		Version: version,
		ToID:    toID,
		FromID:  fromID,
		Payload: payload,
	}
}

// IsValid 验证记录是否有效
func (r *Record) IsValid() bool {
	return r.Version != "" && r.ToID != "" && r.FromID != "" && len(r.Payload) > 0
}
