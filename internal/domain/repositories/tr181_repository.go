// Package repositories 定义领域层的仓储接口
package repositories

import (
	"tr369-wss-client/internal/domain/entities/tr181"
)

// TR181Repository 定义 TR181 数据模型仓储接口
// 该接口定义了对 TR181 数据模型的所有操作
type TR181Repository interface {
	// GetValueByPath 根据路径获取参数值
	GetValueByPath(path string) (interface{}, error)

	// SetValue 设置参数值
	SetValue(path, key, value string) error

	// IsExistPath 检查路径是否存在
	IsExistPath(path string) (bool, string)

	// GetNewInstance 获取新实例路径
	GetNewInstance(path string) string

	// DeleteObject 删除对象
	DeleteObject(path string) (string, bool)

	// SaveData 保存数据到持久化存储
	SaveData() error

	// LoadData 从持久化存储加载数据
	LoadData() error

	// AddListener 添加监听器
	AddListener(paramName string, listener tr181.Listener) error

	// RemoveListener 移除监听器
	RemoveListener(paramName string) error

	// ResetListeners 重置所有监听器
	ResetListeners() error

	// NotifyListeners 通知监听器
	NotifyListeners(paramName string, value interface{})
}
