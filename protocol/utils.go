package protocol

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// ValidateMessage 验证消息的基本有效性
func ValidateMessage(msg *Message) error {
	if msg == nil {
		return fmt.Errorf("message cannot be nil")
	}

	if msg.MessageType == "" {
		return fmt.Errorf("message type cannot be empty")
	}

	if msg.MsgID == "" {
		return fmt.Errorf("message ID cannot be empty")
	}

	if msg.Body == nil {
		return fmt.Errorf("message body cannot be nil")
	}

	// 验证消息类型
	validTypes := map[string]bool{
		MessageTypeGet:    true,
		MessageTypeSet:    true,
		MessageTypeInform: true,
		MessageTypeAdd:    true,
		MessageTypeDelete: true,
		MessageTypeReplace: true,
		MessageTypeOperation: true,
	}

	if !validTypes[msg.MessageType] {
		return fmt.Errorf("invalid message type: %s", msg.MessageType)
	}

	return nil
}

// ValidateParameterName 验证参数名称是否符合TR181格式
func ValidateParameterName(name string) error {
	if name == "" {
		return fmt.Errorf("parameter name cannot be empty")
	}

	// 基本的参数名验证正则表达式
	// TR181参数名通常遵循Device.xxx或InternetGatewayDevice.xxx格式
	validPattern := regexp.MustCompile(`^(Device|InternetGatewayDevice)(\.[A-Za-z0-9_\-]+)+$`)
	if !validPattern.MatchString(name) {
		return fmt.Errorf("invalid parameter name format: %s", name)
	}

	return nil
}

// ValidateParameterValue 验证参数值的类型
func ValidateParameterValue(name string, value interface{}) error {
	if value == nil {
		return fmt.Errorf("parameter value cannot be nil")
	}

	// 根据参数名判断预期的类型
	expectedType := getExpectedParameterType(name)
	actualType := reflect.TypeOf(value).Kind()

	if expectedType != reflect.Invalid && actualType != expectedType {
		// 尝试类型转换
		if !canConvertType(value, expectedType) {
			return fmt.Errorf("parameter %s expects %s type, got %s", 
				name, expectedType, actualType)
		}
	}

	return nil
}

// getExpectedParameterType 根据参数名获取预期的类型
func getExpectedParameterType(name string) reflect.Kind {
	// 这里可以根据不同的参数名映射到预期的类型
	// 例如，所有包含"Interval"的参数可能是整数类型
	if strings.Contains(name, "Interval") || strings.Contains(name, "Uptime") || 
	   strings.Contains(name, "Timeout") || strings.Contains(name, "Port") {
		return reflect.Int
	}

	// 包含"Status"或"Enabled"的参数可能是布尔类型
	if strings.Contains(name, "Status") || strings.Contains(name, "Enabled") ||
	   strings.Contains(name, "Active") || strings.Contains(name, "Locked") {
		return reflect.Bool
	}

	// 包含"Utilization"或"Level"的参数可能是浮点数类型
	if strings.Contains(name, "Utilization") || strings.Contains(name, "Level") ||
	   strings.Contains(name, "Rate") || strings.Contains(name, "Percent") {
		return reflect.Float64
	}

	// 默认返回无效类型，表示没有特定期望
	return reflect.Invalid
}

// canConvertType 检查值是否可以转换为目标类型
func canConvertType(value interface{}, targetKind reflect.Kind) bool {
	v := reflect.ValueOf(value)

	switch targetKind {
	case reflect.Int:
		if v.Kind() == reflect.Float64 || v.Kind() == reflect.Float32 {
			return true
		}
	case reflect.Float64:
		if v.Kind() == reflect.Int || v.Kind() == reflect.Int32 || v.Kind() == reflect.Int64 {
			return true
		}
	}

	return false
}

// DeepCloneMessage 深拷贝消息对象
func DeepCloneMessage(msg *Message) (*Message, error) {
	if msg == nil {
		return nil, nil
	}

	// 使用JSON序列化和反序列化进行深拷贝
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	var cloned Message
	err = json.Unmarshal(data, &cloned)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return &cloned, nil
}

// ExtractParametersFromMessage 从消息中提取参数列表
func ExtractParametersFromMessage(msg *Message) ([]string, error) {
	var params []string

	switch body := msg.Body.(type) {
	case GetRequest:
		for _, param := range body.Parameters {
			params = append(params, param.Name)
		}
	case SetRequest:
		for _, param := range body.Parameters {
			params = append(params, param.Name)
		}
	case InformRequest:
		for _, param := range body.Parameters {
			params = append(params, param.Name)
		}
	default:
		return nil, fmt.Errorf("unsupported message body type for parameter extraction")
	}

	return params, nil
}

// BuildParameterMap 从参数值列表构建映射表
func BuildParameterMap(params []ParameterValue) map[string]interface{} {
	paramMap := make(map[string]interface{})
	for _, param := range params {
		paramMap[param.Name] = param.Value
	}
	return paramMap
}

// MapToParameterValues 将映射表转换为参数值列表
func MapToParameterValues(paramMap map[string]interface{}) []ParameterValue {
	params := make([]ParameterValue, 0, len(paramMap))
	for name, value := range paramMap {
		params = append(params, ParameterValue{
			Name:  name,
			Value: value,
		})
	}
	return params
}

// IsSubParameter 检查paramName是否是baseParam的子参数
func IsSubParameter(baseParam, paramName string) bool {
	if baseParam == paramName {
		return true
	}
	// 如果baseParam以点号结尾，直接检查paramName是否以baseParam开头
	if strings.HasSuffix(baseParam, ".") {
		return strings.HasPrefix(paramName, baseParam)
	}
	// 否则，子参数必须以基参数开头，并且后面跟着一个点
	return strings.HasPrefix(paramName, baseParam+".")
}

// GetBaseParameter 获取参数的基础路径部分
func GetBaseParameter(paramName string) string {
	parts := strings.Split(paramName, ".")
	if len(parts) >= 2 {
		return parts[0] + "." + parts[1]
	}
	return paramName
}

// EscapeParameterValue 转义参数值，确保JSON序列化安全
func EscapeParameterValue(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	// 处理字符串类型的转义
	if str, ok := value.(string); ok {
		// 确保字符串可以被JSON正确序列化
		if _, err := json.Marshal(str); err != nil {
			// 如果序列化失败，返回空字符串
			return ""
		}
	}

	return value
}

// MergeParameterMaps 合并两个参数映射表
func MergeParameterMaps(dest, src map[string]interface{}) {
	for k, v := range src {
		dest[k] = v
	}
}