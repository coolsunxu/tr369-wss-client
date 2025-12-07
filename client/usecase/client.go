package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"tr369-wss-client/client/model"
	"tr369-wss-client/config"
	logger "tr369-wss-client/log"
	"tr369-wss-client/pkg/api"
	tr181Model "tr369-wss-client/tr181/model"
	"tr369-wss-client/utils"
)

type ClientUseCase struct {
	model.ClientUseCase // 嵌入接口以确保实现所有方法
	Config              *config.Config
	TR181DataModel      *tr181Model.TR181DataModel
	ClientRepository    model.ClientRepository
	writeCount          int
	lastWriteTime       int64
	pingTicker          *time.Ticker
	ctx                 context.Context
	cancel              context.CancelFunc
	messageChannel      chan []byte // 消息发送通道
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
		slog.Warn("SendOperateCompleteNotify ", "error", err)
	}

	// 测试
	uc.HandleSubscription("Device.DeviceInfo.sn", "ValueChange")

	uc.ClientRepository.NotifyListeners("Device.DeviceInfo.sn", &api.Notify_ValueChange{
		ParamPath:  "Device.DeviceInfo.sn",
		ParamValue: "1234567890",
	})
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
		slog.Warn("SendOperateCompleteNotify ", "error", err)
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
	}

	uc.ClientRepository.SaveData()

	msg := utils.CreateAddResponseMessage(inComingMsg.Header.MsgId, requestPath, affectedPath, updatedParams)
	logger.Infof("client sent add resp message %s", msg)

	err := uc.HandleMTPMsgTransmit(msg)
	if err != nil {
		slog.Warn("SendOperateCompleteNotify ", "error", err)
	}
}

// HandleDeleteRequest handles incoming DELETE requests
func (uc *ClientUseCase) HandleDeleteRequest(inComingMsg *api.Msg) {
	logger.Infof("client receive delete usp msg %s", inComingMsg.String())

	objPaths := inComingMsg.GetBody().GetRequest().GetDelete().GetObjPaths()
	var affectedPath []string
	var requestPath []string

	for _, objPath := range objPaths {
		requestPath = append(requestPath, objPath)
		nodePath, isFound := uc.ClientRepository.HandleDeleteRequest(objPath)
		if isFound {
			affectedPath = append(affectedPath, nodePath)
		} else {
			affectedPath = append(affectedPath, objPath)
		}
	}

	uc.ClientRepository.SaveData()

	msg := utils.CreateDeleteResponseMessage(inComingMsg.Header.MsgId, requestPath, affectedPath)
	logger.Infof("client sent delete resp message %s", msg)

	err := uc.HandleMTPMsgTransmit(msg)
	if err != nil {
		slog.Warn("SendOperateCompleteNotify ", "error", err)
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
		slog.Warn("SendOperateCompleteNotify ", "error", err)
	}
}

func (uc *ClientUseCase) SendOperateCompleteNotify(objPath string, commandName string, commandKey string, outputArgs map[string]string) {

	msg := utils.CreateOperateCompleteMessage(objPath, commandName, commandKey, outputArgs)
	err := uc.HandleMTPMsgTransmit(msg)
	if err != nil {
		slog.Warn("SendOperateCompleteNotify ", "error", err)
	}
}

func (uc *ClientUseCase) HandleMTPMsgTransmit(msg *api.Msg) error {

	rec := utils.CreateUspRecordNoSession("1.0", uc.Config.WebsocketConfig.EndpointId, uc.Config.WebsocketConfig.ControllerId, msg)
	payload, _ := utils.EncodeUspRecord(rec)
	uc.messageChannel <- payload
	return nil
}

// HandleValueChange 处理 value change 事件
func (uc *ClientUseCase) HandleValueChange(change interface{}) {
	valueChange, ok := change.(*api.Notify_ValueChange)

	if !ok {
		return
	}

	logger.Infof("client receive value change usp msg %s", valueChange.String())

}

// HandleSubscription 处理订阅消息
func (uc *ClientUseCase) HandleSubscription(path string, subscriptionType string) error {
	switch subscriptionType {
	case tr181Model.ValueChange:
		// 处理 value change 订阅
		err := uc.ClientRepository.AddListener(path, uc.HandleValueChange)
		if err != nil {
			return err
		}

	case tr181Model.ObjectCreation:
		// 处理对象创建订阅
	case tr181Model.ObjectDeletion:
		// 处理对象删除订阅
	case tr181Model.OperationComplete:
		// 处理操作完成订阅
	default:
		return fmt.Errorf("unknown subscription type %s", subscriptionType)
	}

	return nil
}
