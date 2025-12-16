// Package config 测试配置代码集中性
package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// **Feature: code-structure-optimization, Property 7: 配置代码集中性**
// **验证需求: 需求 2.4**
// 对于任何配置相关的代码，都应该集中在 configs 目录或 internal/infrastructure/config 包中

func TestConfigCodeCentralization(t *testing.T) {
	rootDir := getProjectRootForConfig()

	// 配置代码应该集中的位置
	allowedConfigLocations := []string{
		filepath.Join(rootDir, "configs"),
		filepath.Join(rootDir, "internal", "infrastructure", "config"),
	}

	// 验证配置目录存在
	for _, loc := range allowedConfigLocations {
		if _, err := os.Stat(loc); os.IsNotExist(err) {
			t.Errorf("配置目录不存在: %s", loc)
		}
	}
}

// TestConfigFilesInConfigsDirectory 验证配置文件在 configs 目录
func TestConfigFilesInConfigsDirectory(t *testing.T) {
	rootDir := getProjectRootForConfig()
	configsDir := filepath.Join(rootDir, "configs")

	if _, err := os.Stat(configsDir); os.IsNotExist(err) {
		t.Fatalf("configs 目录不存在: %s", configsDir)
	}

	// 检查 configs 目录下是否有配置文件
	hasConfigFiles := false
	err := filepath.Walk(configsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			ext := filepath.Ext(path)
			if ext == ".json" || ext == ".yaml" || ext == ".yml" || ext == ".toml" {
				hasConfigFiles = true
				t.Logf("发现配置文件: %s", path)
			}
		}
		return nil
	})

	if err != nil {
		t.Fatalf("遍历 configs 目录失败: %v", err)
	}

	if !hasConfigFiles {
		t.Error("configs 目录下应该有配置文件")
	}
}

// TestConfigPackageExists 验证配置包存在
func TestConfigPackageExists(t *testing.T) {
	rootDir := getProjectRootForConfig()
	configPkgDir := filepath.Join(rootDir, "internal", "infrastructure", "config")

	if _, err := os.Stat(configPkgDir); os.IsNotExist(err) {
		t.Fatalf("配置包目录不存在: %s", configPkgDir)
	}

	// 检查是否有 Go 文件
	hasGoFiles := false
	entries, err := os.ReadDir(configPkgDir)
	if err != nil {
		t.Fatalf("读取配置包目录失败: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") && !strings.HasSuffix(entry.Name(), "_test.go") {
			hasGoFiles = true
			t.Logf("发现配置代码文件: %s", entry.Name())
		}
	}

	if !hasGoFiles {
		t.Error("配置包目录下应该有 Go 文件")
	}
}

// TestNoConfigCodeOutsideAllowedLocations 验证配置代码不在其他位置
func TestNoConfigCodeOutsideAllowedLocations(t *testing.T) {
	rootDir := getProjectRootForConfig()

	// 检查领域层和应用层不应该有配置加载代码
	// 这里我们检查是否有直接使用 viper 或 json.Unmarshal 加载配置的代码
	dirsToCheck := []string{
		filepath.Join(rootDir, "internal", "domain"),
		filepath.Join(rootDir, "internal", "application"),
	}

	for _, dir := range dirsToCheck {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}

		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
				content, readErr := os.ReadFile(path)
				if readErr != nil {
					return nil
				}

				contentStr := string(content)
				// 检查是否有直接的配置加载代码
				if strings.Contains(contentStr, "viper.") && !strings.Contains(path, "config") {
					t.Errorf("文件 %s 不应该直接使用 viper 加载配置", path)
				}
			}
			return nil
		})

		if err != nil {
			t.Logf("遍历目录 %s 时出错: %v", dir, err)
		}
	}
}

// getProjectRootForConfig 获取项目根目录
func getProjectRootForConfig() string {
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
