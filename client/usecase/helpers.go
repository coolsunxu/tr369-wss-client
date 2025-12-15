package usecase

import (
	"tr369-wss-client/client/model"
	"tr369-wss-client/pkg/api"
	"tr369-wss-client/trtree"
)

// extractParamSettings 从参数设置列表中提取键值对
// 支持 []*api.Set_UpdateParamSetting 和 []*api.Add_CreateParamSetting
func extractParamSettings[T model.ParamSetting](settings []T) map[string]string {
	result := make(map[string]string)
	for _, setting := range settings {
		result[setting.GetParam()] = setting.GetValue()
	}
	return result
}

// constructGetResp 构建 GET 响应（业务逻辑）
func (uc *ClientUseCase) constructGetResp(paths []string) api.Response_GetResp {
	params := uc.DataRepo.GetParameters()
	return trtree.ConstructGetResp(params, paths)
}

// isExistPath 检查路径是否存在（包含路径表达式解析）
func (uc *ClientUseCase) isExistPath(path string) (isSuccess bool, nodePath string) {
	params := uc.DataRepo.GetParameters()
	return trtree.IsExistPath(params, path)
}

// getNewInstance 获取新实例路径（包含实例编号生成策略）
func (uc *ClientUseCase) getNewInstance(path string) string {
	params := uc.DataRepo.GetParameters()
	return trtree.GetNewInstance(params, path)
}
