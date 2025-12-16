// Package services 提供应用层服务实现
package services

import (
	"tr369-wss-client/internal/domain/repositories"
	"tr369-wss-client/internal/domain/services"
	"tr369-wss-client/pkg/api"
)

// MessageProcessor 消息处理器
type MessageProcessor struct {
	repository repositories.ClientRepository
	logger     services.Logger
}

// NewMessageProcessor 创建新的消息处理器
func NewMessageProcessor(repo repositories.ClientRepository, logger services.Logger) *MessageProcessor {
	return &MessageProcessor{
		repository: repo,
		logger:     logger,
	}
}

// ProcessGetRequest 处理 GET 请求
func (mp *MessageProcessor) ProcessGetRequest(paths []string) api.Response_GetResp {
	mp.logger.Debug("处理 GET 请求, 路径数量: %d", len(paths))
	return mp.repository.ConstructGetResp(paths)
}

// ProcessSetRequest 处理 SET 请求
func (mp *MessageProcessor) ProcessSetRequest(path, key, value string) {
	mp.logger.Debug("处理 SET 请求, 路径: %s, 键: %s", path, key)
	mp.repository.HandleSetRequest(path, key, value)
}

// ProcessAddRequest 处理 ADD 请求
func (mp *MessageProcessor) ProcessAddRequest(path string) string {
	mp.logger.Debug("处理 ADD 请求, 路径: %s", path)
	return mp.repository.GetNewInstance(path)
}

// ProcessDeleteRequest 处理 DELETE 请求
func (mp *MessageProcessor) ProcessDeleteRequest(path string) (string, bool) {
	mp.logger.Debug("处理 DELETE 请求, 路径: %s", path)
	return mp.repository.HandleDeleteRequest(path)
}

// SaveData 保存数据
func (mp *MessageProcessor) SaveData() {
	mp.repository.SaveData()
}

// IsPathExist 检查路径是否存在
func (mp *MessageProcessor) IsPathExist(path string) (bool, string) {
	return mp.repository.IsExistPath(path)
}
