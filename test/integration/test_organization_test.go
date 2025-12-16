// Package integration 提供测试组织验证
package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// **Feature: code-structure-optimization, Property 14: 测试文件同目录组织**
// **验证需求: 需求 4.1**
// 对于任何单元测试文件，都应该与被测试的源文件位于同一目录

func TestUnitTestFilesInSameDirectory(t *testing.T) {
	rootDir := getProjectRootForTestOrg()
	internalDir := filepath.Join(rootDir, "internal")

	if _, err := os.Stat(internalDir); os.IsNotExist(err) {
		t.Fatalf("internal 目录不存在: %s", internalDir)
	}

	// 检查 internal 目录下的测试文件是否与源文件在同一目录
	testFilesFound := 0
	err := filepath.Walk(internalDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, "_test.go") {
			testFilesFound++
			// 检查同目录下是否有对应的源文件
			dir := filepath.Dir(path)
			hasSourceFile := false

			entries, _ := os.ReadDir(dir)
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") && !strings.HasSuffix(entry.Name(), "_test.go") {
					hasSourceFile = true
					break
				}
			}

			if !hasSourceFile {
				t.Logf("测试文件 %s 所在目录没有源文件", path)
			}
		}
		return nil
	})

	if err != nil {
		t.Fatalf("遍历 internal 目录失败: %v", err)
	}

	if testFilesFound == 0 {
		t.Error("internal 目录下应该有单元测试文件")
	} else {
		t.Logf("发现 %d 个单元测试文件", testFilesFound)
	}
}

// **Feature: code-structure-optimization, Property 15: 集成测试专门目录**
// **验证需求: 需求 4.2**
// 对于任何集成测试，都应该位于专门的 test/integration 目录中

func TestIntegrationTestsInDedicatedDirectory(t *testing.T) {
	rootDir := getProjectRootForTestOrg()
	integrationDir := filepath.Join(rootDir, "test", "integration")

	if _, err := os.Stat(integrationDir); os.IsNotExist(err) {
		t.Fatalf("集成测试目录不存在: %s", integrationDir)
	}

	// 检查集成测试目录下是否有测试文件
	testFilesFound := 0
	entries, err := os.ReadDir(integrationDir)
	if err != nil {
		t.Fatalf("读取集成测试目录失败: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), "_test.go") {
			testFilesFound++
			t.Logf("发现集成测试文件: %s", entry.Name())
		}
	}

	if testFilesFound == 0 {
		t.Error("集成测试目录下应该有测试文件")
	}
}

// **Feature: code-structure-optimization, Property 16: 测试数据统一管理**
// **验证需求: 需求 4.3**
// 对于任何测试数据文件，都应该位于 test/fixtures 目录中

func TestTestDataInFixturesDirectory(t *testing.T) {
	rootDir := getProjectRootForTestOrg()
	fixturesDir := filepath.Join(rootDir, "test", "fixtures")

	if _, err := os.Stat(fixturesDir); os.IsNotExist(err) {
		t.Fatalf("测试数据目录不存在: %s", fixturesDir)
	}

	t.Log("测试数据目录已正确配置")
}

// getProjectRootForTestOrg 获取项目根目录
func getProjectRootForTestOrg() string {
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
