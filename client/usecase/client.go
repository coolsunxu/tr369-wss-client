package usecase

import (
	"context"
	"tr369-wss-client/client/model"
	"tr369-wss-client/config"
	logger "tr369-wss-client/log"
	"tr369-wss-client/pkg/api"
	"tr369-wss-client/utils"
)

// ClientUseCase 客户端业务逻辑处理器
type ClientUseCase struct {
	Config         *config.Config
	DataRepo       model.DataRepository  // 数据访问接口
	ListenerMgr    model.ListenerManager // 监听器管理接口
	ctx            context.Context
	messageChannel chan []byte // 消息发送通道
}

// NewClientUseCase creates a new client use case instance
func NewClientUseCase(
	ctx context.Context,
	cfg *config.Config,
	dataRepo model.DataRepository,
	listenerMgr model.ListenerManager,
	messageChannel chan []byte,
) *ClientUseCase {
	return &ClientUseCase{
		ctx:            ctx,
		Config:         cfg,
		DataRepo:       dataRepo,
		ListenerMgr:    listenerMgr,
		messageChannel: messageChannel,
	}
}

// HandleMessage processes incoming USP messages
func (uc *ClientUseCase) HandleMessage(msg *api.Msg) {
	// 防御性检查：检查 msg 是否为 nil
	if msg == nil {
		logger.Warnf("[USP] received nil message, ignoring")
		return
	}

	// 防御性检查：检查 msg.Header 是否为 nil
	if msg.Header == nil {
		logger.Warnf("[USP] received message with nil header, ignoring")
		return
	}

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
		logger.Warnf("[USP] UNKNOWN: unsupported message type=%v, msgId=%s", msg.Header.MsgType, msg.Header.MsgId)
	}
}

// HandleMTPMsgTransmit 发送 MTP 消息
func (uc *ClientUseCase) HandleMTPMsgTransmit(msg *api.Msg) error {
	rec := utils.CreateUspRecordNoSession(uc.Config.Tr369Config.Version, uc.Config.WebsocketConfig.EndpointId, uc.Config.WebsocketConfig.ControllerId, msg)
	payload, err := utils.EncodeUspRecord(rec)
	if err != nil {
		return err
	}

	// 发送消息到通道
	uc.messageChannel <- payload
	return nil
}
