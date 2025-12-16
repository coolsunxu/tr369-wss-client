// Package client 提供客户端相关的用例实现
package client

import (
	"context"

	"tr369-wss-client/internal/domain/services"
)

// ConnectUseCase 连接用例
type ConnectUseCase struct {
	config services.ConfigProvider
	logger services.Logger
	client services.WebSocketClient
}

// NewConnectUseCase 创建新的连接用例
func NewConnectUseCase(cfg services.ConfigProvider, logger services.Logger, client services.WebSocketClient) *ConnectUseCase {
	return &ConnectUseCase{
		config: cfg,
		logger: logger,
		client: client,
	}
}

// Execute 执行连接操作
func (uc *ConnectUseCase) Execute(ctx context.Context) error {
	if err := uc.client.Connect(); err != nil {
		uc.logger.Error("连接失败: %v", err)
		return err
	}

	uc.logger.Info("连接成功")
	return nil
}

// DisconnectUseCase 断开连接用例
type DisconnectUseCase struct {
	logger services.Logger
	client services.WebSocketClient
}

// NewDisconnectUseCase 创建新的断开连接用例
func NewDisconnectUseCase(logger services.Logger, client services.WebSocketClient) *DisconnectUseCase {
	return &DisconnectUseCase{
		logger: logger,
		client: client,
	}
}

// Execute 执行断开连接操作
func (uc *DisconnectUseCase) Execute() {
	if uc.client != nil {
		uc.client.Disconnect()
		uc.logger.Info("已断开连接")
	}
}
