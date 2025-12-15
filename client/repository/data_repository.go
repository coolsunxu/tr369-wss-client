package repository

import (
	"fmt"
	"strings"

	logger "tr369-wss-client/log"
	"tr369-wss-client/trtree"
)

// DataRepository 实现 model.DataRepository 接口
type DataRepository struct {
	*BaseRepository
}

// NewDataRepository 创建 DataRepository 实例
func NewDataRepository(base *BaseRepository) *DataRepository {
	return &DataRepository{BaseRepository: base}
}

// GetValue 获取指定路径的值
func (repo *DataRepository) GetValue(path string) (interface{}, error) {
	paths := strings.Split(path, ".")
	value, _, found := trtree.FindKeyInMap(repo.TR181DataModel.Parameters, paths, "")
	if !found {
		return path, fmt.Errorf("path not found: %s", path)
	}
	return value, nil
}

// GetParameters 获取底层参数数据（供 UseCase 构建响应使用）
func (repo *DataRepository) GetParameters() map[string]interface{} {
	return repo.TR181DataModel.Parameters
}

// SetValue 设置指定路径的值
// 返回: changed (是否发生变化), oldValue (旧值)
func (repo *DataRepository) SetValue(path string, key string, value string) (changed bool, oldValue string) {
	// 先获取旧值
	oldVal, err := repo.GetValue(path + key)
	if err != nil {
		logger.Debugf("Failed to get old value for path %s%s: %v", path, key, err)
		oldValue = ""
	} else {
		oldValue, _ = oldVal.(string)
	}

	// 保存到数据库
	trtree.HandleSetRequest(repo.TR181DataModel.Parameters, path, key, value)
	repo.SaveData()

	// 判断是否有变化
	changed = oldValue != value
	return changed, oldValue
}

// DeleteNode 删除指定路径的节点
func (repo *DataRepository) DeleteNode(path string) (nodePath string, isFound bool) {
	return trtree.HandleDeleteRequest(repo.TR181DataModel.Parameters, path)
}

// Start 启动数据仓库（初始化和数据同步）
func (repo *DataRepository) Start() {
	repo.BaseRepository.Start()
}
