package repository

import (
	"context"
	"log"
	"time"
	"tr369-wss-client/client/model"
	"tr369-wss-client/common"
	"tr369-wss-client/config"
	"tr369-wss-client/pkg/api"
	"tr369-wss-client/trtree"
)

type clientRepository struct {
	Config         *config.Config
	TR181DataModel *model.TR181DataModel
	writeCount     int
	lastWriteTime  int64
	pingTicker     *time.Ticker
	ctx            context.Context
	cancel         context.CancelFunc
}

func NewClientRepository(
	config *config.Config,
	ctx context.Context,
	cancel context.CancelFunc,
) model.ClientRepository {

	tr181DataModel := &model.TR181DataModel{
		Parameters: make(map[string]interface{}),
		Listeners:  make(map[string][]model.ParameterChangeListener),
	}

	pingTicker := time.NewTicker(time.Duration(config.DataRefreshConfig.IntervalSeconds))

	return &clientRepository{
		Config:         config,
		TR181DataModel: tr181DataModel,
		pingTicker:     pingTicker,
		ctx:            ctx,
		cancel:         cancel,
	}
}

func (repo *clientRepository) StartClientRepository() {

	// 初始化写入计数
	repo.writeCount = 0
	repo.lastWriteTime = time.Now().UnixMilli()

	// 初始化默认参数值
	loadDefaultTR181Nodes(repo.TR181DataModel, repo.Config)

	// 启动数据同步定时器
	go repo.DataSynchronizationTick()

}

func (repo *clientRepository) DataSynchronizationTick() {
	// 防止很久没有写入数据，但是writeCount没有达到写入磁盘的阈值
	for {
		select {
		case <-repo.pingTicker.C:
			// 没有数据写入
			if repo.writeCount == 0 {
				log.Println("No data change detected, skipping save.")
				continue
			}
			// 保存数据
			common.SaveJsonFile(repo.TR181DataModel.Parameters, repo.Config.TR181DataModelPath)
			// 重置写入计数
			repo.writeCount = 0
			// 重置最后写入时间
			repo.lastWriteTime = time.Now().UnixMilli()
		case <-repo.ctx.Done():
			return
		}
	}

}

func (repo *clientRepository) ConstructGetResp(paths []string) api.Response_GetResp {
	return trtree.ConstructGetResp(repo.TR181DataModel.Parameters, paths)
}

func (repo *clientRepository) HandleSetRequest(path string, key string, value string) {
	trtree.HandleSetRequest(repo.TR181DataModel.Parameters, path, key, value)
}

func (repo *clientRepository) IsExistPath(path string) (isSuccess bool, nodePath string) {
	return trtree.IsExistPath(repo.TR181DataModel.Parameters, path)
}

func (repo *clientRepository) GetNewInstance(path string) (nodePath string) {
	return trtree.GetNewInstance(repo.TR181DataModel.Parameters, path)
}

func (repo *clientRepository) HandleDeleteRequest(path string) (nodePath string, isFound bool) {
	return trtree.HandleDeleteRequest(repo.TR181DataModel.Parameters, path)
}

func (repo *clientRepository) SaveData() {

	// 每次写入计数加1
	repo.writeCount++

	// 如果写入次数大于约定的最大次数，开始保存到文件磁盘
	if repo.writeCount >= repo.Config.DataRefreshConfig.WriteCountThreshold {
		common.SaveJsonFile(repo.TR181DataModel.Parameters, repo.Config.TR181DataModelPath)
		// 重置写入计数
		repo.writeCount = 0
		// 重置最后写入时间
		repo.lastWriteTime = time.Now().UnixMilli()
	}

	// todo print log

}

func loadDefaultTR181Nodes(tr181DataModel *model.TR181DataModel, config *config.Config) {

	tr181DataModel.Parameters = common.LoadJsonFile(config.TR181DataModelPath)
}
