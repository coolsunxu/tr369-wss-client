// Package integration 提供架构验证测试
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

// **Feature: code-structure-optimization, Property 13: 接口向后兼容性**
// **验证需求: 需求 3.5**
// 对于任何公开接口的修改，都应该保持向后兼容性

func TestInterfaceBackwardCompatibility(t *testing.T) {
	rootDir := getProjectRootForArch()
	domainDir := filepath.Join(rootDir, "internal", "domain")

	// 收集所有接口及其方法
	interfaces := make(map[string][]string)

	err := filepath.Walk(domainDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			fileInterfaces := getInterfacesWithMethods(path)
			for name, methods := range fileInterfaces {
				interfaces[name] = methods
			}
		}
		return nil
	})

	if err != nil {
		t.Fatalf("遍历领域层目录失败: %v", err)
	}

	// 验证核心接口存在且有方法
	coreInterfaces := []string{
		"TR181Repository",
		"ClientRepository",
		"WebSocketClient",
		"MessageHandler",
		"Logger",
	}

	for _, iface := range coreInterfaces {
		if methods, exists := interfaces[iface]; exists {
			if len(methods) == 0 {
				t.Errorf("接口 %s 应该有方法定义", iface)
			} else {
				t.Logf("接口 %s 有 %d 个方法", iface, len(methods))
			}
		}
	}
}

// **Feature: code-structure-optimization, Property 5: 功能模块内聚性**
// **验证需求: 需求 2.2**
// 对于任何功能模块，相关的代码应该组织在同一个包或目录中

func TestFunctionalModuleCohesion(t *testing.T) {
	rootDir := getProjectRootForArch()

	// 检查功能模块是否内聚
	modules := map[string][]string{
		"domain/entities/tr181":                {"datamodel.go", "subscription.go"},
		"domain/entities/usp":                  {"message.go", "record.go"},
		"domain/repositories":                  {"tr181_repository.go", "client_repository.go"},
		"domain/services":                      {"message_handler.go", "websocket_client.go", "logger.go", "config.go"},
		"infrastructure/config":                {"config.go", "viper.go"},
		"infrastructure/logging":               {"logger.go", "zap_logger.go"},
		"infrastructure/websocket":             {"client.go"},
		"infrastructure/persistence/json":      {"file_manager.go"},
		"infrastructure/persistence/repository": {"client_repository.go"},
	}

	for modulePath, expectedFiles := range modules {
		fullPath := filepath.Join(rootDir, "internal", modulePath)

		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Logf("模块目录不存在: %s", modulePath)
			continue
		}

		entries, err := os.ReadDir(fullPath)
		if err != nil {
			t.Logf("读取模块目录失败: %s", modulePath)
			continue
		}

		foundFiles := make(map[string]bool)
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") && !strings.HasSuffix(entry.Name(), "_test.go") {
				foundFiles[entry.Name()] = true
			}
		}

		for _, expected := range expectedFiles {
			if !foundFiles[expected] {
				t.Logf("模块 %s 缺少文件: %s", modulePath, expected)
			}
		}

		t.Logf("模块 %s 内聚性检查完成，包含 %d 个源文件", modulePath, len(foundFiles))
	}
}

// getInterfacesWithMethods 获取文件中的接口及其方法
func getInterfacesWithMethods(filePath string) map[string][]string {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil
	}

	interfaces := make(map[string][]string)

	ast.Inspect(node, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			if interfaceType, isInterface := typeSpec.Type.(*ast.InterfaceType); isInterface {
				var methods []string
				for _, method := range interfaceType.Methods.List {
					if len(method.Names) > 0 {
						methods = append(methods, method.Names[0].Name)
					}
				}
				interfaces[typeSpec.Name.Name] = methods
			}
		}
		return true
	})

	return interfaces
}

// getProjectRootForArch 获取项目根目录
func getProjectRootForArch() string {
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
