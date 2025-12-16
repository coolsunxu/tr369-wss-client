// Package json 提供 JSON 处理工具
package json

import (
	"encoding/json"
)

// SafeMarshal 安全地将对象序列化为 JSON 字符串
// 如果序列化失败，返回空字符串
func SafeMarshal(v interface{}) string {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

// SafeUnmarshal 安全地将 JSON 字符串反序列化为对象
// 如果反序列化失败，返回错误
func SafeUnmarshal(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

// PrettyMarshal 将对象序列化为格式化的 JSON 字符串
func PrettyMarshal(v interface{}) (string, error) {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
