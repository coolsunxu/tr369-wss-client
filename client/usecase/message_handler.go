package usecase

import (
	logger "tr369-wss-client/log"
	"tr369-wss-client/pkg/api"
	"tr369-wss-client/utils"
)

// HandleGetRequest handles incoming GET requests
func (uc *ClientUseCase) HandleGetRequest(inComingMsg *api.Msg) {
	// 防御性检查
	if inComingMsg == nil || inComingMsg.Header == nil {
		logger.Warnf("[USP] HandleGetRequest received invalid message")
		return
	}

	msgId := inComingMsg.Header.MsgId
	logger.Infof("[USP] receive GET request: %s", inComingMsg.String())

	getNodePaths := inComingMsg.GetBody().GetRequest().GetGet().GetParamPaths()
	resp := uc.constructGetResp(getNodePaths)
	msg := utils.CreateGetResponseMessage(msgId, resp)
	logger.Infof("[USP] send GET response: %s", msg.String())

	err := uc.HandleMTPMsgTransmit(msg)
	if err != nil {
		logger.Warnf("[USP] GET error: msgId=%s, err=%v", msgId, err)
	}
}

// HandleSetRequest handles incoming SET requests
func (uc *ClientUseCase) HandleSetRequest(inComingMsg *api.Msg) {
	// 防御性检查
	if inComingMsg == nil || inComingMsg.Header == nil {
		logger.Warnf("[USP] HandleSetRequest received invalid message")
		return
	}

	msgId := inComingMsg.Header.MsgId
	logger.Infof("[USP] receive SET request: %s", inComingMsg.String())

	getUpdateObjs := inComingMsg.GetBody().GetRequest().GetSet().GetUpdateObjs()

	var affectedPath []string
	var requestPath []string
	var updatedParams []map[string]string

	for _, updateObj := range getUpdateObjs {
		path := updateObj.GetObjPath()
		isSuccess, nodePath := uc.isExistPath(path)
		if !isSuccess {
			continue
		}
		requestPath = append(requestPath, path)
		affectedPath = append(affectedPath, nodePath)

		// 使用辅助函数提取参数设置
		paramSettings := extractParamSettings(updateObj.GetParamSettings())
		// 将参数设置应用到 repository，并处理 value change 通知
		for setKey, setValue := range paramSettings {
			changed, _ := uc.DataRepo.SetValue(nodePath, setKey, setValue)
			if changed {
				// 发送 value change 通知
				uc.notifyValueChange(nodePath+setKey, setValue)
			}
		}
		updatedParams = append(updatedParams, paramSettings)
	}

	msg := utils.CreateSetResponseMessage(msgId, requestPath, affectedPath, updatedParams)
	logger.Infof("[USP] send SET response: %s", msg.String())

	err := uc.HandleMTPMsgTransmit(msg)
	if err != nil {
		logger.Warnf("[USP] SET error: msgId=%s, err=%v", msgId, err)
	}
}

// HandleAddRequest handles incoming ADD requests
func (uc *ClientUseCase) HandleAddRequest(inComingMsg *api.Msg) {
	// 防御性检查
	if inComingMsg == nil || inComingMsg.Header == nil {
		logger.Warnf("[USP] HandleAddRequest received invalid message")
		return
	}

	msgId := inComingMsg.Header.MsgId
	logger.Infof("[USP] receive ADD request: %s", inComingMsg.String())

	getCreateObjs := inComingMsg.GetBody().GetRequest().GetAdd().GetCreateObjs()

	var affectedPath []string
	var requestPath []string
	var updatedParams []map[string]string

	for _, createObj := range getCreateObjs {
		path := createObj.GetObjPath()
		nodePath := uc.getNewInstance(path)

		requestPath = append(requestPath, path)
		affectedPath = append(affectedPath, nodePath)

		// 使用辅助函数提取参数设置
		paramSettings := extractParamSettings(createObj.GetParamSettings())
		// 将参数设置应用到 repository
		for setKey, setValue := range paramSettings {
			uc.DataRepo.SetValue(nodePath, setKey, setValue)
		}
		updatedParams = append(updatedParams, paramSettings)

		// 处理对象创建后的副作用（订阅注册、通知发送）
		uc.handleObjectCreationSideEffects(path, paramSettings)
	}

	msg := utils.CreateAddResponseMessage(msgId, requestPath, affectedPath, updatedParams)
	logger.Infof("[USP] send ADD response: %s", msg.String())

	err := uc.HandleMTPMsgTransmit(msg)
	if err != nil {
		logger.Warnf("[USP] ADD error: msgId=%s, err=%v", msgId, err)
	}
}

// handleObjectCreationSideEffects 处理对象创建后的副作用
// 包括：订阅注册（如果是订阅节点）、发送对象创建通知
func (uc *ClientUseCase) handleObjectCreationSideEffects(path string, paramSettings map[string]string) {
	// 尝试注册订阅（如果是订阅节点）
	if err := uc.HandleAddLocalAgentSubscription(path, paramSettings); err != nil {
		logger.Warnf("[USP] ADD subscription registration error: path=%s, err=%v", path, err)
		return
	}

	// 发送对象创建通知
	uc.notifyObjectCreation(path, paramSettings)
}

// HandleDeleteRequest handles incoming DELETE requests
func (uc *ClientUseCase) HandleDeleteRequest(inComingMsg *api.Msg) {
	// 防御性检查
	if inComingMsg == nil || inComingMsg.Header == nil {
		logger.Warnf("[USP] HandleDeleteRequest received invalid message")
		return
	}

	msgId := inComingMsg.Header.MsgId
	logger.Infof("[USP] receive DELETE request: %s", inComingMsg.String())

	objPaths := inComingMsg.GetBody().GetRequest().GetDelete().GetObjPaths()
	var affectedPath []string
	var requestPath []string

	for _, objPath := range objPaths {
		// 处理对象删除前的副作用（取消订阅）
		uc.handleObjectDeletionPreEffects(objPath)

		// 执行删除操作
		requestPath = append(requestPath, objPath)
		nodePath, isFound := uc.DataRepo.DeleteNode(objPath)
		if isFound {
			affectedPath = append(affectedPath, nodePath)
		} else {
			affectedPath = append(affectedPath, objPath)
		}

		// 处理对象删除后的副作用（发送通知）
		uc.notifyObjectDeletion(objPath)
	}

	msg := utils.CreateDeleteResponseMessage(msgId, requestPath, affectedPath)
	logger.Infof("[USP] send DELETE response: %s", msg.String())

	err := uc.HandleMTPMsgTransmit(msg)
	if err != nil {
		logger.Warnf("[USP] DELETE error: msgId=%s, err=%v", msgId, err)
	}
}

// handleObjectDeletionPreEffects 处理对象删除前的副作用
// 包括：取消订阅（如果是订阅节点），防止节点被删除后导致超时
func (uc *ClientUseCase) handleObjectDeletionPreEffects(objPath string) {
	if err := uc.HandleDeleteLocalAgentSubscription(objPath); err != nil {
		logger.Warnf("[USP] DELETE subscription cleanup error: path=%s, err=%v", objPath, err)
		// 继续执行，不要因为一个错误而中断删除操作
	}
}

// HandleOperateRequest handles incoming OPERATE requests
func (uc *ClientUseCase) HandleOperateRequest(inComingMsg *api.Msg) {
	// 防御性检查
	if inComingMsg == nil || inComingMsg.Header == nil {
		logger.Warnf("[USP] HandleOperateRequest received invalid message")
		return
	}
	logger.Infof("[USP] receive OPERATE request: %s", inComingMsg.String())

	// 根据command名称决定调用operComplete还是event
	uc.HandleCommand(inComingMsg.GetBody().GetRequest().GetOperate(), inComingMsg.Header.MsgId)
}

// HandleNotifyResp handles incoming NOTIFY_RESP messages
func (uc *ClientUseCase) HandleNotifyResp(inComingMsg *api.Msg) {
	// 防御性检查
	if inComingMsg == nil || inComingMsg.Header == nil {
		logger.Warnf("[USP] HandleNotifyResp received invalid message")
		return
	}

	logger.Infof("[USP] receive NOTIFY_RESP: %s", inComingMsg.String())
}
