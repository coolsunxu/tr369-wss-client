package repository

import (
	"context"
	"time"

	"tr369-wss-client/client/model"
	"tr369-wss-client/common"
	"tr369-wss-client/config"
	logger "tr369-wss-client/log"
	tr181Model "tr369-wss-client/tr181/model"
)

// BaseRepository 共享的基础仓库结构
type BaseRepository struct {
	Config         *config.Config
	TR181DataModel *tr181Model.TR181DataModel
	WriteCount     int
	LastWriteTime  int64
	PingTicker     *time.Ticker
	Ctx            context.Context
	Cancel         context.CancelFunc
}

// NewRepository 创建新的仓库实例
// 返回分别实现 model.DataRepository 和 model.ListenerManager 接口的实例
func NewRepository(
	cfg *config.Config,
	ctx context.Context,
	cancel context.CancelFunc,
) (model.DataRepository, model.ListenerManager) {

	tr181DataModel := &tr181Model.TR181DataModel{
		Parameters: make(map[string]interface{}),
		Listeners:  make(map[string][]tr181Model.Listener),
	}

	pingTicker := time.NewTicker(time.Duration(cfg.DataRefreshConfig.IntervalSeconds) * time.Second)

	base := &BaseRepository{
		Config:         cfg,
		TR181DataModel: tr181DataModel,
		PingTicker:     pingTicker,
		Ctx:            ctx,
		Cancel:         cancel,
	}

	dataRepo := NewDataRepository(base)
	listenerMgr := NewListenerManager(base)

	return dataRepo, listenerMgr
}

// Start 启动数据仓库（初始化和数据同步）
func (repo *BaseRepository) Start() {
	// 初始化写入计数
	repo.WriteCount = 0
	repo.LastWriteTime = time.Now().UnixMilli()

	// 初始化默认参数值
	loadDefaultTR181Nodes(repo.TR181DataModel, repo.Config)

	// 启动数据同步定时器
	go repo.DataSynchronizationTick()
}

// DataSynchronizationTick 数据同步定时器
func (repo *BaseRepository) DataSynchronizationTick() {
	for {
		select {
		case <-repo.PingTicker.C:
			if repo.WriteCount == 0 {
				logger.Debugf("No data change detected, skipping save.")
				continue
			}
			common.SaveJsonFile(repo.TR181DataModel.Parameters, repo.Config.DataRefreshConfig.TR181DataModelPath)
			repo.WriteCount = 0
			repo.LastWriteTime = time.Now().UnixMilli()
			logger.Debugf("tick save data synchronized.")
		case <-repo.Ctx.Done():
			return
		}
	}
}

// SaveData 保存数据到磁盘
func (repo *BaseRepository) SaveData() {
	repo.WriteCount++

	if repo.WriteCount >= repo.Config.DataRefreshConfig.WriteCountThreshold {
		common.SaveJsonFile(repo.TR181DataModel.Parameters, repo.Config.DataRefreshConfig.TR181DataModelPath)
		repo.WriteCount = 0
		repo.LastWriteTime = time.Now().UnixMilli()
		logger.Debugf("normal save data synchronized.")
	}

	logger.Debugf("current write count: %d", repo.WriteCount)
}

// loadDefaultTR181Nodes 加载默认 TR181 节点
func loadDefaultTR181Nodes(tr181DataModel *tr181Model.TR181DataModel, config *config.Config) {
	tr181DataModel.Parameters = common.LoadJsonFile(config.DataRefreshConfig.TR181DataModelPath)
}
