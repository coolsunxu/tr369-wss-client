package usecase

import (
	"context"
	"fmt"
	"regexp"
	"tr369-wss-client/client/model"
	"tr369-wss-client/config"
	logger "tr369-wss-client/log"
	"tr369-wss-client/pkg/api"
	tr181Model "tr369-wss-client/tr181/model"
	"tr369-wss-client/utils"
)

type ClientUseCase struct {
	Config           *config.Config
	ClientRepository model.ClientRepository
	ctx              context.Context
	messageChannel   chan []byte // 消息发送通道
}

// NewClientUseCase creates a new client use case instance
func NewClientUseCase(
	ctx context.Context,
	cfg *config.Config,
	repo model.ClientRepository,
	messageChannel chan []byte,
) *ClientUseCase {
	return &ClientUseCase{
		ctx:              ctx,
		Config:           cfg,
		ClientRepository: repo,
		messageChannel:   messageChannel,
	}
}

// HandleMessage processes incoming USP messages
func (uc *ClientUseCase) HandleMessage(msg *api.Msg) {
	// 根据消息类型处理不同的请求
	switch msg.Header.MsgType {
	case api.Header_GET:
		uc.HandleGetRequest(msg)
	case api.Header_SET:
		uc.HandleSetRequest(msg)
	case api.Header_ADD:
		uc.HandleAddRequest(msg)
	case api.Header_DELETE:
		uc.HandleDeleteRequest(msg)
	case api.Header_OPERATE:
		uc.HandleOperateRequest(msg)
	case api.Header_NOTIFY_RESP:
		uc.HandleNotifyResp(msg)
	default:
		logger.Infof("Unknown message type: %v", msg.Header.MsgType)
	}
}

// HandleGetRequest handles incoming GET requests
func (uc *ClientUseCase) HandleGetRequest(inComingMsg *api.Msg) {
	logger.Infof("client receive get usp msg %s", inComingMsg.String())

	getNodePaths := inComingMsg.GetBody().GetRequest().GetGet().GetParamPaths()
	resp := uc.ClientRepository.ConstructGetResp(getNodePaths)
	msg := utils.CreateGetResponseMessage(inComingMsg.Header.MsgId, resp)
	logger.Infof("client sent get resp message %s", msg)

	err := uc.HandleMTPMsgTransmit(msg)
	if err != nil {
		logger.Warnf("SendOperateCompleteNotify error: %v", err)
	}

}

// HandleSetRequest handles incoming SET requests
func (uc *ClientUseCase) HandleSetRequest(inComingMsg *api.Msg) {
	logger.Infof("client receive set usp msg %s", inComingMsg.String())

	getUpdateObjs := inComingMsg.GetBody().GetRequest().GetSet().GetUpdateObjs()

	var affectedPath []string
	var requestPath []string
	var updatedParams []map[string]string
	paramSettings := make(map[string]string)

	for _, UpdateObj := range getUpdateObjs {
		path := UpdateObj.GetObjPath()
		isSuccess, nodePath := uc.ClientRepository.IsExistPath(path)
		if !isSuccess {
			continue
		}
		requestPath = append(requestPath, path)
		affectedPath = append(affectedPath, nodePath)
		for _, paramSetting := range UpdateObj.GetParamSettings() {
			setKey := paramSetting.GetParam()
			setValue := paramSetting.GetValue()
			paramSettings[setKey] = setValue
			uc.ClientRepository.HandleSetRequest(nodePath, setKey, setValue)
		}
		updatedParams = append(updatedParams, paramSettings)
	}

	uc.ClientRepository.SaveData()

	msg := utils.CreateSetResponseMessage(inComingMsg.Header.MsgId, requestPath, affectedPath, updatedParams)
	logger.Infof("client sent set resp message %s", msg)

	err := uc.HandleMTPMsgTransmit(msg)
	if err != nil {
		logger.Warnf("SendOperateCompleteNotify error: %v", err)
	}
}

// HandleAddRequest handles incoming ADD requests
func (uc *ClientUseCase) HandleAddRequest(inComingMsg *api.Msg) {
	logger.Infof("client receive add usp msg %s", inComingMsg.String())

	getCreateObjs := inComingMsg.GetBody().GetRequest().GetAdd().GetCreateObjs()

	var affectedPath []string
	var requestPath []string
	var updatedParams []map[string]string
	paramSettings := make(map[string]string)

	for _, createObj := range getCreateObjs {
		path := createObj.GetObjPath()
		nodePath := uc.ClientRepository.GetNewInstance(path)

		requestPath = append(requestPath, path)
		affectedPath = append(affectedPath, nodePath)

		for _, paramSetting := range createObj.GetParamSettings() {
			setKey := paramSetting.GetParam()
			setValue := paramSetting.GetValue()
			paramSettings[setKey] = setValue
			uc.ClientRepository.HandleSetRequest(nodePath, setKey, setValue)
		}
		updatedParams = append(updatedParams, paramSettings)

		// 判断是否是新建订阅
		err := uc.HandleAddLocalAgentSubscription(path, paramSettings)
		if err != nil {
			continue
		}

		// 判断是否需要发送obj creation
		notifyCreateObj := &api.Notify_ObjCreation{
			ObjCreation: &api.Notify_ObjectCreation{
				ObjPath:    path,
				UniqueKeys: paramSettings,
			},
		}
		uc.ClientRepository.NotifyListeners(path, notifyCreateObj)
	}

	uc.ClientRepository.SaveData()

	msg := utils.CreateAddResponseMessage(inComingMsg.Header.MsgId, requestPath, affectedPath, updatedParams)
	logger.Infof("client sent add resp message %s", msg)

	err := uc.HandleMTPMsgTransmit(msg)
	if err != nil {
		logger.Warnf("SendOperateCompleteNotify error: %v", err)
	}
}

// HandleDeleteRequest handles incoming DELETE requests
func (uc *ClientUseCase) HandleDeleteRequest(inComingMsg *api.Msg) {
	logger.Infof("client receive delete usp msg %s", inComingMsg.String())

	objPaths := inComingMsg.GetBody().GetRequest().GetDelete().GetObjPaths()
	var affectedPath []string
	var requestPath []string

	for _, objPath := range objPaths {

		// 先删除相关Listener,防止节点被删除，导致超时
		err := uc.HandleDeleteLocalAgentSubscription(objPath)
		if err != nil {
			logger.Infof("handleDeleteLocalAgentSubscription error: %v", err)
			return
		}

		requestPath = append(requestPath, objPath)
		nodePath, isFound := uc.ClientRepository.HandleDeleteRequest(objPath)
		if isFound {
			affectedPath = append(affectedPath, nodePath)
		} else {
			affectedPath = append(affectedPath, objPath)
		}

		// 判断是否需要发送obj creation
		notifyDeleteObj := &api.Notify_ObjDeletion{
			ObjDeletion: &api.Notify_ObjectDeletion{
				ObjPath: objPath,
			},
		}
		uc.ClientRepository.NotifyListeners(objPath, notifyDeleteObj)
	}

	uc.ClientRepository.SaveData()

	msg := utils.CreateDeleteResponseMessage(inComingMsg.Header.MsgId, requestPath, affectedPath)
	logger.Infof("client sent delete resp message %s", msg)

	err := uc.HandleMTPMsgTransmit(msg)
	if err != nil {
		logger.Warnf("SendOperateCompleteNotify error: %v", err)
	}
}

// HandleOperateRequest handles incoming OPERATE requests
func (uc *ClientUseCase) HandleOperateRequest(inComingMsg *api.Msg) {
	logger.Infof("client receive operate usp msg %s", inComingMsg.String())

	operate := inComingMsg.GetBody().GetRequest().GetOperate()
	command := operate.GetCommand()

	msg := utils.CreateOperateResponseMessage(inComingMsg.Header.MsgId, command)
	logger.Infof("client sent operate resp message %s", msg)

	err := uc.HandleMTPMsgTransmit(msg)
	if err != nil {
		logger.Warnf("SendOperateCompleteNotify error: %v", err)
	}
}

func (uc *ClientUseCase) HandleNotifyResp(inComingMsg *api.Msg) {
	logger.Infof("client receive notify resp msg %s", inComingMsg.String())
}

func (uc *ClientUseCase) SendOperateCompleteNotify(objPath string, commandName string, commandKey string, outputArgs map[string]string) {

	msg := utils.CreateOperateCompleteMessage(objPath, commandName, commandKey, outputArgs)
	err := uc.HandleMTPMsgTransmit(msg)
	if err != nil {
		logger.Warnf("SendOperateCompleteNotify error: %v", err)
	}
}

func (uc *ClientUseCase) HandleMTPMsgTransmit(msg *api.Msg) error {

	rec := utils.CreateUspRecordNoSession("1.0", uc.Config.WebsocketConfig.EndpointId, uc.Config.WebsocketConfig.ControllerId, msg)
	payload, err := utils.EncodeUspRecord(rec)
	if err != nil {
		logger.Errorf("Failed to encode USP record: %v", err)
		return err
	}
	uc.messageChannel <- payload
	return nil
}

func (uc *ClientUseCase) HandleAddLocalAgentSubscription(requestPath string, paramSettings map[string]string) error {

	// 为空则直接返回
	if requestPath == "" {
		return nil
	}

	// 判断路径是否是订阅节点
	if requestPath != tr181Model.DeviceLocalAgentSubscription {
		return nil
	}

	referenceList := ""
	subscriptionId := ""
	subscriptionType := ""
	for key, value := range paramSettings {
		switch key {
		case "ID":
			subscriptionId = value
		case "ReferenceList":
			referenceList = value
		case "NotifType":
			subscriptionType = value
		default:
		}
	}

	return uc.HandleSubscription(referenceList, subscriptionId, subscriptionType)
}

func (uc *ClientUseCase) HandleDeleteLocalAgentSubscription(requestPath string) error {

	// 为空则直接返回
	if requestPath == "" {
		return nil
	}

	// 判断是否是各级父节点
	if requestPath == "Device." || requestPath == "Device.LocalAgent." || requestPath == "Device.LocalAgent.Subscription." {
		// 清除所有订阅
		logger.Infof("delete device all local agent subscription %s", requestPath)
		return uc.ClientRepository.ResetListener()
	}

	// 判断是否符合 Device.LocalAgent.Subscription.*. 格式，*为正整数
	matched, err := regexp.MatchString(`^Device\.LocalAgent\.Subscription\.([1-9]\d*)\.$`, requestPath)
	if err != nil {
		return err
	}

	if matched {
		// 符合格式，处理删除逻辑
		logger.Infof("Handle delete subscription for path: %s", requestPath)
		pathName, err := uc.ClientRepository.GetValueByPath(requestPath + "ReferenceList")
		if err != nil {
			return err
		}
		// 这里可以添加实际的删除订阅逻辑，比如调用repository的RemoveListener
		_ = uc.ClientRepository.RemoveListener(pathName.(string))
	}

	return nil
}

// HandleSubscription 处理订阅消息
func (uc *ClientUseCase) HandleSubscription(path string, subscriptionId string, subscriptionType string) error {
	switch subscriptionType {
	case tr181Model.ValueChange:
		// 处理 value change 订阅
		err := uc.ClientRepository.AddListener(path, tr181Model.Listener{
			SubscriptionId: subscriptionId,
			Listener:       uc.HandleValueChange,
		})
		if err != nil {
			return err
		}

	case tr181Model.ObjectCreation:
		// 处理对象创建订阅
		err := uc.ClientRepository.AddListener(path, tr181Model.Listener{
			SubscriptionId: subscriptionId,
			Listener:       uc.HandleObjectCreation,
		})
		if err != nil {
			return err
		}
	case tr181Model.ObjectDeletion:
		// 处理对象删除订阅
		err := uc.ClientRepository.AddListener(path, tr181Model.Listener{
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

// HandleValueChange 处理 value change 事件
func (uc *ClientUseCase) HandleValueChange(subscriptionId string, change interface{}) {
	valueChange, ok := change.(*api.Notify_ValueChange_)

	// 判断类型
	if !ok {
		logger.Infof("Receive message is not Notify_ValueChange")
		return
	}

	// 构建notify消息
	notify := &api.Notify{
		SubscriptionId: subscriptionId,
		SendResp:       true,
		Notification:   valueChange,
	}
	msg := utils.CreateNotifyMessage(notify)

	// 发送消息，通过channel
	err := uc.HandleMTPMsgTransmit(msg)
	if err != nil {
		logger.Warnf("SendNotify error: %v", err)
	}

	logger.Infof("client send value change usp msg %s", msg)

}

// HandleObjectCreation 处理 object creation 事件
func (uc *ClientUseCase) HandleObjectCreation(subscriptionId string, change interface{}) {
	objCreation, ok := change.(*api.Notify_ObjCreation)

	// 判断类型
	if !ok {
		logger.Infof("Receive message is not Notify_ObjCreation")
		return
	}

	// 构建notify消息
	notify := &api.Notify{
		SubscriptionId: subscriptionId,
		SendResp:       true,
		Notification:   objCreation,
	}
	msg := utils.CreateNotifyMessage(notify)

	// 发送消息，通过channel
	err := uc.HandleMTPMsgTransmit(msg)
	if err != nil {
		logger.Warnf("SendNotify error: %v", err)
	}

	logger.Infof("client send obj creation usp msg %s", msg)

}

// HandleObjectDeletion 处理 object deletion 事件
func (uc *ClientUseCase) HandleObjectDeletion(subscriptionId string, change interface{}) {
	ObjDeletion, ok := change.(*api.Notify_ObjDeletion)

	// 判断类型
	if !ok {
		logger.Infof("Receive message is not Notify_ObjDeletion")
		return
	}

	// 构建notify消息
	notify := &api.Notify{
		SubscriptionId: subscriptionId,
		SendResp:       true,
		Notification:   ObjDeletion,
	}
	msg := utils.CreateNotifyMessage(notify)

	// 发送消息，通过channel
	err := uc.HandleMTPMsgTransmit(msg)
	if err != nil {
		logger.Warnf("SendNotify error: %v", err)
	}

	logger.Infof("client send obj deletion usp msg %s", msg)

}
