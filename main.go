package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"tr369-wss-client/client"
	"tr369-wss-client/config"
	"tr369-wss-client/datamodel"
)

func main() {
	// 解析命令行参数
	configFile := flag.String("config", "", "Path to configuration file")
	debug := flag.Bool("debug", true, "Enable debug mode")
	flag.Parse()

	// 设置日志格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 加载配置
	var cfg *config.Config
	if *configFile != "" {
		// 从指定文件加载配置
		cfg = loadConfigFromFile(*configFile)
	} else {
		// 使用默认配置路径
		cfg = config.LoadConfig()
	}

	// 打印配置信息
	if *debug {
		log.Printf("Configuration: %v", cfg)
	}

	// 初始化数据模型
	model := datamodel.NewTR181DataModel()

	// 创建WebSocket客户端
	wsClient := client.NewWSClient(cfg, model)

	// 连接到服务器
	log.Printf("Connecting to TR369 server at %s...", cfg.ServerURL)
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

// loadConfigFromFile 从指定文件加载配置
func loadConfigFromFile(filePath string) *config.Config {
	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	// 使用默认配置作为基础
	cfg := config.DefaultConfig()

	// 解析JSON配置
	if err := json.Unmarshal(content, cfg); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	log.Printf("Configuration loaded from %s", filePath)
	return cfg
}
