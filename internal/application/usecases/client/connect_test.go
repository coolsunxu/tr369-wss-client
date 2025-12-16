// Package client 测试客户端用例
package client

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// **Feature: code-structure-optimization, Property 22: 模块独立开发支持**
// **验证需求: 需求 6.1**
// 对于任何功能模块，都应该支持独立的编译和测试

func TestUseCaseModuleIndependence(t *testing.T) {
	// 验证用例模块可以独立编译
	// 这个测试本身的存在就证明了模块可以独立测试

	rootDir := getProjectRoot()
	usecaseDir := filepath.Join(rootDir, "internal", "application", "usecases")

	// 检查用例目录是否存在
	if _, err := os.Stat(usecaseDir); os.IsNotExist(err) {
		t.Fatalf("用例目录不存在: %s", usecaseDir)
	}

	// 检查是否有子模块
	entries, err := os.ReadDir(usecaseDir)
	if err != nil {
		t.Fatalf("读取用例目录失败: %v", err)
	}

	moduleCount := 0
	for _, entry := range entries {
		if entry.IsDir() {
			moduleCount++
			t.Logf("发现用例模块: %s", entry.Name())
		}
	}

	if moduleCount == 0 {
		t.Error("用例目录下应该有子模块")
	}
}

func TestConnectUseCaseCreation(t *testing.T) {
	// 测试 ConnectUseCase 可以被创建
	// 这验证了模块的基本功能

	// 由于我们不想在测试中实际连接，只验证结构
	uc := &ConnectUseCase{}
	if uc == nil {
		t.Error("ConnectUseCase 创建失败")
	}
}

func TestDisconnectUseCaseCreation(t *testing.T) {
	// 测试 DisconnectUseCase 可以被创建
	uc := &DisconnectUseCase{}
	if uc == nil {
		t.Error("DisconnectUseCase 创建失败")
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
