// Package repository 提供仓储实现
package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"tr369-wss-client/internal/domain/entities/tr181"
	"tr369-wss-client/internal/domain/services"
	"tr369-wss-client/internal/infrastructure/config"
	"tr369-wss-client/internal/infrastructure/persistence/json"
	"tr369-wss-client/internal/infrastructure/persistence/trtree"
	"tr369-wss-client/pkg/api"
)

// ClientRepository 客户端仓储实现
type ClientRepository struct {
	config        *config.Config
	dataModel     *tr181.DataModel
	writeCount    int
	lastWriteTime int64
	ticker        *time.Ticker
	ctx           context.Context
	cancel        context.CancelFunc
	logger        services.Logger
	fileManager   *json.FileManager
}

// NewClientRepository 创建新的客户端仓储
func NewClientRepository(
	cfg *config.Config,
	ctx context.Context,
	cancel context.CancelFunc,
	logger services.Logger,
) *ClientRepository {
	dataModel := tr181.NewDataModel()
	ticker := time.NewTicker(time.Duration(cfg.DataRefreshConfig.IntervalSeconds) * time.Second)

	return &ClientRepository{
		config:      cfg,
		dataModel:   dataModel,
		ticker:      ticker,
		ctx:         ctx,
		cancel:      cancel,
		logger:      logger,
		fileManager: json.NewFileManager(logger),
	}
}

// StartClientRepository 启动仓储服务
func (r *ClientRepository) StartClientRepository() {
	r.writeCount = 0
	r.lastWriteTime = time.Now().UnixMilli()

	// 加载默认数据
	r.loadDefaultData()

	// 启动数据同步
	go r.dataSynchronizationTick()
}

// loadDefaultData 加载默认数据
func (r *ClientRepository) loadDefaultData() {
	data, err := r.fileManager.LoadJSONFile(r.config.DataRefreshConfig.TR181DataModelPath, 3, time.Second)
	if err != nil {
		r.logger.Warn("加载默认数据失败: %v", err)
		return
	}
	r.dataModel.Parameters = data
}

// dataSynchronizationTick 数据同步定时器
func (r *ClientRepository) dataSynchronizationTick() {
	for {
		select {
		case <-r.ticker.C:
			if r.writeCount == 0 {
				r.logger.Debug("没有数据变更，跳过保存")
				continue
			}
			r.fileManager.SaveJSONFile(r.dataModel.Parameters, r.config.DataRefreshConfig.TR181DataModelPath)
			r.writeCount = 0
			r.lastWriteTime = time.Now().UnixMilli()
			r.logger.Debug("定时保存数据完成")
		case <-r.ctx.Done():
			return
		}
	}
}

// GetValueByPath 根据路径获取值
func (r *ClientRepository) GetValueByPath(path string) (interface{}, error) {
	paths := strings.Split(path, ".")
	value, _, found := trtree.FindKeyInMap(r.dataModel.Parameters, paths, "")
	if !found {
		return path, fmt.Errorf("路径未找到: %s", path)
	}
	return value, nil
}

// ConstructGetResp 构建 GET 响应
func (r *ClientRepository) ConstructGetResp(paths []string) api.Response_GetResp {
	return trtree.ConstructGetResp(r.dataModel.Parameters, paths)
}

// HandleSetRequest 处理 SET 请求
func (r *ClientRepository) HandleSetRequest(path, key, value string) {
	old, _ := r.GetValueByPath(path)
	trtree.HandleSetRequest(r.dataModel.Parameters, path, key, value)

	if oldValue, ok := old.(string); ok && oldValue != value {
		notifyValueChange := &api.Notify_ValueChange_{
			ValueChange: &api.Notify_ValueChange{
				ParamPath:  path,
				ParamValue: value,
			},
		}
		r.NotifyListeners(path+key, notifyValueChange)
	}
}

// IsExistPath 检查路径是否存在
func (r *ClientRepository) IsExistPath(path string) (bool, string) {
	return trtree.IsExistPath(r.dataModel.Parameters, path)
}

// GetNewInstance 获取新实例路径
func (r *ClientRepository) GetNewInstance(path string) string {
	return trtree.GetNewInstance(r.dataModel.Parameters, path)
}

// HandleDeleteRequest 处理 DELETE 请求
func (r *ClientRepository) HandleDeleteRequest(path string) (string, bool) {
	return trtree.HandleDeleteRequest(r.dataModel.Parameters, path)
}

// SaveData 保存数据
func (r *ClientRepository) SaveData() {
	r.writeCount++

	if r.writeCount >= r.config.DataRefreshConfig.WriteCountThreshold {
		r.fileManager.SaveJSONFile(r.dataModel.Parameters, r.config.DataRefreshConfig.TR181DataModelPath)
		r.writeCount = 0
		r.lastWriteTime = time.Now().UnixMilli()
		r.logger.Debug("正常保存数据完成")
	}

	r.logger.Debug("当前写入计数: %d", r.writeCount)
}

// AddListener 添加监听器
func (r *ClientRepository) AddListener(paramName string, listener tr181.Listener) error {
	r.dataModel.Listeners[paramName] = append(r.dataModel.Listeners[paramName], listener)
	return nil
}

// RemoveListener 移除监听器
func (r *ClientRepository) RemoveListener(paramName string) error {
	delete(r.dataModel.Listeners, paramName)
	return nil
}

// ResetListeners 重置所有监听器
func (r *ClientRepository) ResetListeners() error {
	r.dataModel.Listeners = make(map[string][]tr181.Listener)
	return nil
}

// NotifyListeners 通知监听器
func (r *ClientRepository) NotifyListeners(paramName string, value interface{}) {
	listeners, exists := r.dataModel.Listeners[paramName]
	if exists {
		for _, listener := range listeners {
			go listener.Listener(listener.SubscriptionId, value)
		}
	}
}
