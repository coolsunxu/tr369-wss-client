// Package main 是 TR369 WebSocket 客户端的入口点
package main

import (
	"os"
	"os/signal"
	"syscall"

	"tr369-wss-client/internal/domain/services"
	"tr369-wss-client/internal/infrastructure/config"
	"tr369-wss-client/internal/infrastructure/di"
	"tr369-wss-client/internal/infrastructure/logging"

	"github.com/spf13/cobra"
)

var (
	logger    services.Logger
	container *di.Container
)

var rootCmd = &cobra.Command{
	Use:   "tr369-client",
	Short: "TR369 WebSocket 客户端",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 初始化日志
		logger = logging.NewZapLogger()

		// 初始化配置
		err := config.InitConfig("./configs/environments/development.json", logger)
		if err != nil {
			logger.Fatal("初始化配置失败: %v", err)
			return
		}

		// 创建依赖注入容器
		container = di.NewContainer(&config.GlobalConfig)

		logger.Info("配置加载成功")
	},
	Run: func(cmd *cobra.Command, args []string) {
		startClient()
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func startClient() {
	defer container.Shutdown()

	// 通过容器获取依赖
	wsClient := container.GetWebSocketClient()

	// 初始化仓储
	_ = container.GetClientRepository()

	// 连接服务器
	cfg := container.GetConfig()
	logger.Info("正在连接 TR369 服务器: %s", cfg.WebsocketConfig.ServerURL)
	if err := wsClient.Connect(); err != nil {
		logger.Fatal("连接失败: %v", err)
		return
	}

	logger.Info("成功连接到 TR369 服务器")

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭...")
}
