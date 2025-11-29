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
	"tr369-wss-client/client/repository"
	"tr369-wss-client/common"

	"google.golang.org/protobuf/proto"
	"nhooyr.io/websocket"

	"tr369-wss-client/config"
	"tr369-wss-client/datamodel"
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
	dataModel        *datamodel.TR181DataModel
	conn             *websocket.Conn // nhooyr.io/websocket 已迁移至 github.com/coder/websocket，但类型名保持不变
	ctx              context.Context
	cancel           context.CancelFunc
	connected        bool
	pingTicker       *time.Ticker
	clientRepository *repository.ClientRepository
}

// NewWSClient creates a new WebSocket client instance
func NewWSClient(cfg *config.Config, model *datamodel.TR181DataModel) *WSClient {
	ctx, cancel := context.WithCancel(context.Background())
	return &WSClient{
		config:           cfg,
		dataModel:        model,
		ctx:              ctx,
		cancel:           cancel,
		connected:        false,
		clientRepository: repository.NewClientRepository("datamodel/default_tr181_nodes.json"),
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

	// 根据要求，需要将eid添加为apiL查询参数，格式为 ?eid=deviceid

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

	msg := c.createGetResponseMessage(getNodePaths)
	msg.Header.MsgId = inComingMsg.Header.MsgId
	slog.Debug("sent get resp message", "msg", msg)
	err := c.HandleMTPMsgTransmit(msg)
	if err != nil {
		slog.Warn("HandleGetRequest ", "error", err)
	}
}

func (c *WSClient) createGetResponseMessage(getNodePaths []string) (result *api.Msg) {
	resp := c.clientRepository.ConstructGetResp(getNodePaths)
	result = &api.Msg{
		Header: &api.Header{
			MsgType: api.Header_GET_RESP,
		},
		Body: &api.Body{
			MsgBody: &api.Body_Response{
				Response: &api.Response{
					RespType: &resp,
				},
			},
		},
	}
	return
}

func (c *WSClient) HandleSetRequest(inComingMsg *api.Msg) {
	getUpdateObjs := inComingMsg.GetBody().GetRequest().GetSet().GetUpdateObjs()

	var affectedPath []string
	var requestPath []string
	var updatedParams []map[string]string
	paramSettings := make(map[string]string)

	slog.Debug("get set req from tauc", "req", getUpdateObjs)

	for _, UpdateObj := range getUpdateObjs {
		path := UpdateObj.GetObjPath()
		isSuccess, nodePath := c.clientRepository.IsExistPath(path)
		if !isSuccess {
			continue
		}
		requestPath = append(requestPath, path)
		affectedPath = append(affectedPath, tpath)
		for _, paramSetting := range UpdateObj.GetParamSettings() {
			setKey := paramSetting.GetParam()
			setValue := paramSetting.GetValue()
			paramSettings[setKey] = setValue
			c.clientRepository.HandleSetRequest(nodePath, setKey, setValue)
		}
		updatedParams = append(updatedParams, paramSettings)
	}

	//c.repo.SaveData(dataMap)

	msg := c.createSetResponseMessage(requestPath, affectedPath, updatedParams)
	msg.Header.MsgId = inComingMsg.Header.MsgId
	slog.Debug("sent set resp message")
	err := c.HandleMTPMsgTransmit(msg)
	if err != nil {
		slog.Warn("HandleSetRequest", "error", err)
	}

}

func (c *WSClient) createSetResponseMessage(requestPath []string, affectedPath []string, updatedParams []map[string]string) (result *api.Msg) {
	var updatedObjResults []*api.SetResp_UpdatedObjectResult
	for k, path := range requestPath {
		updatedObjResult := &api.SetResp_UpdatedObjectResult{
			RequestedPath: path,
			OperStatus: &api.SetResp_UpdatedObjectResult_OperationStatus{
				OperStatus: &api.SetResp_UpdatedObjectResult_OperationStatus_OperSuccess{
					OperSuccess: &api.SetResp_UpdatedObjectResult_OperationStatus_OperationSuccess{
						UpdatedInstResults: []*api.SetResp_UpdatedInstanceResult{
							{
								AffectedPath:  affectedPath[k],
								UpdatedParams: updatedParams[k],
							},
						},
					},
				},
			},
		}
		updatedObjResults = append(updatedObjResults, updatedObjResult)
	}
	result = &api.Msg{
		Header: &api.Header{
			MsgType: api.Header_SET_RESP,
		},
		Body: &api.Body{
			MsgBody: &api.Body_Response{
				Response: &api.Response{
					RespType: &api.Response_SetResp{
						SetResp: &api.SetResp{
							UpdatedObjResults: updatedObjResults,
						},
					},
				},
			},
		},
	}
	return
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
	msg := c.createAddResponseMessage(requestPath, affectedPath, updatedParams)
	msg.Header.MsgId = inComingMsg.Header.MsgId
	slog.Debug("sent add resp message")
	err := c.HandleMTPMsgTransmit(msg)
	if err != nil {
		slog.Warn("HandleAddRequest error:", "error", err)
	}
}

func (c *WSClient) createAddResponseMessage(requestPath []string, affectedPath []string, updatedParams []map[string]string) (result *api.Msg) {
	var createdObjResults []*api.AddResp_CreatedObjectResult
	for k, path := range requestPath {
		createdObjResult := &api.AddResp_CreatedObjectResult{
			RequestedPath: path,
			OperStatus: &api.AddResp_CreatedObjectResult_OperationStatus{
				OperStatus: &api.AddResp_CreatedObjectResult_OperationStatus_OperSuccess{
					OperSuccess: &api.AddResp_CreatedObjectResult_OperationStatus_OperationSuccess{
						InstantiatedPath: affectedPath[k],
						UniqueKeys:       updatedParams[k],
					},
				},
			},
		}
		createdObjResults = append(createdObjResults, createdObjResult)
	}

	result = &api.Msg{
		Header: &api.Header{
			MsgType: api.Header_ADD_RESP,
		},
		Body: &api.Body{
			MsgBody: &api.Body_Response{
				Response: &api.Response{
					RespType: &api.Response_AddResp{
						AddResp: &api.AddResp{
							CreatedObjResults: createdObjResults,
						},
					},
				},
			},
		},
	}
	return
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
	msg := c.createDeleteResponseMessage(requestPath, affectedPath)
	msg.Header.MsgId = inComingMsg.Header.MsgId
	slog.Debug("sent delete resp message")
	err := c.HandleMTPMsgTransmit(msg)
	if err != nil {
		slog.Warn("HandleDeleteRequest error:", "error", err)
	}
}

func (c *WSClient) createDeleteResponseMessage(requestPath []string, affectedPath []string) (result *api.Msg) {
	var deletedObjResults []*api.DeleteResp_DeletedObjectResult
	for k, path := range requestPath {
		deletedObjResult := &api.DeleteResp_DeletedObjectResult{
			RequestedPath: path,
			OperStatus: &api.DeleteResp_DeletedObjectResult_OperationStatus{
				OperStatus: &api.DeleteResp_DeletedObjectResult_OperationStatus_OperSuccess{
					OperSuccess: &api.DeleteResp_DeletedObjectResult_OperationStatus_OperationSuccess{
						AffectedPaths: []string{affectedPath[k]},
					},
				},
			},
		}
		deletedObjResults = append(deletedObjResults, deletedObjResult)
	}

	result = &api.Msg{
		Header: &api.Header{
			MsgType: api.Header_DELETE_RESP,
		},
		Body: &api.Body{
			MsgBody: &api.Body_Response{
				Response: &api.Response{
					RespType: &api.Response_DeleteResp{
						DeleteResp: &api.DeleteResp{
							DeletedObjResults: deletedObjResults,
						},
					},
				},
			},
		},
	}
	return
}

func (c *WSClient) HandleOperateRequest(inComingMsg *api.Msg) {
	operate := inComingMsg.GetBody().GetRequest().GetOperate()

	slog.Debug("get operate req", "req", operate)

	command := operate.GetCommand()

	msg := c.createOperateResponseMessage(command)
	msg.Header.MsgId = inComingMsg.Header.MsgId
	slog.Debug("sent operate resp message")
	err := c.HandleMTPMsgTransmit(msg)
	if err != nil {
		slog.Warn("HandleOperateRequest error:", "error", err)
	}
}

func (c *WSClient) createOperateResponseMessage(command string) (result *api.Msg) {
	var operationResults []*api.OperateResp_OperationResult

	operationResult := &api.OperateResp_OperationResult{
		ExecutedCommand: command,
		OperationResp: &api.OperateResp_OperationResult_ReqObjPath{
			ReqObjPath: "Device.LocalAgent.Request.1",
		},
	}
	operationResults = append(operationResults, operationResult)

	result = &api.Msg{
		Header: &api.Header{
			MsgType: api.Header_OPERATE_RESP,
		},
		Body: &api.Body{
			MsgBody: &api.Body_Response{
				Response: &api.Response{
					RespType: &api.Response_OperateResp{
						OperateResp: &api.OperateResp{
							OperationResults: operationResults,
						},
					},
				},
			},
		},
	}
	return
}

func (c *WSClient) SendOperateCompleteNotify(objPath string, commandName string, commandKey string, outputArgs map[string]string) {

	completeNotify := &api.Notify{
		SubscriptionId: "/tpuc/tr369controller",
		SendResp:       false,
		Notification: &api.Notify_OperComplete{
			OperComplete: &api.Notify_OperationComplete{
				ObjPath:     objPath,
				CommandName: commandName,
				CommandKey:  commandKey,
				OperationResp: &api.Notify_OperationComplete_ReqOutputArgs{
					ReqOutputArgs: &api.Notify_OperationComplete_OutputArgs{
						OutputArgs: outputArgs,
					},
				},
			},
		},
	}
	msg := c.createNotifyMessage(completeNotify)
	err := c.HandleMTPMsgTransmit(msg)
	if err != nil {
		slog.Warn("SendOperateCompleteNotify ", "error", err)
	}
}

func (c *WSClient) createNotifyMessage(notify *api.Notify) (result *api.Msg) {
	result = &api.Msg{
		Header: &api.Header{
			MsgType: api.Header_NOTIFY,
			MsgId:   common.RandStr(10),
		},
		Body: &api.Body{
			MsgBody: &api.Body_Request{
				Request: &api.Request{
					ReqType: &api.Request_Notify{
						Notify: notify,
					},
				},
			},
		},
	}
	return
}

func (c *WSClient) HandleMTPMsgTransmit(msg *api.Msg) error {

	rec := c.createUspRecordNoSession("1.0", c.config.EndpointId, c.config.ControllerIdentifier, msg)
	payload := c.encodeUspRecord(rec)
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

func (c *WSClient) encodeUspRecord(rec *api.Record) []byte {
	result, _ := proto.Marshal(rec)
	return result
}

func (c *WSClient) encodeUspMessage(msg *api.Msg) []byte {
	result, _ := proto.Marshal(msg)
	return result
}

func (c *WSClient) createUspRecordNoSession(ver, to, from string, um *api.Msg) (result *api.Record) {
	result = &api.Record{
		Version: ver,
		ToId:    to,
		FromId:  from,
		RecordType: &api.Record_NoSessionContext{
			NoSessionContext: &api.NoSessionContextRecord{
				Payload: c.encodeUspMessage(um),
			},
		},
	}

	return
}
