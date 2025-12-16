// Package tr181 测试 TR181 数据模型
package tr181

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// **Feature: code-structure-optimization, Property 6: 目录结构一致性**
// **验证需求: 需求 2.3**
// 对于任何功能模块，都应该遵循统一的目录结构模式

func TestDirectoryStructureConsistency(t *testing.T) {
	// 定义期望的目录结构
	expectedDirs := []string{
		"cmd/client",
		"internal/domain/entities",
		"internal/domain/repositories",
		"internal/domain/services",
		"internal/domain/valueobjects",
		"internal/application/usecases",
		"internal/infrastructure/config",
		"internal/infrastructure/logging",
		"internal/infrastructure/websocket",
		"internal/infrastructure/persistence",
		"internal/infrastructure/protobuf",
		"pkg/api",
		"pkg/errors",
		"pkg/utils",
		"configs/environments",
		"test/fixtures",
		"test/mocks",
		"docs/architecture",
		"scripts",
		"examples",
	}

	// 获取项目根目录
	rootDir := getProjectRoot()

	for _, dir := range expectedDirs {
		fullPath := filepath.Join(rootDir, dir)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("期望的目录不存在: %s", dir)
		}
	}
}

// TestInternalPackageStructure 测试 internal 包结构
func TestInternalPackageStructure(t *testing.T) {
	rootDir := getProjectRoot()
	internalDir := filepath.Join(rootDir, "internal")

	// 检查 internal 目录下的四层架构
	layers := []string{"domain", "application", "infrastructure"}

	for _, layer := range layers {
		layerPath := filepath.Join(internalDir, layer)
		if _, err := os.Stat(layerPath); os.IsNotExist(err) {
			t.Errorf("internal 目录下缺少层级: %s", layer)
		}
	}
}

// TestDomainLayerStructure 测试领域层结构
func TestDomainLayerStructure(t *testing.T) {
	rootDir := getProjectRoot()
	domainDir := filepath.Join(rootDir, "internal", "domain")

	// 领域层应该包含的子目录
	subDirs := []string{"entities", "repositories", "services", "valueobjects"}

	for _, subDir := range subDirs {
		subDirPath := filepath.Join(domainDir, subDir)
		if _, err := os.Stat(subDirPath); os.IsNotExist(err) {
			t.Errorf("领域层缺少子目录: %s", subDir)
		}
	}
}

// TestConfigEnvironmentsSeparation 测试配置环境分离
func TestConfigEnvironmentsSeparation(t *testing.T) {
	rootDir := getProjectRoot()
	envDir := filepath.Join(rootDir, "configs", "environments")

	// 应该存在不同环境的配置文件
	envFiles := []string{"development.json", "production.json", "test.json"}

	for _, envFile := range envFiles {
		filePath := filepath.Join(envDir, envFile)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("缺少环境配置文件: %s", envFile)
		}
	}
}

// getProjectRoot 获取项目根目录
func getProjectRoot() string {
	// 从当前测试文件位置向上查找项目根目录
	dir, _ := os.Getwd()

	// 向上查找直到找到 go.mod 文件
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// 已经到达根目录
			break
		}
		dir = parent
	}

	// 如果找不到，尝试从测试路径推断
	wd, _ := os.Getwd()
	// 假设测试在 internal/domain/entities/tr181 目录下运行
	if strings.Contains(wd, "internal") {
		parts := strings.Split(wd, "internal")
		return strings.TrimSuffix(parts[0], string(os.PathSeparator))
	}

	return wd
}
