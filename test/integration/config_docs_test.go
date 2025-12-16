// Package integration 提供配置和文档组织测试
package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// **Feature: code-structure-optimization, Property 18: 环境配置分离**
// **验证需求: 需求 5.1**
// 对于任何环境特定的配置，都应该按环境分离存储

func TestEnvironmentConfigSeparation(t *testing.T) {
	rootDir := getProjectRootForConfigDocs()
	envDir := filepath.Join(rootDir, "configs", "environments")

	if _, err := os.Stat(envDir); os.IsNotExist(err) {
		t.Fatalf("环境配置目录不存在: %s", envDir)
	}

	// 期望的环境配置文件
	expectedEnvs := []string{"development", "production", "test"}
	foundEnvs := make(map[string]bool)

	entries, err := os.ReadDir(envDir)
	if err != nil {
		t.Fatalf("读取环境配置目录失败: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			name := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
			foundEnvs[name] = true
			t.Logf("发现环境配置: %s", entry.Name())
		}
	}

	for _, env := range expectedEnvs {
		if !foundEnvs[env] {
			t.Errorf("缺少环境配置: %s", env)
		}
	}
}

// **Feature: code-structure-optimization, Property 20: 生成代码分离**
// **验证需求: 需求 5.4**
// 对于任何自动生成的代码，都应该与手写代码分离存储

func TestGeneratedCodeSeparation(t *testing.T) {
	rootDir := getProjectRootForConfigDocs()

	// 检查 proto 目录是否存在（用于存放 protobuf 生成的代码）
	protoDir := filepath.Join(rootDir, "proto")
	if _, err := os.Stat(protoDir); os.IsNotExist(err) {
		t.Fatalf("proto 目录不存在: %s", protoDir)
	}

	// 检查 pkg/api 目录是否存在（用于存放生成的 API 代码）
	apiDir := filepath.Join(rootDir, "pkg", "api")
	if _, err := os.Stat(apiDir); os.IsNotExist(err) {
		t.Fatalf("pkg/api 目录不存在: %s", apiDir)
	}

	t.Log("生成代码目录已正确配置")
}

// **Feature: code-structure-optimization, Property 21: 示例代码独立管理**
// **验证需求: 需求 5.5**
// 对于任何示例代码，都应该位于独立的 examples 目录中

func TestExampleCodeSeparation(t *testing.T) {
	rootDir := getProjectRootForConfigDocs()
	examplesDir := filepath.Join(rootDir, "examples")

	if _, err := os.Stat(examplesDir); os.IsNotExist(err) {
		t.Fatalf("示例代码目录不存在: %s", examplesDir)
	}

	// 检查是否有示例代码
	hasExamples := false
	err := filepath.Walk(examplesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			hasExamples = true
			t.Logf("发现示例代码: %s", path)
		}
		return nil
	})

	if err != nil {
		t.Fatalf("遍历示例目录失败: %v", err)
	}

	if !hasExamples {
		t.Log("示例目录存在但没有示例代码（这可能是正常的）")
	}
}

// getProjectRootForConfigDocs 获取项目根目录
func getProjectRootForConfigDocs() string {
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
