// Package di 提供依赖注入容器
package di

import (
	"context"

	"tr369-wss-client/internal/domain/repositories"
	"tr369-wss-client/internal/domain/services"
	"tr369-wss-client/internal/infrastructure/config"
	"tr369-wss-client/internal/infrastructure/logging"
	"tr369-wss-client/internal/infrastructure/persistence/repository"
	"tr369-wss-client/internal/infrastructure/websocket"
)

// Container 依赖注入容器
type Container struct {
	config           *config.Config
	logger           services.Logger
	clientRepository repositories.ClientRepository
	wsClient         services.WebSocketClient
	messageChannel   chan []byte
	ctx              context.Context
	cancel           context.CancelFunc
}

// NewContainer 创建新的依赖注入容器
func NewContainer(cfg *config.Config) *Container {
	ctx, cancel := context.WithCancel(context.Background())

	return &Container{
		config: cfg,
		ctx:    ctx,
		cancel: cancel,
	}
}

// InitializeLogger 初始化日志服务
func (c *Container) InitializeLogger() services.Logger {
	if c.logger == nil {
		c.logger = logging.NewZapLogger()
	}
	return c.logger
}

// InitializeClientRepository 初始化客户端仓储
func (c *Container) InitializeClientRepository() repositories.ClientRepository {
	if c.clientRepository == nil {
		logger := c.InitializeLogger()
		repo := repository.NewClientRepository(c.config, c.ctx, c.cancel, logger)
		repo.StartClientRepository()
		c.clientRepository = repo
	}
	return c.clientRepository
}

// InitializeWebSocketClient 初始化 WebSocket 客户端
func (c *Container) InitializeWebSocketClient() services.WebSocketClient {
	if c.wsClient == nil {
		logger := c.InitializeLogger()
		c.messageChannel = make(chan []byte, c.config.WebsocketConfig.MessageChannelSize)
		c.wsClient = websocket.NewClient(c.config, c.messageChannel, logger)
	}
	return c.wsClient
}

// GetLogger 获取日志服务
func (c *Container) GetLogger() services.Logger {
	return c.InitializeLogger()
}

// GetClientRepository 获取客户端仓储
func (c *Container) GetClientRepository() repositories.ClientRepository {
	return c.InitializeClientRepository()
}

// GetWebSocketClient 获取 WebSocket 客户端
func (c *Container) GetWebSocketClient() services.WebSocketClient {
	return c.InitializeWebSocketClient()
}

// GetConfig 获取配置
func (c *Container) GetConfig() *config.Config {
	return c.config
}

// GetMessageChannel 获取消息通道
func (c *Container) GetMessageChannel() chan []byte {
	if c.messageChannel == nil {
		c.messageChannel = make(chan []byte, c.config.WebsocketConfig.MessageChannelSize)
	}
	return c.messageChannel
}

// Shutdown 关闭容器，释放资源
func (c *Container) Shutdown() {
	if c.wsClient != nil {
		c.wsClient.Disconnect()
	}
	if c.messageChannel != nil {
		close(c.messageChannel)
	}
	c.cancel()
}
