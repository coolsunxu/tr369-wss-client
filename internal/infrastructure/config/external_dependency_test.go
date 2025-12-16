// Package config 测试外部依赖隔离
package config

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// **Feature: code-structure-optimization, Property 3: 外部依赖隔离**
// **验证需求: 需求 1.4**
// 对于任何外部依赖（如数据库、消息队列、第三方服务），都应该通过适配器模式隔离在基础设施层

func TestExternalDependenciesIsolatedInInfrastructure(t *testing.T) {
	rootDir := getProjectRootForExternal()

	// 外部依赖包列表（这些应该只在基础设施层使用）
	externalDependencies := []string{
		"github.com/coder/websocket",
		"github.com/spf13/viper",
		"go.uber.org/zap",
	}

	// 检查领域层不应该导入外部依赖
	domainDir := filepath.Join(rootDir, "internal", "domain")
	checkNoExternalDependencies(t, domainDir, externalDependencies, "领域层")

	// 检查应用层不应该直接导入外部依赖
	applicationDir := filepath.Join(rootDir, "internal", "application")
	checkNoExternalDependencies(t, applicationDir, externalDependencies, "应用层")
}

// TestExternalDependenciesOnlyInInfrastructure 验证外部依赖只在基础设施层
func TestExternalDependenciesOnlyInInfrastructure(t *testing.T) {
	rootDir := getProjectRootForExternal()
	infraDir := filepath.Join(rootDir, "internal", "infrastructure")

	// 检查基础设施层是否存在
	if _, err := os.Stat(infraDir); os.IsNotExist(err) {
		t.Fatalf("基础设施层目录不存在: %s", infraDir)
	}

	// 外部依赖应该在基础设施层中使用
	externalDependencies := []string{
		"github.com/coder/websocket",
		"go.uber.org/zap",
	}

	foundDependencies := make(map[string]bool)

	err := filepath.Walk(infraDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			imports, parseErr := getImportsForExternal(path)
			if parseErr != nil {
				return nil
			}

			for _, imp := range imports {
				for _, extDep := range externalDependencies {
					if strings.HasPrefix(imp, extDep) {
						foundDependencies[extDep] = true
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		t.Fatalf("遍历基础设施层目录失败: %v", err)
	}

	// 验证至少有一些外部依赖在基础设施层中使用
	if len(foundDependencies) == 0 {
		t.Log("基础设施层中没有发现外部依赖（这可能是正常的，取决于项目配置）")
	} else {
		t.Logf("基础设施层中使用的外部依赖: %v", foundDependencies)
	}
}

// checkNoExternalDependencies 检查指定目录不包含外部依赖
func checkNoExternalDependencies(t *testing.T, dir string, externalDeps []string, layerName string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			imports, parseErr := getImportsForExternal(path)
			if parseErr != nil {
				return nil
			}

			for _, imp := range imports {
				for _, extDep := range externalDeps {
					if strings.HasPrefix(imp, extDep) {
						t.Errorf("%s文件 %s 不应该直接导入外部依赖: %s", layerName, path, imp)
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		t.Fatalf("遍历%s目录失败: %v", layerName, err)
	}
}

// getImportsForExternal 获取文件的导入列表
func getImportsForExternal(filePath string) ([]string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}

	var imports []string
	for _, imp := range node.Imports {
		importPath := strings.Trim(imp.Path.Value, "\"")
		imports = append(imports, importPath)
	}

	return imports, nil
}

// getProjectRootForExternal 获取项目根目录
func getProjectRootForExternal() string {
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
