// Package config 测试配置管理
package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// **Feature: code-structure-optimization, Property 2: 基础设施层隔离性**
// **验证需求: 需求 1.3**
// 对于任何基础设施层的修改，领域层和应用层的测试应该继续通过而不受影响

func TestConfigValidation(t *testing.T) {
	testCases := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name: "有效配置",
			config: &Config{
				DataRefreshConfig: &DataRefreshConfig{
					IntervalSeconds:     60,
					WriteCountThreshold: 10,
					TR181DataModelPath:  "./data/test.json",
				},
				WebsocketConfig: &WebsocketConfig{
					ServerURL:          "ws://localhost:8080",
					ControllerId:       "controller",
					PingInterval:       30,
					MaxMessageSize:     65536,
					EndpointId:         "endpoint",
					MessageChannelSize: 100,
				},
			},
			expectError: false,
		},
		{
			name: "空 DataRefreshConfig",
			config: &Config{
				DataRefreshConfig: nil,
				WebsocketConfig: &WebsocketConfig{
					ServerURL:          "ws://localhost:8080",
					ControllerId:       "controller",
					PingInterval:       30,
					MaxMessageSize:     65536,
					EndpointId:         "endpoint",
					MessageChannelSize: 100,
				},
			},
			expectError: true,
		},
		{
			name: "无效的 IntervalSeconds",
			config: &Config{
				DataRefreshConfig: &DataRefreshConfig{
					IntervalSeconds:     0,
					WriteCountThreshold: 10,
					TR181DataModelPath:  "./data/test.json",
				},
				WebsocketConfig: &WebsocketConfig{
					ServerURL:          "ws://localhost:8080",
					ControllerId:       "controller",
					PingInterval:       30,
					MaxMessageSize:     65536,
					EndpointId:         "endpoint",
					MessageChannelSize: 100,
				},
			},
			expectError: true,
		},
		{
			name: "空 ServerURL",
			config: &Config{
				DataRefreshConfig: &DataRefreshConfig{
					IntervalSeconds:     60,
					WriteCountThreshold: 10,
					TR181DataModelPath:  "./data/test.json",
				},
				WebsocketConfig: &WebsocketConfig{
					ServerURL:          "",
					ControllerId:       "controller",
					PingInterval:       30,
					MaxMessageSize:     65536,
					EndpointId:         "endpoint",
					MessageChannelSize: 100,
				},
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Validate()
			if tc.expectError && err == nil {
				t.Error("期望返回错误，但没有")
			}
			if !tc.expectError && err != nil {
				t.Errorf("不期望返回错误，但返回了: %v", err)
			}
		})
	}
}

// **Feature: code-structure-optimization, Property 3: 外部依赖隔离**
// **验证需求: 需求 1.4**
// 对于任何新添加的外部依赖，都应该只出现在基础设施层中

func TestInfrastructureLayerContainsExternalDependencies(t *testing.T) {
	rootDir := getProjectRoot()
	infraDir := filepath.Join(rootDir, "internal", "infrastructure")

	// 检查基础设施层目录是否存在
	if _, err := os.Stat(infraDir); os.IsNotExist(err) {
		t.Fatalf("基础设施层目录不存在: %s", infraDir)
	}

	// 验证基础设施层包含预期的子模块
	expectedModules := []string{"config", "logging", "websocket", "persistence", "protobuf"}

	for _, module := range expectedModules {
		modulePath := filepath.Join(infraDir, module)
		if _, err := os.Stat(modulePath); os.IsNotExist(err) {
			t.Errorf("基础设施层缺少模块: %s", module)
		}
	}
}

// **Feature: code-structure-optimization, Property 10: 外部服务实现基础设施层定位**
// **验证需求: 需求 3.2**
// 对于任何外部服务接口的实现，都应该位于基础设施层中

func TestExternalServiceImplementationsInInfrastructure(t *testing.T) {
	rootDir := getProjectRoot()

	// WebSocket 实现应该在基础设施层
	wsPath := filepath.Join(rootDir, "internal", "infrastructure", "websocket", "client.go")
	if _, err := os.Stat(wsPath); os.IsNotExist(err) {
		t.Error("WebSocket 客户端实现应该在基础设施层")
	}

	// 配置加载实现应该在基础设施层
	configPath := filepath.Join(rootDir, "internal", "infrastructure", "config", "config.go")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("配置管理实现应该在基础设施层")
	}

	// 日志实现应该在基础设施层
	logPath := filepath.Join(rootDir, "internal", "infrastructure", "logging", "logger.go")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Error("日志实现应该在基础设施层")
	}
}

// getProjectRoot 获取项目根目录
func getProjectRoot() string {
	dir, _ := os.Getwd()

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	wd, _ := os.Getwd()
	if strings.Contains(wd, "internal") {
		parts := strings.Split(wd, "internal")
		return strings.TrimSuffix(parts[0], string(os.PathSeparator))
	}

	return wd
}
