package model

import "fmt"

// 路径校验相关常量
const (
	// MaxPathLength 路径最大长度
	MaxPathLength = 256

	// PathPrefix 路径必须以此开头
	PathPrefix = "Device."

	// WildcardPlaceholder 实例通配符
	WildcardPlaceholder = "{i}"
)

// 预定义错误原因
const (
	ErrReasonEmpty           = "路径不能为空"
	ErrReasonInvalidPrefix   = "路径必须以 Device. 开头"
	ErrReasonTooLong         = "路径长度超过 256 字符限制"
	ErrReasonIllegalChar     = "路径包含非法字符"
	ErrReasonConsecutiveDots = "路径包含连续的点"
	ErrReasonInvalidSegment  = "路径段格式无效"
)

// PathInfo 路径解析结果
type PathInfo struct {
	FullPath     string   // 完整路径，如 "Device.WiFi.Radio.1.Enabled"
	Segments     []string // 路径段，如 ["Device", "WiFi", "Radio", "1", "Enabled"]
	IsObject     bool     // 是否为对象路径（以.结尾）
	IsParameter  bool     // 是否为参数路径（不以.结尾）
	HasWildcard  bool     // 是否包含通配符 {i}
	InstanceNums []int    // 实例编号位置索引（在 Segments 中的位置）
}

// PathValidationError 路径校验错误
type PathValidationError struct {
	Path     string // 原始路径
	Position int    // 错误位置（-1 表示整体错误）
	Reason   string // 错误原因
}

// Error 实现 error 接口
func (e *PathValidationError) Error() string {
	if e.Position >= 0 {
		return fmt.Sprintf("路径校验失败 [%s] 位置 %d: %s", e.Path, e.Position, e.Reason)
	}
	return fmt.Sprintf("路径校验失败 [%s]: %s", e.Path, e.Reason)
}

// NewPathValidationError 创建路径校验错误
func NewPathValidationError(path string, position int, reason string) *PathValidationError {
	return &PathValidationError{
		Path:     path,
		Position: position,
		Reason:   reason,
	}
}
