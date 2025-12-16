// Package tr181 提供 TR181 相关的用例实现
package tr181

import (
	"tr369-wss-client/internal/domain/repositories"
	"tr369-wss-client/internal/domain/services"
	"tr369-wss-client/pkg/api"
)

// GetParameterUseCase GET 参数用例
type GetParameterUseCase struct {
	repository repositories.ClientRepository
	logger     services.Logger
}

// NewGetParameterUseCase 创建新的 GET 参数用例
func NewGetParameterUseCase(repo repositories.ClientRepository, logger services.Logger) *GetParameterUseCase {
	return &GetParameterUseCase{
		repository: repo,
		logger:     logger,
	}
}

// Execute 执行 GET 参数操作
func (uc *GetParameterUseCase) Execute(paths []string) api.Response_GetResp {
	uc.logger.Debug("执行 GET 参数操作, 路径: %v", paths)
	return uc.repository.ConstructGetResp(paths)
}

// SetParameterUseCase SET 参数用例
type SetParameterUseCase struct {
	repository repositories.ClientRepository
	logger     services.Logger
}

// NewSetParameterUseCase 创建新的 SET 参数用例
func NewSetParameterUseCase(repo repositories.ClientRepository, logger services.Logger) *SetParameterUseCase {
	return &SetParameterUseCase{
		repository: repo,
		logger:     logger,
	}
}

// Execute 执行 SET 参数操作
func (uc *SetParameterUseCase) Execute(path, key, value string) {
	uc.logger.Debug("执行 SET 参数操作, 路径: %s, 键: %s, 值: %s", path, key, value)
	uc.repository.HandleSetRequest(path, key, value)
}

// SaveData 保存数据
func (uc *SetParameterUseCase) SaveData() {
	uc.repository.SaveData()
}
