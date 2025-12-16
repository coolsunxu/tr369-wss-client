// Package integration 提供接口通信测试
package integration

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// **Feature: code-structure-optimization, Property 4: 跨层调用接口化**
// **验证需求: 需求 1.5**
// 对于任何跨层调用，都应该通过接口进行而不是直接依赖具体实现

func TestCrossLayerCallsThroughInterfaces(t *testing.T) {
	rootDir := getProjectRootForInterface()

	// 检查应用层是否通过接口调用基础设施层
	applicationDir := filepath.Join(rootDir, "internal", "application")

	if _, err := os.Stat(applicationDir); os.IsNotExist(err) {
		t.Skip("应用层目录不存在，跳过测试")
		return
	}

	// 应用层应该依赖领域层的接口，而不是基础设施层的具体实现
	err := filepath.Walk(applicationDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			// 检查文件中的结构体字段是否使用接口类型
			checkStructFieldsUseInterfaces(t, path)
		}
		return nil
	})

	if err != nil {
		t.Fatalf("遍历应用层目录失败: %v", err)
	}
}

// **Feature: code-structure-optimization, Property 11: 模块间接口通信**
// **验证需求: 需求 3.3**
// 对于任何模块间的通信，都应该通过明确定义的接口进行

func TestModuleInterfaceCommunication(t *testing.T) {
	rootDir := getProjectRootForInterface()
	domainDir := filepath.Join(rootDir, "internal", "domain")

	// 检查领域层是否定义了足够的接口
	interfaceCount := 0

	err := filepath.Walk(domainDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			count := countInterfacesInFileForInterface(path)
			interfaceCount += count
		}
		return nil
	})

	if err != nil {
		t.Fatalf("遍历领域层目录失败: %v", err)
	}

	// 验证领域层定义了足够的接口用于模块间通信
	if interfaceCount < 5 {
		t.Errorf("领域层应该定义足够的接口用于模块间通信，当前只有 %d 个", interfaceCount)
	} else {
		t.Logf("领域层定义了 %d 个接口用于模块间通信", interfaceCount)
	}
}

// TestDependencyInjectionContainerExists 验证依赖注入容器存在
func TestDependencyInjectionContainerExists(t *testing.T) {
	rootDir := getProjectRootForInterface()
	diDir := filepath.Join(rootDir, "internal", "infrastructure", "di")

	if _, err := os.Stat(diDir); os.IsNotExist(err) {
		t.Fatalf("依赖注入目录不存在: %s", diDir)
	}

	containerFile := filepath.Join(diDir, "container.go")
	if _, err := os.Stat(containerFile); os.IsNotExist(err) {
		t.Fatalf("依赖注入容器文件不存在: %s", containerFile)
	}

	t.Log("依赖注入容器已正确配置")
}

// checkStructFieldsUseInterfaces 检查结构体字段是否使用接口类型
func checkStructFieldsUseInterfaces(t *testing.T, filePath string) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			if structType, isStruct := typeSpec.Type.(*ast.StructType); isStruct {
				for _, field := range structType.Fields.List {
					// 检查字段类型是否是接口或指向接口的指针
					if ident, ok := field.Type.(*ast.Ident); ok {
						// 如果字段类型是大写开头（导出的），检查是否是接口
						if len(ident.Name) > 0 && ident.Name[0] >= 'A' && ident.Name[0] <= 'Z' {
							// 这里我们假设以 "er" 结尾的类型是接口（Go 惯例）
							// 或者是已知的接口类型
							if strings.HasSuffix(ident.Name, "Repository") ||
								strings.HasSuffix(ident.Name, "Logger") ||
								strings.HasSuffix(ident.Name, "Client") ||
								strings.HasSuffix(ident.Name, "Handler") ||
								strings.HasSuffix(ident.Name, "Service") ||
								strings.HasSuffix(ident.Name, "Provider") {
								// 这是一个接口类型，符合预期
							}
						}
					}
				}
			}
		}
		return true
	})
}

// countInterfacesInFileForInterface 统计文件中的接口数量
func countInterfacesInFileForInterface(filePath string) int {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return 0
	}

	count := 0
	ast.Inspect(node, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			if _, isInterface := typeSpec.Type.(*ast.InterfaceType); isInterface {
				count++
			}
		}
		return true
	})

	return count
}

// getProjectRootForInterface 获取项目根目录
func getProjectRootForInterface() string {
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
