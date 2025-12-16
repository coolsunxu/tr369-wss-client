// Package repositories 定义领域层的仓储接口
package repositories

import (
	"tr369-wss-client/pkg/api"
)

// ClientRepository 定义客户端数据仓储接口
// 该接口定义了客户端数据操作的所有方法
type ClientRepository interface {
	// ConstructGetResp 构建 GET 响应
	ConstructGetResp(paths []string) api.Response_GetResp

	// HandleSetRequest 处理 SET 请求
	HandleSetRequest(path, key, value string)

	// HandleDeleteRequest 处理 DELETE 请求
	HandleDeleteRequest(path string) (string, bool)

	// GetNewInstance 获取新实例路径
	GetNewInstance(path string) string

	// IsExistPath 检查路径是否存在
	IsExistPath(path string) (bool, string)

	// GetValueByPath 根据路径获取值
	GetValueByPath(path string) (interface{}, error)

	// SaveData 保存数据
	SaveData()

	// StartClientRepository 启动仓储服务
	StartClientRepository()
}
