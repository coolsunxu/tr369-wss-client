// Package config 提供配置管理功能
package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// DataRefreshConfig 定义数据刷新配置
type DataRefreshConfig struct {
	// IntervalSeconds 刷新间隔时间，单位为秒
	IntervalSeconds int `json:"interval_seconds" mapstructure:"interval_seconds"`

	// WriteCountThreshold 写入次数阈值
	WriteCountThreshold int `json:"write_count_threshold" mapstructure:"write_count_threshold"`

	// TR181DataModelPath TR181 数据模型文件路径
	TR181DataModelPath string `json:"tr181_data_model_path" mapstructure:"tr181_data_model_path"`
}

// WebsocketConfig 定义 WebSocket 配置
type WebsocketConfig struct {
	// ServerURL 服务器地址
	ServerURL string `json:"server_url" mapstructure:"server_url"`

	// ControllerId 控制器 ID
	ControllerId string `json:"controller_id" mapstructure:"controller_id"`

	// PingInterval ping 间隔
	PingInterval int `json:"ping_interval" mapstructure:"ping_interval"`

	// MaxMessageSize 最大消息大小
	MaxMessageSize int64 `json:"max_message_size" mapstructure:"max_message_size"`

	// EndpointId 端点 ID
	EndpointId string `json:"endpoint_id" mapstructure:"endpoint_id"`

	// MessageChannelSize 消息通道容量
	MessageChannelSize int `json:"message_channel_size" mapstructure:"message_channel_size"`
}

// Config 表示客户端配置
type Config struct {
	DataRefreshConfig *DataRefreshConfig `json:"data_refresh_config" mapstructure:"data_refresh_config"`
	WebsocketConfig   *WebsocketConfig   `json:"websocket_config" mapstructure:"websocket_config"`
}

// LoadConfig 从文件加载配置
func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("打开配置文件失败: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return &config, nil
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.DataRefreshConfig == nil {
		return fmt.Errorf("DataRefreshConfig 不能为空")
	}

	if c.DataRefreshConfig.IntervalSeconds <= 0 {
		return fmt.Errorf("IntervalSeconds 必须为正数")
	}

	if c.DataRefreshConfig.WriteCountThreshold <= 0 {
		return fmt.Errorf("WriteCountThreshold 必须为正数")
	}

	if c.DataRefreshConfig.TR181DataModelPath == "" {
		return fmt.Errorf("TR181DataModelPath 不能为空")
	}

	if c.WebsocketConfig == nil {
		return fmt.Errorf("WebsocketConfig 不能为空")
	}

	if c.WebsocketConfig.ServerURL == "" {
		return fmt.Errorf("ServerURL 不能为空")
	}

	if c.WebsocketConfig.EndpointId == "" {
		return fmt.Errorf("EndpointId 不能为空")
	}

	if c.WebsocketConfig.ControllerId == "" {
		return fmt.Errorf("ControllerId 不能为空")
	}

	if c.WebsocketConfig.PingInterval <= 0 {
		return fmt.Errorf("PingInterval 必须为正数")
	}

	if c.WebsocketConfig.MaxMessageSize <= 0 {
		return fmt.Errorf("MaxMessageSize 必须为正数")
	}

	if c.WebsocketConfig.MessageChannelSize < 0 {
		return fmt.Errorf("MessageChannelSize 不能为负数")
	}

	return nil
}

// GetServerURL 实现 ConfigProvider 接口
func (c *Config) GetServerURL() string {
	if c.WebsocketConfig == nil {
		return ""
	}
	return c.WebsocketConfig.ServerURL
}

// GetEndpointID 实现 ConfigProvider 接口
func (c *Config) GetEndpointID() string {
	if c.WebsocketConfig == nil {
		return ""
	}
	return c.WebsocketConfig.EndpointId
}

// GetControllerEndpointID 实现 ConfigProvider 接口
func (c *Config) GetControllerEndpointID() string {
	if c.WebsocketConfig == nil {
		return ""
	}
	return c.WebsocketConfig.ControllerId
}

// GetMessageChannelSize 实现 ConfigProvider 接口
func (c *Config) GetMessageChannelSize() int {
	if c.WebsocketConfig == nil {
		return 0
	}
	return c.WebsocketConfig.MessageChannelSize
}

// GetDataRefreshInterval 实现 ConfigProvider 接口
func (c *Config) GetDataRefreshInterval() int {
	if c.DataRefreshConfig == nil {
		return 0
	}
	return c.DataRefreshConfig.IntervalSeconds
}

// GetWriteCountThreshold 实现 ConfigProvider 接口
func (c *Config) GetWriteCountThreshold() int {
	if c.DataRefreshConfig == nil {
		return 0
	}
	return c.DataRefreshConfig.WriteCountThreshold
}

// GetTR181DataModelPath 实现 ConfigProvider 接口
func (c *Config) GetTR181DataModelPath() string {
	if c.DataRefreshConfig == nil {
		return ""
	}
	return c.DataRefreshConfig.TR181DataModelPath
}

// 确保 Config 实现了 ConfigProvider 接口
var _ interface {
	GetServerURL() string
	GetEndpointID() string
	GetControllerEndpointID() string
	GetMessageChannelSize() int
	GetDataRefreshInterval() int
	GetWriteCountThreshold() int
	GetTR181DataModelPath() string
} = (*Config)(nil)
