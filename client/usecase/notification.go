package usecase

import (
	logger "tr369-wss-client/log"
	"tr369-wss-client/pkg/api"
	"tr369-wss-client/utils"
)

// sendNotification 通用通知发送函数
// 接收 subscriptionId、notification 和 notifyType 参数，统一构建 Notify 消息并发送
// notification 参数必须是实现了 isNotify_Notification 接口的类型
func (uc *ClientUseCase) sendNotification(subscriptionId string, notify *api.Notify, notifyType string) {
	msg := utils.CreateNotifyMessage(notify)

	// 发送消息，通过 channel
	if err := uc.HandleMTPMsgTransmit(msg); err != nil {
		logger.Warnf("[USP] %s notify error: subscriptionId=%s, err=%v", notifyType, subscriptionId, err)
		return
	}

	logger.Infof("[USP] send %s notify: %s", notifyType, msg)
}

// HandleValueChange 处理 value change 事件
func (uc *ClientUseCase) HandleValueChange(subscriptionId string, change interface{}) {
	valueChange, ok := change.(*api.Notify_ValueChange_)
	if !ok {
		logger.Warnf("[USP] VALUE_CHANGE type assertion error: subscriptionId=%s, err=expected Notify_ValueChange_, got %T", subscriptionId, change)
		return
	}

	notify := &api.Notify{
		SubscriptionId: subscriptionId,
		SendResp:       true,
		Notification:   valueChange,
	}
	uc.sendNotification(subscriptionId, notify, "VALUE_CHANGE")
}

// HandleObjectCreation 处理 object creation 事件
func (uc *ClientUseCase) HandleObjectCreation(subscriptionId string, change interface{}) {
	objCreation, ok := change.(*api.Notify_ObjCreation)
	if !ok {
		logger.Warnf("[USP] OBJ_CREATION type assertion error: subscriptionId=%s, err=expected Notify_ObjCreation, got %T", subscriptionId, change)
		return
	}

	notify := &api.Notify{
		SubscriptionId: subscriptionId,
		SendResp:       true,
		Notification:   objCreation,
	}
	uc.sendNotification(subscriptionId, notify, "OBJ_CREATION")
}

// HandleObjectDeletion 处理 object deletion 事件
func (uc *ClientUseCase) HandleObjectDeletion(subscriptionId string, change interface{}) {
	objDeletion, ok := change.(*api.Notify_ObjDeletion)
	if !ok {
		logger.Warnf("[USP] OBJ_DELETION type assertion error: subscriptionId=%s, err=expected Notify_ObjDeletion, got %T", subscriptionId, change)
		return
	}

	notify := &api.Notify{
		SubscriptionId: subscriptionId,
		SendResp:       true,
		Notification:   objDeletion,
	}
	uc.sendNotification(subscriptionId, notify, "OBJ_DELETION")
}

// notifyValueChange 发送值变化通知
func (uc *ClientUseCase) notifyValueChange(paramPath string, newValue string) {
	notifyValueChange := &api.Notify_ValueChange_{
		ValueChange: &api.Notify_ValueChange{
			ParamPath:  paramPath,
			ParamValue: newValue,
		},
	}
	uc.ListenerMgr.NotifyListeners(paramPath, notifyValueChange)
}

// notifyObjectCreation 发送对象创建通知
func (uc *ClientUseCase) notifyObjectCreation(path string, uniqueKeys map[string]string) {
	notifyCreateObj := &api.Notify_ObjCreation{
		ObjCreation: &api.Notify_ObjectCreation{
			ObjPath:    path,
			UniqueKeys: uniqueKeys,
		},
	}
	uc.ListenerMgr.NotifyListeners(path, notifyCreateObj)
}

// notifyObjectDeletion 发送对象删除通知
func (uc *ClientUseCase) notifyObjectDeletion(objPath string) {
	notifyDeleteObj := &api.Notify_ObjDeletion{
		ObjDeletion: &api.Notify_ObjectDeletion{
			ObjPath: objPath,
		},
	}
	uc.ListenerMgr.NotifyListeners(objPath, notifyDeleteObj)
}
