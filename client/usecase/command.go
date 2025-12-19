package usecase

import (
	"tr369-wss-client/client/model"
	logger "tr369-wss-client/log"
	"tr369-wss-client/pkg/api"
	"tr369-wss-client/utils"
)

// HandleCommand handles command
func (uc *ClientUseCase) HandleCommand(operate *api.Operate, msgId string) {
	command := operate.GetCommand()

	switch command {
	case model.DeviceReboot:
		// 处理boot
		uc.HandleDeviceBoot(operate, msgId)
	default:
		logger.Infof("[USP] unknown command: %s", command)

	}
}

func (uc *ClientUseCase) HandleDeviceBoot(operate *api.Operate, msgId string) {

	// 构建resp
	var operationResults []*api.OperateResp_OperationResult

	operationResult := &api.OperateResp_OperationResult{
		ExecutedCommand: operate.GetCommand(),
		OperationResp: &api.OperateResp_OperationResult_ReqObjPath{
			ReqObjPath: "Device.LocalAgent.Request.1",
		},
	}
	operationResults = append(operationResults, operationResult)

	msg := utils.CreateOperateResponseMessage(msgId, operationResults)
	logger.Infof("[USP] send OPERATE response: %s", msg.String())

	err := uc.HandleMTPMsgTransmit(msg)
	if err != nil {
		logger.Warnf("[USP] OPERATE error: msgId=%s, err=%v", msgId, err)
	}

	// 发送Boot! 事件
	// 构建boot事件参数
	params := map[string]string{
		"CommandKey":      "",
		"Cause":           "RemoteReboot",
		"Reason":          "",
		"FirmwareUpdated": "false",
		"ParameterMap":    "",
	}

	uc.notifyEvent(model.DeviceReboot, model.BOOT, params)
}

//func (uc *ClientUseCase) HandleUrlUpgrade(operate *api.Operate, msgId string) {
//
//	// 构建resp
//	var operationResults []*api.OperateResp_OperationResult
//
//	operationResult := &api.OperateResp_OperationResult{
//		ExecutedCommand: operate.GetCommand(),
//		OperationResp: &api.OperateResp_OperationResult_ReqObjPath{
//			ReqObjPath: "Device.LocalAgent.Request.1",
//		},
//	}
//	operationResults = append(operationResults, operationResult)
//
//	msg := utils.CreateOperateResponseMessage(msgId, operationResults)
//	logger.Infof("[USP] send OPERATE response: %s", msg.String())
//
//	err := uc.HandleMTPMsgTransmit(msg)
//	if err != nil {
//		logger.Warnf("[USP] OPERATE error: msgId=%s, err=%v", msgId, err)
//	}
//
//	// 发送升级事件
//	uc.notifyOperComplete("", operate.GetCommandKey(), operate.GetCommand())
//}
