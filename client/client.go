package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"
	"tr369-wss-client/client/model"
	"tr369-wss-client/utils"

	"google.golang.org/protobuf/proto"
	"nhooyr.io/websocket"

	"tr369-wss-client/config"
	"tr369-wss-client/pkg/api"
)

// contains 检查字符串s是否包含子串substr
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// WSClient represents the TR369 WebSocket client
type WSClient struct {
	model.WSClient   // 嵌入接口以确保实现所有方法
	config           *config.Config
	conn             *websocket.Conn // nhooyr.io/websocket 已迁移至 github.com/coder/websocket，但类型名保持不变
	ctx              context.Context
	cancel           context.CancelFunc
	connected        bool
	pingTicker       *time.Ticker
	clientRepository model.ClientRepository
}

// NewWSClient creates a new WebSocket client instance
func NewWSClient(
	cfg *config.Config,
	clientRepository model.ClientRepository,
) *WSClient {
	ctx, cancel := context.WithCancel(context.Background())
	return &WSClient{
		config:           cfg,
		ctx:              ctx,
		cancel:           cancel,
		connected:        false,
		clientRepository: clientRepository,
	}
}

// Connect establishes a WebSocket connection to the server
func (c *WSClient) Connect() error {

	options := &websocket.DialOptions{
		Subprotocols: []string{"v1.usp"},
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	// 设置TR369协议必需的HTTP头
	headers := http.Header{}
	options.HTTPHeader = headers

	connectUrl := c.config.ServerURL

	// 检查apiL是否已经包含查询参数
	connectUrl += "?eid=" + c.config.EndpointId

	log.Printf("Connecting with eid in apiL: %s", connectUrl)

	// 连接服务器
	conn, _, err := websocket.Dial(c.ctx, connectUrl, options)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}

	c.conn = conn
	c.connected = true

	// 设置读消息的最大大小
	c.conn.SetReadLimit(c.config.MaxMessageSize)

	// 配置ping/pong
	c.pingTicker = time.NewTicker(c.config.PingInterval)

	return nil
}

// Disconnect closes the WebSocket connection
func (c *WSClient) Disconnect() {
	if c.pingTicker != nil {
		c.pingTicker.Stop()
	}

	if c.conn != nil {
		// 先发送关闭帧
		c.conn.Close(websocket.StatusNormalClosure, "client disconnecting")
		c.conn = nil
	}

	c.cancel()
	c.connected = false
}

// StartMessageHandler starts the message handling goroutines
func (c *WSClient) StartMessageHandler() {
	// 启动ping goroutine
	go c.pingHandler()

	// 启动消息接收goroutine
	go c.messageHandler()

	// 启动消息写goroutine，顺序写入，便于控制
}

// pingHandler handles periodic ping messages
func (c *WSClient) pingHandler() {
	for {
		select {
		case <-c.pingTicker.C:
			if err := c.conn.Ping(c.ctx); err != nil {
				log.Printf("Ping failed: %v", err)
				c.Disconnect()
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// messageHandler handles incoming messages using protobuf
func (c *WSClient) messageHandler() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			// 读取二进制消息
			_, data, err := c.conn.Read(c.ctx)
			if err != nil {
				if errors.Is(c.ctx.Err(), context.Canceled) {
					return
				}
				if status := websocket.CloseStatus(err); status != -1 {
					log.Printf("Connection closed with status %d: %v", status, err)
					c.Disconnect()
					return
				} else {
					log.Printf("Read error: %v", err)
				}
				return
			}

			// 使用protobuf解码消息，先解析为Record
			record := new(api.Record)
			if err := proto.Unmarshal(data, record); err != nil {
				log.Printf("Failed to unmarshal Record: %v", err)
				continue
			}

			// 提取NoSessionContextRecord
			noSessionContext := record.GetNoSessionContext()
			if noSessionContext == nil {
				log.Printf("Record is not NoSessionContextRecord")
				continue
			}

			// 从NoSessionContextRecord中提取payload并解析为api.Msg
			var msg api.Msg
			if err := proto.Unmarshal(noSessionContext.GetPayload(), &msg); err != nil {
				log.Printf("Failed to unmarshal api.Msg from payload: %v", err)
				continue
			}

			// 处理接收到的消息
			c.handleProtobufMessage(&msg)
		}
	}
}

// handleProtobufMessage processes incoming TR369 protobuf messages
func (c *WSClient) handleProtobufMessage(msg *api.Msg) {
	// 根据消息类型处理不同的请求
	switch msg.Header.MsgType {
	case api.Header_GET:
		c.HandleGetRequest(msg)
	case api.Header_SET:
		c.HandleSetRequest(msg)
	case api.Header_ADD:
		c.HandleAddRequest(msg)
	case api.Header_DELETE:
		c.HandleDeleteRequest(msg)
	case api.Header_OPERATE:
		c.HandleOperateRequest(msg)
	default:
		log.Printf("Unknown message type: %v", msg.Header.MsgType)
	}
}

func (c *WSClient) HandleGetRequest(inComingMsg *api.Msg) {
	getNodePaths := inComingMsg.GetBody().GetRequest().GetGet().GetParamPaths()
	msg := c.createGetResponseMessage(inComingMsg.Header.MsgId, getNodePaths)
	slog.Debug("sent get resp message", "msg", msg)
	err := c.HandleMTPMsgTransmit(msg)
	if err != nil {
		slog.Warn("HandleGetRequest ", "error", err)
	}
}

func (c *WSClient) createGetResponseMessage(msgId string, getNodePaths []string) (result *api.Msg) {
	resp := c.clientRepository.ConstructGetResp(getNodePaths)
	return utils.CreateGetResponseMessage(msgId, resp)
}

func (c *WSClient) HandleSetRequest(inComingMsg *api.Msg) {
	getUpdateObjs := inComingMsg.GetBody().GetRequest().GetSet().GetUpdateObjs()

	var affectedPath []string
	var requestPath []string
	var updatedParams []map[string]string
	paramSettings := make(map[string]string)

	slog.Debug("get set req from server", "req", getUpdateObjs)

	for _, UpdateObj := range getUpdateObjs {
		path := UpdateObj.GetObjPath()
		isSuccess, nodePath := c.clientRepository.IsExistPath(path)
		if !isSuccess {
			continue
		}
		requestPath = append(requestPath, path)
		affectedPath = append(affectedPath, nodePath)
		for _, paramSetting := range UpdateObj.GetParamSettings() {
			setKey := paramSetting.GetParam()
			setValue := paramSetting.GetValue()
			paramSettings[setKey] = setValue
			c.clientRepository.HandleSetRequest(nodePath, setKey, setValue)
		}
		updatedParams = append(updatedParams, paramSettings)
	}

	//c.repo.SaveData(dataMap)

	msg := utils.CreateSetResponseMessage(inComingMsg.Header.MsgId, requestPath, affectedPath, updatedParams)
	slog.Debug("sent set resp message")
	err := c.HandleMTPMsgTransmit(msg)
	if err != nil {
		slog.Warn("HandleSetRequest", "error", err)
	}

}

func (c *WSClient) HandleAddRequest(inComingMsg *api.Msg) {
	getCreateObjs := inComingMsg.GetBody().GetRequest().GetAdd().GetCreateObjs()

	var affectedPath []string
	var requestPath []string
	var updatedParams []map[string]string
	paramSettings := make(map[string]string)

	for _, createObj := range getCreateObjs {
		path := createObj.GetObjPath()
		nodePath := c.clientRepository.GetNewInstance(path)

		requestPath = append(requestPath, path)
		affectedPath = append(affectedPath, nodePath)

		for _, paramSetting := range createObj.GetParamSettings() {
			setKey := paramSetting.GetParam()
			setValue := paramSetting.GetValue()
			paramSettings[setKey] = setValue
			c.clientRepository.HandleSetRequest(nodePath, setKey, setValue)
		}
		updatedParams = append(updatedParams, paramSettings)
	}

	//c.repo.SaveData(dataMap)
	msg := utils.CreateAddResponseMessage(inComingMsg.Header.MsgId, requestPath, affectedPath, updatedParams)
	slog.Debug("sent add resp message")
	err := c.HandleMTPMsgTransmit(msg)
	if err != nil {
		slog.Warn("HandleAddRequest error:", "error", err)
	}
}

func (c *WSClient) HandleDeleteRequest(inComingMsg *api.Msg) {
	objPaths := inComingMsg.GetBody().GetRequest().GetDelete().GetObjPaths()
	var affectedPath []string
	var requestPath []string

	for _, objPath := range objPaths {
		requestPath = append(requestPath, objPath)
		nodePath, isFound := c.clientRepository.HandleDeleteRequest(objPath)
		if isFound {
			affectedPath = append(affectedPath, nodePath)
		} else {
			affectedPath = append(affectedPath, objPath)
		}

	}
	//c.repo.SaveData(dataMap)
	msg := utils.CreateDeleteResponseMessage(inComingMsg.Header.MsgId, requestPath, affectedPath)
	slog.Debug("sent delete resp message")
	err := c.HandleMTPMsgTransmit(msg)
	if err != nil {
		slog.Warn("HandleDeleteRequest error:", "error", err)
	}
}

func (c *WSClient) HandleOperateRequest(inComingMsg *api.Msg) {
	operate := inComingMsg.GetBody().GetRequest().GetOperate()

	slog.Debug("get operate req", "req", operate)

	command := operate.GetCommand()

	msg := utils.CreateOperateResponseMessage(inComingMsg.Header.MsgId, command)
	slog.Debug("sent operate resp message")
	err := c.HandleMTPMsgTransmit(msg)
	if err != nil {
		slog.Warn("HandleOperateRequest error:", "error", err)
	}
}

func (c *WSClient) SendOperateCompleteNotify(objPath string, commandName string, commandKey string, outputArgs map[string]string) {

	msg := utils.CreateOperateCompleteMessage(objPath, commandName, commandKey, outputArgs)
	err := c.HandleMTPMsgTransmit(msg)
	if err != nil {
		slog.Warn("SendOperateCompleteNotify ", "error", err)
	}
}

func (c *WSClient) HandleMTPMsgTransmit(msg *api.Msg) error {

	rec := utils.CreateUspRecordNoSession("1.0", c.config.EndpointId, c.config.ControllerIdentifier, msg)
	payload, _ := utils.EncodeUspRecord(rec)
	c.sendResponse(payload)
	return nil
}

// sendResponse sends a protobuf response message
func (c *WSClient) sendResponse(payload []byte) {
	ctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
	defer cancel()

	// 发送二进制消息
	if err := c.conn.Write(ctx, websocket.MessageBinary, payload); err != nil {
		log.Printf("Failed to send response: %v", err)
	}
}
