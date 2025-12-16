// Package tr181 提供 TR181 相关的用例实现
package tr181

import (
	"tr369-wss-client/internal/domain/repositories"
	"tr369-wss-client/internal/domain/services"
)

// AddObjectUseCase ADD 对象用例
type AddObjectUseCase struct {
	repository repositories.ClientRepository
	logger     services.Logger
}

// NewAddObjectUseCase 创建新的 ADD 对象用例
func NewAddObjectUseCase(repo repositories.ClientRepository, logger services.Logger) *AddObjectUseCase {
	return &AddObjectUseCase{
		repository: repo,
		logger:     logger,
	}
}

// Execute 执行 ADD 对象操作
func (uc *AddObjectUseCase) Execute(path string) string {
	uc.logger.Debug("执行 ADD 对象操作, 路径: %s", path)
	return uc.repository.GetNewInstance(path)
}

// SetParameter 设置参数
func (uc *AddObjectUseCase) SetParameter(path, key, value string) {
	uc.repository.HandleSetRequest(path, key, value)
}

// SaveData 保存数据
func (uc *AddObjectUseCase) SaveData() {
	uc.repository.SaveData()
}

// DeleteObjectUseCase DELETE 对象用例
type DeleteObjectUseCase struct {
	repository repositories.ClientRepository
	logger     services.Logger
}

// NewDeleteObjectUseCase 创建新的 DELETE 对象用例
func NewDeleteObjectUseCase(repo repositories.ClientRepository, logger services.Logger) *DeleteObjectUseCase {
	return &DeleteObjectUseCase{
		repository: repo,
		logger:     logger,
	}
}

// Execute 执行 DELETE 对象操作
func (uc *DeleteObjectUseCase) Execute(path string) (string, bool) {
	uc.logger.Debug("执行 DELETE 对象操作, 路径: %s", path)
	return uc.repository.HandleDeleteRequest(path)
}

// SaveData 保存数据
func (uc *DeleteObjectUseCase) SaveData() {
	uc.repository.SaveData()
}
