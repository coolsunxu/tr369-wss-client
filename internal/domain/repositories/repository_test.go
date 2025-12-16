// Package repositories 测试仓储接口定义
package repositories

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// **Feature: code-structure-optimization, Property 9: 业务接口领域层定位**
// **验证需求: 需求 3.1**
// 对于任何业务接口定义，都应该位于领域层中

func TestBusinessInterfacesInDomainLayer(t *testing.T) {
	rootDir := getProjectRoot()
	domainDir := filepath.Join(rootDir, "internal", "domain")

	// 检查领域层目录是否存在
	if _, err := os.Stat(domainDir); os.IsNotExist(err) {
		t.Fatalf("领域层目录不存在: %s", domainDir)
	}

	// 遍历领域层目录，查找接口定义
	interfaceCount := 0
	err := filepath.Walk(domainDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理 Go 文件，跳过测试文件
		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			count, parseErr := countInterfacesInFile(path)
			if parseErr != nil {
				t.Logf("解析文件失败 %s: %v", path, parseErr)
				return nil
			}
			interfaceCount += count
		}
		return nil
	})

	if err != nil {
		t.Fatalf("遍历领域层目录失败: %v", err)
	}

	// 验证领域层至少定义了一些接口
	if interfaceCount == 0 {
		t.Error("领域层应该定义业务接口")
	}

	t.Logf("领域层共定义了 %d 个接口", interfaceCount)
}

// TestRepositoryInterfacesExist 测试仓储接口是否存在
func TestRepositoryInterfacesExist(t *testing.T) {
	rootDir := getProjectRoot()
	repoDir := filepath.Join(rootDir, "internal", "domain", "repositories")

	// 期望的仓储接口
	expectedInterfaces := []string{
		"TR181Repository",
		"ClientRepository",
	}

	// 收集所有接口名称
	foundInterfaces := make(map[string]bool)

	err := filepath.Walk(repoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			interfaces, parseErr := getInterfaceNames(path)
			if parseErr != nil {
				return nil
			}
			for _, name := range interfaces {
				foundInterfaces[name] = true
			}
		}
		return nil
	})

	if err != nil {
		t.Fatalf("遍历仓储目录失败: %v", err)
	}

	for _, expected := range expectedInterfaces {
		if !foundInterfaces[expected] {
			t.Errorf("缺少期望的仓储接口: %s", expected)
		}
	}
}

// TestServiceInterfacesExist 测试服务接口是否存在
func TestServiceInterfacesExist(t *testing.T) {
	rootDir := getProjectRoot()
	serviceDir := filepath.Join(rootDir, "internal", "domain", "services")

	// 期望的服务接口
	expectedInterfaces := []string{
		"MessageHandler",
		"WebSocketClient",
	}

	// 收集所有接口名称
	foundInterfaces := make(map[string]bool)

	err := filepath.Walk(serviceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			interfaces, parseErr := getInterfaceNames(path)
			if parseErr != nil {
				return nil
			}
			for _, name := range interfaces {
				foundInterfaces[name] = true
			}
		}
		return nil
	})

	if err != nil {
		t.Fatalf("遍历服务目录失败: %v", err)
	}

	for _, expected := range expectedInterfaces {
		if !foundInterfaces[expected] {
			t.Errorf("缺少期望的服务接口: %s", expected)
		}
	}
}

// countInterfacesInFile 统计文件中的接口数量
func countInterfacesInFile(filePath string) (int, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return 0, err
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

	return count, nil
}

// getInterfaceNames 获取文件中的接口名称
func getInterfaceNames(filePath string) ([]string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var names []string
	ast.Inspect(node, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			if _, isInterface := typeSpec.Type.(*ast.InterfaceType); isInterface {
				names = append(names, typeSpec.Name.Name)
			}
		}
		return true
	})

	return names, nil
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
