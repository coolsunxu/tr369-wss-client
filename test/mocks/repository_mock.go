// Package mocks 提供测试用的模拟对象
package mocks

import (
	"tr369-wss-client/internal/domain/entities/tr181"
	"tr369-wss-client/pkg/api"
)

// MockClientRepository 模拟客户端仓储
type MockClientRepository struct {
	GetValueByPathFunc      func(path string) (interface{}, error)
	ConstructGetRespFunc    func(paths []string) api.Response_GetResp
	HandleSetRequestFunc    func(path, key, value string)
	HandleDeleteRequestFunc func(path string) (string, bool)
	GetNewInstanceFunc      func(path string) string
	IsExistPathFunc         func(path string) (bool, string)
	SaveDataFunc            func()
	StartFunc               func()
	AddListenerFunc         func(paramName string, listener tr181.Listener) error
	RemoveListenerFunc      func(paramName string) error
	ResetListenersFunc      func() error
	NotifyListenersFunc     func(paramName string, value interface{})
}

// GetValueByPath 模拟获取路径值
func (m *MockClientRepository) GetValueByPath(path string) (interface{}, error) {
	if m.GetValueByPathFunc != nil {
		return m.GetValueByPathFunc(path)
	}
	return nil, nil
}

// ConstructGetResp 模拟构建 GET 响应
func (m *MockClientRepository) ConstructGetResp(paths []string) api.Response_GetResp {
	if m.ConstructGetRespFunc != nil {
		return m.ConstructGetRespFunc(paths)
	}
	return api.Response_GetResp{}
}

// HandleSetRequest 模拟处理 SET 请求
func (m *MockClientRepository) HandleSetRequest(path, key, value string) {
	if m.HandleSetRequestFunc != nil {
		m.HandleSetRequestFunc(path, key, value)
	}
}

// HandleDeleteRequest 模拟处理 DELETE 请求
func (m *MockClientRepository) HandleDeleteRequest(path string) (string, bool) {
	if m.HandleDeleteRequestFunc != nil {
		return m.HandleDeleteRequestFunc(path)
	}
	return "", false
}

// GetNewInstance 模拟获取新实例
func (m *MockClientRepository) GetNewInstance(path string) string {
	if m.GetNewInstanceFunc != nil {
		return m.GetNewInstanceFunc(path)
	}
	return path + "1."
}

// IsExistPath 模拟检查路径是否存在
func (m *MockClientRepository) IsExistPath(path string) (bool, string) {
	if m.IsExistPathFunc != nil {
		return m.IsExistPathFunc(path)
	}
	return true, path
}

// SaveData 模拟保存数据
func (m *MockClientRepository) SaveData() {
	if m.SaveDataFunc != nil {
		m.SaveDataFunc()
	}
}

// StartClientRepository 模拟启动仓储
func (m *MockClientRepository) StartClientRepository() {
	if m.StartFunc != nil {
		m.StartFunc()
	}
}

// AddListener 模拟添加监听器
func (m *MockClientRepository) AddListener(paramName string, listener tr181.Listener) error {
	if m.AddListenerFunc != nil {
		return m.AddListenerFunc(paramName, listener)
	}
	return nil
}

// RemoveListener 模拟移除监听器
func (m *MockClientRepository) RemoveListener(paramName string) error {
	if m.RemoveListenerFunc != nil {
		return m.RemoveListenerFunc(paramName)
	}
	return nil
}

// ResetListeners 模拟重置监听器
func (m *MockClientRepository) ResetListeners() error {
	if m.ResetListenersFunc != nil {
		return m.ResetListenersFunc()
	}
	return nil
}

// NotifyListeners 模拟通知监听器
func (m *MockClientRepository) NotifyListeners(paramName string, value interface{}) {
	if m.NotifyListenersFunc != nil {
		m.NotifyListenersFunc(paramName, value)
	}
}
