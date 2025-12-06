package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tr369-wss-client/client/repository"
	logger "tr369-wss-client/log"
	"tr369-wss-client/utils"

	"tr369-wss-client/client"
	"tr369-wss-client/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "websocket client",
	Short: "wsClient",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 初始化日志
		logger.InitLogger()

		//  初始化配置
		err := config.InitConfig("./configs/config.json")
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
	// 在命令的 Flag 中定义参数
	rootCmd.Flags().StringP("websocket-server-port", "p", "8081", "websocket server port")
	rootCmd.Flags().BoolP("use-tls", "t", false, "use tls or not")
	rootCmd.Flags().StringP("controllerId", "c", "usp-controller-ws", "controllerId")
	rootCmd.Flags().StringP("controllerPath", "a", "/usp", "controllerPath")

	// Viper 绑定到具体的命令 Flag
	err := viper.BindPFlag("wsPort", rootCmd.Flags().Lookup("websocket-server-port"))
	if err != nil {
		return
	}
	err = viper.BindPFlag("wsTlsConfig.useTls", rootCmd.Flags().Lookup("use-tls"))
	if err != nil {
		return
	}
	err = viper.BindPFlag("controllerId", rootCmd.Flags().Lookup("controllerId"))
	if err != nil {
		return
	}

	err = viper.BindPFlag("controllerPath", rootCmd.Flags().Lookup("controllerPath"))
	if err != nil {
		return
	}
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		return
	}
}

func startClient() {

	ctx, cancel := context.WithCancel(context.Background())

	// 初始化数据操作
	clientRepository := repository.NewClientRepository(config.GlobalConfig, ctx, cancel)
	clientRepository.StartClientRepository()

	// 创建WebSocket客户端
	wsClient := client.NewWSClient(config.GlobalConfig, clientRepository)

	// 连接到服务器
	log.Printf("Connecting to TR369 server at %s...", config.GlobalConfig.ServerURL)
	if err := wsClient.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer wsClient.Disconnect()

	log.Println("Connected to TR369 server successfully")

	// 启动消息处理
	wsClient.StartMessageHandler()

	log.Println("TR369 client is running. Press Ctrl+C to exit.")

	// 等待中断信号优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
}
