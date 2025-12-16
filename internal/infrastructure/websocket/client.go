// Package websocket 提供 WebSocket 客户端实现
package websocket

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/coder/websocket"

	"tr369-wss-client/internal/domain/services"
	"tr369-wss-client/internal/infrastructure/config"
)

// Client 表示 WebSocket 客户端
type Client struct {
	config         *config.Config
	conn           *websocket.Conn
	ctx            context.Context
	cancel         context.CancelFunc
	connected      bool
	pingTicker     *time.Ticker
	messageChannel chan []byte
	logger         services.Logger
}

// NewClient 创建新的 WebSocket 客户端实例
func NewClient(cfg *config.Config, messageChannel chan []byte, logger services.Logger) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	return &Client{
		config:         cfg,
		ctx:            ctx,
		cancel:         cancel,
		connected:      false,
		messageChannel: messageChannel,
		logger:         logger,
	}
}

// 确保 Client 实现了 services.WebSocketClient 接口
var _ services.WebSocketClient = (*Client)(nil)

// Connect 建立 WebSocket 连接
func (c *Client) Connect() error {
	options := &websocket.DialOptions{
		Subprotocols: []string{"v1.usp"},
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	headers := http.Header{}
	options.HTTPHeader = headers

	connectUrl := c.config.WebsocketConfig.ServerURL

	// 检查 URL 是否已经包含查询参数
	if strings.Contains(connectUrl, "?") {
		connectUrl += "&eid=" + c.config.WebsocketConfig.EndpointId
	} else {
		connectUrl += "?eid=" + c.config.WebsocketConfig.EndpointId
	}

	c.logger.Info("正在连接: %s", connectUrl)

	// 连接服务器
	conn, _, err := websocket.Dial(c.ctx, connectUrl, options)
	if err != nil {
		return fmt.Errorf("连接失败: %w", err)
	}

	c.conn = conn
	c.connected = true

	// 设置读消息的最大大小
	c.conn.SetReadLimit(c.config.WebsocketConfig.MaxMessageSize)

	// 配置 ping/pong
	c.pingTicker = time.NewTicker(time.Duration(c.config.WebsocketConfig.PingInterval))

	return nil
}

// Disconnect 关闭 WebSocket 连接
func (c *Client) Disconnect() {
	if c.pingTicker != nil {
		c.pingTicker.Stop()
		c.pingTicker = nil
	}

	if c.conn != nil {
		err := c.conn.Close(websocket.StatusNormalClosure, "客户端断开连接")
		if err != nil {
			c.logger.Error("关闭 WebSocket 连接失败: %v", err)
		}
		c.conn = nil
	}

	c.connected = false
}

// IsConnected 返回连接状态
func (c *Client) IsConnected() bool {
	return c.connected
}

// Send 发送消息
func (c *Client) Send(data []byte) error {
	if c.conn == nil {
		return fmt.Errorf("连接未建立")
	}

	return c.conn.Write(c.ctx, websocket.MessageBinary, data)
}

// Read 读取消息
func (c *Client) Read() ([]byte, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("连接未建立")
	}

	_, data, err := c.conn.Read(c.ctx)
	return data, err
}

// Context 返回客户端上下文
func (c *Client) Context() context.Context {
	return c.ctx
}

// Cancel 取消客户端上下文
func (c *Client) Cancel() {
	c.cancel()
}
