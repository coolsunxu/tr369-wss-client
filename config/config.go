package config

import (
	"time"
)

// DataRefreshConfig 定义数据刷库的参数
type DataRefreshConfig struct {
	// 刷新间隔时间，单位为秒
	IntervalSeconds int `json:"interval_seconds"`

	// 写入次数
	WriteCountThreshold int `json:"write_count_threshold"`
}

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

	DataRefreshConfig *DataRefreshConfig `json:"data_refresh_config"`
}

// DefaultConfig returns the default configuration
var GlobalConfig = &Config{
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
