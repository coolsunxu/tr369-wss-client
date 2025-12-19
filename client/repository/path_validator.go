package repository

import (
	"strconv"
	"strings"
	"tr369-wss-client/client/model"
	"unicode"
)

// PathValidator 路径校验器
type PathValidator struct{}

// NewPathValidator 创建路径校验器实例
func NewPathValidator() *PathValidator {
	return &PathValidator{}
}

// ValidatePath 校验路径是否符合 TR181 Path Name 规范
func (v *PathValidator) ValidatePath(path string) error {
	// 1. 检查空路径
	if path == "" {
		return model.NewPathValidationError(path, -1, model.ErrReasonEmpty)
	}

	// 2. 检查长度
	if len(path) > model.MaxPathLength {
		return model.NewPathValidationError(path, model.MaxPathLength, model.ErrReasonTooLong)
	}

	// 3. 检查前缀
	if !strings.HasPrefix(path, model.PathPrefix) {
		return model.NewPathValidationError(path, 0, model.ErrReasonInvalidPrefix)
	}

	// 4. 检查非法字符（控制字符）
	for i, r := range path {
		if v.isIllegalChar(r) {
			return model.NewPathValidationError(path, i, model.ErrReasonIllegalChar)
		}
	}

	// 5. 检查连续点
	if pos := strings.Index(path, ".."); pos >= 0 {
		return model.NewPathValidationError(path, pos, model.ErrReasonConsecutiveDots)
	}

	return nil
}

// isIllegalChar 检查是否为非法字符
// 非法字符包括：控制字符（0x00-0x1F，除了常见的空白字符如 \t, \n, \r）
func (v *PathValidator) isIllegalChar(r rune) bool {
	// 控制字符范围 0x00-0x1F，但排除 \t(0x09), \n(0x0A), \r(0x0D)
	if r < 0x20 && r != '\t' && r != '\n' && r != '\r' {
		return true
	}
	// 检查是否为其他控制字符
	if unicode.IsControl(r) && r != '\t' && r != '\n' && r != '\r' {
		return true
	}
	return false
}

// IsObjectPath 判断是否为对象路径（以 . 结尾）
func (v *PathValidator) IsObjectPath(path string) bool {
	if path == "" {
		return false
	}
	return strings.HasSuffix(path, ".")
}

// IsParameterPath 判断是否为参数路径（不以 . 结尾）
func (v *PathValidator) IsParameterPath(path string) bool {
	if path == "" {
		return false
	}
	return !strings.HasSuffix(path, ".")
}

// ParsePath 解析路径信息
func (v *PathValidator) ParsePath(path string) (*model.PathInfo, error) {
	// 先校验路径
	if err := v.ValidatePath(path); err != nil {
		return nil, err
	}

	info := &model.PathInfo{
		FullPath:     path,
		IsObject:     v.IsObjectPath(path),
		IsParameter:  v.IsParameterPath(path),
		HasWildcard:  strings.Contains(path, model.WildcardPlaceholder),
		InstanceNums: []int{},
	}

	// 分割路径段
	// 如果以 . 结尾，去掉最后的空字符串
	segments := strings.Split(path, ".")
	if len(segments) > 0 && segments[len(segments)-1] == "" {
		segments = segments[:len(segments)-1]
	}
	info.Segments = segments

	// 识别实例编号位置
	for i, seg := range segments {
		// 检查是否为数字（实例编号）
		if _, err := strconv.Atoi(seg); err == nil {
			info.InstanceNums = append(info.InstanceNums, i)
		}
	}

	return info, nil
}
