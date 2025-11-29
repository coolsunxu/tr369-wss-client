package repository

import (
	"tr369-wss-client/client/model"
	"tr369-wss-client/common"
	"tr369-wss-client/config"
	"tr369-wss-client/pkg/api"
	"tr369-wss-client/trtree"
)

type clientRepository struct {
	Config         *config.Config
	TR181DataModel *model.TR181DataModel
}

func NewClientRepository(
	config *config.Config,
) model.ClientRepository {

	tr181DataModel := &model.TR181DataModel{
		Parameters: make(map[string]interface{}),
		Listeners:  make(map[string][]model.ParameterChangeListener),
	}

	// 初始化默认参数值
	loadDefaultTR181Nodes(tr181DataModel, config)

	return &clientRepository{
		Config:         config,
		TR181DataModel: tr181DataModel,
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

func loadDefaultTR181Nodes(tr181DataModel *model.TR181DataModel, config *config.Config) {

	tr181DataModel.Parameters = common.LoadJsonFile(config.TR181DataModelPath)
}
