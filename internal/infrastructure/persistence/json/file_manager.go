// Package json 提供 JSON 文件持久化功能
package json

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"tr369-wss-client/internal/domain/services"
)

// FileManager 提供 JSON 文件管理功能
type FileManager struct {
	logger services.Logger
}

// NewFileManager 创建新的文件管理器
func NewFileManager(logger services.Logger) *FileManager {
	return &FileManager{
		logger: logger,
	}
}

// LoadJSONFile 从文件加载 JSON 数据
func (fm *FileManager) LoadJSONFile(filePath string, maxRetries int, retryInterval time.Duration) (map[string]interface{}, error) {
	var result map[string]interface{}

	for i := 0; i < maxRetries; i++ {
		file, err := os.Open(filePath)
		if err != nil {
			fm.logger.Warn("打开 JSON 文件失败, 路径: %s, 重试次数: %d, 错误: %v", filePath, i, err)
			time.Sleep(retryInterval)
			continue
		}

		decoder := json.NewDecoder(file)
		err = decoder.Decode(&result)
		file.Close()

		if err != nil {
			fm.logger.Warn("解析 JSON 数据失败, 路径: %s, 重试次数: %d, 错误: %v", filePath, i, err)
			time.Sleep(retryInterval)
			continue
		}

		return result, nil
	}

	return nil, fmt.Errorf("加载 JSON 文件失败，已重试 %d 次", maxRetries)
}

// SaveJSONFile 保存 JSON 数据到文件
func (fm *FileManager) SaveJSONFile(data map[string]interface{}, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		fm.logger.Warn("打开文件失败, 路径: %s", filePath)
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		fm.logger.Warn("编码 JSON 失败, 路径: %s", filePath)
		return fmt.Errorf("编码 JSON 失败: %w", err)
	}

	fm.logger.Debug("保存 JSON 文件成功, 路径: %s", filePath)
	return nil
}
