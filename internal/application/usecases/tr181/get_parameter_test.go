// Package tr181 测试 TR181 用例
package tr181

import (
	"testing"

	"tr369-wss-client/internal/domain/entities/tr181"
	"tr369-wss-client/pkg/api"
)

// **Feature: code-structure-optimization, Property 23: 模块依赖接口化**
// **验证需求: 需求 6.2**
// 对于任何模块间的依赖关系，都应该通过明确的接口定义

// MockLogger 模拟日志器
type MockLogger struct{}

func (m *MockLogger) Debug(msg string, args ...interface{}) {}
func (m *MockLogger) Info(msg string, args ...interface{})  {}
func (m *MockLogger) Warn(msg string, args ...interface{})  {}
func (m *MockLogger) Error(msg string, args ...interface{}) {}
func (m *MockLogger) Fatal(msg string, args ...interface{}) {}

// MockClientRepository 模拟客户端仓储
type MockClientRepository struct {
	constructGetRespCalled bool
	handleSetRequestCalled bool
	getNewInstanceCalled   bool
	handleDeleteCalled     bool
	saveDataCalled         bool
}

func (m *MockClientRepository) GetValueByPath(path string) (interface{}, error) {
	return nil, nil
}

func (m *MockClientRepository) ConstructGetResp(paths []string) api.Response_GetResp {
	m.constructGetRespCalled = true
	return api.Response_GetResp{}
}

func (m *MockClientRepository) HandleSetRequest(path, key, value string) {
	m.handleSetRequestCalled = true
}

func (m *MockClientRepository) HandleDeleteRequest(path string) (string, bool) {
	m.handleDeleteCalled = true
	return path, true
}

func (m *MockClientRepository) GetNewInstance(path string) string {
	m.getNewInstanceCalled = true
	return path + "1."
}

func (m *MockClientRepository) IsExistPath(path string) (bool, string) {
	return true, path
}

func (m *MockClientRepository) SaveData() {
	m.saveDataCalled = true
}

func (m *MockClientRepository) StartClientRepository() {}

func (m *MockClientRepository) AddListener(paramName string, listener tr181.Listener) error {
	return nil
}

func (m *MockClientRepository) RemoveListener(paramName string) error {
	return nil
}

func (m *MockClientRepository) ResetListeners() error {
	return nil
}

func (m *MockClientRepository) NotifyListeners(paramName string, value interface{}) {}

// TestGetParameterUseCaseWithMock 测试 GET 参数用例使用模拟对象
func TestGetParameterUseCaseWithMock(t *testing.T) {
	mockRepo := &MockClientRepository{}
	mockLogger := &MockLogger{}

	uc := NewGetParameterUseCase(mockRepo, mockLogger)

	// 执行用例
	uc.Execute([]string{"Device.DeviceInfo."})

	// 验证仓储方法被调用
	if !mockRepo.constructGetRespCalled {
		t.Error("期望 ConstructGetResp 被调用")
	}
}

// TestSetParameterUseCaseWithMock 测试 SET 参数用例使用模拟对象
func TestSetParameterUseCaseWithMock(t *testing.T) {
	mockRepo := &MockClientRepository{}
	mockLogger := &MockLogger{}

	uc := NewSetParameterUseCase(mockRepo, mockLogger)

	// 执行用例
	uc.Execute("Device.DeviceInfo.", "Manufacturer", "TestValue")

	// 验证仓储方法被调用
	if !mockRepo.handleSetRequestCalled {
		t.Error("期望 HandleSetRequest 被调用")
	}
}

// TestAddObjectUseCaseWithMock 测试 ADD 对象用例使用模拟对象
func TestAddObjectUseCaseWithMock(t *testing.T) {
	mockRepo := &MockClientRepository{}
	mockLogger := &MockLogger{}

	uc := NewAddObjectUseCase(mockRepo, mockLogger)

	// 执行用例
	result := uc.Execute("Device.LocalAgent.Subscription.")

	// 验证仓储方法被调用
	if !mockRepo.getNewInstanceCalled {
		t.Error("期望 GetNewInstance 被调用")
	}

	// 验证返回值
	if result != "Device.LocalAgent.Subscription.1." {
		t.Errorf("期望返回 'Device.LocalAgent.Subscription.1.', 实际返回 '%s'", result)
	}
}

// TestDeleteObjectUseCaseWithMock 测试 DELETE 对象用例使用模拟对象
func TestDeleteObjectUseCaseWithMock(t *testing.T) {
	mockRepo := &MockClientRepository{}
	mockLogger := &MockLogger{}

	uc := NewDeleteObjectUseCase(mockRepo, mockLogger)

	// 执行用例
	path, found := uc.Execute("Device.LocalAgent.Subscription.1.")

	// 验证仓储方法被调用
	if !mockRepo.handleDeleteCalled {
		t.Error("期望 HandleDeleteRequest 被调用")
	}

	// 验证返回值
	if !found {
		t.Error("期望找到路径")
	}

	if path != "Device.LocalAgent.Subscription.1." {
		t.Errorf("期望返回 'Device.LocalAgent.Subscription.1.', 实际返回 '%s'", path)
	}
}

// TestUseCaseDependsOnInterface 验证用例依赖于接口而不是具体实现
func TestUseCaseDependsOnInterface(t *testing.T) {
	// 这个测试验证我们可以使用任何实现了接口的对象
	// 如果用例依赖于具体实现，这个测试将无法编译

	var repo interface{} = &MockClientRepository{}

	// 验证 MockClientRepository 实现了所需的接口方法
	if _, ok := repo.(interface {
		ConstructGetResp(paths []string) api.Response_GetResp
	}); !ok {
		t.Error("MockClientRepository 应该实现 ConstructGetResp 方法")
	}

	if _, ok := repo.(interface {
		HandleSetRequest(path, key, value string)
	}); !ok {
		t.Error("MockClientRepository 应该实现 HandleSetRequest 方法")
	}
}
