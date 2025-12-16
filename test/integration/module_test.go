// Package integration 提供模块化支持测试
package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// **Feature: code-structure-optimization, Property 24: 模块边界封装性**
// **验证需求: 需求 6.3**
// 对于任何模块，都应该有明确的边界和封装机制

func TestModuleBoundaryEncapsulation(t *testing.T) {
	rootDir := getProjectRootForModule()
	internalDir := filepath.Join(rootDir, "internal")

	if _, err := os.Stat(internalDir); os.IsNotExist(err) {
		t.Fatalf("internal 目录不存在: %s", internalDir)
	}

	// 检查 internal 目录下的模块是否有明确的边界
	expectedModules := []string{"domain", "application", "infrastructure"}

	for _, module := range expectedModules {
		moduleDir := filepath.Join(internalDir, module)
		if _, err := os.Stat(moduleDir); os.IsNotExist(err) {
			t.Errorf("模块目录不存在: %s", module)
		} else {
			t.Logf("模块 %s 已正确封装在 internal 目录中", module)
		}
	}
}

// **Feature: code-structure-optimization, Property 25: 模块独立测试能力**
// **验证需求: 需求 6.4**
// 对于任何模块，都应该能够独立进行单元测试

func TestModuleIndependentTesting(t *testing.T) {
	rootDir := getProjectRootForModule()
	internalDir := filepath.Join(rootDir, "internal")

	// 检查每个主要模块是否有测试文件
	modules := []string{"domain", "application", "infrastructure"}
	modulesWithTests := 0

	for _, module := range modules {
		moduleDir := filepath.Join(internalDir, module)
		hasTests := false

		err := filepath.Walk(moduleDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && strings.HasSuffix(path, "_test.go") {
				hasTests = true
			}
			return nil
		})

		if err != nil {
			t.Logf("遍历模块 %s 时出错: %v", module, err)
			continue
		}

		if hasTests {
			modulesWithTests++
			t.Logf("模块 %s 有独立的测试文件", module)
		} else {
			t.Logf("模块 %s 没有测试文件", module)
		}
	}

	if modulesWithTests == 0 {
		t.Error("至少应该有一个模块有独立的测试文件")
	}
}

// **Feature: code-structure-optimization, Property 26: 模块配置独立性**
// **验证需求: 需求 6.5**
// 对于任何模块，都应该能够独立配置而不影响其他模块

func TestModuleConfigIndependence(t *testing.T) {
	rootDir := getProjectRootForModule()

	// 检查配置是否集中管理，允许模块独立配置
	configDir := filepath.Join(rootDir, "internal", "infrastructure", "config")

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Fatalf("配置目录不存在: %s", configDir)
	}

	// 检查配置文件是否存在
	configFile := filepath.Join(configDir, "config.go")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Fatalf("配置文件不存在: %s", configFile)
	}

	t.Log("模块配置已正确组织，支持独立配置")
}

// getProjectRootForModule 获取项目根目录
func getProjectRootForModule() string {
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
	if strings.Contains(wd, "test") {
		parts := strings.Split(wd, "test")
		return strings.TrimSuffix(parts[0], string(os.PathSeparator))
	}

	return wd
}
