// Package main 提供基本客户端使用示例
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"tr369-wss-client/internal/infrastructure/config"
	"tr369-wss-client/internal/infrastructure/logging"
	"tr369-wss-client/internal/infrastructure/websocket"
)

func main() {
	// 初始化日志
	logger := logging.NewLogger()
	logger.Info("启动示例客户端...")

	// 加载配置
	cfg, err := config.LoadConfig("./configs/environments/development.json")
	if err != nil {
		logger.Fatal("加载配置失败: %v", err)
		return
	}

	// 创建消息通道
	messageChannel := make(chan []byte, cfg.WebsocketConfig.MessageChannelSize)
	defer close(messageChannel)

	// 创建 WebSocket 客户端
	client := websocket.NewClient(cfg, messageChannel, logger)

	// 连接服务器
	if err := client.Connect(); err != nil {
		logger.Fatal("连接失败: %v", err)
		return
	}
	defer client.Disconnect()

	logger.Info("连接成功!")

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动消息处理
	go handleMessages(ctx, client, logger)

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭...")
}

func handleMessages(ctx context.Context, client *websocket.Client, logger logging.Logger) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			data, err := client.Read()
			if err != nil {
				logger.Error("读取消息失败: %v", err)
				return
			}
			fmt.Printf("收到消息: %d 字节\n", len(data))
		}
	}
}
