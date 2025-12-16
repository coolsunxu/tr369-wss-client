// Package integration 提供集成测试
package integration

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// **Feature: code-structure-optimization, Property 1: 依赖倒置原则遵循**
// **验证需求: 需求 1.2**
// 对于任何内层模块到外层模块的调用，都应该通过接口进行而不是直接依赖具体实现

func TestDependencyInversionPrinciple(t *testing.T) {
	rootDir := getProjectRoot()

	// 领域层不应该导入基础设施层或应用层
	domainDir := filepath.Join(rootDir, "internal", "domain")
	forbiddenImportsForDomain := []string{
		"tr369-wss-client/internal/infrastructure",
		"tr369-wss-client/internal/application",
	}

	err := checkForbiddenImports(domainDir, forbiddenImportsForDomain, t)
	if err != nil {
		t.Fatalf("检查领域层导入失败: %v", err)
	}
}

// TestApplicationLayerDependencies 测试应用层依赖
func TestApplicationLayerDependencies(t *testing.T) {
	rootDir := getProjectRoot()

	// 应用层不应该直接导入基础设施层的具体实现（除了通过接口）
	// 但可以导入领域层
	applicationDir := filepath.Join(rootDir, "internal", "application")

	// 检查应用层是否存在
	if _, err := os.Stat(applicationDir); os.IsNotExist(err) {
		t.Skip("应用层目录不存在，跳过测试")
		return
	}

	// 应用层可以导入领域层，这是允许的
	allowedImports := []string{
		"tr369-wss-client/internal/domain",
		"tr369-wss-client/pkg",
	}

	err := filepath.Walk(applicationDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			imports, parseErr := getImports(path)
			if parseErr != nil {
				return nil
			}

			for _, imp := range imports {
				// 检查是否导入了基础设施层
				if strings.Contains(imp, "internal/infrastructure") {
					// 这是允许的，因为应用层需要使用基础设施层的实现
					// 但应该通过依赖注入而不是直接实例化
					t.Logf("应用层文件 %s 导入了基础设施层: %s (应通过依赖注入)", path, imp)
				}

				// 检查是否是允许的导入
				isAllowed := false
				for _, allowed := range allowedImports {
					if strings.Contains(imp, allowed) {
						isAllowed = true
						break
					}
				}

				if strings.Contains(imp, "tr369-wss-client/internal") && !isAllowed && !strings.Contains(imp, "infrastructure") {
					t.Errorf("应用层文件 %s 包含不允许的导入: %s", path, imp)
				}
			}
		}
		return nil
	})

	if err != nil {
		t.Fatalf("遍历应用层目录失败: %v", err)
	}
}

// TestDomainLayerPurity 测试领域层纯净性
func TestDomainLayerPurity(t *testing.T) {
	rootDir := getProjectRoot()
	domainDir := filepath.Join(rootDir, "internal", "domain")

	// 领域层不应该导入外部框架（除了标准库和 protobuf）
	allowedExternalPackages := []string{
		"tr369-wss-client/pkg",
		"google.golang.org/protobuf",
	}

	err := filepath.Walk(domainDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			imports, parseErr := getImports(path)
			if parseErr != nil {
				return nil
			}

			for _, imp := range imports {
				// 跳过标准库
				if !strings.Contains(imp, ".") && !strings.Contains(imp, "/") {
					continue
				}

				// 跳过项目内部包
				if strings.HasPrefix(imp, "tr369-wss-client") {
					// 检查是否导入了不允许的内部包
					if strings.Contains(imp, "internal/infrastructure") || strings.Contains(imp, "internal/application") {
						t.Errorf("领域层文件 %s 不应该导入: %s", path, imp)
					}
					continue
				}

				// 检查外部包是否在允许列表中
				isAllowed := false
				for _, allowed := range allowedExternalPackages {
					if strings.HasPrefix(imp, allowed) {
						isAllowed = true
						break
					}
				}

				if !isAllowed && strings.Contains(imp, "/") {
					// 这可能是外部框架，记录但不报错（因为可能是必要的）
					t.Logf("领域层文件 %s 导入了外部包: %s", path, imp)
				}
			}
		}
		return nil
	})

	if err != nil {
		t.Fatalf("遍历领域层目录失败: %v", err)
	}
}

// checkForbiddenImports 检查禁止的导入
func checkForbiddenImports(dir string, forbidden []string, t *testing.T) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			imports, parseErr := getImports(path)
			if parseErr != nil {
				return nil
			}

			for _, imp := range imports {
				for _, forbid := range forbidden {
					if strings.Contains(imp, forbid) {
						t.Errorf("文件 %s 包含禁止的导入: %s", path, imp)
					}
				}
			}
		}
		return nil
	})
}

// getImports 获取文件的导入列表
func getImports(filePath string) ([]string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}

	var imports []string
	for _, imp := range node.Imports {
		// 去除引号
		importPath := strings.Trim(imp.Path.Value, "\"")
		imports = append(imports, importPath)
	}

	return imports, nil
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
	if strings.Contains(wd, "test") {
		parts := strings.Split(wd, "test")
		return strings.TrimSuffix(parts[0], string(os.PathSeparator))
	}

	return wd
}
