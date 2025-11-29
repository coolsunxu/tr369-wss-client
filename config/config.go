package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Config represents the client configuration
type Config struct {
	// 服务器配置
	ServerURL string `json:"server_url"`

	// WebSocket配置
	PingInterval         time.Duration `json:"ping_interval"`
	ReconnectDelay       time.Duration `json:"reconnect_delay"`
	MaxReconnectAttempts int           `json:"max_reconnect_attempts"`
	MaxMessageSize       int64         `json:"max_message_size"`

	// TLS配置
	TLSEnabled         bool   `json:"tls_enabled"`
	InsecureSkipVerify bool   `json:"insecure_skip_verify"`
	CACertFile         string `json:"ca_cert_file"`
	ClientCertFile     string `json:"client_cert_file"`
	ClientKeyFile      string `json:"client_key_file"`

	// 日志配置
	LogLevel string `json:"log_level"`
	LogFile  string `json:"log_file"`

	// 设备信息
	EndpointId string `json:"endpoint_id"`

	// controller Id
	ControllerIdentifier string `json:"controller_identifier"`

	// tr181节点
	TR181DataModelPath string `json:"tr181_data_model_path"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		// 默认服务器URL
		ServerURL: "wss://localhost:7547",

		// 默认WebSocket设置
		PingInterval:         30 * time.Second,
		ReconnectDelay:       5 * time.Second,
		MaxReconnectAttempts: 5,
		MaxMessageSize:       1024 * 1024, // 1MB

		// 默认TLS设置
		TLSEnabled:         true,
		InsecureSkipVerify: false,
		CACertFile:         "",
		ClientCertFile:     "",
		ClientKeyFile:      "",

		// 默认日志设置
		LogLevel: "info",
		LogFile:  "",

		ControllerIdentifier: "usp-controller-ws",
	}
}

// LoadConfig loads the configuration from file or returns default
func LoadConfig() *Config {
	// 尝试从多个位置加载配置文件
	configPaths := []string{
		"./config.json",
	}

	var config *Config
	var err error

	// 首先使用默认配置
	config = DefaultConfig()

	// 尝试从配置文件加载
	for _, path := range configPaths {

		// 检查文件是否存在
		if _, err = os.Stat(path); os.IsNotExist(err) {
			log.Printf("Config file %s does not exist, skipping\n", path)
			continue
		}

		// 加载配置文件
		config, err = loadConfigFromFile(path)
		if err == nil {
			// 成功加载配置
			fmt.Printf("Configuration loaded from %s\n", path)
			return config
		}
		fmt.Printf("Error loading config from %s: %v\n", path, err)
	}

	// 如果没有找到配置文件，使用默认配置
	fmt.Println("No configuration file found, using default settings")
	return config
}

// loadConfigFromFile loads configuration from a JSON file
func loadConfigFromFile(filePath string) (*Config, error) {
	// 读取文件内容
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 使用默认配置作为基础
	config := DefaultConfig()

	// 解析JSON配置
	err = json.Unmarshal(content, config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

// SaveConfig saves the configuration to a file
func (c *Config) SaveConfig(filePath string) error {
	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// 序列化配置为JSON
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// 写入文件
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// String returns a string representation of the configuration
func (c *Config) String() string {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Sprintf("{error: %v}", err)
	}
	return string(data)
}
