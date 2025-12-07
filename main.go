package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"tr369-wss-client/client/repository"
	"tr369-wss-client/client/usecase"
	logger "tr369-wss-client/log"
	"tr369-wss-client/utils"

	"tr369-wss-client/client"
	"tr369-wss-client/config"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "websocket client",
	Short: "wsClient",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 初始化日志
		logger.InitLogger()

		//  初始化配置
		err := config.InitConfig("./config/config.json")
		if err != nil {
			logger.Fatal(err)
			return
		}

		logger.Infof("configs.InitConfig %v", utils.SafeMarshal(config.GlobalConfig))
	},
	Run: func(cmd *cobra.Command, args []string) {
		startClient()
	},
}

func init() {
	// 暂时不做操作
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		return
	}
}

func startClient() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	messageChannel := make(chan []byte, config.GlobalConfig.WebsocketConfig.MessageChannelSize)
	defer close(messageChannel)

	// 初始化数据操作
	clientRepository := repository.NewClientRepository(&config.GlobalConfig, ctx, cancel)
	clientRepository.StartClientRepository()

	// 初始化clientUseCase
	clientUseCase := usecase.NewClientUseCase(ctx, &config.GlobalConfig, clientRepository, messageChannel)

	// 创建WebSocket客户端
	wsClient := client.NewWSClient(&config.GlobalConfig, clientRepository, clientUseCase, messageChannel)

	// 连接到服务器
	logger.Infof("Connecting to TR369 server at %s...", config.GlobalConfig.WebsocketConfig.ServerURL)
	if err := wsClient.Connect(); err != nil {
		logger.Fatalf("Failed to connect: %v", err)
	}
	defer wsClient.Disconnect()

	logger.Infof("Connected to TR369 server successfully")

	// 启动消息处理
	wsClient.StartMessageHandler()

	// 等待中断信号优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Infof("Shutting down...")
}
