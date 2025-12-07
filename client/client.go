package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
	"tr369-wss-client/client/model"
	logger "tr369-wss-client/log"
	"tr369-wss-client/utils"

	"github.com/coder/websocket"

	"tr369-wss-client/config"
)

// WSClient represents the TR369 WebSocket client
type WSClient struct {
	config               *config.Config
	conn                 *websocket.Conn
	ctx                  context.Context
	cancel               context.CancelFunc
	connected            bool
	pingTicker           *time.Ticker
	clientRepository     model.ClientRepository
	clientUseCase        model.ClientUseCase
	messageChannel       chan []byte   // 消息发送通道
	reconnectInterval    time.Duration // 重连间隔
	maxReconnectAttempts int           // 最大重连次数
	reconnectAttempts    int           // 当前重连次数
	reconnectTicker      *time.Ticker  // 重连定时器
}

// NewWSClient creates a new WebSocket client instance
func NewWSClient(
	cfg *config.Config,
	clientRepository model.ClientRepository,
	clientUseCase model.ClientUseCase,
	messageChannel chan []byte,
) *WSClient {
	ctx, cancel := context.WithCancel(context.Background())

	return &WSClient{
		config:           cfg,
		ctx:              ctx,
		cancel:           cancel,
		connected:        false,
		clientRepository: clientRepository,
		clientUseCase:    clientUseCase,
		messageChannel:   messageChannel,
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

	connectUrl := c.config.WebsocketConfig.ServerURL

	// 检查URL是否已经包含查询参数
	if strings.Contains(connectUrl, "?") {
		connectUrl += "&eid=" + c.config.WebsocketConfig.EndpointId
	} else {
		connectUrl += "?eid=" + c.config.WebsocketConfig.EndpointId
	}

	logger.Infof("Connecting with eid in url: %s", connectUrl)

	// 连接服务器
	conn, _, err := websocket.Dial(c.ctx, connectUrl, options)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}

	c.conn = conn
	c.connected = true

	// 设置读消息的最大大小
	c.conn.SetReadLimit(c.config.WebsocketConfig.MaxMessageSize)

	// 配置ping/pong
	c.pingTicker = time.NewTicker(time.Duration(c.config.WebsocketConfig.PingInterval))

	return nil
}

// Disconnect closes the WebSocket connection
func (c *WSClient) Disconnect() {
	if c.pingTicker != nil {
		c.pingTicker.Stop()
		c.pingTicker = nil
	}

	if c.conn != nil {
		// 先发送关闭帧
		err := c.conn.Close(websocket.StatusNormalClosure, "client disconnecting")
		if err != nil {
			logger.Errorf("Failed to close WebSocket connection: %v", err)
		}
		c.conn = nil
	}

	if c.reconnectTicker != nil {
		c.reconnectTicker.Stop()
		c.reconnectTicker = nil
	}

	c.connected = false

	// 启动重连逻辑，除非已经达到最大重连次数
	if c.reconnectAttempts < c.maxReconnectAttempts {
		c.StartReconnectTicker()
	}
}

// StartMessageHandler starts the message handling goroutines
func (c *WSClient) StartMessageHandler() {
	// 启动ping goroutine
	go c.pingHandler()

	// 启动消息接收goroutine
	go c.messageHandler()

	// 启动消息写goroutine，顺序写入，便于控制
	go c.messageSendHandler()
}

// pingHandler handles periodic ping messages
func (c *WSClient) pingHandler() {
	for {
		select {
		case <-c.pingTicker.C:
			if c.conn == nil {
				return
			}
			if err := c.conn.Ping(c.ctx); err != nil {
				logger.Infof("Ping failed: %v", err)
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
			if c.conn == nil {
				return
			}
			// 读取二进制消息
			_, data, err := c.conn.Read(c.ctx)
			if err != nil {
				logger.Infof("Connection closed with status %v", err)
				c.Disconnect()
				return
			}

			record, err := utils.DecodeUSPRecord(data)
			if err != nil {
				logger.Infof("Failed to decode USP Record: %v", err)
				continue
			}

			logger.Infof("Decoded Record - From: %s, To: %s", record.FromId, record.ToId)

			// 提取NoSessionContextRecord
			noSessionContext := record.GetNoSessionContext()
			if noSessionContext == nil {
				logger.Infof("Record is not NoSessionContextRecord")
				continue
			}

			// 从NoSessionContextRecord中提取payload并解析为api.Msg
			msg, err := utils.DecodeUSPMessage(noSessionContext.GetPayload())
			if err != nil {
				logger.Infof("Failed to decode USP Message: %v", err)
				continue
			}

			// 处理接收到的消息，调用usecase层的HandleMessage方法
			c.clientUseCase.HandleMessage(msg)
		}
	}
}

// messageSendHandler handles sending messages from the message channel
func (c *WSClient) messageSendHandler() {
	for {

		select {
		case <-c.ctx.Done():
			return
		case payload, ok := <-c.messageChannel:
			if !ok {
				return
			}

			// 发送二进制消息
			if c.conn == nil {
				logger.Warnf("Cannot send message: connection is nil")
				return
			}

			if err := c.conn.Write(c.ctx, websocket.MessageBinary, payload); err != nil {
				logger.Infof("Failed to send response: %v", err)
				return
			}

			logger.Infof("Send message success")
		}
	}
}

// StartReconnectTicker starts the reconnect ticker
func (c *WSClient) StartReconnectTicker() {
	if c.reconnectTicker != nil {
		c.reconnectTicker.Stop()
	}

	c.reconnectTicker = time.NewTicker(c.reconnectInterval)

	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return
			case <-c.reconnectTicker.C:
				c.Reconnect()
			}
		}
	}()
}

// Reconnect attempts to reconnect to the server
func (c *WSClient) Reconnect() {
	c.reconnectAttempts++
	logger.Infof("Attempting to reconnect (attempt %d/%d)...", c.reconnectAttempts, c.maxReconnectAttempts)

	// 尝试连接服务器
	err := c.Connect()
	if err != nil {
		logger.Errorf("Reconnect failed: %v", err)

		// 如果达到最大重连次数，停止重连
		if c.reconnectAttempts >= c.maxReconnectAttempts {
			logger.Errorf("Max reconnect attempts reached, giving up")
			if c.reconnectTicker != nil {
				c.reconnectTicker.Stop()
				c.reconnectTicker = nil
			}
			c.cancel()
		}
		return
	}

	// 重连成功，重置重连计数器
	logger.Infof("Reconnected successfully")
	c.reconnectAttempts = 0

	// 停止重连定时器
	if c.reconnectTicker != nil {
		c.reconnectTicker.Stop()
		c.reconnectTicker = nil
	}

	// 启动消息处理
	c.StartMessageHandler()
}
