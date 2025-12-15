package usecase

import (
	"fmt"
	"regexp"
	"tr369-wss-client/client/model"
	logger "tr369-wss-client/log"
	tr181Model "tr369-wss-client/tr181/model"
)

// 预编译正则表达式，避免每次调用都重新编译
// 匹配 Device.LocalAgent.Subscription.{正整数}. 格式
var subscriptionInstanceRegex = regexp.MustCompile(`^` + regexp.QuoteMeta(model.PathSubscription) + `([1-9]\d*)\.`)

// HandleAddLocalAgentSubscription 处理添加 LocalAgent 订阅
// 仅当路径为订阅节点时才处理，否则静默返回
func (uc *ClientUseCase) HandleAddLocalAgentSubscription(requestPath string, paramSettings map[string]string) error {
	// 非订阅节点路径，静默返回
	if !uc.isSubscriptionPath(requestPath) {
		return nil
	}

	// 提取并验证订阅参数
	params, err := uc.extractSubscriptionParams(paramSettings)
	if err != nil {
		return fmt.Errorf("add subscription failed: %w", err)
	}

	// 注册订阅监听器
	if err := uc.HandleSubscription(params.ReferenceList, params.Id, params.NotifType); err != nil {
		return fmt.Errorf("register subscription listener failed (id=%s, type=%s): %w",
			params.Id, params.NotifType, err)
	}

	logger.Infof("[USP] ADD_SUBSCRIPTION: success id=%s, type=%s, ref=%s",
		params.Id, params.NotifType, params.ReferenceList)
	return nil
}

// HandleDeleteLocalAgentSubscription 处理删除 LocalAgent 订阅
// 支持三种删除场景：父节点删除（批量）、订阅实例删除（单个）、非订阅路径（忽略）
func (uc *ClientUseCase) HandleDeleteLocalAgentSubscription(requestPath string) error {
	if requestPath == "" {
		return nil
	}

	// 场景1：父节点删除 - 清除所有订阅
	if uc.isSubscriptionParentPath(requestPath) {
		return uc.deleteAllSubscriptions(requestPath)
	}

	// 场景2：订阅实例删除 - 删除单个订阅
	if uc.isSubscriptionInstancePath(requestPath) {
		return uc.deleteSingleSubscription(requestPath)
	}

	// 场景3：非订阅相关路径，静默返回
	return nil
}

// HandleSubscription 处理订阅消息
func (uc *ClientUseCase) HandleSubscription(path string, subscriptionId string, subscriptionType string) error {
	switch subscriptionType {
	case tr181Model.ValueChange:
		// 处理 value change 订阅
		err := uc.ListenerMgr.AddListener(path, tr181Model.Listener{
			SubscriptionId: subscriptionId,
			Listener:       uc.HandleValueChange,
		})
		if err != nil {
			return err
		}

	case tr181Model.ObjectCreation:
		// 处理对象创建订阅
		err := uc.ListenerMgr.AddListener(path, tr181Model.Listener{
			SubscriptionId: subscriptionId,
			Listener:       uc.HandleObjectCreation,
		})
		if err != nil {
			return err
		}
	case tr181Model.ObjectDeletion:
		// 处理对象删除订阅
		err := uc.ListenerMgr.AddListener(path, tr181Model.Listener{
			SubscriptionId: subscriptionId,
			Listener:       uc.HandleObjectDeletion,
		})
		if err != nil {
			return err
		}
	case tr181Model.OperationComplete:
		// 处理操作完成订阅
	default:
		return fmt.Errorf("unknown subscription type %s", subscriptionType)
	}

	return nil
}

// isSubscriptionPath 判断路径是否为订阅节点
func (uc *ClientUseCase) isSubscriptionPath(path string) bool {
	return path == tr181Model.DeviceLocalAgentSubscription
}

// isSubscriptionParentPath 判断是否为订阅的父节点路径
func (uc *ClientUseCase) isSubscriptionParentPath(path string) bool {
	switch path {
	case model.PathDevice, model.PathLocalAgent, model.PathSubscription:
		return true
	}
	return false
}

// isSubscriptionInstancePath 判断是否为订阅实例路径
// 格式：Device.LocalAgent.Subscription.{正整数}.
func (uc *ClientUseCase) isSubscriptionInstancePath(path string) bool {
	return subscriptionInstanceRegex.MatchString(path)
}

// extractSubscriptionParams 从参数设置中提取订阅参数
func (uc *ClientUseCase) extractSubscriptionParams(paramSettings map[string]string) (*model.SubscriptionParams, error) {
	if paramSettings == nil {
		return nil, fmt.Errorf("paramSettings is nil")
	}

	params := &model.SubscriptionParams{
		Id:            paramSettings["ID"],
		ReferenceList: paramSettings["ReferenceList"],
		NotifType:     paramSettings["NotifType"],
	}

	// 验证必要参数
	if params.Id == "" {
		return nil, fmt.Errorf("missing required parameter 'ID'")
	}
	if params.ReferenceList == "" {
		return nil, fmt.Errorf("missing required parameter 'ReferenceList'")
	}
	if params.NotifType == "" {
		return nil, fmt.Errorf("missing required parameter 'NotifType'")
	}

	return params, nil
}

// deleteAllSubscriptions 删除所有订阅
func (uc *ClientUseCase) deleteAllSubscriptions(parentPath string) error {
	logger.Infof("[USP] DELETE_SUBSCRIPTION: deleting all subscriptions for parent path=%s", parentPath)
	return uc.ListenerMgr.ResetListener()
}

// deleteSingleSubscription 删除单个订阅实例
func (uc *ClientUseCase) deleteSingleSubscription(instancePath string) error {
	// 获取订阅的 ReferenceList
	refList, err := uc.getSubscriptionReferenceList(instancePath)
	if err != nil {
		return err
	}

	// 移除监听器
	if err := uc.ListenerMgr.RemoveListener(refList); err != nil {
		logger.Warnf("[USP] DELETE_SUBSCRIPTION remove listener error: path=%s, err=%v", instancePath, err)
		return err
	}

	logger.Infof("[USP] DELETE_SUBSCRIPTION: success path=%s, ReferenceList=%s", instancePath, refList)
	return nil
}

// getSubscriptionReferenceList 获取订阅实例的 ReferenceList
func (uc *ClientUseCase) getSubscriptionReferenceList(instancePath string) (string, error) {
	pathName, err := uc.DataRepo.GetValue(instancePath + "ReferenceList")
	if err != nil {
		logger.Warnf("[USP] DELETE_SUBSCRIPTION get ReferenceList error: path=%s, err=%v", instancePath, err)
		return "", err
	}

	refList, ok := pathName.(string)
	if !ok {
		err := fmt.Errorf("ReferenceList is not a string, got %T", pathName)
		logger.Warnf("[USP] DELETE_SUBSCRIPTION ReferenceList type error: path=%s, err=%v", instancePath, err)
		return "", err
	}

	return refList, nil
}
